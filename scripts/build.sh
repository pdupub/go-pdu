#!/bin/bash
set -e

echo "Building project..."
go build -o build/pdu ./cmd/pdu
echo "Build completed."

