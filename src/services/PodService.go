package services

import "k8s-Management-System/src/models"

//@Service
type PodService struct {
	PodMap *PodMap `inject:"-"`
	Common *CommonService `inject:"-"`
}

func NewPodService() *PodService {
	return &PodService{}
}

func(p *PodService) ListByNamespace(namespace string) []*models.Pod {
	podList := p.PodMap.ListByNamespace(namespace)
	res := make([]*models.Pod, 0)

	for _, pod := range podList {
		res = append(res, &models.Pod{
			Name: pod.Name,
			NameSpace: pod.Namespace,
			Images: p.Common.GetImagesByPod(pod.Spec.Containers),
			NodeName: pod.Spec.NodeName,
			Phase: string(pod.Status.Phase),
			IP: []string{pod.Status.PodIP, pod.Status.HostIP},
			CreateTime: pod.CreationTimestamp.Format("2006-01-02 15:04:05"),

		})
	}

	return res

}