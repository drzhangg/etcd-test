package main

import (
	"flag"
	"fmt"
	"github.com/drzhangg/etcd-test/prepare/crontab/master"
	"runtime"
)

var (
	confFile string //配置文件路径
)

// 解析命令行参数
func initArgs() {
	flag.StringVar(&confFile, "config", "./master.json", "指定master.json")
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
	if err = master.InitConfig(confFile); err != nil {
		goto ERR
	}

	//初始化服务发现模块
	if err = master.InitWorkerMgr(); err != nil {
		goto ERR
	}

	//日志管理器
	if err = master.InitLogMgr(); err != nil {
		goto ERR
	}

	//任务管理器
	master.InitJobMgr()

	//启动api HTTP服务
	master.InitApiServer()

	//正常退出

ERR:
	fmt.Println(err)
}
