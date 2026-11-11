"""Generate in isolation, then replace only generator-owned files."""
import argparse
import filecmp
import hashlib
import json
import os
import pathlib
import shutil
import subprocess
import tempfile
import tomllib

root = pathlib.Path(__file__).resolve().parent.parent
parser = argparse.ArgumentParser(description=__doc__)
parser.add_argument('--check', action='store_true', help='Fail if regeneration differs; leave the checkout untouched')
parser.add_argument('--publishing', action='store_true', help='Preview the publishing group and include its generated workflows')
args = parser.parse_args()
version = tomllib.loads((root / 'Cargo.toml').read_text())['package']['version']
owned = ['Cargo.toml', 'Cargo.lock', 'reference.md', 'src', 'cli', 'tests', 'wiremock']
if args.publishing:
    owned += ['dist-workspace.toml', '.github/workflows/ci.yml', '.github/workflows/release.yml']
provenance = json.loads((root / 'fern/source.json').read_text())
for name, expected in provenance['sha256'].items():
    actual = hashlib.sha256((root / 'fern' / name).read_bytes()).hexdigest()
    if actual != expected:
        raise SystemExit(f'Vendored API input differs from fern/source.json: {name}')


def identical(left, right):
    if left.is_dir() and right.is_dir():
        entries = {p.name for p in left.iterdir()} | {p.name for p in right.iterdir()}
        return all(identical(left / name, right / name) for name in entries)
    return left.is_file() and right.is_file() and filecmp.cmp(left, right, shallow=False)


with tempfile.TemporaryDirectory(prefix='vapi-fern-') as temporary:
    work = pathlib.Path(temporary)
    shutil.copytree(root / 'fern', work / 'fern')
    # asdf resolves shims from the child's working directory. Carry the
    # repository's Node pin into the temporary workspace, too.
    environment = dict(os.environ)
    environment.setdefault('ASDF_NODEJS_VERSION', '22.23.4')
    command = ['fern', 'generate', '--group', 'publish' if args.publishing else 'cli', '--version', version]
    if args.publishing:
        command += ['--preview', '--output', str(work / 'generated')]
    subprocess.run(command,
                   cwd=work, env=environment, check=True)
    generated = work / 'generated' / 'fern-cli-generator' if args.publishing else work / 'generated'
    for name in owned:
        if not (generated / name).exists():
            raise SystemExit(f'Fern did not generate expected output: {name}')
    # Keep temporary upstream fixes reproducible for both output groups. A
    # changed generator must match each patch before anything is adopted.
    for patch in sorted((root / 'patches').glob('*.patch')):
        subprocess.run(['git', 'apply', '--check', str(patch)], cwd=generated, check=True)
        subprocess.run(['git', 'apply', str(patch)], cwd=generated, check=True)
        print(f'Applied temporary Fern fix: {patch.name}', flush=True)
    changed = [name for name in owned if not identical(generated / name, root / name)]
    if args.check:
        if changed:
            raise SystemExit('Generated output differs: ' + ', '.join(changed))
        print('Generated files match the pinned Fern inputs.')
    else:
        for name in owned:
            source, destination = generated / name, root / name
            if source.is_dir():
                if destination.exists():
                    shutil.rmtree(destination)
                shutil.copytree(source, destination)
            else:
                destination.parent.mkdir(parents=True, exist_ok=True)
                shutil.copyfile(source, destination)
        print('Updated generated CLI. Review git diff before committing.')
