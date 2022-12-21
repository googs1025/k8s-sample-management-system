package models

type NodesConfig struct {
	Name string
	Ip string
	User string
	Pass string
}

type K8sConfig struct {
	Nodes []*NodesConfig
}

type SysConfig struct {
	K8s *K8sConfig
}
