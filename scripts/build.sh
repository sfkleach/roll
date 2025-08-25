#!/bin/bash

# Build script for the roll dice application
# Usage: ./scripts/build.sh [version]
# If no version is provided, it will try to get it from git tags or use "dev"

set -e

# Get the version
VERSION="${1:-}"
if [ -z "$VERSION" ]; then
    # Try to get version from git tags
    if git describe --tags --exact-match HEAD 2>/dev/null; then
        VERSION=$(git describe --tags --exact-match HEAD | sed 's/^v//')
    elif git describe --tags 2>/dev/null; then
        # If we're not on a tag, use the latest tag with commit info
        VERSION=$(git describe --tags | sed 's/^v//')
    else
        # If no tags exist, use "dev"
        VERSION="dev"
    fi
fi

echo "Building roll dice application version: $VERSION"

# Build with version injection
go build -ldflags "-X github.com/sfkleach/roll/internal/info.Version=$VERSION" -o roll .

echo "Build complete: ./roll"
echo "Version check:"
./roll --version
