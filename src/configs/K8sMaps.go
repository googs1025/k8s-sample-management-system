package configs

import (
	"k8s-Management-System/src/core"
	"k8s.io/client-go/kubernetes"
)

type K8sMaps struct {
	K8sClient kubernetes.Interface `inject:"-"`
}

func NewK8sMaps() *K8sMaps {
	return &K8sMaps{}
}

func (d *K8sMaps) InitDeploymentMap() *core.DeploymentMap {
	return &core.DeploymentMap{}
}

//初始化 podmap
func(this *K8sMaps) InitPodMap() *core.PodMap{
	return &core.PodMap{}
}
