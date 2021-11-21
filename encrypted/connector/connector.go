package main // import "test_connector"

import (
	"context"
	"crypto/rand"
	"fmt"
	"net"

	"github.com/alecthomas/units"
	ocf "github.com/chris-henderson-alation/blind_grpc_reverse_proxy/servicer_connector_proto"
	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

//go:generate protoc -I . --go_out=plugins=grpc:. tennis.proto

type Tennis struct{}

func (t *Tennis) Big(_ *empty.Empty, server ocf.Tennis_BigServer) error {
	log.Infof("Beginning 1GiB stream of random data")
	b := make([]byte, units.KiB*50)
	rand.Read(b)
	target := units.MiB * 750
	for i := units.GiB; i >= 0; i -= units.KiB * 50 {
		err := server.Send(&ocf.Bytes{
			Inner: b,
		})
		if err != nil {
			return err
		}
		if i < target {
			sent := units.GiB - target
			log.Infof("%s sent.", sent)
			target = target - units.MiB*250
		}
	}
	log.Info("~1GiB of data sent. Stream complete.")
	return nil
}

func (t *Tennis) Rapid(ball *ocf.Ball, server ocf.Tennis_RapidServer) error {
	log.Debug(ball)
	for i := 0; i < 10; i++ {
		err := server.Send(&ocf.Ball{Ball: "pong"})
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Tennis) Ping(ctx context.Context, ball *ocf.Ball) (*ocf.Ball, error) {
	log.Info("Playing tennis. Got ", ball.Ball)
	return &ocf.Ball{Ball: "pong"}, nil
}

func main() {
	//p := os.Getenv("PORT")
	//if p == "" {
	//	p = "8080"
	//}
	p := "10001"
	log.Infof("Begin listening on 0.0.0.0:%s", p)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", p))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	ocf.RegisterTennisServer(s, &Tennis{})
	s.Serve(lis)
}
