"""Generate in isolation, then replace only generator-owned files."""
import argparse
import filecmp
import pathlib
import shutil
import subprocess
import tempfile
import tomllib

root = pathlib.Path(__file__).resolve().parent.parent
parser = argparse.ArgumentParser(description=__doc__)
parser.add_argument('--check', action='store_true', help='Fail if regeneration differs; leave the checkout untouched')
args = parser.parse_args()
version = tomllib.loads((root / 'Cargo.toml').read_text())['package']['version']
owned = ['Cargo.toml', 'Cargo.lock', 'reference.md', 'src', 'cli', 'tests', 'wiremock']


def identical(left, right):
    if left.is_dir() and right.is_dir():
        entries = {p.name for p in left.iterdir()} | {p.name for p in right.iterdir()}
        return all(identical(left / name, right / name) for name in entries)
    return left.is_file() and right.is_file() and filecmp.cmp(left, right, shallow=False)


with tempfile.TemporaryDirectory(prefix='vapi-fern-') as temporary:
    work = pathlib.Path(temporary)
    shutil.copytree(root / 'fern', work / 'fern')
    subprocess.run(['fern', 'generate', '--group', 'cli', '--version', version], cwd=work, check=True)
    generated = work / 'generated'
    for name in owned:
        if not (generated / name).exists():
            raise SystemExit(f'Fern did not generate expected output: {name}')
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
                shutil.copyfile(source, destination)
        print('Updated generated CLI. Review git diff before committing.')
