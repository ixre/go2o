/**
 * Copyright (C) 2009-2022 56X.NET, All rights reserved.
 *
 * name : job_exec_data_service.go
 * author : jarrysix
 * date : 2022/03/06 03:16:21
 * description :
 * history :
 */

package impl

import (
	"context"
	"errors"
	"github.com/ixre/go2o/core/domain/interface/job"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/storage"
)

var _ proto.ExecutionServiceServer = new(executionServiceImpl)

type executionServiceImpl struct {
	repo job.IJobRepo
	s    storage.Interface
	serviceUtil
}

func NewExecDataService(s storage.Interface, repo job.IJobRepo) *executionServiceImpl {
	return &executionServiceImpl{
		s:    s,
		repo: repo,
	}
}

func (j *executionServiceImpl) GetJob(_ context.Context, request *proto.GetJobRequest) (*proto.SExecutionData, error) {
	ij := j.repo.GetJobByName(request.JobName)
	if ij == nil {
		if !request.Create {
			return nil, nil
		}
		ij = j.repo.CreateJob(&job.ExecData{
			JobName:       request.JobName,
			LastExecIndex: 0,
		})
		_ = ij.Save()
	}
	v := ij.GetValue()
	return j.parseExecData(&v), nil
}

func (j *executionServiceImpl) UpdateExecuteCursor(_ context.Context, request *proto.UpdateCursorRequest) (*proto.Result, error) {
	ij := j.repo.GetJobByName(request.JobName)
	if ij == nil {
		return j.error(errors.New("no such job")), nil
	}
	err := ij.UpdateExecCursor(int(request.CursorId))
	if err == nil{
		err = ij.Save()
	}
	if err != nil {
		return j.error(err), nil
	}
	return j.success(nil), nil
}

func (j *executionServiceImpl) AddFail(_ context.Context, request *proto.AddFailRequest) (*proto.Result, error) {
	ij := j.repo.GetJobByName(request.JobName)
	if ij == nil {
		return j.error(errors.New("no such job")), nil
	}
	err := ij.AddFail(int(request.CursorId))
	if err != nil {
		return j.error(err), nil
	}
	return j.success(nil), nil
}

func (j *executionServiceImpl) parseExecData(v *job.ExecData) *proto.SExecutionData {
	return &proto.SExecutionData{
		Id:            v.Id,
		JobName:       v.JobName,
		LastExecuteCursorId: v.LastExecIndex,
		LastExecuteTime:  v.LastExecTime,
	}
}


// RejoinQueue 保存重新加入队列
func (j *executionServiceImpl) RejoinQueue(_ context.Context, r *proto.RejoinQueueRequest) (*proto.RejoinQueueResponse, error) {
	job := j.repo.CreateJob(&job.ExecData{
		JobName: r.JobName,
	})
	id,err := job.RejoinQueue(r.RelateId,r.RelateData)
	ret := &proto.RejoinQueueResponse{
		QueueId: int64(id),
	}
	if err != nil{
		ret.ErrCode = 1
		ret.ErrMsg = err.Error()
	}
	return ret,nil
}
