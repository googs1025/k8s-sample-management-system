package services

import (
	"k8s-Management-System/src/wscore"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"log"
)

// DeploymentHandler 使用informer后 回调的方法
type DeploymentHandler struct {
	DeploymentMap *DeploymentMap `inject:"-"`
	DeploymentService *DeploymentService `inject:"-"`
}

func (d *DeploymentHandler) OnAdd(obj interface{}) {
	if dep, ok := obj.(*v1.Deployment); ok {
		d.DeploymentMap.Add(dep)
	}
	wscore.ClientMap.SendAllDepList(d.DeploymentService.ListAll(obj.(*v1.Deployment).Namespace))
}

func (d *DeploymentHandler) OnDelete(obj interface{}) {
	if dep, ok := obj.(*v1.Deployment); ok {
		d.DeploymentMap.Delete(dep)
	}
	wscore.ClientMap.SendAllDepList(d.DeploymentService.ListAll(obj.(*v1.Deployment).Namespace))
}

func (d *DeploymentHandler) OnUpdate(oldObj, newObj interface{}) {
	err := d.DeploymentMap.Update(newObj.(*v1.Deployment))
	if err != nil {
		log.Println(err)
	} else {
		wscore.ClientMap.SendAllDepList(d.DeploymentService.ListAll(newObj.(*v1.Deployment).Namespace))
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