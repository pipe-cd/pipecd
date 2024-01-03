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

# parse params
while [[ -z "$1" ]]
do
  echo "Missing docs version..."
  exit 1
done

v=$1
VERSION=${v%.*}.x

echo "Prepare version docs ${VERSION}"

CONTENT_DIR=docs/content/en

# Create new $CONTENT_DIR/docs-$VERSION
rm -rf $CONTENT_DIR/docs-$VERSION
cp -rf $CONTENT_DIR/docs-dev $CONTENT_DIR/docs-$VERSION
cp -rf docs/themes/docsy/layouts/docs/ docs/layouts/docs-$VERSION
cat <<EOT > $CONTENT_DIR/docs-$VERSION/_index.md
---
title: "Welcome to PipeCD"
linkTitle: "Documentation [$VERSION]"
type: docs
---
EOT

# Check whether docs/config.toml file contains route for new version docs /docs-$VERSION/
# If it contained already, skip updating docs/config.toml
grep "/docs-$VERSION/" docs/config.toml > /dev/null
if [ $? -eq 0 ]
then
  echo "Version docs has been prepared successfully at $CONTENT_DIR/docs-$VERSION/"
  exit 0
fi

# Update docs/config.toml
LINE_NUM=$(($(grep -Fn "# Append the release versions here." docs/config.toml | cut -f1 -d ':')+5))
head -n $LINE_NUM docs/config.toml >> docs/config.toml.tmp
cat <<EOT >> docs/config.toml.tmp
[[params.versions]]
  version = "$VERSION"
  url = "/docs-$VERSION/"
EOT
tail -n +$LINE_NUM docs/config.toml >> docs/config.toml.tmp
mv docs/config.toml.tmp docs/config.toml

# Update docs/main.go
sed -i '' "s/const latestPath.*/const latestPath = \"\/docs-"$VERSION"\/\"/g" docs/main.go

echo "Version docs has been prepared successfully at $CONTENT_DIR/docs-$VERSION/"
