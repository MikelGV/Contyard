#!/bin/bash
VERSION="0.0.1"
OUTPUT_DIR="../bin"
mkdir -p $OUTPUT_DIR

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $OUTPUT_DIR/contyard-linux-amd64-$VERSION ../main.go
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o $OUTPUT_DIR/contyard-linux-arm64-$VERSION ../main.go

GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o $OUTPUT_DIR/contyard-darwin-amd64-$VERSION ../main.go 
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o $OUTPUT_DIR/contyard-darwin-arm64-$VERSION ../main.go 

GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o $OUTPUT_DIR/contyard-windows-amd64-$VERSION.exe ../main.go 
GOOS=windows GOARCH=arm64 go build -ldflags="-s -w" -o $OUTPUT_DIR/contyard-windows-arm64-$VERSION.exe ../main.go 

cd $OUTPUT_DIR
tar -czf contyard-linux-amd64-$VERSION.tar.gz contyard-linux-amd64-$VERSION
tar -czf contyard-linux-arm64-$VERSION.tar.gz contyard-linux-arm64-$VERSION
tar -czf contyard-darwin-amd64-$VERSION.tar.gz contyard-linux-amd64-$VERSION
tar -czf contyard-darwin-arm64-$VERSION.tar.gz contyard-linux-arm64-$VERSION
zip contyard-windows-amd64-$VERSION.zip contyard-windows-amd64-$VERSION.exe
zip contyard-windows-arm64-$VERSION.zip contyard-windows-arm64-$VERSION.exe
