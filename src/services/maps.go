package services

import (
	"fmt"
	"k8s-Management-System/src/models"
	v1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"reflect"
	"sort"
	"sync"
)

// DeploymentMap 使用informer监听资源变化后，事件变化加入map中
type DeploymentMap struct {
	data sync.Map
}

func (d *DeploymentMap) Add(deployment *v1.Deployment) {

	if deploymentList, ok := d.data.Load(deployment.Namespace); ok {
		deploymentList = append(deploymentList.([]*v1.Deployment), deployment)
		d.data.Store(deployment.Namespace, deploymentList)
	} else {
		newDeploymentList := make([]*v1.Deployment, 0)
		newDeploymentList = append(newDeploymentList, deployment)
		d.data.Store(deployment.Namespace, newDeploymentList)
	}

}

func (d *DeploymentMap) Delete(deployment *v1.Deployment) {

	if deploymentList, ok := d.data.Load(deployment.Namespace); ok {
		list := deploymentList.([]*v1.Deployment)
		for k, needDeleteDeployment := range list {
			if deployment.Name == needDeleteDeployment.Name {
				newList := append(list[:k], list[k+1:]...)
				d.data.Store(deployment.Namespace, newList)
				break
			}
		}
	}


}

func (d *DeploymentMap) Update(deployment *v1.Deployment) error {

	if deploymentList, ok := d.data.Load(deployment.Namespace); ok {
		list := deploymentList.([]*v1.Deployment)
		for k, needUpdateDeployment := range list {
			if deployment.Name == needUpdateDeployment.Name {
				list[k] = deployment
			}
		}
		return nil

	}

	return fmt.Errorf("deployment-%s update error", deployment.Name)

}

// ListDeploymentByNamespace 内存中读取deploymentList
func (d *DeploymentMap) ListDeploymentByNamespace(namespace string) ([]*v1.Deployment, error) {
	if deploymentList, ok := d.data.Load(namespace); ok {
		return deploymentList.([]*v1.Deployment), nil
	}

	return nil, fmt.Errorf("list deployment error, not found")
}

// GetDeployment 内存中读取deployment
func (d *DeploymentMap) GetDeployment(namespace string, deploymentName string) (*v1.Deployment, error) {
	if deploymentList, ok := d.data.Load(namespace); ok {
		list := deploymentList.([]*v1.Deployment)
		for _, dep := range list {
			if dep.Name == deploymentName {
				return dep, nil
			}
		}
	}

	return nil, fmt.Errorf("get deployment error, not found")
}


// 保存Pod集合
type PodMap struct {
	data sync.Map  // [key string] []*v1.Pod    key=>namespace
}

func(p *PodMap) ListByNamespace(namespace string) []*corev1.Pod {
	if podList, ok := p.data.Load(namespace); ok {
		return podList.([]*corev1.Pod)
	}
	return nil
}

func(p *PodMap) GetPod(namespace string, podName string) *corev1.Pod{
	if podList, ok := p.data.Load(namespace); ok {
		list := podList.([]*corev1.Pod)
		for _, pod := range list {
			if pod.Name==podName{
				return pod
			}
		}
	}
	return nil
}

func (p *PodMap) Add(pod *corev1.Pod) {
	if podList, ok := p.data.Load(pod.Namespace); ok {
		podList = append(podList.([]*corev1.Pod), pod)
		p.data.Store(pod.Namespace, podList)
	} else {
		newPodList := make([]*corev1.Pod, 0)
		newPodList = append(newPodList, pod)
		p.data.Store(pod.Namespace, newPodList)
	}
}

func (p *PodMap) Update(pod *corev1.Pod) error {
	if podList, ok := p.data.Load(pod.Namespace); ok {
		list := podList.([]*corev1.Pod)
		for i, needUpdatePod := range list {
			if needUpdatePod.Name == pod.Name {
				list[i] = pod
			}
		}
		return nil
	}
	return fmt.Errorf("pod-%s update error",pod.Name)
}

func (p *PodMap) Delete(pod *corev1.Pod) {
	if podList, ok := p.data.Load(pod.Namespace); ok {
		list := podList.([]*corev1.Pod)
		for i, needDeletePod := range list {
			if needDeletePod.Name == pod.Name {
				newList := append(list[:i], list[i+1:]...)
				p.data.Store(pod.Namespace,newList)
				break
			}
		}
	}
}
//根据标签获取 POD列表
func(p *PodMap) ListByLabels(namespace string,labels []map[string]string) ([]*corev1.Pod,error){
	ret := make([]*corev1.Pod, 0)
	if podList, ok := p.data.Load(namespace); ok {
		list := podList.([]*corev1.Pod)
		for _, pod := range list {
			for _, label := range labels {
				if reflect.DeepEqual(pod.Labels, label) {  //标签完全匹配
					ret = append(ret, pod)
				}
			}
		}
		return ret,nil
	}
	return nil, fmt.Errorf("pods not found ")
}

func(p *PodMap) DEBUG_ListByNS(ns string) []*corev1.Pod {
	ret := make([]*corev1.Pod, 0)
	if podList, ok := p.data.Load(ns); ok {
		list := podList.([]*corev1.Pod)
		for _, pod := range list {
			ret = append(ret, pod)
		}

	}
	return ret
}

type NamespaceMap struct {
	data sync.Map
}

func (n *NamespaceMap) Add(ns *corev1.Namespace) {
	n.data.Store(ns.Name, ns)
}

func (n *NamespaceMap) Update(ns *corev1.Namespace) {
	n.data.Store(ns.Name, ns)
}

func (n *NamespaceMap) Delete(ns *corev1.Namespace) {
	n.data.Delete(ns.Name)
}

func (n *NamespaceMap) Get(namespace string) *corev1.Namespace {
	if item, ok := n.data.Load(namespace); ok {
		ns := item.(*corev1.Namespace)
		return ns
	}
	return nil
}

func (n *NamespaceMap) ListAll() []*models.NamespaceModel {

	items := convertToMapItems(n.data)
	sort.Sort(items)

	res := make([]*models.NamespaceModel, 0)
	n.data.Range(func(key, value any) bool {
		nsName := &models.NamespaceModel{Name: key.(string)}
		res = append(res, nsName)
		return true
	})

	return res
}

// event 事件map 相关
// EventSet 集合 用来保存事件, 只保存最新的一条
type EventMap struct {
	data sync.Map   // [key string] *v1.Event
	// key=>namespace+"_"+kind+"_"+name 这里的name 不一定是pod ,这样确保唯一
}
func(e *EventMap) GetMessage(ns string, kind string, name string) string {
	key := fmt.Sprintf("%s_%s_%s", ns, kind,name)
	if v, ok := e.data.Load(key); ok {
		return v.(*corev1.Event).Message
	}
	return ""
}



type MapItems []*MapItem
type MapItem struct {
	key string
	value interface{}
}
//把sync.map  转为 自定义切片
func convertToMapItems(m sync.Map) MapItems{
	items := make(MapItems,0)
	m.Range(func(key, value interface{}) bool {
		items = append(items, &MapItem{key:key.(string),value:value})
		return true
	})
	return items
}

func(m MapItems) Len() int{
	return len(m)
}
func(m MapItems) Less(i, j int) bool{
	return m[i].key < m[j].key
}
func(m MapItems) Swap(i, j int){
	m[i], m[j] = m[j], m[i]
}

// JobMap 使用informer监听资源变化后，事件变化加入map中
type JobMap struct {
	data sync.Map
}

func (j *JobMap) Add(job *batchv1.Job) {

	if jobList, ok := j.data.Load(job.Namespace); ok {
		jobList = append(jobList.([]*batchv1.Job), job)
		j.data.Store(job.Namespace, jobList)
	} else {
		newJobList := make([]*batchv1.Job, 0)
		newJobList = append(newJobList, job)
		j.data.Store(job.Namespace, newJobList)
	}

}

func (j *JobMap) Delete(job *batchv1.Job) {

	if jobList, ok := j.data.Load(job.Namespace); ok {
		list := jobList.([]*batchv1.Job)
		for k, needDeleteJob := range list {
			if job.Name == needDeleteJob.Name {
				newList := append(list[:k], list[k+1:]...)
				j.data.Store(job.Namespace, newList)
				break
			}
		}
	}
}

func (j *JobMap) Update(job *batchv1.Job) error {

	if jobList, ok := j.data.Load(job.Namespace); ok {
		list := jobList.([]*batchv1.Job)
		for k, needUpdateJob := range list {
			if job.Name == needUpdateJob.Name {
				list[k] = job
			}
		}
		return nil

	}

	return fmt.Errorf("job-%s update error", job.Name)

}

// ListJobByNamespace 内存中读取jobList
func (j *JobMap) ListJobByNamespace(namespace string) ([]*batchv1.Job, error) {
	if jobList, ok := j.data.Load(namespace); ok {
		return jobList.([]*batchv1.Job), nil
	}

	return nil, fmt.Errorf("list job error, not found")
}

// GetJob 内存中读取job
func (j *JobMap) GetJob(namespace string, jobName string) (*batchv1.Job, error) {
	if jobList, ok := j.data.Load(namespace); ok {
		list := jobList.([]*batchv1.Job)
		for _, jj := range list {
			if jj.Name == jobName {
				return jj, nil
			}
		}
	}

	return nil, fmt.Errorf("get job error, not found")
}

// ServiceMap 使用informer监听资源变化后，事件变化加入map中
type ServiceMap struct {
	data sync.Map
}

func (s *ServiceMap) Add(service *corev1.Service) {

	if serviceList, ok := s.data.Load(service.Namespace); ok {
		serviceList = append(serviceList.([]*corev1.Service), service)
		s.data.Store(service.Namespace, serviceList)
	} else {
		newServiceList := make([]*corev1.Service, 0)
		newServiceList = append(newServiceList, service)
		s.data.Store(service.Namespace, newServiceList)
	}

}

func (s *ServiceMap) Delete(service *corev1.Service) {

	if serviceList, ok := s.data.Load(service.Namespace); ok {
		list := serviceList.([]*corev1.Service)
		for k, needDeleteService := range list {
			if service.Name == needDeleteService.Name {
				newList := append(list[:k], list[k+1:]...)
				s.data.Store(service.Namespace, newList)
				break
			}
		}
	}
}

func (s *ServiceMap) Update(service *corev1.Service) error {

	if serviceList, ok := s.data.Load(service.Namespace); ok {
		list := serviceList.([]*corev1.Service)
		for k, needUpdateService := range list {
			if service.Name == needUpdateService.Name {
				list[k] = service
			}
		}
		return nil

	}

	return fmt.Errorf("service-%s update error", service.Name)

}

// ListServiceByNamespace 内存中读取serviceList
func (s *ServiceMap) ListServiceByNamespace(namespace string) ([]*corev1.Service, error) {
	if serviceList, ok := s.data.Load(namespace); ok {
		return serviceList.([]*corev1.Service), nil
	}

	return nil, fmt.Errorf("list service error, not found")
}

// GetService 内存中读取service
func (s *ServiceMap) GetService(namespace string, serviceName string) (*corev1.Service, error) {
	if serviceList, ok := s.data.Load(namespace); ok {
		list := serviceList.([]*corev1.Service)
		for _, service := range list {
			if service.Name == serviceName {
				return service, nil
			}
		}
	}

	return nil, fmt.Errorf("get service error, not found")
}

// ServiceMap 使用informer监听资源变化后，事件变化加入map中
type StatefulSetMap struct {
	data sync.Map
}

func (s *StatefulSetMap) Add(statefulSet *v1.StatefulSet) {

	if statefulSetList, ok := s.data.Load(statefulSet.Namespace); ok {
		statefulSetList = append(statefulSetList.([]*v1.StatefulSet), statefulSet)
		s.data.Store(statefulSet.Namespace, statefulSetList)
	} else {
		newStatefulSetList := make([]*v1.StatefulSet, 0)
		newStatefulSetList = append(newStatefulSetList, statefulSet)
		s.data.Store(statefulSet.Namespace, newStatefulSetList)
	}

}

func (s *StatefulSetMap) Delete(statefulSet *v1.StatefulSet) {

	if statefulSetList, ok := s.data.Load(statefulSet.Namespace); ok {
		list := statefulSetList.([]*v1.StatefulSet)
		for k, needDeleteStatefulSet := range list {
			if statefulSet.Name == needDeleteStatefulSet.Name {
				newList := append(list[:k], list[k+1:]...)
				s.data.Store(statefulSet.Namespace, newList)
				break
			}
		}
	}
}

func (s *StatefulSetMap) Update(statefulSet *v1.StatefulSet) error {

	if statefulSetList, ok := s.data.Load(statefulSet.Namespace); ok {
		list := statefulSetList.([]*v1.StatefulSet)
		for k, needUpdateStatefulSet := range list {
			if statefulSet.Name == needUpdateStatefulSet.Name {
				list[k] = statefulSet
			}
		}
		return nil

	}

	return fmt.Errorf("statefulSet-%s update error", statefulSet.Name)

}

// ListStatefulSetByNamespace 内存中读取statefulSetList
func (s *StatefulSetMap) ListStatefulSetByNamespace(namespace string) ([]*v1.StatefulSet, error) {
	if statefulSetList, ok := s.data.Load(namespace); ok {
		return statefulSetList.([]*v1.StatefulSet), nil
	}

	return nil, fmt.Errorf("list statefulSet error, not found")
}

// GetStatefulSet 内存中读取statefulSet
func (s *StatefulSetMap) GetStatefulSet(namespace string, statefulSetName string) (*v1.StatefulSet, error) {
	if statefulSetList, ok := s.data.Load(namespace); ok {
		list := statefulSetList.([]*v1.StatefulSet)
		for _, statefulSet := range list {
			if statefulSet.Name == statefulSetName {
				return statefulSet, nil
			}
		}
	}

	return nil, fmt.Errorf("get statefulSet error, not found")
}

// ServiceMap 使用informer监听资源变化后，事件变化加入map中
type CronJobMap struct {
	data sync.Map
}

func (s *CronJobMap) Add(cronJob *batchv1beta1.CronJob) {

	if cronJobList, ok := s.data.Load(cronJob.Namespace); ok {
		cronJobList = append(cronJobList.([]*batchv1beta1.CronJob), cronJob)
		s.data.Store(cronJob.Namespace, cronJobList)
	} else {
		newCronJobList := make([]*batchv1beta1.CronJob, 0)
		newCronJobList = append(newCronJobList, cronJob)
		s.data.Store(cronJob.Namespace, newCronJobList)
	}

}

func (s *CronJobMap) Delete(cronJob *batchv1beta1.CronJob) {

	if cronJobList, ok := s.data.Load(cronJob.Namespace); ok {
		list := cronJobList.([]*batchv1beta1.CronJob)
		for k, needDeleteCronJob := range list {
			if cronJob.Name == needDeleteCronJob.Name {
				newList := append(list[:k], list[k+1:]...)
				s.data.Store(cronJob.Namespace, newList)
				break
			}
		}
	}
}

func (s *CronJobMap) Update(cronJob *batchv1beta1.CronJob) error {

	if cronJobList, ok := s.data.Load(cronJob.Namespace); ok {
		list := cronJobList.([]*batchv1beta1.CronJob)
		for k, needUpdateCronJob := range list {
			if cronJob.Name == needUpdateCronJob.Name {
				list[k] = cronJob
			}
		}
		return nil

	}

	return fmt.Errorf("cronJob-%s update error", cronJob.Name)

}

// ListCronJobByNamespace 内存中读取cronJobList
func (s *CronJobMap) ListCronJobByNamespace(namespace string) ([]*batchv1beta1.CronJob, error) {
	if cronJobList, ok := s.data.Load(namespace); ok {
		return cronJobList.([]*batchv1beta1.CronJob), nil
	}

	return nil, fmt.Errorf("list cronJob error, not found")
}

// GetCronJob 内存中读取cronJob
func (s *CronJobMap) GetCronJob(namespace string, cronJobName string) (*batchv1beta1.CronJob, error) {
	if cronJobList, ok := s.data.Load(namespace); ok {
		list := cronJobList.([]*batchv1beta1.CronJob)
		for _, cronJob := range list {
			if cronJob.Name == cronJobName {
				return cronJob, nil
			}
		}
	}

	return nil, fmt.Errorf("get cronJob error, not found")
}

type IngressMap struct {
	data sync.Map   // [ns string] []*v1beta1.Ingress
}
//获取单个Ingress
func(i *IngressMap) Get(namespace string,name string) *networkingv1.Ingress{
	if items,ok := i.data.Load(namespace);ok{
		for _, item := range items.([]*networkingv1.Ingress){
			if item.Name==name{
				return item
			}
		}
	}
	return nil
}

func(i *IngressMap) Add(ingress *networkingv1.Ingress){
	if list, ok := i.data.Load(ingress.Namespace); ok {
		list = append(list.([]*networkingv1.Ingress),ingress)
		i.data.Store(ingress.Namespace, list)
	} else {
		i.data.Store(ingress.Namespace, []*networkingv1.Ingress{ingress})
	}
}

func(i *IngressMap) Update(ingress *networkingv1.Ingress) error {
	if list,ok := i.data.Load(ingress.Namespace); ok {
		for ii, needUpdateIngress := range list.([]*networkingv1.Ingress) {
			if needUpdateIngress.Name == ingress.Name {
				list.([]*networkingv1.Ingress)[ii] = ingress
			}
		}
		return nil
	}
	return fmt.Errorf("ingress-%s not found",ingress.Name)
}

func(i *IngressMap) Delete(ingress *networkingv1.Ingress){
	if list, ok := i.data.Load(ingress.Namespace); ok{
		for ii, needDeleteIngress:=range list.([]*networkingv1.Ingress){
			if needDeleteIngress.Name == ingress.Name {
				newList := append(list.([]*networkingv1.Ingress)[:ii], list.([]*networkingv1.Ingress)[ii+1:]...)
				i.data.Store(ingress.Namespace, newList)
				break
			}
		}
	}
}

func(i *IngressMap) ListAll(ns string)[]*models.IngressModel{
	if list,ok:=i.data.Load(ns);ok{
		ingressList:=list.([]*networkingv1.Ingress)

		ret := make([]*models.IngressModel,len(ingressList))
		for ii, item := range ingressList{
			ret[ii] = &models.IngressModel{
				Name: item.Name,
				CreateTime: item.CreationTimestamp.Format("2006-01-02 15:04:05"),
				NameSpace: item.Namespace,
			}
		}
		return ret
	}
	return []*models.IngressModel{} //返回空列表
}