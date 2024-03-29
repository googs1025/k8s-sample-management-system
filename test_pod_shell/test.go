package main

import (
	"github.com/gin-gonic/gin"
	"k8s-Management-System/src/common"
	"k8s-Management-System/src/helpers"
	"k8s-Management-System/src/wscore"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"log"
)

func main() {
	// 本地config
	path := common.GetWd()
	config, err := clientcmd.BuildConfigFromFlags("", path + "/config" )
	if err != nil {
		log.Fatal(err)
	}
	config.Insecure = true
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.New()
	r.GET("/", func(c *gin.Context) {
		wsClient, err := wscore.Upgrader.Upgrade(c.Writer, c.Request,nil)
		if err != nil {
			return
		}
		shellClient := wscore.NewWsShellClient(wsClient)
		err = helpers.HandleCommand("default","pod1", "container1", client, config, []string{"sh"}).
			Stream(remotecommand.StreamOptions{
				Stdin: shellClient,
				Stdout: shellClient,
				Stderr: shellClient,
				Tty: true,
			})
		if err != nil {
			log.Println(err)
		}

	})
	r.Run(":8088")
}

