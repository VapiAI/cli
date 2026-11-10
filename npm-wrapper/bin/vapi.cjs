#!/usr/bin/env node
const path = require('node:path');
const { spawnSync } = require('node:child_process');

const binary = path.join(__dirname, '..', 'native', process.platform === 'win32' ? 'vapi.exe' : 'vapi');
const result = spawnSync(binary, process.argv.slice(2), { stdio: 'inherit' });
if (result.error) {
  console.error(`Cannot run Vapi CLI: ${result.error.message}. Reinstall @vapi-ai/cli.`);
  process.exitCode = 1;
} else if (result.signal) {
  process.kill(process.pid, result.signal);
} else {
  process.exitCode = result.status ?? 1;
}
