BIN_NAME          ?= "ghh-api"
BIN_ARCH          ?= "amd64"
BIN_VERSION 	   = $(shell bin/ghh-api version)
FILES_TO_FMT      ?= $(shell find . -path ./vendor -prune -o -name '*.go' -print)

DOCKER_IMAGE_NAME ?= ${BIN_NAME}
DOCKER_IMAGE_TAG  ?= $(subst /,-,$(shell git rev-parse --abbrev-ref HEAD))-$(shell date +%Y-%m-%d)-$(shell git rev-parse --short HEAD)

TMP_GOPATH        ?= /tmp/ghh-api-go
GOBIN             ?= ${GOPATH}/bin
GO                := $(shell which go)
GO111MODULE       ?= on
export GO111MODULE
# GOPROXY           ?= https://proxy.golang.org
# export GOPROXY

# Tools.
GOIMPORTS         ?= $(GOBIN)/goimports-$(GOIMPORTS_VERSION)
GOIMPORTS_VERSION ?= 9d4d845e86f14303813298ede731a971dd65b593
CI_LINT           ?= $(GOBIN)/golangci-lint-$(CI_LINT_VERSION)
CI_LINT_VERSION   ?= v1.17.1
GIT               ?= $(shell which git)
ME                ?= $(shell whoami)

AUDIT_ARGS        ?= --logLevel=debug --logFormat=json

# fetch_go_bin_version downloads (go gets) the binary from specific version and installs it in $(GOBIN)/<bin>-<version>
# arguments:
# $(1): Install path. (e.g github.com/campoy/embedmd)
# $(2): Tag or revision for checkout.
define fetch_go_bin_version
  @mkdir -p $(GOBIN)
  @mkdir -p $(TMP_GOPATH)

  @echo ">> fetching $(1)@$(2) revision/version"
  @if [ ! -d '$(TMP_GOPATH)/src/$(1)' ]; then \
    GOPATH='$(TMP_GOPATH)' GO111MODULE='off' go get -d -u '$(1)/...'; \
  else \
    CDPATH='' cd -- '$(TMP_GOPATH)/src/$(1)' && git fetch; \
  fi
  @CDPATH='' cd -- '$(TMP_GOPATH)/src/$(1)' && git checkout -f -q '$(2)'
  @echo ">> installing $(1)@$(2)"
  @GOBIN='$(TMP_GOPATH)/bin' GOPATH='$(TMP_GOPATH)' GO111MODULE='off' go install '$(1)'
  @mv -- '$(TMP_GOPATH)/bin/$(shell basename $(1))' '$(GOBIN)/$(shell basename $(1))-$(2)'
  @echo ">> produced $(GOBIN)/$(shell basename $(1))-$(2)"

endef

.PHONY: all
all: lint build docker package-deb

.PHONY: build
build: check-git go-mod-tidy 
	@echo ">> building binaries bin/$(BIN_NAME)"
	@$(GO) build -o bin/$(BIN_NAME) -v

.PHONY: run
run: build
	@echo ">> running binary bin/$(BIN_NAME) ${KIRBY_ARGS}"
	@bin/$(BIN_NAME) ${KIRBY_ARGS}

.PHONY: docker
docker:
	@echo ">> building Docker image"
	@docker build . -t $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)
	@docker tag $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) $(DOCKER_IMAGE_NAME):latest

.PHONY: lint
lint: $(CI_LINT)
	@echo ">> linting code"
	@$(CI_LINT) run -v --enable-all --disable=misspell --disable=gochecknoglobals --disable=gochecknoinits \
		--disable=lll --disable=interfacer --disable=dupl --no-config --exclude-use-default --exclude ".*Log"

.PHONY: go-mod-tidy
go-mod-tidy: check-git
	@go mod tidy

.PHONY: check-go-mod
check-go-mod:
	@go mod verify

.PHONY: check-git
check-git:
ifneq ($(GIT),)
	@test -x $(GIT) || (echo >&2 "No git executable binary found at $(GIT)."; exit 1)
else
	@echo >&2 "No git binary found."; exit 1
endif

# non-phony targets
$(GOIMPORTS):
	$(call fetch_go_bin_version,golang.org/x/tools/cmd/goimports,$(GOIMPORTS_VERSION))

$(CI_LINT):
	$(call fetch_go_bin_version,github.com/golangci/golangci-lint/cmd/golangci-lint,$(CI_LINT_VERSION))
