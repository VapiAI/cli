# Vapi CLI installation script for Windows
# Usage: iex ((New-Object System.Net.WebClient).DownloadString('https://vapi.ai/install.ps1'))

[CmdletBinding()]
param()

$ErrorActionPreference = "Stop"
$ProgressPreference = "SilentlyContinue"  # Suppress progress bars for faster downloads

# Ensure compatibility with PowerShell 5.1
if (-not (Test-Path Variable:IsWindows)) {
    # We're on PowerShell 5.1 or earlier (which only runs on Windows)
    $IsWindows = $true
}

# Check if running on Windows
if (-not $IsWindows) {
    Write-Host "[ERROR] This installer is for Windows only." -ForegroundColor Red
    Write-Host ""
    if ($IsMacOS -or $IsLinux) {
        Write-Host "For macOS/Linux, use the shell script instead:" -ForegroundColor Yellow
        Write-Host "  curl -sSL https://vapi.ai/install.sh | bash" -ForegroundColor White
    }
    exit 1
}

# Configuration
$Repo = "VapiAI/cli"
$BinaryName = "vapi.exe"
$InstallDir = "$env:LOCALAPPDATA\Programs\Vapi"
$ManDir = "$InstallDir\docs\man"

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
        default { 
            # Also check PROCESSOR_ARCHITEW6432 for 32-bit processes on 64-bit systems
            $arch64 = $env:PROCESSOR_ARCHITEW6432
            if ($arch64 -eq "AMD64") {
                return "Windows_x86_64"
            }
            Write-Error "Unsupported architecture: $arch" 
        }
    }
}

# Get latest release version
function Get-LatestVersion {
    Write-Info "Fetching latest version..."
    
    try {
        $headers = @{}
        # Add User-Agent for better GitHub API compatibility
        $headers["User-Agent"] = "Vapi-CLI-Installer/1.0"
        
        $response = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest" -Headers $headers
        $version = $response.tag_name
        
        if (-not $version) {
            Write-Error "Failed to fetch latest version"
        }
        
        Write-Info "Latest version: $version"
        return $version
    }
    catch {
        Write-Error "Failed to fetch latest version: $($_.Exception.Message)"
    }
}

# Download and install
function Install-Vapi($Version, $Platform) {
    $url = "https://github.com/$Repo/releases/download/$Version/cli_$Platform.tar.gz"
    $tempDir = [System.IO.Path]::GetTempPath() + [System.Guid]::NewGuid().ToString()
    $tarFile = "$tempDir\vapi.tar.gz"
    
    Write-Info "Downloading Vapi CLI..."
    Write-Info "URL: $url"
    
    # Create temp directory
    New-Item -ItemType Directory -Path $tempDir -Force | Out-Null
    
    try {
        # Download the file with progress
        $webClient = New-Object System.Net.WebClient
        $webClient.Headers.Add("User-Agent", "Vapi-CLI-Installer/1.0")
        $webClient.DownloadFile($url, $tarFile)
        $webClient.Dispose()
        
        Write-Info "Extracting..."
        
        # Check for tar.exe (available in Windows 10 1803+)
        $tarCmd = Get-Command tar -ErrorAction SilentlyContinue
        if ($tarCmd) {
            tar -xzf $tarFile -C $tempDir
        } else {
            Write-Error "tar command not found. Please update to Windows 10 version 1803 or later, or install a compatible extraction tool."
        }
        
        # Create install directory
        if (-not (Test-Path $InstallDir)) {
            New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
        }
        
        # Create man page directory
        if (-not (Test-Path $ManDir)) {
            New-Item -ItemType Directory -Path $ManDir -Force | Out-Null
        }
        
        # Move binary to install directory
        $binaryPath = "$tempDir\vapi.exe"
        if (-not (Test-Path $binaryPath)) {
            $binaryPath = "$tempDir\vapi"  # Try without .exe extension
        }
        
        if (Test-Path $binaryPath) {
            Copy-Item $binaryPath "$InstallDir\$BinaryName" -Force
        } else {
            Write-Error "Binary not found after extraction"
        }
        
        # Extract man pages if available
        $manPagesFound = 0
        $manFiles = @()
        
        # Check for man pages in man/ subdirectory
        $manSubdir = "$tempDir\man"
        if (Test-Path $manSubdir) {
            $manFiles += Get-ChildItem "$manSubdir\*.1" -ErrorAction SilentlyContinue
        }
        
        # Check for man pages in root directory
        $manFiles += Get-ChildItem "$tempDir\*.1" -ErrorAction SilentlyContinue
        
        foreach ($manFile in $manFiles) {
            try {
                Copy-Item $manFile.FullName $ManDir -Force
                $manPagesFound++
            }
            catch {
                # Continue on error - man pages are not critical for Windows
            }
        }
        
        if ($manPagesFound -gt 0) {
            Write-Info "Extracted $manPagesFound manual page(s) to $ManDir"
            Write-Info "Note: Manual pages are primarily for Unix systems. On Windows, use 'vapi --help' for documentation."
        }
        
        Write-Info "Vapi CLI installed successfully to: $InstallDir"
    }
    catch {
        Write-Error "Installation failed: $($_.Exception.Message)"
    }
    finally {
        # Cleanup
        if (Test-Path $tempDir) {
            Remove-Item $tempDir -Recurse -Force -ErrorAction SilentlyContinue
        }
    }
}

# Add to PATH
function Add-ToPath {
    $currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    
    if ($currentPath -notlike "*$InstallDir*") {
        Write-Info "Adding Vapi CLI to PATH..."
        
        # Clean up PATH before adding (remove trailing semicolons)
        $currentPath = $currentPath.TrimEnd(';')
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
            $output = & $vapiPath --version 2>&1
            if ($LASTEXITCODE -eq 0) {
                Write-Info "Verification successful: $output"
                Write-Host ""
                Write-Info "Installation complete! 🎉"
                Write-Host ""
                Write-Host "Get started with:" -ForegroundColor Cyan
                Write-Host "  vapi login" -ForegroundColor White
                Write-Host "  vapi --help" -ForegroundColor White
                Write-Host ""
                
                # Check if man pages were installed
                if (Test-Path "$ManDir\vapi.1") {
                    Write-Host "Documentation:" -ForegroundColor Cyan
                    Write-Host "  Manual pages are available in: $ManDir" -ForegroundColor White
                    Write-Host "  Use 'vapi --help' for command-line help" -ForegroundColor White
                    Write-Host ""
                }
                
                Write-Warning "Please restart your terminal or PowerShell session to use 'vapi' command globally."
            } else {
                Write-Warning "Vapi CLI was installed but verification failed (exit code: $LASTEXITCODE)"
                Write-Warning "You may need to restart your terminal"
            }
        }
        catch {
            Write-Warning "Vapi CLI was installed but verification failed: $($_.Exception.Message)"
            Write-Warning "You may need to restart your terminal"
        }
    } else {
        Write-Error "Installation verification failed - binary not found at: $vapiPath"
    }
}

# Main installation flow
function Main {
    Write-Host "===================================" -ForegroundColor Cyan
    Write-Host "    Vapi CLI Installer" -ForegroundColor Cyan
    Write-Host "===================================" -ForegroundColor Cyan
    Write-Host ""
    
    # Check Windows version
    $osVersion = [System.Environment]::OSVersion.Version
    if ($osVersion.Major -lt 10) {
        Write-Warning "Windows 10 or later is recommended for best compatibility"
    }
    
    $platform = Get-Platform
    Write-Info "Detected platform: $platform"
    
    $version = Get-LatestVersion
    Install-Vapi $version $platform
    Add-ToPath
    Test-Installation
}

# Run main function
Main 
