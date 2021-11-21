package main // import "github.com/chris-henderson-alation/blind_grpc_reverse_proxy"

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/chris-henderson-alation/blind_grpc_reverse_proxy/rpc"

	"google.golang.org/grpc"
)

func main() {
	for {
		pollConnectorJobs()
		time.Sleep(time.Second * 5)
	}
}

func pollConnectorJobs() {
	jobs := make([]*rpc.Job, 0)
	for {
		fmt.Println("connecting")
		resp, err := http.Get("http://localhost:8080/connectors")
		fmt.Println("got response")
		if err != nil {
			// Put into exponentional backoff retry loop?
			// Otherwise stop trying anc move on with the jobs
			// we may have gotten so far.
			fmt.Println(err)
			break
		}
		if resp.StatusCode == http.StatusNoContent {
			fmt.Println("no jobs")
			// All jobs queued jobs have been consumed.
			break
		}
		if resp.StatusCode == http.StatusInternalServerError {
			fmt.Println("darn demos")
			break
		}
		job, err := func() (*rpc.Job, error) {
			defer resp.Body.Close()
			job := &rpc.Job{}
			return job, json.NewDecoder(resp.Body).Decode(job)
		}()
		jobs = append(jobs, job)
		break
	}
	for _, job := range jobs {
		go Execute(job)
	}
}

func Execute(job *rpc.Job) {
	port := job.Headers.Connector + 10000
	client, err := grpc.Dial(fmt.Sprintf("0.0.0.0:%v", port), grpc.WithInsecure(), grpc.WithBlock(), grpc.WithDefaultCallOptions(grpc.ForceCodec(NoopCodec{})))
	defer client.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	switch job.Headers.Type {
	case rpc.Unary:
		unary(client, job)
	case rpc.UnaryStream:
		stream(client, job)
	}

}

func unary(client *grpc.ClientConn, job *rpc.Job) {
	dst := make([]byte, 0)

	err := client.Invoke(context.Background(), job.Headers.Method, []uint8(job.Body), &dst)
	if err != nil {
		fmt.Println("darn", err)
		return
	}
	fmt.Println(string(dst))
}

func stream(client *grpc.ClientConn, job *rpc.Job) {
	stream, err := client.NewStream(context.Background(), &grpc.StreamDesc{ServerStreams: true}, job.Headers.Method)
	if err != nil {
		fmt.Println("dammit ", err)
		return
	}
	err = stream.SendMsg([]uint8(job.Body))
	if err != nil {
		fmt.Println("dammit 2", err)
		return
	}
	stream.CloseSend()
	for {
		dst := make([]byte, 0)
		err := stream.RecvMsg(&dst)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(dst))
	}

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
