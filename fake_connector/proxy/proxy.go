package main // import "fake_connector/proxy"
import (
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func main() {
	external := NewInternetFacing("1234")
	internal := grpc.NewServer(
		grpc.UnknownServiceHandler(external.StreamHandler()),
		grpc.ForceServerCodec(NoopCodec{}),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time: time.Second * 30,
		}), grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             time.Second * 10,
			PermitWithoutStream: true,
		}))
	listener, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		panic(err)
	}
	fmt.Printf("starting internal facing server on :8080\n")
	internal.Serve(listener)
}
