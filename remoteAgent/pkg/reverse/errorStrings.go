package reverse

import (
	"fmt"

	"github.com/Alation/alation_connector_manager/docker/remoteAgent/protocol"
	"github.com/golang/protobuf/ptypes/any"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
)

var ConnectorDown = connectorDown{}

type connectorDown struct{}

func (c connectorDown) Fmt(original error, agentId uint64, connectorId uint64) *protocol.Error {
	fmted := fmt.Sprintf(`The connector appears to be down! Please SSH into the node provisioned with your agent 
ID %d and check up on the health of the system. Some helpful commands might be "docker ps -a", "docker ps -a | grep connector_%d", 
"kratos start %d", and "hydra restart" (ALERT! restarts all connectors and agent services!)`, agentId, connectorId, connectorId)
	o, _ := status.FromError(original)
	o2, err := anypb.New(o.Proto())
	if err != nil {
		panic(err)
	}
	return &protocol.Error{
		Code:    int32(codes.Unavailable),
		Message: fmted,
		Details: []*any.Any{o2},
	}
}
