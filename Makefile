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
	gometalinter ./... --enable=goimports --enable=gosimple \
	--enable=unparam --enable=unused --disable=gotype --vendor -t

.PHONY: check
check: setup
	gometalinter ./... --disable-all --enable=vet --enable=vetshadow \
	--enable=errcheck --enable=goimports --vendor -t

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
ci: setup check test

.PHONY: install
install: setup
	go install $(PKGS)

.PHONY: build
build: setup
	go build $(PKGS)

BIN_DIR := $(GOPATH)/bin
GOIMPORTS := $(BIN_DIR)/goimports
GOMETALINTER := $(BIN_DIR)/gometalinter
DEP := $(BIN_DIR)/dep

$(GOIMPORTS):
	go get -u golang.org/x/tools/cmd/goimports

$(GOMETALINTER):
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install &> /dev/null

$(DEP):
	go get -u github.com/golang/dep/cmd/dep

tools: $(GOIMPORTS) $(GOMETALINTER) $(DEP)

vendor: $(DEP)
	dep ensure

setup: tools vendor

updatedeps:
	dep ensure -update

BINARY := yamldiff
VERSION ?= latest
PLATFORMS := darwin/amd64 linux/amd64 windows/amd64

temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))

.PHONY: $(PLATFORMS)
$(PLATFORMS): setup
	mkdir -p $(CURDIR)/release
	GOOS=$(os) GOARCH=$(arch) go build -ldflags="-X main.version=$(VERSION)" \
	-o release/$(BINARY)-v$(VERSION)-$(os)-$(arch)

.PHONY: release
release: $(PLATFORMS)
