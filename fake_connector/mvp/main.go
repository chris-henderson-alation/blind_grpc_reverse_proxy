package main // import "asdasd"
import (
	context "context"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func main() {
	conn, err := grpc.Dial(fmt.Sprintf("0.0.0.0:%d", 10001), grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                time.Second * 30,
		PermitWithoutStream: true,
	}), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	method := "/slow.Slow/Slow"
	downstream, err := conn.NewStream(context.Background(), &grpc.StreamDesc{
		ServerStreams: true,
		ClientStreams: true,
	}, method, grpc.ForceCodec(NoopCodec{}))
	b, err := proto.Marshal(&String{Message: "lolol"})
	if err != nil {
		panic(err)
	}
	err = downstream.SendMsg(b)
	if err != nil {
		panic(err)
	}
	var msg []byte
	err = downstream.RecvMsg(&msg)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(msg))
}

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
