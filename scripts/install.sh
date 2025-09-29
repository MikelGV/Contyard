#!/bin/bash
set -e

REPO="MikelGV/Contyard"
BINARY_NAME="contyard"
INSTALL_PATH="usr/local/bin/${BINARY_NAME}"

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64 | arm64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac
case "$OS" in
    linux | darwin | windows) ;;
    *) echo "Unsuported OS: $OS"; exit 1 ;;
esac

VERSION="${VERSION:-0.0.1-beta}"
VERSION=${VERSION#v}

FILE="${BINARY_NAME}-${OS}-${ARCH}-${VERSION}.tar.gz"
URL="https://github.com/${REPO}/releases/download/v${VERSION}/${FILE}"

echo "Downloading ${FILE} from ${URL}..."
curl -sL "$URL" -o "/tmp/${FILE}" || { echo "Failed to download ${FILE}"; exit 1; }
tar -xzf "/tmp/${FILE}" -C /tmp || { echo "Failed to extract ${FILE}"; exit 1; }
BINARY="/tmp/${BINARY_NAME}-${OS}-${ARCH}-${VERSION}"
if [ ! -f "$BINARY_NAME" ]; then
    echo "Binary ${BINARY} not found in archive"
    exit 1
fi

echo "Installing ${BINARY_NAME} to ${INSTALL_PATH}..."
sudo mv "$BINARY" "$INSTALL_PATH" || { echo "Failed to move binary to ${INSTALL_PATH}. Try without sudo or check permissions."; exit 1; }
sudo chmod +x "$INSTALL_PATH" || { echo "Failed to sdet executable permissions"; exit 1; }

rm "/tmp/${FILE}"
echo "${BINARY_NAME} v${VERSION} installed successfully! Run '${BINARY_NAME} --version' to verify."
${BINARY_NAME} --version
