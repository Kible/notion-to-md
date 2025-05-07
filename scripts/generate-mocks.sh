#!/bin/bash
set -e

echo "Generating mocks with mockgen..."
if ! command -v mockgen &> /dev/null; then
    echo "mockgen not found, installing..."
    go install github.com/golang/mock/mockgen@latest
fi

# Add your mockgen commands here
# Example:
# mockgen -source=internal/service/service.go -destination=internal/mocks/service_mock.go

echo "Mock generation completed"
