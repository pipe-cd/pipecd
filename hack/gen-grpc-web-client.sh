#!/usr/bin/env bash

# Copyright 2020 The PipeCD Authors.
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

ROOT=$(dirname ${BASH_SOURCE})/..
PROTO_FILE_DIR=${ROOT}/pkg/app/server/service
OUTPUT_DIR=${ROOT}/pkg/app/web/src/service/
mkdir -p ${OUTPUT_DIR}
protoc -I=${PROTO_FILE_DIR} service.proto \
  --js_out=import_style=commonjs:${OUTPUT_DIR} \
  --grpc-web_out=import_style=typescript,mode=grpcweb:${OUTPUT_DIR}
