package main

import (
	"fmt"
	"os"
	"path"
)

//生成目录
type DirGenerator struct {
	dirList []string
}

var AllDirList []string = []string{
	"controller",
	"idl",
	"main",
	"script",
	"conf",
	"app/conf",
	"app/router",
	"model",
	"generate",
}

func (d *DirGenerator) Run(opt *Option) (err error) {

	for _, dir := range AllDirList {
		fullDir := path.Join(opt.Output, dir)
		if err = os.MkdirAll(fullDir, 0755); err != nil {
			fmt.Sprintf("mkdir dir %s failed,err:%v\n", fullDir, err)
			return
		}
	}

	return
}

func init() {
	dir := &DirGenerator{dirList: AllDirList}
	Register("dir generator", dir)
}
