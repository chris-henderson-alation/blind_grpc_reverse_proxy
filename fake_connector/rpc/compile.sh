#!/usr/bin/env bash

GO_OUT=$(git rev-parse --show-toplevel)/fake_connector/rpc/go

(
cd ocf/bi/v1/
protoc -I . --go-grpc_out="${GO_OUT}" --go_out="${GO_OUT}" ./*.proto
)

(
cd ocf/bi/v2/
protoc -I . --go-grpc_out="${GO_OUT}" --go_out="${GO_OUT}" ./*.proto
)

(
cd ocf/rdbms
protoc -I . --go-grpc_out="${GO_OUT}" --go_out="${GO_OUT}" ./*.proto
)
