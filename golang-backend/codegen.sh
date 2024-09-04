#!/bin/sh
PROTO_DIR=./proto
GENERATED_DIR=./generated
mkdir -p ${GENERATED_DIR}
protoc --go_out=${GENERATED_DIR} --go_opt=paths=source_relative --go-grpc_out=${GENERATED_DIR} --go-grpc_opt=paths=source_relative ${PROTO_DIR}/*.proto