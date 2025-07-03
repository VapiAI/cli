# Vapi CLI installation script for Windows
# Usage: iex ((New-Object System.Net.WebClient).DownloadString('https://vapi.ai/install.ps1'))

$ErrorActionPreference = "Stop"

# Configuration
$Repo = "VapiAI/cli"
$BinaryName = "vapi.exe"
$InstallDir = "$env:LOCALAPPDATA\Programs\Vapi"

# Helper functions
function Write-Info($Message) {
    Write-Host "[INFO] $Message" -ForegroundColor Green
}

function Write-Error($Message) {
    Write-Host "[ERROR] $Message" -ForegroundColor Red
    exit 1
}

function Write-Warning($Message) {
    Write-Host "[WARN] $Message" -ForegroundColor Yellow
}

# Detect architecture
function Get-Platform {
    $arch = $env:PROCESSOR_ARCHITECTURE
    
    switch ($arch) {
        "AMD64" { return "Windows_x86_64" }
        "ARM64" { return "Windows_arm64" }
        default { Write-Error "Unsupported architecture: $arch" }
    }
}

# Get latest release version
function Get-LatestVersion {
    Write-Info "Fetching latest version..."
    
    try {
        $response = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest"
        $version = $response.tag_name
        
        if (-not $version) {
            Write-Error "Failed to fetch latest version"
        }
        
        Write-Info "Latest version: $version"
        return $version
    }
    catch {
        Write-Error "Failed to fetch latest version: $_"
    }
}

# Download and install
function Install-Vapi($Version, $Platform) {
    $url = "https://github.com/$Repo/releases/download/$Version/vapi_$Platform.tar.gz"
    $tempDir = [System.IO.Path]::GetTempPath() + [System.Guid]::NewGuid().ToString()
    $tarFile = "$tempDir\vapi.tar.gz"
    
    Write-Info "Downloading Vapi CLI..."
    Write-Info "URL: $url"
    
    # Create temp directory
    New-Item -ItemType Directory -Path $tempDir -Force | Out-Null
    
    try {
        # Download the file
        Invoke-WebRequest -Uri $url -OutFile $tarFile
        
        Write-Info "Extracting..."
        
        # Extract tar.gz (requires tar.exe available in Windows 10+)
        if (Get-Command tar -ErrorAction SilentlyContinue) {
            tar -xzf $tarFile -C $tempDir
        } else {
            Write-Error "tar command not found. Please update to Windows 10 version 1803 or later."
        }
        
        # Create install directory
        if (-not (Test-Path $InstallDir)) {
            New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
        }
        
        # Move binary to install directory
        $binaryPath = "$tempDir\vapi.exe"
        if (-not (Test-Path $binaryPath)) {
            $binaryPath = "$tempDir\vapi"  # Try without .exe extension
        }
        
        if (Test-Path $binaryPath) {
            Move-Item $binaryPath "$InstallDir\$BinaryName" -Force
        } else {
            Write-Error "Binary not found after extraction"
        }
        
        Write-Info "Vapi CLI installed successfully!"
    }
    catch {
        Write-Error "Installation failed: $_"
    }
    finally {
        # Cleanup
        if (Test-Path $tempDir) {
            Remove-Item $tempDir -Recurse -Force
        }
    }
}

# Add to PATH
function Add-ToPath {
    $currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    
    if ($currentPath -notlike "*$InstallDir*") {
        Write-Info "Adding Vapi CLI to PATH..."
        $newPath = "$currentPath;$InstallDir"
        [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
        
        # Update current session PATH
        $env:PATH = "$env:PATH;$InstallDir"
        
        Write-Info "Added $InstallDir to your PATH"
    } else {
        Write-Info "Vapi CLI directory already in PATH"
    }
}

# Verify installation
function Test-Installation {
    $vapiPath = "$InstallDir\$BinaryName"
    
    if (Test-Path $vapiPath) {
        # Test if vapi command works
        try {
            $version = & $vapiPath --version 2>$null
            Write-Info "Verification: $version"
            Write-Host ""
            Write-Info "Installation complete! 🎉"
            Write-Host ""
            Write-Host "Get started with:"
            Write-Host "  vapi login"
            Write-Host "  vapi --help"
            Write-Host ""
            Write-Warning "Please restart your terminal or PowerShell session to use 'vapi' command globally."
        }
        catch {
            Write-Warning "Vapi CLI was installed but verification failed"
            Write-Warning "You may need to restart your terminal"
        }
    } else {
        Write-Error "Installation verification failed - binary not found"
    }
}

# Main installation flow
function Main {
    Write-Host "===================================" -ForegroundColor Cyan
    Write-Host "    Vapi CLI Installer" -ForegroundColor Cyan
    Write-Host "===================================" -ForegroundColor Cyan
    Write-Host ""
    
    $platform = Get-Platform
    Write-Info "Detected platform: $platform"
    
    $version = Get-LatestVersion
    Install-Vapi $version $platform
    Add-ToPath
    Test-Installation
}

# Run main function
Main 