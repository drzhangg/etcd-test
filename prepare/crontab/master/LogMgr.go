package master

import (
	"context"
	"github.com/drzhangg/etcd-test/prepare/crontab/common"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
	"time"
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
	var (
		client *mongo.Client
	)
	//建立mongo连接
	if client, err = mongo.Connect(context.TODO(), G_config.MongodbUri, clientopt.ConnectTimeout(time.Duration(G_config.MongodbConnectTimeout)*time.Millisecond)); err != nil {
		return
	}

	G_logMgr = &LogMgr{
		client:        client,
		logCollection: client.Database("cron").Collection("log"),
	}
	return
}

//查看任务日志
func (log *LogMgr) LogList(name string, skip int, limit int) (logArr []*common.JobLog, err error) {
	var (
		jobFilter *common.JobLogFilter
		cursor    mongo.Cursor
		logSort   *common.SortLogByStartTime
		jobLog    *common.JobLog
	)
	logArr = make([]*common.JobLog, 0)
	//查找过滤条件
	jobFilter = &common.JobLogFilter{JobName: name}

	//按照任务开始时间倒排
	logSort = &common.SortLogByStartTime{SortOrder: -1}

	//按照条件过滤查询
	if cursor, err = log.logCollection.Find(context.TODO(), jobFilter, findopt.Sort(logSort), findopt.Skip(int64(skip)), findopt.Limit(int64(limit))); err != nil {
		return
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		jobLog = &common.JobLog{}
		//反序列化bson数据
		if err = cursor.Decode(jobLog); err != nil {
			continue
		}
		logArr = append(logArr, jobLog)
	}

	return
}
