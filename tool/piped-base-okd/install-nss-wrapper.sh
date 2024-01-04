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

BASE_URL="https://ftp.samba.org/pub/cwrap/nss_wrapper"
VERSION="1.1.11"

echo "Installing nss_wrapper-${VERSION}..."
curl -L ${BASE_URL}-${VERSION}.tar.gz | tar xvz

cd nss_wrapper-${VERSION}
mkdir build
cd build
cmake -D CMAKE_INSTALL_PREFIX=/usr/local -D CMAKE_BUILD_TYPE=Release ..
make install
cd ../..
rm -rf nss_wrapper-${VERSION}
echo "Successfully installed nss_wrapper-${VERSION}"
