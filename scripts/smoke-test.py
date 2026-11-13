"""Exercise the shipped CLI against a local HTTP server, never the Vapi API."""
import http.server
import json
import os
import queue
import re
from pathlib import Path
import subprocess
import tempfile
import threading
import tomllib
import unittest
from urllib.parse import parse_qs, unquote, urlsplit

ROOT = Path(__file__).resolve().parent.parent
BINARY = Path(os.environ.get('VAPI_TEST_BINARY', ROOT / 'target/debug' / ('vapi.exe' if os.name == 'nt' else 'vapi')))
VERSION = tomllib.loads((ROOT / 'Cargo.toml').read_text())['package']['version']


class Handler(http.server.BaseHTTPRequestHandler):
    def respond(self):
        size = int(self.headers.get('Content-Length', 0))
        body = self.rfile.read(size)
        self.server.requests.append((self.command, self.path, dict(self.headers), body))
        status = self.server.statuses.pop(0) if self.server.statuses else self.server.status
        self.send_response(status)
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
        self.environment.pop('CI', None)
        self.environment.update(NO_PROXY='127.0.0.1', FERN_CLI_CREDENTIAL_STORE='file',
                                XDG_CONFIG_HOME=self.directory.name, APPDATA=self.directory.name)
        self.server = http.server.ThreadingHTTPServer(('127.0.0.1', 0), Handler)
        self.server.requests = []
        self.server.status = 200
        self.server.statuses = []
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

    def assert_client_headers(self, headers, command, suffix=None, ci=False):
        value = headers.get('x-vapi-client')
        self.assertIsNotNone(value)
        self.assertLessEqual(len(value), 256)
        product, *parameters = [part.strip() for part in value.split(';')]
        self.assertEqual(product, f'cli/{VERSION}')
        fields = dict(part.split('=', 1) for part in parameters if '=' in part)
        self.assertEqual(set(fields), {'command'})
        self.assertEqual(fields['command'], command)
        self.assertEqual([part for part in parameters if '=' not in part], ['ci'] if ci else [])
        user_agent = f'vapi-cli/{VERSION}' + (f' {suffix}' if suffix else '')
        self.assertEqual(headers.get('user-agent'), user_agent)
        return value

    def test_oauth_login_uses_pkce_without_a_client_secret(self):
        spec = json.loads(self.cli('--spec-raw').stdout)
        flow = spec['components']['securitySchemes']['bearer']['flows']['authorizationCode']
        process = subprocess.Popen([str(BINARY), 'auth', 'login', '--no-browser'],
                                   stdin=subprocess.DEVNULL, stdout=subprocess.PIPE,
                                   stderr=subprocess.STDOUT, text=True,
                                   env=self.environment, cwd=self.directory.name)
        lines = queue.Queue()

        def read_output():
            for line in process.stdout:
                lines.put(line)
            lines.put(None)

        reader = threading.Thread(target=read_output, daemon=True)
        reader.start()
        authorization = None
        output = []
        try:
            while authorization is None:
                line = lines.get(timeout=5)
                if line is None:
                    break
                output.append(line)
                for value in re.findall(r'https?://[^\s]+', line):
                    parsed = urlsplit(value)
                    if 'code_challenge' in parse_qs(parsed.query):
                        authorization = parsed
                        break
        finally:
            if process.poll() is None:
                process.terminate()
            process.wait(timeout=5)
            reader.join(timeout=5)
            process.stdout.close()
        self.assertIsNotNone(authorization, ''.join(output))
        expected = urlsplit(flow['authorizationUrl'])
        self.assertEqual((authorization.scheme, authorization.netloc, authorization.path),
                         (expected.scheme, expected.netloc, expected.path))
        parameters = parse_qs(authorization.query)
        self.assertEqual(parameters['response_type'], ['code'])
        self.assertEqual(parameters['code_challenge_method'], ['S256'])
        self.assertTrue(parameters['code_challenge'][0])
        self.assertTrue(parameters['state'][0])
        self.assertTrue(parameters['client_id'][0])
        self.assertNotIn('client_secret', parameters)
        self.assertEqual(set(parameters['scope'][0].split()), set(flow['scopes']))
        callback = urlsplit(parameters['redirect_uri'][0])
        self.assertEqual((callback.scheme, callback.hostname), ('http', '127.0.0.1'))
        self.assertIn(callback.port, [8484, 8485, 8486])
        self.assertEqual(self.server.requests, [])

    def test_client_identity_preserved_with_user_agent_suffix(self):
        self.environment['VAPI_USER_AGENT_SUFFIX'] = 'downstream/1.0'
        self.cli('assistants', 'list')
        first_identity = self.assert_client_headers(self.server.requests[-1][2], 'assistants_list', 'downstream/1.0')
        self.cli('assistants', 'list', '--user-agent-suffix', 'wrapper/2.0')
        second_identity = self.assert_client_headers(self.server.requests[-1][2], 'assistants_list', 'wrapper/2.0')
        self.assertEqual(first_identity, second_identity)

    def test_client_ci_flag(self):
        for value, expected in [('true', True), ('1', True), ('false', False), ('0', False)]:
            with self.subTest(ci=value):
                self.environment['CI'] = value
                self.cli('assistants', 'list')
                self.assert_client_headers(self.server.requests[-1][2], 'assistants_list', ci=expected)

    def test_client_identity_is_stable_across_retries(self):
        self.server.statuses = [503, 200]
        self.cli('assistants', 'list', '--retries', '1')
        self.assertEqual(len(self.server.requests), 2)
        identities = [self.assert_client_headers(request[2], 'assistants_list') for request in self.server.requests]
        self.assertEqual(identities[0], identities[1])

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
        self.assert_client_headers(headers, 'assistants_get')

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
            self.assert_client_headers(headers, 'assistants_create' if path == '/assistant' else 'calls_create')
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
        self.assert_client_headers(self.server.requests[0][2], 'assistants_delete')

    def test_api_errors_exit_nonzero(self):
        self.server.status = 401
        self.server.response = {'message': 'local unauthorized fixture'}
        result = self.cli('assistants', 'list', '--no-retry', success=False)
        self.assertIn('401', result.stdout + result.stderr)


if __name__ == '__main__':
    unittest.main()
