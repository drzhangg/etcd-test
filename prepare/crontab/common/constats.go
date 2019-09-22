package common

const (
	//任务保存目录
	JOB_SAVE_DIR = "/zhang/cron/jobs/"

	//任务强杀目录
	JOB_KILLER_DIR = "/zhang/cron/killer/"

	//任务锁目录
	JOB_LOCK_DIR = "/zhang/cron/lock/"

	// 服务注册目录
	JOB_WORKER_DIR = "/zhang/cron/workers/"

	//保存任务事件
	JOB_EVENT_SAVE = 1

	//删除任务事件
	JOB_EVENT_DELETE = 2

	//强杀任务事件
	JOB_EVENT_KILL = 3
)
