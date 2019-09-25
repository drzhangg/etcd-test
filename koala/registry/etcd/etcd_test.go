package etcd

import (
	"context"
	"fmt"
	"github.com/drzhangg/etcd-test/koala/registry"
	"testing"
	"time"
)

func TestRegister(t *testing.T) {
	registryInst, err := registry.InitRegistry(context.TODO(), "etcd",
		registry.WithAddress([]string{"127.0.0.1:2379"}),
		registry.WithTimeout(time.Second),
		registry.WithRegistryPath("/drzhangg/etcd-test/koala/"),
		registry.WithHearBeat(5),
	)
	if err != nil {
		t.Errorf("init regitry failed,err:%v", err)
		return
	}

	service := &registry.Service{
		Name: "comment_service",
	}

	service.Nodes = append(service.Nodes, &registry.Node{
		Ip:   "127.0.0.1",
		Port: 8081,
	},
		&registry.Node{
			Ip:   "127.0.0.2",
			Port: 8081,
		},
	)
	registryInst.Register(context.TODO(), service)

	go func() {
		time.Sleep(time.Second * 10)
		service.Nodes = append(service.Nodes,
			&registry.Node{Ip: "127.0.0.3", Port: 8081,},
			&registry.Node{Ip: "127.0.0.4", Port: 8081},
		)
		registryInst.Register(context.TODO(), service)
	}()

	for {
		//service, err := registryInst.GetService(context.TODO(), "common_service")
		//if err != nil {
		//	t.Errorf("get service failed,err:%v", err)
		//	return
		//}
		//
		//for _, node := range service.Nodes {
		//	fmt.Printf("service:%s, node:%v\n", service.Name, node)
		//}
		//fmt.Println()

		service, err := registryInst.GetService(context.TODO(), "common_service")
		if err != nil {
			t.Errorf("get service failed,err :%v", err)
			return
		}
		for _, node := range service.Nodes {
			fmt.Printf("service:%s, node:%#v\n", service.Name, node)
		}
		//fmt.Printf("service:%#v\n", service)
		fmt.Println()
		time.Sleep(time.Second * 5)
	}

}
