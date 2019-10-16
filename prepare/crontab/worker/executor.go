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

// 初始化执行器
func InitExecutor() (err error) {
	G_executor = &Executor{}
	return
}
