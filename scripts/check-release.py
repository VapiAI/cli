"""Refuse public release tags while the generated CLI still uses staging OAuth."""
import argparse
import json
from pathlib import Path
import re
from urllib.parse import urlsplit

STAGING_CLIENT = 'client_01M721J47S27BS1GSXK1780TSR'


def validate(root):
    source = (root / 'cli/vapi/main.rs').read_text()
    for field in ('client_id', 'authorization_url', 'token_url'):
        matches = re.findall(r'\.' + field + r'\("([^\"]+)"\)', source)
        if not matches:
            raise ValueError(f'Cannot verify generated OAuth {field}')
        for value in matches:
            if value == STAGING_CLIENT or 'staging' in value.lower():
                raise ValueError(f'Public release blocked: staging OAuth {field}')
    resources = re.findall(r'\("resource",\s*"([^\"]+)"\)', source)
    if not resources:
        raise ValueError('Cannot verify the generated OAuth API resource')
    for value in resources:
        host = urlsplit(value).hostname or ''
        if host.endswith(('.local', '.test')) or 'staging' in host:
            raise ValueError('Public release blocked: staging OAuth API resource')
    specs = sorted((root / 'cli/vapi').glob('openapi*.json'))
    if not specs:
        raise ValueError('Cannot verify the generated API server')
    for spec in specs:
        server = json.loads(spec.read_text()).get('servers', [{}])[0].get('url', '')
        host = urlsplit(server).hostname or ''
        if not host or host.endswith(('.local', '.test')) or 'staging' in host:
            raise ValueError('Public release blocked: staging or missing default API server')


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--root', type=Path, default=Path(__file__).resolve().parent.parent)
    args = parser.parse_args()
    try:
        validate(args.root)
    except (ValueError, OSError, json.JSONDecodeError) as error:
        raise SystemExit(str(error)) from error
    print('Generated OAuth and default API settings contain no staging registration.')
