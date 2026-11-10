"""Exercise the shipped CLI against a local HTTP server, never the Vapi API."""
import http.server
import json
import os
from pathlib import Path
import subprocess
import tempfile
import threading
import unittest

ROOT = Path(__file__).resolve().parent.parent
BINARY = Path(os.environ.get('VAPI_TEST_BINARY', ROOT / 'target/debug' / ('vapi.exe' if os.name == 'nt' else 'vapi')))


class Handler(http.server.BaseHTTPRequestHandler):
    def respond(self):
        size = int(self.headers.get('Content-Length', 0))
        body = self.rfile.read(size)
        self.server.requests.append((self.command, self.path, dict(self.headers), body))
        self.send_response(self.server.status)
        self.send_header('Content-Type', 'application/json')
        self.end_headers()
        self.wfile.write(json.dumps(self.server.response).encode())

    do_GET = do_POST = do_PATCH = do_DELETE = respond

    def log_message(self, *_args):
        pass


class CLITest(unittest.TestCase):
    def setUp(self):
        self.directory = tempfile.TemporaryDirectory(prefix='vapi-smoke-')
        self.addCleanup(self.directory.cleanup)
        self.environment = {k: v for k, v in os.environ.items() if not k.startswith('VAPI_')}
        self.environment.update(VAPI_TOKEN='local-test-token', NO_PROXY='127.0.0.1',
                                XDG_CONFIG_HOME=self.directory.name)
        self.server = http.server.ThreadingHTTPServer(('127.0.0.1', 0), Handler)
        self.server.requests = []
        self.server.status = 200
        self.server.response = {'id': 'local-assistant', 'name': 'Preview'}
        thread = threading.Thread(target=self.server.serve_forever, daemon=True)
        thread.start()
        self.addCleanup(self.server.server_close)
        self.addCleanup(self.server.shutdown)
        self.base = f'http://127.0.0.1:{self.server.server_port}'

    def cli(self, *arguments, stdin=None, success=True):
        result = subprocess.run([str(BINARY), '--base-url', self.base, *arguments],
                                input=stdin, text=True, capture_output=True, timeout=20,
                                env=self.environment, cwd=self.directory.name)
        if success:
            self.assertEqual(result.returncode, 0, result.stdout + result.stderr)
        else:
            self.assertNotEqual(result.returncode, 0, result.stdout + result.stderr)
        return result

    def test_schema_and_removed_commands(self):
        schema = json.loads(self.cli('--schema').stdout)
        operations = {operation['operation'] for operation in schema['operations']}
        self.assertIn('assistants.create', operations)
        self.assertIn('squads.list', operations)
        for command in ['assistant', 'init', 'listen', 'workflow', 'mcp', 'update']:
            self.cli(command, success=False)
        self.assertEqual(self.server.requests, [])

    def test_get_auth_path_and_json_output(self):
        result = self.cli('assistants', 'get', '--id', 'local-assistant')
        self.assertEqual(json.loads(result.stdout), self.server.response)
        method, path, headers, _ = self.server.requests[0]
        self.assertEqual((method, path), ('GET', '/assistant/local-assistant'))
        self.assertEqual(headers['authorization'], 'Bearer local-test-token')

    def test_create_body_from_stdin(self):
        self.cli('assistants', 'create', '--json', '-', stdin='{"name":"From stdin"}')
        method, path, _, body = self.server.requests[0]
        self.assertEqual((method, path), ('POST', '/assistant'))
        self.assertEqual(json.loads(body), {'name': 'From stdin'})

    def test_query_and_response_projection(self):
        self.server.response = [{'id': 'one'}, {'id': 'two'}]
        result = self.cli('calls', 'list', '--limit', '2', '--query', '[].id')
        self.assertEqual(json.loads(result.stdout), ['one', 'two'])
        self.assertEqual(self.server.requests[0][1], '/call?limit=2')

    def test_delete_is_noninteractive_and_dry_run_sends_nothing(self):
        result = self.cli('assistants', 'delete', '--id', 'preview', '--dry-run')
        self.assertEqual(json.loads(result.stdout)['method'], 'DELETE')
        self.assertEqual(self.server.requests, [])
        self.cli('assistants', 'delete', '--id', 'preview', stdin='')
        self.assertEqual(self.server.requests[0][:2], ('DELETE', '/assistant/preview'))

    def test_api_errors_exit_nonzero(self):
        self.server.status = 401
        self.server.response = {'message': 'local unauthorized fixture'}
        result = self.cli('assistants', 'list', '--no-retry', success=False)
        self.assertIn('401', result.stdout + result.stderr)


if __name__ == '__main__':
    unittest.main()
