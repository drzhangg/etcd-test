package master

import (
	"github.com/mongodb/mongo-go-driver/mongo"
)

// mongodb 日志管理
type LogMgr struct {
	client        *mongo.Client     //mongo连接
	logCollection *mongo.Collection //mongo数据库
}

var (
	G_logMgr *LogMgr
)

func InitLogMgr() (err error) {
	return
}

//查看任务日志
