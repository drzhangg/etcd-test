package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"log"
	"time"
)

func main() {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	//申请一个租约
	leaseGrant, err := client.Grant(context.TODO(), 5)
	if err != nil {
		log.Fatal(err)
	}

	//往etcd里put值
	_, err = client.Put(context.TODO(), "/foo", "bar", clientv3.WithLease(leaseGrant.ID))
	if err != nil {
		log.Fatal(err)
	}

	//续约
	ch, err := client.KeepAlive(context.TODO(), leaseGrant.ID)
	if err != nil {
		log.Fatal(err)
	}

	for {
		ka := <-ch
		fmt.Println("ttl:", ka)
	}
}
