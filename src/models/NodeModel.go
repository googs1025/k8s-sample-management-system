package models
type NodeModel struct {
	Name string
	IP string
	HostName string
	Lables []string // 标签
	Taints []string // 污点
	Capacity *NodeCapacity  //容量 包含了cpu 内存和pods数量
	Usage *NodeUsage //资源 使用情况
	CreateTime string
}

type NodeUsage struct {
	Pods int
	Cpu float64
	Memory float64
}

func NewNodeUsage(pods int, cpu float64, memory float64) *NodeUsage {
	return &NodeUsage{Pods: pods, Cpu: cpu, Memory: memory}
}

//容量
type NodeCapacity struct {
	Cpu int64
	Memory int64
	Pods int64
}

func NewNodeCapacity(cpu int64, memory int64, pods int64) *NodeCapacity {
	return &NodeCapacity{Cpu: cpu, Memory: memory, Pods: pods}
}
