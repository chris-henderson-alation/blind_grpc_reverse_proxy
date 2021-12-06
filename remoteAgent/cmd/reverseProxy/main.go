package main

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/Alation/alation_connector_manager/docker/remoteAgent/grpcinverter"

	"github.com/Alation/alation_connector_manager/docker/remoteAgent/grpcinverter/ioc"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

const agentId = 1

func main() {
	client := NewClient("1234")
	jobs, err := client.JobTunnel(grpcinverter.NewHeaderBuilder().SetAgentId(agentId).Build(context.Background()), &empty.Empty{})
	if err != nil {
		panic(err)
	}
	for {
		job, err := jobs.Recv()
		if err != nil {
			panic(err)
		}
		go dispatch(job)
	}
}

func NewClient(port string) ioc.GrpcInverterClient {
	//conn, err := grpc.Dial(fmt.Sprintf("fresh-crane.alation-test.com:%s", port),
	//	grpc.WithInsecure(),
	//	grpc.WithKeepaliveParams(shared.KEEPALIVE_CLIENT_PARAMETERS))
	conn, err := grpc.Dial(fmt.Sprintf("0.0.0.0:%s", port),
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(grpcinverter.KEEPALIVE_CLIENT_PARAMETERS))
	if err != nil {
		panic(err)
	}
	return ioc.NewGrpcInverterClient(conn)
}

func dispatch(job *ioc.Job) {
	logrus.Infof("Received job %v", job)
	headers := grpcinverter.NewHeaderBuilder().
		SetConnectorId(job.Connector).
		SetAgentId(agentId).
		SetJobId(job.JobID).
		Build(context.Background())
	c := NewClient("1234")
	upstream, err := c.Pipe(headers)
	if err != nil {
		panic(err)
	}
	conn, err := grpc.Dial(fmt.Sprintf("0.0.0.0:%d", job.Connector+10000),
		grpc.WithKeepaliveParams(grpcinverter.KEEPALIVE_CLIENT_PARAMETERS),
		grpc.WithInsecure())
	if err != nil {
		fmt.Println("booo tits")
		upstream.SendMsg(err)
		return
	}
	downstream, err := conn.NewStream(context.Background(), &grpc.StreamDesc{
		ServerStreams: true,
		ClientStreams: true,
	}, job.Method, grpc.ForceCodec(grpcinverter.NoopCodec{}))
	if err != nil {
		fmt.Println("booo")
		upstream.SendMsg(err)
		return
	}
	ioc.ReverseProxy(upstream, downstream)
}
