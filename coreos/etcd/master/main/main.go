package main

import (
	"flag"
	"fmt"
	"github.com/drzhangg/etcd-test/coreos/etcd/master/common"
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

	initArgs()

	//初始化配置文件
	if err := common.InitConfig(confFile); err != nil {
		fmt.Println(err)
	}

	//初始化etcd


}