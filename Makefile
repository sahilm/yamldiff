COVERAGE_DIR=coverage
COVER_PROFILE=$(COVERAGE_DIR)/cover.out
TMP_COVER_PROFILE=$(COVERAGE_DIR)/cover.tmp
PKGS=$$(go list ./... | grep -v /vendor)
VERSION?=latest
ALL_OS=\
	darwin \
	linux \
	windows

.PHONY: all
all: setup lint test

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

.PHONY: cover
cover:
	@rm -rf $(COVERAGE_DIR)
	@mkdir -p $(COVERAGE_DIR)
	@echo "mode: set" > $(COVER_PROFILE)
	@for pkg in $(PKGS); do \
		go test -v -coverprofile=$(TMP_COVER_PROFILE) $$pkg; \
		if [ -f $(TMP_COVER_PROFILE) ]; then \
			grep -v 'mode: set' $(TMP_COVER_PROFILE) >> $(COVER_PROFILE); \
			rm $(TMP_COVER_PROFILE); \
		fi; \
	done
	@go tool cover -html=$(COVER_PROFILE) -o $(COVERAGE_DIR)/index.html

.PHONY: ci
ci: setup check test

.PHONY: install
install: setup
	@go install $(PKGS)

.PHONY: build
build: setup
	@go build $(PKGS)

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
