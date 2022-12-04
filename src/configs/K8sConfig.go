package configs

import (
	"k8s-Management-System/src/services"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
)

type K8sConfig struct {
	DepHandler *services.DeploymentHandler `inject:"-"`
	PodHandler *services.PodHandler        `inject:"-"`
	NsHandler  *services.NamespaceHandler  `inject:"-"`
	EventHandler *services.EventHandler `inject:"-"`
	JobHandler *services.JobHandler `inject:"-"`
	ServiceHandler *services.ServiceHandler `inject:"-"`
	StatefulSetHandler *services.StatefulSetHandler `inject:"-"`
	CronJobHandler *services.CronJobHandler `inject:"-"`
	IngressHandler *services.IngressHandler `inject:"-"`
	SecretHandler *services.SecretHandler `inject:"-"`
	ConfigMapHandler *services.ConfigMapHandler `inject:"-"`
}

func NewK8sConfig() *K8sConfig {
	return &K8sConfig{}
}

func (*K8sConfig) InitClient() kubernetes.Interface {
	config := &rest.Config{
		Host: "http://1.14.120.233:8009",
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	
	return client
}

func (k *K8sConfig) InitInformer() informers.SharedInformerFactory {

	fact := informers.NewSharedInformerFactory(k.InitClient(), 0)

	deploymentInformer := fact.Apps().V1().Deployments()
	deploymentInformer.Informer().AddEventHandler(k.DepHandler)

	podInformer := fact.Core().V1().Pods() //监听pod
	podInformer.Informer().AddEventHandler(k.PodHandler)

	nsInformer := fact.Core().V1().Namespaces()
	nsInformer.Informer().AddEventHandler(k.NsHandler)

	eventInformer := fact.Core().V1().Events()
	eventInformer.Informer().AddEventHandler(k.EventHandler)

	jobInformer := fact.Batch().V1().Jobs()
	jobInformer.Informer().AddEventHandler(k.JobHandler)

	serviceInformer := fact.Core().V1().Services()
	serviceInformer.Informer().AddEventHandler(k.ServiceHandler)

	statefulSetInformer := fact.Apps().V1().StatefulSets()
	statefulSetInformer.Informer().AddEventHandler(k.StatefulSetHandler)

	cronJobInformer := fact.Batch().V1beta1().CronJobs()
	cronJobInformer.Informer().AddEventHandler(k.CronJobHandler)

	IngressInformer:=fact.Networking().V1beta1().Ingresses() // 监听 Ingress
	IngressInformer.Informer().AddEventHandler(k.IngressHandler)

	SecretInformer := fact.Core().V1().Secrets() //监听Secret
	SecretInformer.Informer().AddEventHandler(k.SecretHandler)

	ConfigMapInformer := fact.Core().V1().ConfigMaps() //监听Configmap
	ConfigMapInformer.Informer().AddEventHandler(k.ConfigMapHandler)



	fact.Start(wait.NeverStop)

	return fact

}


//var K8sClient *kubernetes.Clientset
//
//func init() {
//
//	// 两个选一个用
//	configs := kubeConfig()
//	//configs := kubeProxyConfig()
//
//
//	clientSet, err := kubernetes.NewForConfig(configs)
//	if err != nil {
//		return
//	}
//	K8sClient = clientSet
//}
//
//func HomeDir() string {
//	if h := os.Getenv("HOME"); h != "" {
//		return h
//	}
//
//	return os.Getenv("USERPROFILE")
//
//}
//
//
//// kubeConfig 集权中执行
//func kubeConfig() *rest.Config {
//	// 法一：直接在k8s上运行的代码
//	var kubeConfig *string
//
//	if home := HomeDir(); home != "" {
//		kubeConfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "configs"), "")
//	} else {
//		kubeConfig = flag.String("kubeconfig", "", "")
//	}
//	//flag.Parse()
//
//	configs, err := clientcmd.BuildConfigFromFlags("", *kubeConfig)
//	if err != nil {
//		log.Panic(err.Error())
//	}
//	return configs
//}
//
//// kubeProxyConfig 利用端口转发，可以在本地执行。
//func kubeProxyConfig() *rest.Config {
//	// 法二：需要用端口转换 kubectl proxy --address="0.0.0.0" --accept-hosts='^*$' --port=8009
//	configs := &rest.Config{
//		Host: "http://1.14.120.233:8009",
//	}
//	return configs
//}