package master

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

/**
1.注册服务
2.保存kv数据
3.创建租约
4.租约过期进行续租
*/

type EtcdService struct {
	client             *clientv3.Client
	lease              clientv3.Lease
	leaseResp          *clientv3.LeaseGrantResponse
	leaseKeepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
}

// 注册etcd
func RegisterService(addr []string) (*EtcdService, error) {
	var (
		config clientv3.Config
		client *clientv3.Client
		err    error
	)
	config = clientv3.Config{
		Endpoints:   addr,
		DialTimeout: 5 * time.Second,
	}

	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &EtcdService{client: client}, nil
}

// 设置租约,并返回租约id和租约时间
func (this *EtcdService) SetLease(timeNum int64) error {
	var (
		lease     clientv3.Lease
		leaseResp *clientv3.LeaseGrantResponse
		err       error
	)
	//通过client获取租约
	lease = clientv3.NewLease(this.client)

	//设置租约时间
	if leaseResp, err = lease.Grant(context.TODO(), timeNum); err != nil {
		return err
	}

	//ctx, cancel := context.WithCancel(context.TODO())
	leaseRespchan, err := lease.KeepAlive(context.TODO(), leaseResp.ID)
	if err != nil {
		return err
	}

	this.lease = lease
	this.leaseResp = leaseResp
	this.leaseKeepAliveChan = leaseRespchan

	return nil
}

// 通过租约 注册服务
func (this *EtcdService) PutService(key, val string) {
	var (
		kv      clientv3.KV
		err     error
		putResp *clientv3.PutResponse
	)
	fmt.Println("k and v ", key, val)
	kv = clientv3.NewKV(this.client)
	if putResp, err = kv.Put(context.TODO(), key, val, clientv3.WithLease(this.leaseResp.ID)); err != nil {
		fmt.Printf("etcd set value failed! key:%s;value:%s", key, val)
	}
	fmt.Println(putResp)
	if putResp.PrevKv != nil {
		fmt.Println("", putResp.PrevKv.Key, putResp.PrevKv.Value)
	}

	//putResp = putResp
}

// ListenLeaseRespChan 监听 续租情况
func (this *EtcdService) ListenLeaseRespChan() {
	for {
		select {
		case leaseRespKeep := <-this.leaseKeepAliveChan:
			if leaseRespKeep == nil { //这里说明续租失败
				fmt.Printf("已经关闭续租功能\n")
				return
			} else {
				fmt.Printf("续租成功\n")
			}
		}
	}
}

// 撤销租约
func (this *EtcdService) RevokeLease() error {
	time.Sleep(2 * time.Second)
	_, err := this.lease.Revoke(context.TODO(), this.leaseResp.ID)
	return err
}
