package ioc // import "github.com/Alation/alation_connector_manager/docker/remoteAgent/grpcinverter"

import (
	"io"
	"sync"

	codes "google.golang.org/grpc/codes"

	"github.com/sirupsen/logrus"

	"google.golang.org/grpc/status"

	grpc "google.golang.org/grpc"
)

type Alation interface {
	grpc.ServerStream
	SendError(*status.Status)
}

type Connector interface {
	grpc.ClientStream
}

type Forward interface {
	GrpcInverter_PipeClient
}

type Reverse interface {
	GrpcInverter_PipeServer
}

func ForwardProxy(alation Alation, agent Reverse) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		var body []byte
		for {
			err := alation.RecvMsg(&body)
			switch err {
			case nil:
				e := agent.Send(&Message{Body: body})
				if e != nil {
					logrus.Errorf("agent.Send(msg) failed with %v, type %T", e, e)
					return
				}
			case io.EOF:
				e := agent.Send(&Message{EOF: true})
				if e != nil {
					logrus.Errorf("agent.Send(EOF) failed with %v, type %T", e, e)
					return
				}
				return
			default:
				logrus.Errorf("other error from alation.RecvMsg(), %v", err)
				return
			}
		}
	}()
	go func() {
		defer wg.Done()
		for {
			body, err := agent.Recv()
			switch err {
			case nil:
				if body.Error != nil {
					alation.SendError(status.New(codes.Code(body.Error.Code), body.Error.Description))
					return
				}
				alation.SendMsg(body.Body)
			case io.EOF:
			case io.ErrClosedPipe:
			default:
			}
		}
	}()
	wg.Wait()
}

func ReverseProxy(upstream Forward, connector Connector) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for {
			body, err := upstream.Recv()
			switch err {
			case nil:
				if body.EOF {
					connector.CloseSend()
					return
				}
				err = connector.SendMsg(body.Body)
				switch err {
				case nil:
				case io.EOF:
					upstream.CloseSend()
					return
				}
			case io.EOF:
				connector.CloseSend()
				return
			default:
				logrus.Errorf("other error from upstream.Recv(), %v", err)
				return
			}
		}
	}()
	go func() {
		defer wg.Done()
		var msg []byte
		for {
			err := connector.RecvMsg(&msg)
			switch err {
			case nil:
				e := upstream.Send(&Message{Body: msg})
				if e != nil {
					logrus.Errorf("upstream.Send(msg) failed with %v, type %T", e, e)
					return
				}
			case io.EOF:
				e := upstream.CloseSend()
				if e != nil {
					logrus.Errorf("upstream.CloseSed() (from io.EOF) failed with %v, type %T", e, e)
					return
				}
				return
			default:
				// We are technically the client in this scenario and clients
				// cannot send errors in the world of gRPC. So we have to represent
				// it ourselves.
				serr, _ := status.FromError(err)
				resp := &Message{
					Error: &Error{
						Code:        uint32(serr.Code()),
						Description: serr.Message(),
					},
				}
				e := upstream.Send(resp)
				if e != nil {
					logrus.Errorf("upstream.Send(resp) failed with %v, type %T", e, e)
					return
				}
				e = upstream.CloseSend()
				if e != nil {
					logrus.Errorf("upstream.CloseSend() failed with %v, type %T", e, e)
					return
				}
				return
			}
		}
	}()
	wg.Wait()
}
