"""Check schema discovery performance without credentials or API requests."""
import json
import os
from pathlib import Path
import subprocess
import sys
import tempfile
import time

ROOT = Path(__file__).resolve().parent.parent
BINARY = Path(os.environ.get('VAPI_TEST_BINARY', ROOT / 'target/debug/vapi')).resolve()
MAX_SECONDS = 5
MAX_RSS_BYTES = 256 * 1024 * 1024
MAX_OUTPUT_BYTES = 4 * 1024 * 1024
COMMANDS = [('--schema',), ('assistants', 'list', '--schema'),
            ('assistants', 'create', '--schema')]


def resident_bytes(process):
    if sys.platform.startswith('linux'):
        try:
            status = Path(f'/proc/{process.pid}/status').read_text()
        except FileNotFoundError:
            return 0  # The child exited between polling and reading its status.
        for line in status.splitlines():
            if line.startswith('VmHWM:'):
                return int(line.split()[1]) * 1024
        if process.poll() is not None:
            return 0
        raise RuntimeError('Child memory usage is unavailable in /proc')
    result = subprocess.run(['ps', '-o', 'rss=', '-p', str(process.pid)],
                            capture_output=True, text=True, timeout=1)
    if result.returncode == 0 and result.stdout.strip():
        return int(result.stdout.strip()) * 1024
    if process.poll() is not None:
        return 0
    raise RuntimeError('Unable to measure child resident memory with ps')


def check_command(arguments):
    with tempfile.TemporaryDirectory(prefix='vapi-performance-') as directory:
        environment = {k: v for k, v in os.environ.items() if not k.startswith('VAPI_')}
        environment.update(FERN_CLI_CREDENTIAL_STORE='file',
                           XDG_CONFIG_HOME=directory, APPDATA=directory)
        with tempfile.TemporaryFile() as stdout, tempfile.TemporaryFile() as stderr:
            started = time.monotonic()
            process = subprocess.Popen([str(BINARY), *arguments], cwd=directory,
                                       env=environment, stdin=subprocess.DEVNULL,
                                       stdout=stdout, stderr=stderr)
            peak_rss = 0
            try:
                while process.poll() is None:
                    peak_rss = max(peak_rss, resident_bytes(process))
                    elapsed = time.monotonic() - started
                    if elapsed > MAX_SECONDS:
                        raise AssertionError(f'exceeded {MAX_SECONDS}s runtime budget')
                    if peak_rss > MAX_RSS_BYTES:
                        raise AssertionError('exceeded 256 MiB resident memory budget')
                    if os.fstat(stdout.fileno()).st_size > MAX_OUTPUT_BYTES:
                        raise AssertionError('exceeded 4 MiB output budget')
                    time.sleep(0.02)
            finally:
                if process.poll() is None:
                    process.kill()
                process.wait(timeout=3)
            elapsed = time.monotonic() - started
            size = os.fstat(stdout.fileno()).st_size
            stderr.seek(0)
            assert process.returncode == 0, stderr.read(4096).decode(errors='replace')
            assert elapsed <= MAX_SECONDS, f'exceeded {MAX_SECONDS}s runtime budget'
            assert size <= MAX_OUTPUT_BYTES, 'exceeded 4 MiB output budget'
            stdout.seek(0)
            json.load(stdout)
            print(f'{" ".join(arguments)}: {elapsed:.2f}s, '
                  f'{peak_rss / 1024 / 1024:.1f} MiB RSS, {size:,} bytes', flush=True)


if __name__ == '__main__':
    if sys.platform != 'darwin' and not sys.platform.startswith('linux'):
        raise SystemExit('Performance checks require Linux or macOS memory accounting')
    for command in COMMANDS:
        try:
            check_command(command)
        except (AssertionError, RuntimeError) as error:
            raise SystemExit(f'{" ".join(command)}: {error}') from error
