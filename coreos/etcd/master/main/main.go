package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"log"
	"time"
)

func main() {
	config := clientv3.Config{
		Endpoints:   []string{"47.103.9.218:2379"},
		DialTimeout: time.Duration(time.Second * 5),
	}

	client, err := clientv3.New(config)
	if err != nil {
		log.Fatal(err)
	}

	kv := clientv3.NewKV(client)
	getResp, err := kv.Get(context.TODO(), "/test/key1")
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range getResp.Kvs {
		fmt.Println(string(v.Key), string(v.Value))
	}

}
