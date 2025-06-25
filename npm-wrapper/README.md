# Vapi CLI

Voice AI for developers - The official Vapi CLI

## Installation

```bash
npm install -g @vapi/cli
# or
yarn global add @vapi/cli
# or
pnpm add -g @vapi/cli
```

## Quick Start

1. **Install the CLI**:

   ```bash
   npm install -g @vapi/cli
   ```

2. **Login to Vapi**:

   ```bash
   vapi login
   ```

3. **Initialize a project**:

   ```bash
   vapi init
   ```

4. **List your assistants**:
   ```bash
   vapi assistants list
   ```

## Features

- 🔐 **Browser-based authentication** - Secure login via your Vapi dashboard
- 🚀 **Quick project setup** - Initialize Vapi in any JavaScript/TypeScript project
- 🤖 **Assistant management** - List, create, and manage your voice AI assistants
- 📞 **Call management** - View and manage your calls
- 🔧 **Framework detection** - Automatically detects React, Vue, Next.js, and more

## Commands

```bash
vapi --help              # Show help
vapi login               # Authenticate with Vapi
vapi init                # Initialize Vapi in your project
vapi assistants list     # List all assistants
vapi assistants get <id> # Get assistant details
vapi calls list          # List recent calls
vapi config get          # View configuration
```

## Documentation

For detailed documentation, visit [docs.vapi.ai](https://docs.vapi.ai)

## Support

- Documentation: [docs.vapi.ai](https://docs.vapi.ai)
- GitHub Issues: [github.com/VapiAI/cli](https://github.com/VapiAI/cli/issues)
- Discord: [discord.gg/vapi](https://discord.gg/vapi)

## License

MIT © [Vapi, Inc.](https://vapi.ai)
