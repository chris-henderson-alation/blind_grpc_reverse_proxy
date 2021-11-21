package main // import "github.com/chris-henderson-alation/blind_grpc_reverse_proxy/servicer"

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"

	"github.com/chris-henderson-alation/blind_grpc_reverse_proxy/rpc"

	"github.com/golang/protobuf/proto"

	"github.com/gin-gonic/gin"
)

//go:generate protoc -I . --go_out=plugins=grpc:. tennis.proto

func main() {
	r := gin.Default()
	r.GET("/connectors", func(c *gin.Context) {
		//if rand.Intn(10) < 3 {
		//	// 30% chance for there to be no jobs
		//	c.Status(http.StatusNoContent)
		//	return
		//}
		// Build a message with some random data in it
		random := randomString()
		fmt.Println("Generated ", random)
		connectorMessage, err := marshal(makeBall())
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.JSON(200, rpc.Job{
			Headers: rpc.Headers{
				// Pick a random connector
				//Connector: uint64(rand.Intn(3) + 1),
				Connector: 1,
				Method:    "main.Tennis/Rapid",
				Type:      rpc.UnaryStream,
				JobID:     jobID(),
			},
			Body: connectorMessage,
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func makeBall() *Ball {
	return &Ball{
		Ball: randomString(),
		Complex: &Complex{
			Array:  []uint32{1, 2, 3, 4, 5},
			Corpus: Complex_IMAGES,
			Result: &Complex_Result{
				Url:      "https://alation.com",
				Title:    "Winners",
				Snippets: []string{randomString(), randomString(), randomString()},
			},
		},
	}
}

var id uint64 = 0
var idLock = sync.Mutex{}

func jobID() uint64 {
	idLock.Lock()
	defer idLock.Unlock()
	id += 1
	return id
}

func marshal(v interface{}) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}

func randomString() string {
	length := rand.Intn(64)
	s := strings.Builder{}
	for i := 0; i < length; i++ {
		s.WriteByte(byte(rand.Intn(42) + 48))
	}
	return s.String()
}
