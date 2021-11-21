package main // import "github.com/chris-henderson-alation/blind_grpc_reverse_proxy/servicer"

import (
	context "context"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"sync"

	"github.com/chris-henderson-alation/blind_grpc_reverse_proxy/rpc"
	ocf "github.com/chris-henderson-alation/blind_grpc_reverse_proxy/servicer_connector_proto"
	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"

	"github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
)

//go:generate protoc -I . --go_out=plugins=grpc:. tennis.proto

type Servicer struct{}

func (s Servicer) PushStream(ctx context.Context, item *rpc.PushStreamItem) (*empty.Empty, error) {
	// just black hole it for the purpose of the PoC, this isn't what we're testing
	go func() {
		for _, item := range item.Body {
			target := &ocf.Bytes{
				Inner: make([]byte, 0),
			}
			err := proto.Unmarshal(item, target)
			if err != nil {
				fmt.Println("darnit Chris")
			}
		}
	}()
	return &empty.Empty{}, nil
}

func (s Servicer) PollJob(ctx context.Context, _ *empty.Empty) (*rpc.JobOptional, error) {
	job, err := makeBigJob()
	if err != nil {
		return nil, err
	}
	log.Infof("Sending job for %v", job)
	return job, nil
}

func makeRapidJob() (*rpc.JobOptional, error) {
	ball, err := marshal(makeBall())
	if err != nil {
		return nil, err
	}
	return &rpc.JobOptional{
		Options: &rpc.JobOptional_Job{
			Job: &rpc.Job{
				Headers: &rpc.Headers{
					Method:    "ocf.Tennis/Rapid",
					Type:      rpc.GrpcCallType_ServerStream,
					JobID:     jobID(),
					Connector: 1,
				},
				Body: ball,
			}},
	}, nil
}

func makeBigJob() (*rpc.JobOptional, error) {
	body, err := marshal(&empty.Empty{})
	if err != nil {
		return nil, err
	}
	return &rpc.JobOptional{
		Options: &rpc.JobOptional_Job{
			Job: &rpc.Job{
				Headers: &rpc.Headers{
					Method:    "ocf.Tennis/Big",
					Type:      rpc.GrpcCallType_ServerStream,
					JobID:     jobID(),
					Connector: 1,
				},
				Body: body,
			}},
	}, nil
}

func main() {
	p := "8080"
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", p))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	rpc.RegisterServicerServer(s, &Servicer{})
	s.Serve(lis)
}

func makeBall() *ocf.Ball {
	return &ocf.Ball{
		Ball: randomString(),
		Complex: &ocf.Complex{
			Array:  []uint32{1, 2, 3, 4, 5},
			Corpus: ocf.Complex_IMAGES,
			Result: &ocf.Complex_Result{
				Url:      "https://alation.com",
				Title:    "Winners",
				Snippets: []string{randomString(), randomString(), randomString()},
			},
		},
	}
}

var id uint64 = 0
var idLock = sync.Mutex{}

func jobID() uint64 {
	idLock.Lock()
	defer idLock.Unlock()
	id += 1
	return id
}

func marshal(v interface{}) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}

func randomString() string {
	length := rand.Intn(64)
	s := strings.Builder{}
	for i := 0; i < length; i++ {
		s.WriteByte(byte(rand.Intn(42) + 48))
	}
	return s.String()
}
