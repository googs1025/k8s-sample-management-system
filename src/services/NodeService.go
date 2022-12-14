package services

import (
	"k8s-Management-System/src/helpers"
	"k8s-Management-System/src/models"
	"k8s.io/api/core/v1"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

// @service
type NodeService struct {
	NodeMap *NodeMap `inject:"-"`
	PodMap *PodMap `inject:"-"`
	Metric *versioned.Clientset `inject:"-"`
}

func NewNodeService() *NodeService {
	return &NodeService{}
}

// ListAllNodes 显示所有节点
func (ns *NodeService) ListAllNodes() []*models.NodeModel{
	list := ns.NodeMap.ListAll()
	ret := make([]*models.NodeModel, len(list))
	for i, node := range list {
		nodeUsage:=helpers.GetNodeUsage(ns.Metric,node)
		ret[i] = &models.NodeModel{
			Name: node.Name,
			IP: node.Status.Addresses[0].Address,
			Labels: helpers.FilterLabels(node.Labels),
			Taints: helpers.FilterTaints(node.Spec.Taints),
			HostName: node.Status.Addresses[1].Address,
			Capacity: models.NewNodeCapacity(node.Status.Capacity.Cpu().Value(),
				node.Status.Capacity.Memory().Value(),node.Status.Capacity.Pods().Value()),
			Usage:models.NewNodeUsage(ns.PodMap.GetNum(node.Name), nodeUsage[0], nodeUsage[1]),
			CreateTime: node.CreationTimestamp.Format("2006-01-02 15:04:05"),
		}
	}
	return ret
}


// LoadOriginNode 保存时用的
func(ns *NodeService) LoadOriginNode(nodeName string ) *v1.Node{
	return ns.NodeMap.Get(nodeName)
}

// LoadNode 加载node信息，给编辑用的
func(ns *NodeService) LoadNode(nodeName string ) *models.NodeModel{
	node := ns.NodeMap.Get(nodeName)
	return &models.NodeModel{
		Name: node.Name,
		IP: node.Status.Addresses[0].Address,
		HostName: node.Status.Addresses[1].Address,
		OriginLabels: node.Labels,
		OriginTaints: node.Spec.Taints,
		CreateTime: node.CreationTimestamp.Format("2006-01-02 15:04:05"),
	}
}