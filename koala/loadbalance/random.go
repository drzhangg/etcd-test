package loadbalance

import (
	"context"
	"github.com/drzhangg/etcd-test/koala/registry"
	"math/rand"
)

type RandomBalance struct {
}

func (r *RandomBalance) Name() string {
	return "random"
}

func (r *RandomBalance) Select(ctx context.Context, nodes []*registry.Node) (node *registry.Node, err error) {

	if len(nodes) == 0 {
		return
	}
	index := rand.Intn(len(nodes))
	node = nodes[index]
	return
}
