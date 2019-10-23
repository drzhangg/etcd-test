package main

import (
	"fmt"
	"github.com/emicklei/proto"
	"os"
	"path"
)

type CtrlGenerator struct {
	service *proto.Service
	message []*proto.Message
	rpc     []*proto.RPC
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

	//fmt.Printf("parse protoc succ, rpc:%v\n", d.rpc)
	return d.generateRpc(opt)
}

func (d *CtrlGenerator) generateRpc(opt *Option) (err error) {
	fileName := path.Join("./", opt.Output, "controller", fmt.Sprintf("%s.go", d.service.Name))
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Sprintf("open file:%s failed, err:%v\n", fileName, err)
		return
	}
	defer file.Close()
	fmt.Fprintf(file, "package controller\n")
	fmt.Fprintf(file, "import(\n")
	fmt.Fprintf(file, `"context"`)
	fmt.Fprintln(file)
	fmt.Fprintf(file, `"github.com/drzhangg/etcd-test/koala/tools/koala/output/generate"`)
	fmt.Fprintln(file)
	fmt.Fprintf(file, ")\n")
	fmt.Fprintf(file, "type Server struct{\n}\n")
	fmt.Fprintf(file, "\n\n")

	for _, rpc := range d.rpc {
		fmt.Fprintf(file, "func (s *Server) %s(ctx context.Context,r *hello.%s)(resp *hello.%s,err error){\nreturn}\n\n", rpc.Name, rpc.RequestType, rpc.ReturnsType)
	}
	return
}

func init() {
	ctrl := &CtrlGenerator{}
	Register("control generator", ctrl)
}

func (d *CtrlGenerator) handleService(s *proto.Service) {
	//fmt.Println(s.Name)
	d.service = s
}

func (d *CtrlGenerator) handleMessage(m *proto.Message) {
	//fmt.Println(m.Name)
	d.message = append(d.message, m)

}

func (d *CtrlGenerator) handleRPC(r *proto.RPC) {
	d.rpc = append(d.rpc, r)
}
