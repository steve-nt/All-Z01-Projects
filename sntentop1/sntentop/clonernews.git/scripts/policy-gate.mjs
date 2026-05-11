/*
 * Purpose: Run policy checks that validate PR metadata, ticket-track ownership scope, and architecture boundaries.
 * Public API: Command-line script used by `npm run ci:policy`.
 * Notes: CLI args are forwarded to validate-pr-template.mjs (`--file` and `--base` are supported).
 */

import { spawnSync } from 'node:child_process';
import path from 'node:path';
import process from 'node:process';
import { fileURLToPath } from 'node:url';

const scriptDirectory = path.dirname(fileURLToPath(import.meta.url));
const forwardedArguments = process.argv.slice(2);

const runStep = (scriptName, args = []) => {
  const scriptPath = path.join(scriptDirectory, scriptName);
  const result = spawnSync(process.execPath, [scriptPath, ...args], {
    stdio: 'inherit',
    env: process.env,
  });

  if (typeof result.status === 'number') {
    if (result.status !== 0) {
      process.exit(result.status);
    }

    return;
  }

  process.exit(1);
};

runStep('validate-pr-template.mjs', forwardedArguments);
runStep('enforce-architecture-boundaries.mjs');

console.log('Policy gate checks passed.');
