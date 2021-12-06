package grpcinverter

import (
	"time"

	"google.golang.org/grpc/keepalive"
)

var KEEPALIVE_SERVER_PARAMETERS = keepalive.ServerParameters{
	Time: time.Second * 30,
}
var KEEPALIVE_ENFORCEMENT_POLICY = keepalive.EnforcementPolicy{
	MinTime:             time.Second * 10,
	PermitWithoutStream: true,
}

var KEEPALIVE_CLIENT_PARAMETERS = keepalive.ClientParameters{
	Time:                time.Second * 30,
	PermitWithoutStream: true,
}
