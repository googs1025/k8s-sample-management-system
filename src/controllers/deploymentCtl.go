package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/shenyisyn/goft-gin/goft"
	"k8s-Management-System/src/services"
	v1 "k8s.io/api/apps/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// deployment控制器
type DeploymentCtl struct {
	K8sClient kubernetes.Interface `inject:"-"`
	DeploymentService *services.DeploymentService  `inject:"-"`
	DeployMap *services.DeploymentMap `inject:"-"`
	RsMap *services.RsMap  `inject:"-"`
	PodMap *services.PodMap `inject:"-"`
}

func NewDeploymentCtl() *DeploymentCtl {
	return &DeploymentCtl{}
}

// Name 实现deployment controller 框架规范
func (*DeploymentCtl) Name() string {
	return "DeploymentCtl"
}

// Build 实现deployment controller 路由 框架规范
func (d *DeploymentCtl) Build(goft *goft.Goft) {
	goft.Handle("GET", "/deployments", d.List)
	goft.Handle("GET","/deployments/:ns/:name", d.LoadDeployment)
	goft.Handle("POST","/deployments", d.SaveDeployment)
	goft.Handle("DELETE","/deployments/:ns/:name", d.DeleteDeployment)
	// 根据Deployment获取PODS
	goft.Handle("GET","/deployments-pods/:ns/:name", d.LoadDeployPods)
}

// List 获取dep列表
func (d *DeploymentCtl) List(c *gin.Context) goft.Json {
	namespace := c.DefaultQuery("namespace", "default") // 请求： GET /deployments?namespace=xxxxxxx

	// 配合前端
	return gin.H{
		"code": 20000,
		"data": d.DeploymentService.ListAll(namespace),
	}
	//return d.DeploymentService.ListAll(namespace)
}

// LoadDeployment 拿到特定deployment
func(d *DeploymentCtl) LoadDeployment(c *gin.Context) goft.Json{
	ns := c.Param("ns")
	name := c.Param("name")
	dep, err := d.DeployMap.GetDeployment(ns, name)  // k8s原生的deployment对象
	goft.Error(err)
	return gin.H{
		"code": 20000,
		"data": dep,
	}
}

func(d *DeploymentCtl) SaveDeployment(c *gin.Context) goft.Json{
	dep := &v1.Deployment{}
	goft.Error(c.ShouldBindJSON(dep))
	if c.Query("fast") != "" {  //代表是快捷创建。预设置label
		d.initLabel(dep)
	}
	// debug用
	//fmt.Println(dep.Spec.Template.ObjectMeta.Labels)
	//fmt.Println(dep.Spec.Selector)

	update := c.Query("update") //代表是更新
	if update != "" {
		_, err := d.K8sClient.AppsV1().Deployments(dep.Namespace).Update(c, dep, v12.UpdateOptions{})
		goft.Error(err)
	} else {
		_, err := d.K8sClient.AppsV1().Deployments(dep.Namespace).Create(c, dep, v12.CreateOptions{})
		goft.Error(err)
	}

	return gin.H{
		"code": 20000,
		"data": "success",
	}
}

// DeleteDeployment 删除deployment
func(d *DeploymentCtl) DeleteDeployment(c *gin.Context) goft.Json{
	ns := c.Param("ns")
	name := c.Param("name")

	err := d.K8sClient.AppsV1().Deployments(ns).Delete(c, name, v12.DeleteOptions{})
	goft.Error(err)
	return gin.H{
		"code": 20000,
		"data": "success",
	}
}

// initLabel 支持前端快捷创建时需要的初始化label
func(d *DeploymentCtl) initLabel(deploy *v1.Deployment) {
	if deploy.Spec.Selector == nil {
		deploy.Spec.Selector = &v12.LabelSelector{
			MatchLabels: map[string]string{
				"kube-manager-app": deploy.Name,
			},
		}
	}

	if deploy.Spec.Selector.MatchLabels == nil {
		deploy.Spec.Selector.MatchLabels = map[string]string{
			"kube-manager-app":deploy.Name,
		}
	}

	if deploy.Spec.Template.ObjectMeta.Labels == nil {
		deploy.Spec.Template.ObjectMeta.Labels = map[string]string{
			"kube-manager-app":deploy.Name,
		}
	}

	deploy.Spec.Selector.MatchLabels["kube-manager-app"] = deploy.Name

	deploy.Spec.Template.ObjectMeta.Labels["kube-manager-app"] = deploy.Name
}

// LoadDeployPods 加载deployment的pods列表
func(d *DeploymentCtl) LoadDeployPods(c *gin.Context) goft.Json{
	ns := c.Param("ns")
	name := c.Param("name")

	// 1. 先拿到所有deployment
	dep, err := d.DeployMap.GetDeployment(ns, name)// k8s原生 deployment对象
	goft.Error(err)

	// 2. 取得deployment 过滤出 rs的标签
	labels, err := d.getLabelsByDep(dep, ns)  // 根据deployment过滤出 rs，然后直接获取标签
	goft.Error(err)

	// 3. 过滤出pods列表
	podList, err := d.PodMap.ListByLabels(ns, labels)
	goft.Error(err)


	return gin.H{
		"code": 20000,
		"data": podList,
	}
}

const (
	Deployment = "Deployment"
)

// isRsFromDep 查看rs是否为deployment的关联对象
func(d *DeploymentCtl) isRsFromDep(dep *v1.Deployment,set v1.ReplicaSet) bool{
	for _, ref := range set.OwnerReferences {
		if ref.Kind == Deployment && ref.Name == dep.Name {
			return true
		}
	}
	return false
}

// getLabelsByDep 获取deployment下的 ReplicaSet的 标签集合
func(d *DeploymentCtl) getLabelsByDep(dep *v1.Deployment, ns string ) ([]map[string]string,error){
	rsList, err := d.RsMap.ListByNameSpace(ns)  // 根据namespace 取到 所有rs
	goft.Error(err)

	ret := make([]map[string]string, 0)
	for _, item := range rsList{
		if d.isRsFromDep(dep, *item) {
			s, err := v12.LabelSelectorAsMap(item.Spec.Selector)
			if err != nil {
				return nil,err
			}
			ret = append(ret,s)
		}
	}
	return ret,nil
}


//// 测试用
//func (d *DeploymentCtl) List(c *gin.Context) goft.Json {
//	deploymentList, err := d.K8sClient.AppsV1().Deployments("default").List(c, metav1.ListOptions{})
//	goft.Error(err)
//	return deploymentList
//}
