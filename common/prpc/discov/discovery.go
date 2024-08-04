package discov

import "context"

type Discovery interface {
	// Name 服务发现用的什么 eg etcd zk
	Name() string
	// Register 注册服务
	Register(ctx context.Context, service *Service)
	// UnRegister 注销服务
	UnRegister(ctx context.Context, service *Service)
	// GetService 获取由 name 指定的服务 IP Port 信息
	GetService(ctx context.Context, name string) *Service
	AddListener(ctx context.Context, f func())
	NotifyListeners()
}
