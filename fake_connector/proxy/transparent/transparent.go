package transparent

//import (
//	"fmt"
//
//	"google.golang.org/grpc"
//	"google.golang.org/grpc/codes"
//)

//func Main() {
//	grpc.NewServer(grpc.UnknownServiceHandler(DoIt), grpc.ForceServerCodec(NoopCodec{}))
//
//}
//
//func DoIt(_ interface{}, stream grpc.ServerStream) error {
//	fullMethodName, ok := grpc.MethodFromServerStream(stream)
//	if !ok {
//		return grpc.Errorf(codes.Internal, "lowLevelServerStream not exists in context")
//	}
//	fmt.Println(fullMethodName)
//	var msg []byte
//	err := stream.RecvMsg(msg)
//	if err != nil {
//		return err
//	}
//	fmt.Println(string(msg))
//	return nil
//}
//
//func RegisterService(server *grpc.Server, serviceName string, methodNames ...string) {
//	grpc.UnknownServiceHandler()
//	streamer := &handler{director}
//	fakeDesc := &grpc.ServiceDesc{
//		ServiceName: serviceName,
//		HandlerType: (*interface{})(nil),
//	}
//	for _, m := range methodNames {
//		streamDesc := grpc.StreamDesc{
//			StreamName:    m,
//			Handler:       streamer.handler,
//			ServerStreams: true,
//			ClientStreams: true,
//		}
//		fakeDesc.Streams = append(fakeDesc.Streams, streamDesc)
//	}
//	grpc.UnknownServiceHandler()
//	server.RegisterService(fakeDesc, streamer)
//}
//
//func TransparentHandler(director StreamDirector) grpc.StreamHandler {
//	streamer := &handler{director}
//	return streamer.handler
//}
//
//type handler struct {
//	director StreamDirector
//}
//
//// handler is where the real magic of proxying happens.
//// It is invoked like any gRPC server stream and uses the gRPC server framing to get and receive bytes from the wire,
//// forwarding it to a ClientStream established against the relevant ClientConn.
//func (s *handler) handler(srv interface{}, serverStream grpc.ServerStream) error {
//	// little bit of gRPC internals never hurt anyone
//	fullMethodName, ok := grpc.MethodFromServerStream(serverStream)
//	if !ok {
//		return grpc.Errorf(codes.Internal, "lowLevelServerStream not exists in context")
//	}
//	grpc.Me
//}
//
//type NoopCodec struct{}
//
//func (cb NoopCodec) Name() string {
//	return "NoopCodec"
//}
//
//func (cb NoopCodec) Marshal(v interface{}) ([]byte, error) {
//	return v.([]byte), nil
//}
//
//func (cb NoopCodec) Unmarshal(data []byte, v interface{}) error {
//	b, _ := v.(*[]byte)
//	*b = data
//	return nil
//}
