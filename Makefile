.PHONY: all
all: setup lint test

PKGS := $(shell go list ./... | grep -v /vendor)
.PHONY: test
test: setup
	go test $(PKGS)

sources = $(shell find . -name '*.go' -not -path './vendor/*')
.PHONY: goimports
goimports: setup
	goimports -w $(sources)

.PHONY: lint
lint: setup
	$(BIN_DIR)/golangci-lint run

COVERAGE := $(CURDIR)/coverage
COVER_PROFILE :=$(COVERAGE)/cover.out
TMP_COVER_PROFILE :=$(COVERAGE)/cover.tmp
.PHONY: cover
cover: setup
	rm -rf $(COVERAGE)
	mkdir -p $(COVERAGE)
	echo "mode: set" > $(COVER_PROFILE)
	for pkg in $(PKGS); do \
		go test -v -coverprofile=$(TMP_COVER_PROFILE) $$pkg; \
		if [ -f $(TMP_COVER_PROFILE) ]; then \
			grep -v 'mode: set' $(TMP_COVER_PROFILE) >> $(COVER_PROFILE); \
			rm $(TMP_COVER_PROFILE); \
		fi; \
	done
	go tool cover -html=$(COVER_PROFILE) -o $(COVERAGE)/index.html

.PHONY: ci
ci: setup lint test

.PHONY: install
install: setup
	go install $(PKGS)

.PHONY: build
build: setup
	go build $(PKGS)

GOPATH ?= $(HOME)/go
BIN_DIR := $(GOPATH)/bin
GOIMPORTS := $(BIN_DIR)/goimports
GOLANG_CI_LINT := $(BIN_DIR)/golangci-lint

$(GOIMPORTS):
	go get -u golang.org/x/tools/cmd/goimports

$(GOLANG_CI_LINT):
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(BIN_DIR) v2.8.0

tools: $(GOIMPORTS) $(GOLANG_CI_LINT)

setup: tools

BINARY := yamldiff
VERSION ?= latest
PLATFORMS := darwin/amd64 linux/amd64 windows/amd64

temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))

.PHONY: $(PLATFORMS)
$(PLATFORMS): setup
	mkdir -p $(CURDIR)/release
	CGO_ENABLED=0 GOOS=$(os) GOARCH=$(arch) go build -ldflags="-X main.version=$(VERSION)" \
	-o release/$(BINARY)-v$(VERSION)-$(os)-$(arch)

.PHONY: release
release: $(PLATFORMS)
