package main

import (
	"flag"
	"fmt"
	"github.com/drzhangg/etcd-test/prepare/crontab/worker"
	"runtime"
)

var (
	confFile string //配置文件路径
)

// 解析命令行参数
func initArgs() {
	flag.StringVar(&confFile, "config", "./worker.json", "worker.json")
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

	//服务注册

	//启动日志协程

	//启动执行器

	//启动调度器

	//初始化任务管理器

	//正常退出
ERR:
	fmt.Println(err)
}
