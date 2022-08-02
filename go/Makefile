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
	@printf "\033[1;33mFIXME - enable examples folder in test. Skipping while it's broken fixed\033[0m\n"
	@go test $(shell go list ./... | grep -v /examples/ ) -coverprofile coverage.txt
	@grep -Ev '\.pb\.' coverage.txt > coverage.txt.tmp && mv coverage.txt.tmp coverage.txt

clean:
	$(call outmsg, $(MODULE))
	@go clean

fmt:
	$(call outmsg, $(MODULE))
	@go fmt -x ./... >/dev/null

run:
	$(call not-supported-msg,$(MODULE))

lint:
ifneq ($(shell gofmt -l . | wc -l | grep -v "0"),)
	@printf "\033[0;31m[make $@] gofmt errors found\033[0m\n" >&2
	@printf "\033[1;33mThe following files are not formatted correctly\n" >&2
	@gofmt -l . >&2
	@printf "To fix, run\033[0m\n" >&2
	@printf "\033[1m$$> make fmt\033[0m\n" >&2
	@exit 1
endif
