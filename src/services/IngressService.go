package services

import (
	"context"
	"k8s-Management-System/src/models"
	"k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"strconv"
	"strings"
)

//@service
type IngressService struct {
	Client *kubernetes.Clientset `inject:"-"`
}

func NewIngressService() *IngressService {
	return &IngressService{}
}

// PostIngress 创建Ingress
func(i *IngressService) PostIngress(post *models.IngressPost) error{
	className := "nginx"
	ingressRules := []v1beta1.IngressRule{}
	// 凑 Rule对象
	for _, r := range post.Rules {
		httpRuleValue := &v1beta1.HTTPIngressRuleValue{}
		rulePaths := make([]v1beta1.HTTPIngressPath, 0)
		for _, pathCfg := range r.Paths {
			port,err:=strconv.Atoi(pathCfg.Port)
			if err!=nil{
				return err
			}
			rulePaths = append(rulePaths, v1beta1.HTTPIngressPath{
				Path: pathCfg.Path,
				Backend: v1beta1.IngressBackend{
					ServiceName: pathCfg.SvcName,
					ServicePort: intstr.FromInt(port), //这里需要FromInt
				},
			})
		}
		httpRuleValue.Paths = rulePaths
		rule := v1beta1.IngressRule{
			Host: r.Host,
			IngressRuleValue: v1beta1.IngressRuleValue{
				HTTP: httpRuleValue,
			},
		}
		ingressRules = append(ingressRules, rule)
	}

	// 凑 Ingress对象
	ingress := &v1beta1.Ingress{
		TypeMeta:v1.TypeMeta{
			Kind: "Ingress",
			APIVersion: "networking.k8s.io/v1beta1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      post.Name,
			Namespace: post.Namespace,
			Annotations: i.parseAnnotations(post.Annotations),

		},
		Spec: v1beta1.IngressSpec{
			IngressClassName: &className,
			Rules:            ingressRules,
		},
	}

	_, err := i.Client.NetworkingV1beta1().Ingresses(post.Namespace).
		Create(context.Background(), ingress, v1.CreateOptions{})
	return err

}


// 解析标签
func(i *IngressService) parseAnnotations(annos string) map[string]string{
	replace := []string{"\t"," ","\n","\r\n"}
	for _, r := range replace{
		annos = strings.ReplaceAll(annos,r,"")
	}
	ret := make(map[string]string)
	list := strings.Split(annos,";")
	for _ , item := range list {
		annos := strings.Split(item,":")
		if len(annos) == 2 {
			ret[annos[0]] = annos[1]
		}
	}
	return ret

}