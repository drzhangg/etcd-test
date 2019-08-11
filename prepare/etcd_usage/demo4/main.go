package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"

	"time"
)

func main() {
	var (
		config  clientv3.Config
		err     error
		client  *clientv3.Client
		kv      clientv3.KV
		delResp *clientv3.DeleteResponse
		prevKvs *mvccpb.KeyValue
	)

	config = clientv3.Config{
		Endpoints:   []string{"47.99.240.52:2379"},
		DialTimeout: 5 * time.Second,
	}

	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	//创建用于操作etcd的kv
	kv = clientv3.NewKV(client)

	if delResp, err = kv.Delete(context.TODO(), "/cron/jobs/job1",clientv3.WithPrevKV()); err != nil {
		fmt.Println(err)
	}

	if len(delResp.PrevKvs) != 0 {
		for _, prevKvs = range delResp.PrevKvs {
			fmt.Println("删除了：", string(prevKvs.Key), string(prevKvs.Value))
		}
	}

}
