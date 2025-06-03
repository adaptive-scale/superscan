#!/bin/bash

# Exit on error
set -e

# Version information
VERSION=${VERSION:-"dev"}
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(git rev-parse HEAD 2>/dev/null || echo "unknown")

# Build flags
LDFLAGS="-X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME -X main.GitCommit=$GIT_COMMIT"

# Create build directory
mkdir -p build

# Build for current platform
echo "Building for $(go env GOOS)/$(go env GOARCH)..."
go build -ldflags "$LDFLAGS" -o build/superscan

# Build for multiple platforms if specified
if [ "$1" == "all" ]; then
    echo "Building for multiple platforms..."
    
    # List of platforms to build for
    PLATFORMS=(
        "darwin/amd64"
        "darwin/arm64"
        "linux/amd64"
        "linux/arm64"
        "windows/amd64"
    )

    for PLATFORM in "${PLATFORMS[@]}"; do
        OS="${PLATFORM%/*}"
        ARCH="${PLATFORM#*/}"
        OUTPUT="build/superscan-$OS-$ARCH"
        
        if [ "$OS" == "windows" ]; then
            OUTPUT="$OUTPUT.exe"
        fi
        
        echo "Building for $OS/$ARCH..."
        GOOS=$OS GOARCH=$ARCH go build -ldflags "$LDFLAGS" -o "$OUTPUT"
    done
fi

echo "Build complete! Output files are in the build directory." 