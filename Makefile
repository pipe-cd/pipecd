####################
# All make commands are following the format as "make action/target"
# "action" can be either:
#   check:   run checks which should be passed before committing
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
build/go: BUILD_VERSION ?= $(shell git describe --tags --always --dirty --abbrev=7 --match 'v[0-9]*.*')
build/go: BUILD_COMMIT ?= $(shell git rev-parse HEAD)
build/go: BUILD_DATE ?= $(shell date -u '+%Y%m%d-%H%M%S')
build/go: BUILD_LDFLAGS_PREFIX := -X github.com/pipe-cd/pipecd/pkg/version
build/go: BUILD_OPTS ?= -ldflags "$(BUILD_LDFLAGS_PREFIX).version=$(BUILD_VERSION) $(BUILD_LDFLAGS_PREFIX).gitCommit=$(BUILD_COMMIT) $(BUILD_LDFLAGS_PREFIX).buildDate=$(BUILD_DATE) -w" -trimpath
build/go: BUILD_OS ?= $(shell go version | cut -d ' ' -f4 | cut -d/ -f1)
build/go: BUILD_ARCH ?= $(shell go version | cut -d ' ' -f4 | cut -d/ -f2)
build/go: BUILD_ENV ?= GOOS=$(BUILD_OS) GOARCH=$(BUILD_ARCH) CGO_ENABLED=0
build/go: BIN_SUFFIX ?=
build/go:
ifndef MOD
	$(BUILD_ENV) go build $(BUILD_OPTS) -o ./.artifacts/control-plane$(BIN_SUFFIX) ./cmd/control-plane
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
build/chart: VERSION ?= $(shell git describe --tags --always --dirty --abbrev=7 --match 'v[0-9]*.*')
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

.PHONY: build/plugin
build/plugin: PLUGINS_BIN_DIR ?= ~/.piped/plugins
build/plugin: PLUGINS_SRC_DIR ?= ./pkg/app/pipedv1/plugin
build/plugin: PLUGINS_OUT_DIR ?= ${PWD}/.artifacts/plugins
build/plugin: PLUGINS ?= $(shell find $(PLUGINS_SRC_DIR) -mindepth 1 -maxdepth 1 -type d | while read -r dir; do basename "$$dir"; done | paste -sd, -) # comma separated list of plugins. eg: PLUGINS=kubernetes,ecs,lambda
build/plugin: BUILD_OPTS ?= -ldflags "-s -w" -trimpath
build/plugin: BUILD_OS ?= $(shell go version | cut -d ' ' -f4 | cut -d/ -f1)
build/plugin: BUILD_ARCH ?= $(shell go version | cut -d ' ' -f4 | cut -d/ -f2)
build/plugin: BUILD_ENV ?= GOOS=$(BUILD_OS) GOARCH=$(BUILD_ARCH) CGO_ENABLED=0
build/plugin: BIN_SUFFIX ?=
build/plugin:
	mkdir -p $(PLUGINS_BIN_DIR)
	@echo "Building plugins..."
	@for plugin in $(shell echo $(PLUGINS) | tr ',' ' '); do \
		echo "Building plugin: $$plugin"; \
		$(BUILD_ENV) go -C $(PLUGINS_SRC_DIR)/$$plugin build $(BUILD_OPTS) -o $(PLUGINS_OUT_DIR)/$${plugin}$(BIN_SUFFIX) . \
			&& cp $(PLUGINS_OUT_DIR)/$${plugin}$(BIN_SUFFIX) $(PLUGINS_BIN_DIR)/$$plugin; \
	done
	@echo "Plugins are built and copied to $(PLUGINS_BIN_DIR)"

.PHONY: push
push/chart: BUCKET ?= charts.pipecd.dev
push/chart: VERSION ?= $(shell git describe --tags --always --dirty --abbrev=7 --match 'v[0-9]*.*')
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
test/go: COVERAGE ?= false
test/go: COVERAGE_OPTS ?= -covermode=atomic
test/go: COVERAGE_OUTPUT ?= ${PWD}/coverage.out
test/go: setup-envtest
# Where to find the setup-envtest binary
test/go: GOBIN ?= ${PWD}/.dev/bin
# We need an absolute path for setup-envtest
test/go: ENVTEST_BIN ?= ${PWD}/.dev/bin
test/go: KUBEBUILDER_ASSETS ?= "$(shell $(GOBIN)/setup-envtest use --bin-dir $(ENVTEST_BIN) -p path)"
test/go: MODULES ?= $(shell find . -name go.mod | while read -r dir; do dirname "$$dir"; done | paste -sd, -) # comma separated list of modules. eg: MODULES=.,pkg/plugin/sdk,tool/actions-gh-release
test/go:
ifeq ($(COVERAGE), true)
	@echo "Run tests with coverage out of CI is not working expectedly. Because the coverage profile is overwritten by other tests."
	@echo "Testing go modules with coverage..."
	@for module in $(shell echo $(MODULES) | tr ',' ' '); do \
		if [ "$$module" = "." ]; then \
			echo "Testing root module"; \
			KUBEBUILDER_ASSETS=$(KUBEBUILDER_ASSETS) go test -failfast -race $(COVERAGE_OPTS) -coverprofile=$(COVERAGE_OUTPUT).tmp ./pkg/... ./cmd/...; \
		else \
			echo "Testing module: $$module"; \
			KUBEBUILDER_ASSETS=$(KUBEBUILDER_ASSETS) go -C $$module test -failfast -race $(COVERAGE_OPTS) -coverprofile=$(COVERAGE_OUTPUT).tmp ./...; \
		fi; \
	done
	cat $(COVERAGE_OUTPUT).tmp | grep -v ".pb.go\|.pb.validate.go" > $(COVERAGE_OUTPUT)
	rm -rf $(COVERAGE_OUTPUT).tmp
else
	@echo "Testing go modules..."
	@for module in $(shell echo $(MODULES) | tr ',' ' '); do \
		if [ "$$module" = "." ]; then \
			echo "Testing root module"; \
			KUBEBUILDER_ASSETS=$(KUBEBUILDER_ASSETS) go test -failfast -race ./pkg/... ./cmd/...; \
		else \
			echo "Testing module: $$module"; \
			KUBEBUILDER_ASSETS=$(KUBEBUILDER_ASSETS) go -C $$module test -failfast -race ./...; \
		fi; \
	done
endif

.PHONY: test/web
test/web:
	yarn --cwd web test:coverage --runInBand

.PHONY: test/integration
test/integration:
	go test ./test/integration/...

# Run commands

.PHONY: run/control-plane
run/control-plane: $(eval TIMESTAMP = $(shell date +%s))
# NOTE: previously `git describe --tags` was used to determine the version for running locally
# However, this does not work on a forked branch, so the decision was made to hardcode at version 0.0.0
# see: https://github.com/pipe-cd/pipecd/issues/4845
run/control-plane: BUILD_VERSION ?= "v0.0.0-$(shell git rev-parse --short HEAD)-$(TIMESTAMP)"
run/control-plane: BUILD_COMMIT ?= $(shell git rev-parse HEAD)
run/control-plane: BUILD_DATE ?= $(shell date -u '+%Y%m%d-%H%M%S')
run/control-plane: BUILD_LDFLAGS_PREFIX := -X github.com/pipe-cd/pipecd/pkg/version
run/control-plane: BUILD_OPTS ?= -ldflags "$(BUILD_LDFLAGS_PREFIX).version=$(BUILD_VERSION) $(BUILD_LDFLAGS_PREFIX).gitCommit=$(BUILD_COMMIT) $(BUILD_LDFLAGS_PREFIX).buildDate=$(BUILD_DATE) -w"
run/control-plane: CONTROL_PLANE_VALUES ?= ./quickstart/control-plane-values.yaml
run/control-plane:
	@echo "Building go binary of Control Plane..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(BUILD_ENV) go build $(BUILD_OPTS) -o ./.artifacts/control-plane ./cmd/control-plane

	@echo "Building web static files..."
	yarn --cwd web build

	@echo "Building docker image and pushing it to local registry..."
	docker build -f cmd/control-plane/Dockerfile -t localhost:5001/control-plane:$(BUILD_VERSION) .
	docker push localhost:5001/control-plane:$(BUILD_VERSION)

	@echo "Installing Control Plane in kind..."
	mkdir -p .artifacts
	helm package manifests/pipecd --version $(BUILD_VERSION) --app-version $(BUILD_VERSION) --dependency-update --destination .artifacts
	helm -n pipecd upgrade --install control-plane .artifacts/pipecd-$(BUILD_VERSION).tgz --create-namespace \
		--set server.image.repository=localhost:5001/control-plane \
		--set ops.image.repository=localhost:5001/control-plane \
		--values $(CONTROL_PLANE_VALUES)

.PHONY: stop/control-plane
stop/control-plane:
	helm -n pipecd uninstall control-plane

.PHONY: run/piped
run/piped: CONFIG_FILE ?=
run/piped: INSECURE ?= false
run/piped: LAUNCHER ?= false
run/piped: LOG_ENCODING ?= humanize
run/piped: EXPERIMENTAL ?= false
run/piped:
ifeq ($(EXPERIMENTAL), true)
	go run cmd/pipedv1/main.go piped --tools-dir=/tmp/piped-bin --config-file=$(CONFIG_FILE) --insecure=$(INSECURE) --log-encoding=$(LOG_ENCODING)
else ifeq ($(LAUNCHER),true)
	go run cmd/launcher/main.go launcher --config-file=$(CONFIG_FILE) --insecure=$(INSECURE) --log-encoding=$(LOG_ENCODING)
else
	go run cmd/piped/main.go piped --tools-dir=/tmp/piped-bin --config-file=$(CONFIG_FILE) --insecure=$(INSECURE) --log-encoding=$(LOG_ENCODING)
endif

.PHONY: run/web
run/web:
	yarn --cwd web dev

.PHONY: run/site
run/site:
	env RELEASE=$(shell grep '^tag:' RELEASE | awk '{print $$2}') hugo server --source=docs

# Lint commands

.PHONY: lint
lint: lint/go lint/web lint/helm

.PHONY: lint/go
lint/go: FIX ?= false
lint/go: VERSION ?= sha256:c2f5e6aaa7f89e7ab49f6bd45d8ce4ee5a030b132a5fbcac68b7959914a5a890 # golangci/golangci-lint:v1.64.7
lint/go: FLAGS ?= --rm -e GOCACHE=/repo/.cache/go-build -e GOLANGCI_LINT_CACHE=/repo/.cache/golangci-lint -v ${PWD}:/repo -it
lint/go: MODULES ?= $(shell find . -name go.mod | while read -r dir; do dirname "$$dir"; done | paste -sd, -) # comma separated list of modules. eg: MODULES=.,pkg/plugin/sdk
lint/go:
	@echo "Linting go modules..."
	@for module in $(shell echo $(MODULES) | tr ',' ' '); do \
		echo "Linting module: $$module"; \
		docker run ${FLAGS} -w /repo/$$module golangci/golangci-lint@${VERSION} golangci-lint run -v --config /repo/.golangci.yml --fix=$(FIX); \
	done

.PHONY: lint/web
lint/web: FIX ?= false
lint/web:
ifeq ($(FIX),true)
	yarn --cwd web lint:fix
	yarn --cwd web typecheck
else
	yarn --cwd web lint
	yarn --cwd web typecheck
endif

.PHONY: lint/helm
lint/helm:
	@for dir in $$(find ./manifests -mindepth 1 -maxdepth 1 -type d); do \
		helm lint $$dir || exit $$?; \
	done

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

.PHONY: update/copyright
update/copyright:
	./hack/update-copyright.sh

# Generate commands

.PHONY: gen/code
gen/code:
	# NOTE: Keep this container image as same as defined in .github/workflows/gen.yml
	docker run --rm -v ${PWD}:/repo -it --entrypoint ./tool/codegen/codegen.sh ghcr.io/pipe-cd/codegen@sha256:3aa25a5abafe40419861ce1f1667580d4274e144370d03ce9f1d00e9b391d7fd /repo # v0.52.0-135-gcefd641

.PHONY: gen/test-tls
gen/test-tls:
	openssl req -x509 -nodes -days 3650 -newkey rsa:2048 \
		-keyout pkg/rpc/testdata/tls.key \
		-out pkg/rpc/testdata/tls.crt \
		-subj "/CN=localhost" \
		-config pkg/rpc/testdata/tls.config

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

.PHONY: setup-envtest
# Where to install the setup-envtest binary
setup-envtest: export GOBIN ?= ${PWD}/.dev/bin
setup-envtest: ## Download setup-envtest locally if necessary.
	test -x $(GOBIN)/setup-envtest || go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest

# Check commands
.PHONY: check
check: build lint test check/gen/code check/dco

.PHONY: check/gen/code
check/gen: gen/code
	git add -N .
	git diff --exit-code --quiet HEAD

.PHONY: check/dco
check/dco:
	./hack/ensure-dco.sh

.PHONY: setup-local-oidc
setup-local-oidc:
	./hack/oidc/run-local-keycloak.sh

.PHONY: delete-local-oidc
delete-local-oidc:
	docker compose -f ./hack/oidc/docker-compose.yml down
