package master

import (
	"testing"
)

func TestMaster(t *testing.T) {
	etcdser, _ := RegisterService([]string{"47.103.9.218:2379"})
	etcdser.SetLease(5)
	etcdser.PutService("/drzhang/zhang", "job")
	select {}
}

func TestEtcd(t *testing.T) {

}
