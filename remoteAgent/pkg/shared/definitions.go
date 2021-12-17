package shared

import (
	"time"

	"github.com/cenkalti/backoff/v4"
	"google.golang.org/grpc"
	grpcBackoff "google.golang.org/grpc/backoff"

	"google.golang.org/grpc/keepalive"
)

var KEEPALIVE_SERVER_PARAMETERS = keepalive.ServerParameters{
	Time:    time.Second * 30,
	Timeout: time.Second * 30,
}
var KEEPALIVE_ENFORCEMENT_POLICY = keepalive.EnforcementPolicy{
	MinTime:             time.Second * 10,
	PermitWithoutStream: true,
}

// https://www.evanjones.ca/tcp-connection-timeouts.html
var KEEPALIVE_CLIENT_PARAMETERS = keepalive.ClientParameters{
	Time:                time.Second * 30,
	Timeout:             time.Second * 30,
	PermitWithoutStream: false,
}

// It might be somewhat attractive to use the same exponential backoff
// settings that gRPC by default uses to manage it's tcp connection.
var RECONNECT_EXP_BACKOFF_CONFIG = &backoff.ExponentialBackOff{
	InitialInterval:     grpcBackoff.DefaultConfig.BaseDelay,
	RandomizationFactor: grpcBackoff.DefaultConfig.Jitter,
	Multiplier:          grpcBackoff.DefaultConfig.Jitter,
	MaxInterval:         grpcBackoff.DefaultConfig.MaxDelay,
	Clock:               backoff.SystemClock,
}

var BIDIRECTIONAL_STREAM_DESC = &grpc.StreamDesc{
	ServerStreams: true,
	ClientStreams: true,
}

const tcp = "tcp"
