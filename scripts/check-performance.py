"""Check release CLI schema performance using OS-reported peak resident memory."""
import argparse
import json
import math
import os
from pathlib import Path
import re
import signal
import statistics
import subprocess
import sys
import tempfile
import time

ROOT = Path(__file__).resolve().parent.parent
BINARY = Path(os.environ.get('VAPI_TEST_BINARY', ROOT / 'target/release/vapi')).resolve()
MAX_SECONDS = 0.5
MAX_RSS_BYTES = 128 * 1024 * 1024
MAX_OUTPUT_BYTES = 2 * 1024 * 1024
COMMANDS = [('--schema',), ('assistants', 'list', '--schema'),
            ('assistants', 'create', '--schema')]


def cli_environment(directory):
    environment = {k: v for k, v in os.environ.items() if not k.startswith('VAPI_')}
    environment.update(FERN_CLI_CREDENTIAL_STORE='file', LC_ALL='C',
                       XDG_CONFIG_HOME=directory, APPDATA=directory)
    return environment


def warm_up():
    with tempfile.TemporaryDirectory(prefix='vapi-performance-warmup-') as directory:
        subprocess.run([str(BINARY), '--version'], cwd=directory,
                       env=cli_environment(directory), stdin=subprocess.DEVNULL,
                       stdout=subprocess.DEVNULL, stderr=subprocess.PIPE,
                       check=True, timeout=2)
    print('Untimed warmup: --version', flush=True)


def check_command(arguments):
    with tempfile.TemporaryDirectory(prefix='vapi-performance-') as directory:
        # BSD time reports bytes; GNU time's %M reports KiB. Both obtain the
        # child's peak RSS from the OS after exit, without active sampling.
        flags = ['-l'] if sys.platform == 'darwin' else ['-f', 'VAPI_MAX_RSS_KIB=%M']
        with tempfile.TemporaryFile() as stdout:
            started = time.perf_counter()
            process = subprocess.Popen(['/usr/bin/time', *flags, str(BINARY), *arguments],
                                       cwd=directory, env=cli_environment(directory),
                                       stdin=subprocess.DEVNULL, stdout=stdout,
                                       stderr=subprocess.PIPE, start_new_session=True)
            try:
                remaining = max(0, MAX_SECONDS - (time.perf_counter() - started))
                _, stderr = process.communicate(timeout=remaining)
            except subprocess.TimeoutExpired:
                # Kill both time and the CLI, including any descendants.
                try:
                    os.killpg(process.pid, signal.SIGKILL)
                except ProcessLookupError:
                    pass
                process.communicate(timeout=3)
                raise AssertionError(f'exceeded {MAX_SECONDS}s runtime budget') from None
            elapsed = time.perf_counter() - started
            error_text = stderr.decode(errors='replace')
            assert process.returncode == 0, error_text
            assert elapsed <= MAX_SECONDS, f'exceeded {MAX_SECONDS}s runtime budget'
            pattern = r'(\d+)\s+maximum resident set size' if sys.platform == 'darwin' else r'VAPI_MAX_RSS_KIB=(\d+)'
            usage = re.search(pattern, error_text)
            assert usage, f'Peak resident memory missing from /usr/bin/time output: {error_text}'
            peak_rss = int(usage.group(1)) * (1 if sys.platform == 'darwin' else 1024)
            assert peak_rss <= MAX_RSS_BYTES, 'exceeded 128 MiB peak resident memory budget'
            size = os.fstat(stdout.fileno()).st_size
            assert size <= MAX_OUTPUT_BYTES, 'exceeded 2 MiB output budget'
            stdout.seek(0)
            json.load(stdout)
            return elapsed, peak_rss, size


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--runs', type=int, default=5, help='Measured runs per command (default: 5)')
    args = parser.parse_args()
    if args.runs < 1:
        parser.error('--runs must be positive')
    if sys.platform != 'darwin' and not sys.platform.startswith('linux'):
        parser.error('Performance checks require Linux or macOS /usr/bin/time')
    warm_up()
    for command in COMMANDS:
        try:
            samples = [check_command(command) for _ in range(args.runs)]
        except (AssertionError, RuntimeError) as error:
            raise SystemExit(f'{" ".join(command)}: {error}') from error
        durations = sorted(sample[0] for sample in samples)
        p95 = durations[math.ceil(0.95 * len(durations)) - 1]
        print(f'{" ".join(command)} ({args.runs} runs): '
              f'p50 {statistics.median(durations) * 1000:.1f} ms, '
              f'p95 {p95 * 1000:.1f} ms, max {max(durations) * 1000:.1f} ms; '
              f'peak {max(sample[1] for sample in samples) / 1024 / 1024:.1f} MiB RSS; '
              f'max {max(sample[2] for sample in samples):,} bytes', flush=True)
