package services

import (
	"github.com/shenyisyn/goft-gin/goft"
	"k8s-Management-System/src/models"
)

type StatefulSetService struct {
	StatefulSetMap *StatefulSetMap   `inject:"-"`
	Common		  *CommonService `inject:"-"`
}

func NewStatefulSetService() *StatefulSetService {
	return &StatefulSetService{}
}



func (s *StatefulSetService) ListAll(namespace string) (res []*models.StatefulSet) {

	statefulSetList, err := s.StatefulSetMap.ListStatefulSetByNamespace(namespace)
	goft.Error(err)

	for _, statefulSet := range statefulSetList {
		res = append(res, &models.StatefulSet{
			Name: statefulSet.Name,
			NameSpace: statefulSet.Namespace,
			Replicas: statefulSet.Status.Replicas,
			Images: s.Common.GetStatefulSetImages(*statefulSet),
			CreateTime: statefulSet.CreationTimestamp.Format("2006-01-02 15:04:05"),
		})

	}

	return

}
