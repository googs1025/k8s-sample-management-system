package models

type Deployment struct {
	Name string
	NameSpace string
	Replicas [3]int32   //3个值，分别是总副本数，可用副本数 ，不可用副本数
	Images string
	IsComplete bool //是否完成
	Message string // 显示错误信息
	CreateTime string
	Pods []*Pod
}

