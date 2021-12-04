package main // import "fake_connector/connector"
import (
	context "context"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.TraceLevel)
	s := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		Time: time.Second * 30,
	}), grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
		MinTime:             time.Second * 10,
		PermitWithoutStream: true,
	}))
	listener, err := net.Listen("tcp", "0.0.0.0:10001")
	if err != nil {
		panic(err)
	}
	RegisterTestServer(s, &Connector{})
	log.Info("starting connector on port 10001")
	s.Serve(listener)
}

type Connector struct{}

func (c Connector) Unary(ctx context.Context, s *String) (*String, error) {
	response := randomString()
	log.Infof("Got Unary request: '%s' adding '%s'", s.Message, response)
	str := &String{Message: fmt.Sprintf("%s %s", s.Message, response)}
	return str, nil
}

func (c Connector) ServerStream(s *String, server Test_ServerStreamServer) error {
	log.Infof("Got ServerStream request: '%s'", s.Message)
	for i := 0; i < 5; i++ {
		time.Sleep(time.Millisecond * time.Duration(rand.Int63n(2000)))
		str := &String{Message: fmt.Sprintf("%s %s", s.Message, randomString())}
		log.Infof("Sending '%s'", str.Message)
		err := server.Send(str)
		if err != nil {
			log.Errorf("ServerStream error: %v", err)
			return err
		}
	}
	return nil
}

func (c Connector) ClientStream(server Test_ClientStreamServer) error {
	log.Info("Got ClientStream request")
	for {
		str := &String{}
		err := server.RecvMsg(str)
		switch err {
		case nil:
			log.Infof("Got response '%s'", str.Message)
		case io.EOF:
			log.Info("Shutdown from client")
			return nil
		case io.ErrClosedPipe:
			log.Info("Closed pipe from client")
			return nil
		default:
			log.Infof("Some other error from client: %v", err)
			return err
		}
	}
}

func (c Connector) Bidirectional(server Test_BidirectionalServer) error {
	for i := 1; i <= 2; i++ {
		message, err := server.Recv()
		if err != nil {
			log.Errorf("BiDirection recv %d err %v", i, err)
			return err
		}
		response := randomString()
		log.Infof("Got Bidirectional request: '%s' adding '%s'", message.Message, response)
		str := &String{Message: fmt.Sprintf("%s %s", message.Message, response)}
		err = server.Send(str)
		if err != nil {
			log.Errorf("BiDirection send %d err %v", i, err)
			return err
		}
	}

	return nil
}

func (c Connector) IntentionalError(s *String, server Test_IntentionalErrorServer) error {
	return status.New(codes.DataLoss, "THIS IS INTENTIONALLY BAD").Err()
}

func (c Connector) Performance(s *String, server Test_PerformanceServer) error {
	b := strings.Builder{}
	for i := b.Len(); i < 1024*50; i = b.Len() {
		b.WriteString(randomString())
	}
	length := b.Len()
	str := &String{Message: b.String()}
	gigabyte := 1024 * 1024 * 1024
	for i := gigabyte; i > 0; i -= length {
		err := server.Send(str)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c Connector) mustEmbedUnimplementedTestServer() {
	//TODO implement me
	panic("implement me")
}

func randomString() string {
	length := rand.Intn(16)
	s := strings.Builder{}
	for i := 0; i < length; i++ {
		s.WriteByte(byte(rand.Intn(42) + 48))
	}
	return s.String()
}
