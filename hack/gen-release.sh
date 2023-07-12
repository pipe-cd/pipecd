#!/usr/bin/env bash

# Copyright 2023 The PipeCD Authors.
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

# parse params
while [[ -z "$1" ]]
do
    echo "Missing release version..."
    exit 1
done

echo "Prepare release for version $1"

# Update release file
RELEASE_FILE=RELEASE
sed -i '' "s/tag:.*/tag: $1/" $RELEASE_FILE
