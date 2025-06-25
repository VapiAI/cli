# Distribution Channel Setup Guide

This guide walks through setting up all distribution channels for the Vapi CLI.

## Prerequisites

- Admin access to VapiAI GitHub organization
- Access to vapi.ai website deployment
- npm account with publishing rights (for npm distribution)

## 1. Install Script Deployment (Manual)

The `scripts/install.sh` should be manually copied to your website when it changes:

### Manual Update Process

1. Copy `scripts/install.sh` to your website repo's `public/` directory
2. Commit and push to trigger Vercel deployment
3. Verify the script is accessible at `https://vapi.ai/install.sh`

**Why manual?** The install script rarely changes after initial setup, making automation unnecessary overhead.

### Testing the Install Script

Before deploying to production:

```bash
# Test locally
bash scripts/install.sh

# Test from a URL (after deploying to staging)
curl -sSL https://staging.vapi.ai/install.sh | bash
```

## 2. Homebrew Setup

### Step 1: Create Tap Repository

```bash
# Create new repository: VapiAI/homebrew-tap
# Initialize with README:
```

**README.md for homebrew-tap:**

````markdown
# Vapi Homebrew Tap

## Installation

```bash
brew tap vapi/tap
brew install vapi-cli
```
````

## Available Formulas

- `vapi-cli` - Voice AI for developers

````

### Step 2: Create Initial Formula
Create `Formula/vapi-cli.rb`:
```ruby
class VapiCli < Formula
  desc "Voice AI for developers - Vapi CLI"
  homepage "https://vapi.ai"
  version "1.0.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/VapiAI/cli/releases/download/v0.0.1/vapi_Darwin_arm64.tar.gz"
      sha256 "PLACEHOLDER_SHA256"
    else
      url "https://github.com/VapiAI/cli/releases/download/v0.0.1/vapi_Darwin_x86_64.tar.gz"
      sha256 "PLACEHOLDER_SHA256"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/VapiAI/cli/releases/download/v0.0.1/vapi_Linux_arm64.tar.gz"
      sha256 "PLACEHOLDER_SHA256"
    else
      url "https://github.com/VapiAI/cli/releases/download/v0.0.1/vapi_Linux_x86_64.tar.gz"
      sha256 "PLACEHOLDER_SHA256"
    end
  end

  def install
    bin.install "vapi"
  end

  test do
    system "#{bin}/vapi", "--version"
  end
end
````

### Step 3: Enable in GoReleaser

In `.goreleaser.yaml`, change:

```yaml
brews:
  - skip_upload: false # Changed from true
```

## 3. Scoop Setup (Windows)

### Step 1: Create Bucket Repository

```bash
# Create new repository: VapiAI/scoop-bucket
```

**README.md for scoop-bucket:**

````markdown
# Vapi Scoop Bucket

## Installation

```powershell
scoop bucket add vapi https://github.com/VapiAI/scoop-bucket
scoop install vapi-cli
```
````

```

### Step 2: Create Bucket Structure
```

scoop-bucket/
├── bucket/
│ └── vapi-cli.json
└── README.md

````

**bucket/vapi-cli.json:**
```json
{
    "version": "1.0.0",
    "description": "Voice AI for developers - Vapi CLI",
    "homepage": "https://vapi.ai",
    "license": "MIT",
    "architecture": {
        "64bit": {
            "url": "https://github.com/VapiAI/cli/releases/download/v0.0.1/vapi_Windows_x86_64.zip",
            "hash": "PLACEHOLDER_SHA256"
        }
    },
    "bin": "vapi.exe",
    "checkver": {
        "github": "https://github.com/VapiAI/cli"
    },
    "autoupdate": {
        "architecture": {
            "64bit": {
                "url": "https://github.com/VapiAI/cli/releases/download/v$version/vapi_Windows_x86_64.zip"
            }
        }
    }
}
````

### Step 3: Enable in GoReleaser

In `.goreleaser.yaml`, change:

```yaml
scoops:
  - skip_upload: false # Changed from true
```

## 4. NPM Distribution

### Step 1: Create npm Package

Create a new directory `npm-wrapper/`:

**npm-wrapper/package.json:**

```json
{
  "name": "@vapi/cli",
  "version": "1.0.0",
  "description": "Voice AI for developers - Vapi CLI",
  "homepage": "https://vapi.ai",
  "repository": {
    "type": "git",
    "url": "https://github.com/VapiAI/cli.git"
  },
  "keywords": ["vapi", "voice", "ai", "cli"],
  "author": "Dan Goosewin <dan@vapi.ai>",
  "license": "MIT",
  "preferGlobal": true,
  "bin": {
    "vapi": "./bin/vapi"
  },
  "files": ["bin/", "install.js", "README.md"],
  "scripts": {
    "postinstall": "node install.js"
  },
  "engines": {
    "node": ">=14.0.0"
  }
}
```

**npm-wrapper/install.js:**

```javascript
#!/usr/bin/env node

const { execSync } = require("child_process");
const fs = require("fs");
const https = require("https");
const path = require("path");
const { promisify } = require("util");
const stream = require("stream");
const tar = require("tar");
const unzipper = require("unzipper");

const pipeline = promisify(stream.pipeline);

const REPO = "VapiAI/cli";
const BIN_NAME = "vapi";

async function getLatestRelease() {
  return new Promise((resolve, reject) => {
    https
      .get(
        {
          hostname: "api.github.com",
          path: `/repos/${REPO}/releases/latest`,
          headers: { "User-Agent": "vapi-cli-installer" },
        },
        (res) => {
          let data = "";
          res.on("data", (chunk) => (data += chunk));
          res.on("end", () => {
            try {
              resolve(JSON.parse(data));
            } catch (e) {
              reject(e);
            }
          });
        }
      )
      .on("error", reject);
  });
}

function getPlatform() {
  const platform = process.platform;
  const arch = process.arch;

  const platforms = {
    "darwin-x64": "Darwin_x86_64",
    "darwin-arm64": "Darwin_arm64",
    "linux-x64": "Linux_x86_64",
    "linux-arm64": "Linux_arm64",
    "win32-x64": "Windows_x86_64",
  };

  const key = `${platform}-${arch}`;
  if (!platforms[key]) {
    throw new Error(`Unsupported platform: ${key}`);
  }

  return platforms[key];
}

async function downloadBinary(url, destPath) {
  const response = await new Promise((resolve, reject) => {
    https
      .get(url, (res) => {
        if (res.statusCode === 302 || res.statusCode === 301) {
          https.get(res.headers.location, resolve).on("error", reject);
        } else {
          resolve(res);
        }
      })
      .on("error", reject);
  });

  const isZip = url.endsWith(".zip");
  const tempFile = path.join(__dirname, isZip ? "temp.zip" : "temp.tar.gz");

  await pipeline(response, fs.createWriteStream(tempFile));

  if (isZip) {
    await pipeline(
      fs.createReadStream(tempFile),
      unzipper.Extract({ path: path.dirname(destPath) })
    );
  } else {
    await tar.x({
      file: tempFile,
      cwd: path.dirname(destPath),
    });
  }

  fs.unlinkSync(tempFile);

  if (process.platform !== "win32") {
    fs.chmodSync(destPath, 0o755);
  }
}

async function install() {
  try {
    console.log("Installing Vapi CLI...");

    const release = await getLatestRelease();
    const version = release.tag_name;
    const platform = getPlatform();

    const ext = process.platform === "win32" ? ".zip" : ".tar.gz";
    const binExt = process.platform === "win32" ? ".exe" : "";

    const assetName = `vapi_${platform}${ext}`;
    const asset = release.assets.find((a) => a.name === assetName);

    if (!asset) {
      throw new Error(`No binary found for platform: ${platform}`);
    }

    const binDir = path.join(__dirname, "bin");
    if (!fs.existsSync(binDir)) {
      fs.mkdirSync(binDir);
    }

    const binPath = path.join(binDir, BIN_NAME + binExt);

    console.log(`Downloading ${asset.name}...`);
    await downloadBinary(asset.browser_download_url, binPath);

    console.log(`Vapi CLI ${version} installed successfully!`);
    console.log('Run "vapi --help" to get started.');
  } catch (error) {
    console.error("Installation failed:", error.message);
    process.exit(1);
  }
}

if (require.main === module) {
  install();
}
```

**npm-wrapper/README.md:**

````markdown
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
````

## Usage

```bash
vapi --help
```

For more information, visit [vapi.ai](https://vapi.ai)

````

### Step 2: Publish Workflow
Create `.github/workflows/npm-publish.yml`:
```yaml
name: Publish to npm

on:
  release:
    types: [published]

jobs:
  publish:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-node@v3
        with:
          node-version: '18'
          registry-url: 'https://registry.npmjs.org'

      - name: Update package version
        working-directory: npm-wrapper
        run: |
          VERSION=${GITHUB_REF#refs/tags/v}
          npm version $VERSION --no-git-tag-version

      - name: Publish to npm
        working-directory: npm-wrapper
        run: npm publish --access public
        env:
          NODE_AUTH_TOKEN: ${{ secrets.NPM_TOKEN }}
````

### Step 3: Setup npm

1. Create npm account if needed
2. Generate automation token: https://www.npmjs.com/settings/YOUR_USERNAME/tokens
3. Add `NPM_TOKEN` to GitHub secrets

## 5. Docker Hub Setup

### Step 1: Create Docker Hub Repository

1. Create repository at https://hub.docker.com/r/vapiai/cli

### Step 2: Enable in GoReleaser

In `.goreleaser.yaml`, change:

```yaml
dockers:
  - skip_push: false # Changed from true
    image_templates:
      - "vapiai/cli:{{ .Tag }}" # Docker Hub
      - "vapiai/cli:v{{ .Major }}"
      - "vapiai/cli:v{{ .Major }}.{{ .Minor }}"
      - "vapiai/cli:latest"
      - "ghcr.io/vapiai/cli:{{ .Tag }}" # GitHub Container Registry
      - "ghcr.io/vapiai/cli:latest"
```

### Step 3: Add Secrets

Add to GitHub secrets:

- `DOCKERHUB_USERNAME`
- `DOCKERHUB_TOKEN`

Update `.github/workflows/release.yml`:

```yaml
- name: Login to Docker Hub
  uses: docker/login-action@v3
  with:
    username: ${{ secrets.DOCKERHUB_USERNAME }}
    password: ${{ secrets.DOCKERHUB_TOKEN }}

- name: Login to GitHub Container Registry
  uses: docker/login-action@v3
  with:
    registry: ghcr.io
    username: ${{ github.actor }}
    password: ${{ secrets.GITHUB_TOKEN }}
```

## 6. Linux Package Managers (Future)

### APT/DEB (Debian/Ubuntu)

- Use `nfpm` in GoReleaser
- Host on packagecloud.io or custom APT repository

### YUM/RPM (RHEL/CentOS)

- Use `nfpm` in GoReleaser
- Host on packagecloud.io or custom YUM repository

### AUR (Arch Linux)

- Community maintained
- Provide PKGBUILD template

## Release Checklist

Before first release:

- [ ] Create VapiAI/homebrew-tap repository
- [ ] Create VapiAI/scoop-bucket repository
- [ ] Set up npm account and token
- [ ] Configure Docker Hub repository
- [ ] Manually copy install.sh to website repo
- [ ] Add all required GitHub secrets
- [ ] Test GoReleaser locally: `goreleaser release --snapshot --clean`

## Testing Distribution Channels

After setup, test each channel:

```bash
# Homebrew
brew tap vapi/tap
brew install vapi-cli
vapi --version

# Scoop
scoop bucket add vapi https://github.com/VapiAI/scoop-bucket
scoop install vapi-cli
vapi --version

# NPM
npm install -g @vapi/cli
vapi --version

# Docker
docker run --rm vapiai/cli --version

# Direct download
curl -sSL https://vapi.ai/install.sh | bash
vapi --version
```
