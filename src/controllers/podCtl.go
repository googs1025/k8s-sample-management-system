package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/shenyisyn/goft-gin/goft"
	"k8s-Management-System/src/services"
)

type PodCtl struct {
	PodService *services.PodService `inject:"-"`
}

func NewPodCtl() *PodCtl {
	return &PodCtl{}
}

func(p *PodCtl) List(c *gin.Context) goft.Json{
	namespace := c.DefaultQuery("namespace", "default")
	return gin.H{
		"code": 20000,
		"data": p.PodService.ListByNamespace(namespace),
	}

}


func(p *PodCtl)  Build(goft *goft.Goft){
	goft.Handle("GET","/pods",p.List)
}


func(*PodCtl) Name() string{
	return "PodCtl"
}
