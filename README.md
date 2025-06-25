# Vapi CLI

A powerful command-line interface for interacting with the Vapi API and seamlessly integrating Vapi into existing web projects.

## Features

- 🔐 **Browser-Based Authentication**: Secure login flow that opens your browser for seamless authentication
- 🎯 **Real API Integration**: Manage assistants, calls, and other Vapi resources using the official SDK
- 🚀 **Smart Project Detection**: Automatically detects React, Vue, Angular, Svelte, Next.js, and more
- 🎨 **Interactive Setup**: Step-by-step guided integration with customizable options
- 📦 **Package Manager Detection**: Works with npm, yarn, pnpm, and bun
- 🔧 **TypeScript & Tailwind Support**: Automatic detection and appropriate code generation

## Installation

### From Source

```bash
git clone https://github.com/VapiAI/cli.git
cd cli
go build -o vapi .
```

### Using Go Install

```bash
go install github.com/VapiAI/cli@latest
```

## Quick Start

### 1. Authenticate with Vapi

```bash
# Login using browser-based authentication
vapi login
```

This will open your browser for secure authentication and save your credentials locally.

### 2. Initialize Vapi in Your Project

Navigate to your web project and run:

```bash
# Auto-detect framework and start interactive setup
vapi init

# Or specify a project path
vapi init /path/to/your/project
```

The CLI will:

- 🔍 Automatically detect your framework (React, Vue, Angular, etc.)
- 📋 Show interactive prompts for customization
- 📦 Install the Vapi Web SDK
- 🎨 Generate framework-specific components
- 📝 Create environment configuration templates

### 3. Manage Your Assistants

```bash
# List all assistants
vapi assistant list

# Get assistant details
vapi assistant get <assistant-id>

# Interactive assistant creation
vapi assistant create
```

```javascript
import { VapiButton } from "./components/vapi/VapiButton";

function App() {
  return (
    <div>
      <h1>My App with Vapi</h1>
      <VapiButton />
    </div>
  );
}
```

## Commands

### Authentication

```bash
# Login with browser-based authentication
vapi login

# Show current configuration
vapi config show
```

### API Management

```bash
# Assistants
vapi assistant list                    # List all assistants
vapi assistant get <id>                # Get assistant details
vapi assistant create                  # Interactive assistant creation
vapi assistant delete <id>             # Delete an assistant

# Calls
vapi call list                         # List recent calls
vapi call get <id>                     # Get call details
```

### Project Integration

```bash
# Auto-detect framework and initialize
vapi init

# Initialize in specific directory
vapi init /path/to/project

# Get help
vapi init --help
```

### Configuration

```bash
# Manual configuration (if not using login)
vapi config set api_key YOUR_API_KEY
vapi config get api_key
vapi config show
```

## Configuration

The Vapi CLI uses a configuration file located at:

- `./.vapi-cli.yaml` (current directory)
- `$HOME/.vapi-cli.yaml` (home directory)

### Configuration Options

```yaml
# Required: Your Vapi API key
api_key: "your-api-key-here"

# Optional: Custom Vapi server URL
base_url: "https://api.vapi.ai"

# Optional: Request timeout in seconds (default: 30)
timeout: 30
```

### Environment Variables

You can also use environment variables:

```bash
export VAPI_API_KEY="your-api-key"
export VAPI_BASE_URL="https://api.vapi.ai"
export VAPI_TIMEOUT="30"
```

## React Integration Details

### Supported Project Types

- ✅ Create React App (CRA)
- ✅ Next.js (App Router & Pages Router)
- ✅ Vite
- ✅ TypeScript projects
- ✅ Tailwind CSS projects

### Generated Components

The CLI generates the following files in your React project:

#### `useVapi` Hook

A custom React hook for managing Vapi connections:

```typescript
const { startCall, endCall, isSessionActive, isLoading, error } = useVapi({
  publicKey: "your-public-key",
  assistantId: "your-assistant-id",
});
```

#### `VapiButton` Component

A ready-to-use button component:

```jsx
<VapiButton />
<VapiButton className="custom-style">Custom Text</VapiButton>
```

#### `VapiExample` Component

A complete example showing how to use Vapi in your app.

### Environment Variables

For React projects, set these environment variables:

**Standard React Apps:**

```bash
REACT_APP_VAPI_PUBLIC_KEY=your_public_key
REACT_APP_VAPI_ASSISTANT_ID=your_assistant_id
```

**Next.js Apps:**

```bash
NEXT_PUBLIC_VAPI_PUBLIC_KEY=your_public_key
NEXT_PUBLIC_VAPI_ASSISTANT_ID=your_assistant_id
```

## Global Flags

- `--config`: Specify custom config file path
- `--api-key`: Override API key from command line
- `--help`: Show help information

## Examples

### Basic API Usage

```bash
# Configure CLI
vapi config set api_key sk-xxx

# List your assistants
vapi assistant list

# Get details of a specific assistant
vapi assistant get assistant-123

# Create a new call
vapi call create
```

### React Integration Workflow

```bash
# Navigate to your React project
cd my-react-app

# Initialize Vapi integration
vapi init react

# Follow the setup instructions
cp .env.example .env
# Edit .env with your keys
npm install
```

## Development

### Building from Source

```bash
git clone https://github.com/VapiAI/cli.git
cd cli
go mod download
go build -o vapi .
```

### Running Tests

```bash
go test ./...
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- 📖 [Vapi Documentation](https://docs.vapi.ai)
- 💬 [Discord Community](https://discord.gg/vapi)
- 🐛 [Issues](https://github.com/VapiAI/cli/issues)
- 📧 [Email Support](mailto:support@vapi.ai)
