package configs

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
)

type K8sConfig struct {
	
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