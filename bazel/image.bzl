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

load("@io_bazel_rules_docker//go:image.bzl", "go_image")
load("@io_bazel_rules_docker//container:container.bzl", "container_bundle")
load("@io_bazel_rules_docker//contrib:push-all.bzl", "docker_push")

def app_image(name, binary, repository, base = None, **kwargs):
    go_image(
        name = "image",
        binary = binary,
        base = base,
        **kwargs
    )

    container_bundle(
        name = "bundle",
        images = {
            "$(DOCKER_REGISTRY)/%s:{STABLE_GIT_COMMIT}" % repository: ":image",
            "$(DOCKER_REGISTRY)/%s:{STABLE_GIT_COMMIT_FULL}" % repository: ":image",
        },
    )

    docker_push(
        name = "push",
        bundle = ":bundle",
        **kwargs
    )
