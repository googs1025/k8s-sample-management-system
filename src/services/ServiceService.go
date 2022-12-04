package services

import (
	"github.com/shenyisyn/goft-gin/goft"
	"k8s-Management-System/src/models"
)

type ServiceService struct {
	ServiceMap *ServiceMap   `inject:"-"`
	Common *CommonService `inject:"-"`
}

func NewServiceService() *ServiceService {
	return &ServiceService{}
}

func(s *ServiceService) ListServiceByNamespace(namespace string) []*models.Service {
	serviceList, err := s.ServiceMap.ListServiceByNamespace(namespace)
	goft.Error(err)
	res := make([]*models.Service, 0)

	for _, service := range serviceList {
		res = append(res, &models.Service{
			Name: service.Name,
			NameSpace: service.Namespace,
			Type: string(service.Spec.Type),
			ClusterIp: service.Spec.ClusterIP,
			ClusterIps: service.Spec.ClusterIPs,
			Ports: s.Common.ServicePort(service.Spec.Ports),
			CreateTime: service.CreationTimestamp.Format("2006-01-02 15:04:05"),

		})
	}

	return res

}



