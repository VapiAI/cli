# Migrating from the Go CLI

The 1.0 beta replaces the Go CLI with generated API commands and OAuth login.
Command names and flags change. Review scripts before upgrading; existing
configuration files and API keys are not imported.

## Install the intended version

Use the beta npm channel while evaluating 1.0:

```bash
npm install --global @vapi-ai/cli@beta
vapi --version
```

Homebrew is available starting with the stable release, not the prerelease.
Native archives are available on the matching GitHub prerelease. The generated
shell installer uses `$CARGO_HOME/bin`, or `~/.cargo/bin` when unset.

An old installation in `~/.vapi/bin` can shadow the new binary. Check which
executable your shell finds:

```bash
type -a vapi
```

Remove the obsolete installation or its `~/.vapi/bin` PATH entry, then open a
new shell or run `hash -r`. On PowerShell, use `Get-Command vapi -All`.

## Sign in again

```bash
vapi auth login
vapi auth status
```

Login uses a public OAuth client with PKCE. `auth status` reports available
credentials; it does not export tokens or establish your account identity.
`auth logout` removes local credentials; it does not revoke a token at the
identity provider. Sign out before switching from a staging OAuth registration
to production, and sign in again afterward.

For automation, set `VAPI_TOKEN` to a short-lived OAuth access token for the
configured API audience. The CLI uses it without saving it or refreshing it.
This variable does not restore support for legacy API keys. Without a selected
profile, the environment token takes precedence over a saved login. A selected
profile's saved credential takes precedence; unset `VAPI_TOKEN` to avoid an
unintended fallback when that profile has no credential.

For SSH login, the browser must reach the CLI's loopback callback. Forward all
configured ports from your local machine, then run login on the remote host:

```bash
ssh -L 8484:127.0.0.1:8484 -L 8485:127.0.0.1:8485 -L 8486:127.0.0.1:8486 user@host
vapi auth login --no-browser
```

Open the printed URL in the local browser. Opening it on an unrelated device
without forwarding the callback cannot complete login. Containers similarly
need reachable loopback forwarding, or an OAuth access token for automation.

## Update commands and flags

Use `vapi --help`, `vapi <group> --help`, and `vapi <group> <command> --schema`
to inspect the generated commands and accepted parameters.

| Go CLI | 1.0 beta |
| --- | --- |
| `assistant list/create/get/update/delete` | `assistants list/create/get/update/delete`; inspect the named flags with `--help` |
| `call`, `chat`, `tool`, `phone-number`, `squad`, `file` groups | `calls`, `chats`, `tools`, `phone-numbers`, `squads`, `files` |
| `--api-key` and saved API-key login | Removed; use OAuth login or an OAuth `VAPI_TOKEN` |
| `--config` | Removed; use generated profiles and their supported environment variables |
| `-v` | Use `--version` |
| `auth whoami` | No identity-query equivalent; `auth status` reports credential status |
| `auth token` | Removed; the CLI does not export stored tokens |
| `auth switch` | Use `profiles use` or `--profile` with a newly configured profile |
| `auth logout --all` | Removed; sign out of the relevant selected profiles individually |
| `config get/set/env` and `config analytics ...` | Removed; see `profiles --help` and the supported global options |
| `call end` | Removed; there is no generated equivalent to this convenience command |
| `chat continue` | Removed; inspect `chats --help` for the API's current chat operations |
| `tool test` and `tool types` | Removed; inspect `tools --help` and operation schemas; test tool servers separately |
| `init` and `listen` | Removed |
| Bundled MCP documentation server | Removed; use the separately maintained Vapi MCP server |
| `install-man-pages` | Use `vapi man > vapi.1`, then install the resulting page using your system's man-page setup |
| `update` and `update check` | Removed; reinstall through the chosen npm channel or native release |

Shared SDK names also simplify commands such as `campaigns list` and
`campaigns list-v2`. Existing server SDK code may require method/request-type
renames when adopting regenerated SDK versions.

## Check scripts and destructive operations

Deletes execute immediately without an interactive confirmation. Preview the
request with `--dry-run` before running destructive commands:

```bash
vapi assistants delete --id ASSISTANT_ID --dry-run
```

Check exit status before consuming command output. A list command without
supported pagination options returns one server page, not necessarily every
matching record. Inspect its `--help`; changing `--limit` alone does not ensure
a complete export.

## Platform and distribution changes

The published npm platforms are Linux x86-64/AArch64, macOS x86-64/AArch64, and
Windows x86-64. Linux ARMv6/ARMv7 builds and the old Docker image are no longer
published. There is no built-in self-update command. Use the package manager or
native archive that installed the CLI.

See [the command reference](../reference.md) for all API operations and
[release instructions](../RELEASING.md) for maintainer rollout requirements.
