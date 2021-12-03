package main

type NoopCodec struct{}

func (cb NoopCodec) Name() string {
	return "NoopCodec"
}

func (cb NoopCodec) Marshal(v interface{}) ([]byte, error) {
	return v.([]byte), nil
}

func (cb NoopCodec) Unmarshal(data []byte, v interface{}) error {
	b, _ := v.(*[]byte)
	*b = data
	return nil
}

//type Proxy struct {
//	outbound *InternetFacing
//	rdbms.UnimplementedRdbmsServer
//}
//
//func (p Proxy) ConfigurationVerification(ctx context.Context, request *verify.Request) (*empty.Empty, error) {
//	connector := connectorId(ctx)
//	response := p.outbound.Send(GrpcCall{
//		Message:   request,
//		Method:    "rdbms.Rdbms/ConfigurationVerification",
//		Connector: connector,
//		CallType:  Unary,
//	})
//	answer := <-response
//	if answer.Error != nil {
//		return nil, answer.Error
//	}
//	e := &empty.Empty{}
//	err := proto.Unmarshal(answer.Message, e)
//	if err != nil {
//		return nil, err
//	}
//	return e, nil
//}
//
//func (p Proxy) FilterExtraction(request *mde.FilterExtractionRequest, server rdbms.Rdbms_FilterExtractionServer) error {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (p Proxy) MetadataExtraction(request *mde.MetadataExtractionRequest, server rdbms.Rdbms_MetadataExtractionServer) error {
//	connector := connectorId(server.Context())
//	response := p.outbound.Send(GrpcCall{
//		Message:   request,
//		Method:    "rdbms.Rdbms/ConfigurationVerification",
//		Connector: connector,
//		CallType:  ServerStream,
//	})
//	for answer := range response {
//		if answer.Error != nil {
//			return answer.Error
//		}
//		md := &mde.MetadataExtractionResponse{}
//		err := proto.Unmarshal(answer.Message, md)
//		if err != nil {
//			return err
//		}
//		err = server.Send(md)
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}
