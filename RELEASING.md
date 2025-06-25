# Release Process

This document describes the release process for the Vapi CLI.

## Overview

The Vapi CLI uses:

- **Semantic Versioning** (v1.2.3)
- **GitHub Releases** for binary distribution
- **GoReleaser** for automated cross-platform builds
- **GitHub Actions** for CI/CD

## Quick Release

To create a new release:

```bash
# Make sure you're on main and up to date
git checkout main
git pull origin main

# Create and push a new tag
git tag -a v1.2.3 -m "Release v1.2.3"
git push origin v1.2.3
```

This will trigger the automated release process.

## Release Checklist

Before releasing:

- [ ] All tests pass (`make test`)
- [ ] No linting issues (`make lint`)
- [ ] README.md is up to date
- [ ] CHANGELOG.md is updated (if maintaining one)
- [ ] Version number follows semver

## Versioning Strategy

We follow [Semantic Versioning](https://semver.org/):

- **MAJOR** (v2.0.0): Breaking changes
- **MINOR** (v1.3.0): New features, backwards compatible
- **PATCH** (v1.2.4): Bug fixes, backwards compatible

### Pre-releases

For testing releases:

```bash
git tag -a v1.2.3-beta.1 -m "Beta release v1.2.3-beta.1"
```

## Distribution Channels

### 1. Direct Download (Available Now)

Users can download pre-built binaries from GitHub Releases:

- macOS (Intel & Apple Silicon)
- Linux (x86_64, ARM)
- Windows (x86_64)

### 2. Homebrew (macOS/Linux)

**Setup Required:**

1. Create `VapiAI/homebrew-tap` repository
2. Update `.goreleaser.yaml` to set `skip_upload: false`
3. Users install with:
   ```bash
   brew tap vapi/tap
   brew install vapi-cli
   ```

### 3. Scoop (Windows)

**Setup Required:**

1. Create `VapiAI/scoop-bucket` repository
2. Update `.goreleaser.yaml` to set `skip_upload: false`
3. Users install with:
   ```powershell
   scoop bucket add vapi https://github.com/VapiAI/scoop-bucket
   scoop install vapi-cli
   ```

### 4. Package Managers (Future)

Consider adding support for:

- **apt/deb** (Debian/Ubuntu)
- **yum/rpm** (RHEL/CentOS)
- **AUR** (Arch Linux)
- **Chocolatey** (Windows)
- **npm/yarn** (Cross-platform via Node.js wrapper)

### 5. Install Script

Create a universal install script:

```bash
curl -sSL https://vapi.ai/install.sh | bash
```

## What Happens During Release

When you push a tag:

1. **GitHub Actions** triggers the release workflow
2. **Tests** run on all platforms
3. **GoReleaser** builds binaries for all platforms
4. **Archives** are created with proper naming
5. **Checksums** are generated for security
6. **GitHub Release** is created with:
   - Pre-built binaries
   - Installation instructions
   - Changelog from commits
7. **Optional:** Updates Homebrew/Scoop formulas

## Manual Release (Emergency)

If automation fails:

```bash
# Install GoReleaser locally
brew install goreleaser

# Create release manually
goreleaser release --clean

# Or create a snapshot (no tag required)
goreleaser release --snapshot --clean
```

## Best Practices for CLI Distribution

### 1. **Binary Naming**

- Simple, memorable name (`vapi` not `vapi-cli`)
- Consistent across platforms

### 2. **Version Information**

```bash
vapi --version
# Output: vapi version 1.2.3 (commit: abc123, built: 2024-01-01)
```

### 3. **Auto-Updates**

Consider adding self-update capability:

```bash
vapi update
```

### 4. **Installation Methods**

Provide multiple options:

- Direct download
- Package managers
- Install script
- Docker image
- Build from source

### 5. **Platform Support**

Essential platforms:

- macOS (Intel + Apple Silicon)
- Linux (x86_64 + ARM)
- Windows (x86_64)

### 6. **Security**

- Sign binaries (code signing certificates)
- Provide checksums
- Use HTTPS for all downloads
- Consider notarization for macOS

### 7. **Documentation**

- Clear installation instructions per platform
- Troubleshooting guide
- Uninstall instructions

### 8. **Backwards Compatibility**

- Maintain config file compatibility
- Provide migration tools
- Clear deprecation warnings

## Setting Up Distribution Channels

### Homebrew Tap

1. Create `VapiAI/homebrew-tap` repository
2. GoReleaser will auto-update the formula
3. Test installation:
   ```bash
   brew tap vapi/tap
   brew install vapi-cli
   brew test vapi-cli
   ```

### NPM Package (Optional)

For Node.js users, create a wrapper:

```json
{
  "name": "@vapi/cli",
  "version": "1.2.3",
  "bin": {
    "vapi": "./bin/vapi"
  },
  "scripts": {
    "postinstall": "node install.js"
  }
}
```

### Docker Image

Already configured in `.goreleaser.yaml`. Enable when ready:

```bash
docker run -it ghcr.io/vapiai/cli:latest
```

## Release Notes Template

```markdown
## What's Changed

### ✨ Features

- Feature 1 (#123)
- Feature 2 (#124)

### 🐛 Bug Fixes

- Fix 1 (#125)
- Fix 2 (#126)

### 📚 Documentation

- Updated README (#127)

### 🔧 Maintenance

- Dependency updates (#128)

**Full Changelog**: https://github.com/VapiAI/cli/compare/v1.2.2...v1.2.3
```

## Monitoring Releases

Track adoption through:

- GitHub Release download counts
- Homebrew analytics (`brew info vapi-cli`)
- Error reporting/telemetry (with user consent)

## Troubleshooting

Common issues:

### Tag already exists

```bash
git tag -d v1.2.3
git push origin :v1.2.3
```

### GoReleaser fails

Check:

- Git is clean (`git status`)
- Tag is pushed to origin
- GitHub token has correct permissions

### Platform-specific issues

Test locally:

```bash
goreleaser build --single-target --snapshot
```
