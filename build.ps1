<#
.SYNOPSIS
    Build script - Generateur d'Infographie v2.1 (Go Edition, WYSIWYG)
.DESCRIPTION
    Telecharge les bibliotheques JavaScript vendorees (versions epinglees)
    puis compile le binaire Go. A executer depuis la racine du projet.

    L'empreinte SHA-256 du binaire est affichee en fin de build :
    la publier avec chaque release permet aux utilisateurs de verifier
    l'integrite du fichier telecharge.
.EXAMPLE
    .\build.ps1
    .\build.ps1 -SkipVendor
    .\build.ps1 -ForceVendor
    .\build.ps1 -Target linux
    .\build.ps1 -Console
#>

param(
    [switch]$SkipVendor,
    [switch]$ForceVendor,
    [switch]$Console,
    [ValidateSet("windows","linux","darwin")]
    [string]$Target = "windows",
    [string]$OutputName = ""
)

$ErrorActionPreference = "Stop"
$ProgressPreference = "SilentlyContinue"

# --- Configuration ---
$vendorDir = "web\editor\vendor"
$version = "2.2.1"

# Versions epinglees : ne pas passer en "latest" — la reproductibilite du
# build et la stabilite de l'editeur en dependent.
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
        Write-Host "    Go detecte: $goVersion" -ForegroundColor Green
        return $true
    } catch {
        Write-Host "    ERREUR: Go n'est pas installe ou absent du PATH." -ForegroundColor Red
        Write-Host "    Telechargement: https://go.dev/dl/" -ForegroundColor Yellow
        return $false
    }
}

function Download-VendorFiles {
    Write-Step "Telechargement des bibliotheques vendorees..."

    if (-not (Test-Path $vendorDir)) {
        New-Item -ItemType Directory -Path $vendorDir -Force | Out-Null
        Write-Host "    Repertoire cree: $vendorDir"
    }

    foreach ($file in $vendorFiles) {
        $outPath = Join-Path $vendorDir $file.Name

        if ((Test-Path $outPath) -and -not $ForceVendor) {
            Write-Host "    [SKIP] $($file.Name) (deja present ; -ForceVendor pour re-telecharger)" -ForegroundColor DarkGray
            continue
        }

        Write-Host "    Telechargement $($file.Name)..." -NoNewline
        try {
            Invoke-WebRequest -Uri $file.URL -OutFile $outPath -UseBasicParsing
            $size = (Get-Item $outPath).Length
            Write-Host " OK ($([math]::Round($size/1024, 1)) KB)" -ForegroundColor Green
        } catch {
            Write-Host " ECHEC" -ForegroundColor Red
            Write-Host "    URL: $($file.URL)" -ForegroundColor Yellow
            Write-Host "    Erreur: $_" -ForegroundColor Red
            exit 1
        }
    }

    Write-Host "    Bibliotheques vendorees pretes." -ForegroundColor Green
}

function Build-GoBinary {
    Write-Step "Compilation du binaire Go..."

    # Resolution des modules.
    # Note: go mod tidy ecrit sa progression sur stderr, que PowerShell
    # traiterait comme des erreurs avec ErrorActionPreference=Stop.
    Write-Host "    go mod tidy..."
    $prevPref = $ErrorActionPreference
    $ErrorActionPreference = "Continue"
    & go mod tidy 2>&1 | ForEach-Object {
        Write-Host "    $($_.ToString())" -ForegroundColor DarkGray
    }
    $modTidyExit = $LASTEXITCODE
    $ErrorActionPreference = $prevPref
    if ($modTidyExit -ne 0) {
        Write-Host "    ERREUR: go mod tidy a echoue (code $modTidyExit)." -ForegroundColor Red
        exit 1
    }

    # Nom de sortie
    $ext = ""
    $goos = $Target
    $goarch = "amd64"

    if ($Target -eq "windows") { $ext = ".exe" }

    if (-not $OutputName) {
        $script:OutputName = "infographic-generator${ext}"
    }

    $outPath = Join-Path "." $script:OutputName

    # Flags de build.
    # -X main.version=... : injecte la version dans le binaire.
    #   IMPORTANT: ne fonctionne que parce que main.go declare `var version`
    #   (une const est figee a la compilation et ignorerait silencieusement -X).
    $ldflags = "-s -w -X main.version=$version"
    if ($Target -eq "windows" -and -not $Console) {
        # -H windowsgui : pas de fenetre console sous Windows.
        # L'application vit en arriere-plan, pilotee par le navigateur.
        # Utiliser -Console pour garder la console (debug).
        $ldflags = "-s -w -H windowsgui -X main.version=$version"
    }
    $env:GOOS = $goos
    $env:GOARCH = $goarch
    $env:CGO_ENABLED = "0"

    Write-Host "    Cible: $goos/$goarch"
    Write-Host "    Sortie: $($script:OutputName)"
    Write-Host "    CGO: desactive"
    Write-Host "    Version: $version"
    Write-Host ""
    Write-Host "    Compilation..." -NoNewline

    $prevPref = $ErrorActionPreference
    $ErrorActionPreference = "Continue"
    & go build -ldflags $ldflags -o $outPath . 2>&1 | ForEach-Object {
        Write-Host "    $($_.ToString())" -ForegroundColor DarkGray
    }
    $buildExit = $LASTEXITCODE
    $ErrorActionPreference = $prevPref

    if ($buildExit -ne 0) {
        Write-Host " ECHEC" -ForegroundColor Red
        exit 1
    }

    $size = (Get-Item $outPath).Length
    Write-Host " OK" -ForegroundColor Green
    Write-Host ""
    Write-Host "    Binaire: $outPath ($([math]::Round($size/1MB, 2)) MB)" -ForegroundColor Green

    # Empreinte a publier avec la release (verification d'integrite,
    # indispensable tant que le binaire n'est pas signe).
    $hash = (Get-FileHash -Algorithm SHA256 -Path $outPath).Hash.ToLower()
    Write-Host "    SHA-256: $hash" -ForegroundColor Yellow
    $script:BinaryHash = $hash

    # Nettoyage de l'environnement
    Remove-Item Env:\GOOS -ErrorAction SilentlyContinue
    Remove-Item Env:\GOARCH -ErrorAction SilentlyContinue
    Remove-Item Env:\CGO_ENABLED -ErrorAction SilentlyContinue

    return $outPath
}

# --- Main ---

Write-Host ""
Write-Host "=============================================" -ForegroundColor White
Write-Host "  Generateur d'Infographie v$version - Build" -ForegroundColor White
Write-Host "  (edition WYSIWYG)" -ForegroundColor White
Write-Host "=============================================" -ForegroundColor White

# Verifier Go
Write-Step "Verification des prerequis..."
if (-not (Test-GoInstalled)) { exit 1 }

# Verifier qu'on est bien a la racine du projet COMPLET.
# Cas classique : le script est lance depuis un dossier ne contenant que des
# fichiers partiels (patch decompresse seul) ou depuis un sous-repertoire.
$requis = @("go.mod", "main.go", "internal", "web")
$manquants = $requis | Where-Object { -not (Test-Path $_) }
if ($manquants.Count -gt 0) {
    Write-Host "    ERREUR: fichiers/dossiers du projet introuvables ici: $($manquants -join ', ')" -ForegroundColor Red
    Write-Host "    Ce script doit etre lance depuis la RACINE du projet complet" -ForegroundColor Yellow
    Write-Host "    (celle qui contient go.mod, main.go, internal\ et web\)." -ForegroundColor Yellow
    Write-Host "    Si vous n'avez que les fichiers modifies, decompressez-les" -ForegroundColor Yellow
    Write-Host "    par-dessus votre clone du depot, puis relancez depuis ce clone." -ForegroundColor Yellow
    exit 1
}

# Bibliotheques vendorees
if (-not $SkipVendor) {
    Download-VendorFiles
} else {
    Write-Step "Telechargement vendor ignore (-SkipVendor)"
    if (-not (Test-Path (Join-Path $vendorDir "alpine.min.js"))) {
        Write-Host "    ATTENTION: bibliotheques manquantes ! Relancer sans -SkipVendor." -ForegroundColor Yellow
    }
}

# Compilation
$binaryPath = Build-GoBinary

# Resume
Write-Host ""
Write-Host "=============================================" -ForegroundColor Green
Write-Host "  BUILD REUSSI" -ForegroundColor Green
Write-Host "=============================================" -ForegroundColor Green
Write-Host ""
Write-Host "  Lancer:" -ForegroundColor White
Write-Host "    .\$OutputName" -ForegroundColor Yellow
Write-Host ""
Write-Host "  Options d'execution:" -ForegroundColor White
Write-Host "    .\$OutputName -port 8080          # Port fixe" -ForegroundColor DarkGray
Write-Host "    .\$OutputName -no-kiosk           # Onglet navigateur classique" -ForegroundColor DarkGray
Write-Host "    .\$OutputName -no-browser         # Serveur seul" -ForegroundColor DarkGray
Write-Host "    .\$OutputName -data C:\mes\donnees # Repertoire de donnees" -ForegroundColor DarkGray
Write-Host "    .\$OutputName -version            # Afficher la version" -ForegroundColor DarkGray
Write-Host ""
Write-Host "  Options de build:" -ForegroundColor White
Write-Host "    .\build.ps1 -Console              # Garder la console (debug)" -ForegroundColor DarkGray
Write-Host "    .\build.ps1 -ForceVendor          # Re-telecharger les bibliotheques" -ForegroundColor DarkGray
Write-Host "    .\build.ps1 -Target linux         # Cross-compilation Linux" -ForegroundColor DarkGray
Write-Host ""
Write-Host "  Publication (GitHub Releases):" -ForegroundColor White
Write-Host "    Joindre l'empreinte SHA-256 ci-dessus a la description de la release." -ForegroundColor DarkGray
Write-Host "    Le binaire ne doit PAS etre commite dans Git (voir .gitignore)." -ForegroundColor DarkGray
Write-Host ""
