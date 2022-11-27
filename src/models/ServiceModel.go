package models

type Service struct {
	Name string
	NameSpace string  //新增一个命名空间
	Type string
	ClusterIp string
	ClusterIps []string
	Ports []int32
	CreateTime string
}
