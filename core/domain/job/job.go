package job

import (
	"errors"
	"github.com/ixre/go2o/core/domain/interface/job"
	"time"
)

var _ job.IJobAggregate = new(jobImpl)

type jobImpl struct {
	repo  job.IJobRepo
	value *job.ExecData
}


func NewJobImpl(repo job.IJobRepo, value *job.ExecData) job.IJobAggregate {
	return &jobImpl{
		repo:  repo,
		value: value,
	}
}

func (j jobImpl) GetAggregateRootId() int64 {
	return j.value.Id
}

func (j jobImpl) GetValue() job.ExecData {
	return *j.value
}

func (j jobImpl) SetValue(data job.ExecData) error {
	if data.Id <= 0 {
		j.value.JobName = data.JobName
	}
	j.value.LastExecIndex = data.LastExecIndex
	j.value.LastExecTime = data.LastExecTime
	return nil
}

func (j jobImpl) AddFail(recordId int) error {
	if j.GetAggregateRootId() == 0 {
		return errors.New("job not exists")
	}
	v := &job.ExecFail{
		JobId:      j.GetAggregateRootId(),
		JobDataId:  int64(recordId),
		RetryCount: 0,
		CreateTime: time.Now().Unix(),
		RetryTime:  0,
	}
	id, err := j.repo.SaveExecFail(v)
	if err == nil {
		v.Id = int64(id)
	}
	return err
}

func (j jobImpl) UpdateExecCursor(id int) error {
	if id <= 0 {
		return errors.New("id error")
	}
	j.value.LastExecIndex = int64(id)
	j.value.LastExecTime = time.Now().Unix()
	return nil
}

func (j jobImpl) Save() error {
	id, err := j.repo.SaveExecData(j.value)
	if j.GetAggregateRootId() == 0 && err == nil {
		j.value.Id = int64(id)
	}
	return err
}


func (j jobImpl) RejoinQueue(relateId int64, relateData string) (int,error){
	if j.value == nil || len(j.value.JobName) == 0{
		return 0,errors.New("no such job")
	}
	if relateId <= 0 && len(relateData) ==  0{
		return 0,errors.New("relate id or data is empty")
	}
	v := &job.ExecRequeue{
		Id:         0,
		QueueName:  j.value.JobName,
		RelateId:   relateId,
		RelateData: relateData,
		CreateTime: time.Now().Unix(),
	}
	return j.repo.SaveRequeue(v)
}
