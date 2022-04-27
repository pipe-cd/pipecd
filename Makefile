####################
# All make commands are following the format as "make action/target"
# "action" can be either:
#   build:  build artifacts such as binary, container image, chart
#   test:   execute test
#   run:    run a module locally
#   lint:   lint the source code
#   update: update packages or dependencies to the newer versions
#   gen:    execute code or docs generation
####################

# Build commands

.PHONY: build
build: build/go build/web

.PHONY: build/go
build/go: BUILD_VERSION ?= $(shell git describe --tags --always --dirty --abbrev=7)
build/go: BUILD_COMMIT ?= $(shell git rev-parse HEAD)
build/go: BUILD_DATE ?= $(shell date -u '+%Y%m%d-%H%M%S')
build/go: BUILD_LDFLAGS_PREFIX := -X github.com/pipe-cd/pipecd/pkg/version
build/go: BUILD_OPTS ?= -ldflags "$(BUILD_LDFLAGS_PREFIX).version=$(BUILD_VERSION) $(BUILD_LDFLAGS_PREFIX).gitCommit=$(BUILD_COMMIT) $(BUILD_LDFLAGS_PREFIX).buildDate=$(BUILD_DATE) -w"
build/go: BUILD_OS ?= linux
build/go: BUILD_ARCH ?= amd64
build/go: BUILD_ENV ?= GOOS=$(BUILD_OS) GOARCH=$(BUILD_ARCH) CGO_ENABLED=0
build/go: BIN_SUFFIX ?=
build/go:
ifndef MOD
	$(BUILD_ENV) go build $(BUILD_OPTS) -o ./.artifacts/pipecd$(BIN_SUFFIX) ./cmd/pipecd
	$(BUILD_ENV) go build $(BUILD_OPTS) -o ./.artifacts/piped$(BIN_SUFFIX) ./cmd/piped
	$(BUILD_ENV) go build $(BUILD_OPTS) -o ./.artifacts/launcher$(BIN_SUFFIX) ./cmd/launcher
	$(BUILD_ENV) go build $(BUILD_OPTS) -o ./.artifacts/pipectl$(BIN_SUFFIX) ./cmd/pipectl
	$(BUILD_ENV) go build $(BUILD_OPTS) -o ./.artifacts/helloworld$(BIN_SUFFIX) ./cmd/helloworld
else
	$(BUILD_ENV) go build $(BUILD_OPTS) -o ./.artifacts/$(MOD)$(BIN_SUFFIX) ./cmd/$(MOD)
endif

.PHONY: build/web
build/web:
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

# Test commands

.PHONY: test
test: test/go test/web

.PHONY: test/go
test/go:
	go test ./pkg/... ./cmd/...

.PHONY: test/web
test/web:
	yarn --cwd web test:coverage --runInBand

.PHONY: test/integration
test/integration:
	@echo "Unimplemented"

# Run commands

.PHONY: run/pipecd
run/pipecd: BUILD_VERSION ?= $(shell git describe --tags --always --dirty --abbrev=7)
run/pipecd: BUILD_COMMIT ?= $(shell git rev-parse HEAD)
run/pipecd: BUILD_DATE ?= $(shell date -u '+%Y%m%d-%H%M%S')
run/pipecd: BUILD_LDFLAGS_PREFIX := -X github.com/pipe-cd/pipecd/pkg/version
run/pipecd: BUILD_OPTS ?= -ldflags "$(BUILD_LDFLAGS_PREFIX).version=$(BUILD_VERSION) $(BUILD_LDFLAGS_PREFIX).gitCommit=$(BUILD_COMMIT) $(BUILD_LDFLAGS_PREFIX).buildDate=$(BUILD_DATE) -w"
run/pipecd:
	@echo "Building go binary of Control Plane..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(BUILD_ENV) go build $(BUILD_OPTS) -o ./.artifacts/pipecd ./cmd/pipecd

	@echo "Building web static files..."
	yarn --cwd web build

	@echo "Building docker image and pushing it to local registry..."
	docker build -f cmd/pipecd/Dockerfile -t localhost:5001/pipecd:$(BUILD_VERSION) .
	docker push localhost:5001/pipecd:$(BUILD_VERSION)

	@echo "Installing Control Plane in kind..."
	mkdir -p .artifacts
	helm package manifests/pipecd --version $(BUILD_VERSION) --app-version $(BUILD_VERSION) --dependency-update --destination .artifacts
	helm -n pipecd install pipecd .artifacts/pipecd-$(BUILD_VERSION).tgz --create-namespace \
		--values ./quickstart/control-plane-values.yaml \
		--set server.image.repository=localhost:5001/pipecd \
		--set ops.image.repository=localhost:5001/pipecd

.PHONY: run/piped
run/piped: CONFIG_FILE ?=
run/piped:
	go run cmd/piped/main.go piped --tools-dir=/tmp/piped-bin --config-file=$(CONFIG_FILE)

.PHONY: run/web
run/web:
	yarn --cwd web dev

.PHONY: run/site
run/site:
	env RELEASE=$(shell cut -c10- release/RELEASE) hugo server --source=docs

# Lint commands

.PHONY: lint/go
lint/go:
	@echo "Unimplemented"

.PHONY: lint/web
lint/web:
	@echo "Unimplemented"

# Update commands

.PHONY: update/go-deps
update/go-deps:
	go mod tidy
	go mod vendor

.PHONY: update/web-deps
update/web-deps:
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

# Other commands

.PHONY: kind-up
kind-up:
	./hack/create-kind-cluster.sh pipecd

.PHONY: kind-down
kind-down:
	kind delete cluster --name pipecd
