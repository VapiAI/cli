# Releasing

The first replacement is `1.0.0-beta.1`, a breaking major over the Go CLI's
0.2.x releases. Opening or merging the migration PR does not publish anything.

## Validate

1. Review `docs/MIGRATION.md`, especially removed workflows, authentication,
   command names, and noninteractive deletion.
2. Set the same version in `Cargo.toml` and `npm-wrapper/package.json`, regenerate
   with `make generate`, and run `make test`.
3. Require CI on Linux, macOS, and Windows to pass. Release binaries also perform
   version and schema smoke checks on all five target platforms.
4. Confirm npm trusted publishing authorizes `VapiAI/cli` and the calling
   `release.yml` workflow for the existing `@vapi-ai/cli` package. This changes
   the previous workflow identity: npm validates the caller when publishing
   through a reusable workflow. See [npm trusted publishing](https://docs.npmjs.com/trusted-publishers/).
   No npm token is passed into build or package jobs.

## Publish

After approval and merge, push a matching tag (for example `v1.0.0-beta.1`).
The release workflow builds all native binaries before creating a GitHub
release with archives and `checksums.txt`. It then calls the npm workflow
explicitly, since releases created with `GITHUB_TOKEN` do not trigger another
release-event workflow.

Prerelease tags create a GitHub prerelease and publish npm under `next`.
Stable versions publish under `latest`. The npm installer downloads the exact
package version's release, verifies its checksum, and installs that binary.
It never substitutes the latest stable release for a beta or pinned version.

Releases retain `cli_Darwin_*`, `cli_Linux_*`, and `cli_Windows_x86_64.zip`
archive names for the npm distribution path. Windows archives are ZIP files;
Unix archives are tar.gz files. Native archives include license notices.

The npm workflow retains separate package and publish jobs: the job with OIDC
publishing permission has no checkout, does not restore caches, and publishes
the prebuilt tarball with lifecycle scripts disabled. Users run the installer
only when they install the published package.

Do not create a stable 1.0.0 until the beta's UX decisions are accepted. The old
Docker images, Go installers, and Go self-update command are retired by this
migration; npm or manual release downloads are the supported v1 channels.

## Fern-managed publishing

The local-file generator works with the current Fern login. An attempted
npm/GitHub preview failed with “You do not have permission to run this generator
for organization 'vapi'.” This PR therefore uses generated Rust plus
repository-owned release/npm workflows. No Fern npm or GitHub publish was run.

After Fern confirms publishing access, a separate change can add a publishing
group to `VapiAI/docs` with `output.location: npm`, package `@vapi-ai/cli`, OIDC,
and `github.mode: pull-request`. Compare its generated workflow with the
credential isolation here before replacing the release path. See
[Fern publishing](https://buildwithfern.com/learn/cli-generator/get-started/publishing).
