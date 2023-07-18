NAME := kelpie
RELEASE_DIR := bin
BUILD_TARGETS := build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64
GOVERSION = $(shell go version)
THIS_GOOS = $(word 1,$(subst /, ,$(lastword $(GOVERSION))))
THIS_GOARCH = $(word 2,$(subst /, ,$(lastword $(GOVERSION))))
GOOS = $(THIS_GOOS)
GOARCH = $(THIS_GOARCH)
VERSION = $(shell cat ./VERSION)
REVISION = $(shell git rev-parse --verify HEAD)

.PHONY: fmt lint test all build clean rebuild mock

fmt: ## format
	@go fmt

lint: ## Examine source code and lint
	@go vet ./...
	@golint -set_exit_status ./...

test: ## run test
	@go test -v -cover -test.v -count 1 ./...

all: $(BUILD_TARGETS) ## build for all platform

build: $(RELEASE_DIR)/$(NAME)_$(GOOS)_$(GOARCH) ## build kelpie

build-linux-amd64: ## build AMD64 linux binary
	@$(MAKE) build GOOS=linux GOARCH=amd64

build-linux-arm64: ## build ARM64 linux binary
	@$(MAKE) build GOOS=linux GOARCH=arm64

build-darwin-amd64: ## build AMD64 darwin binary
	@$(MAKE) build GOOS=darwin GOARCH=amd64

build-darwin-arm64: ## build AMD64 darwin binary
	@$(MAKE) build GOOS=darwin GOARCH=arm64

$(RELEASE_DIR)/$(NAME)_$(GOOS)_$(GOARCH):
	@printf "\e[32m"
	@echo "==> Build kelpie for ${GOOS}-${GOARCH}"
	@printf "\e[90m"
	@GO111MODULE=on go build -tags netgo -a -v -o $(RELEASE_DIR)/$(NAME)_$(GOOS)_$(GOARCH) \
		-ldflags "-X main.version=$(VERSION) -X main.revision=$(REVISION)" \
		./main.go
	@printf "\e[m"

clean: ## Clean up built files
	@printf "\e[32m"
	@echo '==> clean up built files ./${RELEASE_DIR}/...'
	@printf "\e[90m"
	@ls -1 ./${RELEASE_DIR}
	@rm -rf ${RELEASE_DIR}/*
	@printf "\e[m"

rebuild: clean build

mock: vsphere/client.go
	@mockgen -source vsphere/client.go -destination vsphere/mock/client.go -package mock