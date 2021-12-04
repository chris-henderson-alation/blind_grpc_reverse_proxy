package main // import "fake_connector/agent"
import (
	"context"
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"

	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/metadata"

	"github.com/chris-henderson-alation/blind_grpc_reverse_proxy/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func main() {
	client := NewClient("1234")
	jobs, err := client.InitiateJob(context.Background(), &empty.Empty{})
	if err != nil {
		panic(err)
	}
	for {
		job := &rpc.Job{}
		err = jobs.RecvMsg(job)
		if err != nil {
			panic(err)
		}
		go client.dispatch(job)
	}
}

type Client struct {
	rpc.AgentClient
}

func NewClient(port string) *Client {
	conn, err := grpc.Dial(fmt.Sprintf("0.0.0.0:%s", port),
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                time.Second * 30,
			PermitWithoutStream: true,
		}))
	if err != nil {
		panic(err)
	}
	return &Client{rpc.NewAgentClient(conn)}
}

func (c *Client) dispatch(job *rpc.Job) {
	headers := metadata.New(map[string]string{
		"x-alation-job-id": strconv.FormatUint(job.JobID, 10),
	})
	upstream, err := c.JobCallback(metadata.NewOutgoingContext(context.Background(), headers))
	if err != nil {
		panic(err)
	}
	conn, err := grpc.Dial(fmt.Sprintf("0.0.0.0:%d", job.Connector+10000), grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                time.Second * 30,
		PermitWithoutStream: true,
	}), grpc.WithInsecure())
	if err != nil {
		fmt.Println("booo tits")
		upstream.SendMsg(err)
		return
	}
	downstream, err := conn.NewStream(context.Background(), &grpc.StreamDesc{
		ServerStreams: true,
		ClientStreams: true,
	}, job.Method, grpc.ForceCodec(NoopCodec{}))
	if err != nil {
		fmt.Println("booo")
		upstream.SendMsg(err)
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for {
			//fmt.Println("I go")
			body, err := upstream.Recv()
			//fmt.Println("got message from upstream")
			switch err {
			case nil:
				if body.EOF {
					downstream.CloseSend()
					return
				}
				//fmt.Println("and it was nil! Yay!")
				err = downstream.SendMsg(body.Body)
				switch err {
				case nil:
				case io.EOF:
					upstream.CloseSend()
					return
				}
			case io.EOF:
				//fmt.Println("and it was EOF")
				downstream.CloseSend()
				return
			default:
				fmt.Printf("and it was something else %v\n", err)
				downstream.SendMsg(err)
				downstream.CloseSend()
				return
			}
		}
	}()
	go func() {
		defer wg.Done()
		var msg []byte
		for {
			err := downstream.RecvMsg(&msg)
			switch err {
			case nil:
				upstream.Send(&rpc.Body{Body: msg})
			case io.EOF:
				upstream.CloseSend()
				return
			default:
				// We are technically the client in this scenario and clients
				// cannot send errors in the world of gRPC. So we have to represent
				// it ourselves.
				serr, _ := status.FromError(err)
				resp := &rpc.Body{
					Error: &rpc.Error{
						Code:        uint32(serr.Code()),
						Description: serr.Message(),
					},
				}
				upstream.Send(resp)
				upstream.CloseSend()
				return
			}
		}
	}()
	wg.Wait()
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
