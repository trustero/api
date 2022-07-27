#!/bin/bash

which protoc > /dev/null || {
  echo "error: protoc not installed" >&2
  exit 1
}

GO_PATH_BIN=$(go env GOPATH)/bin
PROTOC_FLAGS="--experimental_allow_proto3_optional"
PROTO_MODULES="receptor_v1"
PROTOS_SRC_DIR="../proto/"

export GO111MODULE=on
export PATH=$PATH:$GO_PATH_BIN

outDir="."
mkdir -p $outDir

for module in $PROTO_MODULES; do
  proto_files=$(find ../proto/"${module}" -name "*.proto")

  for proto in $proto_files; do
    protoc $PROTOC_FLAGS \
          --go_out=paths=source_relative:$outDir \
          --go-grpc_out=require_unimplemented_servers=false,paths=source_relative:$outDir \
          --proto_path="${PROTOS_SRC_DIR}" \
          -I "${PROTOS_SRC_DIR}" \
          "${proto}"
  done;
  
done;
