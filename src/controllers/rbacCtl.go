package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/shenyisyn/goft-gin/goft"
	"k8s-Management-System/src/services"
)

type RBACCtl struct {
	RoleService *services.RoleService `inject:"-"`

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
}