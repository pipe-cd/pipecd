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

# parse params
while [[ -z "$1" ]]
do
    echo "Missing docs version..."
    exit 1
done

echo "Sync stable docs with docs at version $1"

CONTENT_DIR=docs/content/en

rm -rf $CONTENT_DIR/docs
cp -rf $CONTENT_DIR/docs-$1 $CONTENT_DIR/docs
cat <<EOT > $CONTENT_DIR/docs/_index.md
---
title: "Welcome to PipeCD"
linkTitle: "Documentation"
weight: 1
menu:
  main:
    weight: 20
---
EOT

echo "Stable version docs has been synced successfully with docs at version $1"
