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

BASE_URL="https://storage.googleapis.com/kubernetes-release/release"

echo "Installing various versions of kubectl..."

for version in "1.16.9" "1.17.8" "1.18.2"
do
  echo "Installing kubectl-${version} into ${PIPED_BIN_DIR}/kubectl-${version}..."
  curl -LO ${BASE_URL}/v${version}/bin/linux/amd64/kubectl
  mv kubectl ${PIPED_BIN_DIR}/kubectl-${version}
  chmod +x ${PIPED_BIN_DIR}/kubectl-${version}
  echo "Successfully installed kubectl-${version} into ${PIPED_BIN_DIR}/kubectl-${version}..."
done

cp ${PIPED_BIN_DIR}/kubectl-1.18.2 ${PIPED_BIN_DIR}/kubectl
echo "Successfully linked kubectl to kubectl-1.18.2"
