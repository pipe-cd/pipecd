# Copyright 2020 The Pipe Authors.
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

load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")
load("@io_bazel_rules_go//proto:compiler.bzl", "go_proto_compiler")

def pgv_go_proto_library(name, proto, deps = [], compilers = [], **kwargs):
    go_proto_compiler(
        name = "pgv_plugin_go",
        suffix = ".pb.validate.go",
        valid_archive = False,
        plugin = "@com_github_envoyproxy_protoc_gen_validate//:protoc-gen-validate",
        options = ["lang=go"],
    )

    go_proto_library(
        name = name,
        proto = proto,
        deps = [
            "@com_github_envoyproxy_protoc_gen_validate//validate:go_default_library",
            "@com_github_golang_protobuf//ptypes:go_default_library_gen",
        ] + deps,
        compilers = [
            "pgv_plugin_go",
        ] + compilers,
        **kwargs
    )
