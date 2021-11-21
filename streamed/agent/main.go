package main // import "github.com/chris-henderson-alation/blind_grpc_reverse_proxy"

import (
	"context"
	"fmt"
	"time"

	"github.com/alecthomas/units"
	"github.com/chris-henderson-alation/blind_grpc_reverse_proxy/rpc"
	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
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
	dst := make([]byte, 0)

	err := client.Invoke(context.Background(), job.Headers.Method, job.Body, &dst)
	if err != nil {
		fmt.Println("darn", err)
		return
	}
	fmt.Println(string(dst))
}

func stream(client *grpc.ClientConn, job *rpc.Job, c *grpc.ClientConn) {
	log.Infof("Beginning ServerStream at %s/%s", client.Target(), job.Headers.Method)
	upstream, err := rpc.NewServicerClient(c).PushStream(context.Background())
	if err != nil {
		log.Error(err)
		return
	}
	stream, err := client.NewStream(context.Background(), &grpc.StreamDesc{ServerStreams: true}, job.Headers.Method)
	if err != nil {
		fmt.Println("dammit ", err)
		return
	}
	err = stream.SendMsg(job.Body)
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
			err = upstream.Send(&rpc.PushStreamItem{
				Headers:  job.Headers,
				Sequence: sequence,
				Body:     buffer,
			})
			if err != nil {
				log.Error(err)
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
		err = upstream.Send(&rpc.PushStreamItem{
			Headers:  job.Headers,
			Sequence: sequence,
			Body:     buffer,
		})
		if err != nil {
			log.Error(err)
		}
		timeSpentPushing += time.Now().Sub(s)
	}
	runtime := time.Now().Sub(start)
	log.Infof("Processed %v bytes in %s (%s/second)", total, runtime, units.Base2Bytes(float64(total)/runtime.Seconds()))
	log.Infof("Time spend pulling data was %s", timeSpentPulling)
	log.Infof("Time spent pushing was %s", timeSpentPushing)
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
