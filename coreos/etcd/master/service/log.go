package service

import (
	"context"
	"github.com/drzhangg/etcd-test/coreos/etcd/common"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
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
func (log *Log) LogList(name string, skip int, limit int) (logArr []Log, err error) {
	var (
		jobFilter *common.JobFilter
	)
	//初始化数组
	logArr = make([]Log, 0)
	//查找过滤条件
	jobFilter = &common.JobFilter{JobName: name}



	return
}
