#!/usr/bin/env bash
# ==========================================================================
# Build du Générateur d'Infographie (équivalent build.ps1 pour Linux/macOS)
# Usage :
#   ./build.sh                  build pour l'OS courant
#   ./build.sh windows          cross-compilation (windows | linux | darwin)
#   SKIP_VENDOR=1 ./build.sh    sans re-télécharger les bibliothèques
# ==========================================================================
set -euo pipefail

VERSION="2.1.0"
TARGET="${1:-$(go env GOOS)}"
VENDOR_DIR="web/editor/vendor"

declare -A VENDOR=(
  ["alpine.min.js"]="https://cdn.jsdelivr.net/npm/alpinejs@3.14.8/dist/cdn.min.js"
  ["pickr.min.js"]="https://cdn.jsdelivr.net/npm/@simonwep/pickr@1.9.1/dist/pickr.min.js"
  ["pickr-nano.min.css"]="https://cdn.jsdelivr.net/npm/@simonwep/pickr@1.9.1/dist/themes/nano.min.css"
  ["Sortable.min.js"]="https://cdn.jsdelivr.net/npm/sortablejs@1.15.6/Sortable.min.js"
)

if [[ "${SKIP_VENDOR:-0}" != "1" ]]; then
  echo "==> Téléchargement des bibliothèques vendorées (versions épinglées)"
  mkdir -p "$VENDOR_DIR"
  for name in "${!VENDOR[@]}"; do
    curl -fsSL "${VENDOR[$name]}" -o "$VENDOR_DIR/$name"
    echo "    $name OK"
  done
fi

OUT="infographic-generator"
[[ "$TARGET" == "windows" ]] && OUT="infographic-generator.exe"

echo "==> Compilation Go ($TARGET, v$VERSION)"
LDFLAGS="-s -w -X main.version=$VERSION"
[[ "$TARGET" == "windows" ]] && LDFLAGS="$LDFLAGS -H=windowsgui"
GOOS="$TARGET" GOARCH=amd64 go build -ldflags "$LDFLAGS" -o "$OUT" .

echo "==> OK : $OUT ($(du -h "$OUT" | cut -f1))"
echo "==> SHA-256 (à publier avec la release) :"
sha256sum "$OUT" 2>/dev/null || shasum -a 256 "$OUT"
