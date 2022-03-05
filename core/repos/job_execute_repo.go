package repos

import (
	"database/sql"
	"github.com/ixre/go2o/core/domain/interface/job"
	jobImpl "github.com/ixre/go2o/core/domain/job"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/storage"
	"log"
)

var _ job.IJobRepo = new(jobRepositoryImpl)
type jobRepositoryImpl struct{
	_orm orm.Orm
}

func (j *jobRepositoryImpl) CreateJob(data *job.ExecData) job.IJobAggregate {
	return jobImpl.NewJobImpl(j,data)
}


var jobRepoImplMapped = false

// NewJobRepository Create new JobExecDataDao
func NewJobRepository(o orm.Orm, sto storage.Interface) job.IJobRepo {
	if !jobRepoImplMapped{
		_ = o.Mapping(job.ExecData{},"job_exec_data")
		_ = o.Mapping(job.ExecFail{},"job_exec_fail")
		jobRepoImplMapped = true
	}
	return &jobRepositoryImpl{
		_orm:o,
	}
}
// GetExecData Get JobExecData
func (j *jobRepositoryImpl) GetExecData(primary interface{})*job.ExecData{
	e := job.ExecData{}
	err := j._orm.Get(primary,&e)
	if err == nil{
		return &e
	}
	if err != sql.ErrNoRows{
		log.Println("[ Orm][ Error]:",err.Error(),"; Entity:ExecData")
	}
	return nil
}

// GetExecDataBy GetBy JobExecData
func (j *jobRepositoryImpl) GetExecDataBy(where string,v ...interface{})*job.ExecData{
	e := job.ExecData{}
	err := j._orm.GetBy(&e,where,v...)
	if err == nil{
		return &e
	}
	if err != sql.ErrNoRows{
		log.Println("[ Orm][ Error]:",err.Error(),"; Entity:ExecData")
	}
	return nil
}

func (j *jobRepositoryImpl) GetJobByName(name string) job.IJobAggregate {
	v := j.GetExecDataBy("job_name = $1",name)
	if v != nil{
		return j.CreateJob(v)
	}
	return nil
}

// SaveExecData Save JobExecData
func (j *jobRepositoryImpl) SaveExecData(v *job.ExecData)(int,error){
	id,err := orm.Save(j._orm,v,int(v.Id))
	if err != nil && err != sql.ErrNoRows{
		log.Println("[ Orm][ Error]:",err.Error(),"; Entity:ExecData")
	}
	return id,err
}

// GetExecFailBy GetBy 任务执行失败
func (j *jobRepositoryImpl) GetExecFailBy(where string,v ...interface{})*job.ExecFail{
	e := job.ExecFail{}
	err := j._orm.GetBy(&e,where,v...)
	if err == nil{
		return &e
	}
	if err != sql.ErrNoRows{
		log.Println("[ Orm][ Error]:",err.Error(),"; Entity:ExecFail")
	}
	return nil
}

// SaveExecFail Save 任务执行失败
func (j *jobRepositoryImpl) SaveExecFail(v *job.ExecFail)(int,error){
	id,err := orm.Save(j._orm,v,int(v.Id))
	if err != nil && err != sql.ErrNoRows{
		log.Println("[ Orm][ Error]:",err.Error(),"; Entity:ExecFail")
	}
	return id,err
}

