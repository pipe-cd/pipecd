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

declare -A pathcases
pathcases["helm-2.16.3"]="/usr/local/piped/helm-2.16.3"
pathcases["helm-2.16"]="/usr/local/piped/helm-2.16"
pathcases["helm-3.1.1"]="/usr/local/piped/helm-3.1.1"
pathcases["helm-3.1"]="/usr/local/piped/helm-3.1"

for h in "${!pathcases[@]}"
do
    got=$(which $h)
    want=${pathcases[$h]}
    if [[ ${got} == ${want} ]]; then
        echo -e "PASSED: Correct path for ${h}."
        echo "  want: ${want}"
        echo "  got : ${got}"
    else
        echo "FAILED: Wrong path for ${h}."
        echo "  want: ${want}"
        echo "  got : ${got}"
        exit 1
    fi
done

declare -A versioncases
versioncases["helm-2.16.3"]="Client: \&version.Version{SemVer:\"v2.16.3\", GitCommit:\"1ee0254c86d4ed6887327dabed7aa7da29d7eb0d\", GitTreeState:\"clean\"}"
versioncases["helm-2.16"]="Client: \&version.Version{SemVer:\"v2.16.3\", GitCommit:\"1ee0254c86d4ed6887327dabed7aa7da29d7eb0d\", GitTreeState:\"clean\"}"
versioncases["helm-3.1.1"]="version.BuildInfo{Version:\"v3.1.1\", GitCommit:\"afe70585407b420d0097d07b21c47dc511525ac8\", GitTreeState:\"clean\", GoVersion:\"go1.13.8\"}"
versioncases["helm-3.1"]="version.BuildInfo{Version:\"v3.1.1\", GitCommit:\"afe70585407b420d0097d07b21c47dc511525ac8\", GitTreeState:\"clean\", GoVersion:\"go1.13.8\"}"

for h in "${!versioncases[@]}"
do
    got=$($h version --client)
    want=${versioncases[$h]}
    if [[ ${got} == ${want} ]]; then
        echo -e "PASSED: Correct version for ${h}."
        echo "  want: ${want}"
        echo "  got : ${got}"
    else
        echo "FAILED: Wrong version for ${h}."
        echo "  want: ${want}"
        echo "  got : ${got}"
        exit 1
    fi
done
