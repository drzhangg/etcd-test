package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func main() {
	var (
		config  clientv3.Config
		client  *clientv3.Client
		kv      clientv3.KV
		err     error
		getResp *clientv3.GetResponse
	)

	config = clientv3.Config{
		Endpoints:   []string{"47.99.240.52:2379"},
		DialTimeout: 5 * time.Second,
	}

	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	//对etcd进行操作的kv
	kv = clientv3.NewKV(client)

	go func() {
		kv.Put(context.TODO(), "/cron/lock/job2", "job2")

		kv.Delete(context.TODO(), "/cron/lock/job2")

		time.Sleep(1 * time.Second)
	}()

	if getResp, err = kv.Get(context.TODO(), "/cron/lock/job2"); err != nil {
		fmt.Println(err)
		return
	}

	if len(getResp.Kvs) != 0 {
		fmt.Println("value:", string(getResp.Kvs[0].Value))
	}

}
