package configs

import (
	"k8s-Management-System/src/services"
)

//注入 回调handler
type K8sHandler struct {}

func NewK8sHandler() *K8sHandler {
	return &K8sHandler{}
}

// deployment handler
func(k *K8sHandler) DepHandlers() *services.DeploymentHandler {
	return &services.DeploymentHandler{}
}

// pod handler
func(k *K8sHandler) PodHandlers() *services.PodHandler {
	return &services.PodHandler{}
}

// namespace handler
func (k *K8sHandler) NamespaceHandler() *services.NamespaceHandler {
	return &services.NamespaceHandler{}
}
