#!/usr/bin/env node
const fs = require('node:fs');
const os = require('node:os');
const path = require('node:path');
const { createHash } = require('node:crypto');
const { execFileSync } = require('node:child_process');

const REPOSITORY = 'https://github.com/VapiAI/cli';
const PLATFORMS = {
  'darwin-x64': 'Darwin_x86_64',
  'darwin-arm64': 'Darwin_arm64',
  'linux-x64': 'Linux_x86_64',
  'linux-arm64': 'Linux_arm64',
  'win32-x64': 'Windows_x86_64',
};

function artifact(version, platform = process.platform, arch = process.arch) {
  if (!/^\d+\.\d+\.\d+(?:-[0-9A-Za-z.-]+)?$/.test(version)) {
    throw new Error(`Invalid CLI version: ${version}`);
  }
  const target = PLATFORMS[`${platform}-${arch}`];
  if (!target) throw new Error(`Unsupported platform: ${platform}-${arch}`);
  return {
    name: `cli_${target}${platform === 'win32' ? '.zip' : '.tar.gz'}`,
    binary: platform === 'win32' ? 'vapi.exe' : 'vapi',
    baseURL: `${REPOSITORY}/releases/download/v${version}`,
  };
}

async function download(url, fetchFile) {
  const response = await fetchFile(url);
  if (!response.ok) throw new Error(`Download failed (${response.status}): ${url}`);
  return Buffer.from(await response.arrayBuffer());
}

async function install({ directory = __dirname, fetchFile = fetch } = {}) {
  const { version } = JSON.parse(fs.readFileSync(path.join(directory, 'package.json'), 'utf8'));
  const release = artifact(version);
  const archive = await download(`${release.baseURL}/${release.name}`, fetchFile);
  const checksums = (await download(`${release.baseURL}/checksums.txt`, fetchFile)).toString('utf8');
  const expected = checksums.split(/\r?\n/).find(line => line.trim().split(/\s+/)[1] === release.name);
  const actual = createHash('sha256').update(archive).digest('hex');
  if (!expected || expected.split(/\s+/)[0] !== actual) throw new Error(`Checksum mismatch for ${release.name}`);
  const temporary = fs.mkdtempSync(path.join(os.tmpdir(), 'vapi-install-'));
  try {
    const archivePath = path.join(temporary, release.name);
    fs.writeFileSync(archivePath, archive);
    execFileSync('tar', ['-xf', archivePath, '-C', temporary], { stdio: 'inherit' });
    const destination = path.join(directory, 'native');
    fs.mkdirSync(destination, { recursive: true });
    const binary = path.join(destination, release.binary);
    fs.copyFileSync(path.join(temporary, release.binary), binary);
    fs.chmodSync(binary, 0o755);
    console.log(`Installed Vapi CLI ${version}. Run vapi --help to get started.`);
  } finally {
    fs.rmSync(temporary, { recursive: true, force: true });
  }
}

module.exports = { artifact, install };
if (require.main === module) {
  install().catch(error => {
    console.error(`Vapi CLI installation failed: ${error.message}`);
    process.exitCode = 1;
  });
}
