module github.com/chris-henderson-alation/blind_grpc_reverse_proxy

go 1.16

require (
	github.com/alecthomas/units v0.0.0-20210927113745-59d0afb8317a
	github.com/chris-henderson-alation/blind_grpc_reverse_proxy/rpc v0.0.0
	github.com/golang/protobuf v1.5.2
	github.com/sirupsen/logrus v1.8.1
	google.golang.org/grpc v1.42.0
)

replace github.com/chris-henderson-alation/blind_grpc_reverse_proxy/rpc => ../rpc
