#!/bin/bash

which protoc > /dev/null || { \
  echo "error: protoc not installed" >&2; \
  exit 1; \
}

[ -z ${DOC_DIR+x} ] && DOC_DIR="."
mkdir -p $DOC_DIR

[ -z ${PROTOS_SRC_DIR+x} ] && PROTOS_SRC_DIR="../proto"
PROTO_PACKAGES=$(ls $PROTOS_SRC_DIR)

export PATH=$PATH:$(go env GOPATH)/bin

for pkg in $PROTO_PACKAGES; do
  proto_files=$(find ${PROTOS_SRC_DIR}/"${pkg}" -name "*.proto")
  doc_dir=${DOC_DIR}/${pkg}
  mkdir -p ${doc_dir}
  for proto in $proto_files; do
    proto_prefix=${proto%".proto"}
    protoc \
          --doc_out=${doc_dir} \
          --doc_opt="markdown,${proto_prefix}.md" \
          --proto_path="${PROTOS_SRC_DIR}" \
          -I "${PROTOS_SRC_DIR}" \
          "${proto}"
  done
done
