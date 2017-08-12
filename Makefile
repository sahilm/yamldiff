.PHONY: all
all: setup lint test

PKGS := $(shell go list ./... | grep -v /vendor)
.PHONY: test
test: setup
	@go test $(PKGS)

.PHONY: lint
lint: setup
	@gometalinter ./... --enable=goimports --enable=gosimple \
	--enable=unparam --enable=unused --vendor -t

.PHONY: check
check: setup
	@gometalinter ./... --disable-all --enable=vet --enable=vetshadow \
	--enable=errcheck --enable=goimports --vendor -t

COVERAGE := $(CURDIR)/coverage
COVER_PROFILE :=$(COVERAGE)/cover.out
TMP_COVER_PROFILE :=$(COVERAGE)/cover.tmp
$(COVERAGE): setup
	@rm -rf $(COVERAGE)
	@mkdir -p $(COVERAGE)
	@echo "mode: set" > $(COVER_PROFILE)
	@for pkg in $(PKGS); do \
		go test -v -coverprofile=$(TMP_COVER_PROFILE) $$pkg; \
		if [ -f $(TMP_COVER_PROFILE) ]; then \
			grep -v 'mode: set' $(TMP_COVER_PROFILE) >> $(COVER_PROFILE); \
			rm $(TMP_COVER_PROFILE); \
		fi; \
	done
	@go tool cover -html=$(COVER_PROFILE) -o $(COVERAGE)/index.html

cover: $(COVERAGE)

.PHONY: ci
ci: setup check test

.PHONY: install
install: setup
	@go install $(PKGS)

.PHONY: build
build: setup
	@go build $(PKGS)

VERSION?=latest
ALL_OS=\
	darwin \
	linux \
	windows
.PHONY: buildrelease
buildrelease:
	@rm -rf release
	@mkdir release
	@for os in $(ALL_OS); do \
		GOOS=$$os GOARCH=amd64 go build -ldflags="-X main.version=$(VERSION)" -o release/yamldiff-v$(VERSION)-$$os-amd64 || exit 1; \
	done

BIN_DIR := $(GOPATH)/bin
GOIMPORTS := $(BIN_DIR)/goimports
GOMETALINTER := $(BIN_DIR)/gometalinter
DEP := $(BIN_DIR)/dep
VENDOR := $(CURDIR)/vendor

$(GOIMPORTS):
	@go get -u golang.org/x/tools/cmd/goimports

$(GOMETALINTER):
	@go get -u github.com/alecthomas/gometalinter
	@gometalinter --install

$(DEP):
	@go get -u github.com/golang/dep/cmd/dep

tools: $(GOIMPORTS) $(GOMETALINTER) $(DEP)

$(VENDOR): $(DEP)
	@dep ensure

setup: tools $(VENDOR)
