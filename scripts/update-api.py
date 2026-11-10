"""Vendor the shared SDK inputs from an immutable VapiAI/docs commit."""
import hashlib
import json
from pathlib import Path
import re
import sys
from urllib.request import urlopen

if len(sys.argv) != 2 or not re.fullmatch(r'[0-9a-f]{40}', sys.argv[1]):
    sys.exit('Usage: python3 scripts/update-api.py <40-character-docs-commit-sha>')
root = Path(__file__).resolve().parent.parent / 'fern'
source = json.loads((root / 'source.json').read_text())
revision = sys.argv[1]
files = {}
for name in source['sha256']:
    url = f'https://raw.githubusercontent.com/{source["repository"]}/{revision}/{source["directory"]}/{name}'
    with urlopen(url, timeout=60) as response:
        files[name] = response.read()
source['revision'] = revision
source['sha256'] = {name: hashlib.sha256(content).hexdigest() for name, content in files.items()}
for name, content in files.items():
    (root / name).write_bytes(content)
(root / 'source.json').write_text(json.dumps(source, indent=2) + '\n')
print(f'Updated SDK inputs to {revision}. Run make generate and make test.')
