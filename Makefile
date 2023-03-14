# Image URL to use all building/pushing image targets
TAG ?= latest
IMG ?= ghcr.io/undistro/marvin:${TAG}

# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: all
all: build

##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: test
test: fmt vet ## Run tests.
	go test ./... -coverprofile cover.out

.PHONY: addlicense
addlicense: ## Add copyright license headers in source code files
	@test -s $(LOCALBIN)/addlicense || GOBIN=$(LOCALBIN) go install github.com/google/addlicense@latest
	$(LOCALBIN)/addlicense -c "Undistro Authors" -l "apache" -ignore ".github/**" -ignore ".idea/**" -ignore "dist/**" -ignore ".goreleaser.yaml" .

.PHONY: checklicense
checklicense: ## Check copyright license headers in source code files
	@test -s $(LOCALBIN)/addlicense || GOBIN=$(LOCALBIN) go install github.com/google/addlicense@latest
	$(LOCALBIN)/addlicense -c "Undistro Authors" -l "apache" -ignore ".github/**" -ignore ".idea/**" -ignore "dist/**" -ignore ".goreleaser.yaml" -check .

##@ Build

.PHONY: build
build: fmt vet ## Build marvin binary.
	go build -ldflags="-s -w -X github.com/undistro/marvin/pkg/version.version=${TAG}" -o bin/marvin main.go

.PHONY: run
run: fmt vet ## Run marvin from your host.
	go run ./main.go

PLATFORMS ?= linux/arm64,linux/amd64,linux/s390x,linux/ppc64le
.PHONY: docker-buildx
docker-buildx: test ## Build and push docker image for cross-platform support.
	sed -e '1 s/\(^FROM\)/FROM --platform=\$$\{BUILDPLATFORM\}/; t' -e ' 1,// s//FROM --platform=\$$\{BUILDPLATFORM\}/' Dockerfile > Dockerfile.cross
	- docker buildx create --name cross-builder
	docker buildx use cross-builder
	- docker buildx build --push --platform=$(PLATFORMS) --tag ${IMG} -f Dockerfile.cross .
	- docker buildx rm cross-builder
	rm Dockerfile.cross

.PHONY: docker-build
docker-build: test ## Build docker image.
	docker build -t ${IMG} .

.PHONY: docker-push
docker-push: ## Push docker image.
	docker push ${IMG}

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)
