# Development

## Requirements

- Rust 1.99.1, pinned in `rust-toolchain.toml`; rustup installs it on first use.
- Python 3.11+ for repository scripts.
- `/usr/bin/time` for performance checks (GNU `time` on Linux; included on macOS).
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

Linux CI builds the release CLI and runs `python3 scripts/check-performance.py`
against root and assistant list/create schema discovery. Each command runs five
times and must finish within 500 ms, use at most 128 MiB of peak resident memory,
and produce at most 2 MiB of valid JSON. `/usr/bin/time` reports peak memory
after exit; there is no active memory sampling. A timeout stops hangs.

The checks also run with `make test` on Linux and macOS. For repeated local
measurements, build with `cargo build --release --locked`, then run
`python3 scripts/check-performance.py --runs 20`. The summary reports p50, p95,
and maximum wall time, plus maximum peak RSS and output size. One untimed
`--version` invocation warms the OS loader and file cache before the measurements
(with a two-second timeout). Every measured run must meet the limits. Each is a
fresh process, so normal per-invocation startup remains part of the measurement.
Use `VAPI_TEST_BINARY=/absolute/path/to/vapi` to compare another build.

## Source of truth

API generation is configured alongside the server SDKs in
[`VapiAI/docs/fern/apis/cli/generators.yml`](https://github.com/VapiAI/docs/blob/main/fern/apis/cli/generators.yml),
under `cli-sdk`. It consumes the shared OpenAPI definition, naming overrides,
parameter-description overlay and import settings directly. This repository
does not vendor a second API definition or generator configuration.

The CLI's separate Fern API workspace selects public-client OAuth with PKCE,
refresh, and the staging default server. It does not change the other SDKs'
authentication. Its `openapi-overrides.yml` also makes the provider-resource
POST's reserved `Content-Type` header optional; explicit `--content-type` remains
available and JSON requests set their content type through the HTTP stack.

## Generate

Run the **Release CLI** workflow in `VapiAI/docs`, supplying the desired version.
Like the SDK workflows, this runs `fern generate` in pull-request mode against
`VapiAI/cli`. Review the generated update and run this repository's checks before
merging it. Generation and merging do not publish a release.

For a local preview, clone `VapiAI/docs`, install its pinned Fern CLI with
Node.js 22+, authenticate to the Vapi Fern organization, and run there:

```bash
fern generate --api cli --group cli-sdk --version 1.0.0-beta.1 --preview
```

The generator remains pinned to `fern-cli-generator` 0.59.1, with profiles and
wire tests enabled and custom Rust commands disabled. Check generated
`cli/vapi/main.rs` for `PkceLoginFlow` when changing authentication. The staging
public client, issuer, resource parameters and default API server must be
promoted together after live acceptance; see [RELEASING.md](RELEASING.md).

## Preserve CLI customizations

`.fernignore` protects this repository's docs, verification workflow, smoke and
performance checks, toolchain and patch references from generation. Generated
runtime, wire tests and publishing workflows remain managed by Fern.

Fern Replay's native `.fern/replay.lock` records the schema-renderer and command-binding changes
on top of a pure `[fern-generated]` baseline. Central pull-request generation
reapplies the schema and command-attribution binding changes without a separate
CLI regeneration script. Commit future runtime customizations here;
review Replay's report and resolve conflicts in the generated PR. Do not disable
Replay for this generator while the schema fix is needed.

`src/http.rs` is also protected by `.fernignore` to retain Vapi request attribution.
Replay rejects that file when its snapshot contains Fern's embedded public test
certificate. Keep generator 0.59.1 pinned; when upgrading, compare this HTTP file
with fresh output and carry forward the attribution changes. The request-header
smoke tests verify the behavior. Other runtime files, including PKCE auth, remain
generated.

[`patches/`](patches/README.md) retains reviewable copies and explains the fixes.
Once Fern supplies equivalent behavior, remove the corresponding customization
with `fern replay forget` from a configured Fern checkout and review the next
generated diff. CI checks observable correctness and performance, not the
presence of a particular patch.

Endpoint names are maintained in
[`VapiAI/docs/fern/apis/api/openapi-overrides.yml`](https://github.com/VapiAI/docs/blob/main/fern/apis/api/openapi-overrides.yml)
for all generated server SDKs and this CLI. Review SDK version compatibility
when changing these names, then regenerate the affected SDKs and the central
`cli-sdk` group. Do not add CLI-only naming patches here.
