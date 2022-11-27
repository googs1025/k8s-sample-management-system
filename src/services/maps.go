package services

import (
	"fmt"
	"k8s-Management-System/src/models"
	v1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
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