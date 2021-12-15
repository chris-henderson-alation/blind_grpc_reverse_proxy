package grpcinverter

import (
	"fmt"

	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/ptypes/any"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/Alation/alation_connector_manager/docker/remoteAgent/grpcinverter/ioc"
)

var ConnectorDown = connectorDown{}

type connectorDown struct{}

func (c connectorDown) Fmt(original error, agentId uint64, connectorId uint64) *ioc.Error {
	fmted := fmt.Sprintf(`The connector appears to be down! Please SSH into the node provisioned with your agent 
ID %d and check up on the health of the system. Some helpful commands might be "docker ps -a", "docker ps -a | grep connector_%d", 
"kratos start %d", and "hydra restart" (ALERT! restarts all connectors and agent services!)`, agentId, connectorId, connectorId)
	o, _ := status.FromError(original)
	o2, err := anypb.New(o.Proto())
	if err != nil {
		panic(err)
	}
	return &ioc.Error{
		Code:    int32(codes.Unavailable),
		Message: fmted,
		Details: []*any.Any{o2},
	}
}

//type ErrorString string
//
//func (e ErrorString) Fmt(strs ...interface{}) error {
//	return fmt.Errorf(string(e), strs...)
//}
//
//const ConnectorDown ErrorString = `The connector appears to be down! Please SSH into the node provisioned with you agent ID %d. Some helpful commands might be "docker ps", "docker ps | grep connector_%d", "docker start connector_%d"`
//
//type lol struct{}
//
//var derp lol = lol{}
