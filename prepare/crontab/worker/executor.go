package worker

import (
	"github.com/drzhangg/etcd-test/prepare/crontab/common"
	"github.com/drzhangg/go-crontab/worker"
	"os/exec"
)

// 任务执行器
type Executor struct {
}

var (
	G_executor *Executor
)

// 执行一个任务
func (executor *Executor) ExecuteJob(info *common.JobExecuteInfo) {
	go func() {
		var (
			cmd     *exec.Cmd
			result  *common.JobExecuteResult
			jobLock *JobLock
		)

		//任务结果
		result = &common.JobExecuteResult{
			ExecuteInfo: info,
			Output:      make([]byte, 0),
		}

		//初始化分布式锁
		jobLock = worker.G_jobMgr.CreateJobLock(info.Job.Name)

		//记录任务开始时间

		// 上锁
		// 随机睡眠(0~1s)

		//判断上锁是否成功
	}()
}

// 初始化执行器
func InitExecutor() (err error) {
	G_executor = &Executor{}
	return
}
