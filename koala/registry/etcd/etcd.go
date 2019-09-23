package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/drzhangg/etcd-test/koala/registry"
	"go.etcd.io/etcd/clientv3"
	"path"
	"sync"
	"sync/atomic"
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

	value              atomic.Value
	lock               sync.Mutex
	registryServiceMap map[string]*RegisterService
}

type AllServiceInfo struct {
	serviceMap map[string]*registry.Service
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
	allServiceInfo := &AllServiceInfo{
		serviceMap: make(map[string]*registry.Service, MaxServiceNum),
	}
	etcdRegistry.value.Store(allServiceInfo)
	registry.RegisterPlugin(etcdRegistry)
	go etcdRegistry.run()
}

// Name 插件名字
func (e *EtcdRegistry) Name() string {
	return "etcd"
}

// 实现接口方法
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

// Register 服务注册 实现接口方法
func (e *EtcdRegistry) Register(ctx context.Context, service *registry.Service) (err error) {
	select {
	case e.serviceCh <- service:
	default:
		err = fmt.Errorf("register chan is full")
		return
	}
	return
}

// Unregister  反注册 实现接口方法
func (e *EtcdRegistry) Unregister(ctx context.Context, service *registry.Service) (err error) {
	return
}

// getServiceFromCache 从缓存中读取数据
func (e *EtcdRegistry) getServiceFromCache(ctx context.Context, name string) (service *registry.Service, ok bool) {
	// 从缓存中加载数据
	allServiceInfo := e.value.Load().(*AllServiceInfo)
	service, ok = allServiceInfo.serviceMap[name]
	return
}

func (e *EtcdRegistry) run() {

	ticker := time.NewTicker(MaxSyncServiceInterval)
	for {
		select {
		case service := <-e.serviceCh:
			registryService, ok := e.registryServiceMap[service.Name]
			if ok {
				for _, node := range service.Nodes {
					registryService.service.Nodes = append(registryService.service.Nodes, node)

				}
				registryService.registered = false
				break
			}
			registryService = &RegisterService{
				service: service,
			}
			e.registryServiceMap[service.Name] = registryService
		case <-ticker.C:
			e.syncServiceFromEtcd()
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
		fmt.Println("key:", key)

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

func (e *EtcdRegistry) servicePath(name string) string {
	return path.Join(e.options.RegistryPath, name)
}

// GetService 实现接口方法
func (e *EtcdRegistry) GetService(ctx context.Context, name string) (service *registry.Service, err error) {
	//一般情况下,都会从缓存中读取
	service, ok := e.getServiceFromCache(ctx, name)
	if ok {
		return
	}

	//如果缓存中没有这个sservice,则从etcd中读取
	e.lock.Lock()
	defer e.lock.Unlock()
	//先检测，是否从etcd中加载成功了
	service, ok = e.getServiceFromCache(ctx, name)
	if ok {
		return
	}

	//从etcd中读取指定服务名字的服务信息
	key := e.servicePath(name)
	resp, err := e.client.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return
	}

	service = &registry.Service{
		Name: name,
	}
	for _, kv := range resp.Kvs {
		value := kv.Value
		var tmpService registry.Service
		err = json.Unmarshal(value, &tmpService)
		if err != nil {
			return
		}

		for _, node := range tmpService.Nodes {
			service.Nodes = append(service.Nodes, node)
		}
	}
	allServiceInfoOld := e.value.Load().(*AllServiceInfo)
	var allServiceInfoNew = &AllServiceInfo{
		serviceMap: make(map[string]*registry.Service, MaxServiceNum),
	}
	for key, val := range allServiceInfoOld.serviceMap {
		allServiceInfoNew.serviceMap[key] = val
	}

	allServiceInfoNew.serviceMap[name] = service
	e.value.Store(allServiceInfoNew)
	return
}

func (e *EtcdRegistry) syncServiceFromEtcd() {
	var allServiceInfoNew = &AllServiceInfo{serviceMap: make(map[string]*registry.Service, MaxServiceNum)}

	ctx := context.TODO()
	allServiceInfo := e.value.Load().(*AllServiceInfo)

	//对于缓存的每一个服务，都需要从etcd中进行更新
	for _, service := range allServiceInfo.serviceMap {
		key := e.servicePath(service.Name)
		resp, err := e.client.Get(ctx, key, clientv3.WithPrefix())
		if err != nil {
			allServiceInfoNew.serviceMap[service.Name] = service
			continue
		}

		serviceNew := &registry.Service{
			Name: service.Name,
		}

		for _, kv := range resp.Kvs {
			value := kv.Value
			var tmpService registry.Service
			err = json.Unmarshal(value, &tmpService)
			if err != nil {
				return
			}

			for _, node := range tmpService.Nodes {
				serviceNew.Nodes = append(serviceNew.Nodes, node)
			}
		}
		allServiceInfoNew.serviceMap[serviceNew.Name] = serviceNew
	}
	e.value.Store(allServiceInfoNew)
}
