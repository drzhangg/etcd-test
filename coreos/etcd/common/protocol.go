package common

import "strings"

//任务执行日志
type JobLog struct {
	JobName      string `json:"jobName" bson:"jobName"`           //任务执行名称
	Command      string `json:"command" bson:"command"`           //脚本执行命令
	Err          string `json:"err" bson:"err"`                   //err
	Output       string `json:"output" bson:"output"`             //输出脚本
	PlanTime     int64  `json:"planTime" bson:"planTime"`         //计划开始时间
	ScheduleTime int64  `json:"scheduleTime" bson:"scheduleTime"` //实际调度时间
	StartTime    int64  `json:"startTime" bson:"startTime"`       //任务开始执行时间
	EndTime      int64  `json:"endTime" bson:"endTime"`           //任务结束执行时间
}

//日志排序规则	1：正排，-1：倒排
type SortLogByStartTime struct {
	SortOrder int `bson:"startTime"`
}

// 任务日志过滤条件
type JobFilter struct {
	JobName string `bson:"jobName"`
}

//截取key后面的ip
func ExtractWorkerIP(regKey string) string {
	return strings.TrimPrefix(regKey, JOB_WORKER_DIR)
}
