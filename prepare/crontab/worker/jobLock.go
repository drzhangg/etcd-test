package worker

import (
	"context"
	"github.com/drzhangg/etcd-test/prepare/crontab/common"
	"go.etcd.io/etcd/clientv3"
)

// 分布式锁
type JobLock struct {
	//etcd客户端
	kv    clientv3.KV
	lease clientv3.Lease

	jobName   string             //任务名
	leaseId   clientv3.LeaseID   //租约id
	isLocked  bool               //是否上锁成功
	cancelFun context.CancelFunc //用于终止自动续租
}

// 初始化一把锁
func InitJobLock(jobName string, kv clientv3.KV, lease clientv3.Lease) (jobLock *JobLock) {
	return &JobLock{
		kv:      kv,
		lease:   lease,
		jobName: jobName,
	}
}

// 尝试上锁
func (jobLock *JobLock) TryLock() (err error) {
	var (
		leaseGrantResp     *clientv3.LeaseGrantResponse
		cancelCtx          context.Context
		cancelFunc         context.CancelFunc
		leaseId            clientv3.LeaseID
		leaseKeepAliveResp <-chan *clientv3.LeaseKeepAliveResponse
		txn                clientv3.Txn
		lockKey            string
		txnResp            *clientv3.TxnResponse
	)
	//1.创建5秒的租约
	if leaseGrantResp, err = jobLock.lease.Grant(context.TODO(), 5); err != nil {
		return
	}

	//context用于取消续租
	cancelCtx, cancelFunc = context.WithCancel(context.TODO())

	//租约id
	leaseId = leaseGrantResp.ID

	//2.自动续租
	if leaseKeepAliveResp, err = jobLock.lease.KeepAlive(cancelCtx, leaseId); err != nil {
		goto FAIL
	}

	//3.处理续租应答的协程
	go func() {
		var (
			keepResp *clientv3.LeaseKeepAliveResponse
		)
		for {
			select {
			case keepResp = <-leaseKeepAliveResp:
				if keepResp == nil {
					goto END
				}
			}
		}
	END:
	}()

	//4.创建事务txn
	txn = jobLock.kv.Txn(context.TODO())

	//获取锁路径
	lockKey = common.JOB_LOCK_DIR + jobLock.jobName

	//5.事务抢锁
	txn.If(clientv3.Compare(clientv3.CreateRevision(lockKey), "=", 0)).Then(clientv3.OpPut(lockKey, "", clientv3.WithLease(leaseId))).Else(clientv3.OpGet(lockKey))

	//提交事务
	if txnResp, err = txn.Commit(); err != nil {
		goto FAIL
	}

	//6.成功返回，失败释放租约
	if !txnResp.Succeeded { //锁被占用，抢锁失败
		err = common.ERR_LOCAL_ALREADY_REQUIRED
		goto FAIL
	}

	//抢锁成功
	jobLock.leaseId = leaseId
	jobLock.cancelFun = cancelFunc
	jobLock.isLocked = true
	return
FAIL:
	cancelFunc()                                  //取消自动续租
	jobLock.lease.Revoke(context.TODO(), leaseId) //释放租约
	return
}

// 释放锁
func (jobLock *JobLock) Unlock() {
	if jobLock.isLocked {
		jobLock.cancelFun()                                   //取消自动续租
		jobLock.lease.Revoke(context.TODO(), jobLock.leaseId) //释放租约
	}
}
