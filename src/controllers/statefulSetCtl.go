package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/shenyisyn/goft-gin/goft"
	"k8s-Management-System/src/services"
	"k8s.io/client-go/kubernetes"
)

// statefulset控制器
type StatefulSetCtl struct {
	K8sClient kubernetes.Interface `inject:"-"`
	StatefulSetService *services.StatefulSetService  `inject:"-"`
}

func NewStatefulSetCtl() *StatefulSetCtl {
	return &StatefulSetCtl{
	}
}

// Name 实现statefulset controller 框架规范
func (*StatefulSetCtl) Name() string {
	return "StatefulSetCtl"
}

// Build 实现job controller 路由 框架规范
func (s *StatefulSetCtl) Build(goft *goft.Goft) {
	goft.Handle("GET", "/statefulsets", s.List)
}

// List 获取列表
func (s *StatefulSetCtl) List(c *gin.Context) goft.Json {
	namespace := c.DefaultQuery("namespace", "default") // 请求： GET /cronjobs?namespace=xxxxxxx

	// 配合前端
	return gin.H{
		"code": 20000,
		"data": s.StatefulSetService.ListAll(namespace),
	}

}

// 新增statefulset crud功能

