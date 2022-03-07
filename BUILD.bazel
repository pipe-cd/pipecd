load("@bazel_gazelle//:def.bzl", "gazelle")

# gazelle:exclude docs
# gazelle:exclude hack
# gazelle:exclude dockers
# gazelle:exclude manifests
# gazelle:exclude terraform
# gazelle:exclude template
# gazelle:exclude vendor
# gazelle:exclude pkg/app/web/node_modules
# gazelle:exclude pkg/plugin/golinter/gofmt/testdata
# gazelle:exclude pkg/app/kapetool/cmd/godifflinter/pkg/linters/unusedparam/testdata
# gazelle:exclude pkg/app/kapetool/cmd/godifflinter/pkg/linters/ineffassign/testdata

# gazelle:build_file_name BUILD.bazel
# gazelle:prefix github.com/pipe-cd/pipecd

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

# gazelle:exclude pkg/model/*.proto
# gazelle:exclude pkg/app/server/service/webservice/*.proto
# gazelle:exclude pkg/app/server/service/pipedservice/*.proto
# gazelle:exclude pkg/app/server/service/apiservice/*.proto
# gazelle:exclude pkg/app/helloworld/service/*.proto
