#!/bin/bash

which protoc > /dev/null || {
  echo "error: protoc not installed" >&2
  exit 1
}

set -euo pipefail

cd $(dirname "$0")

codegen_dir="./lib"
index_js="./index.js"
> $index_js

rm -rf ${codegen_dir}
mkdir ${codegen_dir}

proto_location="$(
  cd ..
  pwd
)/proto/receptor_v1"

proto_files=$(find "${proto_location}" -name '*.proto')

for proto_file in $proto_files; do
  protoc "${proto_file}" \
    -I="${proto_location}" \
    --js_out=import_style=commonjs,binary:"${codegen_dir}" \
    --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:"${codegen_dir}" \
    --experimental_allow_proto3_optional

  protoc "${proto_file}" \
    -I="${proto_location}" \
    --js_out=import_style=typescript,binary:"${codegen_dir}" \
    --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:"${codegen_dir}" \
    --experimental_allow_proto3_optional
done

js_files=$(find "${codegen_dir}" -name '*.js' | tr '\n' ' ')

for js_file in $js_files; do
  echo "export * from \"${js_file}\"" >> ${index_js}
done
