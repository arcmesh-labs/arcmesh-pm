#Requires -Version 5.1
<#
.SYNOPSIS
    Installs apm (arcmesh-pm) for the current user. No admin rights required.
#>
Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

# ── architecture ──────────────────────────────────────────────────────────────

$arch = if ([System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture -eq
            [System.Runtime.InteropServices.Architecture]::Arm64) { 'arm64' } else { 'amd64' }

# ── paths ─────────────────────────────────────────────────────────────────────

$installDir = Join-Path $env:USERPROFILE 'bin'
$exePath    = Join-Path $installDir 'apm.exe'
$zipUrl     = "https://github.com/arcmesh-labs/arcmesh-pm/releases/latest/download/apm-windows-$arch.zip"
$tmpZip     = Join-Path $env:TEMP "apm-windows-$arch.zip"

# ── download ──────────────────────────────────────────────────────────────────

Write-Host "Detected architecture: $arch"
Write-Host "Downloading $zipUrl ..."

try {
    Invoke-WebRequest -Uri $zipUrl -OutFile $tmpZip -UseBasicParsing
} catch {
    Write-Error "Download failed: $_"
    exit 1
}

Write-Host "Download complete."

# ── extract ───────────────────────────────────────────────────────────────────

if (-not (Test-Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir | Out-Null
    Write-Host "Created $installDir"
}

Write-Host "Extracting to $installDir ..."

Add-Type -AssemblyName System.IO.Compression.FileSystem
$zip = [System.IO.Compression.ZipFile]::OpenRead($tmpZip)
try {
    foreach ($entry in $zip.Entries) {
        if ($entry.Name -eq 'apm.exe') {
            [System.IO.Compression.ZipFileExtensions]::ExtractToFile($entry, $exePath, $true)
            break
        }
    }
} finally {
    $zip.Dispose()
    Remove-Item $tmpZip -Force
}

if (-not (Test-Path $exePath)) {
    Write-Error "Extraction failed: apm.exe not found in archive."
    exit 1
}

Write-Host "Installed apm.exe to $exePath"

# ── PATH ──────────────────────────────────────────────────────────────────────

$userPath = [Environment]::GetEnvironmentVariable('PATH', 'User')
$dirs = $userPath -split ';' | Where-Object { $_ -ne '' }

if ($dirs -notcontains $installDir) {
    $newPath = ($dirs + $installDir) -join ';'
    [Environment]::SetEnvironmentVariable('PATH', $newPath, 'User')
    Write-Host "Added $installDir to your PATH."
} else {
    Write-Host "$installDir is already in your PATH."
}

# ── done ──────────────────────────────────────────────────────────────────────

Write-Host ""
Write-Host "apm installed successfully."
Write-Host "Open a new terminal window for the PATH change to take effect, then run: apm --help"
