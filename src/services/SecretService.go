package services

import (
	"k8s-Management-System/src/models"
	"k8s.io/client-go/kubernetes"
)

//@service
type SecretService struct {
	Client *kubernetes.Clientset `inject:"-"`
	SecretMap *SecretMap `inject:"-"`
}
func NewSecretService() *SecretService {
	return &SecretService{}
}
// 前台用于显示Secret列表
func(s *SecretService) ListSecretByNamespace(namespace string) []*models.SecretModel {

	list := s.SecretMap.ListAll(namespace)
	ret := make([]*models.SecretModel,len(list))
	for i, item := range list{
		ret[i] = &models.SecretModel {
			Name: item.Name,
			CreateTime: item.CreationTimestamp.Format("2006-01-02 15:04:05"),
			NameSpace: item.Namespace,
			Type: models.SecretType[string(item.Type)],// 类型的翻译
		}
	}
	return ret
}
