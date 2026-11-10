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

`fern/generators.yml` has two groups:

- `cli`: local-file output used to inspect/build the generated CLI.
- `publish`: npm package `@vapi-ai/cli`, OIDC, GitHub delivery to `VapiAI/cli` in
  pull-request mode, and Homebrew formula `vapi-cli` in `VapiAI/homebrew-tap`.

Inspect publishing files without publishing:

```bash
make preview-publish
```

The latest attempt failed with “You do not have permission to run this generator
for organization 'vapi'.” Local-file generation works. Confirm the logged-in
account's organization membership and CLI publishing-generator access with Fern.
The local group is for evaluation, not an alternative publishing implementation.

Once access is enabled, inspect and adopt Fern's actual generated npm launcher,
platform packages, release workflow, installers, and Homebrew job. Keep
`.github/workflows/verify.yml` for repository checks; Fern owns publishing files.
Review credential boundaries in the generated workflow against the isolation
requirements established in PR #18 before using it.

Publishing setup also requires:

- The Fern GitHub App installed on `VapiAI/cli`.
- Existing npm platform packages or their initial setup, followed by OIDC
  trusted-publisher configuration matching Fern's generated workflow. Creating
  those packages is a separate publishing action, not part of previewing.
- Homebrew tap write credentials (the default secret is `HOMEBREW_TAP_TOKEN`,
  or Fern's GitHub App option). The tap already exists, but its old formula/cask
  entries need to be reconciled with the new generated formula before rollout.

Fern skips Homebrew updates for prereleases. The beta is evaluated through npm
and native artifacts; the first stable release updates the Homebrew formula.
Confirm npm prerelease dist-tag handling in the actual generated workflow before
creating a beta tag. Nothing in this draft has published or bootstrapped packages.

See [Fern publishing](https://buildwithfern.com/learn/cli-generator/get-started/publishing)
and [Fern OAuth support](https://buildwithfern.com/learn/cli-generator/get-started/authentication).
