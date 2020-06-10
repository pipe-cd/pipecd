workspace(
    name = "pipe",
    managed_directories = {"@npm": ["pkg/app/web/node_modules"]},
)

load(
    "@bazel_tools//tools/build_defs/repo:http.bzl",
    "http_archive",
)
load(
    "@bazel_tools//tools/build_defs/repo:git.bzl",
    "git_repository",
)

### Rules_go and gazelle
http_archive(
    name = "io_bazel_rules_go",
    sha256 = "a8d6b1b354d371a646d2f7927319974e0f9e52f73a2452d2b3877118169eb6bb",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.23.3/rules_go-v0.23.3.tar.gz",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.23.3/rules_go-v0.23.3.tar.gz",
    ],
)

load(
    "@io_bazel_rules_go//go:deps.bzl",
    "go_register_toolchains",
    "go_rules_dependencies",
)

go_rules_dependencies()

go_register_toolchains(
    go_version = "1.14.4",
)

load(
    "@io_bazel_rules_go//extras:embed_data_deps.bzl",
    "go_embed_data_dependencies",
)

go_embed_data_dependencies()

http_archive(
    name = "bazel_gazelle",
    sha256 = "cdb02a887a7187ea4d5a27452311a75ed8637379a1287d8eeb952138ea485f7d",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.21.1/bazel-gazelle-v0.21.1.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.21.1/bazel-gazelle-v0.21.1.tar.gz",
    ],
)

load(
    "@bazel_gazelle//:deps.bzl",
    "gazelle_dependencies",
    "go_repository",
)

gazelle_dependencies()

### Google Protobuf
git_repository(
    name = "com_google_protobuf",
    commit = "d09d649aea36f02c03f8396ba39a8d4db8a607e4",
    remote = "https://github.com/protocolbuffers/protobuf",
    shallow_since = "1571943965 -0700",
)

load(
    "@com_google_protobuf//:protobuf_deps.bzl",
    "protobuf_deps",
)

protobuf_deps()

### BuildTools
http_archive(
    name = "com_github_bazelbuild_buildtools",
    strip_prefix = "buildtools-0.29.0",
    url = "https://github.com/bazelbuild/buildtools/archive/0.29.0.zip",
)

load(
    "@com_github_bazelbuild_buildtools//buildifier:deps.bzl",
    "buildifier_dependencies",
)

buildifier_dependencies()

### Docker
http_archive(
    name = "io_bazel_rules_docker",
    sha256 = "14ac30773fdb393ddec90e158c9ec7ebb3f8a4fd533ec2abbfd8789ad81a284b",
    strip_prefix = "rules_docker-0.12.1",
    urls = ["https://github.com/bazelbuild/rules_docker/releases/download/v0.12.1/rules_docker-v0.12.1.tar.gz"],
)

load(
    "@io_bazel_rules_docker//repositories:repositories.bzl",
    container_repositories = "repositories",
)

container_repositories()

load(
    "@io_bazel_rules_docker//go:image.bzl",
    _go_image_repos = "repositories",
)

_go_image_repos()

load(
    "@io_bazel_rules_docker//container:container.bzl",
    "container_pull",
)

container_pull(
    name = "piped-base",
    digest = "sha256:561e2e38fb33f2b7800759f9c3cd065ea79bc5b906c25e45d7e527b8e0320ae8",
    registry = "gcr.io",
    repository = "kapetanios/piped-base",
    tag = "0.0.2",
)

container_pull(
    name = "debug-base",
    digest = "sha256:7bc48b9cef27dffaa66617a39f649b6f4b3359e40fc2cac6a51b4879439aa638",
    registry = "gcr.io",
    repository = "kapetanios/debug-base",
    tag = "0.0.1",
)

### Protoc-gen-validate
git_repository(
    name = "com_github_envoyproxy_protoc_gen_validate",
    commit = "9eff07ddfcb4001aa1aab280648153f46e1a8ddc",
    remote = "https://github.com/envoyproxy/protoc-gen-validate.git",
    shallow_since = "1560436592 +0000",
)

### web

http_archive(
    name = "build_bazel_rules_nodejs",
    sha256 = "d14076339deb08e5460c221fae5c5e9605d2ef4848eee1f0c81c9ffdc1ab31c1",
    urls = ["https://github.com/bazelbuild/rules_nodejs/releases/download/1.6.1/rules_nodejs-1.6.1.tar.gz"],
)

load(
    "@build_bazel_rules_nodejs//:index.bzl",
    "node_repositories",
    "yarn_install",
)

### https://bazelbuild.github.io/rules_nodejs/Built-ins.html#usage
node_repositories(
    node_version = "12.13.0",
    package_json = ["//pkg/app/web:package.json"],
    yarn_version = "1.22.4",
)

yarn_install(
    name = "npm",
    package_json = "//pkg/app/web:package.json",
    yarn_lock = "//pkg/app/web:yarn.lock",
)

load("@npm//:install_bazel_dependencies.bzl", "install_bazel_dependencies")

install_bazel_dependencies()

load("@npm_bazel_labs//:package.bzl", "npm_bazel_labs_dependencies")

npm_bazel_labs_dependencies()

load("@npm_bazel_typescript//:index.bzl", "ts_setup_workspace")

ts_setup_workspace()

# gazelle:repository_macro repositories.bzl%go_repositories

### Load dependencies.
load("//:repositories.bzl", "go_repositories")

go_repositories()
