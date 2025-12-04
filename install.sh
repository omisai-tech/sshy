#!/bin/bash
set -e

REPO="omisai-tech/sshy"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

detect_os() {
    case "$(uname -s)" in
        Linux*)  echo "linux" ;;
        Darwin*) echo "darwin" ;;
        MINGW*|MSYS*|CYGWIN*) echo "windows" ;;
        *) echo "unknown" ;;
    esac
}

detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64) echo "amd64" ;;
        aarch64|arm64) echo "arm64" ;;
        *) echo "unknown" ;;
    esac
}

get_latest_version() {
    curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | cut -d'"' -f4
}

main() {
    OS=$(detect_os)
    ARCH=$(detect_arch)
    VERSION=$(get_latest_version)

    if [ "$OS" = "unknown" ] || [ "$ARCH" = "unknown" ]; then
        echo "Error: Unsupported OS or architecture"
        echo "OS: $(uname -s), Arch: $(uname -m)"
        exit 1
    fi

    if [ -z "$VERSION" ]; then
        echo "Error: Could not determine latest version"
        exit 1
    fi

    echo "Installing sshy ${VERSION} for ${OS}/${ARCH}..."

    EXT="tar.gz"
    if [ "$OS" = "windows" ]; then
        EXT="zip"
    fi

    FILENAME="sshy_${VERSION}_${OS}_${ARCH}.${EXT}"
    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${FILENAME}"

    TMP_DIR=$(mktemp -d)
    trap "rm -rf ${TMP_DIR}" EXIT

    echo "Downloading ${DOWNLOAD_URL}..."
    curl -sL "$DOWNLOAD_URL" -o "${TMP_DIR}/${FILENAME}"

    echo "Extracting..."
    if [ "$EXT" = "tar.gz" ]; then
        tar xzf "${TMP_DIR}/${FILENAME}" -C "${TMP_DIR}"
    else
        unzip -q "${TMP_DIR}/${FILENAME}" -d "${TMP_DIR}"
    fi

    BINARY="sshy"
    if [ "$OS" = "windows" ]; then
        BINARY="sshy.exe"
    fi

    if [ -w "$INSTALL_DIR" ]; then
        mv "${TMP_DIR}/${BINARY}" "${INSTALL_DIR}/"
    else
        echo "Installing to ${INSTALL_DIR} (requires sudo)..."
        sudo mv "${TMP_DIR}/${BINARY}" "${INSTALL_DIR}/"
    fi

    chmod +x "${INSTALL_DIR}/${BINARY}"

    echo "Successfully installed sshy ${VERSION} to ${INSTALL_DIR}/${BINARY}"
    echo ""
    echo "Run 'sshy init' to get started!"
}

main
