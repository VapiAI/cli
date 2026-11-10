# Releasing

The replacement targets `1.0.0-beta.1`, a breaking major over the Go CLI's 0.2.x
releases. Opening or merging this draft must not be treated as release readiness.
The agreed direction is **OAuth authentication and Fern-managed publishing**.
There is no repository-owned npm/release workflow fallback.

## OAuth prerequisite

Fern supports authorization-code login with PKCE and automatic refresh through
`auth-schemes` in `generators.yml`. The CLI needs a public client registration,
issuer/authorization/token endpoints, API audience/scopes, and registered
loopback callback behavior. It must not embed a client secret.

Vapi's WorkOS Connect work currently delegates MCP access through server-side
token exchange. The API auth provider binds those exchanged tokens to its
configured exchange client. A native CLI's public-client token cannot be assumed
to satisfy that contract. Confirm direct CLI token issuance and API/gateway
acceptance with the OAuth work, then encode the actual contract and regenerate.
Do not copy the MCP server's exchange credentials into the CLI.

The current generated inspection build still has bearer authentication. It is
not an implementation of this prerequisite, and API-key login is not an accepted
release substitute. Validate browser login, org selection, token refresh,
logout, wrong audience, permission denial, and org isolation before releasing.

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
