#!/usr/bin/env bash

# Copyright 2021 The PipeCD Authors.
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

REGISTRY=""
VERSION=""
OUT=".rendered-manifests"

if [ $# -eq 0 ]
then
  VERSION="$(git describe --tags --always --dirty --abbrev=7)"
  REGISTRY="localhost:5001"
else
  VERSION=$1
fi

echo "Start rendering manifests for version ${VERSION}"

echo "Cleaning old ${OUT}..."
rm -rf ${OUT}

echo "Copying manifests to ${OUT}..."
cp -rf manifests ${OUT}

echo "Updating version to ${VERSION}..."
sed -i'' -e 's/{{ .VERSION }}/'"${VERSION}"'/g' ${OUT}/pipecd/Chart.yaml
sed -i'' -e 's/{{ .VERSION }}/'"${VERSION}"'/g' ${OUT}/piped/Chart.yaml

if [ ! -z "${REGISTRY}" ]
then
  echo "Updating image registry to ${REGISTRY}..."
  sed -i'' -e 's/gcr.io\/pipecd/'"${REGISTRY}"'/g' ${OUT}/pipecd/values.yaml
  sed -i'' -e 's/gcr.io\/pipecd/'"${REGISTRY}"'/g' ${OUT}/piped/values.yaml
fi

echo "Updating dependencies..."
helm dependency update ${OUT}/pipecd
