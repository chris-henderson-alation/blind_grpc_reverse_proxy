package main // import "testclient"
import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
)

var port = "8080"

func main() {
	perf()
	//clientStream()
}

func clientStream() {
	conn, err := grpc.Dial(fmt.Sprintf("0.0.0.0:%s", port),
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                time.Second * 30,
			PermitWithoutStream: true,
		}))
	if err != nil {
		panic(err)
	}
	client := NewTestClient(conn)
	md := metadata.New(map[string]string{"x-alation-connector": "1"})
	stream, err := client.ClientStream(metadata.NewOutgoingContext(context.Background(), md))
	if err != nil {
		panic(err)
	}
	if err := stream.Send(&String{Message: "hi"}); err != nil {
		panic(err)
	}
	if err := stream.Send(&String{Message: "hello"}); err != nil {
		panic(err)
	}
	stream.CloseSend()
	stream.RecvMsg(nil)
	fmt.Println("dun")
}

func perf() {
	conn, err := grpc.Dial(fmt.Sprintf("0.0.0.0:%s", port),
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                time.Second * 30,
			PermitWithoutStream: true,
		}))
	if err != nil {
		panic(err)
	}
	client := NewTestClient(conn)
	md := metadata.New(map[string]string{"x-alation-connector": "1"})
	stream, err := client.Performance(metadata.NewOutgoingContext(context.Background(), md), &String{Message: "hi"})
	if err != nil {
		panic(err)
	}
	if _, err := stream.Recv(); err != nil {
		panic(err)
	}
	start := time.Now()
	for _, err := stream.Recv(); err == nil; _, err = stream.Recv() {
	}
	fmt.Println(time.Now().Sub(start))
}
