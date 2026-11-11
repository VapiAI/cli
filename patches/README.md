# Temporary Fern fixes

`scripts/generate.py` applies these patches to the pinned generator's isolated
output before checking or copying it, for both `cli` and `publish`. It checks
each patch before application and fails if the context has changed. Review
these patches when upgrading Fern; remove a patch once the upstream generator
has the corresponding fix. CI checks schema discovery performance rather than
requiring this particular implementation.

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

`python3 scripts/check-performance.py` checks the compiled CLI's schema discovery
time, memory use, and output size. It runs in Linux CI and supports Linux and
macOS locally. These checks continue to apply after the patch is replaced by
an upstream fix.
