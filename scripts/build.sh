#!/bin/bash
set -e

VERSION="${TAG:-0.0.1-beta}"
OUTPUT_DIR="bin"
mkdir -p $OUTPUT_DIR

command -v tar >/dev/null 2>&1 || { echo "tar is required but not installed"; exit 1; }
command -v zip >/dev/null 2>&1 || { echo "zip is required but not installed"; exit 1; }

rm -f "$OUTPUT_DIR"/.tar.gz*

PLATFORMS=("linux/amd64" "linux/arm64" "darwin/amd64" "darwin/arm64" "windows/amd64" "windows/arm64")
for platform in "${PLATFORMS[@]}"; do
    GOOS=${platform%/*}
    GOARCH=${platform#*/}
    OUTPUT="${OUTPUT_DIR}/contyard-${GOOS}-${GOARCH}-${VERSION}"
    if [ "$GOOS" = "windows" ]; then
        OUTPUT="${OUTPUT}.exe"
    fi

    echo "Building for $GOOS/$GOARCH..."
    CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="-s -w" -o "$OUTPUT" ./main.go
done

cd "$OUTPUT_DIR"
for platform in "${PLATFORMS[@]}"; do
    GOOS=${platform%/*}
    GOARCH=${platform#*/}
    BINARY="contyard-${GOOS}-${GOARCH}-${VERSION}"
    ARCHIVE="contyard-${GOOS}-${GOARCH}-${VERSION}"
    if [ "$GOOS" = "windows" ]; then
        ARCHIVE="${ARCHIVE}.zip"
        BINARY="${BINARY}.exe"
        echo "Creating $ARCHIVE..."
        zip "$ARCHIVE" "$BINARY" || { echo "Failed to create $ARCHIVE"; exit 1; }
    else
        ARCHIVE="${ARCHIVE}.tar.gz"
        echo "Creating $ARCHIVE..."
        tar -czf "$ARCHIVE" "$BINARY" || { echo "Failed to create $ARCHIVE"; exit 1; }

    fi
done



