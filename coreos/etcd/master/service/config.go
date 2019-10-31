package service

import (
	"fmt"
	"github.com/ghodss/yaml"
	"io/ioutil"
)

type Config struct {
	Etcd    Etcd    `yaml:"etcd"`
	MongoDb MongoDB `yaml:"mongodb"`
}

type Etcd struct {
	Endpoints   []string `yaml:"endpoints"`
	DialTimeout int64      `yaml:"dialTimeout"`
}

type MongoDB struct {
	Url            string `yaml:"url"`
	ConnectTimeout int64    `yaml:"connectTimeout"`
}

var (
	//单例模式，全局变量
	G_config *Config
)

func InitConfig(fileName string) (err error) {
	var (
		config Config
		bytes  []byte
	)

	if bytes, err = ioutil.ReadFile(fileName); err != nil {
		return
	}

	if err = yaml.Unmarshal(bytes, &config); err != nil {
		return
	}
	fmt.Println(config)
	G_config = &config
	return
}
