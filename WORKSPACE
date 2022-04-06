workspace(
    name = "pipecd",
    managed_directories = {"@npm": ["web/node_modules"]},
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
    sha256 = "d6b2513456fe2229811da7eb67a444be7785f5323c6708b38d851d2b51e54d83",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.30.0/rules_go-v0.30.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.30.0/rules_go-v0.30.0.zip",
    ],
)

load(
    "@io_bazel_rules_go//go:deps.bzl",
    "go_register_toolchains",
    "go_rules_dependencies",
)

go_rules_dependencies()

go_register_toolchains(
    go_version = "1.17.6",
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "de69a09dc70417580aabf20a28619bb3ef60d038470c7cf8442fafcf627c21cb",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.24.0/bazel-gazelle-v0.24.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.24.0/bazel-gazelle-v0.24.0.tar.gz",
    ],
)

### Load dependencies.
load("//:repositories.bzl", "go_repositories")

go_repositories()

load(
    "@bazel_gazelle//:deps.bzl",
    "gazelle_dependencies",
)

gazelle_dependencies()

### Google Protobuf
http_archive(
    name = "com_google_protobuf",
    sha256 = "3bd7828aa5af4b13b99c191e8b1e884ebfa9ad371b0ce264605d347f135d2568",
    strip_prefix = "protobuf-3.19.4",
    urls = [
        "https://mirror.bazel.build/github.com/protocolbuffers/protobuf/archive/v3.19.4.tar.gz",
        "https://github.com/protocolbuffers/protobuf/archive/v3.19.4.tar.gz",
    ],
)

load(
    "@com_google_protobuf//:protobuf_deps.bzl",
    "protobuf_deps",
)

protobuf_deps()

### BuildTools
http_archive(
    name = "com_github_bazelbuild_buildtools",
    strip_prefix = "buildtools-4.2.3",
    url = "https://github.com/bazelbuild/buildtools/archive/4.2.3.zip",
)

load(
    "@com_github_bazelbuild_buildtools//buildifier:deps.bzl",
    "buildifier_dependencies",
)

buildifier_dependencies()

### Docker
http_archive(
    name = "io_bazel_rules_docker",
    sha256 = "59536e6ae64359b716ba9c46c39183403b01eabfbd57578e84398b4829ca499a",
    strip_prefix = "rules_docker-0.22.0",
    urls = ["https://github.com/bazelbuild/rules_docker/releases/download/v0.22.0/rules_docker-v0.22.0.tar.gz"],
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
    "@io_bazel_rules_docker//repositories:deps.bzl",
    container_deps = "deps",
)

container_deps()

load(
    "@io_bazel_rules_docker//container:container.bzl",
    "container_pull",
)

container_pull(
    name = "piped-base",
    digest = "sha256:be303a0bc87480a26ee90c91288d47498e5742e90e7803cc4f2e11bfcbffb118",
    registry = "gcr.io",
    repository = "pipecd/piped-base",
    tag = "0.2.2",
)

container_pull(
    name = "piped-base-okd",
    digest = "sha256:54f11a2701a5ad8c9d9fbf1f1c3232fa02f30c4fa399c98e7c2df1640fdb4f0d",
    registry = "gcr.io",
    repository = "pipecd/piped-base-okd",
    tag = "0.1.0",
)

container_pull(
    name = "pipectl-base",
    digest = "sha256:0cf7eacedb0cc8d759248f0e25bd8eddf659de6f2c1db315ac95a272ec2e60cc",
    registry = "gcr.io",
    repository = "pipecd/pipectl-base",
    tag = "0.2.0",
)

container_pull(
    name = "pipecd-base",
    digest = "sha256:f3e98a27b85b8ead610c4f93cec8d936c760a43866cf817d32563daf9b198358",
    registry = "gcr.io",
    repository = "pipecd/pipecd-base",
    tag = "0.1.0",
)

### web

http_archive(
    name = "build_bazel_rules_nodejs",
    sha256 = "cfc289523cf1594598215901154a6c2515e8bf3671fd708264a6f6aefe02bf39",
    urls = ["https://github.com/bazelbuild/rules_nodejs/releases/download/4.4.6/rules_nodejs-4.4.6.tar.gz"],
)

load(
    "@build_bazel_rules_nodejs//:index.bzl",
    "node_repositories",
    "yarn_install",
)

### https://bazelbuild.github.io/rules_nodejs/Built-ins.html#usage
node_repositories(
    node_version = "16.13.0",
    package_json = ["//web:package.json"],
    yarn_version = "1.22.4",
)

yarn_install(
    name = "npm",
    frozen_lockfile = False,
    package_json = "//web:package.json",
    yarn_lock = "//web:yarn.lock",
)

load("@npm//@bazel/labs:package.bzl", "npm_bazel_labs_dependencies")

npm_bazel_labs_dependencies()

# gazelle:repository_macro repositories.bzl%go_repositories
