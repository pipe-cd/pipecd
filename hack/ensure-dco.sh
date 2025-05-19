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

MERGE_BASE=$(git merge-base HEAD origin/master)
COMMIT_HASHES=$(git log --format=format:"%H" "$MERGE_BASE..HEAD")

# check if the commit message contains the DCO sign-off
for commit_hash in $COMMIT_HASHES; do
  sign_off=$(git log -1 --format="%(trailers:key=Signed-off-by,valueonly)%-C()" "$commit_hash")
  if [ -z "$sign_off" ]; then
    echo "Error: Commit $commit_hash is missing Signed-off-by line"
    exit 1
  fi
done

echo "All commits have the DCO sign-off"
