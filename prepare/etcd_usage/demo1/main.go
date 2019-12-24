package main

import (
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func main() {
	var (
		config clientv3.Config
		client *clientv3.Client
		err    error
	)
	config = clientv3.Config{
		Endpoints:   []string{"47.99.240.51:29"},
		DialTimeout: 5 * time.Second,
	}

	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("连接成功")
	client = client
}
