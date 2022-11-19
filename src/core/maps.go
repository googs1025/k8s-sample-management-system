package core

import (
	"fmt"
	v1 "k8s.io/api/apps/v1"
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
		deploymentList = make([]*v1.Deployment, 0)
		deploymentList = append(deploymentList.([]*v1.Deployment), deployment)
		d.data.Store(deployment.Namespace, deploymentList)
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