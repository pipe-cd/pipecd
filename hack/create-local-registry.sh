#!/usr/bin/env bash

# Copyright 2025 The PipeCD Authors.
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

REG_NAME='kind-registry'
REG_PORT='5001'

# Create registry container unless it already exists
echo "Creating local registry container..."
running="$(docker inspect -f '{{.State.Running}}' "${REG_NAME}" 2>/dev/null || true)"
if [ "${running}" != 'true' ]; then
  docker run \
    -e REGISTRY_HTTP_ADDR=0.0.0.0:5001 \
    -d --restart=always -p "127.0.0.1:${REG_PORT}:5001" --name "${REG_NAME}" \
    registry:2
fi
