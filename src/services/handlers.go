package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s-Management-System/src/wscore"
	v1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
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
	ns := obj.(*v1.Deployment).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"deployments",
			"result":gin.H{"ns":ns,"data":d.DeploymentService.ListAll(obj.(*v1.Deployment).Namespace)},
		},
	)
}

func (d *DeploymentHandler) OnDelete(obj interface{}) {
	if dep, ok := obj.(*v1.Deployment); ok {
		d.DeploymentMap.Delete(dep)
	}
	ns := obj.(*v1.Deployment).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"deployments",
			"result":gin.H{"ns":ns,"data":d.DeploymentService.ListAll(obj.(*v1.Deployment).Namespace)},
		},
	)
}

func (d *DeploymentHandler) OnUpdate(oldObj, newObj interface{}) {
	err := d.DeploymentMap.Update(newObj.(*v1.Deployment))
	if err != nil {
		log.Println(err)
	} else {
		ns := newObj.(*v1.Deployment).Namespace
		wscore.ClientMap.SendAll(
			gin.H{
				"type":"deployments",
				"result":gin.H{"ns":ns,"data":d.DeploymentService.ListAll(newObj.(*v1.Deployment).Namespace)},
			},
		)
	}
}

// pod相关的回调handler
type PodHandler struct {
	PodMap *PodMap `inject:"-"`
	PodService *PodService `inject:"-"`
}

func(p *PodHandler) OnAdd(obj interface{}){
	p.PodMap.Add(obj.(*corev1.Pod))
	ns := obj.(*corev1.Pod).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"pods",
			"result":gin.H{
				"ns": ns,
				"data":p.PodService.ListByNamespace(obj.(*corev1.Pod).Namespace),
			},
		},
	)
}

func(p *PodHandler) OnUpdate(oldObj, newObj interface{}){
	err := p.PodMap.Update(newObj.(*corev1.Pod))
	if err != nil {
		log.Println(err)
	}
	ns := newObj.(*corev1.Pod).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"pods",
			"result":gin.H{
				"ns": ns,
				"data":p.PodService.ListByNamespace(newObj.(*corev1.Pod).Namespace),
			},
		},
	)
}

func(p *PodHandler)	OnDelete(obj interface{}){
	if d, ok := obj.(*corev1.Pod); ok {
		p.PodMap.Delete(d)
	}
	ns := obj.(*corev1.Pod).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"pods",
			"result":gin.H{
				"ns": ns,
				"data":p.PodService.ListByNamespace(obj.(*corev1.Pod).Namespace),
			},
		},
	)
}

// namespace 相关的回调
type NamespaceHandler struct {
	NamespaceMap *NamespaceMap `inject:"-"`
}

func (n *NamespaceHandler) OnAdd(obj interface{}) {
	n.NamespaceMap.Add(obj.(*corev1.Namespace))
}

func (n *NamespaceHandler) OnUpdate(oldObj, newObj interface{}) {
	n.NamespaceMap.Add(newObj.(*corev1.Namespace))
}

func (n *NamespaceHandler) OnDelete(obj interface{}) {
	if d, ok := obj.(*corev1.Namespace); ok {
		n.NamespaceMap.Delete(d)
	}
}

// event 事件相关的handler
type EventHandler struct {
	EventMap *EventMap  `inject:"-"`
}

func(e *EventHandler) storeData(obj interface{},isdelete bool){
	if event, ok := obj.(*corev1.Event); ok {
		key := fmt.Sprintf("%s_%s_%s", event.Namespace, event.InvolvedObject.Kind, event.InvolvedObject.Name)
		if !isdelete {
			e.EventMap.data.Store(key,event)
		} else {
			e.EventMap.data.Delete(key)
		}
	}
}

func(e *EventHandler) OnAdd(obj interface{}){
	e.storeData(obj,false)
}
func(e *EventHandler) OnUpdate(oldObj, newObj interface{}){
	e.storeData(newObj,false)
}
func(e *EventHandler) OnDelete(obj interface{}){
	e.storeData(obj,true)
}


// JobHandler 使用informer后 回调的方法
type JobHandler struct {
	JobMap *JobMap `inject:"-"`
	JobService *JobService `inject:"-"`
}

func (j *JobHandler) OnAdd(obj interface{}) {
	if jj, ok := obj.(*batchv1.Job); ok {
		j.JobMap.Add(jj)
	}
	ns := obj.(*batchv1.Job).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"jobs",
			"result":gin.H{"ns":ns,"data":j.JobService.ListAll(obj.(*batchv1.Job).Namespace)},
		},
	)
}

func (j *JobHandler) OnDelete(obj interface{}) {
	if jj, ok := obj.(*batchv1.Job); ok {
		j.JobMap.Delete(jj)
	}
	ns := obj.(*batchv1.Job).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"jobs",
			"result":gin.H{"ns": ns,"data":j.JobService.ListAll(obj.(*batchv1.Job).Namespace)},
		},
	)
}

func (j *JobHandler) OnUpdate(oldObj, newObj interface{}) {
	err := j.JobMap.Update(newObj.(*batchv1.Job))
	if err != nil {
		log.Println(err)
	} else {
		ns := newObj.(*batchv1.Job).Namespace
		wscore.ClientMap.SendAll(
			gin.H{
				"type":"jobs",
				"result":gin.H{"ns":ns,"data":j.JobService.ListAll(newObj.(*batchv1.Job).Namespace)},
			},
		)
	}
}

// ServiceHandler 使用informer后 回调的方法
type ServiceHandler struct {
	ServiceMap *ServiceMap `inject:"-"`
	ServiceService *ServiceService `inject:"-"`
}

func (s *ServiceHandler) OnAdd(obj interface{}) {
	if ss, ok := obj.(*corev1.Service); ok {
		s.ServiceMap.Add(ss)
	}
	ns := obj.(*corev1.Service).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"services",
			"result":gin.H{"ns":ns,"data":s.ServiceService.ListServiceByNamespace(obj.(*corev1.Service).Namespace)},
		},
	)
}

func (s *ServiceHandler) OnDelete(obj interface{}) {
	if ss, ok := obj.(*corev1.Service); ok {
		s.ServiceMap.Delete(ss)
	}
	ns := obj.(*corev1.Service).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"services",
			"result":gin.H{"ns": ns,"data":s.ServiceService.ListServiceByNamespace(obj.(*corev1.Service).Namespace)},
		},
	)
}

func (s *ServiceHandler) OnUpdate(oldObj, newObj interface{}) {
	err := s.ServiceMap.Update(newObj.(*corev1.Service))
	if err != nil {
		log.Println(err)
	} else {
		ns := newObj.(*corev1.Service).Namespace
		wscore.ClientMap.SendAll(
			gin.H{
				"type":"services",
				"result":gin.H{"ns":ns,"data":s.ServiceService.ListServiceByNamespace(newObj.(*corev1.Service).Namespace)},
			},
		)
	}
}