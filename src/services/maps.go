package services

import (
	"fmt"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"reflect"
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