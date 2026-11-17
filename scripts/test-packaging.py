"""Pack the generated npm artifacts locally; all registry operations are stubbed."""
import argparse
import json
import os
from pathlib import Path
import re
import shutil
import subprocess
import sys
import tempfile
import unittest

ROOT = Path(__file__).resolve().parent.parent
parser = argparse.ArgumentParser(add_help=False)
parser.add_argument('--workflow', type=Path, default=ROOT / '.github/workflows/ci.yml')
options, remaining = parser.parse_known_args()


def step_script(name):
    text = options.workflow.read_text()
    start = text.index('      - name: ' + name + '\n')
    match = re.search(r'^        run: \|\n((?:          .*\n|\n)+)', text[start:], re.M)
    if not match:
        raise AssertionError('No shell body for ' + name)
    return '\n'.join(line[10:] if line.startswith('          ') else line for line in match[1].splitlines())


class NpmPackagingTest(unittest.TestCase):
    def test_platforms_and_launcher_include_license_notices(self):
        npm = shutil.which('npm')
        self.assertIsNotNone(npm, 'npm is required for the offline packaging check')
        cases = [('linux-x64', 'x86_64-unknown-linux-musl'),
                 ('linux-arm64', 'aarch64-unknown-linux-musl'),
                 ('darwin-x64', 'x86_64-apple-darwin'),
                 ('darwin-arm64', 'aarch64-apple-darwin'),
                 ('win32-x64', 'x86_64-pc-windows-msvc'), ('launcher', None)]
        for suffix, target in cases:
            with self.subTest(package=suffix), tempfile.TemporaryDirectory(prefix='vapi-npm-pack-') as directory:
                root = Path(directory)
                for name in ['LICENSE', 'LICENSE-APACHE', 'NOTICE']:
                    shutil.copyfile(ROOT / name, root / name)
                shim = root / 'shim'
                shim.mkdir()
                npm_shim = shim / 'npm'
                npm_shim.write_text('''#!/bin/bash
set -e
case "$1" in
  view) echo 'npm error code E404' >&2; exit 1 ;;
  publish) "$VAPI_PACK_REAL_NPM" pack --offline --ignore-scripts --json > "$VAPI_PACK_REPORT" ;;
  dist-tag) exit 0 ;;
  *) echo "Unexpected registry command: $1" >&2; exit 99 ;;
esac
''')
                npm_shim.chmod(0o755)
                npx = shim / 'npx'
                npx.write_text('#!/bin/bash\necho "Unexpected network-dependent npx invocation" >&2\nexit 99\n')
                npx.chmod(0o755)
                report = root / 'pack-report.json'
                env = dict(os.environ, PATH=str(shim) + os.pathsep + os.environ['PATH'],
                           GITHUB_REF_NAME='v1.0.0-beta.1', NPM_REGISTRY_URL='https://registry.invalid',
                           VAPI_PACK_REAL_NPM=npm, VAPI_PACK_REPORT=str(report))
                if target:
                    binary = root / 'target' / target / 'dist' / ('vapi.exe' if suffix.startswith('win32') else 'vapi')
                    binary.parent.mkdir(parents=True)
                    binary.write_text('fixture binary\n')
                    script = step_script('Package and publish npm platform package')
                    script = script.replace('${{ matrix.rust-target }}', target)
                    script = script.replace('${{ matrix.npm-platform-suffix }}', suffix)
                    package = root / 'npm-pkg' / '@vapi-ai' / ('cli-' + suffix)
                else:
                    script = step_script('Publish launcher package')
                    package = root / 'npm-pkg/launcher'
                result = subprocess.run(['bash', '-c', script], cwd=root, env=env,
                                        capture_output=True, text=True, timeout=60)
                self.assertEqual(result.returncode, 0, result.stdout + result.stderr)
                packed = json.loads(report.read_text())[0]
                paths = {entry['path'] for entry in packed['files']}
                self.assertTrue({'LICENSE', 'LICENSE-APACHE', 'NOTICE'} <= paths, paths)
                manifest = json.loads((package / 'package.json').read_text())
                self.assertEqual(manifest['license'], 'MIT AND Apache-2.0')


if __name__ == '__main__':
    unittest.main(argv=[sys.argv[0], *remaining])
