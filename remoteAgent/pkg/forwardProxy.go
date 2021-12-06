package grpcinverter

import (
	"fmt"

	"github.com/Alation/alation_connector_manager/docker/remoteAgent/grpcinverter/ioc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sirupsen/logrus"
)

type ForwardProxy struct {
	agents *Agents
	ioc.UnimplementedGrpcInverterServer
}

func (proxy *ForwardProxy) Submit(call *GrpcCall) error {
	return proxy.agents.Submit(call)
}

func (proxy *ForwardProxy) JobTunnel(empty *empty.Empty, server ioc.GrpcInverter_JobTunnelServer) error {
	agentId, err := ExtractAgentId(server.Context())
	if err != nil {
		logrus.Error(err)
		return err
	}
	jobs, alreadyRegistered := proxy.agents.Register(agentId)
	if alreadyRegistered {
		return fmt.Errorf(
			"an agent with ID %d is already listeners for events, "+
				"this agent has been rejected as being a duplicate", agentId)
	}
	for {
		select {
		case <-server.Context().Done():
			proxy.agents.Unregister(agentId)
			return nil
		case job := <-jobs:
			// let's goooo
			err := server.Send(&ioc.Job{
				JobID:     job.JobId,
				Method:    job.Method,
				Connector: job.Connector,
			})
			if err != nil {
				s, ok := status.FromError(err)
				if !ok {
					job.SendError(status.New(codes.Unavailable, err.Error()))
				} else {
					job.SendError(s)
				}
				continue
			}
			proxy.agents.Enqueue(job)
		}
	}

}

func (proxy *ForwardProxy) Pipe(agent ioc.GrpcInverter_PipeServer) error {
	agentId, err := ExtractAgentId(agent.Context())
	if err != nil {
		return err
	}
	jobId, err := ExtractJobId(agent.Context())
	if err != nil {
		return err
	}
	job := proxy.agents.Retrieve(agentId, jobId)
	if job == nil {
		return fmt.Errorf("nice catch blanco nino, but too bad your ass got saaaaaacked")
	}
	ioc.ForwardProxy(job.Alation, agent)
	ok := true
	if !ok {
		// put the job back in just case the agent calls back.
		proxy.agents.Enqueue(job)
	}
	return nil
}

func (proxy *ForwardProxy) StreamHandler() grpc.StreamHandler {
	return func(_ interface{}, upstream grpc.ServerStream) error {
		fullMethodName, ok := grpc.MethodFromServerStream(upstream)
		if !ok {
			return status.Error(codes.Internal, "lowLevelServerStream not exists in context")
		}
		angentId, err := ExtractAgentId(upstream.Context())
		if err != nil {
			return err
		}
		connectorId, err := ExtractConnectorId(upstream.Context())
		if err != nil {
			return err
		}
		jobId, err := ExtractJobId(upstream.Context())
		if err != nil {
			return err
		}
		logrus.Infof("Received %s dispatch to connector %d", fullMethodName, connectorId)
		errors := make(chan error)
		call := &GrpcCall{
			Method:    fullMethodName,
			Agent:     angentId,
			Connector: connectorId,
			JobId:     jobId,
			Alation: &UpstreamAlation{
				ServerStream: upstream,
				Error:        errors,
			},
		}
		err = proxy.Submit(call)
		if err != nil {
			return err
		}
		return <-errors
	}
}

type UpstreamAlation struct {
	grpc.ServerStream
	Error chan error
}

func (a *UpstreamAlation) SendError(s *status.Status) {
	a.Error <- status.Error(s.Code(), s.Message())
}

type GrpcCall struct {
	ioc.Alation
	Method    string
	Connector uint64
	Agent     uint64
	JobId     uint64
	Error     chan error
}
