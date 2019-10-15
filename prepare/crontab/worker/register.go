package worker

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/drzhangg/etcd-test/prepare/crontab/common"
	"net"
	"time"
)

type Register struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease

	localIP string //本机IP
}

var (
	G_register *Register
)

// 获取本地网卡IP
func getLocalIP() (ipv4 string, err error) {
	var (
		addrs   []net.Addr
		addr    net.Addr
		ipNet   *net.IPNet //ip地址
		isIpNet bool
	)

	//获取所有的网卡
	if addrs, err = net.InterfaceAddrs(); err != nil {
		return
	}

	//取第一个非lo的网卡IP
	for _, addr = range addrs {
		//网络地址包括ipv4和ipv6
		if ipNet, isIpNet = addr.(*net.IPNet); isIpNet && !ipNet.IP.IsLoopback() {
			//跳过ipv6
			if ipNet.IP.To4() != nil {
				ipv4 = ipNet.IP.String()
				return
			}
		}
	}

	err = common.ERR_NO_LOCAL_IP_FOUND
	return
}

// 注册到/zhang/cron/workers/IP，并自动续租
func (register *Register) keepOnline() {
	var (
		regKey         string
		leaseGrantResp *clientv3.LeaseGrantResponse
		keepAliveChan  <-chan *clientv3.LeaseKeepAliveResponse
		keepAliveResp  *clientv3.LeaseKeepAliveResponse
		err            error
		cancelCtx      context.Context
		cancelFunc     context.CancelFunc
	)

	//获取注册路径
	regKey = common.JOB_WORKER_DIR + register.localIP
	cancelFunc = nil

	//获取租约
	if leaseGrantResp, err = register.lease.Grant(context.TODO(), 10); err != nil {
		goto RETRY
	}

	//自动续租
	if keepAliveChan, err = register.lease.KeepAlive(context.TODO(), leaseGrantResp.ID); err != nil {
		goto RETRY
	}

	//这个如果不加上会报panic
	cancelCtx, cancelFunc = context.WithCancel(context.TODO())

	//put操作
	if _, err = register.kv.Put(cancelCtx, regKey, "", clientv3.WithLease(leaseGrantResp.ID)); err != nil {
		goto RETRY
	}
	fmt.Println("key----", regKey)

	//通过for查看租约是否时效
	for {
		select {
		case keepAliveResp = <-keepAliveChan:
			//说明续租失败
			if keepAliveResp == nil {
				goto RETRY
			}
		}
	}

	//goto进行错误处理
RETRY:
	time.Sleep(1 * time.Second)
	if cancelFunc != nil {
		cancelFunc()
	}
	fmt.Println("keep online err-----", err)
}

func InitRegister() (err error) {
	var (
		config  clientv3.Config
		client  *clientv3.Client
		kv      clientv3.KV
		lease   clientv3.Lease
		localIP string
	)

	//初始化配置
	config = clientv3.Config{
		Endpoints:   G_config.EtcdEndpoints,
		DialTimeout: time.Duration(G_config.EtcdDialTimeout) * time.Millisecond,
	}

	//建立etcd连接
	if client, err = clientv3.New(config); err != nil {
		return
	}

	//获取本机ip
	if localIP, err = getLocalIP(); err != nil {
		return
	}

	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)

	G_register = &Register{
		client:  client,
		kv:      kv,
		lease:   lease,
		localIP: localIP,
	}

	//服务注册，并自动续租
	go G_register.keepOnline()

	return
}
