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

set -o errexit
set -o nounset
set -o pipefail

if [ "$#" -ne 2 ]; then
	echo "usage: $0 <organization> <repository>"
	exit 1
fi

OS="$(go env GOHOSTOS)"
ARCH="$(go env GOARCH)"
ROOT=$(dirname ${BASH_SOURCE})/..
BUILD_TEMPLATE_FILE=${ROOT}/BUILD.bazel.tpl
BUILD_FILE=${ROOT}/BUILD.bazel
GENERATED_BUILD_FILE=${ROOT}/BUILD.bazel.generated

ORGANIZATION=$1
REPOSITORY=$2

cat ${BUILD_TEMPLATE_FILE} > ${GENERATED_BUILD_FILE}

expose_package () {
	local out_path=$1
	local package=$2
	local old_links=$(eval echo \$$3)
	local generated_files=$(eval echo \$$4)

	# Compute the relative_path from this package to the bazel-bin.
	local count_paths="$(echo -n "${package}" | tr '/' '\n' | wc -l)"
	local relative_path=""
	for i in $(seq 0 ${count_paths}); do
		relative_path="../${relative_path}"
	done

	# Delete all old links.
	for f in ${old_links}; do
		if [[ -f "${f}" ]]; then
			echo "Deleting old link: ${f}"
			rm ${f}
		fi
	done

	# Link to the generated files and add them to excluding list in the root BUILD file.
	local found=0
	for f in ${generated_files}; do
		if [[ -f "${f}" ]]; then
			found=1
			local base=${f##*/}
			echo "Adding a new link: ${package}/${base}"
			ln -nsf "${relative_path}${f}" "${package}/"
			if [[ ${base} == *.mock.go ]] || [[ ${base} == *.pb.go ]]; then
				continue
			fi
			echo "# gazelle:exclude ${package}/${base}" >> ${GENERATED_BUILD_FILE}
		fi
	done
	if [[ "${found}" == "0" ]]; then
		echo "Error: No generated file was found inside ${out_path} for the package ${package}"
		exit 1
	fi
}

# Build all packages.
bazelisk build --noincompatible_strict_action_env -- //...

####################
# For proto go giles
####################

# Link to the generated files and add them to excluding list in the root BUILD file.
for label in $(bazelisk query 'kind(go_proto_library, //...)'); do
	package="${label%%:*}"
	package="${package##//}"
	packageName="${package##*/}"
	target="${label##*:}"
	[[ -d "${package}" ]] || continue

	# Compute the path where Bazel puts the files.
	out_path="bazel-bin/${package}/${packageName}_go_proto_/github.com/${ORGANIZATION}/${REPOSITORY}/${package}"

	old_links=$(eval echo ${package}/*{.pb.go,.pb.validate.go})
	generated_files=$(eval echo ${out_path}/*{.pb.go,.pb.validate.go})
	expose_package ${out_path} ${package} old_links generated_files
done

###################
# For mock go files
###################

# Link to the generated files and add them to excluding list in the root BUILD file.
for package in $(bazelisk query 'kind(gomock, //...)' --output package); do
	# Compute the path where Bazel puts the files.
	out_path="bazel-bin/${package}"

	old_links=${package}/*.mock.go
	generated_files=${out_path}/*.mock.go
	expose_package ${out_path} ${package} old_links generated_files
done

# Reset the root BUILD file
cat ${GENERATED_BUILD_FILE} > ${BUILD_FILE}
rm ${GENERATED_BUILD_FILE}
