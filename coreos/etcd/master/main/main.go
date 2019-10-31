package main

import (
	"flag"
	"fmt"
	"github.com/drzhangg/etcd-test/coreos/etcd/master/service"
)

var (
	confFile string
)

// 解析命令行参数
func initArgs() {
	flag.StringVar(&confFile, "config", "./master.yaml", "指定master.yaml")
	flag.Parse()
}

func main() {
	var (
		err error
	)

	initArgs()

	//初始化配置文件
	if err = service.InitConfig(confFile); err != nil {
		goto ERR
	}

	//初始化etcd
	if err = service.InitWork(); err != nil {
		goto ERR
	}

	//初始化日志文件


	//初始化任务管理器，对etcd进行操作
ERR:
	fmt.Println(err)
}
