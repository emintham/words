#!/bin/bash

# Build and start the API server
echo "Building API server..."
go build -o api cmd/api/main.go

if [ ! -f "api" ]; then
    echo "Error: Failed to build API server"
    exit 1
fi

echo "Starting API server on port 8080..."
./api
