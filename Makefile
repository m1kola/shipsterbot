IMPORT_PATH := github.com/m1kola/shipsterbot
OUTPUT_BIN  := shipsterbot


# Current directory based on the Makefile location withut a trailing slash
ROOT    := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))


GOPATH  := $(ROOT)/.gopath
BIN     := $(GOPATH)/bin
BASE    := $(GOPATH)/src/$(IMPORT_PATH)


$(BASE):
	@mkdir -p $(dir $@)
	@ln -sf $(ROOT) $@


# TODO: Replace with downloading a specific version
GODEP := $(BIN)/dep
$(GOPATH)/bin/dep: | $(BASE)
	go get github.com/golang/dep/cmd/dep


.PHONY: vendor
vendor: Gopkg.lock Gopkg.toml | $(GOPATH)/bin/dep
	cd $(BASE) && \
	${GODEP} ensure


# TODO: Download go-bindata
.PHONY: migrations
migrations: vendor
	cd $(BASE) && \
	cd ./migrations && go-bindata -pkg migrations .


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
