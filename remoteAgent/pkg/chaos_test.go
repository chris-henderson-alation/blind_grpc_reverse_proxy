package remoteAgent

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync/atomic"
	"testing"
	"time"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// Test Reconnect

type ReconnectServerStream struct {
	UnimplementedTestServer
}

func (r ReconnectServerStream) ServerStream(s *String, server Test_ServerStreamServer) error {
	for i := 0; i < 20; i++ {
		sent := fmt.Sprintf("%d", i)
		err := server.Send(&String{Message: sent})
		if err != nil {
			return err
		}
		fmt.Println(sent)
		time.Sleep(time.Second)
	}
	return nil
}

func TestChaosReconnect(t *testing.T) {
	if os.Geteuid() != 0 {
		t.Log("This test requires escalated privileges in order to manipulate the network")
		t.SkipNow()
	}
	link := NewLocalLink(t.Name())
	defer link.Destroy()
	stack := newCustomNetworkingStack(&ReconnectServerStream{}, nil, &link.addr, nil)
	defer stack.Stop()
	stream, err := stack.alation.ServerStream(stack.Headers(context.Background()), &String{Message: ""})
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		want := fmt.Sprintf("%d", i)
		s, err := stream.Recv()
		if err != nil {
			t.Fatal(err)
		}
		if s.Message != want {
			t.Fatalf("wanted '%s' got '%s'", want, s.Message)
		}
	}
	link.Down()
	t.Log("sleeping!")
	time.Sleep(time.Second * 30)
	t.Log("finihsings")
	link.Up()
	for i := 10; i < 20; i++ {
		want := fmt.Sprintf("%d", i)
		s, err := stream.Recv()
		if err != nil {
			t.Fatal(err)
		}
		if s.Message != want {
			t.Fatalf("wanted '%s' got '%s'", want, s.Message)
		}
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// Test Reconnect a Job Tunnel

type BrokenTunnelTestConnector struct {
	UnimplementedTestServer
}

func (r *BrokenTunnelTestConnector) ServerStream(s *String, server Test_ServerStreamServer) error {
	for i := 0; i < 20; i++ {
		err := server.Send(&String{Message: fmt.Sprintf("%d", i)})
		if err != nil {
			return err
		}
	}
	return nil
}

func TestChaosBrokenTunnel(t *testing.T) {
	if os.Geteuid() != 0 {
		t.Log("This test requires escalated privileges in order to manipulate the network")
		t.SkipNow()
	}
	t.Log("what is all this then come now")
	link := NewLocalLink(t.Name())
	defer link.Destroy()
	stack := newCustomNetworkingStack(&BrokenTunnelTestConnector{}, nil, &link.addr, nil)
	defer stack.Stop()
	link.Down()
	time.Sleep(time.Second * 60)
	link.Up()
	iters := 0
	for {
		if iters > 100 {
			t.Fatal("not in time")
		}
		iters += 1
		if stack.forward.Listeners.Listening(stack.agent.Id) {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}
	t.Log("yaaaas")
	stream, err := stack.alation.ServerStream(stack.Headers(context.Background()), &String{Message: ""})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("yus")
	for i := 0; i < 20; i++ {
		want := fmt.Sprintf("%d", i)
		s, err := stream.Recv()
		if err != nil {
			t.Fatal(err)
		}
		if s.Message != want {
			t.Fatalf("wanted '%s' got '%s'", want, s.Message)
		}
	}
}

// Shameless copying in order to figure out how to have mulitple dummy links in order
// to simulate down networks.
//
// https://linuxconfig.org/configuring-virtual-network-interfaces-in-linux
type LocalLink struct {
	addr     string
	addrCIDR string
	name     string
	label    string
}

var nextAddrSpace uint32 = 0

func NewLocalLink(name string) LocalLink {
	next := atomic.AddUint32(&nextAddrSpace, 1)
	name = fmt.Sprintf("ocf%d", next)
	addr := fmt.Sprintf("10.0.%d.0", next)
	addrCIDR := fmt.Sprintf("%s/24", addr)
	link := LocalLink{
		addr:     addr,
		addrCIDR: addrCIDR,
		name:     name,
		label:    fmt.Sprintf("%s:0", name),
	}
	link.Destroy()
	if err := exec.Command("modprobe", "dummy").Run(); err != nil {
		panic(err)
	}
	if err := exec.Command("ip", "link", "add", link.name, "type", "dummy").Run(); err != nil {
		panic(err)
	}
	if err := exec.Command("ip", "addr", "add", link.addrCIDR, "brd", "+", "dev", link.name, "label", link.label).Run(); err != nil {
		panic(err)
	}
	link.Up()
	return link
}

func (l LocalLink) Up() {
	if err := exec.Command("ip", "link", "set", "dev", l.name, "up").Run(); err != nil {
		panic(err)
	}
}

func (l LocalLink) Down() {
	if err := exec.Command("ip", "link", "set", "dev", l.name, "down").Run(); err != nil {
		panic(err)
	}
}

func (l LocalLink) Destroy() {
	exec.Command("ip", "addr", "del", l.addrCIDR, "brd", "+", "dev", l.name, "label", l.label).Run()
	exec.Command("ip", "link", "delete", l.name, "type", "dummy").Run()
}
