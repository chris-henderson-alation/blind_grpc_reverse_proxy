package grpcinverter

import (
	"context"
	"math/rand"
	"strings"
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
	logrus.SetLevel(logrus.ErrorLevel)
}

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
	testWithChunkSize(b, megabyte*4-kilobyte)
}

func testWithChunkSize(b *testing.B, chunkSize int) {
	b.ReportAllocs()
	b.SetBytes(gigabyte)
	forward := NewForwardProxyFacade()
	agent := NewAgent(atomic.AddUint64(&agentId, 1), host, forward.external)
	go agent.EventLoop()
	defer agent.Stop()
	defer forward.Stop()
	connector := NewConnector(&PerfAgent{until: gigabyte})
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

func randomStringLength(l int) string {
	length := rand.Intn(l)
	s := strings.Builder{}
	for i := 0; i < length; i++ {
		s.WriteByte(byte(rand.Intn(42) + 48))
	}
	return s.String()
}
