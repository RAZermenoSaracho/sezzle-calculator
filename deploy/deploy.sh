#!/usr/bin/env bash
# Builds the current commit inside the shared Docker VM (via Vagrant) and
# replaces the sezzle-calculator container. Docker is never used on the host.
#
# Environment:
#   DOCKER_VM_DIR  Vagrant directory of the Docker VM (default: ~/scripts/docker-test-vm)
#   HOST_PORT      Port published inside the VM (default: 8080)
set -euo pipefail

DOCKER_VM_DIR="${DOCKER_VM_DIR:-$HOME/scripts/docker-test-vm}"
HOST_PORT="${HOST_PORT:-8080}"
CONTAINER_NAME="sezzle-calculator"
IMAGE="sezzle-calculator"
BUILD_DIR="/tmp/sezzle-calculator-build"

if [[ ! -f "$DOCKER_VM_DIR/Vagrantfile" ]]; then
  echo "No Vagrantfile in DOCKER_VM_DIR: $DOCKER_VM_DIR" >&2
  exit 1
fi
if [[ ! "$HOST_PORT" =~ ^[0-9]+$ ]]; then
  echo "HOST_PORT must be a number, got: $HOST_PORT" >&2
  exit 1
fi
if ! command -v vagrant >/dev/null; then
  echo "vagrant is not on PATH" >&2
  exit 1
fi

repo_root="$(git rev-parse --show-toplevel)"
sha="$(git -C "$repo_root" rev-parse --short HEAD)"

echo "Deploying $sha to $CONTAINER_NAME (port $HOST_PORT)"

# Ship the committed tree into the VM without relying on shared folders.
git -C "$repo_root" archive HEAD | (
  cd "$DOCKER_VM_DIR"
  vagrant ssh -c "rm -rf $BUILD_DIR && mkdir -p $BUILD_DIR && tar -x -C $BUILD_DIR" -- -T
)

# Only the container named $CONTAINER_NAME is touched; other workloads are not.
(
  cd "$DOCKER_VM_DIR"
  vagrant ssh -c "set -eu
    cd $BUILD_DIR
    docker build -t $IMAGE:$sha -t $IMAGE:latest .
    docker rm -f $CONTAINER_NAME >/dev/null 2>&1 || true
    docker run -d --restart unless-stopped --name $CONTAINER_NAME -p $HOST_PORT:8080 $IMAGE:$sha
    rm -rf $BUILD_DIR" -- -T
)

echo "Deployed $sha"
