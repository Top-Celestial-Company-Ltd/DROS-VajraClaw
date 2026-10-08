#!/usr/bin/env bash
set -euo pipefail

PREFIX="/opt/dros-muse-hacker"

echo "=== DROS-Muse Hacker uninstall ==="

if [ -d "$PREFIX" ]; then
    sudo rm -rf "$PREFIX"
    echo "Removed: $PREFIX"
else
    echo "Already absent: $PREFIX"
fi

echo "UNINSTALL=PASS"
