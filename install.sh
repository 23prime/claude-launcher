#!/bin/bash
set -e

# claude-launcher installation script
# Usage: curl -fsSL https://raw.githubusercontent.com/23prime/claude-launcher/main/install.sh | bash

REPO="23prime/claude-launcher"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
BINARY_NAME="claude-launcher"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

# Detect OS
detect_os() {
    case "$(uname -s)" in
        Linux*)     echo "linux" ;;
        Darwin*)    echo "darwin" ;;
        *)          error "Unsupported operating system: $(uname -s)" ;;
    esac
}

# Detect architecture
detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64)   echo "amd64" ;;
        aarch64|arm64)  echo "arm64" ;;
        *)              error "Unsupported architecture: $(uname -m)" ;;
    esac
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Download file
download() {
    local url="$1"
    local output="$2"

    if command_exists curl; then
        curl -fsSL "$url" -o "$output"
    elif command_exists wget; then
        wget -q "$url" -O "$output"
    else
        error "Neither curl nor wget is available. Please install one of them."
    fi
}

main() {
    info "Installing ${BINARY_NAME}..."

    # Detect platform
    OS=$(detect_os)
    ARCH=$(detect_arch)
    info "Detected platform: ${OS}/${ARCH}"

    # Determine download URL
    ARCHIVE_NAME="${BINARY_NAME}-${OS}-${ARCH}.tar.gz"
    DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${ARCHIVE_NAME}"

    # Create temporary directory
    TMP_DIR=$(mktemp -d)
    trap 'rm -rf "${TMP_DIR}"' EXIT

    # Download archive
    info "Downloading ${ARCHIVE_NAME}..."
    download "${DOWNLOAD_URL}" "${TMP_DIR}/${ARCHIVE_NAME}"

    # Extract archive
    info "Extracting archive..."
    tar -xzf "${TMP_DIR}/${ARCHIVE_NAME}" -C "${TMP_DIR}"

    # Find the binary (it should match the pattern claude-launcher-*)
    BINARY_PATH=$(find "${TMP_DIR}" -name "${BINARY_NAME}-${OS}-${ARCH}" -type f -print -quit)
    if [ -z "${BINARY_PATH}" ]; then
        error "Binary not found in archive"
    fi

    # Verify we found exactly one binary
    BINARY_COUNT=$(find "${TMP_DIR}" -name "${BINARY_NAME}-${OS}-${ARCH}" -type f | wc -l)
    if [ "${BINARY_COUNT}" -gt 1 ]; then
        error "Multiple binaries found in archive. Please report this issue."
    fi

    # Make it executable
    chmod +x "${BINARY_PATH}"

    # Install binary
    info "Installing to ${INSTALL_DIR}..."
    if [ -w "${INSTALL_DIR}" ]; then
        mv "${BINARY_PATH}" "${INSTALL_DIR}/${BINARY_NAME}"
    else
        warn "Insufficient permissions. Using sudo..."
        sudo mv "${BINARY_PATH}" "${INSTALL_DIR}/${BINARY_NAME}"
    fi

    # Verify installation
    if command_exists "${BINARY_NAME}"; then
        info "Successfully installed ${BINARY_NAME}!"
        info "Run '${BINARY_NAME} --help' to get started."
    else
        warn "Installation completed, but ${BINARY_NAME} is not in PATH."
        warn "Make sure ${INSTALL_DIR} is in your PATH."
    fi
}

main "$@"
