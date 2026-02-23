#!/bin/bash
set -e

# Configuration parameters
VERSION=${1:-"latest"}
HOST=${host_address:-${h:-${host:-${HOST:-""}}}}
PORT=${port_number:-${p:-${port:-${PORT:-""}}}}

# Detect system and architecture
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
  EXPORT_LINE="export PATH=\"\$PATH:$BIN_DIR\""
  if ! grep -Fxq "$EXPORT_LINE" "$SHELL_PROFILE"; then
    echo "Adding $BIN_DIR to PATH in $SHELL_PROFILE"
    echo "$EXPORT_LINE" >> "$SHELL_PROFILE"
  fi
  export PATH="$PATH:$BIN_DIR"
fi

C_RESET="\033[0m"
C_GREEN="\033[1;32m"

echo -e "${C_GREEN}tModLoader-sync was installed successfully!${C_RESET}"

# Execute immediately after installation
if [ -n "$HOST" ] && [ -n "$PORT" ]; then
  PORT_VAL=$(echo "$PORT" | sed 's/^://')
  echo -e "${C_GREEN}Starting tml-sync connected to $HOST:$PORT_VAL...${C_RESET}"
  "$EXE_PATH" --host "$HOST" --port "$PORT_VAL"
else
  echo -e "${C_GREEN}Starting tml-sync${C_RESET}"
  "$EXE_PATH"
fi