package services

import corev1 "k8s.io/api/core/v1"

//@Service
type SaService struct {
	SaMap *SaMap  `inject:"-"`
}

func NewSaService() *SaService {
	return &SaService{}
}

func(ss *SaService) ListSa(ns string) []*corev1.ServiceAccount {
	//sa:=corev1.ServiceAccount{}

	return ss.SaMap.ListAll(ns)
}
