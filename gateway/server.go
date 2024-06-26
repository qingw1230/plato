package gateway

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/qingw1230/plato/common/config"
	"github.com/qingw1230/plato/common/tcp"
)

// RunMain 开启网关服务
func RunMain(path string) {
	config.Init(path)
	ln, err := net.ListenTCP("tcp", &net.TCPAddr{
		Port: config.GetGatewayServerPort(),
	})
	if err != nil {
		log.Fatalf("StartTCPEpollServer err:%s", err.Error())
		panic(err)
	}
	initWorkPool()
	initEpoll(ln, runProc)
	fmt.Println("------------------- im gateway started ----------------")
	select {}
}

func runProc(c *connection, ep *epoller) {
	dataBuf, err := tcp.ReadData(c.conn)
	if err != nil {
		// 读取时发现连接已断开，则将该连接从 ep 从删除
		if errors.Is(err, io.EOF) {
			ep.remove(c)
		}
		return
	}

	err = wPool.Submit(func() {
		bytes := tcp.DataPkg{
			Len:  uint32(len(dataBuf)),
			Data: dataBuf,
		}
		tcp.SendData(c.conn, bytes.Marshal())
	})
	if err != nil {
		fmt.Errorf("runProc:err:%+v", err.Error())
	}
}
