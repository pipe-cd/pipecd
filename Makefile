# Build commands

.PHONY: build/backend
build/backend: BUILD_VERSION ?= $(shell git describe --tags --always --dirty --abbrev=7)
build/backend: BUILD_COMMIT ?= $(shell git rev-parse HEAD)
build/backend: BUILD_DATE ?= $(shell date -u '+%Y%m%d-%H%M%S')
build/backend: BUILD_LDFLAGS_PREFIX := -X github.com/pipe-cd/pipecd/pkg/version
build/backend: BUILD_OPTS ?= -ldflags "$(BUILD_LDFLAGS_PREFIX).version=$(BUILD_VERSION) $(BUILD_LDFLAGS_PREFIX).gitCommit=$(BUILD_COMMIT) $(BUILD_LDFLAGS_PREFIX).buildDate=$(BUILD_DATE) -w"
build/backend: BUILD_ARCH ?= GOOS=linux GOARCH=amd64
build/backend: BUILD_ENV ?= $(BUILD_ARCH) CGO_ENABLED=0
build/backend:
ifndef MOD
	$(BUILD_ENV) go build $(BUILD_OPTS) -o ./.artifacts/pipecd ./cmd/pipecd
	$(BUILD_ENV) go build $(BUILD_OPTS) -o ./.artifacts/piped ./cmd/piped
	$(BUILD_ENV) go build $(BUILD_OPTS) -o ./.artifacts/launcher ./cmd/launcher
	$(BUILD_ENV) go build $(BUILD_OPTS) -o ./.artifacts/pipectl ./cmd/pipectl
	$(BUILD_ENV) go build $(BUILD_OPTS) -o ./.artifacts/helloworld ./cmd/helloworld
else
	$(BUILD_ENV) go build $(BUILD_OPTS) -o ./.artifacts/$(MOD) ./cmd/$(MOD)
endif

.PHONY: build/frontend
build/frontend:
	yarn --cwd web build

.PHONY: build/chart
build/chart: VERSION ?= $(shell git describe --tags --always --dirty --abbrev=7)
build/chart:
	mkdir -p .artifacts
ifndef MOD
	helm package manifests/pipecd --version $(VERSION) --app-version $(VERSION) --dependency-update --destination .artifacts
	helm package manifests/piped --version $(VERSION) --app-version $(VERSION) --dependency-update --destination .artifacts
	helm package manifests/site --version $(VERSION) --app-version $(VERSION) --dependency-update --destination .artifacts
	helm package manifests/helloworld --version $(VERSION) --app-version $(VERSION) --dependency-update --destination .artifacts
else
	helm package manifests/$(MOD) --version $(VERSION) --app-version $(VERSION) --dependency-update --destination .artifacts
endif

.PHONY: build/image
build/image:
	@echo "Unimplemented"

# Test commands

.PHONY: test/backend
test/backend:
	go test ./pkg/... ./cmd/...

.PHONY: test/frontend
test/frontend:
	yarn --cwd web test:coverage --runInBand

.PHONY: test/integration
test/integration:
	@echo "Unimplemented"

# Run commands

.PHONY: run/pipecd
run/pipecd:
	@echo "Unimplemented"

# .PHONY: load-piped-image
# load-piped-image:
# 	bazelisk ${BAZEL_FLAGS} run ${BAZEL_COMMAND_FLAGS} --config=linux --config=stamping -- //cmd/piped:piped_app_image --norun
#
# .PHONY: kind-up
# kind-up:
# 	./hack/create-kind-cluster.sh pipecd
#
# .PHONY: kind-down
# kind-down:
# 	kind delete cluster --name pipecd

.PHONY: run/piped
run/piped:
	@echo "Unimplemented"

.PHONY: run/frontend
run/frontend:
	yarn --cwd web dev

.PHONY: run/site
run/site:
	env RELEASE=$(shell cut -c10- release/RELEASE) hugo server --source=docs

# Lint commands

.PHONY: lint/backend
lint/backend:
	@echo "Unimplemented"

.PHONY: lint/frontend
lint/frontend:
	@echo "Unimplemented"

# Update commands

.PHONY: update/backend-deps
update/backend-deps:
	go mod tidy
	go mod vendor

.PHONY: update/frontend-deps
update/frontend-deps:
	yarn --cwd web install --prefer-offline

.PHONY: update/docsy
update/docsy:
	rm -rf docs/themes/docsy
	git clone --recurse-submodules --depth 1 https://github.com/google/docsy.git docs/themes/docsy

# Generate commands

.PHONY: gen/code
gen/code:
	docker run --rm -v ${PWD}:/repo -it gcr.io/pipecd/codegen:0.7.0 /repo

.PHONY: gen/release
gen/release:
	./hack/gen-release.sh $(version)

.PHONY: gen/release-docs
gen/release-docs:
	./hack/gen/gen-release-docs.sh $(version)

.PHONY: gen/stable-docs
gen/stable-docs:
	./hack/gen-stable-docs.sh $(version)

.PHONY: gen/test-tls
gen/test-tls:
	openssl req -x509 -nodes -days 3650 -newkey rsa:2048 \
		-keyout pkg/rpc/testdata/tls.key \
		-out pkg/rpc/testdata/tls.crt \
		-subj "/CN=localhost" \
		-config pkg/rpc/testdata/tls.config
