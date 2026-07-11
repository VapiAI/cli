# Development Guide

This guide is for developers working on the Vapi CLI or integrating with different Vapi environments.

## Environment Configuration

The Vapi CLI supports multiple environments for development and testing purposes. This functionality is designed to be **hidden from end users** but accessible to developers.

### Available Environments

| Environment   | API URL                       | Dashboard URL                       | Use Case               |
| ------------- | ----------------------------- | ----------------------------------- | ---------------------- |
| `production`  | `https://api.vapi.ai`         | `https://dashboard.vapi.ai`         | Default, end users     |
| `staging`     | `https://api.staging.vapi.ai` | `https://dashboard.staging.vapi.ai` | Pre-production testing |
| `development` | `http://localhost:3000`       | `http://localhost:3001`             | Local development      |

### Configuration Methods

#### 1. Environment Variables (Recommended for Development)

```bash
# Set environment
export VAPI_ENV=staging

# Or set URLs directly (overrides environment)
export VAPI_API_BASE_URL=http://localhost:3000
export VAPI_DASHBOARD_URL=http://localhost:3001
```

#### 2. CLI Configuration

```bash
# Hidden command for developers
vapi config env staging

# Or set directly
vapi config set environment staging
```

#### 3. Development Script (Easiest)

```bash
# Quick setup for local development
./scripts/dev-env.sh local

# Set up staging environment
./scripts/dev-env.sh setup staging

# Check current status
./scripts/dev-env.sh status

# Reset to production
./scripts/dev-env.sh reset
```

### Priority Order

The CLI determines which environment to use in this order:

1. **Direct URL overrides** (`VAPI_API_BASE_URL`, `VAPI_DASHBOARD_URL`)
2. **Environment variable** (`VAPI_ENV`)
3. **Config file** (`environment` field)
4. **Default** (production)

### Development Workflow

#### Setting Up Local Development

1. **Start your local services:**

   ```bash
   # Start API server on localhost:3000
   # Start dashboard on localhost:3001
   ```

2. **Configure CLI:**

   ```bash
   ./scripts/dev-env.sh local
   ```

3. **Verify setup:**

   ```bash
   vapi version  # Should show environment info
   vapi config get
   ```

4. **Test authentication:**
   ```bash
   vapi login  # Opens localhost:3001/auth/cli
   ```

#### Testing Against Staging

```bash
# Switch to staging
export VAPI_ENV=staging
vapi config set environment staging

# Or use the script
./scripts/dev-env.sh setup staging

# Test commands
vapi assistant list
```

#### Switching Back to Production

```bash
./scripts/dev-env.sh reset
# or
unset VAPI_ENV
vapi config set environment production
```

### Environment Detection

The CLI shows environment information when not in production:

```bash
$ vapi version
vapi version 0.0.3
  commit: abc123
  built at: 2025-01-27
  built by: dev
  go version: go1.24.4
  platform: darwin/arm64
  environment: staging          # Only shown for non-production
  api url: https://api.staging.vapi.ai
```

### For End Users

End users will **never** see environment-related functionality:

- No environment flags in help output
- Commands are hidden (`vapi config env` is hidden)
- Default behavior is always production
- No environment information shown in version (unless non-production)

### Configuration File

The CLI saves configuration to `~/.vapi-cli.yaml`:

```yaml
api_key: "your-api-key"
environment: "staging"
base_url: "https://api.staging.vapi.ai"
dashboard_url: "https://dashboard.staging.vapi.ai"
timeout: 30
```

### Environment Variables Reference

| Variable             | Description            | Example                  |
| -------------------- | ---------------------- | ------------------------ |
| `VAPI_ENV`           | Environment name       | `staging`, `development` |
| `VAPI_API_KEY`       | API key                | `vapi_abc123...`         |
| `VAPI_API_BASE_URL`  | Override API URL       | `http://localhost:3000`  |
| `VAPI_DASHBOARD_URL` | Override dashboard URL | `http://localhost:3001`  |

### Testing Authentication Flow

When testing authentication against different environments:

```bash
# Local development
export VAPI_DASHBOARD_URL=http://localhost:3001
vapi login  # Opens localhost:3001/auth/cli

# Staging
export VAPI_ENV=staging
vapi login  # Opens dashboard.staging.vapi.ai/auth/cli
```

### Troubleshooting

#### Check Current Configuration

```bash
vapi config get  # Shows all settings
./scripts/dev-env.sh status  # Comprehensive status
```

#### Reset Everything

```bash
./scripts/dev-env.sh reset
rm ~/.vapi-cli.yaml  # Nuclear option
```

#### Debug Environment Detection

```bash
# See what environment variables are set
env | grep VAPI

# Check config file
cat ~/.vapi-cli.yaml
```

### Adding New Environments

To add a new environment (e.g., `testing`):

1. Update `pkg/config/config.go`:

   ```go
   "testing": {
       Name:         "testing",
       APIBaseURL:   "https://api.testing.vapi.ai",
       DashboardURL: "https://dashboard.testing.vapi.ai",
   },
   ```

2. Update scripts and documentation as needed.

### Security Considerations

- Environment switching is intentionally hidden from end users
- API keys are masked in all output
- Local development uses localhost URLs only
- Production is always the default and safest option
