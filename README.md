# Vapi CLI

The official command-line interface for [Vapi](https://vapi.ai) - Voice AI for developers.

## Features

- 🔐 **Browser-based Authentication** - Secure OAuth-style login flow
- 🤖 **Assistant Management** - List, create, update, and delete voice assistants
- 📞 **Call Management** - Monitor and control phone calls
- 🔧 **Project Integration** - Auto-detect and integrate with existing projects
- 🚀 **Framework Support** - React, Vue, Angular, Next.js, Node.js, Python, Go, and more
- 📦 **SDK Installation** - Automatic SDK setup for your project type
- 🎨 **Code Generation** - Generate components, hooks, and examples

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

## Project Structure

```
cli/
├── cmd/                    # Command implementations
│   ├── root.go            # Main CLI setup
│   ├── assistant.go       # Assistant commands
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
