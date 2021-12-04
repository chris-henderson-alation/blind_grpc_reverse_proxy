package main

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"google.golang.org/grpc/metadata"

	"github.com/chris-henderson-alation/blind_grpc_reverse_proxy/rpc"
	"github.com/golang/protobuf/ptypes/empty"
)

type CallId = uint64

type InternetFacing struct {
	inbound  chan GrpcCall
	outbound map[CallId]GrpcCall

	counter uint64

	lock sync.Mutex
}

func (i *InternetFacing) StreamHandler() grpc.StreamHandler {
	do := func(_ interface{}, upstream grpc.ServerStream) error {
		fmt.Println("interesting")
		fullMethodName, ok := grpc.MethodFromServerStream(upstream)
		if !ok {
			return grpc.Errorf(codes.Internal, "lowLevelServerStream not exists in context")
		}
		md, ok := metadata.FromIncomingContext(upstream.Context())
		if !ok {
			return status.Error(codes.FailedPrecondition, "no friggin' headers my ")
		}
		var connector uint64
		if id, err := strconv.ParseUint(md.Get("x-alation-connector")[0], 10, 64); err != nil {
			return status.Errorf(codes.FailedPrecondition, "come on, bad connector id %v", err)
		} else {
			connector = id
		}
		err := make(chan error)
		call := GrpcCall{
			Method:    fullMethodName,
			Connector: connector,
			Stream:    upstream,
			Error:     err,
		}
		i.lock.Lock()
		i.counter += 1
		i.outbound[i.counter] = call
		call.JobId = i.counter
		i.lock.Unlock()
		i.inbound <- call
		return <-err
	}
	return do
}

func (i *InternetFacing) InitiateJob(_ *empty.Empty, server rpc.Agent_InitiateJobServer) error {
	fmt.Println("...someone's knocking...")
	for {
		call := <-i.inbound
		fmt.Println("let's goooo")
		go func() {
			job := &rpc.Job{
				Method:    call.Method,
				JobID:     call.JobId,
				Connector: call.Connector,
			}
			err := server.Send(job)
			if err != nil {
				fmt.Println(err)
				i.lock.Lock()
				delete(i.outbound, call.JobId)
				i.lock.Unlock()
				call.Error <- status.Errorf(codes.Unavailable, "Agent is unavailable because of %v", err)
				close(call.Error)
			}
		}()
	}
}

func (i *InternetFacing) JobCallback(downstream rpc.Agent_JobCallbackServer) error {
	md, ok := metadata.FromIncomingContext(downstream.Context())
	if !ok {
		panic("wat")
	}
	id, err := strconv.ParseUint(md.Get("x-alation-job-id")[0], 10, 64)
	if err != nil {
		panic(err)
	}
	i.lock.Lock()
	upstreamClient := i.outbound[id]
	delete(i.outbound, id)
	i.lock.Unlock()
	upstream := upstreamClient.Stream
	upstreamError := upstreamClient.Error
	go func() {
		// Pipe upstreamClient messages to downstream.
		var msg []byte
		for {
			err := upstream.RecvMsg(&msg)
			switch err {
			case nil:
				// Account for downstream hung up.
				downstream.Send(&rpc.Body{Body: msg})
			case io.EOF:
				// Account for downstream hung up.
				downstream.Send(&rpc.Body{EOF: true})
				return
			default:
				return
			}
		}
	}()
	defer close(upstreamError)
	for {
		body, err := downstream.Recv()
		switch err {
		case nil:
			if body.Error != nil {
				upstreamError <- status.Error(codes.Code(body.Error.Code), body.Error.Description)
				return nil
			}
			upstream.SendMsg(body.Body)
		case io.EOF:
			return nil
		case io.ErrClosedPipe:
			upstreamClient.Error <- status.Error(codes.Canceled, "it hung up, just do something about this")
			return nil
		default:
			upstreamClient.Error <- err
			return nil
		}
	}

}

type GrpcCall struct {
	Method    string
	Connector uint64
	JobId     uint64
	Stream    grpc.ServerStream
	Error     chan error
}

func NewInternetFacing(port string) *InternetFacing {
	s := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		Time: time.Second * 30,
	}), grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
		MinTime:             time.Second * 10,
		PermitWithoutStream: true,
	}))
	i := &InternetFacing{
		inbound:  make(chan GrpcCall),
		outbound: map[CallId]GrpcCall{},
		counter:  0,
		lock:     sync.Mutex{},
	}
	rpc.RegisterAgentServer(s, i)
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	go func() {
		fmt.Printf("starting external facing server on :%s\n", port)
		s.Serve(listener)
	}()
	return i
}
