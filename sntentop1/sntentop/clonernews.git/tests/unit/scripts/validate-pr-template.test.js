// @vitest-environment node

import { spawnSync } from 'node:child_process';
import { mkdirSync, mkdtempSync, rmSync, writeFileSync } from 'node:fs';
import { tmpdir } from 'node:os';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
// These tests execute the policy script in temporary git repos so branch-specific behavior is checked end to end.
// Public API under test: scripts/validate-pr-template.mjs command-line behavior.
// Constraints: keep the repos tiny and deterministic so the script sees real git metadata without depending on the workspace state.
import { describe, expect, it } from 'vitest';

const validateScriptPath = fileURLToPath(
  new URL('../../../scripts/validate-pr-template.mjs', import.meta.url),
);

const createValidPrBody = () => `# PR Gate Checklist

## Required checks

- [x] I read AGENTS.md and the agentic workflow guide.
- [x] I ran npm run ci:quality locally.
- [x] I ran npm run ci:policy locally.
- [x] I verified my branch commits reference at least one ticket ID from docs/tickets.md.
- [x] I confirmed changed files stay within the declared ticket track ownership scope.
- [x] I ran the applicable local checks for this change.
- [x] I listed the audit IDs affected by this change.
- [x] I checked security sinks and trust boundaries.
- [x] I checked architecture boundaries.
- [x] I checked dependency and lockfile impact.
- [x] I requested human review.

## Layer boundary confirmation

- [x] \`src/core/\` has no DOM or fetch references.
- [x] \`src/infra/\` has no DOM references.
- [x] All HN API HTML content is sanitized via DOMPurify before insertion.
- [x] All fetch calls use AbortController with timeout.
- [x] Live-data polling respects 5-second throttle minimum.

## What changed
- Example change.

## Why
- Example reason.

## Tests
- Example test.

## Audit questions affected
- Example audit ID.

## Security notes
- Example note.

## Architecture / dependency notes
- Example note.

## Risks
- Example risk.
`;

const createIsolatedEnv = (overrides = {}) => {
  const env = { ...process.env, ...overrides };

  delete env.GITHUB_EVENT_PATH;
  delete env.GITEA_EVENT_PATH;
  delete env.EVENT_PATH;
  delete env.GITHUB_BASE_REF;
  delete env.GITEA_BASE_REF;
  delete env.BASE_REF;

  return env;
};

const runGit = (cwd, args) => {
  const result = spawnSync('git', args, {
    cwd,
    encoding: 'utf8',
    stdio: ['ignore', 'pipe', 'pipe'],
  });

  if (result.status !== 0) {
    throw new Error(
      `git ${args.join(' ')} failed in ${cwd}\n${result.stderr || result.stdout || 'Unknown git error'}`,
    );
  }

  return result;
};

const createRepo = ({ branchName, files, commitMessage }) => {
  const repoDir = mkdtempSync(path.join(tmpdir(), 'clonernews-policy-'));

  try {
    runGit(repoDir, ['init']);
    runGit(repoDir, ['checkout', '-b', branchName]);
    runGit(repoDir, ['config', 'user.name', 'Policy Test']);
    runGit(repoDir, ['config', 'user.email', 'policy@test.invalid']);

    for (const [relativePath, content] of Object.entries(files)) {
      const absolutePath = path.join(repoDir, relativePath);
      mkdirSync(path.dirname(absolutePath), { recursive: true });
      writeFileSync(absolutePath, content);
    }

    runGit(repoDir, ['add', '.']);
    runGit(repoDir, ['commit', '-m', commitMessage]);

    return repoDir;
  } catch (error) {
    rmSync(repoDir, { recursive: true, force: true });
    throw error;
  }
};

const runPolicyScript = (cwd, args = []) =>
  spawnSync(process.execPath, [validateScriptPath, ...args], {
    cwd,
    encoding: 'utf8',
    env: createIsolatedEnv(),
    stdio: ['ignore', 'pipe', 'pipe'],
  });

const expectSuccess = (result) => {
  expect(result.status).toBe(0);
  expect(`${result.stdout}${result.stderr}`).not.toMatch(/Error:/i);
};

const expectFailure = (result, messagePattern) => {
  expect(result.status).not.toBe(0);
  expect(`${result.stdout}${result.stderr}`).toMatch(messagePattern);
};

describe('validate-pr-template policy script', () => {
  it('allows the default branch to skip PR-body sourcing while still validating shared scope', () => {
    const repoDir = createRepo({
      branchName: 'main',
      files: {
        'src/app.js': 'export const app = true;\n',
      },
      commitMessage: 'chore: baseline shared change',
    });

    try {
      const result = runPolicyScript(repoDir);
      expectSuccess(result);
      expect(result.stdout).toMatch(/default branch/i);
      expect(result.stdout).toMatch(/Shared branch scope checks passed/i);
    } finally {
      rmSync(repoDir, { recursive: true, force: true });
    }
  });

  it('still rejects feature branches that do not carry a ticket ID', () => {
    const repoDir = createRepo({
      branchName: 'feature/no-ticket',
      files: {
        'src/app.js': 'export const feature = true;\n',
        'pr.md': createValidPrBody(),
      },
      commitMessage: 'feat: branch without ticket',
    });

    try {
      const result = runPolicyScript(repoDir, ['--file', 'pr.md']);
      expectFailure(result, /No ticket ID found/i);
    } finally {
      rmSync(repoDir, { recursive: true, force: true });
    }
  });

  it('accepts a feature branch with ticket metadata in commits even when no PR body source exists', () => {
    const repoDir = createRepo({
      branchName: 'feature/no-ticket-in-branch',
      files: {
        'docs/tickets.md': '- [ ] **TA-7:** Inferred ticket test\n',
        'src/core/inferred-ticket.js': 'export const inferredTicket = true;\n',
      },
      commitMessage: 'feat(TA-7): infer ticket from commit metadata',
    });

    try {
      const result = runPolicyScript(repoDir);
      expectSuccess(result);
      expect(result.stdout).toMatch(/skipping PR template section\/checklist validation/i);
      expect(result.stdout).toMatch(/Ticket association and scope checks passed/i);
    } finally {
      rmSync(repoDir, { recursive: true, force: true });
    }
  });

  it('keeps the default branch scope bounded to the shared ticket track', () => {
    const repoDir = createRepo({
      branchName: 'main',
      files: {
        'tmp/rogue.txt': 'outside the shared scope\n',
      },
      commitMessage: 'chore: out-of-scope shared change',
    });

    try {
      const result = runPolicyScript(repoDir);
      expectFailure(result, /shared main-branch scope|Out-of-scope files/i);
    } finally {
      rmSync(repoDir, { recursive: true, force: true });
    }
  });
});
