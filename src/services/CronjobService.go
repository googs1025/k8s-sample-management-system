package services

import (
	"github.com/shenyisyn/goft-gin/goft"
	"k8s-Management-System/src/models"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//@Service
type CronJobService struct {
	CronJobMap *CronJobMap   `inject:"-"`
	Common *CommonService `inject:"-"`
}

func NewCronJobService() *CronJobService {
	return &CronJobService{}
}


func (*CronJobService) getCronJobLastScheduleTime(cronjob *batchv1beta1.CronJob) metav1.Time {



	return *cronjob.Status.LastScheduleTime

}


func (cj *CronJobService) ListAll(namespace string) (res []*models.CronJob) {

	cronJobList, err := cj.CronJobMap.ListCronJobByNamespace(namespace)
	goft.Error(err)

	for _, cjj := range cronJobList {
		res = append(res, &models.CronJob{
			Name: cjj.Name,
			NameSpace: cjj.Namespace,
			Images: cj.Common.GetCronJobImages(*cjj),
			LastScheduleTime: cj.getCronJobLastScheduleTime(cjj).Format("2006-01-02 15:04:05"),
			CreateTime: cjj.CreationTimestamp.Format("2006-01-02 15:04:05"),
		})

	}

	return

}
