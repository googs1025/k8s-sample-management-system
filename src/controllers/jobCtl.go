package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/shenyisyn/goft-gin/goft"
	"k8s-Management-System/src/services"
	"k8s.io/client-go/kubernetes"
)

// job控制器
type JobCtl struct {
	K8sClient kubernetes.Interface `inject:"-"`
	JobService *services.JobService  `inject:"-"`
}

func NewJobCtl() *JobCtl {
	return &JobCtl{
	}
}

// Name 实现job controller 框架规范
func (*JobCtl) Name() string {
	return "JobCtl"
}

// Build 实现job controller 路由 框架规范
func (j *JobCtl) Build(goft *goft.Goft) {
	goft.Handle("GET", "/jobs", j.List)
}

func (j *JobCtl) List(c *gin.Context) goft.Json {
	namespace := c.DefaultQuery("namespace", "default") // 请求： GET /jobs?namespace=xxxxxxx

	// 配合前端
	return gin.H{
		"code": 20000,
		"data": j.JobService.ListAll(namespace),
	}

}
// TODO: 实现job的crud