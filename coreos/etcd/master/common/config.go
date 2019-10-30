package common

import (
	"fmt"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"log"
)

type Config struct {
	Etcd `yaml:"etcd,inline"`
}

type Etcd struct {
	Endpoints   string `yaml:"endpoints"`
	DialTimeout int    `yaml:"dialTimeout"`
}

func InitConfig(fileName string) (err error) {
	var etcdConf Etcd
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	yaml.Unmarshal(bytes, &etcdConf)
	fmt.Println(etcdConf)
	return
}
