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

load("@io_bazel_rules_docker//go:image.bzl", "go_image")
load("@io_bazel_rules_docker//container:container.bzl", "container_bundle")
load("@io_bazel_rules_docker//contrib:push-all.bzl", "docker_push")

# TODO: Use go and container rules directly in each app directory and remove this helper.
def app_image(name, binary, repository, base = None, **kwargs):
    go_image(
        name = "%s_image" % name,
        binary = binary,
        base = base,
        **kwargs
    )

    container_bundle(
        name = "%s_bundle" % name,
        images = {
            "$(DOCKER_REGISTRY)/%s:{STABLE_VERSION}" % repository: ":%s_image" % name,
        },
    )

    docker_push(
        name = "%s_push" % name,
        bundle = ":%s_bundle" % name,
        **kwargs
    )
