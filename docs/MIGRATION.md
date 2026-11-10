# Evaluate the Fern CLI beta

Decision: adopt Fern conventions and the generated API foundation. OAuth and
Fern-managed publishing are required before this candidate is released.
It replaces much of the handwritten API plumbing and several dashboard-only
placeholders with scriptable API operations. It is not yet a drop-in replacement
for users relying on browser login, project scaffolding, workflows, or local
webhook forwarding. `init` and `listen` remain removed; command naming should be fixed upstream.

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
| Browser login | Vapi dashboard callback and local account records | Current inspection build is bearer-only; target is generated OAuth login/refresh | OAuth is required; API client registration and token contract are pending |
| Multiple accounts | `auth switch` and YAML config | Fern profiles; no legacy credential migration | Accept Fern conventions; finalize profile behavior with OAuth |
| Output | Custom tables/prose and JSON depending on command | JSON when piped, tables in a terminal; explicit formats and JMESPath | Better for automation; existing output parsers must change |
| Request control | Uneven or limited flags | JSON/stdin input, dry run, base URL override, request flags, retries, schema and spec export | Strong gain; retry and pagination semantics still depend on spec metadata |
| Deletion | Interactive confirmation | Immediate API request; `--dry-run` available | Intentional breaking behavior; scripts and operators must account for it |
| Project setup | `init` detected frameworks, installed SDKs and wrote examples/env files | Removed | Accepted: initialization belongs in a separate project |
| Webhook forwarding | `listen` hosted a local forwarding HTTP server; it did not establish the advertised public tunnel | Removed | Accepted: rebuild later only if missed |
| Webhook management and logs | Commands mostly printed guidance instead of invoking APIs | No standalone webhook/log command groups in this snapshot | Removes misleading promises, not working CRUD/log retrieval |
| IDE integration | MCP/skills setup already removed from main | No IDE installer; Fern adds `generate-skills` to export API command skills | Preserve the setup removal; decide whether the new export belongs in the public help |
| Updates | Custom self-update and SDK-update automation | Regenerate from pinned shared API inputs; upgrade binaries via npm/releases | Smaller maintenance surface; self-update command removed |
| Telemetry | Handwritten PostHog analytics and opt-out configuration | No Vapi PostHog analytics package or analytics config commands | Old telemetry configuration is removed |
| Distribution | GoReleaser, custom npm wrapper, stale Homebrew, installers, Docker | Fern-managed npm, native downloads and Homebrew configured | Custom wrapper/publisher removed; Fern publishing access is blocked |
| MCP docs server | Separate TypeScript package in the repository | Removed from the new CLI tree | Existing published versions are unaffected; future server maintenance needs its own home |

## Breaking command examples

| Before | Candidate |
| --- | --- |
| `export VAPI_API_KEY=...` | OAuth target; no legacy API-key migration |
| `vapi login` | Target: `vapi auth login` with OAuth; not implemented yet |
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

## Recorded product decisions

| Area | Decision | Implementation status |
| --- | --- | --- |
| Command grammar | Keep Fern's groups, flags, and conventions | Present in generated output |
| Long controller names | Shorten through annotations in `vapi`, propagated through OpenAPI | Mechanism confirmed; upstream source change remains to be made |
| Authentication | OAuth, aligned with the other Vapi surfaces | Required; current inspection build is still bearer-only |
| Old credentials | No API-key/YAML compatibility migration | Removed; no new migration path planned |
| `init` | Separate project | Removed |
| `listen` | Delete; rebuild only if missed | Removed |
| Deletion | Keep immediate deletes and document dry run | Present and tested with local mocks |
| Output | Keep Fern defaults | Present |
| Installation | npm, native downloads, and Homebrew | Fern publish group configured; generator permission blocked |
| MCP docs server | Remove | Removed |
| Publishing | Fern-managed | Custom publisher and npm wrapper removed; generated publisher pending access |
| `generate-skills` | User requested more information; keep/remove undecided | Built-in export remains in the inspection build |

## Skill export: what the command does

The compiled inspection build was run with:

```bash
vapi generate-skills --output-dir /tmp/vapi-cli-skills
```

It wrote 24 local files: `vapi-shared/SKILL.md` plus one for each API group,
such as `vapi-assistants/SKILL.md`. The shared guide documents auth sources,
global flags, formats and dry run. Group guides list operations and explain
how an agent can discover parameters through `--schema`.

It uses deterministic templates over the embedded spec, not an LLM. It does
not install files into Codex/Claude, configure an IDE, read account resources,
or embed the user's credential values. Its auth guidance currently reflects
the bearer inspection build; regenerate and inspect it after OAuth changes.

Recommendation remains to keep this built-in exporter because it tracks the
exact installed CLI, while any automatic agent setup stays outside this repo.
The user's decision on retaining it is still pending.

## Remaining integration work

1. Obtain the native CLI's public-client OAuth registration and API token
   contract. The MCP server's confidential exchange-client credentials are
   not usable in a distributed CLI. Generate OAuth login and test refresh,
   org selection/isolation, logout and permission denial.
2. Put shorter names in `vapi` source annotations. Nest `@ApiExtension` can
   emit `x-fern-sdk-group-name` and `x-fern-sdk-method-name`. Review SDK naming
   changes and V1/V2 collisions, then update the shared API snapshot here.
3. Enable Fern's publishing-generator access for the `vapi` organization and
   adopt its actual generated workflow/package files. Configure npm OIDC and
   Homebrew credentials; do not replace this with another custom publisher.
4. Resolve the skill-export decision, run cross-platform checks and perform
   disposable live API acceptance before publishing. Homebrew begins with
   stable releases because Fern skips prerelease formula updates.

The generated provider-resource wire tests initially failed because the spec
marked reserved `Content-Type` as a required standalone header while Fern's IR
removed it. `fern/cli-overrides.yml` makes it optional without editing the shared
SDK files or generated Rust. Keep that exception visible when updating inputs.
