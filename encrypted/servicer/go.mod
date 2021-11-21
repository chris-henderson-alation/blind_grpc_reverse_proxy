module github.com/chris-henderson-alation/blind_grpc_reverse_proxy/servicer

go 1.16

require (
	github.com/Alation/alation_connector_manager/security/common v0.0.0
	github.com/chris-henderson-alation/blind_grpc_reverse_proxy/rpc v0.0.0
	github.com/chris-henderson-alation/blind_grpc_reverse_proxy/servicer_connector_proto v0.0.0
	github.com/golang/protobuf v1.5.0
	github.com/sirupsen/logrus v1.8.1
	golang.org/x/crypto v0.0.0-20211117183948-ae814b36b871
	google.golang.org/grpc v1.42.0
	google.golang.org/protobuf v1.27.1 // indirect
)

replace (
	github.com/Alation/alation_connector_manager/security/common => ../security/common
	github.com/chris-henderson-alation/blind_grpc_reverse_proxy/rpc => ../rpc
	github.com/chris-henderson-alation/blind_grpc_reverse_proxy/servicer_connector_proto => ../servicer_connector_proto
)
