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
	"github.com/ixre/go2o/core/domain/interface/domain/enum"
	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/go2o/core/infrastructure/fw/types"
)

var _ approval.IApprovalAggregateRoot = new(ApprovalImpl)

type ApprovalImpl struct {
	_value *approval.Approval
	_repo  approval.IApprovalRepository
	_flow  *approval.ApprovalFlow
}

func NewApproval(value *approval.Approval, repo approval.IApprovalRepository) approval.IApprovalAggregateRoot {
	return &ApprovalImpl{
		_value: value,
		_repo:  repo,
	}
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

// 获取节点
func (a *ApprovalImpl) getNode(nodeId int) *approval.ApprovalFlowNode {
	flow := a.Flow()
	if flow == nil {
		return nil
	}
	for _, node := range flow.Nodes {
		if node.Id == nodeId {
			return node
		}
	}
	return nil
}

// Approve 审批通过
func (a *ApprovalImpl) Approve() error {
	current, err := a.GetCurrentNodeLog()
	if err != nil {
		return err
	}
	node := a.getNode(current.NodeId)
	if node == nil {
		return errors.New("找不到节点")
	}
	// 更改审批日志信息
	current.ApprovalStatus = int(enum.ReviewApproved)
	current.ApprovalTime = int(time.Now().Unix())
	_, err = a._repo.GetLogRepo().Save(current)
	if err == nil {
		// 处理工作流
		err = a.Process(node.NodeKey, current)
		if err != nil {
			return err
		}
		// 更新审批终态
		if node.NodeType == approval.NodeTypeEnd {
			a._value.FinalStatus = int(enum.ReviewApproved)
			if err = a.Save(); err != nil {
				return err
			}
		}
		// 不是结束节点,则更新到下一个节点
		if !a.IsFinal() {
			err = a.toNextNode()
		}
	}
	return err
}

// Assign 更改分配人
func (a *ApprovalImpl) Assign(uid int, name string) error {
	current, err := a.GetCurrentNodeLog()
	if err != nil {
		return err
	}
	current.AssignUid = uid
	current.AssignName = name
	_, err = a._repo.GetLogRepo().Save(current)
	if err == nil {
		a._value.AssignUid = uid
		a._value.AssignName = name
		err = a.Save()
	}
	return err
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
	current, err := a.GetCurrentNodeLog()
	if err != nil {
		return err
	}
	node := a.getNode(current.NodeId)
	if node == nil {
		return errors.New("找不到节点")
	}
	// 更改审批日志信息
	current.ApprovalStatus = int(enum.ReviewRejected)
	current.ApprovalRemark = remark
	current.ApprovalTime = int(time.Now().Unix())
	_, err = a._repo.GetLogRepo().Save(current)
	if err == nil {
		// 处理工作流
		err = a.Process(node.NodeKey, current)
		if err != nil {
			return err
		}
		// 更新审批终态
		a._value.FinalStatus = int(enum.ReviewRejected)
		err = a.Save()
	}
	return err
}

func (a *ApprovalImpl) Save() error {
	if a.GetAggregateRootId() <= 0 {
		return a.submitApproval()
	}
	_, err := a._repo.Save(a._value)
	return err
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
		err = a.toNextNode()
	}
	return err
}

func (a *ApprovalImpl) toNextNode() error {
	if a.IsFinal() {
		// 已完成审批
		return nil
	}
	node := a.getNextNode()
	if node == nil {
		return errors.New("找不到下一个节点")
	}
	a._value.NodeId = node.Id
	a._value.UpdateTime = int(time.Now().Unix())
	_, err := a._repo.Save(a._value)
	if err == nil {
		approvalLog := &approval.ApprovalLog{
			Id:             0,
			ApprovalId:     a.GetAggregateRootId(),
			NodeId:         node.Id,
			NodeName:       node.NodeName,
			AssignUid:      0,
			AssignName:     "",
			ApprovalStatus: int(enum.ReviewPending),
			ApprovalRemark: "",
			ApprovalTime:   0,
			CreateTime:     int(time.Now().Unix()),
		}
		_, err = a._repo.GetLogRepo().Save(approvalLog)
	}
	return err
}

// 获取下一个节点
func (a *ApprovalImpl) getNextNode() *approval.ApprovalFlowNode {
	flow := a.Flow()
	for i, node := range flow.Nodes {
		if a._value.NodeId <= 0 && node.NodeType == approval.NodeTypeStart {
			// 如果还没有节点，则返回开始节点
			return node
		}
		if node.Id != a._value.NodeId {
			// 不是当前节点
			continue
		}
		if i+1 < len(flow.Nodes) {
			// 返回下一个节点
			return flow.Nodes[i+1]
		}
		break
	}
	return nil
}

// / IsFinal implements approval.IApprovalAggregateRoot.
func (a *ApprovalImpl) IsFinal() bool {
	return a._value.FinalStatus != approval.FinalAwaitStatus
}

// Process implements approval.IApprovalAggregateRoot.
func (a *ApprovalImpl) Process(nodeKey string, tx *approval.ApprovalLog) error {
	panic("工作流处理程序应用实现类进行处理")
}
