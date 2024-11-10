package source

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/qingw1230/plato/common/config"
	"github.com/qingw1230/plato/common/discovery"
)

func testServiceRegister(ctx *context.Context, port, node string) {
	go func() {
		// 创建一台机器，注册到 etcd
		ed := discovery.EndpointInfo{
			IP:   "127.0.0.1",
			Port: port,
			MetaData: map[string]interface{}{
				"connect_num":   float64(rand.Int63n(100000)),
				"message_bytes": float64(rand.Int63n(100 << 20)),
			},
		}
		s, err := discovery.NewServiceRegister(ctx, fmt.Sprintf("%s/%s", config.GetServicePathForIPConf(), node), &ed, 5)
		if err != nil {
			panic(err)
		}

		go s.ListenLeaseResponseChan()

		for {
			// 更新机器状态
			ed = discovery.EndpointInfo{
				IP:   "127.0.0.1",
				Port: port,
				MetaData: map[string]interface{}{
					"connect_num":   float64(rand.Int63n(100000)),
					"message_bytes": float64(rand.Int63n(100 << 20)),
				},
			}
			s.UpdateValue(&ed)
			time.Sleep(5 * time.Second)
		}
	}() // go funv() {
}
