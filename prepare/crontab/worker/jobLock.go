package worker

import (
	"context"
	"github.com/coreos/etcd/clientv3"
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

// 尝试上锁
func (jobLock *JobLock) TryLock() (err error) {

}

// 释放锁
