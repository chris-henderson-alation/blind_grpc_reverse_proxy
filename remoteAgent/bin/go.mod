module github.com/Alation/alation_connector_manager/docker/remoteAgent/forwardProxy

go 1.17

require (
	github.com/Alation/alation_connector_manager/docker/remoteAgent/grpcinverter v0.0.0
	google.golang.org/grpc v1.42.0
)

replace github.com/Alation/alation_connector_manager/docker/remoteAgent/grpcinverter => ../pkg
