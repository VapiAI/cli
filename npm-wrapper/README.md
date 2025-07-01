# Vapi CLI

Voice AI for developers - Official CLI tool for the Vapi platform.

## Installation

```bash
npm install -g @vapi-ai/cli
```

## Usage

```bash
# Get help
vapi --help

# Login to your Vapi account
vapi login

# List your assistants
vapi assistant list

# List your calls
vapi call list

# Start webhook listener for development
vapi listen --forward-to localhost:3000/webhook
```

## Features

- **Assistant Management** - Create, update, and manage AI assistants
- **Call Management** - Handle inbound/outbound calls
- **Phone Number Management** - Manage your Vapi phone numbers
- **Chat Management** - Text-based conversations with assistants
- **Webhook Development** - Local webhook forwarding for testing
- **Tool Management** - Custom functions and integrations
- **Logs & Debugging** - Comprehensive logging and troubleshooting

## Documentation

Visit [docs.vapi.ai](https://docs.vapi.ai) for complete documentation.

## Support

- [GitHub Issues](https://github.com/VapiAI/cli/issues)
- [Discord Community](https://discord.gg/vapi)
- [Documentation](https://docs.vapi.ai)

---

Built with ❤️ by the [Vapi](https://vapi.ai) team.
