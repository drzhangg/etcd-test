package main

import (
	"github.com/drzhangg/etcd-test/koala/tools/koala/output/controller"
	hello "github.com/drzhangg/etcd-test/koala/tools/koala/output/generate"
	"google.golang.org/grpc"
	"log"
	"net"
)

var server = &controller.Server{}

var port = "12345"

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("failed to listen:%!v(MISSING)", err)
	}
	s := grpc.NewServer()
	hello.RegisterHelloServer(s, server)
	s.Serve(lis)
}
