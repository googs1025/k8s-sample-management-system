package core

import (
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"log"
)

// DeploymentHandler 使用informer后 回调的方法
type DeploymentHandler struct {
	DeploymentMap *DeploymentMap `inject:"-"`
}

func (d *DeploymentHandler) OnAdd(obj interface{}) {
	if dep, ok := obj.(*v1.Deployment); ok {
		d.DeploymentMap.Add(dep)
	}
}

func (d *DeploymentHandler) OnDelete(obj interface{}) {
	if dep, ok := obj.(*v1.Deployment); ok {
		d.DeploymentMap.Delete(dep)
	}
}

func (d *DeploymentHandler) OnUpdate(oldObj, newObj interface{}) {
	err := d.DeploymentMap.Update(newObj.(*v1.Deployment))
	if err != nil {
		log.Println(err)
	}
}

// pod相关的回调handler
type PodHandler struct {
	PodMap *PodMap `inject:"-"`
}

func(p *PodHandler) OnAdd(obj interface{}){
	p.PodMap.Add(obj.(*corev1.Pod))
}

func(p *PodHandler) OnUpdate(oldObj, newObj interface{}){
	err := p.PodMap.Update(newObj.(*corev1.Pod))
	if err != nil {
		log.Println(err)
	}
}

func(p *PodHandler)	OnDelete(obj interface{}){
	if d, ok := obj.(*corev1.Pod); ok {
		p.PodMap.Delete(d)
	}
}