#!/bin/bash
set -e

echo "Running golangci-lint..."
if ! command -v golangci-lint &> /dev/null; then
    echo "golangci-lint not found, installing..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
fi
golangci-lint run ./...

echo "Running go vet..."
go vet ./...

echo "Lint completed successfully"
