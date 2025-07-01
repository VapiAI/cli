#!/usr/bin/env node

const https = require('https');
const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');
const { pipeline } = require('stream/promises');

const REPO = 'VapiAI/cli';
const BIN_NAME = 'vapi';

async function httpsGet(url, options = {}) {
  return new Promise((resolve, reject) => {
    https.get(url, options, (res) => {
      if (res.statusCode === 301 || res.statusCode === 302) {
        httpsGet(res.headers.location, options).then(resolve).catch(reject);
      } else {
        resolve(res);
      }
    }).on('error', reject);
  });
}

async function getLatestRelease() {
  console.log('Fetching latest release...');
  const res = await httpsGet(`https://api.github.com/repos/${REPO}/releases/latest`, {
    headers: { 'User-Agent': 'vapi-cli-installer' }
  });
  
  let data = '';
  for await (const chunk of res) {
    data += chunk;
  }
  
  return JSON.parse(data);
}

function getPlatform() {
  const platform = process.platform;
  const arch = process.arch;
  
  const platforms = {
    'darwin-x64': 'Darwin_x86_64',
    'darwin-arm64': 'Darwin_arm64',
    'linux-x64': 'Linux_x86_64',
    'linux-arm64': 'Linux_arm64',
    'linux-arm': 'Linux_armv7',
    'win32-x64': 'Windows_x86_64'
  };
  
  const key = `${platform}-${arch}`;
  if (!platforms[key]) {
    throw new Error(`Unsupported platform: ${key}`);
  }
  
  return { 
    platform: platforms[key], 
    isWindows: platform === 'win32' 
  };
}

async function downloadAndExtract(url, destDir) {
  const isWindows = process.platform === 'win32';
  const ext = isWindows ? '.zip' : '.tar.gz';
  const tempFile = path.join(destDir, `temp${ext}`);
  
  console.log(`Downloading from ${url}...`);
  
  // Download file
  const response = await httpsGet(url);
  await pipeline(response, fs.createWriteStream(tempFile));
  
  // Extract file
  console.log('Extracting...');
  if (isWindows) {
    // Use PowerShell to extract zip on Windows
    execSync(`powershell -command "Expand-Archive -Path '${tempFile}' -DestinationPath '${destDir}' -Force"`, { stdio: 'inherit' });
  } else {
    // Use tar for Unix systems
    execSync(`tar -xzf "${tempFile}" -C "${destDir}"`, { stdio: 'inherit' });
  }
  
  // Clean up temp file
  fs.unlinkSync(tempFile);
}

async function install() {
  try {
    const release = await getLatestRelease();
    const version = release.tag_name;
    const { platform, isWindows } = getPlatform();
    
    const ext = isWindows ? '.zip' : '.tar.gz';
    const binExt = isWindows ? '.exe' : '';
    const assetName = `cli_${platform}${ext}`;
    
    const asset = release.assets.find(a => a.name === assetName);
    if (!asset) {
      throw new Error(`No binary found for platform: ${platform}\nAvailable assets: ${release.assets.map(a => a.name).join(', ')}`);
    }
    
    // Create bin directory
    const binDir = path.join(__dirname, 'bin');
    if (!fs.existsSync(binDir)) {
      fs.mkdirSync(binDir, { recursive: true });
    }
    
    // Download and extract
    await downloadAndExtract(asset.browser_download_url, binDir);
    
    // Set executable permissions on Unix
    if (!isWindows) {
      const binPath = path.join(binDir, BIN_NAME);
      fs.chmodSync(binPath, 0o755);
    }
    
    console.log(`\n✅ Vapi CLI ${version} installed successfully!`);
    console.log('Run "vapi --help" to get started.\n');
    
  } catch (error) {
    console.error('\n❌ Installation failed:', error.message);
    if (error.stack) {
      console.error('\nStack trace:', error.stack);
    }
    process.exit(1);
  }
}

// Run installation
if (require.main === module) {
  install();
} 