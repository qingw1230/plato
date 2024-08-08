package gateway

import (
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/qingw1230/plato/common/config"
)

var (
	path = "/home/qgw/code/plato/im.yaml"
)

func BenchmarkEpoll(b *testing.B) {
	config.Init(path)
	fmt.Println(config.GetGatewayEpollerChanNum())
	go func() {
		ln, err := net.ListenTCP("tcp", &net.TCPAddr{Port: config.GetGatewayTCPServerPort()})
		if err != nil {
			log.Fatalf("StartTCPEpollServer err:%s", err.Error())
			panic(err)
		}
		initWorkPool()
		initEpoll(ln, runProc)
	}()

	timeout := time.After(20 * time.Second)
	<-timeout
	b.Error("Benchmark timed out after 1 seconds")
}

func BenchmarkRunMain(b *testing.B) {
	config.Init(path)
	timeout := time.After(5 * time.Second)

	go func() {
		RunMain(path)
	}()

	<-timeout
	b.Error("Benchmark timed out after 5 seconds")
}
