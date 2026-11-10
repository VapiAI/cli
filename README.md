# Vapi CLI

The Vapi command-line interface, generated from Vapi's shared API definitions
with [Fern's CLI generator](https://buildwithfern.com/learn/cli-generator/get-started/quickstart).
This branch is the **1.0.0-beta.1 evaluation candidate**, replacing the old Go CLI.
See the [comparison and recorded decisions](docs/MIGRATION.md).

**This draft is not ready to release.** OAuth integration and access to Fern's
npm/GitHub publishing generator are outstanding. The checked-in inspection build
still embeds the earlier bearer-auth API spec; it does not implement the agreed
OAuth login yet. An API-key-based release is not the intended fallback.

## Inspect the generated CLI

Install the pinned [Rust toolchain](https://rustup.rs/), then:

```bash
cargo build --locked
./target/debug/vapi --help
./target/debug/vapi --schema
./target/debug/vapi assistants create --json '{"name":"Example"}' --dry-run
./target/debug/vapi completion zsh
```

The current snapshot generates 145 operations in 23 groups. Commands use Fern's
conventions, such as `assistants get --id <id>` and `phone-numbers list`.
Controller-derived command names will be shortened at their source in `vapi`'s
OpenAPI annotations, then flow through the shared API definition.

Deletion sends the API request immediately. Use `--dry-run` to inspect a request
without sending it. Output defaults to tables in terminals and JSON when piped;
use explicit `--format json` in scripts. See [reference.md](reference.md) and
`vapi <group> --help` for the current generated surface.

## Distribution target

Fern will manage npm, native GitHub release downloads, and Homebrew through
`VapiAI/homebrew-tap` (`vapi-cli`). The `publish` generator group uses pull-request
mode, OIDC npm publishing, and Fern's Homebrew distribution integration. Its
publishing preview is currently denied by Fern's organization permissions;
generated publishing files have not been adopted or tested yet.

The old custom npm wrapper and publishing workflows are removed. `init`,
`listen`, legacy configuration migration, and the bundled MCP documentation
server are also removed. Project initialization belongs in a separate project.

## Agent skill export

`vapi generate-skills --output-dir ./skills` exports local Markdown instructions
for coding agents: one shared guide and one file per API group. It describes
commands, authentication sources, flags, and output. It does not install the
files into an IDE or modify agent settings. Keeping this generated command is
still a product decision; its output must be regenerated after OAuth changes.

## Development

[Development](DEVELOPMENT.md) describes regeneration and API provenance;
[releasing](RELEASING.md) records OAuth and publishing prerequisites.
Vapi-specific code is MIT licensed. The generated Fern runtime is attributed in
[NOTICE](NOTICE) and distributed with [Apache License 2.0](LICENSE-APACHE).
