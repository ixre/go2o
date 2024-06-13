#!/usr/bin/env sh

bun add --dev ts-proto && rm -rf output

GO2O_HOME=.
PROTO_PATH=$(find "${GO2O_HOME}" -name "idl" -print -quit)
mkdir output && protoc -I "$PROTO_PATH" --plugin=./node_modules/.bin/protoc-gen-ts_proto \
  --ts_proto_out=./output "$PROTO_PATH"/*.proto "$PROTO_PATH"/**/*.proto
