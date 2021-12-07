package ioc // import "github.com/Alation/alation_connector_manager/docker/remoteAgent/grpcinverter"

import (
	"io"
	"sync"

	"github.com/sirupsen/logrus"

	codes "google.golang.org/grpc/codes"
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
	replay := make(chan *Message, 1)
	alationFailed := make(chan struct{})
	agentFailed := make(chan struct{})
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		var body []byte
		for {
			err := alation.RecvMsg(&body)
			switch err {
			case nil:
				// Just another message to send.
				err = agent.Send(&Message{Body: body})
				if err != nil {
					// Uhhhh...uh oh. This is likely an internet connection failure. here is where we have
					// the opportunity to reconnect.
					//
					// So the interesting thing here is that we have a message loaded into the barrel to send
					// downstream. So let's hold onto that just in case the agent manages to reconnect.
					replay <- &Message{Body: body}
					close(agentFailed)
					logrus.Errorf("agent.Send(EOF) failed with %v, type %T", err, err)
					return
				}
			case io.EOF:
				// Alation is signaling that it is done sending
				// messages down the client side of its pipe.
				err = agent.Send(&Message{EOF: true})
				if err != nil {
					// Uhhhh...uh oh. This is likely an internet connection failure. here is where we have
					// the opportunity to reconnect.
					//
					// So the interesting thing here is that we have a message loaded into the barrel to send
					// downstream. So let's hold onto that just in case the agent manages to reconnect.
					replay <- &Message{EOF: true}
					close(agentFailed)
					logrus.Errorf("agent.Send(EOF) failed with %v, type %T", err, err)
					return
				}
				return
			default:
				// Alation croaked!
				close(alationFailed)
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
				break
			case io.EOF:
				// The stream has successfully completed.
				return
			default:
				// Uhhhh...uh oh. This is likely an internet connection failure. here is where we have
				// the opportunity to reconnect.
				logrus.Error(err)
				close(agentFailed)
				return
			}
			if body.Error != nil {
				// This is actually where a connector to send ITS error messsage back upstream.
				// That is, this is where things like "wrong password" and such passthrough.
				//
				// Essentially this is a sad day for the connector, but this is still within the
				// happy path for this little piece of routing code.
				alation.SendError(status.New(codes.Code(body.Error.Code), body.Error.Description))
				return
			}
			e := alation.SendMsg(body.Body)
			if e != nil {
				// Comms to upstream Alation has failed in some way. We have no way to reconnect
				// the job at this point as Alation as already hung up, so let's shut things down.
				logrus.Error(err)
				close(alationFailed)
			}
		}
	}()
	wg.Wait()
}

//func ForwardProxy1(alation Alation, agent Reverse) {
//	upstream := make(chan *Message)
//	outbound := make(chan *Message)
//	inbound := make(chan *Message)
//	go func() {
//		// covers upstream
//		var body []byte
//		for {
//			err := alation.RecvMsg(&body)
//			switch err {
//			case nil:
//				outbound <- &Message{Body: body}
//			case io.EOF:
//				outbound <- &Message{EOF: true}
//				return
//			default:
//				// this is an actual bad error from alation so send a shutdown signal
//				// uhhhh but we need to get the connector to just stop as well.
//				//
//				// this EOF would just shutdown the client stream
//				// @TODO
//				logrus.Error(err)
//				outbound <- &Message{EOF: true}
//				return
//			}
//		}
//	}()
//	go func() {
//		// covers inbound
//		for {
//			message, err := agent.Recv()
//			switch err {
//			case nil:
//				inbound <- message
//			case io.EOF:
//
//			}
//			if err != nil {
//				// @TODO
//				return
//			}
//			inbound <- message
//		}
//	}()
//	go func() {
//		for {
//			select {
//			case <-alation.Context().Done():
//				outbound <- &Message{EOF: true}
//			case <-agent.Context().Done():
//				// @TODO
//			case body := <-upstream:
//				err := agent.Send(&Message{Body: body})
//				if err != nil {
//					// @TODO
//					return
//				}
//			}
//		}
//	}()
//	go func() {
//		for {
//			select {
//			case <-alation.Context().Done():
//				outbound <- &Message{EOF: true}
//			case <-agent.Context().Done():
//				// @TODO
//			case message := <-inbound:
//				if message.Error != nil {
//					alation.SendError(status.New(codes.Code(message.Error.Code), message.Error.Description))
//					return
//				}
//				// could panic if EOF instead of body
//				alation.SendMsg(message.Body)
//			}
//			// @TODO
//		}
//	}()
//}

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
