# Distribution Setup Guide

This guide will help you publish the Vapi CLI to all major package managers.

## 🚀 Quick Setup Overview

1. **npm** ✅ Ready (just need to publish)
2. **Homebrew** 🛠️ Need to create tap repository
3. **Scoop** 🛠️ Need to create bucket repository
4. **Docker** 🛠️ Optional (can enable easily)

---

## 1. 📦 npm Distribution

### Setup (✅ Already Done)

- Package configuration: `npm-wrapper/package.json`
- Install script: `npm-wrapper/install.js`
- Documentation: `npm-wrapper/README.md`

### Publishing Steps:

```bash
cd npm-wrapper
npm login  # Login to npm with @vapi organization access
npm publish --access public
```

### Users can then install with:

```bash
npm install -g @vapi/cli
```

---

## 2. 🍺 Homebrew Distribution

### Setup Required:

1. **Create the tap repository:**

   ```bash
   # Create a new repository: VapiAI/homebrew-tap
   # Initialize it with a README
   ```

2. **Set up repository structure:**

   ```bash
   mkdir -p Formula
   # GoReleaser will auto-create the formula file
   ```

3. **Add GitHub token permissions:**
   - Go to GitHub Settings → Personal Access Tokens
   - Create token with `repo` and `write:packages` permissions
   - Add as repository secret: `GITHUB_TOKEN`

### Users can then install with:

```bash
brew tap vapi/tap
brew install vapi-cli
```

---

## 3. 🪣 Scoop Distribution (Windows)

### Setup Required:

1. **Create the bucket repository:**

   ```bash
   # Create a new repository: VapiAI/scoop-bucket
   # Initialize it with a README
   ```

2. **Set up repository structure:**
   ```bash
   mkdir -p bucket
   # GoReleaser will auto-create the manifest file
   ```

### Users can then install with:

```powershell
scoop bucket add vapi https://github.com/VapiAI/scoop-bucket
scoop install vapi-cli
```

---

## 4. 🐳 Docker Distribution (Optional)

### To Enable:

Uncomment the `dockers` section in `.goreleaser.yaml` and set `skip_push: false`

### Users can then use with:

```bash
docker run -it ghcr.io/vapiai/cli:latest
```

---

## 🔧 Automated Publishing Process

Once repositories are created, every release automatically:

1. **Builds** cross-platform binaries
2. **Publishes** to GitHub Releases
3. **Updates** Homebrew formula
4. **Updates** Scoop manifest
5. **Builds** Docker images (if enabled)

### Triggering a Release:

```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

---

## 📋 Repository Creation Checklist

### VapiAI/homebrew-tap

- [ ] Create repository
- [ ] Add `Formula/` directory
- [ ] Set repository as public
- [ ] Add description: "Homebrew formulae for Vapi CLI"

### VapiAI/scoop-bucket

- [ ] Create repository
- [ ] Add `bucket/` directory
- [ ] Set repository as public
- [ ] Add description: "Scoop bucket for Vapi CLI"

---

## 🎯 Final Installation Methods

Once everything is set up, users can install via:

```bash
# npm (cross-platform)
npm install -g @vapi/cli

# Homebrew (macOS/Linux)
brew tap vapi/tap && brew install vapi-cli

# Scoop (Windows)
scoop bucket add vapi https://github.com/VapiAI/scoop-bucket
scoop install vapi-cli

# Direct download (any platform)
curl -sSL https://github.com/VapiAI/cli/releases/download/v0.0.3/vapi_$(uname -s)_$(uname -m).tar.gz | tar xz

# Docker (optional)
docker run -it ghcr.io/vapiai/cli:latest
```

---

## 🔍 Verification

After publishing, test installations:

```bash
# Test npm
npm install -g @vapi/cli
vapi --version

# Test Homebrew (on macOS)
brew install vapi/tap/vapi-cli
vapi --version

# Test Scoop (on Windows)
scoop install vapi/vapi-cli
vapi --version
```

---

## 📈 Distribution Analytics

Track adoption through:

- npm download stats
- GitHub release download counts
- Homebrew analytics: `brew info vapi-cli`
- Package manager metrics

---

## 🚨 Troubleshooting

### Common Issues:

1. **Permission errors**: Ensure GitHub token has correct permissions
2. **Repository not found**: Verify repository names match GoReleaser config
3. **Formula/manifest issues**: Check GoReleaser logs for syntax errors

### Support:

- Check [GoReleaser documentation](https://goreleaser.com)
- Review GitHub Actions logs for build issues
- Monitor repository issues for user feedback
