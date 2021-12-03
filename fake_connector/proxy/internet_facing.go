package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"sync"
	"time"

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
	do := func(_ interface{}, stream grpc.ServerStream) error {
		fmt.Println("interesting")
		fullMethodName, ok := grpc.MethodFromServerStream(stream)
		if !ok {
			return grpc.Errorf(codes.Internal, "lowLevelServerStream not exists in context")
		}
		md, ok := metadata.FromIncomingContext(stream.Context())
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
			Stream:    stream,
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

func (i *InternetFacing) JobCallback(server rpc.Agent_JobCallbackServer) error {
	fmt.Println("is back")
	md, ok := metadata.FromIncomingContext(server.Context())
	if !ok {
		panic("wat")
	}
	id, err := strconv.ParseUint(md.Get("x-alation-job-id")[0], 10, 64)
	if err != nil {
		panic(err)
	}
	i.lock.Lock()
	upstream := i.outbound[id]
	delete(i.outbound, id)
	i.lock.Unlock()
	defer close(upstream.Error)
	go func() {
		for {
			var msg []byte
			err := upstream.Stream.RecvMsg(&msg)
			switch err {
			case nil:
				fmt.Println("okay sending")
				server.Send(&rpc.Body{Body: msg})
			case io.EOF, io.ErrClosedPipe:
				return
			}
		}

	}()
	for {
		batch, err := server.Recv()
		switch err {
		case nil:
			fmt.Println("is data")
			for _, message := range batch.Bodies {
				upstream.Stream.SendMsg(message)
			}
			if batch.Error != nil {
				upstream.Error <- status.Error(codes.Code(batch.Error.Code), batch.Error.Description)
				return nil
			}
		case io.EOF:
			fmt.Println("is done")
			return nil
		case io.ErrClosedPipe:
			fmt.Println("is closed pip")
			upstream.Error <- status.Error(codes.Canceled, "it hung up, just do something about this")
			return nil
		default:
			fmt.Println("is other error")
			upstream.Error <- err
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

//func (i *InternetFacing) Send(call GrpcCall) <-chan error {
//	result := make(chan error)
//	i.lock.Lock()
//	i.counter += 1
//	i.outbound[i.counter] = result
//	call.JobId = i.counter
//	i.lock.Unlock()
//	i.inbound <- call
//	return result
//}

//type Initiate struct {
//}

//type GrpcResponse struct {
//	Message []byte
//	Error   error
//}

//type GrpcCallType int
//
//const (
//	Unary GrpcCallType = iota // Primary target
//	ServerStream
//	ClientStream
//	BiDirectionalStream
//)

//func (receiver GrpcCallType) ToProto() rpc.GrpcCallType {
//	switch receiver {
//	case Unary:
//		return rpc.GrpcCallType_Unary
//	case ServerStream:
//		return rpc.GrpcCallType_ServerStream
//	case ClientStream:
//		return rpc.GrpcCallType_ClientStream
//	case BiDirectionalStream:
//		return rpc.GrpcCallType_BiDirectionalStream
//	}
//	panic("non-exhaustive")
//}

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
