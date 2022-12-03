package models

type IngressModel struct {
	Name string
	NameSpace string
	CreateTime string
}

// path 配置
type IngressPath struct {
	Path string  `json:"path"`
	SvcName string `json:"svc_name"`
	Port string `json:"port"`
}

// 规则
type  IngressRules struct {
	Host string `json:"host"`
	Paths []*IngressPath `json:"paths"`
}

//提交Ingress 对象
type  IngressPost struct{
	Name string
	Namespace string
	Rules []*IngressRules
	Annotations string // 标签
}

