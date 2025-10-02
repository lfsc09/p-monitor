#!/bin/bash

set -e

echo "Building p-monitor .deb package..."

# Check for required dependencies
echo "Checking dependencies..."
MISSING_DEPS=()
for dep in libgl1-mesa-dev libx11-dev libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev libxxf86vm-dev libglfw3-dev pkg-config; do
    if ! dpkg -l | grep -q "^ii.*$dep"; then
        MISSING_DEPS+=("$dep")
    fi
done

if [ ${#MISSING_DEPS[@]} -ne 0 ]; then
    echo "Missing dependencies: ${MISSING_DEPS[*]}"
    echo "Please install them with:"
    echo "sudo apt install ${MISSING_DEPS[*]}"
    exit 1
fi

# Clean previous builds
rm -rf debian/usr/bin/p-monitor
rm -rf debian/usr/share/p-monitor/assets/*

# Build the application
echo "Building Go application..."
go build -o debian/usr/bin/p-monitor ./cmd

# Copy assets
echo "Copying assets..."
cp assets/*.png debian/usr/share/p-monitor/assets/

# Set permissions
chmod +x debian/usr/bin/p-monitor

# Build the .deb package
echo "Creating .deb package..."
dpkg-deb --build debian p-monitor_1.0.0_amd64.deb

echo "Package created: p-monitor_1.0.0_amd64.deb"
echo "To install: sudo dpkg -i p-monitor_1.0.0_amd64.deb"
echo "To fix dependencies: sudo apt-get install -f"
