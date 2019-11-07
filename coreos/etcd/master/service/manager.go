package service

import (
	"go.etcd.io/etcd/clientv3"
	"time"
)

type Manager struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

var (
	G_manager *Manager
)

func InitManager() (err error) {
	var (
		config clientv3.Config
		client *clientv3.Client
		kv     clientv3.KV
		lease  clientv3.Lease
	)
	config = clientv3.Config{
		Endpoints:   G_config.Etcd.Endpoints,
		DialTimeout: time.Duration(G_config.Etcd.DialTimeout),
	}

	if client, err = clientv3.New(config); err != nil {
		return
	}
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)

	G_manager = &Manager{
		client: client,
		kv:     kv,
		lease:  lease,
	}

	return
}
