#!/bin/bash
# Build the Docker image and run the apm-add integration test.
# Usage: ./tests/integration/run.sh [--no-cache]
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
IMAGE="apm-integration-test"
NO_CACHE="${1:-}"

if ! command -v docker &>/dev/null; then
    echo "error: docker not found in PATH" >&2
    echo "  On WSL: enable Docker Desktop WSL integration at" >&2
    echo "  Docker Desktop → Settings → Resources → WSL Integration" >&2
    exit 1
fi

echo "==> Building Docker image ($IMAGE)..."
BUILD_ARGS=()
[ "$NO_CACHE" = "--no-cache" ] && BUILD_ARGS+=(--no-cache)
docker build "${BUILD_ARGS[@]}" -t "$IMAGE" "$SCRIPT_DIR"

echo
echo "==> Running integration test..."
docker run --rm "$IMAGE"

echo
echo "==> Cleaning up image..."
docker rmi "$IMAGE" &>/dev/null || true
