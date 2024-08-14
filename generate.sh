#!/bin/bash
PROTO_DIR="./api"
OUT_DIR="./internal/service"

protoc --go_out=paths=source_relative:${OUT_DIR} \
       --go-grpc_out=paths=source_relative:${OUT_DIR} \
       --proto_path=${PROTO_DIR} \
       ${PROTO_DIR}/optionhub_service.proto

protoc --doc_out=. --doc_opt=markdown,README.md --proto_path=${PROTO_DIR} ${PROTO_DIR}/optionhub_service.proto