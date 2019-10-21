package worker

import (
	"fmt"
	"github.com/drzhangg/etcd-test/prepare/crontab/common"
	"time"
)

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
	var (
		jobSchedulePlan *common.JobSchedulePlan
		//jobExecuting    bool
		jobExisted bool
		err        error
	)
	switch jobEvent.EventType {
	case common.JOB_EVENT_SAVE: //保存任务事件
		if jobSchedulePlan, err = common.BuildJobSchedulePlan(jobEvent.Job); err != nil {
			return
		}
		scheduler.jobPlanTable[jobEvent.Job.Name] = jobSchedulePlan
	case common.JOB_EVENT_DELETE:
		if jobSchedulePlan, jobExisted = scheduler.jobPlanTable[jobEvent.Job.Name]; jobExisted {
			delete(scheduler.jobPlanTable, jobEvent.Job.Name)
		}
	case common.JOB_EVENT_KILL:

	}
}

// 尝试执行任务
func (scheduler *Scheduler) TryStartJob(jobPlan *common.JobSchedulePlan) {
	var (
		jobExecuteInfo *common.JobExecuteInfo
		jobExecuting   bool
	)

	if jobExecuteInfo, jobExecuting = scheduler.jobExecutingTable[jobPlan.Job.Name]; jobExecuting {
		return
	}

	//构建执行状态信息
	jobExecuteInfo = common.BuildJobExecuteInfo(jobPlan)

	//保存执行状态
	scheduler.jobExecutingTable[jobPlan.Job.Name] = jobExecuteInfo

	fmt.Println("执行任务：", jobExecuteInfo.Job.Name, jobExecuteInfo.PlanTime, jobExecuteInfo.RealTime)
	G_executor.ExecuteJob(jobExecuteInfo)

}

// 重新计算任务调度任务
func (scheduler *Scheduler) TryScheduler() (schedulerAfter time.Duration) {
	var (
		now      time.Time
		jobPlan  *common.JobSchedulePlan
		nearTime *time.Time
	)

	//如果任务表为空，直接返回
	if len(scheduler.jobPlanTable) == 0 {
		schedulerAfter = 1 * time.Second
		return
	}

	now = time.Now()

	for _, jobPlan = range scheduler.jobPlanTable {
		if jobPlan.NextTime.Before(now) || jobPlan.NextTime.Equal(now) {
			//
			scheduler.TryStartJob(jobPlan)
			//更新下次执行时间
			jobPlan.NextTime = jobPlan.Expr.Next(now)
		}
		if nearTime == nil || jobPlan.NextTime.Before(*nearTime) {
			nearTime = &jobPlan.NextTime
		}
	}

	schedulerAfter = (*nearTime).Sub(now)
	return
}

// 处理任务结果
func (scheduler *Scheduler) handleJobResult(result *common.JobExecuteResult) {
	var (
		jobLog *common.JobLog
	)
	//删除执行任务
	delete(scheduler.jobExecutingTable, result.ExecuteInfo.Job.Name)

	// 生成执行日志
	if result.Err != common.ERR_LOCAL_ALREADY_REQUIRED {
		jobLog = &common.JobLog{
			JobName:      result.ExecuteInfo.Job.Name,
			Command:      result.ExecuteInfo.Job.Command,
			Output:       string(result.Output),
			PlanTime:     result.ExecuteInfo.PlanTime.UnixNano() / 1000 / 1000,
			ScheduleTime: result.ExecuteInfo.RealTime.UnixNano() / 1000 / 1000,
			StartTime:    result.StartTime.UnixNano() / 1000 / 1000,
			EndTime:      result.EndTime.UnixNano() / 1000 / 1000,
		}

		if result.Err != nil {
			jobLog.Err = result.Err.Error()
		} else {
			jobLog.Err = ""
		}

		G_logSink.Append(jobLog)
	}
}

// 调度协程
func (scheduler *Scheduler) scheduleLoop() {
	var (
		schedulerAfter time.Duration
		schedulerTimer *time.Timer
		jobEvent       *common.JobEvent
		jobResultChan  *common.JobExecuteResult
	)

	//初始化一次
	schedulerAfter = scheduler.TryScheduler()

	//调度的延迟定时器
	schedulerTimer = time.NewTimer(schedulerAfter)

	//定时任务
	for {
		select {
		case jobEvent = <-scheduler.jobEventChan: //监听任务的变化事件，任务的crud
			//对内存中维护的任务列表进行crud操作
			scheduler.handleJobEvent(jobEvent)
		case <-schedulerTimer.C: //最近的任务到期了
		case jobResultChan = <-scheduler.jobResultChan: //监听任务执行结果
			scheduler.handleJobResult(jobResultChan)
		}

		//调度一次任务
		schedulerAfter = scheduler.TryScheduler()
		//重置调度器
		schedulerTimer.Reset(schedulerAfter)
	}

}

//推送任务变化事件
func (scheduler *Scheduler) PushJobEvent(jobEvent *common.JobEvent) {
	scheduler.jobEventChan <- jobEvent
}

//初始化调度器
func InitScheduler() (err error) {
	G_scheduler = &Scheduler{
		jobEventChan:      make(chan *common.JobEvent, 1000),
		jobPlanTable:      make(map[string]*common.JobSchedulePlan),
		jobExecutingTable: make(map[string]*common.JobExecuteInfo),
		jobResultChan:     make(chan *common.JobExecuteResult, 1000),
	}

	//启动调度协程
	go G_scheduler.scheduleLoop()
	return
}

//回传任务执行结果
func (scheduler *Scheduler) PushJobResult(jobResult *common.JobExecuteResult) {
	scheduler.jobResultChan <- jobResult
}
