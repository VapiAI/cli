# Vapi CLI

The Vapi command-line interface, generated from Vapi's shared API definitions
with [Fern's CLI generator](https://buildwithfern.com/learn/cli-generator/get-started/quickstart).
This branch is the **1.0.0-beta.1 evaluation candidate**, replacing the old Go CLI.
See the [comparison and recorded decisions](docs/MIGRATION.md).

**Staging-only evaluation build; not ready to release.** Fern generates browser
OAuth login with PKCE and automatic refresh using the registered WorkOS public
client. The default API server is `https://staging-api.vapi.ai`. API acceptance
changes, a real-user login/refresh check, production OAuth promotion and Fern
publishing credential isolation are required before release. No client secret is embedded.

## Inspect the generated CLI

Install the pinned [Rust toolchain](https://rustup.rs/), then:

```bash
cargo build --locked
./target/debug/vapi --help
./target/debug/vapi --schema
./target/debug/vapi auth login --no-browser
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

Fern generated the checked-in `ci.yml` and `release.yml` workflows for npm,
GitHub Release downloads/installers, and Homebrew through `VapiAI/homebrew-tap`
(`vapi-cli`). Version-tag pushes trigger publishing; merging alone does not.
The npm launcher and platform package manifests are created by the generated
workflow at release time.

The publishing preview now succeeds. OAuth acceptance, the generated workflow's
credential isolation, npm trusted publishing, Homebrew access, and target-runner
validation remain release gates; see [releasing](RELEASING.md).

The old custom npm wrapper and publishing workflows are removed. `init`,
`listen`, legacy configuration migration, and the bundled MCP documentation
server are also removed. Project initialization belongs in a separate project.

## Agent skill export

`vapi generate-skills --output-dir ./skills` exports local Markdown instructions
for coding agents: one shared guide and one file per API group. It describes
commands, authentication sources, flags, and output. It does not install the
files into an IDE or modify agent settings. Its OAuth setup guidance is currently
incomplete. Keeping this generated command is still a product decision.

## Development

[Development](DEVELOPMENT.md) describes regeneration and API provenance;
[releasing](RELEASING.md) records OAuth and publishing prerequisites.
Vapi-specific code is MIT licensed. The generated Fern runtime is attributed in
[NOTICE](NOTICE) and distributed with [Apache License 2.0](LICENSE-APACHE).
