package logging

import (
	"context"
	"sync"

	"google.golang.org/grpc/peer"

	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

var LOGGER *zap.Logger
var lock = sync.Mutex{}

func init() {
	var err error
	LOGGER, err = zap.NewProduction(zap.AddCaller(), zap.AddStacktrace(zapcore.FatalLevel))
	if err != nil {
		panic(err)
	}
}

// DisableLogging is useful for making tests a bit quieter if you want to. This is pretty much
// mandatory for benchmark test, though, otherwise you will never see the output.
//
// This function is threadsafe with other class to DisableLogging but it is NOT thread safe
// when using the logger itself!
func DisableLogging() {
	lock.Lock()
	defer lock.Unlock()
	LOGGER = zap.New(zapcore.NewNopCore())
}

func Address(addr string) zap.Field {
	return zap.String("address", addr)
}

func Port(port uint16) zap.Field {
	return zap.Uint16("port", port)
}

func Protocol(protocol string) zap.Field {
	return zap.String("protocol", protocol)
}

func Method(method string) zap.Field {
	return zap.String("method", method)
}

func PeerAlation(ctx context.Context) zap.Field {
	return zap.String("peerAlation", PeerAddr(ctx))
}

func PeerProxy(ctx context.Context) zap.Field {
	return zap.String("peerProxy", PeerAddr(ctx))
}

func PeerAgent(ctx context.Context) zap.Field {
	return zap.String("peerAgent", PeerAddr(ctx))
}

func PeerConnector(ctx context.Context) zap.Field {
	return zap.String("peerConnector", PeerAddr(ctx))
}

func Connector(id uint64) zap.Field {
	return zap.Uint64("connector", id)
}

func Agent(id uint64) zap.Field {
	return zap.Uint64("agent", id)
}

func Job(id uint64) zap.Field {
	return zap.Uint64("job", id)
}

func Error(err error) zap.Field {
	return zap.Error(err)
}

func Errors(errs ...error) zap.Field {
	return zap.Errors("allErrors", errs)
}

func Wanted(value string) zap.Field {
	return zap.String("wanted", value)
}

func PeerAddr(ctx context.Context) string {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return "<unknown>"
	}
	return p.Addr.String()
}
