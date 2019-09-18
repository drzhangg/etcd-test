package registry

import (
	"fmt"
	"sync"
)

var (
	pluginMgr = &PluginMgr{
		plugins: make(map[string]Registry),
	}
)

type PluginMgr struct {
	plugins map[string]Registry
	lock    sync.Mutex //是为了进行多线程操作
}

func (p *PluginMgr) registerPlugin(plugin Registry) (err error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	_, ok := p.plugins[plugin.Name()]
	if ok {
		err = fmt.Errorf("duplicat registry plugin")
		return
	}
	p.plugins[plugin.Name()] = plugin
	return
}
