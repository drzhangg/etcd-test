package main

import "fmt"

var genMgr *GeneratorMgr = &GeneratorMgr{
	genMap: map[string]Generator{},
}

type GeneratorMgr struct {
	genMap map[string]Generator
}

func (gen *GeneratorMgr) Run(opt *Option) (err error) {
	for _, gen := range gen.genMap {
		if err = gen.Run(opt); err != nil {
			return
		}
	}
	return
}

func Register(name string, gen Generator) (err error) {
	_, ok := genMgr.genMap[name]
	if ok {
		err = fmt.Errorf("generator %s is exists", name)
		return
	}
	genMgr.genMap[name] = gen
	return
}
