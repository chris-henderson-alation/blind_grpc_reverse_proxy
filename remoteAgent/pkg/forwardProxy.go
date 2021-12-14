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

func NewForwardProxy() *ForwardProxy {
	return &ForwardProxy{
		agents: NewAgents(),
	}
}

func (proxy *ForwardProxy) Submit(call *GrpcCall) error {
	return proxy.agents.Submit(call)
}

func (proxy *ForwardProxy) JobTunnel(_ *empty.Empty, server ioc.GrpcInverter_JobTunnelServer) error {
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
	logrus.Infof("Agent ID %d connected", agentId)
	for {
		select {
		case <-server.Context().Done():
			logrus.Info("it's gone my guy context cancelled")
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
			proxy.agents.Enqueue(job, nil)
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
	job, _ := proxy.agents.Retrieve(agentId, jobId)
	if job == nil {
		return fmt.Errorf("nice catch blanco nino, but too bad your ass got saaaaaacked")
	}
	ioc.ForwardProxy(job.Alation, agent)
	return nil
}

func (proxy *ForwardProxy) StreamHandler() grpc.StreamHandler {
	return func(_ interface{}, upstream grpc.ServerStream) error {
		fullMethodName, ok := grpc.MethodFromServerStream(upstream)
		if !ok {
			return status.Error(codes.Internal, "lowLevelServerStream not exists in context")
		}
		agentId, err := ExtractAgentId(upstream.Context())
		if err != nil {
			agentId = 1
			//return err
		}
		connectorId, err := ExtractConnectorId(upstream.Context())
		if err != nil {
			return err
		}
		jobId, err := ExtractJobId(upstream.Context())
		if err != nil {
			jobId = 1
			//return err
		}
		logrus.Infof("Received %s dispatch to connector %d", fullMethodName, connectorId)
		errors := make(chan error)
		call := &GrpcCall{
			Method:    fullMethodName,
			Agent:     agentId,
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
	a.Error <- s.Err()
}

type GrpcCall struct {
	ioc.Alation
	Method    string
	Connector uint64
	Agent     uint64
	JobId     uint64
	Error     chan error
}
