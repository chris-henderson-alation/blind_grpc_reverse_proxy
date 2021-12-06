module github.com/Alation/alation_connector_manager/docker/remoteAgent/reverseProxy

go 1.17

require (
	github.com/Alation/alation_connector_manager/docker/remoteAgent/grpcinverter v0.0.0-00010101000000-000000000000
	github.com/golang/protobuf v1.5.2
	google.golang.org/appengine v1.4.0
	google.golang.org/grpc v1.42.0
)

require (
	github.com/sirupsen/logrus v1.8.1 // indirect
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4 // indirect
	golang.org/x/sys v0.0.0-20210510120138-977fb7262007 // indirect
	golang.org/x/text v0.3.3 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/protobuf v1.26.0 // indirect
)

replace github.com/Alation/alation_connector_manager/docker/remoteAgent/grpcinverter => ../../pkg
