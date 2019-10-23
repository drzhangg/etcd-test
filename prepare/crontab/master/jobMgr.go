package master

import (
	"context"
	"encoding/json"
	"github.com/drzhangg/etcd-test/prepare/crontab/common"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"time"
)

// 任务管理器，进行任务的crud
type JobMgr struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

var (
	G_jobMgr *JobMgr
)

//初始化管理器，进行一些etcd的初始化
func InitJobMgr() (err error) {
	var (
		config clientv3.Config
		client *clientv3.Client
		kv     clientv3.KV
		lease  clientv3.Lease
	)
	config = clientv3.Config{
		Endpoints:   G_config.EtcdEndpoints,
		DialTimeout: time.Duration(G_config.EtcdDialTimeout) * time.Millisecond,
	}

	//建立etcd连接
	if client, err = clientv3.New(config); err != nil {
		return
	}

	//得到kv和lease
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)

	//赋值单例
	G_jobMgr = &JobMgr{
		client: client,
		kv:     kv,
		lease:  lease,
	}
	return
}

// 保存任务
func (jobMgr *JobMgr) SaveJob(job *common.Job) (oldJob *common.Job, err error) {
	var (
		jobKey    string
		jobValue  []byte
		putResp   *clientv3.PutResponse
		oldJobObj common.Job
	)

	//etcd保存的key
	jobKey = common.JOB_SAVE_DIR + job.Name

	//任务信息
	if jobValue, err = json.Marshal(job); err != nil {
		return
	}

	//保存到etcd,withPrevKV()是为了获取操作前已经有的key-value
	if putResp, err = jobMgr.kv.Put(context.TODO(), jobKey, string(jobValue), clientv3.WithPrevKV()); err != nil {
		return
	}
	//如果当前是更新操作，那么返回旧的值
	if putResp.PrevKv != nil { //putResp.PrevKv获取put前的值
		if err = json.Unmarshal(putResp.PrevKv.Value, &oldJobObj); err != nil {
			err = nil
			return
		}
		oldJob = &oldJobObj
	}
	return
}

// 删除任务
func (jobMgr *JobMgr) DeleteJob(name string) (oldJob *common.Job, err error) {
	var (
		jobKey    string
		delResp   *clientv3.DeleteResponse
		oldJobObj common.Job
	)
	//etcd中保存任务的key
	jobKey = common.JOB_SAVE_DIR + name

	if delResp, err = jobMgr.kv.Delete(context.TODO(), jobKey, clientv3.WithPrevKV()); err != nil {
		return
	}

	//返回被删除的任务信息
	if delResp.PrevKvs != nil {
		if err = json.Unmarshal(delResp.PrevKvs[0].Value, &oldJobObj); err != nil {
			err = nil
			return
		}
		oldJob = &oldJobObj
	}
	return
}

// 列举全部任务
func (jobMgr *JobMgr) ListJobs() (jobList []*common.Job, err error) {
	var (
		jobKey  string
		getResp *clientv3.GetResponse
		kvPair  *mvccpb.KeyValue
		job     *common.Job
	)

	jobKey = common.JOB_SAVE_DIR

	//get withPrefix
	if getResp, err = jobMgr.kv.Get(context.TODO(), jobKey, clientv3.WithPrefix()); err != nil {
		return
	}

	jobList = make([]*common.Job, 0)
	for _, kvPair = range getResp.Kvs {
		if err = json.Unmarshal(kvPair.Value, &job); err != nil {
			err = nil
			return
		}
		jobList = append(jobList, job)
	}
	return
}

//杀死任务
func (jobMgr *JobMgr) KillJob(name string) (err error) {
	var (
		killKey        string
		leaseGrantResp *clientv3.LeaseGrantResponse
	)
	killKey = common.JOB_KILLER_DIR + name

	//让租约自动过期
	if leaseGrantResp, err = jobMgr.lease.Grant(context.TODO(), 1); err != nil {
		return
	}

	//标记任务为kill状态
	if _, err = jobMgr.kv.Put(context.TODO(), killKey, "", clientv3.WithLease(leaseGrantResp.ID)); err != nil {
		return
	}
	return
}
