package rpc // import "github.com/chris-henderson-alation/blind_grpc_reverse_proxy/rpc"

//go:generate protoc -I . --go_out=plugins=grpc:. rpc.proto
