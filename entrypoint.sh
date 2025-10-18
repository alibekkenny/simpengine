#!/bin/sh
set -e

echo "Running migrations..."
migrate -path ./migrations -database "$DSN" up

echo "Starting server..."
./main