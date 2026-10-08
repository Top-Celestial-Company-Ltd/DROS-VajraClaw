#!/usr/bin/env bash
set -euo pipefail

PREFIX="/opt/dros-muse-hacker"

echo "=== DROS-Muse Hacker status ==="

if [ ! -d "$PREFIX" ]; then
    echo "STATUS=NOT_INSTALLED"
    exit 0
fi

echo "PREFIX=$PREFIX"

for f in \
    "$PREFIX/VERSION" \
    "$PREFIX/dros_guard.py" \
    "$PREFIX/dros_hook/sitecustomize.py" \
    "$PREFIX/config/policy.json"
do
    if [ -f "$f" ]; then
        echo "PRESENT=$f"
    else
        echo "MISSING=$f"
    fi
done

echo
echo "IMPORTANT:"
echo "Presence of the package does NOT prove that a Muse runtime is using the hook."
echo "STATUS=INSTALLED"
