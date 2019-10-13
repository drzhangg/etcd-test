package worker

import "github.com/drzhangg/etcd-test/prepare/crontab/common"

// 任务执行器
type Executor struct {

}

var (
	G_executor *Executor
)

func (executor *Executor) ExecuteJob(info *common.JobExecuteInfo) {

}

func InitExecutor() (err error) {
	return
}
