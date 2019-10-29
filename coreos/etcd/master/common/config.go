package common

type Etcd struct {
	Endpoints []*string`yaml:endpoints`
}
