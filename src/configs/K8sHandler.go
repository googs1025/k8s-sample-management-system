package configs

import "k8s-Management-System/src/core"

//注入 回调handler
type K8sHandler struct {}

func NewK8sHandler() *K8sHandler {
	return &K8sHandler{}
}

// deployment handler
func(k *K8sHandler) DepHandlers() *core.DeploymentHandler{
	return &core.DeploymentHandler{}
}

// pod handler
func(k *K8sHandler) PodHandlers() *core.PodHandler{
	return &core.PodHandler{}
}
