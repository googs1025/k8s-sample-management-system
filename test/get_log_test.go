package test

import (
	"context"
	"fmt"
	"io"
	"k8s-Management-System/src/configs"
	v1 "k8s.io/api/core/v1"
	"testing"
)

func TestPodLog(t *testing.T) {

	client := configs.NewK8sConfig().InitClient()


	// 单次查看log日志
	res := client.CoreV1().Pods("default").GetLogs("myredis-0", &v1.PodLogOptions{})
	ret := res.Do(context.Background())
	b, _ := ret.Raw()
	fmt.Println(string(b))


	// 流方式watch日志
	res_watch := client.CoreV1().Pods("default").GetLogs("myredis-0", &v1.PodLogOptions{
		Follow: true,
	})
	// 流式监听
	reader, _ := res_watch.Stream(context.Background())
	for {
		buf := make([]byte, 1024)
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			break
		}
		fmt.Println(string(buf[0:n]))
	}




}