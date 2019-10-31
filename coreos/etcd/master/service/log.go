package service

import (
	"context"
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
// skip
func (log *Log) LogList(name string, skip int, limit int) (logArr []Log, err error) {
	var ()
	return
}
