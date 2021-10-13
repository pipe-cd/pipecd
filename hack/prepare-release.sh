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
    echo "Missing release version..."
    exit 1
done

CONTENT_DIR=docs/content/en

echo "Prepare release for version $1"

# Render release note template
TEMP="---\ntitle: \"Release $1\"\nlinkTitle: \"Release $1\"\ndate: $(date +"%Y-%m-%d")\ndescription: >\n Release $1\n---\n\n"
OUTPUT_FILE=$CONTENT_DIR/blog/releases/$1.md
echo -e $TEMP > $OUTPUT_FILE

# Update release file
RELEASE_FILE=release/RELEASE
echo -e "version: $1" > $RELEASE_FILE

echo "Your new release note is located at $OUTPUT_FILE"

echo "Prepare docs for new release"

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
cp -rf $CONTENT_DIR/docs-dev $CONTENT_DIR/docs-$1
cp -rf docs/themes/docsy/layouts/docs/ docs/layouts/docs-$1
cat <<EOT > $CONTENT_DIR/docs-$1/_index.md
---
title: "Welcome to PipeCD"
linkTitle: "Documentation [$1]"
type: docs
---
EOT

# Update docs/config.toml
tail -r docs/config.toml | tail -n +5 | tail -r >> docs/config.toml.tmp
cat <<EOT >> docs/config.toml.tmp

[[params.versions]]
  version = "$1"
  url = "/docs-$1/"
EOT
tail -4 docs/config.toml >> docs/config.toml.tmp
mv docs/config.toml.tmp docs/config.toml

echo "Docs has been prepared successfully"
