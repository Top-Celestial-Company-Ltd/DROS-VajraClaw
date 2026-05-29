#!/bin/bash

# DROS VajraClaw iOS SDK Release Pipeline
# Designed to run on macOS (Requires Xcode & gomobile)

set -e

echo "Starting VajraClaw iOS SDK RELEASE compilation pipeline..."
echo "==========================================================="
echo "[SECURITY] Engaging Anti-Reverse Engineering Protections..."
echo "  -> Stripping DWARF debug symbols (-s)"
echo "  -> Stripping DWARF symbol tables (-w)"
echo "  -> Erasing local path identities (-trimpath)"
echo "==========================================================="

# Ensure gomobile is installed
if ! command -v gomobile &> /dev/null
then
    echo "Error: gomobile could not be found. Please run 'go install golang.org/x/mobile/cmd/gomobile@latest'"
    exit 1
fi

gomobile init

# Define Paths
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
BASE_DIR="$(dirname "$SCRIPT_DIR")"
MOBILE_PKG="$BASE_DIR/vajraclaw_sdk/mobile"
TARGET_XCFRAMEWORK="$BASE_DIR/vajraclaw_sdk/ios/VajraClawRelease.xcframework"

# Ensure output directory exists
mkdir -p "$BASE_DIR/vajraclaw_sdk/ios"

cd "$MOBILE_PKG"

echo "Compiling iOS XCFramework dynamic package..."
# Execute the release build with aggressive stripping
gomobile bind -target=ios,iossimulator -ldflags="-s -w" -trimpath -o "$TARGET_XCFRAMEWORK" -v

if [ -d "$TARGET_XCFRAMEWORK" ]; then
    echo "✅ iOS SDK RELEASE compilation succeeded: $TARGET_XCFRAMEWORK"
else
    echo "❌ Compilation failed to produce RELEASE XCFramework!"
    exit 1
fi

echo "🎉 VajraClaw iOS SDK Release Build Pipeline Completed Successfully!"
