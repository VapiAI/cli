# Vapi CLI

The Vapi command-line interface, generated from the same API definitions used by
Vapi's server SDKs with [Fern's CLI generator](https://buildwithfern.com/learn/cli-generator/get-started/quickstart).

This branch is the **1.0.0-beta.1 candidate** replacing the handwritten Go CLI.
It has breaking command, authentication, and output changes. See the
[comparison and migration guide](docs/MIGRATION.md) before upgrading.

## Try the candidate

The beta has not been published yet. Install the pinned [Rust toolchain](https://rustup.rs/), then:

```bash
cargo build --locked
./target/debug/vapi --help
export VAPI_TOKEN='<your private Vapi API key>'
./target/debug/vapi assistants list --limit 5
./target/debug/vapi assistants get --id '<assistant-id>'
```

After a beta release is published, Node.js 22+ users can install the same native CLI:

```bash
npm install -g @vapi-ai/cli@next
```

Supported release platforms: macOS Intel and Apple Silicon, Linux x64 and ARM64
(static musl binaries), and Windows x64. Windows installation requires `tar.exe`
(available on Windows 10 1803+). Releases provide native archives and SHA-256 checksums.

## Authentication

`VAPI_TOKEN` replaces `VAPI_API_KEY`. The old `~/.vapi-cli.yaml` is not imported.
The generated `vapi auth login --with-token` reads an API key from stdin and stores
it in the OS credential store. This spec declares bearer authentication, so
`vapi auth login` does **not** reproduce the old dashboard browser login flow.

Named profiles replace account switching:

```bash
vapi profiles create staging --from-env
vapi profiles use staging
vapi assistants list --profile staging
vapi auth status
```

A selected profile's stored credential takes precedence over `VAPI_TOKEN`.
Use `vapi profiles current` to inspect which profile applies.

## API commands

```bash
vapi assistants create --json '{"name":"Example"}' --dry-run
vapi assistants create --json - < assistant.json
vapi calls list --limit 10 --format json
vapi phone-numbers list
vapi squads list
vapi assistants get --id '<assistant-id>' --query name
vapi --schema
vapi completion zsh
```

API deletions execute without the old confirmation prompt. Use `--dry-run` to
inspect requests before sending them. Use explicit `--format json` in scripts;
interactive output defaults to a table, while piped output defaults to JSON.

See [reference.md](reference.md) for generated endpoint documentation and run
`vapi <group> --help` for exact commands. Some endpoints currently inherit long
controller names from the API spec; the migration guide records these rough edges.

## Development

The checked-in Rust runtime, embedded API spec, reference, and wire tests are
Fern output. [Development](DEVELOPMENT.md) describes regeneration and API provenance;
[releasing](RELEASING.md) describes beta and stable publishing.

Vapi-specific code is MIT licensed. The generated Fern runtime is attributed in
[NOTICE](NOTICE) and distributed with [Apache License 2.0](LICENSE-APACHE).
