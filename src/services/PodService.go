package services

import "k8s-Management-System/src/core"

//@Service
type PodService struct {
	PodMap *core.PodMap `inject:"-"`
}

func NewPodService() *PodService {
	return &PodService{}
}

func(p *PodService) ListByNamespace(ns string ) interface{}{
	return p.PodMap.ListByNamespace(ns)
}