#!/bin/bash

# Exit on error
set -e

# Create bin directory if it doesn't exist
mkdir -p bin

# Build the application
echo "Building superscan..."
go build -o bin/superscan cmd/superscan/main.go

echo "Build complete! Binary is available at bin/superscan" 