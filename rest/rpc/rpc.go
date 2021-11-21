package rpc // import "github.com/chris-henderson-alation/blind_grpc_reverse_proxy/rpc"

import "encoding/base64"

type Job struct {
	Headers Headers `json:"headers"`
	Body    Body    `json:"body"`
}

type Headers struct {
	Connector uint64      `json:"connector"`
	Method    string      `json:"method"`
	Type      GrpCallType `json:"type"`
	JobID     uint64      `json:"job_id"`
}

type Body []byte

type GrpCallType uint

const (
	Unary GrpCallType = iota
	UnaryStream
	StreamUnary
	StreamStream
)

func (b *Body) UnmarshalJSON(src []byte) error {
	src = src[1 : len(src)-1]
	*b = make([]byte, base64.StdEncoding.DecodedLen(len(src)))
	written, err := base64.StdEncoding.Decode(*b, src)
	if err != nil {
		return err
	}
	*b = (*b)[0:written]
	return nil
}

func (b *Body) MarshalJSON() ([]byte, error) {
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(*b))+2)
	dst[0] = '"'
	dst[len(dst)-1] = '"'
	base64.StdEncoding.Encode(dst[1:len(dst)-1], *b)
	return dst, nil
}
