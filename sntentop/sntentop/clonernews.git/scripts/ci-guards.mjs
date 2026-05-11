/*
 * Purpose: Fail CI early when source code violates non-negotiable AGENTS legacy constraints.
 * Public API: Command-line script executed by `npm run ci:guards`.
 */

import { readdirSync, readFileSync } from 'node:fs';
import path from 'node:path';
import process from 'node:process';

const ROOT_DIRECTORIES = Object.freeze(['src', 'tests']);
const JAVASCRIPT_EXTENSIONS = new Set(['.js', '.mjs', '.cjs']);

const FORBIDDEN_RULES = Object.freeze([
  {
    label: 'Legacy var declaration',
    pattern: /\bvar\b/g,
  },
  {
    label: 'Legacy Date API usage',
    pattern: /\bnew\s+Date\b|\bDate\.(now|parse)\b/g,
  },
  {
    label: 'CommonJS require call',
    pattern: /\brequire\s*\(/g,
  },
]);

const collectJavaScriptFiles = (directoryPath) => {
  const entries = readdirSync(directoryPath, { withFileTypes: true });
  const files = [];

  for (const entry of entries) {
    const fullPath = path.join(directoryPath, entry.name);

    if (entry.isDirectory()) {
      files.push(...collectJavaScriptFiles(fullPath));
      continue;
    }

    if (JAVASCRIPT_EXTENSIONS.has(path.extname(entry.name))) {
      files.push(fullPath);
    }
  }

  return files;
};

const getViolationsForFile = (filePath) => {
  const content = readFileSync(filePath, 'utf8');
  const violations = [];

  for (const rule of FORBIDDEN_RULES) {
    const hasViolation = rule.pattern.test(content);
    rule.pattern.lastIndex = 0;

    if (hasViolation) {
      violations.push(rule.label);
    }
  }

  return violations;
};

const findAllViolations = () => {
  const violations = [];

  for (const rootDirectory of ROOT_DIRECTORIES) {
    const files = collectJavaScriptFiles(rootDirectory);

    for (const filePath of files) {
      const fileViolations = getViolationsForFile(filePath);

      if (fileViolations.length > 0) {
        violations.push({
          filePath,
          rules: fileViolations,
        });
      }
    }
  }

  return violations;
};

const violations = findAllViolations();

if (violations.length === 0) {
  console.log('CI guideline guards passed.');
  process.exit(0);
}

console.error('CI guideline guard violations detected:');
for (const violation of violations) {
  const relativePath = path.relative(process.cwd(), violation.filePath);
  console.error(`- ${relativePath}: ${violation.rules.join(', ')}`);
}

process.exit(1);
