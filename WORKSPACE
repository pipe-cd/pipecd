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
    sha256 = "2b1641428dff9018f9e85c0384f03ec6c10660d935b750e3fa1492a281a53b0f",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.29.0/rules_go-v0.29.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.29.0/rules_go-v0.29.0.zip",
    ],
)

load(
    "@io_bazel_rules_go//go:deps.bzl",
    "go_register_toolchains",
    "go_rules_dependencies",
)

go_rules_dependencies()

go_register_toolchains(
    go_version = "1.17.2",
)

load(
    "@io_bazel_rules_go//extras:embed_data_deps.bzl",
    "go_embed_data_dependencies",
)

go_embed_data_dependencies()

http_archive(
    name = "bazel_gazelle",
    sha256 = "de69a09dc70417580aabf20a28619bb3ef60d038470c7cf8442fafcf627c21cb",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.24.0/bazel-gazelle-v0.24.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.24.0/bazel-gazelle-v0.24.0.tar.gz",
    ],
)

### Protoc-gen-validate
git_repository(
    name = "com_github_envoyproxy_protoc_gen_validate",
    commit = "9eff07ddfcb4001aa1aab280648153f46e1a8ddc",
    remote = "https://github.com/envoyproxy/protoc-gen-validate.git",
    shallow_since = "1560436592 +0000",
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
    sha256 = "d0f5f605d0d656007ce6c8b5a82df3037e1d8fe8b121ed42e536f569dec16113",
    strip_prefix = "protobuf-3.14.0",
    urls = [
        "https://mirror.bazel.build/github.com/protocolbuffers/protobuf/archive/v3.14.0.tar.gz",
        "https://github.com/protocolbuffers/protobuf/archive/v3.14.0.tar.gz",
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
    strip_prefix = "buildtools-3.5.0",
    url = "https://github.com/bazelbuild/buildtools/archive/3.5.0.zip",
)

load(
    "@com_github_bazelbuild_buildtools//buildifier:deps.bzl",
    "buildifier_dependencies",
)

buildifier_dependencies()

### Docker
http_archive(
    name = "io_bazel_rules_docker",
    sha256 = "4521794f0fba2e20f3bf15846ab5e01d5332e587e9ce81629c7f96c793bb7036",
    strip_prefix = "rules_docker-0.14.4",
    urls = ["https://github.com/bazelbuild/rules_docker/releases/download/v0.14.4/rules_docker-v0.14.4.tar.gz"],
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
    "@io_bazel_rules_docker//repositories:pip_repositories.bzl",
    "pip_deps",
)

pip_deps()

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
    name = "debug-base",
    digest = "sha256:b0ec52fbde95be09074badc8298b6e94d61a9066e9637d75610267f1646fb0a1",
    registry = "gcr.io",
    repository = "pipecd/debug-base",
    tag = "0.0.1",
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
    sha256 = "55a25a762fcf9c9b88ab54436581e671bc9f4f523cb5a1bd32459ebec7be68a8",
    urls = ["https://github.com/bazelbuild/rules_nodejs/releases/download/3.2.2/rules_nodejs-3.2.2.tar.gz"],
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

load("@npm//@bazel/labs:package.bzl", "npm_bazel_labs_dependencies")

npm_bazel_labs_dependencies()

# gazelle:repository_macro repositories.bzl%go_repositories

