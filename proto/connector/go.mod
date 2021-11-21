module github.com/chris-henderson-alation/blind_grpc_reverse_proxy/connector

go 1.16

require (
	github.com/alecthomas/units v0.0.0-20210927113745-59d0afb8317a
	github.com/chris-henderson-alation/blind_grpc_reverse_proxy/servicer_connector_proto v0.0.0
	github.com/golang/protobuf v1.5.0
	github.com/sirupsen/logrus v1.8.1
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/grpc v1.42.0
	google.golang.org/protobuf v1.27.1 // indirect
)

replace github.com/chris-henderson-alation/blind_grpc_reverse_proxy/servicer_connector_proto => ../servicer_connector_proto
