# Vapi CLI

Create voice assistants, make calls, and manage your Vapi resources from the
terminal.

## Install

Install with npm (requires Node.js):

```bash
npm install --global @vapi-ai/cli
```

Or install with Homebrew:

```bash
brew install VapiAI/tap/vapi-cli
```

You can also download a native binary from [GitHub Releases](https://github.com/VapiAI/cli/releases).

## Sign in

```bash
vapi auth login
```

Complete sign-in in your browser, then check your session:

```bash
vapi auth status
```

For a terminal without browser access, use `vapi auth login --no-browser` and
open the printed URL on a device with a browser. Use `vapi auth logout` to remove
your stored credentials.

## Create an assistant

Create an assistant with a voice, a model, and a short system prompt:

```bash
vapi assistants create --json '{
  "name": "My first assistant",
  "firstMessage": "Hi! I am Alex, your AI assistant. How can I help?",
  "model": {
    "provider": "openai",
    "model": "gpt-4.1-mini",
    "messages": [
      {
        "role": "system",
        "content": "You are Alex, a friendly voice assistant. Keep your answers brief and conversational."
      }
    ]
  },
  "voice": {
    "provider": "vapi",
    "voiceId": "Clara"
  }
}'
```

Save the returned assistant `id`. You can list your assistants or inspect one:

```bash
vapi assistants list
vapi assistants get --id YOUR_ASSISTANT_ID
```

## Make a call

Add or import an outbound phone number in the
[Vapi dashboard](https://dashboard.vapi.ai), then find its Vapi ID:

```bash
vapi phone-numbers list
```

Replace the assistant ID, phone number ID, and destination below with your own
values. The destination should be a phone number you control, in E.164 format.
Running this command starts an outbound call:

```bash
vapi calls create --json '{
  "assistantId": "YOUR_ASSISTANT_ID",
  "phoneNumberId": "YOUR_PHONE_NUMBER_ID",
  "customer": {
    "number": "+14155550123"
  }
}'
```

Save the returned call `id` to inspect its status and results:

```bash
vapi calls get --id YOUR_CALL_ID
vapi calls list --limit 10
```

See the [outbound calling guide](https://docs.vapi.ai/calls/outbound-calling) for
phone number setup and other call options.

## Explore and automate

Explore available commands and their options:

```bash
vapi --help
vapi assistants create --help
vapi --schema
```

Output defaults to tables in a terminal and JSON when piped. Use `--format json`
for explicit JSON output, or `--query` to select fields:

```bash
vapi assistants list --format json
vapi calls list --limit 10 --query '[].id'
```

Use JSON files for larger configurations, and `--dry-run` to preview a request
without sending it:

```bash
vapi assistants create --json - --dry-run < assistant.json
```

Delete commands send the request immediately, without a confirmation prompt.
Preview them with `--dry-run` first:

```bash
vapi assistants delete --id YOUR_ASSISTANT_ID --dry-run
```

The CLI also includes named profiles, shell completions, and local skill export
for coding agents. See [the command reference](reference.md),
`vapi profiles --help`, `vapi completion --help`, and `vapi generate-skills --help`.

## Develop

The CLI is generated from Vapi's shared API definitions with
[Fern](https://buildwithfern.com/learn/cli-generator/get-started/quickstart).
The `cli-sdk` group in [VapiAI/docs](https://github.com/VapiAI/docs/blob/main/fern/apis/cli/generators.yml)
uses the same API definition and naming overrides as the server SDKs, then opens
generated pull requests here. CLI-specific settings retain OAuth with PKCE.

Clone the repository and install the pinned Rust toolchain and
Python 3.11 or later; see [development requirements](DEVELOPMENT.md).

```bash
git clone https://github.com/VapiAI/cli.git
cd cli
cargo build --locked
make test
```

Generate updates from the shared configuration in `VapiAI/docs`; this repo
maintains the CLI runtime, tests and release gates.
[DEVELOPMENT.md](DEVELOPMENT.md) covers regeneration and customization.

## CI and releases

- `verify.yml` runs Rust tests, API wire tests, and local smoke checks on Linux,
  macOS, and Windows.
- Fern's `ci.yml` checks builds on pushes and publishes npm packages from
  version tags.
- Fern's `release.yml` plans releases on pull requests and builds GitHub release
  archives and installers from version tags. Stable releases also update
  Homebrew.

Maintainers should follow [RELEASING.md](RELEASING.md) for publishing setup and
the release checklist. [Migration notes](docs/MIGRATION.md) describe changes
from the previous CLI.

## License

Vapi-specific code is [MIT licensed](LICENSE). The generated Fern runtime is
attributed in [NOTICE](NOTICE) and distributed under [Apache License 2.0](LICENSE-APACHE).
