package grpcinverter

// NoopCodec is used via the grpc.ForceCodec call option and grpc.ForceServerCodec server option.
//
// This codec does exactly zero work to marshal (serialize) and unmarshal (deserialize data).
//
// That is to say, this is the codec that we must use when are blindly shuffling bytes between
// Alation and the target connector.
//
// Please Marshal and Unmarshal for directions on how to use this codec and type safety warnings.
type NoopCodec struct{}

// Name fulfills the code interface.
func (cb NoopCodec) Name() string {
	return "NoopCodec"
}

// Marshal takes in an empty interface that MUST be a []byte containing the raw wire binary
// of an outgoing message and simply returns the same []byte.
//
// This is useful when you want to SEND raw data that is already in the gRPC wire protocol format
// (for example, the forward proxy wants to send a message from Alation to the connector, or the reverse
// proxy wants to send a message from the connector to Alation.)
//
// ALERT! The provided src MUST be []byte. This function does NO conditional checks on the provided type!
// Any type provided other than a []byte will PANIC the program on a failed type assertion!
func (cb NoopCodec) Marshal(src interface{}) ([]byte, error) {
	return src.([]byte), nil
}

// Unmarshal takes in the wire protocol binary contained within the provided data argument and assigns
// its pointer value to dst empty interface. That is to say that NO deserialization occurs within this
// codec at all - recipients shall receive the raw wire binary of the incoming message.
//
// This is useful when you want to RECEIVE raw data that is never gets deserialized from its raw gRPC wire protocol format
// (for example, the forward proxy wants to receive a message from Alation to send to the connector, or the reverse
// proxy wants to receive a message from the connector to send to Alation.)
//
// ALERT! The provided dst MUST be a *[]byte. This function does NO conditional checks on the provided dst type!
// Any type provided other than a *[]byte will PANIC the program on a failed type assertion!
func (cb NoopCodec) Unmarshal(src []byte, dst interface{}) error {
	*dst.(*[]byte) = src
	return nil
}
