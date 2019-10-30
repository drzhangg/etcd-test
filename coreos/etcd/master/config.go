package master

import (
	"fmt"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"log"
)

type EtcdConf struct {
	Etcd Etcd `yaml:"etcd"`
}

type Etcd struct {
	Endpoints   string `yaml:"endpoints"`
	DialTimeout int    `yaml:"dialTimeout"`
}

func InitConfig(fileName string) (err error) {
	var etcdConf EtcdConf
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(bytes, &etcdConf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(etcdConf)
	return
}
