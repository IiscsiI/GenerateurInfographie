<#
.SYNOPSIS
    Build script for Infographic Generator v2.0 (Go Edition)
.DESCRIPTION
    Downloads vendor JavaScript libraries then compiles the Go binary.
    Run from the project root directory.
.EXAMPLE
    .\build.ps1
    .\build.ps1 -SkipVendor
    .\build.ps1 -Target linux
#>

param(
    [switch]$SkipVendor,
    [switch]$Console,
    [ValidateSet("windows","linux","darwin")]
    [string]$Target = "windows",
    [string]$OutputName = ""
)

$ErrorActionPreference = "Stop"
$ProgressPreference = "SilentlyContinue"

# --- Configuration ---
$vendorDir = "web\editor\vendor"
$version = "2.0.0"

$vendorFiles = @(
    @{
        Name = "alpine.min.js"
        URL  = "https://cdn.jsdelivr.net/npm/alpinejs@3.14.8/dist/cdn.min.js"
    },
    @{
        Name = "pickr.min.js"
        URL  = "https://cdn.jsdelivr.net/npm/@simonwep/pickr@1.9.1/dist/pickr.min.js"
    },
    @{
        Name = "pickr-nano.min.css"
        URL  = "https://cdn.jsdelivr.net/npm/@simonwep/pickr@1.9.1/dist/themes/nano.min.css"
    },
    @{
        Name = "Sortable.min.js"
        URL  = "https://cdn.jsdelivr.net/npm/sortablejs@1.15.6/Sortable.min.js"
    }
)

# --- Functions ---

function Write-Step {
    param([string]$Message)
    Write-Host ""
    Write-Host "==> $Message" -ForegroundColor Cyan
}

function Test-GoInstalled {
    try {
        $goVersion = & go version 2>&1
        Write-Host "    Go detected: $goVersion" -ForegroundColor Green
        return $true
    } catch {
        Write-Host "    ERROR: Go is not installed or not in PATH." -ForegroundColor Red
        Write-Host "    Download from: https://go.dev/dl/" -ForegroundColor Yellow
        return $false
    }
}

function Download-VendorFiles {
    Write-Step "Downloading vendor JavaScript libraries..."

    if (-not (Test-Path $vendorDir)) {
        New-Item -ItemType Directory -Path $vendorDir -Force | Out-Null
        Write-Host "    Created directory: $vendorDir"
    }

    foreach ($file in $vendorFiles) {
        $outPath = Join-Path $vendorDir $file.Name

        if (Test-Path $outPath) {
            Write-Host "    [SKIP] $($file.Name) (already exists)" -ForegroundColor DarkGray
            continue
        }

        Write-Host "    Downloading $($file.Name)..." -NoNewline
        try {
            Invoke-WebRequest -Uri $file.URL -OutFile $outPath -UseBasicParsing
            $size = (Get-Item $outPath).Length
            Write-Host " OK ($([math]::Round($size/1024, 1)) KB)" -ForegroundColor Green
        } catch {
            Write-Host " FAILED" -ForegroundColor Red
            Write-Host "    URL: $($file.URL)" -ForegroundColor Yellow
            Write-Host "    Error: $_" -ForegroundColor Red
            exit 1
        }
    }

    Write-Host "    All vendor files ready." -ForegroundColor Green
}

function Build-GoBinary {
    Write-Step "Building Go binary..."

    # Resolve go modules
    # Note: go mod tidy writes download progress to stderr, which PowerShell
    # treats as errors with ErrorActionPreference=Stop. We temporarily relax this.
    Write-Host "    Running go mod tidy..."
    $prevPref = $ErrorActionPreference
    $ErrorActionPreference = "Continue"
    & go mod tidy 2>&1 | ForEach-Object {
        $line = $_.ToString()
        if ($line -match "^go:") {
            Write-Host "    $line" -ForegroundColor DarkGray
        } else {
            Write-Host "    $line" -ForegroundColor DarkGray
        }
    }
    $modTidyExit = $LASTEXITCODE
    $ErrorActionPreference = $prevPref
    if ($modTidyExit -ne 0) {
        Write-Host "    ERROR: go mod tidy failed (exit code $modTidyExit)." -ForegroundColor Red
        exit 1
    }

    # Determine output name
    $ext = ""
    $goos = $Target
    $goarch = "amd64"

    if ($Target -eq "windows") { $ext = ".exe" }

    if (-not $OutputName) {
        $OutputName = "infographic-generator${ext}"
    }

    $outPath = Join-Path "." $OutputName

    # Build flags
    $ldflags = "-s -w -X main.version=$version"
    if ($Target -eq "windows" -and -not $Console) {
        # -H windowsgui: no console window on Windows.
        # The app runs as a background process, controlled entirely via the browser.
        # Use -Console flag to keep the console for debugging.
        $ldflags = "-s -w -H windowsgui -X main.version=$version"
    }
    $env:GOOS = $goos
    $env:GOARCH = $goarch
    $env:CGO_ENABLED = "0"

    Write-Host "    Target: $goos/$goarch"
    Write-Host "    Output: $OutputName"
    Write-Host "    CGO: disabled"
    Write-Host "    Version: $version"
    Write-Host ""
    Write-Host "    Compiling..." -NoNewline

    $prevPref = $ErrorActionPreference
    $ErrorActionPreference = "Continue"
    & go build -ldflags $ldflags -o $outPath . 2>&1 | ForEach-Object {
        Write-Host "    $($_.ToString())" -ForegroundColor DarkGray
    }
    $buildExit = $LASTEXITCODE
    $ErrorActionPreference = $prevPref

    if ($buildExit -ne 0) {
        Write-Host " FAILED" -ForegroundColor Red
        exit 1
    }

    $size = (Get-Item $outPath).Length
    Write-Host " OK" -ForegroundColor Green
    Write-Host ""
    Write-Host "    Binary: $outPath ($([math]::Round($size/1MB, 2)) MB)" -ForegroundColor Green

    # Clean up env
    Remove-Item Env:\GOOS -ErrorAction SilentlyContinue
    Remove-Item Env:\GOARCH -ErrorAction SilentlyContinue
    Remove-Item Env:\CGO_ENABLED -ErrorAction SilentlyContinue

    return $outPath
}

# --- Main ---

Write-Host ""
Write-Host "============================================" -ForegroundColor White
Write-Host "  Infographic Generator v$version - Build" -ForegroundColor White
Write-Host "============================================" -ForegroundColor White

# Check Go
Write-Step "Checking prerequisites..."
if (-not (Test-GoInstalled)) { exit 1 }

# Vendor JS
if (-not $SkipVendor) {
    Download-VendorFiles
} else {
    Write-Step "Skipping vendor download (-SkipVendor)"
    if (-not (Test-Path (Join-Path $vendorDir "alpine.min.js"))) {
        Write-Host "    WARNING: Vendor files missing! Run without -SkipVendor first." -ForegroundColor Yellow
    }
}

# Build
$binaryPath = Build-GoBinary

# Summary
Write-Host ""
Write-Host "============================================" -ForegroundColor Green
Write-Host "  BUILD SUCCESSFUL" -ForegroundColor Green
Write-Host "============================================" -ForegroundColor Green
Write-Host ""
Write-Host "  To run:" -ForegroundColor White
Write-Host "    .\$OutputName" -ForegroundColor Yellow
Write-Host ""
Write-Host "  Options:" -ForegroundColor White
Write-Host "    .\$OutputName -port 8080        # Fixed port" -ForegroundColor DarkGray
Write-Host "    .\$OutputName -no-kiosk          # Normal browser tab" -ForegroundColor DarkGray
Write-Host "    .\$OutputName -no-browser         # Server only" -ForegroundColor DarkGray
Write-Host "    .\$OutputName -data C:\my\data   # Custom data dir" -ForegroundColor DarkGray
Write-Host ""
Write-Host "  Build options:" -ForegroundColor White
Write-Host "    .\build.ps1 -Console              # Keep console window (debug)" -ForegroundColor DarkGray
Write-Host "    .\build.ps1 -Target linux          # Cross-compile for Linux" -ForegroundColor DarkGray
Write-Host ""
