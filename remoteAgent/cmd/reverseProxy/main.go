package main

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/Alation/alation_connector_manager/docker/remoteAgent/grpcinverter"

	"github.com/Alation/alation_connector_manager/docker/remoteAgent/grpcinverter/ioc"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

const agentId = 1
const proxyAddress = "0.0.0.0"
const proxyPort = "1234"

func main() {
	client := NewClient()
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

func NewClient() ioc.GrpcInverterClient {
	//conn, err := grpc.Dial(fmt.Sprintf("fresh-crane.alation-test.com:%s", port),
	//	grpc.WithInsecure(),
	//	grpc.WithKeepaliveParams(shared.KEEPALIVE_CLIENT_PARAMETERS))
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", proxyAddress, proxyPort),
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
	c := NewClient()
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
	var bookmark *ioc.Message
	var shouldRetry bool
	for {
		bookmark, shouldRetry = ioc.ReverseProxy(upstream, downstream, bookmark)
		if shouldRetry {
			logrus.Error("Failed! Trying again in a minute!")
			time.Sleep(time.Minute)
			c = NewClient()
			upstream, err = c.Pipe(headers)
			if err != nil {
				panic("dunno, loops on loops")
			}
			logrus.Info("yoooo reconnected!")
		} else {
			return
		}
	}

}
