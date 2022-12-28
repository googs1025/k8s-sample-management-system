package services

import (
	"k8s-Management-System/src/models"
	rbacv1 "k8s.io/api/rbac/v1"
)

//@Service
type RoleService struct {
	RoleMap *RoleMap  `inject:"-"`
	RoleBindingMap *RoleBindingMap  `inject:"-"`
	ClusterRoleMap *ClusterRoleMap  `inject:"-"`
}

func NewRoleService() *RoleService {
	return &RoleService{}
}

func(rs *RoleService) ListRoles(ns string) []*models.RoleModel {
	list := rs.RoleMap.ListAll(ns)
	ret := make([]*models.RoleModel,len(list))
	for i, item := range list {
		ret[i] = &models.RoleModel {
			Name: item.Name,
			CreateTime: item.CreationTimestamp.Format("2006-01-02 15:04:05"),
			NameSpace: item.Namespace,
		}
	}
	return ret
}

func(rs *RoleService) ListRoleBindings(ns string) []*models.RoleBindingModel {
	list := rs.RoleBindingMap.ListAll(ns)
	ret := make([]*models.RoleBindingModel,len(list))
	for i, item := range list {
		ret[i] = &models.RoleBindingModel {

			Name: item.Name,
			CreateTime: item.CreationTimestamp.Format("2006-01-02 15:04:05"),
			NameSpace: item.Namespace,
			Subject: item.Subjects,

		}
	}
	return ret
}

func(rs *RoleService) ListClusterRoles() []*rbacv1.ClusterRole {
	return rs.ClusterRoleMap.ListAll()
}

func(rs *RoleService) GetRole(ns, name string) *rbacv1.Role{
	rb := rs.RoleMap.Get(ns,name)
	if rb == nil {
		panic("no such role")
	}
	return rb
}

func(rs *RoleService) GetClusterRole(name string) *rbacv1.ClusterRole{
	rb := rs.ClusterRoleMap.Get(name)
	if rb == nil {
		panic("no such cluster-role")
	}
	return rb
}

func(rs *RoleService) GetRoleBinding(ns ,name string) *rbacv1.RoleBinding{
	rb := rs.RoleBindingMap.Get(ns,name)
	if rb == nil {
		panic("no such rolebinding")
	}
	return rb
}