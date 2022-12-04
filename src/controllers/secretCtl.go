package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/shenyisyn/goft-gin/goft"
	"k8s-Management-System/src/services"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

//@rest controller
type SecretCtl struct {
	SecretMap *services.SecretMap `inject:"-"`
	SecretService *services.SecretService `inject:"-"`
	Client *kubernetes.Clientset  `inject:"-"`
}
func NewSecretCtl() *SecretCtl{
	return &SecretCtl{}
}
func(*SecretCtl)  Name() string{
	return "SecretCtl"
}

// DELETE /ingress?ns=xx&name=xx
func(s *SecretCtl) RmSecret(c *gin.Context) goft.Json{
	ns := c.DefaultQuery("namespace","default")
	name := c.DefaultQuery("name","")
	goft.Error(s.Client.CoreV1().Secrets(ns).
		Delete(c,name,v1.DeleteOptions{}))
	return gin.H{
		"code":20000,
		"data":"OK",
	}
}

func(s *SecretCtl) ListAll(c *gin.Context) goft.Json{
	ns := c.DefaultQuery("namespace","default")
	return gin.H{
		"code":20000,
		"data": s.SecretService.ListSecretByNamespace(ns), //暂时 不分页
	}
}
func(s *SecretCtl)  Build(goft *goft.Goft){
	goft.Handle("GET","/secrets", s.ListAll)
	goft.Handle("DELETE","/secrets", s.RmSecret)

}
