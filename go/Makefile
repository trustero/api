include ../make/globals.mk
GO_PATH_BIN := $(shell go env GOPATH)/bin
.PHONY: foo
MODULE := "trustero-api-go"

deps:
	$(call outmsg, $(MODULE))
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest;
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest;

build: deps proto-gen

proto-gen:
	$(call outmsg, $(MODULE), "Generating golang proto bindings")
	@./proto-gen.sh

test:
	$(call outmsg, $(MODULE))
	@go test  ./... -failfast -short --coverprofile=coverage.txt

clean:
	$(call outmsg, $(MODULE))
	@go clean
ifneq ($(wildcard api/.*),)
	@find ./receptor_v1 -name "*.pb.*" -type f -delete
	@find ./receptor_v1 -type d -empty -delete
endif

fmt:
	$(call outmsg, $(MODULE))
	@go fmt -x ./... > /dev/null

run:
	$(call not-supported-msg,$(MODULE))

lint:
	@golangci-lint run