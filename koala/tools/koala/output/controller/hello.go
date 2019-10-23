package controller
import(
"context"
"github.com/drzhangg/etcd-test/koala/tools/koala/output/generate"
)
type Server struct{
}


func (s *Server) SayHello(ctx context.Context,r *hello.DataReq)(resp *hello.DataRsp,err error){
return}

