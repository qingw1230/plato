package plugin

import (
	"fmt"
	"time"

	"github.com/qingw1230/plato/common/prpc/config"
	"github.com/qingw1230/plato/common/prpc/discov"
	"github.com/qingw1230/plato/common/prpc/discov/etcd"
)

// GetDiscovInstance 获取一个服务发现实例
func GetDiscovInstance() (discov.Discovery, error) {
	name := config.GetDiscovName()
	switch name {
	case "etcd":
		return etcd.NewETCDRegister(
			etcd.WithEndpoints(config.GetDiscovEndpoints()),
			etcd.WithRegisterServiceOrKeepAliveInterval(10*time.Millisecond),
		)
	}

	return nil, fmt.Errorf("not exist plugin:%s", name)
}
