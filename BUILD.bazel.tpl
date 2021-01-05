load("@bazel_gazelle//:def.bzl", "gazelle")

# gazelle:exclude docs
# gazelle:exclude dockers
# gazelle:exclude manifests
# gazelle:exclude terraform
# gazelle:exclude template
# gazelle:exclude vendor
# gazelle:exclude pkg/app/web/node_modules
# gazelle:exclude pkg/plugin/golinter/gofmt/testdata

# gazelle:build_file_name BUILD.bazel
# gazelle:prefix github.com/pipe-cd/pipe

gazelle(
    name = "gazelle",
)

load("@com_github_bazelbuild_buildtools//buildifier:def.bzl", "buildifier")

buildifier(
    name = "buildifier",
    exclude_patterns = [
        "./docs/*",
        "./dockers/*",
        "./manifests/*",
        "./terraform/*",
        "./template/*",
        "./vendor/*",
    ],
)

genrule(
    name = "copy_piped",
    srcs = ["//cmd/piped"],
    outs = ["piped"],
    cmd = "cp $< $@",
)

genrule(
    name = "copy_pipectl",
    srcs = ["//cmd/pipectl"],
    outs = ["pipectl"],
    cmd = "cp $< $@",
)
