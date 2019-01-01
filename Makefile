IMPORT_PATH := github.com/m1kola/shipsterbot
OUTPUT_BIN  := shipsterbot


# We support only Darwin and Linux and amd64, at the moment,
# but we do not check OS and arch here for simplicity
# Also no need in compilation for different target archs,
# at the moment: so we compile just for the current OS and arch
GOOS   := $(shell uname -s|tr A-Z a-z)
GOARCH := amd64


# Current directory based on the Makefile location withut a trailing slash
ROOT := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))

GOPATH  := $(ROOT)/.gopath
BIN     := $(GOPATH)/bin
BASE    := $(GOPATH)/src/$(IMPORT_PATH)


# Create go path dir
$(GOPATH):
	@mkdir -p $@

# Create base dir for the project
$(BASE): | $(GOPATH)
	@mkdir -p $(dir $@)
	@ln -sf $(ROOT) $@

# Create a dir for binaries
# We need it to be able to install dep and other 3rd party binaries
$(BIN): | $(BASE)
	@mkdir -p $@


# Download a specific version of dep
GODEP_VERSION   := 0.5.0
GODEP_URL       := https://github.com/golang/dep/releases/download/v$(GODEP_VERSION)/dep-$(GOOS)-$(GOARCH)
GODEP           := $(BIN)/dep
$(GODEP): | $(BIN)
	curl -fsSL -o $@ $(GODEP_URL) && \
	chmod +x $@


# Install mockgen
# NOTE that we don't call mockgen it directly: we use it via go generate
GOMOCKGEN := $(BIN)/mockgen
$(GOMOCKGEN): | $(BIN)
	go get -u github.com/golang/mock/mockgen


# Install dependencies using dep
.PHONY: vendor
vendor: Gopkg.lock Gopkg.toml | $(GODEP)
	cd $(BASE) && \
	$(GODEP) ensure -vendor-only


# Run go generate (generates mocks, etc)
.PHONY: go_generate
go_generate: vendor $(GOMOCKGEN)
	cd $(BASE) && \
	go generate ./...


# Build the application
.PHONY: build
build: vendor migrations
	cd $(BASE) && \
	go build -o $(OUTPUT_BIN) ./cmd/shipsterbot


# Run tests for all pages
.PHONY: test
test: vendor migrations
	cd $(BASE) && \
	go test -race -cover ./...


# Cleanup working directory
.PHONY: clean
clean:
	rm -rf $(BASE)/$(OUTPUT_BIN)
	rm -rf $(GOPATH)
