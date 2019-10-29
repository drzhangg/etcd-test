package main

import (
	"fmt"
	"html/template"
	"os"
	"path"
)

type MainGenerator struct {
}

func (d *MainGenerator) Run(opt *Option,metaData *ServiceMetaData) (err error) {
	filename := path.Join("./", opt.Output, "main/main.go")
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Sprintf("open file:%s failed,err: %v\n", filename, err)
		return
	}
	defer file.Close()

	err = d.render(file, main_template)
	if err != nil {
		return
	}

	return
}

func (d *MainGenerator) render(file *os.File, data string) (err error) {
	t := template.New("main")
	t, err = t.Parse(data)
	if err != nil {
		return
	}
	err = t.Execute(file, nil)
	return
}

func (d *MainGenerator) generateRpc(opt *Option) (err error) {

	return
}

func init() {

	mainGenerator := &MainGenerator{}

	Register("main generator", mainGenerator)
}
