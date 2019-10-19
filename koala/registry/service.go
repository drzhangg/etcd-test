package registry

// 服务注册名称和服务注册节点
type Service struct {
	Name  string  `json:"name"`
	Nodes []*Node `json:"nodes"`
}

//服务注册节点信息
type Node struct {
	Id     string `json:"id"`
	Ip     string `json:"ip"`
	Port   int    `json:"port"`
	Weight int    `json:"weight"` //权重
}
