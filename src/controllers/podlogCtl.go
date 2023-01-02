package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/shenyisyn/goft-gin/goft"
	"io"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"net/http"
	"time"
)

type PodLogsCtl struct {
	Client *kubernetes.Clientset  `inject:"-"`
}

func NewPodLogsCtl() *PodLogsCtl {
	return &PodLogsCtl{}
}

// GetLogs 使用stream流的方式 进行log日志推送
func(pl *PodLogsCtl) GetLogs(c *gin.Context) (v goft.Void) {
	//namespace := c.DefaultQuery("namespace","default")
	//podName := c.DefaultQuery("podname","")
	//cname := c.DefaultQuery("cname","")
	//var tailLine int64=100
	//// 使用调用一次api server的方式获取。 弃用
	////req := pl.Client.CoreV1().Pods(namespace).GetLogs(podName, &v1.PodLogOptions{Container: cname})
	////ret := req.Do(c)
	////b, err := ret.Raw()
	////goft.Error(err)
	////return gin.H{
	////	"code":20000,
	////	"data":string(b),
	////}
	//
	////cc, _ := context.WithTimeout(c,time.Minute*30) //设置半小时超时时间。否则会造成内存泄露
	//
	//req := pl.Client.CoreV1().Pods(namespace).GetLogs(podName, &v1.PodLogOptions{
	//	Follow: true, Container: cname, TailLines: &tailLine})
	//reader, err := req.Stream(c)
	//goft.Error(err)
	//for{
	//	buf := make([]byte,1024)
	//	n, err := reader.Read(buf)
	//	if err != nil && err!=io.EOF {
	//		break
	//	}
	//	if n > 0 {
	//		c.Writer.Write([]byte(string(buf[0:n])))
	//		c.Writer.(http.Flusher).Flush()
	//	}
	//}
	//
	//return

	ns := c.DefaultQuery("ns","default")
	podname := c.DefaultQuery("podname","")
	cname := c.DefaultQuery("cname","")
	var tailLine int64=100
	opt := &v1.PodLogOptions{
		Follow:true ,
		Container:cname,
		TailLines:&tailLine,
	}

	cc, _ := context.WithTimeout(c, time.Minute*30) //设置半小时超时时间。否则会造成内存泄露
	req := pl.Client.CoreV1().Pods(ns).GetLogs(podname,opt)
	reader, err := req.Stream(cc)
	goft.Error(err)
	defer reader.Close()

	for {
		buf := make([]byte, 1024)
		n, err := reader.Read(buf) // 如果 当前日志 读完了。 会阻塞

		if err != nil && err != io.EOF{ //一旦超时 会进入 这个程序 ,,此时一定要break 掉
			break
		}

		w, err := c.Writer.Write([]byte(string(buf[0:n])))
		if w == 0 || err != nil {
			break
		}
		c.Writer.(http.Flusher).Flush()
	}

	return

}

func(*PodLogsCtl) Name() string{
	return "PodLogsCtl"
}

func(pl *PodLogsCtl) Build(goft *goft.Goft){
	goft.Handle("GET","/pods/logs", pl.GetLogs)
}