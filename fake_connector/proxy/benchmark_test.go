package main

import (
	"context"
	"crypto/rand"
	"io"
	"sync"
	"testing"

	"github.com/chris-henderson-alation/blind_grpc_reverse_proxy/rpc"
	"google.golang.org/grpc/metadata"
)

// If we test exactly 1K blocks, we would generate exact multiples of

// the cipher's block size, and the cipher stream fragments would

// always be wordsize aligned, whereas non-aligned is a more typical

// use-case.

const almost1K = 1024 - 5

const almost8K = 8*1024 - 5

const kilobyte = 1024
const megabyte = kilobyte * 1024
const gigabyte = megabyte * 1024

func BenchmarkKilobyte(b *testing.B) {
	testWithChunkSize(b, kilobyte)
}

func Benchmark15Kilobyte(b *testing.B) {
	testWithChunkSize(b, kilobyte*15)
}

func Benchmark50Kilobyte(b *testing.B) {
	testWithChunkSize(b, kilobyte*50)
}

func BenchmarkMegabyte(b *testing.B) {
	testWithChunkSize(b, megabyte)
}

func BenchmarkFourMegabyte(b *testing.B) {
	testWithChunkSize(b, megabyte*4)
}

func testWithChunkSize(b *testing.B, chunkSize int) {
	data := make([]byte, chunkSize)
	b.ReportAllocs()
	n, err := rand.Read(data)
	if err != nil {
		b.Fatal(err)
	}
	if n != chunkSize {
		b.Fatalf("nedd a 1k buffer, got %d", n)
	}
	b.SetBytes(gigabyte)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		agent := &PerfAgent{
			buf:   data,
			until: gigabyte,
		}
		server := &InternetFacing{
			inbound: make(chan GrpcCall),
			outbound: map[CallId]GrpcCall{
				0: {
					Stream: &FakeAlation{},
					Error:  make(chan error),
				},
			},
			lock: sync.Mutex{},
		}
		err = server.JobCallback(agent)
		if err != nil {
			b.Fatal(err)
		}
	}
}

type FakeAlation struct{}

func (f FakeAlation) SendMsg(m interface{}) error {
	return nil
}

func (f FakeAlation) RecvMsg(m interface{}) error {
	return io.EOF
}

func (f FakeAlation) SetHeader(md metadata.MD) error {
	//TODO implement me
	panic("implement me")
}

func (f FakeAlation) SendHeader(md metadata.MD) error {
	//TODO implement me
	panic("implement me")
}

func (f FakeAlation) SetTrailer(md metadata.MD) {
	//TODO implement me
	panic("implement me")
}

func (f FakeAlation) Context() context.Context {
	//TODO implement me
	panic("implement me")
}

type PerfAgent struct {
	buf   []byte
	until int
	total int
}

func (p *PerfAgent) Send(body *rpc.Body) error {
	return nil
}

func (p *PerfAgent) Recv() (*rpc.Body, error) {
	if p.total >= p.until {
		return nil, io.EOF
	}
	would := p.total + len(p.buf)
	if would > p.until {
		left := p.until - p.total
		b := make([]byte, left)
		p.total += len(b)
		return &rpc.Body{Body: b}, nil
	}
	p.total = would
	return &rpc.Body{Body: p.buf}, nil
}

func (p *PerfAgent) Context() context.Context {
	return metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
		"x-alation-job-id": "0",
	}))
}

func (p *PerfAgent) SetHeader(md metadata.MD) error {
	//TODO implement me
	panic("implement me")
}

func (p *PerfAgent) SendHeader(md metadata.MD) error {
	//TODO implement me
	panic("implement me")
}

func (p *PerfAgent) SetTrailer(md metadata.MD) {
	//TODO implement me
	panic("implement me")
}

func (p *PerfAgent) SendMsg(m interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (p *PerfAgent) RecvMsg(m interface{}) error {
	//TODO implement me
	panic("implement me")
}
