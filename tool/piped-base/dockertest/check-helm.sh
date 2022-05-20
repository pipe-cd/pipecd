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
#pathcases["helm-2.16.7"]="/home/piped/.piped/tools/helm-2.16.7"
pathcases["helm"]="/home/piped/.piped/tools/helm"

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
#versioncases["helm-2.16.7"]="Client: \&version.Version{SemVer:\"v2.16.7\", GitCommit:\"5f2584fd3d35552c4af26036f0c464191287986b\", GitTreeState:\"clean\"}"
versioncases["helm"]="version.BuildInfo{Version:\"v3.8.2\", GitCommit:\"6e3701edea09e5d55a8ca2aae03a68917630e91b\", GitTreeState:\"clean\", GoVersion:\"go1.17.5\"}"

for h in "${!versioncases[@]}"
do
    chmod +r -R $HOME
    got=$($h version --client --registry-config ${HELM_REGISTRY_CONFIG})
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
