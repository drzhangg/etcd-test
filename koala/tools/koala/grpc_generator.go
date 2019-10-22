package main

import (
	"fmt"
	"os"
	"os/exec"
)

type GrpcGenerator struct {
}

func (d *GrpcGenerator) Run(opt *Option) (err error) {
	// protoc --go_out=plugins=grpc:. hello.proto
	outputName := fmt.Sprintf("plugins=grpc:%s/generate/", opt.Output)
	cmd := exec.Command("protoc", "--go_out", outputName, opt.Proto3Filename)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err = cmd.Run(); err != nil {
		fmt.Sprintf("grpc generator failed, err:%v\n", err)
		return
	}
	return
}

func init() {
	grpcRegister := &GrpcGenerator{}

	Register("grpc generator", grpcRegister)
}
