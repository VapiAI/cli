# @vapi-ai/mcp-docs-server

[![npm version](https://badge.fury.io/js/@vapi-ai%2Fmcp-docs-server.svg)](https://badge.fury.io/js/@vapi-ai%2Fmcp-docs-server)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**Turn your IDE into a Vapi expert!** 🎯

This is a Model Context Protocol (MCP) server that provides direct access to live Vapi documentation from docs.vapi.ai. It fetches real documentation, examples, and API references, integrating seamlessly with AI-powered IDEs like Cursor, Windsurf, and VSCode to give you instant access to accurate Vapi expertise.

## ✨ Features

- 📚 **Live Documentation Access** - Fetches real docs from docs.vapi.ai
- 💻 **Real Examples** - Links to actual working code in Vapi docs
- 🔧 **Current API Reference** - Always up-to-date API documentation
- 📖 **Official Guides** - Step-by-step tutorials from Vapi
- 📋 **Real Changelog** - Latest updates from actual Vapi changelog
- 🔍 **Smart Search** - Search across all live Vapi documentation
- 🎯 **Always Fresh** - Documentation cached and auto-refreshed

## 🚀 Quick Start

### Installation

```bash
npm install -g @vapi-ai/mcp-docs-server
```

### IDE Configuration

The server is automatically configured when you run `vapi mcp setup` from the [Vapi CLI](https://github.com/VapiAI/cli).

#### Manual Configuration

<details>
<summary>Cursor IDE</summary>

Add to your `.cursor/mcp.json`:

```json
{
  "mcpServers": {
    "vapi": {
      "command": "vapi-mcp-docs-server",
      "args": []
    }
  }
}
```

</details>

<details>
<summary>Windsurf IDE</summary>

Add to your `.codeium/windsurf/mcp_config.json`:

```json
{
  "mcpServers": {
    "vapi": {
      "command": "vapi-mcp-docs-server",
      "args": []
    }
  }
}
```

</details>

<details>
<summary>VSCode</summary>

Add to your `.vscode/mcp.json`:

```json
{
  "servers": {
    "vapi": {
      "command": "vapi-mcp-docs-server",
      "args": []
    }
  }
}
```

</details>

## 🛠️ Available Tools

Once configured, your IDE's AI assistant will have access to these tools:

### `search_documentation`

Search Vapi documentation for specific topics, features, or concepts.

```
Search for "phone calls" → Get comprehensive guides on making calls
Search for "webhooks" → Learn about real-time event handling
Search for "voice settings" → Configure voice providers and parameters
```

### `get_examples`

Get code examples for specific Vapi features or use cases.

```
Get examples for "assistants" in TypeScript
Get examples for "function calling" in Python
Get examples for "voice configuration" in all languages
```

### `get_guides`

Access step-by-step guides for implementing Vapi features.

```
Get guides for "getting started"
Get guides for "phone calls"
Get guides for "tools and functions"
```

### `get_api_reference`

Get detailed API reference information for Vapi endpoints.

```
Get API reference for "assistants"
Get API reference for "calls" with POST method
Get API reference for "webhooks" with examples
```

### `get_changelog`

Get recent changes, updates, and new features in Vapi.

```
Get latest 5 changelog entries
Get changelog for version "1.8.0"
Get only "features" type changes
```

## 🎯 Example Usage

Once configured, you can ask your IDE's AI assistant questions like:

- **"How do I create a voice assistant with Vapi?"**
- **"Show me examples of making phone calls"**
- **"What are the latest Vapi features?"**
- **"How do I set up webhooks for function calling?"**
- **"Give me the API reference for creating assistants"**

The AI will automatically use the MCP tools to provide accurate, up-to-date information!

## 🔧 Development

### Prerequisites

- Node.js 18+
- npm or yarn

### Setup

```bash
# Clone the repository
git clone https://github.com/VapiAI/mcp-docs-server.git
cd mcp-docs-server

# Install dependencies
npm install

# Build the project
npm run build

# Start development server
npm run dev
```

### Testing

```bash
# Run tests
npm test

# Lint code
npm run lint

# Fix linting issues
npm run lint:fix
```

### Project Structure

```
src/
├── index.ts              # Main entry point
├── server.ts             # MCP server implementation
├── tools/                # MCP tools
│   ├── search.ts         # Documentation search
│   ├── examples.ts       # Code examples
│   ├── guides.ts         # Step-by-step guides
│   ├── api-reference.ts  # API documentation
│   └── changelog.ts      # Version history
├── resources/            # MCP resources
│   └── documentation.ts  # Resource handlers
└── utils/                # Utilities
    └── documentation-data.ts # Data and types
```

## 📖 How It Works

This MCP server implements the [Model Context Protocol](https://spec.modelcontextprotocol.io/) to provide structured access to Vapi's knowledge base. When your IDE's AI assistant needs information about Vapi:

1. **Query Processing** - The AI identifies what information is needed
2. **Tool Selection** - Chooses the appropriate MCP tool
3. **Content Retrieval** - Fetches relevant documentation/examples
4. **Response Generation** - Provides accurate, contextual answers

## 🤝 Contributing

We welcome contributions! Here's how you can help:

1. **Add Content** - Expand documentation, examples, or guides
2. **Improve Search** - Enhance search algorithms and relevance
3. **Fix Bugs** - Report and fix issues
4. **Add Features** - Propose new tools or capabilities

### Contributing Guidelines

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests if applicable
5. Run linting and tests (`npm run lint && npm test`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Links

- **[Vapi Website](https://vapi.ai)** - Voice AI platform
- **[Vapi Documentation](https://docs.vapi.ai)** - Complete docs
- **[Vapi CLI](https://github.com/VapiAI/cli)** - Command-line tool
- **[Vapi Discord](https://discord.gg/vapi)** - Community support
- **[Model Context Protocol](https://spec.modelcontextprotocol.io/)** - MCP specification

## 🆘 Support

- **Issues** - [GitHub Issues](https://github.com/VapiAI/mcp-docs-server/issues)
- **Discord** - [Vapi Community](https://discord.gg/vapi)
- **Email** - [support@vapi.ai](mailto:support@vapi.ai)

---

Made with ❤️ by the [Vapi](https://vapi.ai) team
