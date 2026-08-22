#!/usr/bin/env bash
# Copyright 2026 The PipeCD Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${REPO_ROOT}"

kind_arch() {
  case "$(uname -m)" in
    x86_64) echo "amd64" ;;
    aarch64 | arm64) echo "arm64" ;;
    *)
      echo "unsupported architecture for kind: $(uname -m)" >&2
      return 1
      ;;
  esac
}

install_kind() {
  local kind_version="v0.27.0"
  local arch
  arch="$(kind_arch)"
  local kind_bin="/tmp/kind"
  curl -fsSL "https://kind.sigs.k8s.io/dl/${kind_version}/kind-linux-${arch}" -o "${kind_bin}"
  chmod +x "${kind_bin}"
  sudo mv "${kind_bin}" /usr/local/bin/kind
}

setup_yarn() {
  if command -v yarn >/dev/null 2>&1; then
    return
  fi

  if command -v corepack >/dev/null 2>&1; then
    corepack enable
    corepack prepare yarn@1.22.22 --activate
    return
  fi

  npm install -g yarn
}

echo "Installing kind..."
install_kind

echo "Setting up Yarn..."
setup_yarn

echo "Updating Go dependencies..."
make update/go-deps

echo "Updating web dependencies..."
make update/web-deps

echo "Running smoke checks..."
docker info >/dev/null
go version
node --version
yarn --version
helm version --short
kubectl version --client=true
kind version

echo "Dev container setup complete."
