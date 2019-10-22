package main

import (
	"fmt"
	"github.com/emicklei/proto"
	"os"
)

type CtrlGenerator struct {
	service *proto.Service
	message []*proto.Message
	rpc     *proto.RPC
}

func (d *CtrlGenerator) Run(opt *Option) (err error) {
	reader, err := os.Open(opt.Proto3Filename)
	if err != nil {
		fmt.Sprintf("open file:%s failed,err: %v\n", opt.Proto3Filename, err)
		return
	}
	defer reader.Close()

	parser := proto.NewParser(reader)
	defintion, err := parser.Parse()
	if err != nil {
		fmt.Sprintf("parse file:%s failed, err:%v\n", opt.Proto3Filename, err)
		return
	}

	proto.Walk(defintion,
		proto.WithService(d.handleService),
		proto.WithMessage(d.handleMessage),
		proto.WithRPC(d.handleRPC),
	)

	fmt.Printf("parse protoc succ, rpc:%v\n", d.rpc)
	return
}

func init() {
	ctrl := &CtrlGenerator{}
	Register("control generator", ctrl)
}

func (d *CtrlGenerator) handleService(s *proto.Service) {
	fmt.Println(s.Name)
	d.service = s
}

func (d *CtrlGenerator) handleMessage(m *proto.Message) {
	fmt.Println(m.Name)
	d.message = append(d.message, m)

}

func (d *CtrlGenerator) handleRPC(r *proto.RPC) {
	fmt.Println(r.Name)
	fmt.Println(r.RequestType)
	fmt.Println(r.ReturnsType)
	fmt.Printf("rpc:%#v, comment:%#v\n", r, r.Comment)
	d.rpc = r
}
