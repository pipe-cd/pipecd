#!/usr/bin/env bash

# Copyright 2024 The PipeCD Authors.
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

set -o errexit
set -o nounset
set -o pipefail

# Check if rebuild flag is provided
REBUILD=""
if [[ "${1:-}" == "--rebuild" ]]; then
    REBUILD="true"
fi

echo "Generating code using local Docker image..."

# Check if Docker is running
if ! docker version >/dev/null 2>&1; then
    echo "Error: Docker is not running. Please start Docker and try again."
    exit 1
fi

# If rebuild is requested or image doesn't exist, build it
CODEGEN_IMAGE_TAG="pipecd-local-codegen:latest"
if [[ "$REBUILD" == "true" ]] || ! docker image inspect "$CODEGEN_IMAGE_TAG" >/dev/null 2>&1; then
    echo "Building codegen Docker image..."
    make build/codegen
fi

echo "Running code generation..."
make gen/code-local

echo "Code generation completed successfully."
