const { test } = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const os = require('node:os');
const path = require('node:path');
const { createHash } = require('node:crypto');
const { execFileSync } = require('node:child_process');
const { artifact, install } = require('./install');

test('maps every supported platform and pins prerelease downloads', () => {
  for (const [platform, arch, target, extension] of [
    ['darwin', 'x64', 'Darwin_x86_64', '.tar.gz'],
    ['darwin', 'arm64', 'Darwin_arm64', '.tar.gz'],
    ['linux', 'x64', 'Linux_x86_64', '.tar.gz'],
    ['linux', 'arm64', 'Linux_arm64', '.tar.gz'],
    ['win32', 'x64', 'Windows_x86_64', '.zip'],
  ]) {
    const result = artifact('1.0.0-beta.1', platform, arch);
    assert.equal(result.name, `cli_${target}${extension}`);
    assert.equal(result.baseURL, 'https://github.com/VapiAI/cli/releases/download/v1.0.0-beta.1');
  }
  assert.throws(() => artifact('1.0.0', 'linux', 'arm'), /Unsupported/);
  assert.throws(() => artifact('../latest'), /Invalid/);
});

test('installs the verified archive for the npm package version', async t => {
  const directory = fs.mkdtempSync(path.join(os.tmpdir(), 'vapi-package-test-'));
  t.after(() => fs.rmSync(directory, { recursive: true, force: true }));
  fs.writeFileSync(path.join(directory, 'package.json'), JSON.stringify({ version: '1.0.0-beta.1' }));
  const release = artifact('1.0.0-beta.1');
  const fixture = path.join(directory, 'fixture');
  fs.mkdirSync(fixture);
  fs.writeFileSync(path.join(fixture, release.binary), 'fixture binary');
  const archivePath = path.join(directory, release.name);
  const flags = process.platform === 'win32' ? ['-a', '-cf'] : ['-czf'];
  execFileSync('tar', [...flags, archivePath, '-C', fixture, release.binary]);
  const archive = fs.readFileSync(archivePath);
  const checksum = createHash('sha256').update(archive).digest('hex');
  const urls = [];
  const fetchFile = async url => {
    urls.push(url);
    return new Response(url.endsWith('/checksums.txt') ? `${checksum}  ${release.name}\n` : archive);
  };
  await install({ directory, fetchFile });
  assert.equal(fs.readFileSync(path.join(directory, 'native', release.binary), 'utf8'), 'fixture binary');
  assert.deepEqual(urls, [`${release.baseURL}/${release.name}`, `${release.baseURL}/checksums.txt`]);
});

test('rejects corrupt archives before installation', async t => {
  const directory = fs.mkdtempSync(path.join(os.tmpdir(), 'vapi-package-test-'));
  t.after(() => fs.rmSync(directory, { recursive: true, force: true }));
  fs.writeFileSync(path.join(directory, 'package.json'), JSON.stringify({ version: '1.0.0' }));
  await assert.rejects(install({ directory, fetchFile: async () => new Response('corrupt') }), /Checksum mismatch/);
  assert.equal(fs.existsSync(path.join(directory, 'native')), false);
});

test('reports missing releases instead of extracting an HTTP error', async () => {
  await assert.rejects(install({ fetchFile: async () => new Response('not found', { status: 404 }) }), /Download failed \(404\)/);
});
