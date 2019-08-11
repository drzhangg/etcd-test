package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func main() {
	var (
		config    clientv3.Config
		client    *clientv3.Client
		err       error
		lease     clientv3.Lease
		leaseResp *clientv3.LeaseGrantResponse
		leaseId   clientv3.LeaseID
		kv        clientv3.KV
		putResp   *clientv3.PutResponse
		getResp   *clientv3.GetResponse
	)
	config = clientv3.Config{
		Endpoints:   []string{"47.99.240.52:2379"},
		DialTimeout: 5 * time.Second,
	}

	client, err = clientv3.New(config)
	if err != nil {
		fmt.Println(err)
		return
	}

	//申请一个租约
	lease = clientv3.NewLease(client)
	//设置租约时间
	if leaseResp, err = lease.Grant(context.TODO(), 10); err != nil {
		fmt.Println(err)
		return
	}

	//获取的租约id,用于关联进行操作
	leaseId = leaseResp.ID

	kv = clientv3.NewKV(client)

	//对etcd通过关联的租约进行put操作
	if putResp, err = kv.Put(context.TODO(), "/cron/lock/job1", "lock1", clientv3.WithLease(leaseId)); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("写入成功：", putResp.Header.Revision)

	for {
		if getResp, err = kv.Get(context.TODO(), "/cron/lock/job1"); err != nil {
			fmt.Println(err)
			return
		}

		if getResp.Count == 0 {
			fmt.Println("kv过期了")
			break
		}
		fmt.Println("还没过期：", getResp.Kvs)

		time.Sleep(2 * time.Second)
	}

	//建立一个etcd客户端后要申请一个lease租约

	//申请租约后	给租约加时

	//想要对租约进行操作要获取到租约的id

	//对kv进行操作，要与租约进行关联
}
