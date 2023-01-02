package services

import (
	"k8s.io/client-go/kubernetes"
	"k8s-Management-System/src/models"
)

//@service
type ConfigMapService struct {
	Client *kubernetes.Clientset `inject:"-"`
	ConfigMap *ConfigMap  `inject:"-"`
}

func NewConfigMapService() *ConfigMapService {
	return &ConfigMapService{}
}

// 前台用于显示Secret列表
func(cms *ConfigMapService) ListConfigMapByNamespace(namespace string) []*models.ConfigMapModel {
	list := cms.ConfigMap.ListAll(namespace)
	ret := make([]*models.ConfigMapModel, len(list))
	for i, item := range list {
		ret[i] = &models.ConfigMapModel {
			Name: item.Name,
			CreateTime: item.CreationTimestamp.Format("2006-01-02 15:04:05"),
			NameSpace: item.Namespace,
		}
	}
	return ret
}
