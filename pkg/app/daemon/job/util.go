package job

import (
	"log"

	"github.com/ixre/go2o/pkg/infrastructure/locker"
	"github.com/ixre/go2o/pkg/inject"
	"github.com/ixre/go2o/pkg/interface/domain/job"
)

func getJob(jobName string) job.IJobAggregate {
	jobRepo := inject.GetJobRepo()
	j := jobRepo.GetJobByName(jobName)
	if j == nil {
		key := "CreateTable_" + jobName
		if locker.Lock(key, 10) {
			defer locker.Unlock(key)
			j = jobRepo.CreateJob(&job.ExecData{
				JobName: jobName,
			})
			if err := j.Save(); err != nil {
				log.Println("[ Job]:", err.Error())
			}
		}
	}
	return j
}
