COVERAGE_DIR=coverage
COVER_PROFILE=$(COVERAGE_DIR)/cover.out
TMP_COVER_PROFILE=$(COVERAGE_DIR)/cover.tmp
PKGS=$$(go list ./... | grep -v /vendor)
VERSION?=latest
ALL_OS=\
	darwin \
	linux \
	windows
BIN_DIR := $(GOPATH)/bin
GOIMPORTS := $(BIN_DIR)/goimports
GOMETALINTER := $(BIN_DIR)/gometalinter
DEP := $(BIN_DIR)/dep
VENDOR := $(CURDIR)/vendor

.PHONY: all
all: $(SETUP) lint test

.PHONY: test
test: $(SETUP)
	@go test $(PKGS)

.PHONY: lint
lint: $(SETUP)
	@gometalinter ./... --enable=goimports --enable=gosimple \
	--enable=unparam --enable=unused --vendor -t

.PHONY: check
check: $(SETUP)
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
ci: $(SETUP) check test

.PHONY: install
install: $(SETUP)
	@go install $(PKGS)

.PHONY: build
build: $(SETUP)
	@go build $(PKGS)

.PHONY: buildrelease
buildrelease:
	@rm -rf release
	@mkdir release
	@for os in $(ALL_OS); do \
		GOOS=$$os GOARCH=amd64 go build -ldflags="-X main.version=$(VERSION)" -o release/yamldiff-v$(VERSION)-$$os-amd64 || exit 1; \
	done

$(GOIMPORTS):
	@go get -u golang.org/x/tools/cmd/goimports

$(GOMETALINTER):
	@go get -u github.com/alecthomas/gometalinter
	@gometalinter --install

$(DEP):
	@go get -u github.com/golang/dep/cmd/dep

$(TOOLS): $(GOIMPORTS) $(GOMETALINTER) $(DEP)

$(VENDOR): $(DEP)
	@dep ensure

$(SETUP): $(TOOLS) $(VENDOR)
