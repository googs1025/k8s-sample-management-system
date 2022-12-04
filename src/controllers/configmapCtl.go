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

//@ rest controller
type ConfigMapCtl struct {
	ConfigMap *services.ConfigMap `inject:"-"`
	ConfigService *services.ConfigMapService `inject:"-"`
	Client *kubernetes.Clientset  `inject:"-"`
}

func NewConfigMapCtl() *ConfigMapCtl{
	return &ConfigMapCtl{}
}

func(*ConfigMapCtl)  Name() string{
	return "ConfigMapCtl"
}

// 提交 config map
func(cm *ConfigMapCtl) PostConfigmap(c *gin.Context) goft.Json{
	postModel := &models.PostConfigMapModel{}
	err := c.ShouldBindJSON(postModel)

	goft.Error(err)
	_,err = cm.Client.CoreV1().ConfigMaps(postModel.NameSpace).Create(
		c,
		&corev1.ConfigMap{
			ObjectMeta: v1.ObjectMeta{
				Name: postModel.Name,
				Namespace: postModel.NameSpace,
			},
			Data: postModel.Data,
		},
		v1.CreateOptions{},
	)
	goft.Error(err)
	return gin.H {
		"code":20000,
		"data":"OK",
	}
}

// 列出config map
func(cm *ConfigMapCtl) ListAll(c *gin.Context) goft.Json{
	ns := c.DefaultQuery("namespace","default")
	return gin.H{
		"code": 20000,
		"data": cm.ConfigService.ListConfigMapByNamespace(ns), //暂时 不分页
	}
}

// DELETE /configmaps?ns=xx&name=xx
func(cm *ConfigMapCtl) RmCm(c *gin.Context) goft.Json{
	ns := c.DefaultQuery("namespace","default")
	name := c.DefaultQuery("name","")
	goft.Error(cm.Client.CoreV1().ConfigMaps(ns).
		Delete(c,name,v1.DeleteOptions{}))
	return gin.H{
		"code": 20000,
		"data": "OK",
	}
}

// 查看Configmap详细
func(cm *ConfigMapCtl) Detail(c *gin.Context) goft.Json{
	ns := c.Param("ns")
	name := c.Param("name")
	if ns == "" || name == "" {
		panic("error param:ns or name")
	}
	configmap, err := cm.Client.CoreV1().ConfigMaps(ns).Get(c,name,v1.GetOptions{})
	goft.Error(err)

	return gin.H{
		"code":20000,
		"data":&models.ConfigMapModel {
			Name: configmap.Name,
			NameSpace: configmap.Namespace,
			CreateTime: configmap.CreationTimestamp.Format("2006-01-02 15:04:05"),
			Data: configmap.Data,
		},
	}
}

func(cm *ConfigMapCtl)  Build(goft *goft.Goft){
	goft.Handle("GET","/configmaps", cm.ListAll)
	goft.Handle("GET","/configmaps/:ns/:name", cm.Detail)
	goft.Handle("DELETE","/configmaps", cm.RmCm)
	goft.Handle("POST","/configmaps", cm.PostConfigmap)
}
