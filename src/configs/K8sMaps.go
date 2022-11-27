package configs

import (
	"k8s-Management-System/src/services"
	"k8s.io/client-go/kubernetes"
)

type K8sMaps struct {
	K8sClient kubernetes.Interface `inject:"-"`
}

func NewK8sMaps() *K8sMaps {
	return &K8sMaps{}
}

func (k *K8sMaps) InitDeploymentMap() *services.DeploymentMap {
	return &services.DeploymentMap{}
}

//初始化 pod map
func(k *K8sMaps) InitPodMap() *services.PodMap {
	return &services.PodMap{}
}

// 初始化 namespace map
func (k *K8sMaps) InitNamespaceMap() *services.NamespaceMap {
	return &services.NamespaceMap{}
}