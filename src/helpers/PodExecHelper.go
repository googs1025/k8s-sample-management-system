package helpers

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

// HandleCommand 执行k8s调用pod容器的命令行
func HandleCommand(namespace, pod, container string, client *kubernetes.Clientset, config *rest.Config, command []string) remotecommand.Executor {
	option := &v1.PodExecOptions{
		Container: container,
		Command: command,
		Stdin:   true,
		Stdout:  true,
		Stderr:  true,
		TTY:true,
	}

	req := client.CoreV1().RESTClient().Post().Resource("pods").
		Namespace(namespace).
		Name(pod).
		SubResource("exec").
		Param("color","false").
		VersionedParams(
			option,
			scheme.ParameterCodec,
		)

	exec, err := remotecommand.NewSPDYExecutor(config,"POST", req.URL())
	if err != nil {
		panic(err)
	}

	return exec
}
