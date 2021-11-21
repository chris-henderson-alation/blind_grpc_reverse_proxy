package main // import "test_connector"

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"math/rand"
	"net"
	"path/filepath"

	"google.golang.org/grpc"
)

//go:generate protoc -I . --go_out=plugins=grpc:. tennis.proto

type Tennis struct{}

func (t *Tennis) Rapid(ball *Ball, server Tennis_RapidServer) error {
	fmt.Println(ball)
	for i := 0; i < 10; i++ {
		err := server.Send(&Ball{Ball: "pong"})
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Tennis) Ping(ctx context.Context, ball *Ball) (*Ball, error) {
	fmt.Println("Playing tennis. Got ", ball.Ball)
	return &Ball{Ball: "pong"}, nil
}

var storage = filepath.Join("/", "storage")

// Writes some data to disk to test storage facilities
func (t *Tennis) Write(ctx context.Context, msg *Incoming) (*Status, error) {
	fmt.Printf("Got request to write %s\n", msg.Message)
	fname := fmt.Sprintf("%X", big.NewInt(rand.Int63n(int64(math.Pow(2, 32)))).Bytes())
	err := ioutil.WriteFile(filepath.Join(storage, fname), []byte(msg.Message), 0660)
	if err != nil {
		fmt.Println(err)
	}
	return &Status{Fname: fname, Ok: err == nil, Error: fmt.Sprintf("%s", err)}, nil
}

func main() {
	//p := os.Getenv("PORT")
	//if p == "" {
	//	p = "8080"
	//}
	p := "10001"
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", p))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	RegisterTennisServer(s, &Tennis{})
	s.Serve(lis)
}
