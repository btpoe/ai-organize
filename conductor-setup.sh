#!/bin/bash
set -e

echo "🔧 Setting up File Organizer workspace..."

# Check for Go
if ! command -v go &> /dev/null; then
    echo "❌ Error: Go is not installed. Please install Go 1.21 or later."
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "✓ Found Go $GO_VERSION"

# Check for Node.js
if ! command -v node &> /dev/null; then
    echo "❌ Error: Node.js is not installed. Please install Node.js 18 or later."
    exit 1
fi

NODE_VERSION=$(node --version)
echo "✓ Found Node.js $NODE_VERSION"

# Check for Wails CLI
if ! command -v wails &> /dev/null; then
    echo "❌ Error: Wails CLI is not installed."
    echo "   Install it with: go install github.com/wailsapp/wails/v2/cmd/wails@latest"
    exit 1
fi

WAILS_VERSION=$(wails version | head -n 1)
echo "✓ Found $WAILS_VERSION"

# Install frontend dependencies
echo ""
echo "📦 Installing frontend dependencies..."
cd frontend
npm install
cd ..

# Download Go dependencies
echo ""
echo "📦 Downloading Go dependencies..."
go mod download

echo ""
echo "✅ Workspace setup complete!"
echo "   Run 'wails dev' to start the development server."