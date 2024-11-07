package discovery

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/qingw1230/plato/common/config"
)

func init() {
	config.Init("/home/qgw/code/plato/im.yaml")
}

func TestServiceRegister(t *testing.T) {
	ctx := context.Background()
	e := EndpointInfo{
		IP:   "127.0.0.1",
		Port: "9999",
	}
	server, err := NewServiceRegister(&ctx, "/web/node1", &e, 5)
	if err != nil {
		log.Panicln(err)
	}

	go server.ListenLeaseResponseChan()
	<-time.After(10 * time.Second)
	server.Close()
}
