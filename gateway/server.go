package gateway

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/qingw1230/plato/common/config"
	"github.com/qingw1230/plato/common/prpc"
	"github.com/qingw1230/plato/common/tcp"
	"github.com/qingw1230/plato/gateway/rpc/client"
	"github.com/qingw1230/plato/gateway/rpc/service"
	"google.golang.org/grpc"
)

var cmdChannel chan *service.CmdContext

// RunMain 开启网关服务
func RunMain(path string) {
	config.Init(path)
	ln, err := net.ListenTCP("tcp", &net.TCPAddr{Port: config.GetGatewayTCPServerPort()})
	if err != nil {
		log.Fatalf("StartTCPEpollServer err:%s", err.Error())
		panic(err)
	}
	initWorkPool()
	initEpoll(ln, runProc)
	fmt.Println("------------------- im gateway started ----------------")
	cmdChannel = make(chan *service.CmdContext, config.GetGatewayCmdChannelNum())
	s := prpc.NewPServer(
		prpc.WithServiceName(config.GetGatewayServiceName()),
		prpc.WithIP(config.GetGatewayServiceAddr()),
		prpc.WithPort(config.GetGatewayRPCServerPort()),
		prpc.WithWeight(config.GetGatewayRPCWeight()),
	)
	fmt.Println(config.GetGatewayServiceName(),
		config.GetGatewayServiceAddr(),
		config.GetGatewayRPCServerPort(),
		config.GetGatewayRPCWeight(),
	)
	fmt.Println("------------- im gateway started ------------")
	s.RegisterService(func(server *grpc.Server) {
		service.RegisterGatewayServer(server, &service.Service{CmdChannel: cmdChannel})
	})

	client.Init()
	go cmdHandler()
	s.Start(context.TODO())
}

func runProc(c *connection, ep *epoller) {
	ctx := context.Background()
	dataBuf, err := tcp.ReadData(c.conn)
	if err != nil {
		// 读取时发现连接已断开，则将该连接从 ep 从删除
		if errors.Is(err, io.EOF) {
			ep.remove(c)
			// client.CancelConn(&ctx, getEndpoint(), int32(c.fd), nil)
		}
		return
	}

	err = wPool.Submit(func() {
		client.SendMsg(&ctx, getEndpoint(), int32(c.fd), dataBuf)
	})
	if err != nil {
		fmt.Errorf("runProc:err:%+v", err.Error())
	}
}

func cmdHandler() {
	for cmd := range cmdChannel {
		// 异步提交到协程池中完成发送任务
		switch cmd.Cmd {
		case service.DelConnCmd:
			wPool.Submit(func() { closeConn(cmd) })
		case service.PushCmd:
			wPool.Submit(func() { sendMsgByCmd(cmd) })
		default:
			panic("command undefined")
		}
	}
}

func closeConn(cmd *service.CmdContext) {
	if connPtr, ok := ep.tables.Load(cmd.FD); ok {
		conn, _ := connPtr.(*connection)
		conn.Close()
		ep.tables.Delete(cmd.FD)
	}
}

func sendMsgByCmd(cmd *service.CmdContext) {
	if connPtr, ok := ep.tables.Load(cmd.FD); ok {
		conn, _ := connPtr.(*connection)
		dp := tcp.DataPkg{
			Len:  uint32(len(cmd.Payload)),
			Data: cmd.Payload,
		}
		tcp.SendData(conn.conn, dp.Marshal())
	}
}

// getEndpoint 获取网关服务的 IP、Port
func getEndpoint() string {
	return fmt.Sprintf("%s:%d", config.GetGatewayServiceAddr(), config.GetGatewayRPCServerPort())
}
