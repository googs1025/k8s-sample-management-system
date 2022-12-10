package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/shenyisyn/goft-gin/goft"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type PodLogsCtl struct {
	Client *kubernetes.Clientset  `inject:"-"`
}

func pNewPodLogsCtl() *PodLogsCtl {
	return &PodLogsCtl{}
}

func(pl *PodLogsCtl) GetLogs(c *gin.Context) goft.Json{
	namespace := c.DefaultQuery("namespace","default")
	podName := c.DefaultQuery("podname","")
	cname := c.DefaultQuery("cname","")
	req := pl.Client.CoreV1().Pods(namespace).GetLogs(podName, &v1.PodLogOptions{Container: cname})
	ret := req.Do(c)
	b, err := ret.Raw()
	goft.Error(err)
	return gin.H{
		"code":20000,
		"data":string(b),
	}
}

func(*PodLogsCtl)  Name() string{
	return "PodLogsCtl"
}

func(pl *PodLogsCtl)  Build(goft *goft.Goft){
	goft.Handle("GET","/pods/logs",pl.GetLogs )
}