package main

import (
	"flag"
	"fmt"
	"github.com/drzhangg/etcd-test/prepare/crontab/worker"
	"runtime"
	"time"
)

var (
	confFile string //配置文件路径
)

// 解析命令行参数
func initArgs() {
	flag.StringVar(&confFile, "config", "/Users/drzhang/go/src/github.com/drzhangg/etcd-test/prepare/crontab/worker/main/worker.json", "worker.json")
	flag.Parse()
}

// 初始化线程数量
func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		err error
	)

	//初始化命令行参数
	initArgs()

	//初始化线程
	initEnv()

	//加载配置
	if err = worker.InitConfig(confFile); err != nil {
		goto ERR
	}

	//进行etcd服务注册
	if err = worker.InitRegister(); err != nil {
		goto ERR
	}

	//启动日志协程
	if err = worker.InitLogSink(); err != nil {
		goto ERR
	}

	//启动执行器

	//启动调度器

	//初始化任务管理器

	//正常退出
	for {
		time.Sleep(1 * time.Second)
	}
	return
ERR:
	fmt.Println("init err------",err)
}
