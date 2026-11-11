# Development

## Requirements

- Rust 1.99.1, pinned in `rust-toolchain.toml`; rustup installs it on first use.
- Node.js 22+ for the Fern CLI.
- Python 3.11+ for repository scripts.
- On Linux, `pkg-config`, OpenSSL development headers, and D-Bus development
  headers for the local native credential-store build. Fern owns release-target
  selection; inspect its publishing output before promising a Linux ABI.

```bash
cargo build --locked
make test
```

CI runs the generated Rust tests, generated mock-server wire tests, Vapi smoke
checks on Linux, macOS, and Windows. The verification workflow is configured;
GitHub has not provided PR workflow results for this candidate. No Vapi account or
live API writes are required. The smoke script can use another compiled binary
via `VAPI_TEST_BINARY=/absolute/path/to/vapi`.

## Source of truth

`fern/source.json` pins the `VapiAI/docs` commit, source directory, and SHA-256
hashes for the three vendored SDK inputs. They are copied unchanged from that
commit: the OpenAPI spec, SDK overrides, and parameter-description overlay.
`fern/generators.yml` retains the SDK spec import settings. The CLI does not
maintain its own API client or manually track SDK package releases.

`fern/cli-overrides.yml` adapts authentication and the default server for the
staging OAuth evaluation, and removes the unused SDK API-key scheme. It also
contains a CLI-specific API correction: the
provider-resource POST's reserved `Content-Type` header is optional, allowing
Fern's generated test inputs to invoke it. Explicit `--content-type` remains
available and JSON requests set their content type through the HTTP stack.

Update all three upstream inputs together from a reviewed docs commit:

```bash
python3 scripts/update-api.py <40-character-docs-commit-sha>
```

## Generate

Install Fern 5.202.0 (also pinned in `fern/fern.config.json`) and authenticate to
Vapi's Fern organization. Local-file output still uses Fern's cloud generation.

```bash
npm install -g fern-api@5.202.0
fern login
make generate
make test
```

The generator is pinned to `fern-cli-generator` 0.59.1. It emits a spec-only CLI
with named profiles and wire tests enabled. Custom Rust commands are disabled
so the beta evaluates the generated product directly. `auth-schemes` and
`api.auth` select the generated public PKCE login flow. The matching OpenAPI
OAuth scheme is required too; check that regenerated `cli/vapi/main.rs` contains
`PkceLoginFlow` before treating login as implemented.

**Use `make generate`, not direct generation into the checkout.** The script
copies Fern inputs into a temporary workspace, generates there, validates the
expected output, and replaces only these paths:

- `Cargo.toml`, `Cargo.lock`, `reference.md`
- `src/`, `cli/`, `tests/`, `wiremock/`

This isolates Fern's local-output Git handling from the checkout and preserves
repository-owned verification, docs, licenses, and smoke checks. Do not
hand-edit generated Rust or wire tests. Temporary upstream runtime fixes live in
[`patches/`](patches/README.md) and are applied to both output groups before
comparison or adoption. Regeneration fails if a patch no longer applies; review
and remove it once Fern supplies the fix. To check reproducibility without
writing files, run `python3 scripts/generate.py --check`.

The upstream spec currently validates with 12 warnings and no errors. Keep
warnings visible; they are not evidence that every API contract has been tested
against the live Vapi service.

The generated runtime uses the registered staging OAuth client. Local mock
tests do not establish successful OAuth login or live API acceptance. See RELEASING.md for the required integration.

The `publish` group delegates npm, native releases and Homebrew to Fern. Use
`make preview-publish` to inspect its output in `generated-publish/`. It does
not publish in preview mode. `make generate-publish` adopts the runtime plus
`dist-workspace.toml`, `.github/workflows/ci.yml` and
`.github/workflows/release.yml` from that isolated preview. Run
`python3 scripts/generate.py --publishing --check` to check all generated files.
The generated npm packaging scripts live inside `ci.yml`.

Command names belong in source OpenAPI annotations in the `vapi` repository.
For example, Nest's `@ApiExtension('x-fern-sdk-group-name', ['campaigns'])`
and `@ApiExtension('x-fern-sdk-method-name', 'list')` on a controller operation
produce `vapi campaigns list` when propagated through the shared API inputs.
Versioned endpoints must receive distinct group/method paths to avoid collisions.
Those annotations also affect generated SDK names, so review the SDK surface
when making the upstream change. Do not add command-naming patches here.
