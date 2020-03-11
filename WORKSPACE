workspace(
    name = "pipe",
    managed_directories = {"@npm": ["node_modules"]},
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
    sha256 = "221cf4a69eb30aa9b05a87b585f7b06a44cf3caff782e486a344fce8251dd4e3",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.21.4/rules_go-v0.21.4.tar.gz",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.21.4/rules_go-v0.21.4.tar.gz",
    ],
)

load(
    "@io_bazel_rules_go//go:deps.bzl",
    "go_register_toolchains",
    "go_rules_dependencies",
)

go_rules_dependencies()

go_register_toolchains(
    go_version = "1.14",
)

load(
    "@io_bazel_rules_go//extras:embed_data_deps.bzl",
    "go_embed_data_dependencies",
)

go_embed_data_dependencies()

http_archive(
    name = "bazel_gazelle",
    sha256 = "d8c45ee70ec39a57e7a05e5027c32b1576cc7f16d9dd37135b0eddde45cf1b10",
    urls = [
        "https://storage.googleapis.com/bazel-mirror/github.com/bazelbuild/bazel-gazelle/releases/download/v0.20.0/bazel-gazelle-v0.20.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.20.0/bazel-gazelle-v0.20.0.tar.gz",
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

### Protoc-gen-validate
git_repository(
    name = "com_github_envoyproxy_protoc_gen_validate",
    commit = "9eff07ddfcb4001aa1aab280648153f46e1a8ddc",
    remote = "https://github.com/envoyproxy/protoc-gen-validate.git",
    shallow_since = "1560436592 +0000",
)

# gazelle:repository_macro repositories.bzl%go_repositories

### Load dependencies.
load("//:repositories.bzl", "go_repositories")

go_repositories()
