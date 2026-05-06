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
$stageRoot = Join-Path $distRoot "go-down-textbook-windows-amd64"
$iconPath = Join-Path $repoRoot "assets\app.ico"
$resourcePath = Join-Path $repoRoot "cmd\go-down-textbook\rsrc_windows_amd64.syso"
$exePath = Join-Path $stageRoot "go-down-textbook.exe"
$zipPath = Join-Path $distRoot "go-down-textbook-windows-amd64.zip"

New-Item -ItemType Directory -Force -Path $distRoot | Out-Null
Remove-Item -Recurse -Force $stageRoot -ErrorAction SilentlyContinue
New-Item -ItemType Directory -Force -Path $stageRoot | Out-Null

python (Join-Path $repoRoot "scripts\make_windows_icon.py") (Join-Path $repoRoot "logo.jpg") $iconPath

& $goExe run github.com/akavel/rsrc@latest -ico $iconPath -o $resourcePath
& $goExe build -ldflags "-s -w -X main.version=$Version" -o $exePath .\cmd\go-down-textbook
Remove-Item -Force $resourcePath -ErrorAction SilentlyContinue

Copy-Item (Join-Path $repoRoot "README.md") (Join-Path $stageRoot "README.md")

if (Test-Path $zipPath) {
    Remove-Item -Force $zipPath
}
Compress-Archive -Path (Join-Path $stageRoot "*") -DestinationPath $zipPath

Write-Output $exePath
Write-Output $zipPath
