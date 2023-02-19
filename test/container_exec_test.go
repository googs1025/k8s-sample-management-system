package test

import (
	"fmt"
	"k8s-Management-System/src/common"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"log"
	"os"
	"testing"
)

// 本地执行远端pod的容器命令
func TestContainerExec(t *testing.T) {
	// 使用本地的config
	path := common.GetWd()
	config, err := clientcmd.BuildConfigFromFlags("", path + "/config" )
	if err!=nil{
		log.Fatal(err)
	}
	config.Insecure = true
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	option := &v1.PodExecOptions{
		Container:"myapp",
		Command: []string{"sh","-c","ls"},
		Stdin:   true,
		Stdout:  true,
		Stderr:  true,
	}

	req := client.CoreV1().RESTClient().Post().Resource("pods").
		Namespace("default").
		Name("myapp-rs-kr487").
		SubResource("exec").VersionedParams(
		option,
		scheme.ParameterCodec,
	)
	fmt.Println(req.URL())


	exec, err := remotecommand.NewSPDYExecutor(config,"POST",
		req.URL())
	if err != nil {
		log.Fatal(err)
	}
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin: os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tty: true,
	})
	if err != nil {
		log.Fatal(err)
	}

}


//func TestContainerExecServer(t *testing.T) {
//	// 使用本地的config
//	config, err := clientcmd.BuildConfigFromFlags("","/Users/zhenyu.jiang/go/src/golanglearning/new_project/k8s-Management-System/config" )
//	if err!=nil{
//		log.Fatal(err)
//	}
//	config.Insecure = true
//	client, err := kubernetes.NewForConfig(config)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// 测试server
//	r := gin.New()
//	r.POST("/", func(c *gin.Context) {
//		body, _ := c.GetRawData() //获取body 原始内容
//		cmd := strings.Split(string(body)," ")
//		err = helpers.HandleCommand(client, config, cmd).Stream(remotecommand.StreamOptions{
//			Stdout: c.Writer,
//			Stderr: os.Stderr,
//			Tty: true,
//		})
//	})
//	r.Run(":8088")
//
//}

//func TestContainerExecServerWsSocket(t *testing.T) {
//	// 本地config
//	config, err := clientcmd.BuildConfigFromFlags("","/Users/zhenyu.jiang/go/src/golanglearning/new_project/k8s-Management-System/config" )
//	if err != nil {
//		log.Fatal(err)
//	}
//	config.Insecure = true
//	client, err := kubernetes.NewForConfig(config)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	r := gin.New()
//	r.GET("/", func(c *gin.Context) {
//		wsClient, err := wscore.Upgrader.Upgrade(c.Writer, c.Request,nil)
//		if err != nil {
//			return
//		}
//		shellClient := wscore.NewWsShellClient(wsClient)
//		err = helpers.HandleCommand(client, config, []string{"sh"}).
//			Stream(remotecommand.StreamOptions{
//				Stdin: shellClient,
//				Stdout: shellClient,
//				Stderr: shellClient,
//				Tty: true,
//			})
//		if err != nil {
//			log.Println(err)
//		}
//
//	})
//	r.Run(":8088")
//}
