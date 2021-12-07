package main

import (
	"net"

	"github.com/Alation/alation_connector_manager/docker/remoteAgent/grpcinverter/ioc"

	"github.com/Alation/alation_connector_manager/docker/remoteAgent/grpcinverter"
	"google.golang.org/grpc"
)

func main() {
	forwardProxy := grpcinverter.NewForwardProxy()
	internalServer := grpc.NewServer(
		grpc.KeepaliveParams(grpcinverter.KEEPALIVE_SERVER_PARAMETERS),
		grpc.KeepaliveEnforcementPolicy(grpcinverter.KEEPALIVE_ENFORCEMENT_POLICY),
		grpc.ForceServerCodec(grpcinverter.NoopCodec{}),
		grpc.UnknownServiceHandler(forwardProxy.StreamHandler()))
	externalServer := grpc.NewServer(
		grpc.KeepaliveParams(grpcinverter.KEEPALIVE_SERVER_PARAMETERS),
		grpc.KeepaliveEnforcementPolicy(grpcinverter.KEEPALIVE_ENFORCEMENT_POLICY))
	ioc.RegisterGrpcInverterServer(externalServer, forwardProxy)
	internalListener, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		panic(err)
	}
	externalListener, err := net.Listen("tcp", "0.0.0.0:1234")
	if err != nil {
		panic(err)
	}
	go func() {
		internalServer.Serve(internalListener)
	}()
	externalServer.Serve(externalListener)
}
