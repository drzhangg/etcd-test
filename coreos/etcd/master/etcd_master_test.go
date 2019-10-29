package master

import "testing"

func TestMaster(t *testing.T) {
	etcdser, _ := RegisterService([]string{"47.99.240.52:2379"})
	etcdser.SetLease(5)
	etcdser.PutService("/drzhang/zhang", "job")
	select {}
}
