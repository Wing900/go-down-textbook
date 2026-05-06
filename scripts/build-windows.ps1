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
$buildRoot = Join-Path $distRoot "build-tmp"
$iconPath = Join-Path $repoRoot "assets\app.ico"
$resourcePath = Join-Path $repoRoot "cmd\go-down-textbook\rsrc_windows_amd64.syso"
$rawExePath = Join-Path $buildRoot "go-down-textbook.exe"
$exePath = Join-Path $stageRoot "go-down-textbook.exe"
$zipPath = Join-Path $distRoot "go-down-textbook-windows-amd64.zip"
$repoUpx = Join-Path $repoRoot "tools\upx\upx.exe"
$upxExe = $null

if (Test-Path $repoUpx) {
    $upxExe = $repoUpx
} else {
    $upxCmd = Get-Command upx -ErrorAction SilentlyContinue
    if ($upxCmd) {
        $upxExe = $upxCmd.Source
    }
}

New-Item -ItemType Directory -Force -Path $distRoot | Out-Null
Remove-Item -Recurse -Force $buildRoot -ErrorAction SilentlyContinue
New-Item -ItemType Directory -Force -Path $buildRoot | Out-Null
Remove-Item -Recurse -Force $stageRoot -ErrorAction SilentlyContinue
New-Item -ItemType Directory -Force -Path $stageRoot | Out-Null

python (Join-Path $repoRoot "scripts\make_windows_icon.py") (Join-Path $repoRoot "logo.jpg") $iconPath

& $goExe run github.com/akavel/rsrc@latest -ico $iconPath -o $resourcePath
& $goExe build -trimpath -buildvcs=false -ldflags "-s -w -X main.version=$Version" -o $rawExePath .\cmd\go-down-textbook
Remove-Item -Force $resourcePath -ErrorAction SilentlyContinue

if ($upxExe) {
    $packed = $false
    for ($i = 0; $i -lt 5 -and -not $packed; $i++) {
        try {
            & $upxExe --best --lzma $rawExePath
            if ($LASTEXITCODE -eq 0) {
                $packed = $true
                break
            }
        } catch {
        }
        Start-Sleep -Milliseconds 800
    }

    if (-not $packed) {
        throw "UPX 压缩失败: $rawExePath"
    }
}

Copy-Item $rawExePath $exePath -Force

Copy-Item (Join-Path $repoRoot "README.md") (Join-Path $stageRoot "README.md")

if (Test-Path $zipPath) {
    Remove-Item -Force $zipPath
}
Compress-Archive -Path (Join-Path $stageRoot "*") -DestinationPath $zipPath

Write-Output $exePath
Write-Output $zipPath
