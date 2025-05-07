#!/bin/bash
set -e

# Install gosec if not available
if ! command -v gosec &> /dev/null; then
    echo "gosec not found, installing..."
    go install github.com/securego/gosec/v2/cmd/gosec@latest
fi

echo "Running gosec security scanner..."
gosec -quiet ./...

# Install nancy if not available
if ! command -v nancy &> /dev/null; then
    echo "nancy not found, installing..."
    go install github.com/sonatype-nexus-community/nancy@latest
fi

echo "Running nancy for dependency vulnerabilities..."
go list -json -deps ./... | nancy sleuth

echo "Security scan completed"
