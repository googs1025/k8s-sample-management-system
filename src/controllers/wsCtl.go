package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/shenyisyn/goft-gin/goft"
	"k8s-Management-System/src/helpers"
	"k8s-Management-System/src/models"
	"k8s-Management-System/src/wscore"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"log"
)

//@Controller
type WsCtl struct {
	Client *kubernetes.Clientset  `inject:"-"`
	Config *rest.Config  `inject:"-"`
	SysConfig *models.SysConfig   `inject:"-"`
}


func NewWsCtl() *WsCtl {
	return &WsCtl{}
}

func(w *WsCtl) Connect(c *gin.Context) string  {
	client, err := wscore.Upgrader.Upgrade(c.Writer,c.Request,nil)  //升级
	if err != nil {
		log.Println(err)
		return err.Error()
	} else {
		wscore.ClientMap.Store(client)
		return "success"
	}

}

// PodConnect 实现前端 pod 内部容器 命令行执行
func(w *WsCtl) PodConnect(c *gin.Context) (v goft.Void) {
	namespace := c.Query("namespace")
	pod := c.Query("pod")
	container := c.Query("c")

	wsClient, err := wscore.Upgrader.Upgrade(c.Writer,c.Request,nil)
	if err != nil {
		return
	}
	shellClient := wscore.NewWsShellClient(wsClient)
	err = helpers.HandleCommand(namespace, pod, container, w.Client, w.Config, []string{"sh"}).
		Stream(remotecommand.StreamOptions{
			Stdin: shellClient,
			Stdout: shellClient,
			Stderr: shellClient,
			Tty: true,
		})
	return
}

// NodeConnect 实现前端 node介面 命令行执行
func(w *WsCtl) NodeConnect(c *gin.Context) (v goft.Void){
	nodeName := c.Query("node")
	wsClient, err := wscore.Upgrader.Upgrade(c.Writer, c.Request,nil)
	nodeConfig := helpers.GetNodeConfig(w.SysConfig, nodeName) //读取配置文件
	goft.Error(err)

	shellClient := wscore.NewWsShellClient(wsClient)
	// session, err := helpers.SSHConnect(helpers.TempSSHUser,  helpers.TempSSHPWD, helpers.TempSSHIP ,22)
	session, err := helpers.SSHConnect(nodeConfig.User, nodeConfig.Pass, nodeConfig.Ip ,22 )
	//fmt.Println("error:!!!!", err)
	goft.Error(err)
	defer session.Close()
	session.Stdout = shellClient
	session.Stderr = shellClient
	session.Stdin = shellClient
	err = session.RequestPty("xterm-256color", 300, 500, helpers.NodeShellModes)
	goft.Error(err)

	err = session.Run("sh")
	goft.Error(err)
	return
}

func(w *WsCtl)  Build(goft *goft.Goft){
	goft.Handle("GET","/ws", w.Connect)
	goft.Handle("GET","/podws", w.PodConnect)
	goft.Handle("GET","/nodews", w.NodeConnect)
}

func(w *WsCtl) Name() string{
	return "WsCtl"
}