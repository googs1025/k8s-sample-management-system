package models

type Job struct {
	Name string
	NameSpace string
	Images string
	IsComplete bool //是否完成
	Message string // 显示错误信息
	CreateTime string
	Pods []*Pod
}
