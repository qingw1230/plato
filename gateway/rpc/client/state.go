package client

import (
	"context"
	"time"

	"github.com/qingw1230/plato/common/config"
	"github.com/qingw1230/plato/common/prpc"
	"github.com/qingw1230/plato/state/rpc/service"
)

var stateClient service.StateClient

func initStateClient() {
	pCli, err := prpc.NewPClient(config.GetStateServiceName())
	if err != nil {
		panic(err)
	}
	stateClient = service.NewStateClient(pCli.Conn())
}

func CancelConn(ctx *context.Context, endpoint string, fd int32, playLoad []byte) error {
	rpcCtx, _ := context.WithTimeout(*ctx, 100*time.Millisecond)
	stateClient.CancelConn(rpcCtx, &service.StateRequest{
		Endpoint: endpoint,
		Fd:       fd,
		Data:     playLoad,
	})
	return nil
}

func SendMsg(ctx *context.Context, endpoint string, fd int32, playLoad []byte) error {
	rpcCtx, _ := context.WithTimeout(*ctx, 100*time.Millisecond)
	_, err := stateClient.SendMsg(rpcCtx, &service.StateRequest{
		Endpoint: endpoint,
		Fd:       fd,
		Data:     playLoad,
	})
	if err != nil {
		panic(err)
	}
	return nil
}
