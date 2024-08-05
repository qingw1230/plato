package state

import (
	"context"
	"fmt"

	"github.com/qingw1230/plato/common/config"
	"github.com/qingw1230/plato/common/prpc"

	"github.com/qingw1230/plato/state/rpc/client"
	"github.com/qingw1230/plato/state/rpc/service"
	"google.golang.org/grpc"
)

var cmdChannel chan *service.CmdContext

// RunMain 启动网关服务
func RunMain(path string) {
	config.Init(path)
	cmdChannel = make(chan *service.CmdContext, config.GetSateCmdChannelNum())

	s := prpc.NewPServer(
		prpc.WithServiceName(config.GetStateServiceName()),
		prpc.WithIP(config.GetSateServiceAddr()),
		prpc.WithPort(config.GetSateServerPort()),
		prpc.WithWeight(config.GetSateRPCWeight()),
	)

	s.RegisterService(func(server *grpc.Server) {
		service.RegisterStateServer(server, &service.Service{CmdChannel: cmdChannel})
	})
	// 初始化 RPC 客户端
	client.Init()
	// 启动命令处理写协程
	go cmdHandler()
	// 启动 rpc server
	s.Start(context.TODO())
}

func cmdHandler() {
	for cmd := range cmdChannel {
		switch cmd.Cmd {
		case service.CancelConnCmd:
			fmt.Printf("cancelconn endpoint:%s, fd:%d, data:%+v", cmd.Endpoint, cmd.FD, cmd.Playload)
		case service.SendMsgCmd:
			fmt.Println("cmdHandler", string(cmd.Playload))
			client.Push(cmd.Ctx, int32(cmd.FD), cmd.Playload)
		}
	}
}
