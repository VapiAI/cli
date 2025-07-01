# Vapi CLI

The official command-line interface for [Vapi](https://vapi.ai) - Voice AI for developers.

## Features

- 🔐 **Browser-based Authentication** - Secure OAuth-style login flow
- 🤖 **Assistant Management** - List, create, update, and delete voice assistants
- 💬 **Chat Management** - Text-based conversations and chat history
- 📞 **Enhanced Call Management** - Full call lifecycle control and monitoring
- 📱 **Phone Number Management** - Purchase, configure, and manage phone numbers
- 🔄 **Workflow Management** - Manage visual conversation flows and branching logic
- 📣 **Campaign Management** - Create and manage AI phone call campaigns at scale
- 🛠️ **Tool Management** - Custom functions and API integrations
- 🔗 **Webhook Management** - Configure and manage real-time event delivery
- 🎧 **Webhook Testing** - Local webhook forwarding without ngrok
- 📋 **Logs & Debugging** - System logs, call logs, and error tracking
- 🔧 **Project Integration** - Auto-detect and integrate with existing projects
- 🚀 **Framework Support** - React, Vue, Angular, Next.js, Node.js, Python, Go, and more
- 📦 **SDK Installation** - Automatic SDK setup for your project type
- 🎨 **Code Generation** - Generate components, hooks, and examples
- ⬆️ **Auto-Updates** - Keep your CLI up-to-date with the latest features

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/VapiAI/cli.git
cd cli

# Install dependencies
make deps

# Build the CLI
make build

# Install to ~/.local/bin
make install
```

### Binary Releases

Coming soon: Pre-built binaries for macOS, Linux, and Windows.

## Development Requirements

- **Go 1.21+** - [Install Go](https://golang.org/doc/install)
- **golangci-lint** - For code linting

  ```bash
  # macOS
  brew install golangci-lint

  # Linux/Windows
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
  ```

## Usage

### Authentication

First, authenticate with your Vapi account:

```bash
vapi login
```

This will open your browser for secure authentication. Your API key will be saved locally.

### Assistant Management

```bash
# List all assistants
vapi assistant list

# Get assistant details
vapi assistant get <assistant-id>

# Create a new assistant (interactive)
vapi assistant create

# Delete an assistant
vapi assistant delete <assistant-id>
```

### Workflow Management

```bash
# List all workflows
vapi workflow list

# Get workflow details
vapi workflow get <workflow-id>

# Create a new workflow (basic)
vapi workflow create

# Delete a workflow
vapi workflow delete <workflow-id>
```

**Note**: For visual workflow building with nodes and edges, use the [Vapi Dashboard](https://dashboard.vapi.ai/workflows).

### Campaign Management

```bash
# List all campaigns
vapi campaign list

# Get campaign details
vapi campaign get <campaign-id>

# Create a new campaign
vapi campaign create

# Update/end a campaign
vapi campaign update <campaign-id>

# Delete a campaign
vapi campaign delete <campaign-id>
```

**Note**: For advanced campaign features (customer lists, scheduling), use the [Vapi Dashboard](https://dashboard.vapi.ai).

### Project Integration

Initialize Vapi in your existing project:

```bash
# Auto-detect project type and set up Vapi
vapi init

# Initialize in a specific directory
vapi init /path/to/project
```

The `init` command will:

- Detect your project framework/language
- Install the appropriate Vapi SDK
- Generate example code and components
- Create environment configuration templates

### Configuration

```bash
# View current configuration
vapi config get

# Set configuration values
vapi config set <key> <value>

# List all configuration options
vapi config list
```

### Chat Management

Manage text-based chat conversations with Vapi assistants:

```bash
# List all chat conversations
vapi chat list

# Get chat conversation details
vapi chat get <chat-id>

# Create a new chat (guided setup)
vapi chat create

# Continue an existing chat conversation
vapi chat continue <chat-id> "Your message here"

# Delete a chat conversation
vapi chat delete <chat-id>
```

### Phone Number Management

Manage your Vapi phone numbers for calls:

```bash
# List all phone numbers
vapi phone list

# Get phone number details
vapi phone get <phone-number-id>

# Purchase a new phone number (guided)
vapi phone create

# Update phone number configuration
vapi phone update <phone-number-id>

# Release a phone number
vapi phone delete <phone-number-id>
```

### Enhanced Call Management

Enhanced call operations and monitoring:

```bash
# List all calls
vapi call list

# Get call details
vapi call get <call-id>

# Create a new call (guided)
vapi call create

# Update a call in progress
vapi call update <call-id>

# End an active call
vapi call end <call-id>
```

### Logs and Debugging

View system logs for debugging and monitoring:

```bash
# List recent system logs
vapi logs list

# View call-specific logs
vapi logs calls [call-id]

# View recent error logs
vapi logs errors

# View webhook delivery logs
vapi logs webhooks
```

### Tool Management

Manage custom tools and functions that connect your voice agents to external APIs:

```bash
# List all tools
vapi tool list

# Get tool details
vapi tool get <tool-id>

# Create a new tool (guided)
vapi tool create

# Update tool configuration
vapi tool update <tool-id>

# Delete a tool
vapi tool delete <tool-id>

# Test a tool with sample input
vapi tool test <tool-id>

# List available tool types
vapi tool types
```

### Webhook Management

Manage webhook endpoints and configurations for real-time event delivery:

```bash
# List all webhook endpoints
vapi webhook list

# Get webhook details
vapi webhook get <webhook-id>

# Create a new webhook endpoint
vapi webhook create [url]

# Update webhook configuration
vapi webhook update <webhook-id>

# Delete a webhook endpoint
vapi webhook delete <webhook-id>

# Test a webhook endpoint
vapi webhook test <webhook-id>

# List available webhook event types
vapi webhook events
```

### Webhook Testing

Test your webhook integrations locally without needing ngrok or other tunneling tools:

```bash
# Forward webhooks to your local development server
vapi listen --forward-to localhost:3000/webhook

# Use a different port for the webhook listener
vapi listen --forward-to localhost:8080/api/webhooks --port 4242

# Skip TLS verification (for development only)
vapi listen --forward-to localhost:3000/webhook --skip-verify
```

The `listen` command will:

- Start a local webhook server (default port 4242)
- Forward all incoming Vapi webhooks to your specified endpoint
- Display webhook events in real-time for debugging
- Add helpful headers to identify forwarded requests

### Staying Updated

Keep your CLI up-to-date with the latest features and bug fixes:

```bash
# Check for available updates
vapi update check

# Update to the latest version
vapi update
```

The CLI will automatically check for updates periodically and notify you when a new version is available.

## Project Structure

```
cli/
├── cmd/                    # Command implementations
│   ├── root.go            # Main CLI setup
│   ├── assistant.go       # Assistant commands
│   ├── workflow.go        # Workflow commands
│   ├── campaign.go        # Campaign commands
│   ├── call.go           # Call commands
│   ├── config.go         # Configuration commands
│   ├── init.go           # Project initialization
│   └── login.go          # Authentication
├── pkg/                   # Core packages
│   ├── auth/             # Authentication logic
│   ├── client/           # Vapi API client
│   ├── config/           # Configuration management
│   ├── integrations/     # Framework integrations
│   └── output/           # Output formatting
├── build/                # Build artifacts (git-ignored)
├── main.go              # Entry point
├── Makefile             # Build automation
└── README.md            # This file
```

## Development

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run without building
go run main.go
```

### Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage
```

### Code Quality

```bash
# Run linters
make lint

# Format code
go fmt ./...
```

### Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Configuration

The CLI stores configuration in `~/.vapi-cli.yaml`. You can also use environment variables:

- `VAPI_API_KEY` - Your Vapi API key
- `VAPI_BASE_URL` - API base URL (for development)

## Supported Frameworks

### Frontend

- React (Create React App, Vite)
- Vue.js
- Angular
- Svelte
- Next.js
- Nuxt.js
- Remix
- Vanilla JavaScript

### Mobile

- React Native
- Flutter

### Backend

- Node.js/TypeScript
- Python
- Go
- Ruby
- Java
- C#/.NET

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Support

- 📚 [Documentation](https://docs.vapi.ai)
- 💬 [Discord Community](https://discord.gg/vapi)
- 🐛 [Issue Tracker](https://github.com/VapiAI/cli/issues)

---

Built with ❤️ by the Vapi team

## Version Management

The Vapi CLI uses a simple and discoverable version management system:

### Current Version

The current version is stored in the `VERSION` file at the project root. This makes it easy to find and update.

### Managing Versions

#### Using Make (Recommended)

```bash
# Show current version
make version

# Set a specific version
make version-set VERSION=1.2.3

# Bump versions automatically
make version-bump-patch    # 1.2.3 -> 1.2.4
make version-bump-minor    # 1.2.3 -> 1.3.0
make version-bump-major    # 1.2.3 -> 2.0.0
```

#### Using the Script Directly

```bash
# Show current version
./scripts/version.sh get

# Set a specific version
./scripts/version.sh set 1.2.3

# Bump versions
./scripts/version.sh bump patch
./scripts/version.sh bump minor
./scripts/version.sh bump major
```

### How It Works

1. **Development**: The CLI reads the version from the `VERSION` file
2. **Release Builds**: GoReleaser overrides the version using git tags and ldflags
3. **Priority**: Build-time version (from releases) takes priority over the VERSION file

This approach provides:

- ✅ Easy version discovery (just check the `VERSION` file)
- ✅ Automated version bumping with semantic versioning
- ✅ Consistent versioning across development and releases
- ✅ No need to manually edit code files
