#!/usr/bin/env bash

# Copyright 2022 The PipeCD Authors.
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

echo "Updating the contributors list on README.md ..."

LINE_NUM=$(($(grep -Fn "### Thanks to the contributors of PipeCD" README.md | cut -f1 -d ':')+1))
head -n $LINE_NUM README.md >> README.md.tmp

while read -r line
do
    cat <<EOT >> README.md.tmp
$line
EOT
done < <(gh api -XGET /repos/pipe-cd/pipecd/contributors -F per_page=100 | jq -r '.[] | "<a href=\"\(.html_url)\"><img src=\"\(.avatar_url)\" title=\"\(.login)\" width=\"80\" height=\"80\"></a>"')

mv README.md.tmp README.md

echo "Successfully update the contributions list on README.md"
