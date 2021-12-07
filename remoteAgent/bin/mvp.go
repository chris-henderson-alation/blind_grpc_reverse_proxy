package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc/status"

	"github.com/Alation/alation_connector_manager/docker/remoteAgent/grpcinverter"

	"google.golang.org/grpc"
)

func main() {
	s := grpc.NewServer(grpc.UnknownServiceHandler(lol()), grpc.KeepaliveParams(grpcinverter.KEEPALIVE_SERVER_PARAMETERS), grpc.KeepaliveEnforcementPolicy(grpcinverter.KEEPALIVE_ENFORCEMENT_POLICY))
	listener, err := net.Listen("tcp", "192.168.40.1:1234")
	if err != nil {
		panic(err)
	}
	go s.Serve(listener)
	conn, err := grpc.Dial("192.168.40.1:1234", grpc.WithInsecure(), grpc.WithKeepaliveParams(grpcinverter.KEEPALIVE_CLIENT_PARAMETERS))
	if err != nil {
		panic(err)
	}
	_, err = grpc.NewClientStream(context.Background(), &grpc.StreamDesc{
		ServerStreams: true,
		ClientStreams: true,
	}, conn, "/lol.Lol/Lol", grpc.ForceCodec(grpcinverter.NoopCodec{}))
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Second * 999999)
}

func alation()      {}
func forwardProxy() {}
func reverseProxy() {}
func connector() {
	s := grpc.NewServer(
		grpc.KeepaliveParams(grpcinverter.KEEPALIVE_SERVER_PARAMETERS),
		grpc.KeepaliveEnforcementPolicy(grpcinverter.KEEPALIVE_ENFORCEMENT_POLICY),
		grpc.UnknownServiceHandler(func(srv interface{}, stream grpc.ServerStream) error {
			var counter uint64 = 0
			for {
				str := &String{Message: strconv.FormatUint(counter, 10)}
				err := stream.SendMsg(str)
				if err != nil {
					fmt.Println(err)
					return err
				}
				counter += 1
				time.Sleep(time.Second)
			}
			return nil
		}),
	)
	socket, err := net.Listen("tcp", "127.0.0.1:10001")
	if err != nil {
		panic(err)
	}
	s.Serve(socket)
}

func lol() grpc.StreamHandler {
	return func(_ interface{}, stream grpc.ServerStream) error {
		fmt.Println("connected")
		err := stream.RecvMsg(&[]byte{})
		switch err {
		case nil:
			fmt.Println("you're kidding")
		case io.EOF:
			fmt.Println("EOF")
		default:
			e, was := status.FromError(err)
			fmt.Println("was it a status? %v", was)
			fmt.Println(e)
		}
		return nil
	}
}
