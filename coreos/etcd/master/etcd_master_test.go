package master

import "testing"

func TestMaster(t *testing.T) {
	etcdser, _ := RegisterService([]string{"127.0.0.1:2379"})
	etcdser.SetLease(5)
	etcdser.PutService("/drzhang/zhang", "job")
	select {}
}
