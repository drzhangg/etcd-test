package main

import "github.com/emicklei/proto"

type Generator interface {
	Run(opt *Option,metaData *ServiceMetaData) error
}

type ServiceMetaData struct {
	service  *proto.Service
	messages []*proto.Message
	rpc      []*proto.RPC
}
