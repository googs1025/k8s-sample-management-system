package services

import "k8s-Management-System/src/models"

//@Service
type RoleService struct {
	RoleMap *RoleMap  `inject:"-"`
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