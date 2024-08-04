package discovery

import (
	"context"
	"log"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/qingw1230/plato/common/config"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// ServiceRegister 服务注册
type ServiceRegister struct {
	client        *clientv3.Client
	leaseID       clientv3.LeaseID
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse // 用于 lease 保活
	key           string                                  // 服务的 Web 访问路径
	value         string                                  // 序列化后的 EndpointInfo
	ctx           *context.Context
}

// NewServiceRegister 创建服务，向 etcd 中添加键值对
// key 服务的 Web 访问路径，作为 etcd 的 key，例如 /im/ip_dispatcher/node1
// endportinfo 机器信息，作为 etcd 的 value
// ttl 添加到 etcd 中键值对的 TTL
func NewServiceRegister(ctx *context.Context, key string, endpointinfo *EndpointInfo, ttl int64) (*ServiceRegister, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   config.GetEndpointsForDiscovery(),
		DialTimeout: config.GetTimeoutForDiscovery(),
	})
	if err != nil {
		log.Fatal(err)
	}

	service := &ServiceRegister{
		client: client,
		key:    key,
		value:  endpointinfo.Marshal(),
		ctx:    ctx,
	}

	if err := service.putKeyWithLease(ttl); err != nil {
		return nil, err
	}
	return service, nil
}

// putKeyWithLease 创建 lease，向 etcd 中添加 s 保存的键值对，并将 lease 绑定到该键值对上
func (s *ServiceRegister) putKeyWithLease(ttl int64) error {
	response, err := s.client.Grant(*s.ctx, ttl)
	if err != nil {
		return err
	}

	_, err = s.client.Put(*s.ctx, s.key, s.value, clientv3.WithLease(response.ID))
	if err != nil {
		return err
	}

	// 使给定租约永久有效
	leaseResponseChan, err := s.client.KeepAlive(*s.ctx, response.ID)
	if err != nil {
		return err
	}

	s.leaseID = response.ID
	s.keepAliveChan = leaseResponseChan
	return nil
}

// UpdateValue 更新机器信息
func (s *ServiceRegister) UpdateValue(val *EndpointInfo) error {
	value := val.Marshal()
	_, err := s.client.Put(*s.ctx, s.key, value, clientv3.WithLease(s.leaseID))
	if err != nil {
		return err
	}
	s.value = value
	logger.CtxInfof(*s.ctx, "ServiceRegister.UpdateValue leaseID=%d key=%s,value=%s, success!", s.leaseID, s.key, s.value)
	return nil
}

// ListenLeaseResponChan 监听用于 lease 保活的 chan
func (s *ServiceRegister) ListenLeaseResponChan() {
	for leaseKeepRespon := range s.keepAliveChan {
		logger.CtxInfof(*s.ctx, "lease success leaseID:%d, key:%s, value:%s reps:+%v", s.leaseID, s.key, s.value, leaseKeepRespon)
	}
	logger.CtxInfof(*s.ctx, "lease faild! leaseID %d, key:%s, value:%s", s.leaseID, s.key, s.value)
}

// Close 注销服务，从 etcd 中删除键值对
func (s *ServiceRegister) Close() error {
	// 撤销给定的 lease，附加到该 lease 的 key 会过期并被删除
	if _, err := s.client.Revoke(context.Background(), s.leaseID); err != nil {
		return err
	}
	logger.CtxInfof(*s.ctx, "lease close! leaseID:%d, key:%s, value:%s", s.leaseID, s.key, s.value)
	return s.client.Close()
}
