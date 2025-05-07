#!/bin/bash
set -e

echo "Checking for outdated dependencies..."
if ! command -v go-mod-outdated &> /dev/null; then
    echo "go-mod-outdated not found, installing..."
    go install github.com/psampaz/go-mod-outdated@latest
fi
go list -u -m -json all | go-mod-outdated -update -direct

echo "Dependency check completed"
