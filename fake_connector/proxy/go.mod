module fake_connector/proxy

go 1.16

require (
	alation.com/ocf v0.0.0
	github.com/chris-henderson-alation/blind_grpc_reverse_proxy/rpc v0.0.0
	github.com/golang/protobuf v1.5.2
	google.golang.org/grpc v1.42.0
)

replace (
	alation.com/ocf => ../rpc/go/alation.com/ocf
	github.com/chris-henderson-alation/blind_grpc_reverse_proxy/rpc => ./../rpc/rpc
)
