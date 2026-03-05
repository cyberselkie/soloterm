#!/bin/bash

set -e

echo "Building soloterm for all platforms..."

CGO_ENABLED=0 GOOS=darwin  GOARCH=arm64 go build -o bin/soloterm_mac   . && echo "  macOS (arm64)"
CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -o bin/soloterm_linux  . && echo "  Linux (amd64)"
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/soloterm.exe    . && echo "  Windows (amd64)"

echo "Done."
