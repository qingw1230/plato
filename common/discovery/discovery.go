package discovery

import (
	"context"
	"sync"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/qingw1230/plato/common/config"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// ServiceDiscovery 服务发现
type ServiceDiscovery struct {
	client *clientv3.Client
	lock   sync.Mutex
	ctx    *context.Context
}

// NewServiceDiscovery 创建服务发现
func NewServiceDiscovery(ctx *context.Context) *ServiceDiscovery {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   config.GetEndpointsForDiscovery(),
		DialTimeout: config.GetTimeoutForDiscovery(),
	})
	if err != nil {
		logger.Fatal(err)
	}

	return &ServiceDiscovery{
		client: client,
		ctx:    ctx,
	}
}

// WatchService 初始化服务列表并监视
func (s *ServiceDiscovery) WatchService(prefix string, set, del func(key, value string)) error {
	// 获取具有 prefix 前缀的 key
	resp, err := s.client.Get(*s.ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	for _, ev := range resp.Kvs {
		set(string(ev.Key), string(ev.Value))
	}

	s.watcher(prefix, resp.Header.Revision+1, set, del)
	return nil
}

func (s *ServiceDiscovery) watcher(prefix string, rev int64, set, del func(key, value string)) {
	rch := s.client.Watch(*s.ctx, prefix, clientv3.WithPrefix(), clientv3.WithRev(rev))
	logger.CtxInfof(*s.ctx, "watching prefix:%s now ...", prefix)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case mvccpb.PUT:
				set(string(ev.Kv.Key), string(ev.Kv.Value))
			case mvccpb.DELETE:
				del(string(ev.Kv.Key), string(ev.Kv.Value))
			}
		}
	}
}

// Close 关闭服务
func (s *ServiceDiscovery) Close() error {
	return s.client.Close()
}
