package service

import (
	"context"
	"encoding/json"
	"github.com/coreos/etcd/mvcc/mvccpb"
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
func (manager *Manager) ListJobs() (jobList []*common.Job, err error) {
	var (
		jobKey      string
		getResponse *clientv3.GetResponse
		keyVal      *mvccpb.KeyValue
		job         *common.Job
	)
	jobKey = common.JOB_SAVE_DIR

	if getResponse, err = manager.kv.Get(context.TODO(), jobKey, clientv3.WithPrefix()); err != nil {
		return
	}

	jobList = make([]*common.Job, 0)
	for _, keyVal = range getResponse.Kvs {

		if err = json.Unmarshal(keyVal.Value, &job); err != nil {
			return
		}
		jobList = append(jobList, job)
	}
	return
}

//删除任务
func (manager *Manager) DelJob(name string) (oldJob *common.Job, err error) {
	var (
		jobKey       string
		delResp      *clientv3.DeleteResponse
		oldJobObject common.Job
	)

	jobKey = common.JOB_SAVE_DIR
	if delResp, err = manager.kv.Delete(context.TODO(), jobKey, clientv3.WithPrevKV()); err != nil {
		return
	}

	//返回被删除的数据
	if delResp.PrevKvs != nil {
		if err = json.Unmarshal(delResp.PrevKvs[0].Value, &oldJobObject); err != nil {
			return
		}
	}
	oldJob = &oldJobObject
	return
}

//杀死任务 kill
func (manager *Manager) KillJob(name string) (err error) {
	var (
		jobKey         string
		leaseGrantResp *clientv3.LeaseGrantResponse
	)

	jobKey = common.JOB_KILLER_DIR + name

	//让租约自动过期
	if leaseGrantResp, err = manager.lease.Grant(context.TODO(), 1); err != nil {
		return
	}

	//标记任务为kill状态
	if _, err = manager.kv.Put(context.TODO(), jobKey, "", clientv3.WithLease(leaseGrantResp.ID)); err != nil {
		return
	}
	return
}
