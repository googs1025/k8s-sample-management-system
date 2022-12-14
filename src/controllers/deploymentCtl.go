package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/shenyisyn/goft-gin/goft"
	"k8s-Management-System/src/services"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/api/apps/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// deployment控制器
type DeploymentCtl struct {
	K8sClient kubernetes.Interface `inject:"-"`
	DeploymentService *services.DeploymentService  `inject:"-"`
	DeployMap *services.DeploymentMap `inject:"-"`
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
	goft.Handle("GET","/deployments/:ns/:name", d.LoadDeployment)
	goft.Handle("POST","/deployments", d.SaveDeployment)
	goft.Handle("DELETE","/deployments/:ns/:name", d.RmDeployment)
}

// List 获取dep列表
func (d *DeploymentCtl) List(c *gin.Context) goft.Json {
	namespace := c.DefaultQuery("namespace", "default") // 请求： GET /deployments?namespace=xxxxxxx

	// 配合前端
	return gin.H{
		"code": 20000,
		"data": d.DeploymentService.ListAll(namespace),
	}
	//return d.DeploymentService.ListAll(namespace)
}

// LoadDeployment 拿到特定dep
func(d *DeploymentCtl) LoadDeployment(c *gin.Context) goft.Json{
	ns := c.Param("ns")
	name := c.Param("name")
	dep, err := d.DeployMap.GetDeployment(ns, name)// 原生
	goft.Error(err)
	return gin.H{
		"code": 20000,
		"data": dep,
	}
}

func(d *DeploymentCtl) SaveDeployment(c *gin.Context) goft.Json{
	dep := &v1.Deployment{}
	goft.Error(c.ShouldBindJSON(dep))
	_, err := d.K8sClient.AppsV1().Deployments(dep.Namespace).Create(c, dep,v12.CreateOptions{})
	goft.Error(err)
	return gin.H{
		"code": 20000,
		"data": "success",
	}
}

func(d *DeploymentCtl) RmDeployment(c *gin.Context) goft.Json{
	ns := c.Param("ns")
	name := c.Param("name")

	err := d.K8sClient.AppsV1().Deployments(ns).Delete(c, name, v12.DeleteOptions{})
	goft.Error(err)
	return gin.H{
		"code": 20000,
		"data": "success",
	}
}



//// 测试用
//func (d *DeploymentCtl) List(c *gin.Context) goft.Json {
//	deploymentList, err := d.K8sClient.AppsV1().Deployments("default").List(c, metav1.ListOptions{})
//	goft.Error(err)
//	return deploymentList
//}
