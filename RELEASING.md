# Releasing

The replacement targets `1.0.0-beta.1`, a breaking major over the Go CLI's 0.2.x
releases. Opening or merging this draft must not be treated as release readiness.
The agreed direction is **OAuth authentication and Fern-managed publishing**.
There is no repository-owned npm/release workflow fallback.

## OAuth beta setup and release gate

The inspection build now uses Fern 5.202.0 / CLI generator 0.59.1 to generate
public-client authorization-code login with PKCE and automatic refresh.
`auth-schemes`, `api.auth`, and the matching OAuth OpenAPI security scheme are
configured together; no handwritten Rust auth or embedded client secret is used.

The staging WorkOS application was created through the Applications API:

| Setting | Value |
| --- | --- |
| Application | Vapi CLI beta (`connect_app_01M721J47TMJ2T5S0YFMGXS3SS`) |
| Public client ID | `client_01M721J47S27BS1GSXK1780TSR` |
| Issuer | `https://nourishing-mulberry-12-staging.authkit.app` |
| Authorize / token | Issuer plus `/oauth2/authorize` / `/oauth2/token` |
| Requested scopes | `openid offline_access` (standard OIDC scopes, not application permission slugs) |
| API resource | `https://api-mcp-poc.vapi.local` (the existing staging API audience) |
| Callbacks | `http://127.0.0.1:8484/callback`, `:8485/callback`, `:8486/callback` |
| Default API | `https://staging-api.vapi.ai` |

WorkOS reports first-party and PKCE enabled. An authorization probe redirects to
Vapi's staging `/oauth/continue`; invalid-code token probes without a client
secret return `invalid_grant`. These checks do not establish a completed user
login or successful API access.

The [companion API/gateway draft](https://app.graphite.com/github/pr/VapiAI/vapi/26299)
adds explicit acceptance of direct tokens for this
client through `WORKOS_CONNECT_CLI_CLIENT_ID`; the existing MCP exchange actor
check remains for delegated tokens. That companion change must deploy before
live acceptance. Leaving the CLI client setting empty disables direct tokens.
Do not copy the MCP exchange client's secret into the CLI.

Vapito validation was attempted, but its controller could not check out the
private `VapiAI/vapi` repository ("Repository not found"). Restore that source
access before treating live API acceptance as verified.

Before release, validate a real user's login, org selection, authorized read,
automatic refresh, logout, wrong audience/client, permission denial and org
isolation. Then provision and verify the production public client and audience,
and promote the issuer, client ID, resource parameters and default API server
together. Staging credentials must not become the production CLI configuration.
API-key login is not a release substitute. Track evidence on
[PAL-1311](https://linear.app/vapi/issue/PAL-1311/replace-the-vapi-cli-with-a-fern-generated-oauth-beta).

## Fern-managed distribution

The publishing preview succeeded after Fern organization access and the Fern
GitHub App installation for `VapiAI/cli` were approved. The checked-in output is
unchanged from Fern CLI 5.202.0 / CLI generator 0.59.1:

| File | Behavior |
| --- | --- |
| `.github/workflows/ci.yml` | Checks builds on pushes; version tags build and publish five npm platform packages, then the launcher |
| `.github/workflows/release.yml` | Plans releases on PRs; version tags build archives/installers, create or update the GitHub Release, and update Homebrew for stable releases |
| `dist-workspace.toml` | cargo-dist targets, installers, Homebrew tap and formula |

The npm launcher and platform package manifests are generated inside `ci.yml`
at release time. They are not separate checked-in npm package files. npm Linux
packages use musl; GitHub Release artifacts include both GNU and musl builds.

`fern/generators.yml` has two groups: `cli` for local runtime output and
`publish` for npm/GitHub delivery in pull-request mode. All preview commands
below avoid remote PR creation, package publication and release tags:

```bash
make preview-publish  # inspect full output under generated-publish/
make generate-publish # adopt runtime and publishing files from an isolated preview
python3 scripts/generate.py --publishing --check
```

Keep `verify.yml` for repository checks. Keep Fern's generated workflows in
sync through `make generate-publish`; do not replace them with another custom
publisher. Runtime files from the publishing preview match local generation.

### Release gates

1. Complete live OAuth acceptance and production promotion above.
2. Resolve the generated publishing credential boundary before enabling npm
   trusted publishing. In `ci.yml`, the platform publish jobs grant
   `id-token: write` during Rust compilation; they also execute third-party
   setup actions and, for stable releases, `npx semver`. This does not preserve
   PR #18's separation of build/packing from a minimal privileged publisher.
   The generated release workflow also grants write permissions to planning
   and build jobs. These files are checked in for review, not release approval.
3. Bootstrap the five npm platform packages if still absent, then configure
   trusted publishing for all six packages. The public registry currently has
   `@vapi-ai/cli` at `latest: 0.2.1`; the five platform names returned 404.
   Their names are `@vapi-ai/cli-{linux-x64,linux-arm64,darwin-x64,darwin-arm64,win32-x64}`.
   Use owner `VapiAI`, repository `cli`, workflow `ci.yml`, and no environment
   for the workflow as currently generated. Recheck these values after the
   credential-boundary fix. Creating packages is a publishing action; this
   preview did not do it.
4. Configure `HOMEBREW_TAP_TOKEN` with tap write access, or configure Fern's
   GitHub App distribution option and regenerate. Reconcile the tap's old
   formula/cask entries with `Formula/vapi-cli.rb` before the stable rollout.
5. Validate release builds on the target runners. npm and cargo-dist build
   independently; cargo-dist enables `rustls` explicitly while the npm build
   uses target defaults, so matching profile names alone do not establish
   byte-identical artifacts. macOS signing in `release.yml` requires all three
   `CODESIGN_*` secrets; without them those artifacts stay unsigned.

Once these gates are satisfied, publish by pushing a version tag matching
`Cargo.toml`, such as `v1.0.0-beta.1`. Merging alone does not publish. The actual
generated workflows listen for tag pushes; creating a GitHub release can also
create that tag, but no `release: published` event is configured.

The beta's npm dist-tag is `beta`, so it does not replace `latest`. Stable
backports use `backport`; other prereleases use their first identifier.
Homebrew skips prereleases. No package, tag or GitHub release was published
while preparing this draft.

See [Fern publishing](https://buildwithfern.com/learn/cli-generator/get-started/publishing)
and [Fern OAuth support](https://buildwithfern.com/learn/cli-generator/get-started/authentication).
