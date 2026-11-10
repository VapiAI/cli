# @vapi-ai/mcp-docs-server (deprecated)

> **Deprecated and no longer maintained.** This legacy local documentation server
> has been replaced by the [Vapi MCP server](https://github.com/VapiAI/mcp-server)
> and the [Fern documentation MCP](https://docs.vapi.ai/_mcp/server).
> All published versions are deprecated on npm. This directory is retained for
> historical reference and is excluded from the CLI build and installation.

## Use the current MCP server

Connect your agent to **`https://mcp.vapi.ai/mcp`** using the
[Vapi MCP setup guide](https://docs.vapi.ai/sdk/mcp-server).
For local stdio use, the current package is **`@vapi-ai/mcp-server`**.
Follow its [README](https://github.com/VapiAI/mcp-server#readme) for authentication
and configuration.

The current MCP package includes **`ask_vapi_docs`**, which forwards documentation
questions to Fern and returns documentation passages and source URLs. It also
provides tools for managing Vapi assistants, calls, phone numbers, and other
resources. Check your connection's tool list for availability.

## Documentation-only access

For documentation without Vapi account access, connect directly to Fern's
**`https://docs.vapi.ai/_mcp/server`** endpoint. It exposes `searchDocs` and does
not require a Vapi API key. For clients that need a local stdio adapter:

```json
{
  "mcpServers": {
    "vapi-docs": {
      "command": "npx",
      "args": ["-y", "mcp-remote", "https://docs.vapi.ai/_mcp/server"]
    }
  }
}
```

## Migrate an existing installation

Remove any MCP client entry that launches `@vapi-ai/mcp-docs-server`, including
entries written by older versions of `vapi mcp setup`. Add one of the supported
connections above, then restart your client. Existing configurations are not
changed automatically.

## License

[MIT](./LICENSE).
