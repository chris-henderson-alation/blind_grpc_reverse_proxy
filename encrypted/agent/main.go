package main // import "github.com/chris-henderson-alation/blind_grpc_reverse_proxy"

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"time"

	"github.com/Alation/alation_connector_manager/security/common"

	"github.com/alecthomas/units"
	"github.com/chris-henderson-alation/blind_grpc_reverse_proxy/rpc"
	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var myPrivateKey [32]byte
var theirPublicKey [32]byte

func main() {
	key, err := common.GetPrivateKey("/home/chris/alation/blind_grpc_reverse_proxy/encrypted/agent/private_key.pem", []byte("password@1234"))
	if err != nil {
		panic(err)
	}
	private := key.(*ecdsa.PrivateKey)
	pubkey, err := common.GetPublicCert("/home/chris/alation/blind_grpc_reverse_proxy/encrypted/servicer/public_certificate.pem")
	if err != nil {
		panic(err)
	}
	public := pubkey.PublicKey.(*ecdsa.PublicKey)
	myPrivateKey = common.Marshal(private)
	theirPublicKey = common.Marshal(public)
	for {
		pollConnectorJobs()
		time.Sleep(time.Second * 10)
	}
}

func pollConnectorJobs() {
	log.Info("Begin polling for jobs.")
	c, err := grpc.Dial("0.0.0.0:8080", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		fmt.Println("Oh noes! ", err)
		return
	}
	defer c.Close()
	client := rpc.NewServicerClient(c)
	job, err := client.PollJob(context.Background(), &empty.Empty{})
	if err != nil {
		fmt.Println("Oh noes some more! ", err)
		return
	}
	switch j := job.Options.(type) {
	case *rpc.JobOptional_Job:
		Execute(j.Job, c)
	case *rpc.JobOptional_Nothing:
		log.Info("No jobs enqueued.")
	default:
		fmt.Printf("...what %t\n", j)
	}
}

func Execute(job *rpc.Job, c *grpc.ClientConn) {
	port := job.Headers.Connector + 10000
	client, err := grpc.Dial(fmt.Sprintf("0.0.0.0:%v", port), grpc.WithInsecure(), grpc.WithBlock(), grpc.WithDefaultCallOptions(grpc.ForceCodec(NoopCodec{})))
	defer client.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	switch job.Headers.Type {
	case rpc.GrpcCallType_Unary:
		unary(client, job, c)
	case rpc.GrpcCallType_ServerStream:
		stream(client, job, c)
	}

}

func unary(client *grpc.ClientConn, job *rpc.Job, c *grpc.ClientConn) {
	log.Infof("Beginning Unary at %s/%s", client.Target(), job.Headers.Method)
	upstream := rpc.NewServicerClient(c)
	configuration, err := common.Decrypt(job.Body.Ciphertext, job.Body.Nonce, &myPrivateKey, &theirPublicKey)
	if err != nil {
		log.Error("Failed configuration decryption")
	}
	dst := make([]byte, 0)
	err = client.Invoke(context.Background(), job.Headers.Method, configuration, &dst)
	if err != nil {
		fmt.Println("darn", err)
		return
	}
	ciphertext, nonce := common.Encrypt(dst, &theirPublicKey, &myPrivateKey)
	_, err = upstream.PushStream(context.Background(), &rpc.PushStreamItem{
		Headers:  job.Headers,
		Sequence: 1,
		Body: []*rpc.EncryptedBody{{
			Ciphertext: ciphertext,
			Nonce:      nonce,
		}},
	})
	if err != nil {
		log.Error(err)
	}
}

func stream(client *grpc.ClientConn, job *rpc.Job, c *grpc.ClientConn) {
	log.Infof("Beginning ServerStream at %s/%s", client.Target(), job.Headers.Method)
	upstream := rpc.NewServicerClient(c)
	stream, err := client.NewStream(context.Background(), &grpc.StreamDesc{ServerStreams: true}, job.Headers.Method)
	if err != nil {
		fmt.Println("dammit ", err)
		return
	}
	configuration, err := common.Decrypt(job.Body.Ciphertext, job.Body.Nonce, &myPrivateKey, &theirPublicKey)
	if err != nil {
		log.Error("Failed configuration decryption")
	}
	err = stream.SendMsg(configuration)
	if err != nil {
		fmt.Println("dammit 2", err)
		return
	}
	stream.CloseSend()
	var total units.Base2Bytes = 0
	var running units.Base2Bytes = 0
	var sequence uint64 = 0
	buffer := make([][]byte, 0)
	threePointFiveMiB := units.MiB*3 + units.MiB/2
	start := time.Now()
	var timeSpentPulling time.Duration
	var timeSpentPushing time.Duration
	var timeSpentEncrypting time.Duration
	for {
		dst := make([]byte, 0)
		s := time.Now()
		err := stream.RecvMsg(&dst)
		if err != nil {
			break
		}
		timeSpentPulling += time.Now().Sub(s)
		received := units.Base2Bytes(len(dst))
		if running+received >= threePointFiveMiB {
			s := time.Now()
			total += running
			envelopes := make([]*rpc.EncryptedBody, len(buffer))
			for i, b := range buffer {
				es := time.Now()
				ciphertext, nonce := common.Encrypt(b, &theirPublicKey, &myPrivateKey)
				timeSpentEncrypting += time.Now().Sub(es)
				envelopes[i] = &rpc.EncryptedBody{
					Ciphertext: ciphertext,
					Nonce:      nonce,
				}
			}
			_, err := upstream.PushStream(context.Background(), &rpc.PushStreamItem{
				Headers:  job.Headers,
				Sequence: sequence,
				Body:     envelopes,
			})
			if err != nil {
				fmt.Println("noooo")
				return
			}
			sequence += 1
			timeSpentPushing += time.Now().Sub(s)
			buffer = make([][]byte, 0)
			running = 0
		}
		running += received
		buffer = append(buffer, dst)
	}
	if len(buffer) > 0 {
		total += running
		s := time.Now()
		envelopes := make([]*rpc.EncryptedBody, len(buffer))
		for i, b := range buffer {
			es := time.Now()
			ciphertext, nonce := common.Encrypt(b, &theirPublicKey, &myPrivateKey)
			timeSpentEncrypting += time.Now().Sub(es)
			envelopes[i] = &rpc.EncryptedBody{
				Ciphertext: ciphertext,
				Nonce:      nonce,
			}
		}
		_, err := upstream.PushStream(context.Background(), &rpc.PushStreamItem{
			Headers:  job.Headers,
			Sequence: sequence,
			Body:     envelopes,
		})
		if err != nil {
			fmt.Println("noooo")
			return
		}
		timeSpentPushing += time.Now().Sub(s)
	}
	runtime := time.Now().Sub(start)
	log.Infof("Processed %v bytes in %s (%s/second)", total, runtime, units.Base2Bytes(float64(total)/runtime.Seconds()))
	log.Infof("Time spend pulling data was %s", timeSpentPulling)
	log.Infof("Time spent pushing was %s", timeSpentPushing)
	log.Infof("Time spent encrypting was %s", timeSpentEncrypting)
}

type NoopCodec struct{}

func (cb NoopCodec) Name() string {
	return "NoopCodec"
}

func (cb NoopCodec) Marshal(v interface{}) ([]byte, error) {
	return v.([]byte), nil
}

func (cb NoopCodec) Unmarshal(data []byte, v interface{}) error {
	b, _ := v.(*[]byte)
	*b = data
	return nil
}
