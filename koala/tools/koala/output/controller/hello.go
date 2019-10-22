package controller

import (
	"context"
	"github.com/drzhangg/etcd-test/koala/tools/koala/output/generate"
)

type server struct {
}

func (s *server) SayHello(ctx context.Context, r *hello.DataReq) (resp *hello.DataRsp, err error) {
	return
}
