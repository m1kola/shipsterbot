IMPORT_PATH := github.com/m1kola/shipsterbot
OUTPUT_BIN  := shipsterbot


# Current directory based on the Makefile location withut a trailing slash
ROOT        := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
TARGET_OS   := $(shell uname -s|tr A-Z a-z)
TARGET_ARCH := amd64


GOPATH  := $(ROOT)/.gopath
BIN     := $(GOPATH)/bin
BASE    := $(GOPATH)/src/$(IMPORT_PATH)


$(GOPATH):
	@mkdir -p $@

$(BASE): | $(GOPATH)
	@mkdir -p $(dir $@)
	@ln -sf $(ROOT) $@

$(BIN): | $(BASE)
	@mkdir -p $@


GODEP_VERSION   := 0.3.2
GODEP_URL       := https://github.com/golang/dep/releases/download/v$(GODEP_VERSION)/dep-$(TARGET_OS)-$(TARGET_ARCH)
GODEP           := $(BIN)/dep
$(GOPATH)/bin/dep: | $(BIN)
	curl -fsSL -o $@ $(GODEP_URL) && \
	chmod +x $@


GOBINDATA := $(BIN)/go-bindata
$(BIN)/go-bindata: | $(BIN)
	go get -u github.com/jteeuwen/go-bindata/...


.PHONY: vendor
vendor: Gopkg.lock Gopkg.toml | $(GOPATH)/bin/dep
	cd $(BASE) && \
	${GODEP} ensure -vendor-only


.PHONY: migrations
migrations: vendor $(BIN)/go-bindata
	cd $(BASE) && \
	cd ./migrations && $(GOBINDATA) -pkg migrations .


.PHONY: build
build: vendor
	cd $(BASE) && \
	go build -o ${OUTPUT_BIN}


.PHONY: test
test: vendor
	cd $(BASE) && \
	go test ./...


.PHONY: clean
clean:
	rm -rf $(BASE)/${OUTPUT_BIN}
	rm -rf ${GOPATH}
