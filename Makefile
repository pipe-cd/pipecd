####################
# All make commands are following the format as "make action/target"
# "action" can be either:
#   build:   build artifacts such as binary, container image, chart
#   test:    execute test
#   run:     run a module locally
#   stop:    stop a locally running module
#   lint:    lint the source code
#   update:  update packages or dependencies to the newer versions
#   gen:     execute code or docs generation
#   release: commands used in release flow
#   push:    push artifacts such as helm chart
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
build/go: BUILD_OS ?= $(shell go version | cut -d ' ' -f4 | cut -d/ -f1)
build/go: BUILD_ARCH ?= $(shell go version | cut -d ' ' -f4 | cut -d/ -f2)
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

.PHONY: push
push/chart: BUCKET ?= charts.pipecd.dev
push/chart: VERSION ?= $(shell git describe --tags --always --dirty --abbrev=7)
push/chart: CREDENTIALS_FILE ?= ~/.config/gcloud/application_default_credentials.json
push/chart:
	@yq -i '.version = "${VERSION}" | .appVersion = "${VERSION}"' manifests/pipecd/Chart.yaml
	@yq -i '.version = "${VERSION}" | .appVersion = "${VERSION}"' manifests/piped/Chart.yaml
	@yq -i '.version = "${VERSION}" | .appVersion = "${VERSION}"' manifests/site/Chart.yaml
	@yq -i '.version = "${VERSION}" | .appVersion = "${VERSION}"' manifests/helloworld/Chart.yaml
	docker run --rm -it -v ${CREDENTIALS_FILE}:/secret -v ${PWD}:/repo gcr.io/pipecd/chart-releaser@sha256:fc432431b411a81d7658355c27ebaa924afe190962ab11d46f5a6cdff0833cc3 /chart-releaser --bucket=${BUCKET} --manifests-dir=repo/manifests --credentials-file=secret #v0.13.0
	@git checkout manifests/

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
	go test ./test/integration/...

# Run commands

.PHONY: run/pipecd
run/pipecd: $(eval DATE ?= $(shell date +%s))
run/pipecd: BUILD_VERSION ?= "$(shell git describe --tags --always --abbrev=7)-$(DATE)"
run/pipecd: BUILD_COMMIT ?= $(shell git rev-parse HEAD)
run/pipecd: BUILD_DATE ?= $(shell date -u '+%Y%m%d-%H%M%S')
run/pipecd: BUILD_LDFLAGS_PREFIX := -X github.com/pipe-cd/pipecd/pkg/version
run/pipecd: BUILD_OPTS ?= -ldflags "$(BUILD_LDFLAGS_PREFIX).version=$(BUILD_VERSION) $(BUILD_LDFLAGS_PREFIX).gitCommit=$(BUILD_COMMIT) $(BUILD_LDFLAGS_PREFIX).buildDate=$(BUILD_DATE) -w"
run/pipecd: CONTROL_PLANE_VALUES ?= ./quickstart/control-plane-values.yaml
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
	helm -n pipecd upgrade --install pipecd .artifacts/pipecd-$(BUILD_VERSION).tgz --create-namespace \
		--set server.image.repository=localhost:5001/pipecd \
		--set ops.image.repository=localhost:5001/pipecd \
		--values $(CONTROL_PLANE_VALUES)

.PHONY: stop/pipecd
stop/pipecd:
	helm -n pipecd uninstall pipecd

.PHONY: run/piped
run/piped: CONFIG_FILE ?=
run/piped: INSECURE ?= false
run/piped: LAUNCHER ?= false
run/piped:
ifeq ($(LAUNCHER),true)
	go run cmd/launcher/main.go launcher --config-file=$(CONFIG_FILE) --insecure=$(INSECURE)
else
	go run cmd/piped/main.go piped --tools-dir=/tmp/piped-bin --config-file=$(CONFIG_FILE) --insecure=$(INSECURE)
endif

.PHONY: run/web
run/web:
	yarn --cwd web dev

.PHONY: run/site
run/site:
	env RELEASE=$(shell head -n 1 RELEASE | cut -d ' ' -f 2) hugo server --source=docs

# Lint commands

.PHONY: lint/go
lint/go: FIX ?= false
lint/go: VERSION ?= sha256:78d1bbd01a9886a395dc8374218a6c0b7b67694e725dd76f0c8ac1de411b85e8 #v1.46.2
lint/go: FLAGS ?= --rm --platform linux/amd64 -e GOLANGCI_LINT_CACHE=/repo/.cache/golangci-lint -v ${PWD}:/repo -w /repo -it
lint/go:
ifeq ($(FIX),true)
	docker run ${FLAGS} golangci/golangci-lint@${VERSION} golangci-lint run --fix
else
	docker run ${FLAGS} golangci/golangci-lint@${VERSION} golangci-lint run
endif

.PHONY: lint/web
lint/web: FIX ?= false
lint/web:
ifeq ($(FIX),true)
	yarn --cwd web lint:fix
else
	yarn --cwd web lint
endif

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
	# NOTE: Specify a specific version temporally until the next release.
	docker run --rm -v ${PWD}:/repo -it --entrypoint ./tool/codegen/codegen.sh ghcr.io/pipe-cd/codegen@sha256:16766336bd7fd7d7e24eabf29aabe471bfda0631a5c51e0a8d1470a249139a32 /repo #v0.37.0-10-ge573dda

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

.PHONY: gen/contributions
gen/contributions:
	./hack/gen-contributions.sh

.PHONY: release
release: release/init release/docs

.PHONY: release/init
release/init:
	./hack/gen-release.sh $(version)

.PHONY: release/pick
release/pick:
	./hack/cherry-pick.sh $(branch) $(pull_numbers)

.PHONY: release/docs
release/docs:
	./hack/gen-release-docs.sh $(version)

# Other commands

.PHONY: kind-up
kind-up:
	./hack/create-kind-cluster.sh pipecd

.PHONY: kind-down
kind-down:
	kind delete cluster --name pipecd
