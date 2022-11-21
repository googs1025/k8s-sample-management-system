package services

//@Service
type PodService struct {
	PodMap *PodMap `inject:"-"`
}

func NewPodService() *PodService {
	return &PodService{}
}

func(p *PodService) ListByNamespace(ns string ) interface{}{
	return p.PodMap.ListByNamespace(ns)
}