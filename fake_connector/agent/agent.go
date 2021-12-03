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

	"google.golang.org/grpc/keepalive"

	"github.com/chris-henderson-alation/blind_grpc_reverse_proxy/rpc"
	"google.golang.org/grpc"
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
			fmt.Println("I go")
			body, err := upstream.Recv()
			fmt.Println("got message from upstream")
			switch err {
			case nil:
				fmt.Println("and it was nil! Yay!")
				downstream.SendMsg(body.Body)
			case io.EOF:
				fmt.Println("and it was EOF")
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
		for {
			fmt.Println("wat")
			var msg []byte
			err := downstream.RecvMsg(&msg)
			fmt.Println("wooo tits")
			switch err {
			case nil:
				upstream.Send(&rpc.Bodies{
					Bodies: [][]byte{msg},
				})
			case io.EOF:
				upstream.CloseSend()
				return
			default:
				serr, _ := status.FromError(err)
				resp := &rpc.Bodies{
					Error: &rpc.Error{
						Code:        uint32(serr.Code()),
						Description: serr.Message(),
					},
				}
				upstream.Send(resp)
				fmt.Printf("I think this is happening first %v\n", err)
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
