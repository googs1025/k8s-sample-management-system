package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/shenyisyn/goft-gin/goft"
	"k8s-Management-System/src/models"
	"k8s.io/client-go/kubernetes"
	"strings"
)

type ResourcesCtl struct {
	Client *kubernetes.Clientset `inject:"-"`
}


func NewResourcesCtl() *ResourcesCtl {
	return &ResourcesCtl{}
}


func(rc *ResourcesCtl) GetGroupVersion(str string) (group, version string) {
	list := strings.Split(str,"/")

	if len(list) == 1 {
		return "core", list[0]
	} else if len(list) == 2 {
		return list[0], list[1]
	}

	panic("error GroupVersion"+str)
}

// ListResources 获取所有资源: ex kubectl api-resources
func(rc *ResourcesCtl) ListResources(c *gin.Context) goft.Json{

	_, res, err := rc.Client.ServerGroupsAndResources()
	goft.Error(err)
	gRes := make([]*models.GroupResources, 0)
	for _, r := range res {
		group, version := rc.GetGroupVersion(r.GroupVersion)

		gr := &models.GroupResources{
			Group: group,
			Version: version,
			Resources: make([]*models.Resources, 0),
		}

		for _, rr := range r.APIResources {
			res := &models.Resources{Name: rr.Name,Verbs: rr.Verbs}
			gr.Resources = append(gr.Resources,res)
		}
		gRes = append(gRes,gr)
	}
	return gin.H{
		"code": 20000,
		"data": gRes,
	}
}

func(*ResourcesCtl) Name() string{
	return "Resources"
}

func(rc *ResourcesCtl) Build(goft *goft.Goft){
	goft.Handle("GET","/resources", rc.ListResources)
}