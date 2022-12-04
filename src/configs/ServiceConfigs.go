package configs

import "k8s-Management-System/src/services"

//@Config
type ServiceConfig struct {}

func NewServiceConfig() *ServiceConfig {
	return &ServiceConfig{}
}

func(*ServiceConfig) CommonService() *services.CommonService{
	return services.NewCommonService()
}

func(*ServiceConfig) DeploymentService() *services.DeploymentService{
	return services.NewDeploymentService()
}

func (*ServiceConfig) PodService() *services.PodService {
	return services.NewPodService()
}

func (*ServiceConfig) JobService() *services.JobService {
	return services.NewJobService()
}

func (*ServiceConfig) ServiceService() *services.ServiceService {
	return services.NewServiceService()
}

func (*ServiceConfig) StatefulSetService() *services.StatefulSetService {
	return services.NewStatefulSetService()
}

func (*ServiceConfig) CronJobService() *services.CronJobService {
	return services.NewCronJobService()
}

func(*ServiceConfig) IngressService() *services.IngressService{
	return services.NewIngressService()
}

func(*ServiceConfig) SecretService() *services.SecretService{
	return services.NewSecretService()
}

func(*ServiceConfig) ConfigMapService() *services.ConfigMapService{
	return services.NewConfigMapService()
}

//func (*ServiceConfig) CRDService() *services.ServiceService {
//	return services.NewCRDService()
//}