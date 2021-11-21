package main // import "github.com/chris-henderson-alation/blind_grpc_reverse_proxy/servicer"

import (
	context "context"
	"crypto/ecdsa"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"sync"

	"github.com/Alation/alation_connector_manager/security/common"
	"google.golang.org/grpc"

	"github.com/chris-henderson-alation/blind_grpc_reverse_proxy/rpc"
	ocf "github.com/chris-henderson-alation/blind_grpc_reverse_proxy/servicer_connector_proto"
	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"

	"github.com/golang/protobuf/proto"
)

//go:generate protoc -I . --go_out=plugins=grpc:. tennis.proto

var myPrivateKey [32]byte
var theirPublicKey [32]byte

func main() {
	key, err := common.GetPrivateKey("/home/chris/alation/blind_grpc_reverse_proxy/encrypted/servicer/private_key.pem", []byte("password@1234"))
	if err != nil {
		panic(err)
	}
	private := key.(*ecdsa.PrivateKey)
	pubkey, err := common.GetPublicCert("/home/chris/alation/blind_grpc_reverse_proxy/encrypted/agent/public_certificate.pem")
	if err != nil {
		panic(err)
	}
	public := pubkey.PublicKey.(*ecdsa.PublicKey)
	myPrivateKey = common.Marshal(private)
	theirPublicKey = common.Marshal(public)
	p := "8080"
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", p))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	rpc.RegisterServicerServer(s, &Servicer{})
	s.Serve(lis)
}

type Servicer struct{}

func (s Servicer) PushStream(ctx context.Context, item *rpc.PushStreamItem) (*empty.Empty, error) {
	// just black hole it for the purpose of the PoC, this isn't what we're testing
	go func() {
		messages := make([][]byte, len(item.Body))
		for i, envelope := range item.Body {
			plaintext, err := common.Decrypt(envelope.Ciphertext, envelope.Nonce, &myPrivateKey, &theirPublicKey)
			if err != nil {
				log.Error("Failed decryption")
				return
			}
			messages[i] = plaintext
		}
		switch item.Headers.Method {
		case "ocf.Tennis/Big":
			// unmarshal it but don't print it, just junk stuff
			target := &ocf.Bytes{
				Inner: make([]byte, 0),
			}
			err := proto.Unmarshal(messages[0], target)
			if err != nil {
				fmt.Println("darnit Chris")
			}
		case "ocf.Tennis/Ping":
			target := &ocf.Ball{}
			err := proto.Unmarshal(messages[0], target)
			if err != nil {
				fmt.Println("darnit Chris")
				return
			}
			log.Info(target)
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

func makePingJob() (*rpc.JobOptional, error) {
	ball, err := marshal(makeBall())
	if err != nil {
		return nil, err
	}
	return &rpc.JobOptional{
		Options: &rpc.JobOptional_Job{
			Job: &rpc.Job{
				Headers: &rpc.Headers{
					Method:    "ocf.Tennis/Ping",
					Type:      rpc.GrpcCallType_Unary,
					JobID:     jobID(),
					Connector: 1,
				},
				Body: ball,
			}},
	}, nil
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

func marshal(v interface{}) (*rpc.EncryptedBody, error) {
	plaintext, err := proto.Marshal(v.(proto.Message))
	if err != nil {
		return nil, err
	}
	ciphertext, nonce := common.Encrypt(plaintext, &theirPublicKey, &myPrivateKey)
	return &rpc.EncryptedBody{
		Ciphertext: ciphertext,
		Nonce:      nonce[:],
	}, nil
}

func randomString() string {
	length := rand.Intn(64)
	s := strings.Builder{}
	for i := 0; i < length; i++ {
		s.WriteByte(byte(rand.Intn(42) + 48))
	}
	return s.String()
}
