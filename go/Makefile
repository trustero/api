include ../make/globals.mk
GO_PATH_BIN := $(shell go env GOPATH)/bin
.PHONY: foo
MODULE := "trustero-api-go"

deps:
	$(call outmsg, $(MODULE))
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest;
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest;

protoc-build: deps
ifneq ($(shell git diff --quiet HEAD ${REF} -- ../proto || echo "changed"),)
	$(call outmsg, $(MODULE), "API Change detected")
else
	@$(MAKE) proto-gen
endif

build: protoc-build

proto-gen:
	$(call outmsg, $(MODULE), "Generating golang proto bindings")
	@./proto-gen.sh

test:
	$(call outmsg, $(MODULE))
	@go test  -failfast -short ./pkg/...

clean:
	$(call outmsg, $(MODULE))
	@go clean

fmt:
	$(call outmsg, $(MODULE))
	@go fmt -x ./...

run:
	$(call not-supported-msg,$(MODULE))

lint:
	$(call not-supported-msg,$(MODULE))