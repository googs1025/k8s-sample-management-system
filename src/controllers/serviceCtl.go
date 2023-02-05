package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/shenyisyn/goft-gin/goft"
	"k8s-Management-System/src/services"
	"k8s.io/client-go/kubernetes"
)

// ServiceCtl service控制器
type ServiceCtl struct {
	K8sClient kubernetes.Interface `inject:"-"`
	ServiceService *services.ServiceService  `inject:"-"`
}

func NewServiceCtl() *ServiceCtl {
	return &ServiceCtl{
	}
}

// Name 实现service controller 框架规范
func (*ServiceCtl) Name() string {
	return "ServiceCtl"
}

// Build 实现service controller 路由 框架规范
func (s *ServiceCtl) Build(goft *goft.Goft) {
	goft.Handle("GET", "/services", s.List)
}

// List 获取service列表
func (s *ServiceCtl) List(c *gin.Context) goft.Json {
	namespace := c.DefaultQuery("namespace", "default") // 请求： GET /jobs?namespace=xxxxxxx

	// 配合前端
	return gin.H{
		"code": 20000,
		"data": s.ServiceService.ListServiceByNamespace(namespace),
	}

}
