$ErrorActionPreference = "Stop"

$Version = "1.1"
$Root = Split-Path -Parent $PSScriptRoot
$Build = Join-Path $Root "build"
$Dist = Join-Path $Root "dist"
$ReleaseDir = Join-Path $Dist "MiniBin - $Version"
$Zip = Join-Path $Dist "MiniBin - $Version.zip"

Remove-Item $Build -Recurse -Force -ErrorAction SilentlyContinue
Remove-Item $ReleaseDir -Recurse -Force -ErrorAction SilentlyContinue
Remove-Item $Zip -Force -ErrorAction SilentlyContinue
New-Item -ItemType Directory -Force $Build, $ReleaseDir | Out-Null

$env:GOOS = "windows"
$env:GOARCH = "386"
$env:CGO_ENABLED = "0"

Push-Location $Root
try {
    go test -buildvcs=false ./...
    go build -buildvcs=false -trimpath -ldflags "-s -w -H=windowsgui -buildid=" -o "$Build\MiniBin.exe" .\src
} finally {
    Pop-Location
}

Copy-Item "$Build\MiniBin.exe" $ReleaseDir
Copy-Item "$Root\minibin.ini" $ReleaseDir
Copy-Item "$Root\assets\empty.ico" $ReleaseDir
Copy-Item "$Root\assets\25.ico" $ReleaseDir
Copy-Item "$Root\assets\50.ico" $ReleaseDir
Copy-Item "$Root\assets\75.ico" $ReleaseDir
Copy-Item "$Root\assets\full.ico" $ReleaseDir
Copy-Item "$Root\docs\USER_GUIDE_RU.md" "$ReleaseDir\README.md"

$Hash = (Get-FileHash "$ReleaseDir\MiniBin.exe" -Algorithm SHA256).Hash.ToLowerInvariant()
"MiniBin.exe  $Hash" | Set-Content "$ReleaseDir\SHA256SUMS.txt" -Encoding ASCII
"1.1" | Set-Content "$ReleaseDir\VERSION.txt" -Encoding ASCII

Compress-Archive -Path "$ReleaseDir\*" -DestinationPath $Zip -CompressionLevel Optimal
Write-Host "Built: $Zip"
Write-Host "SHA-256 MiniBin.exe: $Hash"
