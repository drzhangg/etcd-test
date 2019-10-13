package worker

import "github.com/drzhangg/etcd-test/prepare/crontab/common"

// 任务调度
type Scheduler struct {
	jobEventChan      chan *common.JobEvent              //etcd任务事件队列
	jobPlanTable      map[string]*common.JobSchedulePlan //任务调度计划表
	jobExecutingTable map[string]*common.JobExecuteInfo  //任务执行表
	jobResultChan     chan *common.JobExecuteResult      //任务结果队列
}

var (
	G_scheduler *Scheduler
)

// 处理任务事件
func (scheduler *Scheduler) handleJobEvent(jobEvent *common.JobEvent) {
	var ()
	switch jobEvent.EventType {
	case common.JOB_EVENT_SAVE:
		common.BuildResponse()
	case common.JOB_EVENT_DELETE:
	case common.JOB_EVENT_KILL:

	}
}

// 尝试执行任务
func (scheduler *Scheduler) TryStartJob(jobPlan *common.JobSchedulePlan) {

}

// 重新计算任务调度任务

// 处理任务结果

// 调度协程

//推送任务变化事件

//初始化调度器

//回传任务执行结果
