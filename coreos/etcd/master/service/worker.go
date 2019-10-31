package service

import (
	"context"
	"github.com/drzhangg/etcd-test/coreos/etcd/common"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"time"
)

type Worker struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

var (
	G_worker *Worker
)

func InitWork() (err error) {
	var (
		config clientv3.Config
		client *clientv3.Client
		kv     clientv3.KV
		lease  clientv3.Lease
	)

	config = clientv3.Config{
		Endpoints:   G_config.Etcd.Endpoints,
		DialTimeout: time.Duration(G_config.Etcd.DialTimeout) * time.Second,
	}

	if client, err = clientv3.New(config); err != nil {
		return
	}

	kv = client.KV
	lease = client.Lease

	G_worker = &Worker{
		client: client,
		kv:     kv,
		lease:  lease,
	}

	return
}

// 获取节点列表
func (work *Worker) ListWorkers() (workers []string, err error) {
	var (
		getResp  *clientv3.GetResponse
		kv       *mvccpb.KeyValue
		workerIP string
	)

	//初始化数组
	workers = make([]string, 0)

	if getResp, err = work.kv.Get(context.TODO(), common.JOB_WORKER_DIR, clientv3.WithPrefix()); err != nil {
		return
	}

	for _, kv = range getResp.Kvs {
		workerIP = common.ExtractWorkerIP(string(kv.Key))
		workers = append(workers, workerIP)
	}
	return
}
