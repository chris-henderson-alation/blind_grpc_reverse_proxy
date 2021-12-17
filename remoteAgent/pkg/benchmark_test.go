package remoteAgent

import (
	"context"
	"math/rand"
	"testing"

	"github.com/Alation/alation_connector_manager/remoteAgent/logging"
	"github.com/Alation/alation_connector_manager/remoteAgent/shared"
)

const kilobyte = 1024
const megabyte = kilobyte * 1024
const gigabyte = megabyte * 1024

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
	logging.DisableLogging()
	b.ReportAllocs()
	b.SetBytes(gigabyte)
	stack := newStack(&PerfAgent{until: gigabyte})
	defer stack.Stop()
	buf := make([]byte, chunkSize)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkPerfStream(stack.alation, stack.agent.Id, stack.connector.id, buf, b)
	}
}

func benchmarkPerfStream(alation TestClient, agentId, connectorId uint64, body []byte, b *testing.B) {
	perf, err := alation.PerformanceBytes(shared.NewHeaderBuilder().
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
