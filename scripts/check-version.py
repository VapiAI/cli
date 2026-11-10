"""Require one version for Cargo, npm, and the release tag."""
import json
import pathlib
import sys
import tomllib

root = pathlib.Path(__file__).resolve().parent.parent
cargo = tomllib.loads((root / 'Cargo.toml').read_text())['package']['version']
npm = json.loads((root / 'npm-wrapper/package.json').read_text())['version']
expected = sys.argv[1] if len(sys.argv) > 1 else cargo
if cargo != npm or cargo != expected:
    sys.exit(f'Version mismatch: Cargo={cargo}, npm={npm}, expected={expected}')
print(f'CLI version: {cargo}')
