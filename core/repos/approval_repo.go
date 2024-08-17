/**
 * Copyright (C) 2007-2024 fze.NET, All rights reserved.
 *
 * name: approval_repo.go
 * author: jarrysix (jarrysix#gmail.com)
 * date: 2024-08-17 12:24:58
 * description: 审批仓储
 * history:
 */

package repos

import (
	"time"

	approvalImpl "github.com/ixre/go2o/core/domain/approval"
	"github.com/ixre/go2o/core/domain/interface/approval"
	"github.com/ixre/go2o/core/domain/interface/domain/enum"
	"github.com/ixre/go2o/core/infrastructure/fw"
)

var _ approval.IApprovalRepository = new(approvalRepositoryImpl)

type approvalRepositoryImpl struct {
	fw.BaseRepository[approval.Approval]
	logRepo fw.Repository[approval.ApprovalLog]
}

// Create implements approval.IApprovalRepository.
// Subtle: this method shadows the method (BaseRepository).Create of approvalRepositoryImpl.BaseRepository.
func (a *approvalRepositoryImpl) Create(flowId int, bizId int) approval.IApprovalAggregateRoot {
	return approvalImpl.NewApproval(&approval.Approval{
		Id:          0,
		ApprovalNo:  "",
		FlowId:      flowId,
		BizId:       bizId,
		NodeId:      0,
		AssignUid:   0,
		AssignName:  "",
		FinalStatus: 0,
		CreateTime:  int(time.Now().Unix()),
		UpdateTime:  0,
	}, a)
}

// FlowManager implements approval.IApprovalRepository.
func (a *approvalRepositoryImpl) FlowManager() approval.IFlowManager {
	return approvalImpl.NewFlowManager(a)
}

// GetCurrentNodeLog implements approval.IApprovalRepository.
func (a *approvalRepositoryImpl) GetCurrentNodeLog(approvalId int) *approval.ApprovalLog {
	return a.GetLogRepo().FindBy("approval_id = ? AND approval_status = ?", approvalId, enum.ReviewPending)
}

// GetApproval implements approval.IApprovalRepository.
func (a *approvalRepositoryImpl) GetApproval(id int) approval.IApprovalAggregateRoot {
	return approvalImpl.NewApproval(a.Get(id), a)
}

// GetLogRepo implements approval.IApprovalRepository.
func (a *approvalRepositoryImpl) GetLogRepo() fw.Repository[approval.ApprovalLog] {
	if a.logRepo == nil {
		a.logRepo = &fw.BaseRepository[approval.ApprovalLog]{
			ORM: a.ORM,
		}
	}
	return a.logRepo
}

func NewApprovalRepository(o fw.ORM) approval.IApprovalRepository {
	return &approvalRepositoryImpl{
		BaseRepository: fw.BaseRepository[approval.Approval]{
			ORM: o,
		},
		logRepo: nil,
	}
}
