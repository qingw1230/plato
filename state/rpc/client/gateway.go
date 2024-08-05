package client

import (
	"context"
	"fmt"
	"time"

	"github.com/qingw1230/plato/common/config"
	"github.com/qingw1230/plato/common/prpc"
	"github.com/qingw1230/plato/gateway/rpc/service"
)

var gatewayClient service.GatewayClient

func initGatewayClient() {
	pCli, err := prpc.NewPClient(config.GetGatewayServiceName())
	if err != nil {
		panic(err)
	}
	gatewayClient = service.NewGatewayClient(pCli.Conn())
}

func DelConn(ctx *context.Context, fd int32, playLoad []byte) error {
	rpcCtx, _ := context.WithTimeout(*ctx, 100*time.Millisecond)
	gatewayClient.DelConn(rpcCtx, &service.GatewayRequest{Fd: fd, Data: playLoad})
	return nil
}

func Push(ctx *context.Context, fd int32, playLoad []byte) error {
	rpcCtx, _ := context.WithTimeout(*ctx, 100*time.Second)
	resp, err := gatewayClient.Push(rpcCtx, &service.GatewayRequest{Fd: fd, Data: playLoad})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)
	return nil
}
