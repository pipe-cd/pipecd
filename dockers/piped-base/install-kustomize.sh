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

echo "Installing various versions of kustomize..."

for version in "3.4.0" "3.5.4" "3.5.5"
do
  echo "Installing kustomize-${version} into ${PIPED_BIN_DIR}/kustomize-${version}..."
  curl -L ${BASE_URL}/v${version}/kustomize_v${version}_linux_amd64.tar.gz | tar xvz
  mv kustomize ${PIPED_BIN_DIR}/kustomize-${version}
  chmod +x ${PIPED_BIN_DIR}/kustomize-${version}
  echo "Successfully installed kustomize-${version} into ${PIPED_BIN_DIR}/kustomize-${version}..."
done

cp ${PIPED_BIN_DIR}/kustomize-3.5.5 ${PIPED_BIN_DIR}/kustomize
echo "Successfully linked kustomize to kustomize-3.5.5"
