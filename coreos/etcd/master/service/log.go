package service

import (
	"context"
	"github.com/drzhangg/etcd-test/coreos/etcd/common"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
	"time"
)

type Log struct {
	client        *mongo.Client
	logCollection *mongo.Collection
}

var (
	G_log *Log
)

func InitLog() (err error) {
	var (
		client *mongo.Client
	)

	if client, err = mongo.Connect(context.TODO(), G_config.MongoDb.Url, clientopt.ConnectTimeout(time.Duration(G_config.MongoDb.ConnectTimeout)*time.Second)); err != nil {
		return
	}

	G_log = &Log{
		client:        client,
		logCollection: client.Database("cron").Collection("log"),
	}
	return
}

// 获取日志列表
// name 过滤条件
// skip	跳过指定数量的数据
// limit	读取指定数量的数据
func (log *Log) LogList(name string, skip int, limit int) (logArr []*common.JobLog, err error) {
	var (
		jobFilter *common.JobFilter
		logSort   *common.SortLogByStartTime
		cursor    mongo.Cursor
		jobLog    *common.JobLog
	)
	//初始化数组
	logArr = make([]*common.JobLog, 0)
	//查找过滤条件
	jobFilter = &common.JobFilter{JobName: name}

	//按照任务开始时间倒排
	logSort = &common.SortLogByStartTime{SortOrder: -1}

	//按照条件过滤查询
	if cursor, err = log.logCollection.Find(context.TODO(), jobFilter, findopt.Sort(logSort), findopt.Skip(int64(skip)), findopt.Limit(int64(limit))); err != nil {
		return
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		jobLog = &common.JobLog{}
		if err = cursor.Decode(jobLog); err != nil {
			continue
		}
		logArr = append(logArr, jobLog)
	}
	return
}
