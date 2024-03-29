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

// RS handler
func(k *K8sHandler) RsHandler() *services.RsHandler{
	return &services.RsHandler{}
}

// pod handler
func(k *K8sHandler) PodHandlers() *services.PodHandler {
	return &services.PodHandler{}
}

// namespace handler
func (k *K8sHandler) NamespaceHandler() *services.NamespaceHandler {
	return &services.NamespaceHandler{}
}

// event handler
func(k *K8sHandler) EventHandlers() *services.EventHandler{
	return &services.EventHandler{}
}

func(k *K8sHandler) JobHandlers() *services.JobHandler{
	return &services.JobHandler{}
}

func(k *K8sHandler) ServiceHandlers() *services.ServiceHandler{
	return &services.ServiceHandler{}
}

func(k *K8sHandler) CronJobHandlers() *services.CronJobHandler{
	return &services.CronJobHandler{}
}

func(k *K8sHandler) StatefulSetHandlers() *services.StatefulSetHandler{
	return &services.StatefulSetHandler{}
}

// IngressHandler
func(k *K8sHandler) IngressHandler() *services.IngressHandler{
	return &services.IngressHandler{}
}

// SecretHandler
func(k *K8sHandler) SecretHandler() *services.SecretHandler{
	return &services.SecretHandler{}
}

// ConfigMapHandler
func(k *K8sHandler) ConfigMapHandler() *services.ConfigMapHandler{
	return &services.ConfigMapHandler{}
}

func(k *K8sHandler) NodeHandler() *services.NodeHandler{
	return &services.NodeHandler{}
}

func(k *K8sHandler) RoleHandler() *services.RoleHandler{
	return &services.RoleHandler{}
}

// RoleBindingHandler
func(k *K8sHandler) RoleBindingHandler() *services.RoleBindingHander{
	return &services.RoleBindingHander{}
}

// SaHandler
func(k *K8sHandler) SaHandler() *services.SaHandler{
	return &services.SaHandler{}
}

// ClusterRoleHandler
func(k *K8sHandler) ClusterRoleHandler() *services.ClusterRoleHandler{
	return &services.ClusterRoleHandler{}
}