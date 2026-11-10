# Development

## Requirements

- Rust 1.99.1, pinned in `rust-toolchain.toml`; rustup installs it on first use.
- Node.js 22+ for the npm wrapper and Fern CLI.
- Python 3.11+ for repository scripts.
- On Linux, `pkg-config`, OpenSSL development headers, and D-Bus development
  headers for the native credential-store build. Release Linux binaries instead
  target musl and use the generated file credential store.

```bash
cargo build --locked
make test
```

CI runs the generated Rust tests, generated mock-server wire tests, Vapi smoke
checks, and npm installer tests on Linux, macOS, and Windows. No Vapi account or
live API writes are required. The smoke script can use another compiled binary
via `VAPI_TEST_BINARY=/absolute/path/to/vapi`.

## Source of truth

`fern/source.json` pins the `VapiAI/docs` commit, source directory, and SHA-256
hashes for the three vendored SDK inputs. They are copied unchanged from that
commit: the OpenAPI spec, SDK overrides, and parameter-description overlay.
`fern/generators.yml` retains the SDK spec import settings. The CLI does not
maintain its own API client or manually track SDK package releases.

`fern/cli-overrides.yml` contains the only CLI-specific API correction: the
provider-resource POST's reserved `Content-Type` header is optional, allowing
Fern's generated test inputs to invoke it. Explicit `--content-type` remains
available and JSON requests set their content type through the HTTP stack.

Update all three upstream inputs together from a reviewed docs commit:

```bash
python3 scripts/update-api.py <40-character-docs-commit-sha>
```

## Generate

Install Fern 5.151.1 (also pinned in `fern/fern.config.json`) and authenticate to
Vapi's Fern organization. Local-file output still uses Fern's cloud generation.

```bash
npm install -g fern-api@5.151.1
fern login
make generate
make test
```

The generator is pinned to `fern-cli-generator` 0.58.5. It emits a spec-only CLI
with named profiles and wire tests enabled. Custom Rust commands are disabled
so the beta evaluates the generated product directly.

**Use `make generate`, not direct generation into the checkout.** The script
copies Fern inputs into a temporary workspace, generates there, validates the
expected output, and replaces only these paths:

- `Cargo.toml`, `Cargo.lock`, `reference.md`
- `src/`, `cli/`, `tests/`, `wiremock/`

This isolates Fern's local-output Git handling from the checkout and preserves
repository-owned packaging, workflows, docs, licenses, and smoke checks. Do not
hand-edit generated Rust or wire tests. To check reproducibility without writing
files, run `python3 scripts/generate.py --check`.

The upstream spec currently validates with 12 warnings and no errors. Keep
warnings visible; they are not evidence that every API contract has been tested
against the live Vapi service.

`Cargo.toml` and `npm-wrapper/package.json` must have the same version. Change
the npm version, regenerate with the matching Cargo version (the generation
script reads it), and run `python3 scripts/check-version.py` before releasing.
