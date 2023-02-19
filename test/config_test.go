package test

import (
	"context"
	"fmt"
	"k8s-Management-System/src/common"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"testing"
)

// TestConfig 使用本地config登入远程k8s
func TestConfig(t *testing.T) {
	path := common.GetWd()
	config, err := clientcmd.BuildConfigFromFlags("", path + "/config")
	if err != nil {
		log.Fatal(err)
	}
	config.Insecure = true
	client, err := kubernetes.NewForConfig(config)
	if err!=nil{
		log.Fatal(err)
	}
	list, err := client.CoreV1().Pods("default").List(context.Background(),v1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	for _, pod := range list.Items {
		fmt.Println(pod.Name)
	}
}
