#!/bin/bash

set -e

echo "Building Golang Service..."

# Build the application
echo "Building Go application..."
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/server cmd/server/main.go

echo "Build completed successfully!"

# Optional: Run tests
if [ "$RUN_TESTS" = "true" ]; then
    echo "Running tests..."
    go test -v ./...
fi

# Optional: Build Docker image
if [ "$BUILD_DOCKER" = "true" ]; then
    echo "Building Docker image..."
    docker build -f deployments/docker/Dockerfile -t golang-service:latest .
fi

echo "All build steps completed!"