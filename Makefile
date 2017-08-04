EXTERNAL_TOOLS=\
	golang.org/x/tools/cmd/goimports \
	github.com/alecthomas/gometalinter \
	github.com/golang/dep/cmd/dep
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
all: lint test

.PHONY: test
test:
	@go test $(PKGS)

.PHONY: lint
lint:
	@gometalinter --enable=goimports --enable=gosimple \
	--enable=unparam --enable=unused --vendor

.PHONY: check
check:
	@gometalinter --disable-all --enable=vet --enable=vetshadow \
	--enable=errcheck --enable=goimports --vendor

.PHONY: bootstrap
bootstrap:
	@for tool in  $(EXTERNAL_TOOLS) ; do \
    	echo "Installing/Updating $$tool" ; \
    	go get -u $$tool; \
    done
	@dep ensure
	@gometalinter --install

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
ci: depensure check test

.PHONY: install
install:
	@go install $(PKGS)

.PHONY: build
build:
	@go build $(PKGS)

.PHONY: depensure
depensure:
	@dep ensure

.PHONY: depupdate
depupdate:
	@dep ensure -update

.PHONY: buildrelease
buildrelease:
	@rm -rf release
	@mkdir release
	@for os in $(ALL_OS); do \
		GOOS=$$os GOARCH=amd64 go build -ldflags="-X main.version=$(VERSION)" -o release/yamldiff-v$(VERSION)-$$os-amd64 || exit 1; \
	done
