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
    sha256 = "b725e6497741d7fc2d55fcc29a276627d10e43fa5d0bb692692890ae30d98d00",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.24.3/rules_go-v0.24.3.tar.gz",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.24.3/rules_go-v0.24.3.tar.gz",
    ],
)

load(
    "@io_bazel_rules_go//go:deps.bzl",
    "go_register_toolchains",
    "go_rules_dependencies",
)

go_rules_dependencies()

go_register_toolchains(
    go_version = "1.14.7",
)

load(
    "@io_bazel_rules_go//extras:embed_data_deps.bzl",
    "go_embed_data_dependencies",
)

go_embed_data_dependencies()

http_archive(
    name = "bazel_gazelle",
    sha256 = "72d339ff874a382f819aaea80669be049069f502d6c726a07759fdca99653c48",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.22.1/bazel-gazelle-v0.22.1.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.22.1/bazel-gazelle-v0.22.1.tar.gz",
    ],
)

load(
    "@bazel_gazelle//:deps.bzl",
    "gazelle_dependencies",
    "go_repository",
)

gazelle_dependencies()

### Google Protobuf
http_archive(
    name = "com_google_protobuf",
    sha256 = "1c744a6a1f2c901e68c5521bc275e22bdc66256eeb605c2781923365b7087e5f",
    strip_prefix = "protobuf-3.13.0",
    urls = ["https://github.com/protocolbuffers/protobuf/archive/v3.13.0.zip"],
)

load(
    "@com_google_protobuf//:protobuf_deps.bzl",
    "protobuf_deps",
)

protobuf_deps()

### BuildTools
http_archive(
    name = "com_github_bazelbuild_buildtools",
    strip_prefix = "buildtools-3.4.0",
    url = "https://github.com/bazelbuild/buildtools/archive/3.4.0.zip",
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
    digest = "sha256:e357a240b1e8668c851a1dcb3e3244b4e037aab56029f394dda6925070d473a9",
    registry = "gcr.io",
    repository = "pipecd/piped-base",
    tag = "0.1.0",
)

container_pull(
    name = "debug-base",
    digest = "sha256:b0ec52fbde95be09074badc8298b6e94d61a9066e9637d75610267f1646fb0a1",
    registry = "gcr.io",
    repository = "pipecd/debug-base",
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
