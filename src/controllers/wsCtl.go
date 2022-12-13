package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/shenyisyn/goft-gin/goft"
	"k8s-Management-System/src/helpers"
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

func(w *WsCtl) PodConnect(c *gin.Context) (v goft.Void) {
	namespace := c.Query("namespace")
	pod := c.Query("pod")
	container := c.Query("c")
	wsClient,err:=wscore.Upgrader.Upgrade(c.Writer,c.Request,nil)
	if err!=nil {
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

func(w *WsCtl)  Build(goft *goft.Goft){
	goft.Handle("GET","/ws", w.Connect)
	goft.Handle("GET","/podws",w.PodConnect)
}

func(w *WsCtl) Name() string{
	return "WsCtl"
}