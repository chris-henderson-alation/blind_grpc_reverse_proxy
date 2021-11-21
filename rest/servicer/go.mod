module github.com/chris-henderson-alation/blind_grpc_reverse_proxy/servicer

go 1.16

require (
	github.com/gin-gonic/gin v1.7.4
	github.com/golang/protobuf v1.5.0
	google.golang.org/grpc v1.42.0
	google.golang.org/protobuf v1.27.1
	github.com/chris-henderson-alation/blind_grpc_reverse_proxy/rpc v0.0.0
)

replace (
	"github.com/chris-henderson-alation/blind_grpc_reverse_proxy/rpc" => ../rpc
)
