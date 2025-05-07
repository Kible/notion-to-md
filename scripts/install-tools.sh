#!/bin/bash
set -e

echo "Installing development tools..."

# Get GOBIN path
GOBIN=$(go env GOPATH)/bin
echo "Tools will be installed to $GOBIN"

# Install golangci-lint
echo "Installing golangci-lint..."
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Install mockgen
echo "Installing mockgen..."
go install github.com/golang/mock/mockgen@latest

# Install go-mod-outdated
echo "Installing go-mod-outdated..."
go install github.com/psampaz/go-mod-outdated@latest

# Install gosec
echo "Installing gosec..."
go install github.com/securego/gosec/v2/cmd/gosec@latest

# Install nancy
echo "Installing nancy..."
go install github.com/sonatype-nexus-community/nancy@latest

# Install genqlient (GraphQL client generator)
echo "Installing genqlient..."
go install github.com/Khan/genqlient@latest

# Install gqlgen (GraphQL server generator)
echo "Installing gqlgen..."
go install github.com/99designs/gqlgen@latest

# Install cff (code format fixer)
echo "Installing cff..."
go install github.com/fwojciec/cff@latest

# Install oapi-codegen (OpenAPI generator)
echo "Installing oapi-codegen..."
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

echo "All tools have been installed!"
