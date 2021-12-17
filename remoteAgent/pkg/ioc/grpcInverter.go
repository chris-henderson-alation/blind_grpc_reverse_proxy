package ioc // import "github.com/Alation/alation_connector_manager/docker/remoteAgent/grpcinverter"

import (
	"io"
	"sync"

	"github.com/Alation/alation_connector_manager/docker/remoteAgent/grpcinverter/logging"
	"go.uber.org/zap"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ServerStreamWithError serves as a way for functions deeper in the callstack (namely the ForwardProxy function) to return
// both values and a single optional error back upstream to Alation.
//
// This is because the ONLY opportunity to actually send an error to Alation is from the `error` return type of
// ForwardProxy.GenericStreamHandler - you cannot send an error using just the provided grpc.ServerStream!
type ServerStreamWithError struct {
	grpc.ServerStream
	Error chan error
}

func (a *ServerStreamWithError) SendError(s *status.Status) {
	a.Error <- s.Err()
}

func ForwardProxy(alation ServerStreamWithError, agent GrpcInverter_PipeServer, logger *zap.Logger) {
	var bytesIn uint64 = 0
	var bytesOut uint64 = 0
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		defer logger.Info("AlationRecv coroutine shutting down")
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
				logger.Info("Alation sent EOF for its client stream")
				msg = &Message{EOF: true}
			default:
				// Alation croaked!
				logger.Error("Alation's client stream failed", logging.Error(err))
				e2 := agent.Send(&Message{Error: ErrorFromGoError(err)})
				if e2 != nil {
					logger.Error("Failed to send Alation's client stream failure to the agent", logging.Error(e2))
				}
				return
			}
			err = agent.Send(msg)
			if err != nil {
				// @TODO
				return
			}
			bytesIn += uint64(len(body))
			if msg.EOF {
				return
			}
		}
	}()
	go func() {
		defer wg.Done()
		defer logger.Info("AgentRecv coroutine shutting down")
		for {
			body, err := agent.Recv()
			switch err {
			case nil:
				break
			case io.EOF:
				// The stream has successfully completed.
				alation.SendError(nil)
				logger.Info("Agent cleanly shutdown its send channel.")
				return
			default:
				alation.SendError(status.New(codes.Unavailable, err.Error()))
				logger.Error("Unexpected disconnect from agent", logging.Error(err))
				return
			}
			if body.Error != nil {
				// This is not an error in the sense that any networking or such has
				// failed. This is the raw error as is being returned by the connector.
				//
				// So this is things like "bad password" and "could not reach database"
				alation.SendError(body.Error.ToStatus())
				logger.Info("Connector failed its request", logging.Error(ErrorToGoError(body.Error)))
				return
			}
			e := alation.SendMsg(body.Body)
			if e != nil {
				// Comms to upstream Alation has failed in some way. We have no way to reconnect
				// the job at this point as Alation as already hung up, so let's shut things down.
				//logrus.Error(err)
				logger.Error("Unexpected disconnect from Alation", logging.Error(err))
				return
			}
			bytesOut += uint64(len(body.Body))
		}
	}()
	wg.Wait()
	logger.Info("Stream complete",
		zap.Uint64("bytesFromAlation", bytesIn),
		zap.Uint64("bytesFromConnector", bytesOut))
	return
}

func ReverseProxy(upstream GrpcInverter_PipeClient, connector grpc.ClientStream, logger *zap.Logger) {
	var bytesIn uint64 = 0
	var bytesOut uint64 = 0
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		defer logger.Info("ProxyRecv coroutine shutting down")
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
			bytesIn += uint64(len(body.Body))
			err = connector.SendMsg(body.Body)
			if err != nil {
				// This can only be, for example EOF. It cannot contain actual error messages
				//logrus.Error(err)
				return
			}
		}
	}()
	go func() {
		defer wg.Done()
		defer logger.Info("ConnectorRecv coroutine shutting down")
		var msg []byte
		for {
			err := connector.RecvMsg(&msg)
			switch err {
			case nil:
				bytesOut += uint64(len(msg))
				e := upstream.Send(&Message{Body: msg})
				if e != nil {
					logger.Error("Failed to send connector bytes due to a proxy disconnect", logging.Error(e))
					return
				}
			case io.EOF:
				// This actually can't fail.
				//
				// https://github.com/grpc/grpc-go/blob/ac4edd2a03b9124d2ceda2a7c205396b31200351/stream.go#L857
				_ = upstream.CloseSend()
				logger.Info("Connector cleanly shutdown its send channel.")
				return
			default:
				// We are technically the client in this scenario and clients
				// cannot send errors in the world of gRPC. So we have to represent
				// it ourselves.
				logger.Info("Connector failed its request")
				e := upstream.Send(&Message{Error: ErrorFromGoError(err)})
				if e != nil {
					logger.Error("Failed to send connector error back upstream", logging.Error(e))
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
	logger.Info("Stream complete",
		zap.Uint64("bytesFromAlation", bytesIn),
		zap.Uint64("bytesFromConnector", bytesOut))
}
