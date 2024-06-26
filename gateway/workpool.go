package gateway

import (
	"fmt"

	"github.com/panjf2000/ants"
	"github.com/qingw1230/plato/common/config"
)

var wPool *ants.Pool

// initWorkPool 初始化协程池
func initWorkPool() {
	var err error
	if wPool, err = ants.NewPool(config.GetGatewayWorkerPoolNum()); err != nil {
		fmt.Printf("InitWorkPoll.err: %s num:%d\n", err.Error(), config.GetGatewayWorkerPoolNum())
	}
}
