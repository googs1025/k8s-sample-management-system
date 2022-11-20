package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/shenyisyn/goft-gin/goft"
)

/*
	对 login logout info 进行静态数据的返回。
 */

type UserCtl struct {

}

func NewUserCtl() *UserCtl {
	return &UserCtl{}
}

func(u *UserCtl) login(c *gin.Context) goft.Json  {
	return gin.H{
		"code":20000,
		"data":gin.H{
			"token":"admin-token",
		},
	}
}

func(u *UserCtl) logout(c *gin.Context)  goft.Json  {
	return gin.H{
		"code":20000,
		"data":"success",
	}
}

func(u *UserCtl) info(c *gin.Context) string{
	c.Header("Content-type","application/json")
	return `{"code":20000,"data":{"roles":["admin"],
		"introduction":"I am a super administrator","avatar":"https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif","name":"Super Admin"}}`

}
func(u *UserCtl)  Build(goft *goft.Goft){
	goft.Handle("POST","/vue-admin-template/user/login", u.login)
	goft.Handle("POST","/vue-admin-template/user/logout", u.logout)
	goft.Handle("GET","/vue-admin-template/user/info", u.info)
}
func(*UserCtl) Name() string{
	return "UserCtl"
}

