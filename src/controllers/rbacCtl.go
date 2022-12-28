package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shenyisyn/goft-gin/goft"
	"k8s-Management-System/src/services"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type RBACCtl struct {
	RoleService *services.RoleService `inject:"-"`
	SaService *services.SaService `inject:"-"`
	Client *kubernetes.Clientset `inject:"-"`
}

func NewRBACCtl() *RBACCtl {
	return &RBACCtl{}
}

func(r *RBACCtl) Roles(c *gin.Context) goft.Json{
	ns := c.DefaultQuery("ns","default")
	return gin.H{
		"code": 20000,
		"data": r.RoleService.ListRoles(ns),
	}
}

func(*RBACCtl) Name() string{
	return "RBACCtl"
}

func(r *RBACCtl) Build(goft *goft.Goft) {
	goft.Handle("GET","/roles", r.Roles)
	goft.Handle("GET","/roles/:ns/:rolename", r.RolesDetail)
	goft.Handle("POST","/roles/:ns/:rolename", r.UpdateRolesDetail) //修改角色
	goft.Handle("GET","/rolebindings", r.RoleBindingList)
	goft.Handle("PUT","/rolebindings", r.AddUserToRoleBinding) //添加用户到binding
	goft.Handle("POST","/roles", r.CreateRole)
	goft.Handle("DELETE","/roles", r.DeleteRole)
	goft.Handle("POST","/rolebindings", r.CreateRoleBinding)
	goft.Handle("DELETE","/rolebindings", r.DeleteRoleBinding)

	goft.Handle("GET","/sa", r.SaList)

	goft.Handle("GET","/clusterroles", r.ClusterRoles)
	goft.Handle("DELETE","/clusterroles", r.DeleteClusterRole)
	goft.Handle("POST","/clusterroles", r.CreateClusterRole) //创建集群角色
	goft.Handle("POST","/clusterroles/:cname", r.UpdateClusterRolesDetail)
	goft.Handle("GET","/clusterroles/:cname", r.ClusterRolesDetail)

}

func(r *RBACCtl) RoleBindingList(c *gin.Context) goft.Json{
	ns:=c.DefaultQuery("ns","default")
	return gin.H{
		"code": 20000,
		"data": r.RoleService.ListRoleBindings(ns),
	}
}

func(r *RBACCtl) DeleteRole(c *gin.Context) goft.Json{
	ns := c.DefaultQuery("ns","default")
	name := c.DefaultQuery("name","")

	err := r.Client.RbacV1().Roles(ns).Delete(c,name,metav1.DeleteOptions{})
	goft.Error(err)
	return gin.H{
		"code": 20000,
		"data": "success",
	}
}

func(r *RBACCtl) CreateRole(c *gin.Context) goft.Json{

	role := rbacv1.Role{} //原生的k8s role 对象
	goft.Error(c.ShouldBindJSON(&role))
	role.APIVersion = "rbac.authorization.k8s.io/v1"
	role.Kind = "Role"

	_, err := r.Client.RbacV1().Roles(role.Namespace).Create(c, &role, metav1.CreateOptions{})
	goft.Error(err)
	return gin.H{
		"code": 20000,
		"data": "success",
	}
}

func(r *RBACCtl) CreateRoleBinding(c *gin.Context) goft.Json{
	rb := &rbacv1.RoleBinding{}
	goft.Error(c.ShouldBindJSON(rb))
	_, err := r.Client.RbacV1().RoleBindings(rb.Namespace).Create(c, rb, metav1.CreateOptions{})
	goft.Error(err)
	return gin.H{
		"code": 20000,
		"data": "success",
	}
}

func(r *RBACCtl) DeleteRoleBinding(c *gin.Context) goft.Json{
	ns := c.DefaultQuery("ns","default")
	name := c.DefaultQuery("name","")
	err := r.Client.RbacV1().RoleBindings(ns).Delete(c, name, metav1.DeleteOptions{})
	goft.Error(err)
	return gin.H{
		"code": 20000,
		"data": "success",
	}
}

//获取角色详细
func(r *RBACCtl) RolesDetail(c *gin.Context) goft.Json{
	ns := c.Param("ns")
	rname := c.Param("rolename")
	return gin.H{
		"code": 20000,
		"data": r.RoleService.GetRole(ns,rname),
	}
}

//更新角色
func(r *RBACCtl) UpdateRolesDetail(c *gin.Context) goft.Json{
	ns := c.Param("ns")
	rname := c.Param("rolename")
	role := r.RoleService.GetRole(ns,rname)
	postRole := rbacv1.Role{}
	goft.Error(c.ShouldBindJSON(&postRole))  //获取提交过来的对象

	role.Rules = postRole.Rules //目前修改只允许修改 rules，其他不允许。大家可以自行扩展，如标签也允许修改
	_, err := r.Client.RbacV1().Roles(role.Namespace).Update(c, role, metav1.UpdateOptions{})
	goft.Error(err)
	return gin.H{
		"code":20000,
		"data":"success",
	}
}

// AddUserToRoleBinding 从rolebinding中 增加或删除用户
func(r *RBACCtl) AddUserToRoleBinding(c *gin.Context) goft.Json{
	ns := c.DefaultQuery("ns","default")
	name := c.DefaultQuery("name","") //rolebinding 名称
	t := c.DefaultQuery("type","") //如果没传值就是增加，传值（不管什么代表删除)
	subject := rbacv1.Subject{}// 传过来
	goft.Error(c.ShouldBindJSON(&subject))
	if subject.Kind == "ServiceAccount" {
		subject.APIGroup = ""
	}
	rb := r.RoleService.GetRoleBinding(ns,name) //通过名称获取 rolebinding对象
	if t != "" { //代表删除
		for i, sub := range rb.Subjects {
			if sub.Kind == subject.Kind && sub.Name == subject.Name {
				rb.Subjects = append(rb.Subjects[:i], rb.Subjects[i+1:]...)
				break //确保只删一个（哪怕有同名同kind用户)
			}
		}
		fmt.Println(rb.Subjects)
	} else {
		rb.Subjects = append(rb.Subjects,subject)
	}
	_, err := r.Client.RbacV1().RoleBindings(ns).Update(c, rb, metav1.UpdateOptions{})
	goft.Error(err)
	return gin.H{
		"code": 20000,
		"data": "success",
	}
}

func(r *RBACCtl) SaList(c *gin.Context) goft.Json{
	ns := c.DefaultQuery("ns","default")
	return gin.H{
		"code": 20000,
		"data": r.SaService.ListSa(ns),
	}
}

func(r *RBACCtl) ClusterRoles(c *gin.Context) goft.Json{
	return gin.H{
		"code":20000,
		"data": r.RoleService.ListClusterRoles(),
	}
}

func(r *RBACCtl) DeleteClusterRole(c *gin.Context) goft.Json{
	name := c.DefaultQuery("name","")

	err := r.Client.RbacV1().ClusterRoles().Delete(c, name, metav1.DeleteOptions{})
	goft.Error(err)
	return gin.H{
		"code": 20000,
		"data": "success",
	}
}

//创建集群角色
func(r *RBACCtl) CreateClusterRole(c *gin.Context) goft.Json{
	clusterRole := rbacv1.ClusterRole{} //原生的k8s role 对象
	goft.Error(c.ShouldBindJSON(&clusterRole))
	clusterRole.APIVersion = "rbac.authorization.k8s.io/v1"
	clusterRole.Kind = "ClusterRole"
	_, err := r.Client.RbacV1().ClusterRoles().Create(c,&clusterRole,metav1.CreateOptions{})
	goft.Error(err)
	return gin.H{
		"code": 20000,
		"data": "success",
	}
}

//更新集群角色
func(r *RBACCtl) UpdateClusterRolesDetail(c *gin.Context) goft.Json{
	cname := c.Param("cname") //集群角色名
	clusterRole := r.RoleService.GetClusterRole(cname)
	postRole := rbacv1.ClusterRole{}
	goft.Error(c.ShouldBindJSON(&postRole))  //获取提交过来的对象

	clusterRole.Rules = postRole.Rules //目前修改只允许修改 rules，其他不允许。大家可以自行扩展，如标签也允许修改
	_, err := r.Client.RbacV1().ClusterRoles().Update(c,clusterRole,metav1.UpdateOptions{})
	goft.Error(err)
	return gin.H{
		"code": 20000,
		"data": "success",
	}
}

// 获取集群角色详细
func(r *RBACCtl) ClusterRolesDetail(c *gin.Context) goft.Json{

	rname := c.Param("cname") //集群角色名
	return gin.H{
		"code":20000,
		"data": r.RoleService.GetClusterRole(rname),
	}
}