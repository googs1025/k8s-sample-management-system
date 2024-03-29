package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/shenyisyn/goft-gin/goft"
	"k8s-Management-System/src/services"
	"k8s.io/client-go/kubernetes"
)

// CronJobCtl job控制器
type CronJobCtl struct {
	K8sClient kubernetes.Interface `inject:"-"`
	CronJobService *services.CronJobService  `inject:"-"`
}

func NewCronJobCtl() *CronJobCtl {
	return &CronJobCtl{}
}

// Name 实现job controller 框架规范
func (*CronJobCtl) Name() string {
	return "CronJobCtl"
}

// Build 实现job controller 路由 框架规范
func (cj *CronJobCtl) Build(goft *goft.Goft) {
	goft.Handle("GET", "/cronjobs", cj.List)
}

// List 列出所有cronjobs
// 请求： GET /cronjobs?namespace=xxxxxxx
func (cj *CronJobCtl) List(c *gin.Context) goft.Json {
	namespace := c.DefaultQuery("namespace", "default")

	// 配合前端
	return gin.H{
		"code": 20000,
		"data": cj.CronJobService.ListAll(namespace),
	}

}

// TODO: 实现 crud的cronjob
