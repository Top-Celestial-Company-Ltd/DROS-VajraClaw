#!/usr/bin/env bash
set -euo pipefail

PREFIX="/opt/dros-muse-hacker"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "=== DROS-Muse Hacker install ==="

sudo mkdir -p "$PREFIX/dros_hook"
sudo mkdir -p "$PREFIX/config"

sudo cp "$SCRIPT_DIR/VERSION" "$PREFIX/VERSION"
sudo cp "$SCRIPT_DIR/dros_guard.py" "$PREFIX/dros_guard.py"
sudo cp "$SCRIPT_DIR/dros_hook/sitecustomize.py" "$PREFIX/dros_hook/sitecustomize.py"
sudo cp "$SCRIPT_DIR/config/policy.json" "$PREFIX/config/policy.json"

sudo chmod 755 "$PREFIX"
sudo chmod 644 "$PREFIX/VERSION" \
               "$PREFIX/dros_guard.py" \
               "$PREFIX/dros_hook/sitecustomize.py" \
               "$PREFIX/config/policy.json"

echo
echo "Installed to: $PREFIX"
echo
echo "NOTE:"
echo "This installer does not modify any unknown Muse systemd unit."
echo "Runtime activation must be explicitly configured for the target Muse process."
echo
echo "INSTALL=PASS"
