package worker

import (
	"go.etcd.io/etcd/clientv3"
	"time"
)

/**
服务发现步骤：
1.watch节点
2.监听对节点进行了哪些操作
*/

type ClientDis struct {
	client *clientv3.Client
}

func NewService(addr []string) (*ClientDis, error) {
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
		return nil, err
	}
	return &ClientDis{client: client}, nil
}

// 获取
//func (cli *ClientDis) extractAddrs(resp *clientv3.GetResponse) []string {
//	addrs := make([]string, 0)
//	if resp == nil || resp.Kvs == nil {
//		return addrs
//	}
//	for i := range resp.Kvs {
//		//resp.Kvs[i].Value
//	}
//}
