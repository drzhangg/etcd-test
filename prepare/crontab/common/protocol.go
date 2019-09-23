package common

import (
	"github.com/gorhill/cronexpr"
	"time"
)

//定时任务
type Job struct {
	Name     string `json:"name"`     //任务名
	Command  string `json:"command"`  //shell命令
	CronExpr string `json:"cronExpr"` //cron表达式
}

//任务调度计划
type JobSchedulePlan struct {
	Job      *Job                 //要调度的信息
	Expr     *cronexpr.Expression //解析好的cronexpr表达式
	NextTime time.Time            //下次调度时间
}
