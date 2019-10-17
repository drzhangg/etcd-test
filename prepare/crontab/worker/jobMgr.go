package worker

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/drzhangg/etcd-test/prepare/crontab/common"
)

// 任务管理器
type JobMgr struct {
	client  *clientv3.Client
	kv      clientv3.KV
	lease   clientv3.Lease
	watcher clientv3.Watcher
}

var (
	G_jobMgr *JobMgr
)

//监听任务变化
func (jobMgr *JobMgr) watchJobs() (err error) {
	var (
		getResp *clientv3.GetResponse
	)
	//1.get一下/zhang/cron/jobs/目录的后续变化
	if getResp, err = jobMgr.kv.Get(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithPrefix()); err != nil {
		return
	}

	//查看当前都有哪些任务
	go func() {
		//从GET时刻的后续版本开始监听变化

		//监听/zhang/cron/jobs/目录的后续变化

		//处理监听事件
	}()

	//2.从该revision向后监听变化事件

	return
}

//监听强杀任务通知
func (jobMgr *JobMgr) watchKiller() {

}

//初始化管理器
func InitJobMgr() (err error) {

	//初始化etcd配置

	//建立etcd连接

	//得到kv和api子集

	//赋值单例

	//启动任务监听

	//启动监听killer

	return
}

//创建任务执行锁
func (jobMgr *JobMgr) CreateJobLock(jobName string) (jobLock *JobLock) {
	jobLock = InitJobLock(jobName, jobMgr.kv, jobMgr.lease)
	return jobLock
}
