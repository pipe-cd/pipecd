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

BASE_URL="https://get.helm.sh"
# Do not forget to update the version number at the following file when changing this.
# https://github.com/pipe-cd/pipecd/blob/master/pkg/app/piped/toolregistry/install.go
VERSION="3.8.2"

echo "Installing helm-${VERSION} into ${PIPED_TOOLS_DIR}/helm..."
curl -L ${BASE_URL}/helm-v${VERSION}-linux-amd64.tar.gz | tar xvz
mv linux-amd64/helm ${PIPED_TOOLS_DIR}/helm
chmod +x ${PIPED_TOOLS_DIR}/helm
rm -rf linux-amd64
echo "Successfully installed helm-${VERSION} into ${PIPED_TOOLS_DIR}/helm"
