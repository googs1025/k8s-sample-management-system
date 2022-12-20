package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/shenyisyn/goft-gin/goft"
	"k8s-Management-System/src/services"
)

//@controller
type NodeCtl struct {
	NodeService *services.NodeService `inject:"-"`
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

func(n *NodeCtl)  Build(goft *goft.Goft){
	goft.Handle("GET","/nodes", n.ListAll)
}

func(*NodeCtl) Name() string{
	return "NodeCtl"
}
