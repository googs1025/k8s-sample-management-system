package services

import (
	"github.com/shenyisyn/goft-gin/goft"
	"k8s-Management-System/src/core"
	"k8s-Management-System/src/models"
	v1 "k8s.io/api/apps/v1"
)

type DeploymentService struct {
	DeploymentMap *core.DeploymentMap `inject:"-"`
	Common		  *CommonService	`inject:"-"`
}

func NewDeploymentService() *DeploymentService {
	return &DeploymentService{}
}

func (*DeploymentService) getDeploymentCondition(deployment *v1.Deployment) string {

	for _, item := range deployment.Status.Conditions {
		if string(item.Type) == "Available" && string(item.Status) != "True" {
			return item.Message
		}
	}

	return ""

}

func (*DeploymentService) getDeploymentIsComplete(deployment *v1.Deployment) bool {
	return deployment.Status.Replicas == deployment.Status.AvailableReplicas
}

func (d *DeploymentService) ListAll(namespace string) (res []*models.Deployment) {

	deploymentList, err := d.DeploymentMap.ListDeploymentByNamespace(namespace)
	goft.Error(err)

	for _, deployment := range deploymentList {
		res = append(res, &models.Deployment{
			Name: deployment.Name,
			NameSpace:deployment.Namespace,
			Replicas:[3]int32{deployment.Status.Replicas,deployment.Status.AvailableReplicas,deployment.Status.UnavailableReplicas},
			Images:d.Common.GetImages(*deployment),
			IsComplete:d.getDeploymentIsComplete(deployment),
			Message:d.getDeploymentCondition(deployment),
		})

	}

	return

}