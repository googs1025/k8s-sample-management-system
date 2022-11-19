package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/shenyisyn/goft-gin/goft"
	"k8s-Management-System/src/services"
	"k8s.io/client-go/kubernetes"
)

// deployment控制器
type DeploymentCtl struct {
	K8sClient kubernetes.Interface `inject:"-"`
	DeploymentService *services.DeploymentService  `inject:"-"`
}

func NewDeploymentCtl() *DeploymentCtl {
	return &DeploymentCtl{
	}
}

// Name 实现deployment controller 框架规范
func (*DeploymentCtl) Name() string {
	return "DeploymentCtl"
}

// Build 实现deployment controller 路由 框架规范
func (d *DeploymentCtl) Build(goft *goft.Goft) {
	goft.Handle("GET", "/deployments", d.List)
}

func (d *DeploymentCtl) List(c *gin.Context) goft.Json {
	return d.DeploymentService.ListAll("default")
}



//// 测试用
//func (d *DeploymentCtl) List(c *gin.Context) goft.Json {
//	deploymentList, err := d.K8sClient.AppsV1().Deployments("default").List(c, metav1.ListOptions{})
//	goft.Error(err)
//	return deploymentList
//}
