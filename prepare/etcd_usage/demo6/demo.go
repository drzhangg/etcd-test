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
		config             clientv3.Config
		client             *clientv3.Client
		kv                 clientv3.KV
		err                error
		getResp            *clientv3.GetResponse
		watchStartRevision int64
		watcher            clientv3.Watcher
		watchChan          clientv3.WatchChan
		watchResp          clientv3.WatchResponse
		event              *clientv3.Event
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

	//当前etcd集群事务id revision，单调递增的
	watchStartRevision = getResp.Header.Revision + 1
	fmt.Println("从该版本向后监听：", watchStartRevision)

	//创建一个watcher
	watcher = clientv3.NewWatcher(client)

	//启动watcher监听,clientv3.WithRev从哪个版本开始监听
	watchChan = watcher.Watch(context.TODO(), "/cron/lock/job2", clientv3.WithRev(watchStartRevision))

	//处理kv变化事件
	for watchResp = range watchChan {
		for _, event = range watchResp.Events {
			switch event.Type {
			case mvccpb.PUT:
				fmt.Println("修改为：", string(event.Kv.Value), "Revision:", event.Kv.Version, event.Kv.ModRevision)
			case mvccpb.DELETE:
				fmt.Println("删除了，Revision:", event.Kv.ModRevision)
			}
		}
	}
}
