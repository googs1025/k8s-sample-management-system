package models

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CronJob struct {
	Name string
	NameSpace string
	Images string
	LastScheduleTime metav1.Time
	CreateTime string
	Pods []*Pod
}
