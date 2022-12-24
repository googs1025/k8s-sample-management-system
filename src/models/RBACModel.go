package models

import (
	rbacv1 "k8s.io/api/rbac/v1"
)

type RoleModel struct {
	Name string
	NameSpace string
	CreateTime string
}

type RoleBindingModel struct {
	Name string
	NameSpace string
	CreateTime string
	RoleRef rbacv1.RoleRef
	Subject []rbacv1.Subject  //包含了 绑定用户 数据
}
