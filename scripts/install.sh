#!/usr/bin/env bash
# MobileOps CLI installer
# Usage: curl -fsSL https://raw.githubusercontent.com/MobileOps-Team/mobileops-cli/main/scripts/install.sh | bash
set -euo pipefail

REPO="MobileOps-Team/mobileops-cli"
BINARY_NAME="mobileops"

echo ""
echo "  MobileOps CLI Installer"
echo "  ========================"
echo ""

# Detect OS
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "$OS" in
  linux)  OS="linux" ;;
  darwin) OS="darwin" ;;
  mingw*|msys*|cygwin*)
    echo "Error: Windows detected. Download the zip for your architecture from"
    echo "  https://github.com/${REPO}/releases/latest"
    echo "and put mobileops.exe somewhere on your PATH."
    exit 1
    ;;
  *)
    echo "Error: Unsupported operating system: $OS"
    exit 1
    ;;
esac

# Detect architecture
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64|amd64)  ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *)
    echo "Error: Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

echo "  Platform: ${OS}/${ARCH}"

# Get latest version
echo "  Fetching latest version..."
VERSION=$(curl -sSf "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
if [ -z "$VERSION" ]; then
  echo "Error: Could not detect latest version."
  echo "Check https://github.com/${REPO}/releases for available versions."
  exit 1
fi
echo "  Version:  ${VERSION}"

# Build download URL
FILENAME="mobileops-cli_${VERSION#v}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/${VERSION}/${FILENAME}"

# Choose install directory
INSTALL_DIR="/usr/local/bin"
NEED_SUDO=false
if [ ! -w "$INSTALL_DIR" ]; then
  if command -v sudo &> /dev/null; then
    NEED_SUDO=true
  else
    INSTALL_DIR="${HOME}/.local/bin"
    mkdir -p "$INSTALL_DIR"
  fi
fi

# Download and install
TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT

echo "  Downloading..."
curl -sSfL "$URL" -o "${TMPDIR}/${FILENAME}"

echo "  Installing to ${INSTALL_DIR}/${BINARY_NAME}..."
tar -xzf "${TMPDIR}/${FILENAME}" -C "$TMPDIR"

if [ "$NEED_SUDO" = true ]; then
  sudo mv "${TMPDIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
  sudo chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
else
  mv "${TMPDIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
  chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
fi

# Verify
echo ""
echo "  Installed successfully!"
echo "  $("${INSTALL_DIR}/${BINARY_NAME}" version)"
echo ""
echo "  Get started:"
echo "    mobileops auth login"
echo "    mobileops vessels list"
echo ""
echo "  For AI agents:"
echo "    mobileops --help --agent"
echo ""

# Check PATH
if ! echo "$PATH" | tr ':' '\n' | grep -qx "$INSTALL_DIR"; then
  echo "  Note: Add ${INSTALL_DIR} to your PATH:"
  echo "    export PATH=\"${INSTALL_DIR}:\$PATH\""
  echo ""
fi
