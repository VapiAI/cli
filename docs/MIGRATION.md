# Evaluate the Fern CLI beta

Recommendation: adopt the generated API foundation and evaluate it as a beta.
It replaces much of the handwritten API plumbing and several dashboard-only
placeholders with scriptable API operations. It is not yet a drop-in replacement
for users relying on browser login, project scaffolding, workflows, or local
webhook forwarding. Decide on those losses and command naming before stable 1.0.

## What was compared

- Go CLI at `4149790e4621407c903f73554130fb31a97d9334`, the latest main used for
  this migration, including the removal of MCP/skills setup.
- Fern CLI generator 0.58.5, Fern CLI 5.151.1, package version 1.0.0-beta.1.
- Vapi SDK API inputs pinned in `fern/source.json`.
- The compiled candidate's `--help`, `--schema`, dry runs, generated tests, and
  local mock API requests. No production calls, purchases, or deletions were run.

The candidate reports **145 API operations in 23 groups**, plus auth, profiles,
completion, man pages, and skill export. Those counts describe this pinned API
snapshot, not a guarantee of all current Vapi API functionality.

## Capability comparison

| Area | Go CLI behavior | Generated candidate | Assessment |
| --- | --- | --- | --- |
| API breadth | Handwritten commands for a small subset of the old Go SDK | 145 operations, including squads, files, sessions, simulations, structured outputs, scorecards and traffic allocation | Clear gain; future coverage follows the shared spec |
| Assistant CRUD | Working API operations with interactive creation and deletion confirmation | `assistants list/create/get/update/delete`, JSON/stdin bodies and flags | Keep generated API commands; guided creation is removed |
| Calls | List/get worked; create/update/end printed SDK or dashboard guidance | Real list/create/update/delete plus spec-defined call operations | Gain; do not equate deleting a call record with ending a live call |
| Chat | List/get/delete worked; creation and continuation printed guidance | Real create/list/get/delete and OpenAI-compatible response operation | Gain; inspect streaming behavior against a real test assistant before promising parity |
| Tools and phones | List/get/delete worked; several create/update/test paths printed guidance | Spec-driven create/update and other operations | More useful automation; no replacement for custom tool-testing helpers |
| Campaigns | Working list/get/create/update/delete with guided selections | V1/V2 API operations and contact retrieval | Coverage gain, but long controller-derived command names need a product decision |
| Workflows | Working list/get/create/delete in the old SDK | No workflow endpoints in this API snapshot | Real removal; use an existing client/dashboard, or add supported endpoints upstream before claiming parity |
| Browser login | Vapi dashboard callback and local account records | Bearer API key from `VAPI_TOKEN`, stdin login, credential store | Browser flow removed; spec has no OAuth flow to generate |
| Multiple accounts | `auth switch` and YAML config | Named profiles with per-profile credential slots and settings | Good replacement, with different commands and credential precedence |
| Output | Custom tables/prose and JSON depending on command | JSON when piped, tables in a terminal; explicit formats and JMESPath | Better for automation; existing output parsers must change |
| Request control | Uneven or limited flags | JSON/stdin input, dry run, base URL override, request flags, retries, schema and spec export | Strong gain; retry and pagination semantics still depend on spec metadata |
| Deletion | Interactive confirmation | Immediate API request; `--dry-run` available | Intentional breaking behavior; scripts and operators must account for it |
| Project setup | `init` detected frameworks, installed SDKs and wrote examples/env files | No equivalent with custom commands disabled | Real removal; decide whether users still need a separately maintained setup tool |
| Webhook forwarding | `listen` hosted a local forwarding HTTP server; it did not establish the advertised public tunnel | No equivalent | Real removal of local forwarding; use a separate tunnel and receiver |
| Webhook management and logs | Commands mostly printed guidance instead of invoking APIs | No standalone webhook/log command groups in this snapshot | Removes misleading promises, not working CRUD/log retrieval |
| IDE integration | MCP/skills setup already removed from main | No IDE installer; Fern adds `generate-skills` to export API command skills | Preserve the setup removal; decide whether the new export belongs in the public help |
| Updates | Custom self-update and SDK-update automation | Regenerate from pinned shared API inputs; upgrade binaries via npm/releases | Smaller maintenance surface; self-update command removed |
| Telemetry | Handwritten PostHog analytics and opt-out configuration | No Vapi PostHog analytics package or analytics config commands | Old telemetry configuration is removed |
| Distribution | GoReleaser, npm download wrapper, install scripts, Docker | Generated Rust CLI; adapted npm wrapper and five native release targets | Keep existing package and executable names; Docker and 32-bit ARM are dropped |
| MCP docs server | Separate TypeScript package in the repository | Removed from the new CLI tree | Existing published versions are unaffected; future server maintenance needs its own home |

## Breaking command examples

| Before | Candidate |
| --- | --- |
| `export VAPI_API_KEY=...` | `export VAPI_TOKEN=...` |
| `vapi login` | `vapi auth login --with-token` (API key via stdin) |
| `vapi auth switch staging` | `vapi profiles use staging` |
| `vapi assistant list` | `vapi assistants list` |
| `vapi assistant get ID` | `vapi assistants get --id ID` |
| `vapi phone list` | `vapi phone-numbers list` |
| `vapi call get ID` | `vapi calls get --id ID` |
| `vapi chat list` | `vapi chats list` |
| `vapi tool get ID` | `vapi tools get --id ID` |
| `vapi campaign list` | `vapi campaigns campaign-controller-find-all` |
| `vapi manual` | `vapi man` |
| `vapi version` | `vapi --version` |

Old YAML credentials/settings are not imported. A selected profile's stored
credential overrides the exported token. Inspect `profiles current` and
`auth status` when switching environments. Resource IDs use named flags; singular
resource command aliases are not installed.

## Decisions before stable 1.0

1. **Command names:** accept plural groups and `--id`, and choose concise method
   names for campaign/provider/other controller-derived operations in the API
   spec. Avoid hand-editing generated command code.
2. **Authentication:** accept API-key login and profiles, or model a supported
   OAuth flow upstream. Existing dashboard login is not reproduced automatically.
3. **Custom workflows:** explicitly accept removing `init`, `listen`, and legacy
   workflows, or restore only the ones people use through Fern custom commands.
4. **Destructive operations:** accept immediate deletion, or add a separately
   reviewed UX policy. Dry run exists but is opt-in.
5. **Real API verification:** test a disposable assistant, chat streaming, file
   upload, pagination on paginated endpoints, and profiles on each supported OS.
   Mock success does not prove the shared API spec perfectly describes production.
6. **Publishing:** validate the native release matrix and existing npm OIDC
   configuration before tagging. Fern-managed npm/GitHub preview access is still
   unresolved; the current publication path does not require it.

The generated provider-resource wire tests initially failed because the spec
marked reserved `Content-Type` as a required standalone header while Fern's IR
removed it. `fern/cli-overrides.yml` makes it optional without editing the shared
SDK files or generated Rust. Keep that exception visible when updating inputs.
