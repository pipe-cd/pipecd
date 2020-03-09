#!/usr/bin/env bash

# Copyright 2020 The Dianomi Authors.
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

# Copyright 2017 The Kubernetes Authors.
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

GIT_COMMIT="$(git describe --tags --always --dirty --abbrev=7)"
GIT_COMMIT_FULL="$(git rev-parse HEAD)"
BUILD_DATE="$(date -u '+%Y%m%d')"
VERSION="${BUILD_DATE}-${GIT_COMMIT}"

cat <<EOF
STABLE_GIT_COMMIT ${GIT_COMMIT}
STABLE_GIT_COMMIT_FULL ${GIT_COMMIT_FULL}
STABLE_BUILD_DATE ${BUILD_DATE}
STABLE_VERSION ${VERSION}
gitCommit ${GIT_COMMIT}
gitCommitFull ${GIT_COMMIT_FULL}
buildDate ${BUILD_DATE}
version ${VERSION}
EOF
