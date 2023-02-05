package configs

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"k8s-Management-System/src/models"
	"k8s-Management-System/src/services"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/metrics/pkg/client/clientset/versioned"
	"log"
)

// K8sConfig 所有资源的handler
type K8sConfig struct {
	DepHandler *services.DeploymentHandler `inject:"-"`
	RsHandler *services.RsHandler `inject:"-"`
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
	NodeHandler *services.NodeHandler `inject:"-"`
	RoleHandler *services.RoleHandler `inject:"-"`
	RoleBindingHandler *services.RoleBindingHander `inject:"-"`
	SaHandler *services.SaHandler `inject:"-"`
	ClusterRoleHandler *services.ClusterRoleHandler `inject:"-"`
}

func NewK8sConfig() *K8sConfig {
	return &K8sConfig{}
}

// InitSysConfig 系统配置初始化
func(*K8sConfig) InitSysConfig() *models.SysConfig{
	b, err := ioutil.ReadFile("app.yaml")
	if err != nil {
		log.Fatal(err)
	}
	config := &models.SysConfig{}
	err = yaml.Unmarshal(b,config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

// K8sRestConfig 默认读取项目根目录的config文件
func(*K8sConfig) K8sRestConfig() *rest.Config{
	config, err := clientcmd.BuildConfigFromFlags("","/Users/zhenyu.jiang/go/src/golanglearning/new_project/k8s-Management-System/config" )
	config.Insecure = true // 不使用认证的方式
	if err != nil {
		log.Fatal(err)
	}
	return config
}

// metric客户端
func(k *K8sConfig) InitMetricClient() *versioned.Clientset {

	c, err := versioned.NewForConfig(k.K8sRestConfig())
	if err != nil {
		log.Fatal(err)
	}
	return c
}

func (k *K8sConfig) InitClient() kubernetes.Interface {

	client, err := kubernetes.NewForConfig(k.K8sRestConfig())
	if err != nil {
		log.Fatal(err)
	}
	
	return client
}

// InitInformer informer初始化
func (k *K8sConfig) InitInformer() informers.SharedInformerFactory {

	fact := informers.NewSharedInformerFactory(k.InitClient(), 0)

	deploymentInformer := fact.Apps().V1().Deployments()
	deploymentInformer.Informer().AddEventHandler(k.DepHandler)

	rsInformer := fact.Apps().V1().ReplicaSets()
	rsInformer.Informer().AddEventHandler(k.RsHandler)

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

	SecretInformer := fact.Core().V1().Secrets() //监听 Secret
	SecretInformer.Informer().AddEventHandler(k.SecretHandler)

	ConfigMapInformer := fact.Core().V1().ConfigMaps() //监听 Configmap
	ConfigMapInformer.Informer().AddEventHandler(k.ConfigMapHandler)

	NodeInformer := fact.Core().V1().Nodes()
	NodeInformer.Informer().AddEventHandler(k.NodeHandler)

	RoleInformer := fact.Rbac().V1().Roles()
	RoleInformer.Informer().AddEventHandler(k.RoleHandler)

	RolesBindingInformer := fact.Rbac().V1().RoleBindings()
	RolesBindingInformer.Informer().AddEventHandler(k.RoleBindingHandler)

	SaInformer := fact.Core().V1().ServiceAccounts()
	SaInformer.Informer().AddEventHandler(k.SaHandler)

	ClusterRoleInformer := fact.Rbac().V1().ClusterRoles()
	ClusterRoleInformer.Informer().AddEventHandler(k.ClusterRoleHandler)

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