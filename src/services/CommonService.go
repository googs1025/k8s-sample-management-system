package services

import (
	"fmt"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

type CommonService struct {

}

func NewCommonService() *CommonService {
	return &CommonService{}
}


func (c *CommonService) GetImages(deployment v1.Deployment) string {
	return c.GetImagesByPod(deployment.Spec.Template.Spec.Containers)
}

func (c *CommonService) GetImagesByPod(containers []corev1.Container) string {
	images := containers[0].Image
	if imagesLen := len(containers); imagesLen > 1 {
		images += fmt.Sprintf("其他%d个镜像", imagesLen-1)
	}
	return images
}
