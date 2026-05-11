param(
    [string]$Version = "dev"
)

$ErrorActionPreference = "Stop"

$repoRoot = Split-Path -Parent $PSScriptRoot
$goExe = "C:\Program Files\Go\bin\go.exe"
if (-not (Test-Path $goExe)) {
    $goExe = (Get-Command go | Select-Object -ExpandProperty Source)
}

$distRoot = Join-Path $repoRoot "dist"
$binRoot = Join-Path $repoRoot "bin"

New-Item -ItemType Directory -Force -Path $distRoot | Out-Null
New-Item -ItemType Directory -Force -Path $binRoot | Out-Null

function New-TarGzPackage {
    param(
        [string]$TargetOS,
        [string]$Arch
    )

    $platformName = "BoooookDown-$TargetOS-$Arch"
    $stageRoot = Join-Path $distRoot $platformName
    $archivePath = Join-Path $distRoot "$platformName.tar.gz"
    $binaryName = "BoooookDown"
    $binarySource = Join-Path $binRoot $platformName
    $binaryTarget = Join-Path $stageRoot $binaryName

    Remove-Item -Recurse -Force $stageRoot -ErrorAction SilentlyContinue
    New-Item -ItemType Directory -Force -Path $stageRoot | Out-Null

    $env:GOOS = $TargetOS
    $env:GOARCH = $Arch
    & $goExe build -trimpath -buildvcs=false -ldflags "-s -w -X main.version=$Version" -o $binarySource .\cmd\boooookdown

    Copy-Item $binarySource $binaryTarget -Force
    Copy-Item (Join-Path $repoRoot "README.md") (Join-Path $stageRoot "README.md")

    if (Test-Path $archivePath) {
        Remove-Item -Force $archivePath
    }

    tar -czf $archivePath -C $distRoot $platformName
}

& (Join-Path $repoRoot "scripts\build-windows.ps1") -Version $Version

New-TarGzPackage -TargetOS "darwin" -Arch "amd64"
New-TarGzPackage -TargetOS "darwin" -Arch "arm64"
New-TarGzPackage -TargetOS "linux" -Arch "amd64"

Remove-Item Env:GOOS -ErrorAction SilentlyContinue
Remove-Item Env:GOARCH -ErrorAction SilentlyContinue

Get-ChildItem $distRoot | Select-Object FullName, Length, LastWriteTime
