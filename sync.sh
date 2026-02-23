#!/bin/bash
set -e

# Params
VERSION=${1:-"latest"}
HOST=${h:-""}
PORT=${p:-""}

# Detect architecture
ARCH=$(uname -m)
OS=$(uname -s)

if [ "$OS" != "Linux" ]; then
    echo "This script is for Linux. Please use sync.ps1 for Windows."
    exit 1
fi

case $ARCH in
    x86_64)
        ARCH_LABEL="x64"
        ;;
    aarch64|arm64)
        ARCH_LABEL="arm64"
        ;;
    *)
        echo "Install Failed: Linux architecture $ARCH is not supported."
        exit 1
        ;;
esac

INSTALL_ROOT="$HOME/.tml-sync"
BIN_DIR="$INSTALL_ROOT/bin"
mkdir -p "$BIN_DIR"

REPO="Ashu11-A/tModLoader-sync"
BASE_URL="https://github.com/$REPO/releases"
# Binary naming convention: client-linux-x64 or client-linux-arm64
TARGET="client-linux-$ARCH_LABEL"

if [ "$VERSION" = "latest" ]; then
    URL="$BASE_URL/latest/download/$TARGET"
else
    URL="$BASE_URL/download/$VERSION/$TARGET"
fi

EXE_PATH="$BIN_DIR/tml-sync"

echo "Downloading tModLoader-sync ($VERSION) for linux/$ARCH_LABEL..."
if command -v curl >/dev/null 2>&1; then
    curl -#fLo "$EXE_PATH" "$URL"
elif command -v wget >/dev/null 2>&1; then
    wget -q --show-progress -O "$EXE_PATH" "$URL"
else
    echo "Error: curl or wget is required."
    exit 1
fi

chmod +x "$EXE_PATH"

# Add to PATH if not already there
SHELL_PROFILE=""
case $SHELL in
*/zsh)
    SHELL_PROFILE="$HOME/.zshrc"
    ;;
*/bash)
    SHELL_PROFILE="$HOME/.bashrc"
    ;;
*)
    SHELL_PROFILE="$HOME/.profile"
    ;;
esac

if [[ ":$PATH:" != *":$BIN_DIR:"* ]]; then
    echo "Adding $BIN_DIR to PATH in $SHELL_PROFILE"
    echo "export PATH=\"\$PATH:$BIN_DIR\"" >> "$SHELL_PROFILE"
    export PATH="$PATH:$BIN_DIR"
fi

C_RESET="\033[0m"
C_GREEN="\033[1;32m"

echo -e "${C_GREEN}tModLoader-sync was installed successfully!${C_RESET}"

if [ -n "$HOST" ] && [ -n "$PORT" ]; then
    PORT_VAL=$(echo "$PORT" | sed 's/^://')
    echo "To start syncing, run:"
    echo -e "${C_GREEN}tml-sync --host $HOST --port $PORT_VAL${C_RESET}"
else
    echo "To start syncing, run: tml-sync --host <IP> --port <PORT>"
fi
