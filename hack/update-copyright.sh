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

ROOT_DIR='./'

# Detect OS for sed -i compatibility
if [[ "$OSTYPE" == "darwin"* ]]; then
  SED_IN_PLACE="sed -i ''"
else
  SED_IN_PLACE="sed -i"
fi

grep -rl "Copyright [0-9]\{4\} The PipeCD Authors" $ROOT_DIR | while IFS= read -r i; do
  echo "Updating copyright year in: $i"
  $SED_IN_PLACE "s/Copyright [0-9]\{4\} The PipeCD Authors/Copyright $(date +%Y) The PipeCD Authors/g" "$i"
done

echo "Copyright year update completed."
