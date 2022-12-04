package models

type CronJob struct {
	Name string
	NameSpace string
	Images string
	LastScheduleTime string
	CreateTime string
	Pods []*Pod
}
