package grpcinverter

import (
	"context"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
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

func init() {
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.TraceLevel)
}

func BenchmarkKilobyte(b *testing.B) {
	testWithChunkSize(b, kilobyte)
}

func Benchmark2Kilobyte(b *testing.B) {
	testWithChunkSize(b, kilobyte*2)
}

func Benchmark4Kilobyte(b *testing.B) {
	testWithChunkSize(b, kilobyte*4)
}

func Benchmark8Kilobyte(b *testing.B) {
	testWithChunkSize(b, kilobyte*8)
}

func Benchmark16Kilobyte(b *testing.B) {
	testWithChunkSize(b, kilobyte*16)
}

func Benchmark32Kilobyte(b *testing.B) {
	testWithChunkSize(b, kilobyte*32)
}

func Benchmark64Kilobyte(b *testing.B) {
	testWithChunkSize(b, kilobyte*64)
}

func BenchmarkMegabyte(b *testing.B) {
	testWithChunkSize(b, megabyte)
}

func BenchmarkFourMegabyte(b *testing.B) {
	testWithChunkSize(b, megabyte*4-kilobyte)
}

func testWithChunkSize(b *testing.B, chunkSize int) {
	b.ReportAllocs()
	b.SetBytes(gigabyte)
	forward := NewForwardProxyFacade("0.0.0.0", "0.0.0.0")
	agent := NewAgent(atomic.AddUint64(&agentId, 1), host, forward.external)
	go agent.EventLoop()
	defer agent.Stop()
	defer forward.Stop()
	connector := NewConnector(&PerfAgent{until: gigabyte}, "0.0.0.0")
	connector.Start()
	defer connector.Stop()
	for {
		if forward.agents.Listening(agentId) {
			break
		}
		time.Sleep(time.Millisecond * 200)
	}
	alation := NewAlationClient(forward.internal)
	buf := make([]byte, chunkSize)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkPerfStream(alation, agent.id, connector.id, buf, b)
	}
}

func benchmarkPerfStream(alation TestClient, agentId, connectorId uint64, body []byte, b *testing.B) {
	perf, err := alation.PerformanceBytes(NewHeaderBuilder().
		SetJobId(1).
		SetAgentId(agentId).
		SetConnectorId(connectorId).
		Build(context.Background()),
		&TestBytes{Body: body})
	if err != nil {
		b.Fatal(err)
	}
	for _, err := perf.Recv(); err == nil; _, err = perf.Recv() {
	}
}

type PerfAgent struct {
	until int
	total int
	UnimplementedTestServer
}

func (p *PerfAgent) PerformanceBytes(s *TestBytes, server Test_PerformanceBytesServer) error {
	size := len(s.Body)
	for {
		would := p.total + size
		if would > p.until {
			return nil
		}
		p.total = would
		err := server.Send(s)
		if err != nil {
			panic(err)
		}
	}
}
