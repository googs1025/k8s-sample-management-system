package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s-Management-System/src/wscore"
	v1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
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

// StatefulSetHandler 使用informer后 回调的方法
type StatefulSetHandler struct {
	StatefulSetMap *StatefulSetMap `inject:"-"`
	StatefulSetService *StatefulSetService `inject:"-"`
}

func (s *StatefulSetHandler) OnAdd(obj interface{}) {
	if ss, ok := obj.(*v1.StatefulSet); ok {
		s.StatefulSetMap.Add(ss)
	}
	ns := obj.(*v1.StatefulSet).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"statefulsets",
			"result":gin.H{"ns":ns,"data":s.StatefulSetService.ListAll(obj.(*v1.StatefulSet).Namespace)},
		},
	)
}

func (s *StatefulSetHandler) OnDelete(obj interface{}) {
	if ss, ok := obj.(*v1.StatefulSet); ok {
		s.StatefulSetMap.Delete(ss)
	}
	ns := obj.(*v1.StatefulSet).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"statefulsets",
			"result":gin.H{"ns": ns,"data":s.StatefulSetService.ListAll(obj.(*v1.StatefulSet).Namespace)},
		},
	)
}

func (s *StatefulSetHandler) OnUpdate(oldObj, newObj interface{}) {
	err := s.StatefulSetMap.Update(newObj.(*v1.StatefulSet))
	if err != nil {
		log.Println(err)
	} else {
		ns := newObj.(*v1.StatefulSet).Namespace
		wscore.ClientMap.SendAll(
			gin.H{
				"type":"statefulsets",
				"result":gin.H{"ns":ns,"data":s.StatefulSetService.ListAll(newObj.(*v1.StatefulSet).Namespace)},
			},
		)
	}
}

// CronJobHandler 使用informer后 回调的方法
type CronJobHandler struct {
	CronJobMap *CronJobMap `inject:"-"`
	CronJobService *CronJobService `inject:"-"`
}

func (cj *CronJobHandler) OnAdd(obj interface{}) {
	if ss, ok := obj.(*batchv1beta1.CronJob); ok {
		cj.CronJobMap.Add(ss)
	}
	ns := obj.(*batchv1beta1.CronJob).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"cronjobs",
			"result":gin.H{"ns":ns,"data": cj.CronJobService.ListAll(obj.(*batchv1beta1.CronJob).Namespace)},
		},
	)
}

func (cj *CronJobHandler) OnDelete(obj interface{}) {
	if ss, ok := obj.(*batchv1beta1.CronJob); ok {
		cj.CronJobMap.Delete(ss)
	}
	ns := obj.(*batchv1beta1.CronJob).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"cronjobs",
			"result":gin.H{"ns": ns,"data": cj.CronJobService.ListAll(obj.(*batchv1beta1.CronJob).Namespace)},
		},
	)
}

func (cj *CronJobHandler) OnUpdate(oldObj, newObj interface{}) {
	err := cj.CronJobMap.Update(newObj.(*batchv1beta1.CronJob))
	if err != nil {
		log.Println(err)
	} else {
		ns := newObj.(*batchv1beta1.CronJob).Namespace
		wscore.ClientMap.SendAll(
			gin.H{
				"type":"cronjobs",
				"result":gin.H{"ns":ns,"data": cj.CronJobService.ListAll(newObj.(*batchv1beta1.CronJob).Namespace)},
			},
		)
	}
}

// ingress相关handler
type IngressHandler struct {
	IngressMap *IngressMap `inject:"-"`
}

func(i *IngressHandler) OnAdd(obj interface{}){
	i.IngressMap.Add(obj.(*networkingv1.Ingress))
	ns := obj.(*networkingv1.Ingress).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"ingress",
			"result":gin.H{"ns": ns,
				"data": i.IngressMap.ListAll(ns)},
		},
	)
}
func(i *IngressHandler) OnUpdate(oldObj, newObj interface{}){
	err := i.IngressMap.Update(newObj.(*networkingv1.Ingress))
	if err != nil {
		log.Println(err)
		return
	}
	ns := newObj.(*networkingv1.Ingress).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"ingress",
			"result":gin.H{"ns": ns,
				"data": i.IngressMap.ListAll(ns)},
		},
	)

}
func(i *IngressHandler) OnDelete(obj interface{}){
	i.IngressMap.Delete(obj.(*networkingv1.Ingress))
	ns := obj.(*networkingv1.Ingress).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"ingress",
			"result":gin.H{"ns": ns,
				"data": i.IngressMap.ListAll(ns)},
		},
	)
}

//Secret相关的handler
type SecretHandler struct {
	SecretMap *SecretMap  `inject:"-"`
	SecretService *SecretService  `inject:"-"`
}
func(s *SecretHandler) OnAdd(obj interface{}){
	s.SecretMap.Add(obj.(*corev1.Secret))
	ns:=obj.(*corev1.Secret).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"secret",
			"result":gin.H{"ns": ns,
				"data": s.SecretService.ListSecretByNamespace(ns)},
		},
	)
}
func(s *SecretHandler) OnUpdate(oldObj, newObj interface{}){
	err:=s.SecretMap.Update(newObj.(*corev1.Secret))
	if err!=nil{
		log.Println(err)
		return
	}
	ns:=newObj.(*corev1.Secret).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"secret",
			"result":gin.H{"ns": ns,
				"data": s.SecretService.ListSecretByNamespace(ns)},
		},
	)
}
func(s *SecretHandler) OnDelete(obj interface{}){
	s.SecretMap.Delete(obj.(*corev1.Secret))
	ns:=obj.(*corev1.Secret).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"secret",
			"result":gin.H{"ns": ns,
				"data": s.SecretService.ListSecretByNamespace(ns)},
		},
	)
}

// ConfigMap相关的handler
type ConfigMapHandler struct {
	ConfigMap *ConfigMap  `inject:"-"`
	ConfigMapService *ConfigMapService  `inject:"-"`
}

func(cm *ConfigMapHandler) OnAdd(obj interface{}){
	cm.ConfigMap.Add(obj.(*corev1.ConfigMap))
	ns := obj.(*corev1.ConfigMap).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"cm",
			"result":gin.H{"ns":ns,
				"data": cm.ConfigMapService.ListConfigMapByNamespace(ns)},
		},
	)
}
func(cm *ConfigMapHandler) OnUpdate(oldObj, newObj interface{}){
	//重点： 只要update返回true 才会发送 。否则不发送
	if cm.ConfigMap.Update(newObj.(*corev1.ConfigMap)){
		ns := newObj.(*corev1.ConfigMap).Namespace
		wscore.ClientMap.SendAll(
			gin.H{
				"type":"cm",
				"result":gin.H{"ns":ns,
					"data": cm.ConfigMapService.ListConfigMapByNamespace(ns)},
			},
		)
	}
}
func(cm *ConfigMapHandler) OnDelete(obj interface{}){
	cm.ConfigMap.Delete(obj.(*corev1.ConfigMap))
	ns := obj.(*corev1.ConfigMap).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"cm",
			"result":gin.H{"ns":ns,
				"data": cm.ConfigMapService.ListConfigMapByNamespace(ns)},
		},
	)
}

//Node相关的handler
type NodeHandler struct {
	NodeMap *NodeMap  `inject:"-"`
	NodeService *NodeService `inject:"-"`
}

func(nm *NodeHandler) OnAdd(obj interface{}){
	nm.NodeMap.Add(obj.(*corev1.Node))

	wscore.ClientMap.SendAll(
		gin.H{
			"type":"node",
			"result":gin.H{"ns":"node",
				"data":nm.NodeService.ListAllNodes()},
		},
	)
}
func(nm *NodeHandler) OnUpdate(oldObj, newObj interface{}){
	//重点： 只要update返回true 才会发送 。否则不发送
	if nm.NodeMap.Update(newObj.(*corev1.Node)){
		wscore.ClientMap.SendAll(
			gin.H{
				"type":"node",
				"result":gin.H{"ns":"node",
					"data":nm.NodeService.ListAllNodes()},
			},
		)
	}
}
func(nm *NodeHandler) OnDelete(obj interface{}){
	nm.NodeMap.Delete(obj.(*corev1.Node))

	wscore.ClientMap.SendAll(
		gin.H{
			"type":"node",
			"result":gin.H{"ns":"node",
				"data":nm.NodeService.ListAllNodes()},
		},
	)
}

type RoleHandler struct {
	RoleMap *RoleMap  `inject:"-"`
	RoleService *RoleService  `inject:"-"`
}

func(rm *RoleHandler) OnAdd(obj interface{}){
	rm.RoleMap.Add(obj.(*rbacv1.Role))
	ns := obj.(*rbacv1.Role).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"role",
			"result":gin.H{"ns": ns,
				"data": rm.RoleService.ListRoles(ns)},
		},
	)
}
func(rm *RoleHandler) OnUpdate(oldObj, newObj interface{}){
	err := rm.RoleMap.Update(newObj.(*rbacv1.Role))
	if err != nil {
		log.Println(err)
		return
	}

	ns := newObj.(*rbacv1.Role).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type": "role",
			"result": gin.H{"ns": ns,
				"data": rm.RoleService.ListRoles(ns)},
		},
	)
}
func(rm *RoleHandler) OnDelete(obj interface{}) {
	rm.RoleMap.Delete(obj.(*rbacv1.Role))
	ns := obj.(*rbacv1.Role).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type": "role",
			"result": gin.H{"ns": ns,
				"data": rm.RoleService.ListRoles(ns)},
		},
	)
}

