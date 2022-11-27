package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/shenyisyn/goft-gin/goft"
	"k8s-Management-System/src/services"
)

type NamespaceCtl struct {
	NamespaceMap *services.NamespaceMap `inject:"-"`
}

func NewNamespaceCtl() *NamespaceCtl {
	return &NamespaceCtl{}
}

func (n *NamespaceCtl) ListAll(c *gin.Context) goft.Json {
	return gin.H{
		"code": 20000,
		"data": n.NamespaceMap.ListAll(),
	}
}

func (n *NamespaceCtl) Build(goft *goft.Goft) {
	goft.Handle("GET", "/namespaces", n.ListAll)
}

func (*NamespaceCtl) Name() string {
	return "NamespaceCtl"
}


