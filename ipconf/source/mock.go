package source

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/qingw1230/plato/common/config"
	"github.com/qingw1230/plato/common/discovery"
)

// testServiceRegister 注册服务，模拟机器状态不断更新，用于测试 /im/ip_dispatcher 正确性
func testServiceRegister(ctx *context.Context, port, node string) {
	go func() {
		// 创建一台机器，并注册到 etcd 中
		ed := discovery.EndportInfo{
			IP:   "127.0.0.1",
			Port: port,
			MetaData: map[string]interface{}{
				"connect_num":   float64(rand.Int63n(1000000)),
				"message_bytes": float64(rand.Int63n(10000000000)),
			},
		}
		s, err := discovery.NewServiceRegister(ctx, fmt.Sprintf("%s/%s", config.GetServicePathForIPConf(), node), &ed, time.Now().Unix())
		if err != nil {
			panic(err)
		}

		go s.ListenLeaseResponChan()

		for {
			// 每 5 秒更新一下机器状态
			ed = discovery.EndportInfo{
				IP:   "127.0.0.1",
				Port: port,
				MetaData: map[string]interface{}{
					"connect_num":   float64(rand.Int63n(1000000)),
					"message_bytes": float64(rand.Int63n((10000000000))),
				},
			}
			s.UpdateValue(&ed)
			time.Sleep(5 * time.Second)
		}
	}() // go func() {
}
