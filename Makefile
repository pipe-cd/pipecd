BAZEL_FLAGS =
BAZEL_COMMAND_FLAGS =
CLOUDSDK_PYTHON = "/usr/bin/python"

ifdef EXTENDED_BAZEL_RC
	BAZEL_FLAGS += --bazelrc=${EXTENDED_BAZEL_RC}
endif

ifdef IS_CI
	BAZEL_FLAGS += --output_base=/workspace/bazel_out
	BAZEL_COMMAND_FLAGS += --config=ci
endif

ifdef BUILD_PLATFORM
	BAZEL_COMMAND_FLAGS += --config=${BUILD_PLATFORM}
endif

ifdef PYTHON_PATH
	CLOUDSDK_PYTHON=${PYTHON_PATH}
endif

.PHONY: build
build:
	bazelisk ${BAZEL_FLAGS} build ${BAZEL_COMMAND_FLAGS} -- //...

.PHONY: build-images
build-images:
	bazelisk ${BAZEL_FLAGS} build ${BAZEL_COMMAND_FLAGS} --config=linux --config=stamping -- //cmd/...

.PHONY: push
push:
	bazelisk ${BAZEL_FLAGS} run ${BAZEL_COMMAND_FLAGS} --config=linux --config=stamping -- //cmd/pipecd:pipecd_app_push
	bazelisk ${BAZEL_FLAGS} run ${BAZEL_COMMAND_FLAGS} --config=linux --config=stamping -- //cmd/piped:piped_app_push

.PHONY: render-manifests
render-manifests:
	./hack/render-manifests.sh $(VERSION)

.PHONY: load-piped-image
load-piped-image:
	bazelisk ${BAZEL_FLAGS} run ${BAZEL_COMMAND_FLAGS} --config=linux --config=stamping -- //cmd/piped:piped_app_image --norun

.PHONY: test
test:
	bazelisk ${BAZEL_FLAGS} test ${BAZEL_COMMAND_FLAGS} -- //pkg/...

.PHONY: test-debug
test-debug:
	bazelisk ${BAZEL_FLAGS} test ${BAZEL_COMMAND_FLAGS} --test_output=all -- //pkg/...

.PHONY: test-integration
test-integration:
	bazelisk ${BAZEL_FLAGS} test ${BAZEL_COMMAND_FLAGS} --action_env=CLOUDSDK_PYTHON=${CLOUDSDK_PYTHON} -- //test/integration/...

.PHONY: coverage
coverage:
	bazelisk ${BAZEL_FLAGS} coverage ${BAZEL_COMMAND_FLAGS} -- //pkg/...

.PHONY: dep
dep:
	go mod tidy
	go mod vendor
	bazelisk run //:gazelle -- update-repos -from_file=go.mod -prune -build_file_proto_mode=disable -to_macro=repositories.bzl%go_repositories

.PHONY: gazelle
gazelle:
	bazelisk run //:gazelle

.PHONY: buildifier
buildifier:
	bazelisk run //:buildifier

.PHONY: clean
clean:
	bazelisk clean --expunge

.PHONY: site
site:
	env RELEASE=$(shell cut -c10- release/RELEASE) hugo server --source=docs

.PHONY: web-dep
web-dep:
	bazelisk build //web:build_api //web:build_model

.PHONY: web-dev
web-dev:
	cd web; yarn dev

.PHONY: web-test
web-test:
	cd web; yarn test:coverage --runInBand

.PHONY: web-lint
web-lint:
	cd web; yarn lint:fix

.PHONY: generate-test-tls
generate-test-tls:
	openssl req -x509 -nodes -days 3650 -newkey rsa:2048 \
		-keyout pkg/rpc/testdata/tls.key \
		-out pkg/rpc/testdata/tls.crt \
		-subj "/CN=localhost" \
		-config pkg/rpc/testdata/tls.config

.PHONY: kind-up
kind-up:
	./hack/create-kind-cluster.sh pipecd

.PHONY: kind-down
kind-down:
	kind delete cluster --name pipecd

.PHONY: prepare-release
prepare-release:
	./hack/prepare-release.sh $(version)

.PHONY: prepare-version-docs
prepare-version-docs:
	./hack/prepare-version-docs.sh $(version)

.PHONY: sync-stable-docs
sync-stable-docs:
	./hack/sync-stable-docs.sh $(version)

.PHONY: update-docsy
update-docsy:
	rm -rf docs/themes/docsy
	git clone --recurse-submodules --depth 1 https://github.com/google/docsy.git docs/themes/docsy

.PHONY: codegen
codegen:
	docker run --rm -v ${PWD}:/repo -it gcr.io/pipecd/codegen:0.7.0 /repo

.PHONY: build-frontend
build-frontend:
	rm -rf .artifacts/web-static
	mkdir -p .artifacts/web-static
	bazelisk ${BAZEL_FLAGS} build ${BAZEL_COMMAND_FLAGS} -- //web:static
	cp -rf bazel-bin/web/static/. .artifacts/web-static


GO_BUILD_VERSION ?= $(shell git describe --tags --always --dirty --abbrev=7)
GO_BUILD_COMMIT ?= $(shell git rev-parse HEAD)
GO_BUILD_DATE ?= $(shell date -u '+%Y%m%d-%H%M%S')
GO_BUILD_LDFLAGS_PREFIX := -X github.com/pipe-cd/pipecd/pkg/version
GO_BUILD_OPTS ?= -ldflags "$(GO_BUILD_LDFLAGS_PREFIX).Version=$(GO_BUILD_VERSION) $(GO_BUILD_LDFLAGS_PREFIX).GitCommit=$(GO_BUILD_COMMIT) $(GO_BUILD_LDFLAGS_PREFIX).BuildDate=$(GO_BUILD_DATE) -w"
GO_BUILD_COMPONENT ?= pipecd
GO_BUILD_ENV ?= GOOS=linux GOARCH=amd64

.PHONY: build-go
build-go:
	$(GO_BUILD_ENV) CGO_ENABLED=0 go build $(GO_BUILD_OPTS) -o ./.artifacts/$(GO_BUILD_COMPONENT) ./cmd/$(GO_BUILD_COMPONENT)
