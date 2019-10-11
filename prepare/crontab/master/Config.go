package master

import (
	"encoding/json"
	"io/ioutil"
)

// 程序配置
type Config struct {
	ApiPort               int      `json:"apiPort"`
	ApiReadTimeout        int      `json:"apiReadTimeout"`
	ApiWriteTimeout       int      `json:"apiWriteTimeout"`
	EtcdEndpoints         []string `json:"etcdEndpoints"`
	EtcdDialTimeout       int      `json:"etcdDialTimeout"`
	WebRoot               string   `json:"webRoot"`
	MongodbUri            string   `json:"mongodbUri"`
	MongodbConnectTimeout int      `json:"mongodbConnectTimeout"`
}

var (
	//单例模式
	G_config *Config
	config   Config
)

//加载配置
func InitConfig(filename string) (err error) {
	var (
		content []byte
	)
	//1.读取配置文件
	if content, err = ioutil.ReadFile(filename); err != nil {
		return
	}
	//2.反序列化
	if err = json.Unmarshal(content, &config); err != nil {
		return
	}

	//3.将反序列化后的值传给单例
	G_config = &config
	return
}
