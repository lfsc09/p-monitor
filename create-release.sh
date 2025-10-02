#!/bin/bash

# Script to help create a new release
# Usage: ./create-release.sh <version>
# Example: ./create-release.sh v1.0.0

set -e

if [ $# -eq 0 ]; then
    echo "Usage: $0 <version>"
    echo "Example: $0 v1.0.0"
    exit 1
fi

VERSION=$1

echo "Creating release $VERSION..."

# Ensure we're on main branch
CURRENT_BRANCH=$(git branch --show-current)
if [ "$CURRENT_BRANCH" != "main" ]; then
    echo "Warning: You're not on the main branch (current: $CURRENT_BRANCH)"
    read -p "Continue anyway? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Check if there are uncommitted changes
if ! git diff-index --quiet HEAD --; then
    echo "Error: You have uncommitted changes. Please commit or stash them first."
    exit 1
fi

# Update version in control file
sed -i "s/Version: .*/Version: ${VERSION#v}/" debian/DEBIAN/control

# Commit version update
git add debian/DEBIAN/control
git commit -m "Update version to $VERSION"

# Create and push tag
git tag -a "$VERSION" -m "Release $VERSION"
git push origin main
git push origin "$VERSION"

echo "Release $VERSION created successfully!"
echo "The GitHub workflow will now build and attach the .deb package to the release."
echo "You can check the progress at: https://github.com/lfsc09/p-monitor/actions"
