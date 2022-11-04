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

# parse params
while [[ -z "$1" ]]
do
  echo "Missing docs version..."
  exit 1
done

echo "Prepare version docs"

CONTENT_DIR=docs/content/en

# Update $CONTENT_DIR/docs
rm -rf $CONTENT_DIR/docs
cp -rf $CONTENT_DIR/docs-dev $CONTENT_DIR/docs
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

# Create new $CONTENT_DIR/docs-$1
rm -rf $CONTENT_DIR/docs-$1
cp -rf $CONTENT_DIR/docs-dev $CONTENT_DIR/docs-$1
cp -rf docs/themes/docsy/layouts/docs/ docs/layouts/docs-$1
cat <<EOT > $CONTENT_DIR/docs-$1/_index.md
---
title: "Welcome to PipeCD"
linkTitle: "Documentation [$1]"
type: docs
---
EOT

# Check whether docs/config.toml file contains route for new version docs /docs-$1/
# If it contained already, skip updating docs/config.toml
grep "/docs-$1/" docs/config.toml > /dev/null
if [ $? -eq 0 ]
then
  echo "Version docs has been prepared successfully at $CONTENT_DIR/docs-$1/"
  exit 0
fi

# Update docs/config.toml
tail -r docs/config.toml | tail -n +5 | tail -r >> docs/config.toml.tmp
cat <<EOT >> docs/config.toml.tmp

[[params.versions]]
  version = "$1"
  url = "/docs-$1/"
EOT
tail -4 docs/config.toml >> docs/config.toml.tmp
mv docs/config.toml.tmp docs/config.toml

echo "Version docs has been prepared successfully at $CONTENT_DIR/docs-$1/"
