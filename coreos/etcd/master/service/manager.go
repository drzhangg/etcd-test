package service

import (
	"context"
	"encoding/json"
	"github.com/drzhangg/etcd-test/coreos/etcd/common"
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

//保存任务，把旧的任务信息返回
func (manager *Manager) SaveManager(job *common.Job) (oldJob *common.Job, err error) {
	var (
		jobKey    string
		datas     []byte
		putResp   *clientv3.PutResponse
		objectJob common.Job
	)
	//etcd保存任务的key
	jobKey = common.JOB_SAVE_DIR + job.Name
	//序列化job数据
	if datas, err = json.Marshal(&job); err != nil {
		return
	}

	if putResp, err = manager.kv.Put(context.TODO(), jobKey, string(datas), clientv3.WithPrevKV()); err != nil {
		return
	}

	if putResp != nil {
		if err = json.Unmarshal(putResp.PrevKv.Value, &objectJob); err != nil {
			return
		}
		oldJob = &objectJob
	}
	return
}

//全部任务

//删除任务

//杀死任务
