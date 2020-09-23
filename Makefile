BAZEL_FLAGS =
BAZEL_COMMAND_FLAGS =

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

.PHONY: build
build:
	bazelisk ${BAZEL_FLAGS} build ${BAZEL_COMMAND_FLAGS} -- //...

.PHONY: test
test:
	bazelisk ${BAZEL_FLAGS} test ${BAZEL_COMMAND_FLAGS} --config=unit -- //pkg/...

.PHONY: test-debug
test-debug:
	bazelisk ${BAZEL_FLAGS} test ${BAZEL_COMMAND_FLAGS} --test_output=all -- //pkg/...

.PHONY: test-integration
test-integration:
	bazelisk ${BAZEL_FLAGS} test ${BAZEL_COMMAND_FLAGS} --config=integration -- //pkg/...

.PHONY: coverage
coverage:
	bazelisk ${BAZEL_FLAGS} coverage ${BAZEL_COMMAND_FLAGS} //pkg/...

.PHONY: dep
dep:
	GO111MODULE=on go mod tidy
	GO111MODULE=on go mod vendor
	bazelisk run //:gazelle -- update-repos -from_file=go.mod -to_macro=repositories.bzl%go_repositories

.PHONY: gazelle
gazelle:
	bazelisk run //:gazelle

.PHONY: buildifier
buildifier:
	bazelisk run //:buildifier

.PHONY: clean
clean:
	bazelisk clean --expunge

.PHONY: expose-generated-go
expose-generated-go:
	./hack/expose-generated-go.sh pipe-cd pipe

.PHONY: site
site:
	hugo server --source=docs
