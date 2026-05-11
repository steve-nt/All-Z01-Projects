/*
 * Purpose: Enforce layer boundaries by scanning core and infra files for forbidden DOM/fetch usage.
 * Public API: Command-line script used by `npm run ci:policy` and CI policy gate workflows.
 */

import { existsSync, readdirSync, readFileSync } from 'node:fs';
import path from 'node:path';
import process from 'node:process';

const SOURCE_EXTENSIONS = new Set(['.js', '.mjs', '.cjs', '.ts', '.tsx', '.jsx']);

const DOM_PATTERNS = Object.freeze([
  /\bdocument\./,
  /\bwindow\./,
  /\bquerySelector(All)?\s*\(/,
  /\bcreateElement(NS)?\s*\(/,
  /\bappendChild\s*\(/,
  /\binsertBefore\s*\(/,
  /\baddEventListener\s*\(/,
  /\binnerHTML\b/,
  /\bouterHTML\b/,
  /\binsertAdjacentHTML\b/,
]);

const FETCH_PATTERN = /\bfetch\s*\(/;

const walkFiles = (rootDirectory) => {
  const files = [];
  const stack = [rootDirectory];

  while (stack.length > 0) {
    const currentDirectory = stack.pop();

    for (const entry of readdirSync(currentDirectory, { withFileTypes: true })) {
      const fullPath = path.join(currentDirectory, entry.name);

      if (entry.isDirectory()) {
        stack.push(fullPath);
        continue;
      }

      if (SOURCE_EXTENSIONS.has(path.extname(entry.name))) {
        files.push(fullPath);
      }
    }
  }

  return files;
};

const scanFiles = (files, patterns, scope) => {
  const violations = [];

  for (const filePath of files) {
    const content = readFileSync(filePath, 'utf8');

    for (const pattern of patterns) {
      if (pattern.test(content)) {
        violations.push(`${scope} violation in ${filePath}: ${pattern}`);
      }
    }
  }

  return violations;
};

const hasCoreDirectory = existsSync('src/core');
const hasInfraDirectory = existsSync('src/infra');

if (!hasCoreDirectory && !hasInfraDirectory) {
  console.log('No src/core or src/infra found; skipping architecture boundary scan.');
  process.exit(0);
}

const violations = [];

if (hasCoreDirectory) {
  const coreFiles = walkFiles('src/core');
  violations.push(...scanFiles(coreFiles, [...DOM_PATTERNS, FETCH_PATTERN], 'src/core'));
}

if (hasInfraDirectory) {
  const infraFiles = walkFiles('src/infra');
  violations.push(...scanFiles(infraFiles, DOM_PATTERNS, 'src/infra'));
}

if (violations.length > 0) {
  throw new Error(`Architecture boundary violations found:\n- ${violations.join('\n- ')}`);
}

console.log('Architecture boundary checks passed.');
