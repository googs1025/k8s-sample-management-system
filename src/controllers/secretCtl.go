package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/shenyisyn/goft-gin/goft"
	"k8s-Management-System/src/models"
	"k8s-Management-System/src/services"
	corev1 "k8s.io/api/core/v1"
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

// PostSecret 提交 Secret
func(s *SecretCtl) PostSecret(c *gin.Context) goft.Json{
	postModel := &models.PostSecretModel{}
	err := c.ShouldBindJSON(postModel)

	goft.Error(err)
	_, err = s.Client.CoreV1().Secrets(postModel.NameSpace).Create(
		c,
		&corev1.Secret{
			ObjectMeta:v1.ObjectMeta{
				Name:postModel.Name,
				Namespace:postModel.NameSpace,
			},
			Type:corev1.SecretType(postModel.Type),
			StringData:postModel.Data,
		},
		v1.CreateOptions{},
	)
	goft.Error(err)
	return gin.H{
		"code": 20000,
		"data": "OK",
	}
}

func(s *SecretCtl) ListAll(c *gin.Context) goft.Json{
	ns := c.DefaultQuery("namespace","default")
	return gin.H{
		"code":20000,
		"data": s.SecretService.ListSecretByNamespace(ns), //不分页
	}
}

// 查看Secret详细
func(s *SecretCtl) Detail(c *gin.Context) goft.Json{
	ns := c.Param("namespace")
	name := c.Param("name")
	if ns == "" || name == "" {
		panic("error param: namespace or name")
	}
	secret, err := s.Client.CoreV1().Secrets(ns).Get(c,name,v1.GetOptions{})
	goft.Error(err)

	return gin.H{
		"code": 20000,
		"data": &models.SecretModel {
			Name: secret.Name,
			NameSpace: secret.Namespace,
			Type: string(secret.Type),
			CreateTime: secret.CreationTimestamp.Format("2006-01-02 15:04:05"),
			Data: secret.Data,
		},
	}
}

func(s *SecretCtl)  Build(goft *goft.Goft){
	goft.Handle("GET","/secrets", s.ListAll)
	goft.Handle("GET","/secrets/:namespace/:name", s.Detail)
	goft.Handle("DELETE","/secrets", s.RmSecret)
	goft.Handle("POST","/secrets", s.PostSecret)
}
