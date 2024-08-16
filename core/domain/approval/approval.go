/**
 * Copyright (C) 2007-2024 fze.NET, All rights reserved.
 *
 * name: approval.go
 * author: jarrysix (jarrysix#gmail.com)
 * date: 2024-08-16 21:48:20
 * description: 审批聚合实现
 * history:
 */

package approval

import (
	"errors"
	"time"

	"github.com/ixre/go2o/core/domain/interface/approval"
	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/go2o/core/infrastructure/fw/types"
)

var _ approval.IApprovalAggregateRoot = new(ApprovalImpl)

type ApprovalImpl struct {
	_value *approval.Approval
	_repo  approval.IApprovalRepository
	_flow  *approval.ApprovalFlow
}

// GetAggregateRootId implements approval.IApprovalAggregateRoot.
func (a *ApprovalImpl) GetAggregateRootId() int {
	return a._value.Id
}

func (a *ApprovalImpl) GetCurrentNodeLog() (*approval.ApprovalLog, error) {
	if a._value.FinalStatus != approval.FinalAwaitStatus {
		return nil, errors.New("审批单已完成审核")
	}
	v := a._repo.GetCurrentNodeLog(a.GetAggregateRootId())
	if v == nil {
		return nil, errors.New("审核单缺少审核节点或已完成审核")
	}
	return v, nil
}

func (a *ApprovalImpl) Approve() error {
	// 实现审批逻辑
	return nil
}

func (a *ApprovalImpl) ChangeAssign(uid int, name string) error {

	// 实现更改分配逻辑
	return nil
}

func (a *ApprovalImpl) FlowId() int {
	return a._value.FlowId
}

func (a *ApprovalImpl) Flow() *approval.ApprovalFlow {
	if a._flow == nil {
		a._flow = a._repo.FlowManager().GetFlow(a.FlowId())
	}
	return a._flow
}

func (a *ApprovalImpl) GetApproval() *approval.Approval {
	return types.DeepClone(a._value)
}

func (a *ApprovalImpl) Reject(remark string) error {
	// 实现拒绝逻辑
	return nil
}

func (a *ApprovalImpl) Save() error {
	if a.GetAggregateRootId() <= 0 {
		return a.submitApproval()
	}
	// 实现保存逻辑
	return nil
}

// submitApproval 提交审批单
func (a *ApprovalImpl) submitApproval() error {
	flow := a.Flow()
	if flow == nil {
		return errors.New("工作流未定义")
	}
	if a._value.NodeId > 0 {
		return errors.New("提交审批单不能包含节点信息")
	}
	// 生成交易号
	txNo := domain.NewTradeNo(0, 0)
	if flow.TxPrefix != "" {
		txNo = flow.TxPrefix + "-" + txNo
	}
	a._value.ApprovalNo = txNo
	a._value.CreateTime = int(time.Now().Unix())
	a._value.UpdateTime = int(time.Now().Unix())
	v, err := a._repo.Save(a._value)
	if err == nil {
		a._value.Id = v.Id
	}
	return err
}
