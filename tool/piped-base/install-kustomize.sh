#!/usr/bin/env bash

# Copyright 2020 The PipeCD Authors.
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

BASE_URL="https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize"
# Do not forget to update the version number at the following file when changing this.
# https://github.com/pipe-cd/pipecd/blob/master/pkg/app/piped/toolregistry/install.go
VERSION="3.8.1"

echo "Installing kustomize-${VERSION} into ${PIPED_TOOLS_DIR}/kustomize..."
curl -L ${BASE_URL}/v${VERSION}/kustomize_v${VERSION}_linux_amd64.tar.gz | tar xvz
mv kustomize ${PIPED_TOOLS_DIR}/kustomize
chmod +x ${PIPED_TOOLS_DIR}/kustomize
echo "Successfully installed kustomize-${VERSION} into ${PIPED_TOOLS_DIR}/kustomize..."
