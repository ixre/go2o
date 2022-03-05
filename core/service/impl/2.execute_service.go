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
	context2 "golang.org/x/net/context"
)

var _ proto.ExecuteServiceServer = new(executeServiceImpl)

type executeServiceImpl struct {
	repo job.IJobRepo
	s    storage.Interface
	serviceUtil
}

func NewExecDataService(s storage.Interface, repo job.IJobRepo) *executeServiceImpl {
	return &executeServiceImpl{
		s:    s,
		repo: repo,
	}
}

func (j *executeServiceImpl) GetJob(_ context.Context, request *proto.GetJobRequest) (*proto.SExecData, error) {
	ij := j.repo.GetJobByName(request.JobName)
	if ij == nil {
		if !request.Create {
			return nil, nil
		}
		ij = j.repo.CreateJob(&job.ExecData{
			JobName:       request.JobName,
			LastExecIndex: 0,
		})
	}
	v := ij.GetValue()
	return j.parseExecData(&v), nil
}

func (j *executeServiceImpl) UpdateExecCursor(_ context.Context, request *proto.UpdateCursorRequest) (*proto.Result, error) {
	ij := j.repo.GetJobByName(request.JobName)
	if ij == nil {
		return j.error(errors.New("no such job")), nil
	}
	err := ij.UpdateExecCursor(int(request.RecordId))
	if err != nil {
		return j.error(err), nil
	}
	return j.success(nil), nil
}

func (j *executeServiceImpl) AddFail(_ context2.Context, request *proto.AddFailRequest) (*proto.Result, error) {
	ij := j.repo.GetJobByName(request.JobName)
	if ij == nil {
		return j.error(errors.New("no such job")), nil
	}
	err := ij.AddFail(int(request.RecordId))
	if err != nil {
		return j.error(err), nil
	}
	return j.success(nil), nil
}

func (j *executeServiceImpl) parseExecData(v *job.ExecData) *proto.SExecData {
	return &proto.SExecData{
		Id:            v.Id,
		JobName:       v.JobName,
		LastExecIndex: v.LastExecIndex,
		LastExecTime:  v.LastExecTime,
	}
}
