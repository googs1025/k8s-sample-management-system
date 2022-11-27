package services

import (
	"github.com/shenyisyn/goft-gin/goft"
	"k8s-Management-System/src/models"
	batchv1 "k8s.io/api/batch/v1"
)

//@Service
type JobService struct {
	JobMap *JobMap   `inject:"-"`
	Common *CommonService `inject:"-"`
}

func NewJobService() *JobService {
	return &JobService{}
}


func (*JobService) getJobCondition(job *batchv1.Job) string {

	for _, item := range job.Status.Conditions {
		if string(item.Type) == "Available" && string(item.Status) != "True" {
			return item.Message
		}
	}

	return ""

}

func (*JobService) getJobIsComplete(job *batchv1.Job) bool {
	return job.Status.Succeeded == 1
}

func (j *JobService) ListAll(namespace string) (res []*models.Job) {

	jobList, err := j.JobMap.ListJobByNamespace(namespace)
	goft.Error(err)

	for _, jj := range jobList {
		res = append(res, &models.Job{
			Name: jj.Name,
			NameSpace: jj.Namespace,
			Images: j.Common.GetJobImages(*jj),
			IsComplete: j.getJobIsComplete(jj),
			Message: j.getJobCondition(jj),
			CreateTime: jj.CreationTimestamp.Format("2006-01-02 15:04:05"),
		})

	}

	return

}