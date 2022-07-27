#!/bin/bash

which protoc > /dev/null || { \
  echo "error: protoc not installed" >&2; \
  exit 1; \
}

[ -z ${DOC_DIR+x} ] && DOC_DIR="./md"
mkdir -p $DOC_DIR

[ -z ${PROTOS_SRC_DIR+x} ] && PROTOS_SRC_DIR="../proto"
PROTO_PACKAGES=$(ls $PROTOS_SRC_DIR)

export PATH=$PATH:$(go env GOPATH)/bin

for pkg in $PROTO_PACKAGES; do
  proto_files=$(find ${PROTOS_SRC_DIR}/"${pkg}" -name "*.proto")

  for proto in $proto_files; do
    protoc \
          --doc_out=$DOC_DIR \
          --doc_opt="markdown,${pkg}.md" \
          --proto_path="${PROTOS_SRC_DIR}" \
          -I "${PROTOS_SRC_DIR}" \
          "${proto}"
  done
done
