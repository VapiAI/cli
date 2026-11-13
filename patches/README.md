# Fern output patches

These files retain reviewable copies of the CLI runtime fixes. Fern Replay's
native `.fern/replay.lock` tracks the equivalent customization commit above a
pure `[fern-generated]` baseline and carries it into centrally generated PRs.
Generation is configured in `VapiAI/docs`, alongside the server SDKs; no local
patch-application or API-vendoring script is required.

Review conflicts when upgrading Fern. Remove a customization from Replay once
the upstream generator supplies equivalent behavior. CI checks correctness and
performance rather than requiring a particular implementation.

## Schema reference expansion (CLI generator 0.59.1)

`fern-cli-0.59.1-schema-refs.patch` fixes unbounded memory growth in scoped
schema discovery, including `vapi assistants list --schema` and
`vapi assistants create --schema`. The original renderer recursively inlines
every named reference. Its cycle guard prevents infinite recursion along one
path, but shared types are expanded again in every sibling branch. A small
schema graph can therefore produce an exponentially large output tree.

The patched renderer keeps the response root and input flag metadata inline.
Nested named types use local `$ref` URI fragments into a shared `$defs` table
at the operation document's root. An iterative worklist renders each reachable
named definition once, preserving recursive types without truncating them.
Consumers must resolve references against the complete returned document.
Root and group command catalogs retain their existing format.

The patch includes regressions for a shared acyclic graph, mutually recursive
request and response types, self-references, and escaping reference names. Run:

```bash
cargo test --locked --lib openapi::help::tests
cargo build --locked
./target/debug/vapi assistants list --schema
./target/debug/vapi assistants create --schema
```

Schema discovery requires no login or API requests.

`python3 scripts/check-performance.py` checks the release CLI's schema discovery
time, OS-reported peak resident memory, and output size. Build with
`cargo build --release --locked` first. It runs in Linux CI and supports Linux
and macOS locally. These checks continue to apply after the patch is replaced
by an upstream fix.

## Vapi client attribution (CLI generator 0.59.1)

`fern-cli-0.59.1-vapi-client.patch` adds Vapi's first-party attribution to the
HTTP client for the `vapi` binary:

```text
X-Vapi-Client: cli/<version>; command=assistants_list[; ci]
User-Agent: vapi-cli/<version>
```

The version comes from `CARGO_PKG_VERSION`, keeping it synchronized with
`--version` without a second version constant. Command tokens are joined from
clap's resolved subcommand names with underscores; argument values are excluded.
The token is omitted if it does not meet the API's lowercase ASCII token and
64-character limit. A truthy `CI` variable (`1`, `true`, `yes`, or
`on`, case insensitive) adds the bare `ci` flag.

Fern still composes the User-Agent and optional downstream suffix. No
`X-Vapi-Source` header is added. This customization derives the header's version and invocation metadata at runtime.
Replace it with generator configuration when that supports these values.

`scripts/smoke-test.py` checks the actual GET, POST, and DELETE request headers,
retry attribution, CI detection, and suffix precedence
against a local HTTP server.
