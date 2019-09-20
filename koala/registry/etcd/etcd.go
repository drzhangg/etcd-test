package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/drzhangg/etcd-test/koala/registry"
	"go.etcd.io/etcd/clientv3"
	"path"
	"time"
)

const (
	MaxServiceNum          = 8
	MaxSyncServiceInterval = time.Second * 10
)

type EtcdRegistry struct {
	options   *registry.Options
	client    *clientv3.Client
	serviceCh chan *registry.Service

	registryServiceMap map[string]*RegisterService
}

type RegisterService struct {
	id      clientv3.LeaseID
	service *registry.Service

	registered  bool
	keepAliveCh <-chan *clientv3.LeaseKeepAliveResponse
}

var (
	etcdRegistry *EtcdRegistry = &EtcdRegistry{
		serviceCh:          make(chan *registry.Service, 8),
		registryServiceMap: make(map[string]*RegisterService, MaxServiceNum),
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

	//ticker := time.NewTicker(MaxSyncServiceInterval)
	for {
		select {
		case service := <-e.serviceCh:
			_, ok := e.registryServiceMap[service.Name]
			if ok {
				break
			}
			registryService := &RegisterService{
				service: service,
			}
			e.registryServiceMap[service.Name] = registryService
		default:
			e.registerOrKeepAlive()
			time.Sleep(time.Millisecond * 500)
		}
	}
}

// registerOrKeepAlive 注册或者续租
func (e *EtcdRegistry) registerOrKeepAlive() {
	for _, registryService := range e.registryServiceMap {
		if registryService.registered {
			e.keepAlive(registryService)
			continue
		}

		e.registerService(registryService)
	}
}

// keepAlive 续租
func (e *EtcdRegistry) keepAlive(registryService *RegisterService) {
	select {
	case resp := <-registryService.keepAliveCh:
		if resp == nil {
			registryService.registered = false
			return
		}
		fmt.Printf("services:%s node:%s ttl:%v\n", registryService.service.Name, registryService.service.Nodes[0].Ip, registryService.service.Nodes[0].Port)
	}
	return

	return
}

// registerService 注册服务
func (e *EtcdRegistry) registerService(registryService *RegisterService) (err error) {
	resp, err := e.client.Grant(context.TODO(), e.options.HeartBeat)
	if err != nil {
		return
	}

	// 获取租约id
	registryService.id = resp.ID
	for _, node := range registryService.service.Nodes {
		tmp := &registry.Service{
			Name: registryService.service.Name,
			Nodes: []*registry.Node{
				node,
			},
		}

		data, err := json.Marshal(tmp)
		if err != nil {
			continue
		}

		//获取节点路径
		key := e.serviceNodePath(tmp)

		_, err = e.client.Put(context.TODO(), key, string(data), clientv3.WithLease(resp.ID))
		if err != nil {
			continue
		}

		ch, err := e.client.KeepAlive(context.TODO(), resp.ID)
		if err != nil {
			continue
		}
		registryService.keepAliveCh = ch
		registryService.registered = true
	}
	return
}

// serviceNodePath 获取节点路径
func (e *EtcdRegistry) serviceNodePath(service *registry.Service) string {
	nodeIP := fmt.Sprintf("%s:%d", service.Nodes[0].Ip, service.Nodes[0].Port)
	return path.Join(e.options.RegistryPath, nodeIP)
}

func (e *EtcdRegistry) servicePath(service *registry.Service) string {
	return path.Join(e.options.RegistryPath, service.Name)
}
