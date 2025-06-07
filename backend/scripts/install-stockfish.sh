#!/bin/bash

set -e

echo "Installing Stockfish chess engine..."

# Detect OS and architecture
OS=$(uname -s)
ARCH=$(uname -m)

# Set download URL based on OS
if [[ "$OS" == "Linux" ]]; then
    if [[ "$ARCH" == "x86_64" ]]; then
        DOWNLOAD_URL="https://github.com/official-stockfish/Stockfish/releases/download/sf_16/stockfish-ubuntu-x86-64-avx2.tar"
        BINARY_NAME="stockfish-ubuntu-x86-64-avx2"
    else
        echo "Unsupported Linux architecture: $ARCH"
        exit 1
    fi
elif [[ "$OS" == "Darwin" ]]; then
    if [[ "$ARCH" == "arm64" ]]; then
        DOWNLOAD_URL="https://github.com/official-stockfish/Stockfish/releases/download/sf_16/stockfish-macos-m1-apple-silicon.tar"
        BINARY_NAME="stockfish-macos-m1-apple-silicon"
    elif [[ "$ARCH" == "x86_64" ]]; then
        DOWNLOAD_URL="https://github.com/official-stockfish/Stockfish/releases/download/sf_16/stockfish-macos-x86-64-avx2.tar"
        BINARY_NAME="stockfish-macos-x86-64-avx2"
    else
        echo "Unsupported macOS architecture: $ARCH"
        exit 1
    fi
else
    echo "Unsupported operating system: $OS"
    exit 1
fi

# Create temporary directory
TEMP_DIR=$(mktemp -d)
cd "$TEMP_DIR"

echo "Downloading Stockfish from: $DOWNLOAD_URL"

# Download Stockfish
if command -v wget >/dev/null 2>&1; then
    wget "$DOWNLOAD_URL" -O stockfish.tar
elif command -v curl >/dev/null 2>&1; then
    curl -L "$DOWNLOAD_URL" -o stockfish.tar
else
    echo "Error: wget or curl is required to download Stockfish"
    exit 1
fi

# Extract the archive
echo "Extracting Stockfish..."
tar -xf stockfish.tar

# Find the stockfish binary
STOCKFISH_PATH=""
if [[ -f "$BINARY_NAME" ]]; then
    STOCKFISH_PATH="$BINARY_NAME"
elif [[ -f "stockfish" ]]; then
    STOCKFISH_PATH="stockfish"
else
    # Search for any executable file that might be stockfish
    for file in *; do
        if [[ -x "$file" && -f "$file" ]]; then
            STOCKFISH_PATH="$file"
            break
        fi
    done
fi

if [[ -z "$STOCKFISH_PATH" ]]; then
    echo "Error: Could not find Stockfish binary in extracted files"
    ls -la
    exit 1
fi

echo "Found Stockfish binary: $STOCKFISH_PATH"

# Make it executable
chmod +x "$STOCKFISH_PATH"

# Test the binary
echo "Testing Stockfish..."
if ! echo "quit" | ./"$STOCKFISH_PATH" >/dev/null 2>&1; then
    echo "Error: Stockfish binary is not working correctly"
    exit 1
fi

# Install to system
INSTALL_DIR="/usr/local/bin"
if [[ ! -w "$INSTALL_DIR" ]]; then
    echo "Installing Stockfish to $INSTALL_DIR (requires sudo)..."
    sudo cp "$STOCKFISH_PATH" "$INSTALL_DIR/stockfish"
    sudo chmod +x "$INSTALL_DIR/stockfish"
else
    echo "Installing Stockfish to $INSTALL_DIR..."
    cp "$STOCKFISH_PATH" "$INSTALL_DIR/stockfish"
    chmod +x "$INSTALL_DIR/stockfish"
fi

# Cleanup
cd /
rm -rf "$TEMP_DIR"

# Verify installation
echo "Verifying installation..."
if command -v stockfish >/dev/null 2>&1; then
    echo "✅ Stockfish installed successfully!"
    echo "Version information:"
    echo "quit" | stockfish | head -n 5
else
    echo "❌ Installation failed: stockfish command not found"
    echo "You may need to add /usr/local/bin to your PATH"
    exit 1
fi

echo ""
echo "🎉 Installation complete!"
echo "You can now use 'stockfish' command or set ENGINE_BINARY_PATH=/usr/local/bin/stockfish" 