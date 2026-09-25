#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
VERSION="1.1"
BUILD="$ROOT/build"
DIST="$ROOT/dist/MiniBin - $VERSION"
rm -rf "$BUILD" "$ROOT/dist"
mkdir -p "$BUILD" "$DIST"
cd "$ROOT"
GOOS=windows GOARCH=386 CGO_ENABLED=0 go test -c -buildvcs=false -o "$BUILD/compilecheck.test.exe" ./src
rm -f "$BUILD/compilecheck.test.exe"
GOOS=windows GOARCH=386 CGO_ENABLED=0 go build -buildvcs=false -trimpath -ldflags='-s -w -H=windowsgui -buildid=' -o "$BUILD/MiniBin.exe" ./src
cp "$BUILD/MiniBin.exe" "$DIST/"
cp minibin.ini "$DIST/"
cp assets/{empty,25,50,75,full}.ico "$DIST/"
cp docs/USER_GUIDE_RU.md "$DIST/README.md"
sha256sum "$DIST/MiniBin.exe" | sed 's# .*/#  #' > "$DIST/SHA256SUMS.txt"
printf '%s\n' "$VERSION" > "$DIST/VERSION.txt"
(cd "$ROOT/dist" && zip -qr "MiniBin - $VERSION.zip" "MiniBin - $VERSION")
echo "Built $ROOT/dist/MiniBin - $VERSION.zip"
