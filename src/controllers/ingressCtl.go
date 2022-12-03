package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/shenyisyn/goft-gin/goft"
	"k8s-Management-System/src/services"
)

type IngressCtl struct{
	IngressMap *services.IngressMap `inject:"-"`
}
func NewIngressCtl() *IngressCtl{
	return &IngressCtl{}
}
func(*IngressCtl)  Name() string{
	return "IngressCtl"
}

func(i *IngressCtl) ListAll(c *gin.Context) goft.Json{
	ns:=c.DefaultQuery("namespace","default")
	return gin.H{
		"code":20000,
		"data":i.IngressMap.ListAll(ns), //暂时 不分页
	}
}

func(i *IngressCtl)  Build(goft *goft.Goft){
	goft.Handle("GET","/ingress", i.ListAll)
}
