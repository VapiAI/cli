# @vapi-ai/cli

The native Vapi CLI, generated with Fern. Requires Node.js 22+.

After the beta is published:

```bash
npm install -g @vapi-ai/cli@next
export VAPI_TOKEN='<private Vapi API key>'
vapi assistants list
```

The installer downloads the exact npm version's GitHub release and verifies its
SHA-256 checksum. macOS Intel/Apple Silicon, Linux x64/ARM64, and Windows x64 are
supported. Windows requires tar.exe (Windows 10 1803+).

Version 1 changes command names, credentials, and output and removes legacy
project scaffolding and webhook forwarding. See the
[migration guide](https://github.com/VapiAI/cli/blob/main/docs/MIGRATION.md).
