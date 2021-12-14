package ioc // import "github.com/Alation/alation_connector_manager/docker/remoteAgent/grpcinverter"

import (
	"io"
	"sync"

	"google.golang.org/grpc/codes"

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
		var msg *Message
		for {
			err := alation.RecvMsg(&body)
			switch err {
			case nil:
				// Just another message to send.
				msg = &Message{Body: body}
			case io.EOF:
				// Alation is signaling that it is done sending
				// messages down the client side of its pipe.
				msg = &Message{EOF: true}
			default:
				// Alation croaked!
				e2 := agent.Send(&Message{Error: ErrorFromGoError(err)})
				if e2 != nil {
					// and so did the agent, lol everything is on fire.
					// I guess this can happen if the network is down for THIS
					// node and Alation is on a different node entirely.
					//logrus.Errorf("1")
					logrus.Error(e2)
				}
				//logrus.Errorf("2")
				logrus.Errorf("other error from alation.RecvMsg(), %v", err)
				return
			}
			err = agent.Send(msg)
			if err != nil {
				// Uhhhh...uh oh. This is likely an internet connection failure. here is where we have
				// the opportunity to reconnect.
				//
				// So the interesting thing here is that we have a message loaded into the barrel to send
				// downstream. So let's hold onto that just in case the agent manages to reconnect.
				logrus.Errorf("agent.Send(EOF) failed with %v, type %T", err, err)
			}
			if msg.EOF {
				return
			}
		}
	}()
	go func() {
		defer wg.Done()
		for {
			//s := time.Now()
			//dead, ok := agent.Context().Deadline()
			//logrus.Errorf("%v %v", dead, ok)
			body, err := agent.Recv()
			switch err {
			case nil:
				break
			case io.EOF:
				// The stream has successfully completed.
				alation.SendError(nil)
				return
			default:
				// Uhhhh...uh oh. This is likely an internet connection failure. Here is where we have
				// the opportunity to reconnect.
				//logrus.Errorf("3 %f seconds", time.Now().Sub(s).Seconds())
				logrus.Error(err)
				alation.SendError(status.New(codes.Unavailable, "it aint up"))
				return
			}
			if body.Error != nil {
				// This is actually where a connector to send ITS error messsage back upstream.
				// That is, this is where things like "wrong password" and such passthrough.
				//
				// Essentially this is a sad day for the connector, but this is still within the
				// happy path for this little piece of routing code.
				alation.SendError(body.Error.ToStatus())
				return
			}
			//logrus.Errorf("received %s", string(body.Body))
			e := alation.SendMsg(body.Body)
			if e != nil {
				// Comms to upstream Alation has failed in some way. We have no way to reconnect
				// the job at this point as Alation as already hung up, so let's shut things down.
				logrus.Error(err)
				return
			}
		}
	}()
	wg.Wait()
	logrus.Errorf("out")
	return
}

func ReverseProxy(upstream Forward, connector Connector) {

	wg := sync.WaitGroup{}
	wg.Add(2)
	// alation recv, connector send
	//
	// This is from the perrspective of a client (alation) sending to a server (the connector)
	go func() {
		defer wg.Done()
		for {
			body, err := upstream.Recv()
			if err != nil {
				// @TODO forward proxy is down
				return
			}
			if body.EOF {
				// This actually can't fail.
				//
				// https://github.com/grpc/grpc-go/blob/ac4edd2a03b9124d2ceda2a7c205396b31200351/stream.go#L857
				_ = connector.CloseSend()
				return
			}
			err = connector.SendMsg(body.Body)
			if err != nil {
				// This can only be, for example EOF. It cannot contain actual error messages
				logrus.Error(err)
				return
			}
		}
	}()
	// connector recv, alation send
	//
	// This is from the perspective of a server (the connector) sending to a client (alation)
	go func() {
		defer wg.Done()
		var msg []byte
		for {
			err := connector.RecvMsg(&msg)
			switch err {
			case nil:
				e := upstream.Send(&Message{Body: msg})
				if e != nil {
					// @TODO forward proxy is down
					logrus.Errorf("upstream.Send(msg) failed with %v, type %T", e, e)
					return
				}
				//logrus.Errorf("Sent %v", string(msg))
			case io.EOF:
				// This actually can't fail.
				//
				// https://github.com/grpc/grpc-go/blob/ac4edd2a03b9124d2ceda2a7c205396b31200351/stream.go#L857
				_ = upstream.CloseSend()
				return
			default:
				// We are technically the client in this scenario and clients
				// cannot send errors in the world of gRPC. So we have to represent
				// it ourselves.
				e := upstream.Send(&Message{Error: ErrorFromGoError(err)})
				if e != nil {
					// @TODO forward proxy is down
					logrus.Errorf("upstream.Send(resp) failed with %v, type %T", e, e)
					return
				}
				// This actually can't fail.
				//
				// https://github.com/grpc/grpc-go/blob/ac4edd2a03b9124d2ceda2a7c205396b31200351/stream.go#L857
				_ = upstream.CloseSend()
				return
			}
		}
	}()
	wg.Wait()
	logrus.Error("good night")
}
