package worker

import (
	"context"
	"github.com/drzhangg/etcd-test/prepare/crontab/common"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
	"time"
)

type LogSink struct {
	client         *mongo.Client
	logCollection  *mongo.Collection
	logChan        chan *common.JobLog   //日志信息
	autoCommitChan chan *common.LogBatch //日志批次
}

var (
	G_logSink *LogSink
)

// 批量写入日志
func (logSink *LogSink) saveLogs(batch *common.LogBatch) {
	logSink.logCollection.InsertMany(context.TODO(), batch.Logs)
	return
}

// 发送日志
func (logSink *LogSink) Append(jobLog *common.JobLog) {
	select {
	//把jobLog传入通道
	case logSink.logChan <- jobLog:
	default:
		//队列满了就丢弃
	}
	return
}

// 日志存储协程
func (logSink *LogSink) writeLoop() {
	var (
		log          *common.JobLog
		logBatch     *common.LogBatch //当前批次
		timeOutBatch *common.LogBatch //超时批次
		commitTimer  *time.Timer
	)
	for {
		select {
		case log = <-logSink.logChan:
			if logBatch == nil {
				logBatch = &common.LogBatch{}
				commitTimer = time.AfterFunc(time.Duration(G_config.JobLogBatchSize)*time.Millisecond, func(batch *common.LogBatch) func() {
					return func() {
						logSink.autoCommitChan <- batch
					}
				}(logBatch),
				)
			}

			//把新日志追加到批次中
			logBatch.Logs = append(logBatch.Logs, log)

			//如果批次满了就立刻保存
			if len(logBatch.Logs) >= G_config.JobLogBatchSize {
				//保存日志
				logSink.saveLogs(logBatch)
				//清空logBatch
				logBatch = nil
				//取消定时器
				commitTimer.Stop()
			}
		case timeOutBatch = <-logSink.autoCommitChan: //过期的日志批次
			//判断过期批次是否仍旧是当前的批次
			if timeOutBatch != logBatch {
				continue //跳过已经提交过的批次
			}
			//保存未提交的批次
			logSink.saveLogs(timeOutBatch)
			//把提交过的批次清空
			logBatch = nil
		}
	}
}

func InitLogSink() (err error) {
	var (
		client *mongo.Client
	)

	//建立mongo连接
	if client, err = mongo.Connect(context.TODO(), G_config.MongodbUri, clientopt.ConnectTimeout(time.Duration(G_config.MongodbConnectTimeout)*time.Millisecond)); err != nil {
		return
	}

	//选择db和collection
	G_logSink = &LogSink{
		client:         client,
		logCollection:  client.Database("cron").Collection("log"),
		logChan:        make(chan *common.JobLog, 1000),
		autoCommitChan: make(chan *common.LogBatch, 1000),
	}

	//启动一个协程，来处理日志
	go G_logSink.writeLoop()

	return
}
