/*
 * Purpose: Validate PR body content when a source is supplied, then enforce ticket association and scope.
 * Public API: Command-line script used by `npm run ci:policy` in CI and local pre-PR checks.
 * Notes: CI validates the PR body from event payload; local runs may skip PR-body validation and still rely on branch/commit metadata for ticket and scope checks.
 */

import { spawnSync } from 'node:child_process';
import { existsSync, readFileSync } from 'node:fs';
import path from 'node:path';
import process from 'node:process';

const REQUIRED_SECTIONS = Object.freeze([
  'What changed',
  'Why',
  'Tests',
  'Audit questions affected',
  'Security notes',
  'Architecture / dependency notes',
  'Risks',
]);

const REQUIRED_CHECKBOXES = Object.freeze([
  'I read AGENTS.md and the agentic workflow guide.',
  'I ran npm run ci:quality locally.',
  'I ran npm run ci:policy locally.',
  'I verified my branch commits reference at least one ticket ID from docs/tickets.md.',
  'I confirmed changed files stay within the declared ticket track ownership scope.',
  'I ran the applicable local checks for this change.',
  'I listed the audit IDs affected by this change.',
  'I checked security sinks and trust boundaries.',
  'I checked architecture boundaries.',
  'I checked dependency and lockfile impact.',
  'I requested human review.',
]);

const REQUIRED_LAYER_CHECKBOXES = Object.freeze([
  '`src/core/` has no DOM or fetch references.',
  '`src/infra/` has no DOM references.',
  'All HN API HTML content is sanitized via DOMPurify before insertion.',
  'All fetch calls use AbortController with timeout.',
  'Live-data polling respects 5-second throttle minimum.',
]);

const TICKET_ID_PATTERN = /\b(?:T0|TA|TB|TC)-\d+\b/gi;

const GENERIC_ALLOWED_PATHS = Object.freeze([
  '.gitea/workflows/',
  '.github/workflows/',
  '.github/pull_request_template.md',
  '.gitignore',
  'AGENTS.md',
  'README.md',
  'biome.json',
  'docs/',
  'index.html',
  'package.json',
  'package-lock.json',
  'playwright.config.js',
  'scripts/',
  'vite.config.js',
]);

const TRACK_OWNERSHIP_BY_TRACK = Object.freeze({
  T0: Object.freeze({
    owner: 'Shared',
    allowedPaths: Object.freeze([
      '.gitea/workflows/',
      '.github/workflows/',
      'index.html',
      'public/',
      'src/',
      'tests/',
    ]),
  }),
  TA: Object.freeze({
    owner: 'Dev 1',
    allowedPaths: Object.freeze([
      '.gitea/workflows/',
      '.github/workflows/',
      'src/core/',
      'src/infra/hn-api-adapter.js',
      'tests/',
      'vite.config.js',
    ]),
  }),
  TB: Object.freeze({
    owner: 'Dev 2',
    allowedPaths: Object.freeze([
      'index.html',
      'public/',
      'src/app.css',
      'src/features/feed/',
      'src/main.js',
      'src/shared/',
      'src/styles/',
    ]),
  }),
  TC: Object.freeze({
    owner: 'Dev 3',
    allowedPaths: Object.freeze([
      'src/features/comments/',
      'src/features/live-banner/',
      'src/features/polls/',
      'src/features/post-detail/',
      'src/infra/cache-adapter.js',
      'src/infra/throttle.js',
    ]),
  }),
});

const escapeRegex = (text) => text.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');

const DEFAULT_BRANCH_NAMES = Object.freeze(new Set(['main', 'master']));

const isDefaultBranch = (branchName = '') =>
  DEFAULT_BRANCH_NAMES.has(branchName.trim().toLowerCase());

const PR_MESSAGE_DIRECTORY = path.join('docs', 'pr-messages');

const normalizePath = (filePath) =>
  filePath.replaceAll('\\', '/').replace(/^\.\//, '').replace(/^\/+/, '').trim();

const parseOptionArgument = (argv, optionName) => {
  const assignmentPrefix = `${optionName}=`;
  for (let index = 2; index < argv.length; index += 1) {
    const token = argv[index];

    if (token === optionName) {
      const value = argv[index + 1];

      if (!value) {
        throw new Error(`Expected a value after ${optionName}.`);
      }

      return value;
    }

    if (token.startsWith(assignmentPrefix)) {
      const value = token.slice(assignmentPrefix.length);

      if (value.length === 0) {
        throw new Error(`Expected a non-empty value for ${optionName} option.`);
      }

      return value;
    }
  }

  return undefined;
};

const runGit = (gitArguments, { allowFailure = false } = {}) => {
  const result = spawnSync('git', gitArguments, {
    encoding: 'utf8',
    stdio: ['ignore', 'pipe', 'pipe'],
  });

  if (result.status === 0) {
    return result.stdout.trim();
  }

  if (allowFailure) {
    return undefined;
  }

  const gitErrorOutput = (result.stderr || result.stdout || 'Unknown git error').trim();
  throw new Error(`Git command failed: git ${gitArguments.join(' ')}\n${gitErrorOutput}`);
};

const hasRef = (reference) =>
  Boolean(
    runGit(['rev-parse', '--verify', '--quiet', `${reference}^{commit}`], { allowFailure: true }),
  );

const resolveBaseReference = (argv) => {
  const cliBaseReference = parseOptionArgument(argv, '--base');
  const candidateReferences = [
    cliBaseReference,
    process.env.GITHUB_BASE_REF,
    process.env.GITEA_BASE_REF,
    process.env.BASE_REF,
    'origin/main',
    'origin/master',
    'main',
    'master',
  ];

  for (const reference of candidateReferences) {
    if (!reference) {
      continue;
    }

    if (hasRef(reference)) {
      return reference;
    }
  }

  return undefined;
};

const resolveCommitRange = (argv) => {
  const baseReference = resolveBaseReference(argv);

  if (baseReference) {
    const mergeBase = runGit(['merge-base', 'HEAD', baseReference], { allowFailure: true });

    if (mergeBase) {
      return {
        baseReference,
        range: `${mergeBase}..HEAD`,
      };
    }
  }

  if (hasRef('HEAD~1')) {
    return {
      baseReference: 'HEAD~1',
      range: 'HEAD~1..HEAD',
    };
  }

  return {
    baseReference: 'HEAD',
    range: 'HEAD',
  };
};

const extractTicketIds = (text) => {
  const ids = text.match(TICKET_ID_PATTERN) ?? [];
  return [...new Set(ids.map((id) => id.toUpperCase()))].toSorted();
};

const extractTrackKey = (ticketId) => ticketId.split('-')[0];

const loadKnownTicketIds = () => {
  const ticketFilePath = 'docs/tickets.md';

  if (!existsSync(ticketFilePath)) {
    throw new Error('Unable to validate ticket IDs because docs/tickets.md was not found.');
  }

  const ticketFileContent = readFileSync(ticketFilePath, 'utf8');
  const listedTicketIds = ticketFileContent.match(/\*\*((?:T0|TA|TB|TC)-\d+):?\*\*/g) ?? [];
  const normalizedIds = listedTicketIds.map((token) =>
    token.replaceAll('*', '').replace(/:$/, '').toUpperCase(),
  );

  return new Set(normalizedIds);
};

const pathMatchesScope = (relativeFilePath, scopePath) =>
  scopePath.endsWith('/') ? relativeFilePath.startsWith(scopePath) : relativeFilePath === scopePath;

const loadChangedFiles = (range) => {
  const diffOutput = runGit(['diff', '--name-only', '--diff-filter=ACDMRTUXB', range], {
    allowFailure: true,
  });

  const rawPaths =
    diffOutput && diffOutput.length > 0
      ? diffOutput.split(/\r?\n/)
      : (runGit(['show', '--pretty=', '--name-only', 'HEAD'], { allowFailure: true }) ?? '').split(
          /\r?\n/,
        );

  return [...new Set(rawPaths.map(normalizePath).filter(Boolean))];
};

const resolveTrackOwnership = (ticketIds) => {
  const trackKeys = [...new Set(ticketIds.map(extractTrackKey))].toSorted();
  const unknownTracks = trackKeys.filter((trackKey) => !TRACK_OWNERSHIP_BY_TRACK[trackKey]);

  if (unknownTracks.length > 0) {
    throw new Error(
      `Unsupported ticket track(s) found: ${unknownTracks.join(', ')}. Update TRACK_OWNERSHIP_BY_TRACK in scripts/validate-pr-template.mjs.`,
    );
  }

  const featureTracks = trackKeys.filter((trackKey) => trackKey !== 'T0');

  if (featureTracks.length > 1) {
    throw new Error(
      [
        `Mixed feature tracks detected in branch context: ${featureTracks.join(', ')}.`,
        'A PR must stay focused on one feature track owner scope at a time (T0 may be combined as shared scope).',
      ].join('\n'),
    );
  }

  return trackKeys.map((trackKey) => ({
    track: trackKey,
    owner: TRACK_OWNERSHIP_BY_TRACK[trackKey].owner,
    allowedPaths: TRACK_OWNERSHIP_BY_TRACK[trackKey].allowedPaths,
  }));
};

const validateTicketAssociationAndScope = (argv, branchName) => {
  const { baseReference, range } = resolveCommitRange(argv);
  const commitMessages = runGit(['log', '--format=%B', range], { allowFailure: true }) ?? '';
  const sharedBranch = isDefaultBranch(branchName);
  const ticketIds = sharedBranch ? [] : extractTicketIds(`${branchName}\n${commitMessages}`);

  let ownershipScopes = [{ track: 'T0', ...TRACK_OWNERSHIP_BY_TRACK.T0 }];

  if (!sharedBranch) {
    if (ticketIds.length === 0) {
      throw new Error(
        [
          'No ticket ID found in current branch name or commit messages.',
          'Include at least one ID from docs/tickets.md (for example: TA-7) in the branch name or commit message.',
        ].join('\n'),
      );
    }

    const knownTicketIds = loadKnownTicketIds();
    const unknownTicketIds = ticketIds.filter((ticketId) => !knownTicketIds.has(ticketId));

    if (unknownTicketIds.length > 0) {
      throw new Error(
        `Unknown ticket IDs found in branch context: ${unknownTicketIds.join(', ')}. Update docs/tickets.md or fix commit metadata.`,
      );
    }

    ownershipScopes = resolveTrackOwnership(ticketIds);
  }

  const ownershipLabels = ownershipScopes.map(({ track, owner }) => `${track} (${owner})`);

  const allowedPaths = [
    ...new Set([
      ...GENERIC_ALLOWED_PATHS,
      ...ownershipScopes.flatMap(({ allowedPaths: trackAllowedPaths }) => trackAllowedPaths),
    ]),
  ];
  const changedFiles = loadChangedFiles(range);

  const outOfScopeFiles = changedFiles.filter(
    (relativeFilePath) =>
      !allowedPaths.some((scopePath) => pathMatchesScope(relativeFilePath, scopePath)),
  );

  if (outOfScopeFiles.length > 0) {
    throw new Error(
      [
        sharedBranch
          ? 'Changed files do not match the shared main-branch scope.'
          : `Changed files do not match the declared ticket scope for ${ticketIds.join(', ')}.`,
        `Track ownership scope: ${ownershipLabels.join(', ')}.`,
        `Base reference: ${baseReference}.`,
        `Out-of-scope files: ${outOfScopeFiles.join(', ')}`,
      ].join('\n'),
    );
  }

  console.log(
    sharedBranch
      ? `Shared branch scope checks passed (ownership: ${ownershipLabels.join(', ')}, base: ${baseReference}).`
      : `Ticket association and scope checks passed (tickets: ${ticketIds.join(', ')}, ownership: ${ownershipLabels.join(', ')}, base: ${baseReference}).`,
  );
};

const parseFileArgument = (argv) => {
  const filePath = parseOptionArgument(argv, '--file');
  return filePath ? normalizePath(filePath) : undefined;
};

const resolveEventPath = () =>
  process.env.GITHUB_EVENT_PATH || process.env.GITEA_EVENT_PATH || process.env.EVENT_PATH;

const getCurrentBranchName = () => {
  const result = spawnSync('git', ['rev-parse', '--abbrev-ref', 'HEAD'], {
    encoding: 'utf8',
  });

  if (typeof result.status === 'number' && result.status === 0) {
    return result.stdout.trim();
  }

  return '';
};

const extractTicketFromBranchName = (branchName) => {
  const match = branchName.match(/\b(T0-\d+|TA-\d+|TB-\d+|TC-\d+)\b/i);
  return match ? match[1].toUpperCase() : undefined;
};

const resolveAutoLocalFilePath = () => {
  const branchName = getCurrentBranchName();
  const ticketId = extractTicketFromBranchName(branchName);

  if (!ticketId) {
    return undefined;
  }

  const candidatePath = path.join(PR_MESSAGE_DIRECTORY, `${ticketId}.md`);
  return existsSync(candidatePath) ? candidatePath : undefined;
};

const loadBodyFromEvent = (eventPath) => {
  if (!existsSync(eventPath)) {
    throw new Error(`Event payload was not found at "${eventPath}".`);
  }

  const eventPayload = JSON.parse(readFileSync(eventPath, 'utf8'));
  return eventPayload.pull_request?.body ?? '';
};

const loadBodyFromFile = (filePath) => {
  if (!existsSync(filePath)) {
    throw new Error(`PR body file was not found at "${filePath}".`);
  }

  return readFileSync(filePath, 'utf8');
};

const resolveBodySource = (argv) => {
  const eventPath = resolveEventPath();

  if (eventPath) {
    return {
      body: loadBodyFromEvent(eventPath),
      sourceLabel: `event payload (${eventPath})`,
    };
  }

  const filePath = parseFileArgument(argv);

  if (filePath) {
    return {
      body: loadBodyFromFile(filePath),
      sourceLabel: `file (${path.relative(process.cwd(), filePath)})`,
    };
  }

  if (isDefaultBranch(branchName)) {
    return undefined;
  }

  const autoFilePath = resolveAutoLocalFilePath();

  if (autoFilePath) {
    return {
      body: loadBodyFromFile(autoFilePath),
      sourceLabel: `auto-file (${path.relative(process.cwd(), autoFilePath)})`,
    };
  }

  return undefined;
};

const findMissingSections = (body) =>
  REQUIRED_SECTIONS.filter((sectionTitle) => {
    const sectionPattern = new RegExp(String.raw`^##\s+${escapeRegex(sectionTitle)}\s*$`, 'im');
    return !sectionPattern.test(body);
  });

const findMissingCheckedItems = (body, labels) =>
  labels.filter((label) => {
    const checkedItemPattern = new RegExp(
      String.raw`^\s*- \[[xX]\]\s+${escapeRegex(label)}\s*$`,
      'im',
    );
    return !checkedItemPattern.test(body);
  });

const branchName = getCurrentBranchName();
const source = resolveBodySource(process.argv);

if (source) {
  const missingSections = findMissingSections(source.body);
  const missingChecks = findMissingCheckedItems(source.body, REQUIRED_CHECKBOXES);
  const missingLayerChecks = findMissingCheckedItems(source.body, REQUIRED_LAYER_CHECKBOXES);

  if (missingSections.length > 0 || missingChecks.length > 0 || missingLayerChecks.length > 0) {
    const reportLines = [];

    if (missingSections.length > 0) {
      reportLines.push(`Missing required PR sections: ${missingSections.join(', ')}`);
    }

    if (missingChecks.length > 0) {
      reportLines.push(`Missing required checklist items: ${missingChecks.join(', ')}`);
    }

    if (missingLayerChecks.length > 0) {
      reportLines.push(
        `Missing required layer-boundary checklist items: ${missingLayerChecks.join(', ')}`,
      );
    }

    throw new Error(reportLines.join('\n'));
  }

  console.log(`PR template validation passed (${source.sourceLabel}).`);
} else {
  console.log(
    isDefaultBranch(branchName)
      ? 'No PR body source detected on the default branch; skipping PR template section/checklist validation.'
      : 'No PR body source detected; skipping PR template section/checklist validation.',
  );
}

validateTicketAssociationAndScope(process.argv, branchName);
