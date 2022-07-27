include ../make/globals.mk
GO_PATH_BIN := $(shell go env GOPATH)/bin
.PHONY: foo
MODULE := "api-gen-doc"
export DOC_DIR = ./md
export PROTOS_SRC_DIR = ../proto

deps:
	$(call outmsg, $(MODULE))
	@go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@latest;

build: deps proto-gen-doc
	$(call outmsg, $(MODULE), "Build done")

proto-gen-doc:
	$(call outmsg, $(MODULE), "Generate API reference documentation")
	@./proto-gen-doc.sh

test:
	$(call not-supported-msg, $(MODULE))

clean:
	$(call outmsg, $(MODULE))
	@rm -rf ${DOC_DIR}

fmt:
	$(call not-supported-msg, $(MODULE))

run:
	$(call not-supported-msg,$(MODULE))
