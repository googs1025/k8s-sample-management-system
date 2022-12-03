package models

type StatefulSet struct {
	Name string
	NameSpace string
	Replicas int32   //3个值，分别是总副本数，可用副本数 ，不可用副本数
	Images string
	CreateTime string
	Pods []*Pod
}
