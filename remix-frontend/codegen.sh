#!/bin/sh

mkdir -p app/grpc/generated

GENERATED_CODE_DIR=./app/grpc/generated
PROTO_DIR=../golang-backend/proto

protoc \
--plugin=protoc-gen-ts_proto=./node_modules/.bin/protoc-gen-ts_proto \
--ts_proto_out=$GENERATED_CODE_DIR \
--ts_proto_opt=outputServices=nice-grpc,outputServices=generic-definitions,useExactTypes=false,esModuleInterop=true,importSuffix=.js \
--proto_path=$PROTO_DIR \
$PROTO_DIR/*.proto
