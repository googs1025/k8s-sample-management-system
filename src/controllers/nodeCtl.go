package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/shenyisyn/goft-gin/goft"
	"k8s-Management-System/src/models"
	"k8s-Management-System/src/services"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

//@controller
type NodeCtl struct {
	NodeService *services.NodeService `inject:"-"`
	Client *kubernetes.Clientset  `inject:"-"`
}

func NewNodeCtl() *NodeCtl {
	return &NodeCtl{}
}

func(n *NodeCtl) ListAll(c *gin.Context) goft.Json{
	return gin.H{
		"code": 20000,
		"data": n.NodeService.ListAllNodes(),
	}

}

// SaveNode 保存node
func(n *NodeCtl) SaveNode(c *gin.Context) goft.Json{
	nodeModel := &models.PostNodeModel{}
	_ = c.ShouldBindJSON(nodeModel)
	node := n.NodeService.LoadOriginNode(nodeModel.Name) //取出原始node 信息
	if node == nil {
		panic("no such node")
	}
	// 需要用到覆盖，不能直接建立新的。
	node.Labels = nodeModel.OriginLabels  //覆盖标签
	node.Spec.Taints = nodeModel.OriginTaints //覆盖 污点
	_, err := n.Client.CoreV1().Nodes().Update(c, node,v1.UpdateOptions{})

	goft.Error(err)

	return gin.H{
		"code":20000,
		"data":"success",
	}
}

func(n *NodeCtl) LoadDetail(c *gin.Context) goft.Json{
	nodeName := c.Param("node")
	return gin.H{
		"code": 20000,
		"data": n.NodeService.LoadNode(nodeName),
	}

}


func(n *NodeCtl)  Build(goft *goft.Goft){
	goft.Handle("GET","/nodes", n.ListAll)
	goft.Handle("GET","/nodes/:node", n.LoadDetail) //加载详细
	goft.Handle("POST","/nodes", n.SaveNode)  //保存
}

func(*NodeCtl) Name() string{
	return "NodeCtl"
}
