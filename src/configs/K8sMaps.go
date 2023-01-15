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

//初始化 ReplicaSetMap
func(k *K8sMaps) InitRsMap() *services.RsMap{
	return &services.RsMap{}
}

// 初始化 pod map
func(k *K8sMaps) InitPodMap() *services.PodMap {
	return &services.PodMap{}
}

// 初始化 namespace map
func (k *K8sMaps) InitNamespaceMap() *services.NamespaceMap {
	return &services.NamespaceMap{}
}

// 初始化 event map
func(k *K8sMaps) InitEventMap() *services.EventMap {
	return &services.EventMap{}
}

// 初始化 job map
func(k *K8sMaps) InitJobMap() *services.JobMap {
	return &services.JobMap{}
}

// 初始化 service map
func(k *K8sMaps) InitServiceMap() *services.ServiceMap {
	return &services.ServiceMap{}
}

// 初始化 service map
func(k *K8sMaps) InitStatefulSetMap() *services.StatefulSetMap {
	return &services.StatefulSetMap{}
}

// 初始化 service map
func(k *K8sMaps) InitCronJobMap() *services.CronJobMap {
	return &services.CronJobMap{}
}

// 初始化 ingress map
func(k *K8sMaps) InitIngressMap() *services.IngressMap{
	return &services.IngressMap{}
}

// 初始化 Secret map
func(k *K8sMaps) InitSecretMap() *services.SecretMap{
	return &services.SecretMap{}
}

// 初始化 ConfigMap map
func(k *K8sMaps) InitConfigMap() *services.ConfigMap{
	return &services.ConfigMap{}
}

//初始化NodeMap
func(k *K8sMaps) InitNodeMap() *services.NodeMap{
	return &services.NodeMap{}
}

func(k *K8sMaps) InitRoleMap() *services.RoleMap{
	return &services.RoleMap{}
}

func(k *K8sMaps) InitRoleBindingMap() *services.RoleBindingMap{
	return &services.RoleBindingMap{}
}

func(k *K8sMaps) InitSaMap() *services.SaMap{
	return &services.SaMap{}
}

func(k *K8sMaps) InitClusterRoleMap() *services.ClusterRoleMap{
	return &services.ClusterRoleMap{}
}

// 初始化 service map
//func(k *K8sMaps) InitCRDMap() *services.CRDMap {
//	return &services.CRDMap{}
//}