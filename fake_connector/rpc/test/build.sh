#!/usr/bin/env bash

ROOT=$(git rev-parse --show-toplevel)
python -m grpc_tools.protoc -I . --python_out="${ROOT}"/fake_connector/alation --grpc_python_out="${ROOT}"/fake_connector/alation test.proto
protoc -I . --go_out="${ROOT}"/fake_connector/connector --go-grpc_out="${ROOT}"/fake_connector/connector test.proto
protoc -I . --go_out=. --go-grpc_out=. test.proto