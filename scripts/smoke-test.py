"""Exercise the shipped CLI against a local HTTP server, never the Vapi API."""
import http.server
import json
import os
from pathlib import Path
import subprocess
import tempfile
import threading
import unittest
from urllib.parse import unquote

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
        self.environment.update(NO_PROXY='127.0.0.1', FERN_CLI_CREDENTIAL_STORE='file',
                                XDG_CONFIG_HOME=self.directory.name, APPDATA=self.directory.name)
        self.server = http.server.ThreadingHTTPServer(('127.0.0.1', 0), Handler)
        self.server.requests = []
        self.server.status = 200
        self.server.response = {'id': 'local-assistant', 'name': 'Preview'}
        thread = threading.Thread(target=self.server.serve_forever, daemon=True)
        thread.start()
        self.addCleanup(self.server.server_close)
        self.addCleanup(self.server.shutdown)
        self.base = f'http://127.0.0.1:{self.server.server_port}'
        self.cli('auth', 'login', '--with-token', stdin='local-test-token\n')

    def cli(self, *arguments, stdin=None, success=True, timeout=20):
        result = subprocess.run([str(BINARY), '--base-url', self.base, *arguments],
                                input=stdin, text=True, capture_output=True, timeout=timeout,
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

    def schema_target(self, document, schema):
        seen = set()
        while isinstance(schema, dict) and '$ref' in schema:
            reference = schema['$ref']
            self.assertTrue(reference.startswith('#/'), f'Nonlocal schema reference: {reference}')
            self.assertNotIn(reference, seen, f'Schema alias cycle: {reference}')
            seen.add(reference)
            schema = document
            for token in unquote(reference[2:]).split('/'):
                token = token.replace('~1', '/').replace('~0', '~')
                schema = schema[int(token)] if isinstance(schema, list) else schema[token]
        return schema

    def assert_schema_references_resolve(self, document, node):
        if isinstance(node, dict):
            if '$ref' in node:
                self.schema_target(document, node)
            for key, value in node.items():
                # Literal example/default values can contain arbitrary JSON.
                if key not in {'example', 'examples', 'default', 'serverDefault', 'enum'}:
                    self.assert_schema_references_resolve(document, value)
        elif isinstance(node, list):
            for value in node:
                self.assert_schema_references_resolve(document, value)

    def test_scoped_schema_and_json_help_without_auth(self):
        self.cli('auth', 'logout')
        for operation, method in [('list', 'GET'), ('create', 'POST')]:
            with self.subTest(operation=operation):
                schema = json.loads(self.cli('assistants', operation, '--schema', timeout=5).stdout)
                json_help = json.loads(self.cli('assistants', operation, '--help', '--format', 'json', timeout=5).stdout)
                self.assertEqual(json_help, schema)
                self.assertEqual(schema['operation'], f'assistants.{operation}')
                self.assertEqual((schema['httpMethod'], schema['path']), (method, '/assistant'))
                response = self.schema_target(schema, schema['output'])
                if operation == 'list':
                    self.assertEqual(response['type'], 'array')
                    response = self.schema_target(schema, response['items'])
                    limit = schema['input']['properties']['limit']
                    self.assertEqual((limit['flag'], limit['location']), ('--limit', 'query'))
                else:
                    name = schema['input']['properties']['name']
                    self.assertEqual((name['type'], name['flag'], name['location']),
                                     ('string', '--name', 'body'))
                self.assertEqual(response['type'], 'object')
                self.assertEqual(response['properties']['id']['type'], 'string')
                self.assert_schema_references_resolve(schema, schema)
        self.assertEqual(self.server.requests, [])

    def test_get_auth_path_and_json_output(self):
        result = self.cli('assistants', 'get', '--id', 'local-assistant')
        self.assertEqual(json.loads(result.stdout), self.server.response)
        method, path, headers, _ = self.server.requests[0]
        self.assertEqual((method, path), ('GET', '/assistant/local-assistant'))
        self.assertEqual(headers['authorization'], 'Bearer local-test-token')

    def test_logout_removes_access_and_legacy_env_does_not_restore_it(self):
        status = self.cli('auth', 'status')
        self.assertIn('bearer', status.stdout)
        self.assertNotIn('local-test-token', status.stdout + status.stderr)
        self.cli('auth', 'logout')
        self.environment.update(VAPI_TOKEN='ignored-token', VAPI_API_KEY='ignored-api-key')
        result = self.cli('assistants', 'list', success=False)
        self.assertIn('Authentication credentials are missing', result.stdout + result.stderr)
        self.assertEqual(self.server.requests, [])

    def test_create_body_from_stdin(self):
        self.cli('assistants', 'create', '--json', '-', stdin='{"name":"From stdin"}')
        method, path, _, body = self.server.requests[0]
        self.assertEqual((method, path), ('POST', '/assistant'))
        self.assertEqual(json.loads(body), {'name': 'From stdin'})

    def test_create_assistant_then_call_preserves_nested_json(self):
        assistant = {
            'name': 'Quickstart',
            'model': {'provider': 'openai', 'model': 'gpt-4.1-mini',
                      'messages': [{'role': 'system', 'content': 'Keep your answers brief.'}]},
            'voice': {'provider': 'vapi', 'voiceId': 'Clara'},
        }
        created = json.loads(self.cli('assistants', 'create', '--json', json.dumps(assistant)).stdout)
        call = {'assistantId': created['id'], 'phoneNumberId': 'local-phone',
                'customer': {'number': '+14155550123'}}
        self.server.response = {'id': 'local-call', 'status': 'queued'}
        result = self.cli('calls', 'create', '--json', json.dumps(call))
        self.assertEqual(json.loads(result.stdout), self.server.response)
        self.assertEqual(len(self.server.requests), 2)
        for request, path, expected in zip(self.server.requests, ['/assistant', '/call'], [assistant, call]):
            method, actual_path, headers, body = request
            self.assertEqual((method, actual_path), ('POST', path))
            self.assertEqual(headers['authorization'], 'Bearer local-test-token')
            self.assertEqual(json.loads(body), expected)

    def test_invalid_input_fails_before_sending_requests(self):
        missing_id = self.cli('assistants', 'get', success=False)
        self.assertIn('--id', missing_id.stdout + missing_id.stderr)
        malformed = self.cli('assistants', 'create', '--json', '{', success=False)
        self.assertIn('Invalid --json body', malformed.stdout + malformed.stderr)
        self.assertEqual(self.server.requests, [])

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
