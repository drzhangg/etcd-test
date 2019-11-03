package common

import "strings"

// 任务日志过滤条件
type JobFilter struct {
	JobName string `bson:"jobName"`
}

//截取key后面的ip
func ExtractWorkerIP(regKey string) string {
	return strings.TrimPrefix(regKey, JOB_WORKER_DIR)
}
