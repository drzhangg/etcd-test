package main

var main_template = `package main
import(
	"net"
	"log"
	"google.golang.org/grpc"
	"github.com/drzhangg/etcd-test/koala/tools/koala/output/controller"

	pb "github.com/drzhangg/etcd-test/koala/tools/koala/output/generate"
)

var server = &controller.Server{}
var port = "12345"
		
func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("failed to listen:%v", err)
	}
	s := grpc.NewServer()
	pb.RegisterHelloServer(s, server)
	s.Serve(lis)
}

`
