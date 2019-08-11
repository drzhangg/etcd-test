package demo6

import "github.com/drzhangg/etcd-test/prepare/etcd_usage/demo5"

type SonController struct {
	BaseController demo5.Controller
}

func main() {
	son := SonController{}
	son.BaseController.BaseController()
}
