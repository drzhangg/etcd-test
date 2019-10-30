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

	if err := common.InitConfig(confFile); err != nil {
		fmt.Println(err)
	}

}
