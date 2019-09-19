package etcd

import (
	"context"
	"fmt"
	"github.com/drzhangg/etcd-test/koala/registry"
	"go.etcd.io/etcd/clientv3"
	"path"
)

type EtcdRegistry struct {
	options   *registry.Options
	client    *clientv3.Client
	serviceCh chan *registry.Service
}

var (
	etcdRegistry *EtcdRegistry = &EtcdRegistry{
		serviceCh: make(chan *registry.Service, 8),
	}
)

func init() {
	registry.RegisterPlugin(etcdRegistry)
	go etcdRegistry.run()
}

// Name 插件名字
func (e *EtcdRegistry) Name() string {
	return "etcd"
}

func (e *EtcdRegistry) Init(ctx context.Context, opts ...registry.Option) (err error) {
	e.options = &registry.Options{}
	for _, opt := range opts {
		opt(e.options)
	}

	e.client, err = clientv3.New(clientv3.Config{
		Endpoints:   e.options.Address,
		DialTimeout: e.options.Timeout,
	})
	if err != nil {
		return
	}
	return
}

// Register 服务注册
func (e *EtcdRegistry) Register(ctx context.Context, service *registry.Service) (err error) {
	select {
	case e.serviceCh <- service:
	default:
		err = fmt.Errorf("register chan is full")
		return
	}
	return
}

// Unregister 反注册
func (e *EtcdRegistry) Unregister(ctx context.Context, service *registry.Service) (err error) {
	return
}

func (e *EtcdRegistry) run() {

}

func (e *EtcdRegistry) servicePath(service *registry.Service) string {
	return path.Join(e.options.RegistryPath, service.Name)
}
