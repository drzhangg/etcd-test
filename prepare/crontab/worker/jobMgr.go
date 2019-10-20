package worker

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
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
		getResp            *clientv3.GetResponse
		kvpair             *mvccpb.KeyValue
		job                *common.Job
		jobEvent           *common.JobEvent
		watchStartRevision int64
		watchChan          clientv3.WatchChan
		watchResp          clientv3.WatchResponse
		watchEvent         mvccpb.Event
		jobName            string
	)
	//1.get一下/zhang/cron/jobs/目录的后续变化
	if getResp, err = jobMgr.kv.Get(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithPrefix()); err != nil {
		return
	}

	//查看当前都有哪些任务
	for _, kvpair = range getResp.Kvs {
		if job, err = common.UnpackJob(kvpair.Value); err != nil {
			//判断执行哪种操作
			jobEvent = common.BuildJobEvent(common.JOB_EVENT_SAVE, job)
			//将调度操作同步给调度协程
			G_scheduler.PushJobEvent(jobEvent)
		}
	}

	//2.从该revision向后监听变化事件
	go func() { //监听协程
		//从GET时刻的后续版本开始监听变化，要在当前版本+1
		watchStartRevision = getResp.Header.Revision + 1

		//监听/zhang/cron/jobs/目录的后续变化
		watchChan = jobMgr.watcher.Watch(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithRev(watchStartRevision), clientv3.WithPrefix())

		//处理监听事件
		for watchResp = range watchChan {
			for _, watchEvent = range watchResp.Events {
				switch watchEvent.Type {
				case mvccpb.PUT:
					//任务保存事件
					if job, err = common.UnpackJob(watchEvent.Kv.Value); err != nil {
						continue
					}
					// 构建一个更新Event
					jobEvent = common.BuildJobEvent(common.JOB_EVENT_SAVE, job)
				case mvccpb.DELETE:
					// 删除任务
					jobName = common.ExtractJobName(string(watchEvent.Kv.Key))

					job = &common.Job{Name: jobName}

					// 构建一个删除Event
					jobEvent = common.BuildJobEvent(common.JOB_EVENT_DELETE, job)
				}
				// 变化推给Scheduler
				G_scheduler.PushJobEvent(jobEvent)
			}
		}
	}()
	return
}

//监听强杀任务通知
func (jobMgr *JobMgr) watchKiller() {
	var (
		watchChan     clientv3.WatchChan
		watchChanResp clientv3.WatchResponse
		watchEvent    mvccpb.Event
		jobName       string
		job           *common.Job
		jobEvent      *common.JobEvent
	)

	go func() {
		watchChan = jobMgr.watcher.Watch(context.TODO(), common.JOB_KILLER_DIR, clientv3.WithPrefix())
		for watchChanResp = range watchChan {

			for _, watchEvent = range watchChanResp.Events {
				switch watchEvent.Type {
				case mvccpb.PUT:	//杀死任务事件
					jobName = common.ExtractKillerName(string(watchEvent.Kv.Key))
					job = &common.Job{Name: jobName}
					jobEvent = common.BuildJobEvent(common.JOB_EVENT_KILL, job)

					G_scheduler.PushJobEvent(jobEvent)
				case mvccpb.DELETE:

				}
			}
		}

	}()

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
