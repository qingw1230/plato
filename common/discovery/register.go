package discovery

import (
	"context"
	"log"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/qingw1230/plato/common/config"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// ServiceRegister 用于注册服务
type ServiceRegister struct {
	ctx           *context.Context
	client        *clientv3.Client
	leaseID       clientv3.LeaseID                        // 服务的 Lease ID
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse // 保存保活响应
	key           string                                  // 服务的 Web 访问路径
	value         string                                  // 序列化后的 EndpointInfo
}

// NewServiceRegister 注册服务，即向 etcd 中添加键值对
func NewServiceRegister(ctx *context.Context, key string, endpointInfo *EndpointInfo, ttl int64) (*ServiceRegister, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   config.GetEndpointsForDiscovery(),
		DialTimeout: config.GetTimeouForDiscovery(),
	})
	if err != nil {
		log.Panic(err)
	}

	service := &ServiceRegister{
		ctx:    ctx,
		client: client,
		key:    key,
		value:  endpointInfo.Marshal(),
	}

	if err := service.putKeyWithLease(ttl); err != nil {
		return nil, err
	}
	return service, nil
}

func (s *ServiceRegister) putKeyWithLease(ttl int64) error {
	// Grant 创建一个租约，当服务器在给定 ttl 内没有收到 keepAlive 时租约过期
	// 如果租约过期所有附加在租约上的 key 将过期被删除，即服务下线
	response, err := s.client.Grant(*s.ctx, ttl)
	if err != nil {
		return err
	}

	_, err = s.client.Put(*s.ctx, s.key, s.value, clientv3.WithLease(response.ID))
	if err != nil {
		return err
	}

	leaseResponseChan, err := s.client.KeepAlive(*s.ctx, response.ID)
	if err != nil {
		return err
	}

	s.leaseID = response.ID
	s.keepAliveChan = leaseResponseChan
	return err
}

// Close 注销服务
func (s *ServiceRegister) Close() error {
	// 附加到该 leaseID 的 key 会过期被并删除
	if _, err := s.client.Revoke(context.Background(), s.leaseID); err != nil {
		return err
	}
	logger.CtxInfof(*s.ctx, "lease close! leaseID:%d, key:%s, value:%s", s.leaseID, s.key, s.value)
	return s.client.Close()
}

// UpdateValue 更新机器信息
func (s *ServiceRegister) UpdateValue(val *EndpointInfo) error {
	value := val.Marshal()
	_, err := s.client.Put(*s.ctx, s.key, value, clientv3.WithLease(s.leaseID))
	if err != nil {
		return err
	}
	s.value = value
	logger.CtxInfof(*s.ctx, "ServiceRegister.UpdateValue leaseID:%d, key:%s, value:%s, success!", s.leaseID, s.key, s.value)
	return nil
}

func (s *ServiceRegister) ListenLeaseResponseChan() {
	for response := range s.keepAliveChan {
		logger.CtxInfof(*s.ctx, "lease success leaseID:%d, key:%s, value:%s, reps:%+v", s.leaseID, s.key, s.value, response)
	}
	logger.CtxInfof(*s.ctx, "lease faild! leaseID:%d, key:%s, value:%s", s.leaseID, s.key, s.value)
}
