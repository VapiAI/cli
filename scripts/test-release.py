"""Exercise release refusal against generated-CLI fixtures without publishing."""
import importlib.util
import json
from pathlib import Path
import tempfile
import unittest

path = Path(__file__).with_name('check-release.py')
spec = importlib.util.spec_from_file_location('check_release', path)
guard = importlib.util.module_from_spec(spec)
spec.loader.exec_module(guard)


class ReleaseGuardTest(unittest.TestCase):
    def setUp(self):
        self.directory = tempfile.TemporaryDirectory()
        self.addCleanup(self.directory.cleanup)
        self.root = Path(self.directory.name)
        folder = self.root / 'cli/vapi'
        folder.mkdir(parents=True)
        self.source = '.client_id("public-production-client").authorization_url("https://auth.vapi.ai/authorize").token_url("https://auth.vapi.ai/token").authorization_params([("resource", "https://api.vapi.ai")])'
        (folder / 'main.rs').write_text(self.source)
        (folder / 'openapi0.json').write_text(json.dumps({'servers': [{'url': 'https://api.vapi.ai'}]}))

    def test_production_shaped_registration_passes(self):
        guard.validate(self.root)

    def test_staging_client_cannot_be_hidden_by_promoting_only_urls(self):
        (self.root / 'cli/vapi/main.rs').write_text(self.source.replace('public-production-client', guard.STAGING_CLIENT))
        with self.assertRaisesRegex(ValueError, 'staging OAuth client_id'):
            guard.validate(self.root)

    def test_each_staging_surface_blocks_release(self):
        for original, replacement in [('auth.vapi.ai', 'staging-auth.vapi.ai'), ('https://api.vapi.ai', 'https://api-mcp-poc.vapi.local')]:
            with self.subTest(replacement=replacement):
                (self.root / 'cli/vapi/main.rs').write_text(self.source.replace(original, replacement))
                with self.assertRaises(ValueError):
                    guard.validate(self.root)
        (self.root / 'cli/vapi/main.rs').write_text(self.source)
        (self.root / 'cli/vapi/openapi0.json').write_text(json.dumps({'servers': [{'url': 'https://staging-api.vapi.ai'}]}))
        with self.assertRaisesRegex(ValueError, 'default API server'):
            guard.validate(self.root)

    def test_unrecognized_generated_auth_fails_closed(self):
        (self.root / 'cli/vapi/main.rs').write_text('fn main() {}')
        with self.assertRaisesRegex(ValueError, 'Cannot verify'):
            guard.validate(self.root)


if __name__ == '__main__':
    unittest.main()
