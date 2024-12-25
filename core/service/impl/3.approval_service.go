/**
 * Copyright (C) 2007-2024 fze.NET, All rights reserved.
 *
 * name: 3.approval_service.go
 * author: jarrysix (jarrysix#gmail.com)
 * date: 2024-08-17 18:05:53
 * description: 审批服务
 * history:
 */

package impl

import (
	"context"
	"errors"

	"github.com/ixre/go2o/core/domain/interface/approval"
	"github.com/ixre/go2o/core/service/proto"
)

var _ proto.ApprovalServiceServer = new(ApprovalServiceImpl)

type ApprovalServiceImpl struct {
	_repo approval.IApprovalRepository
	serviceUtil
	proto.UnimplementedApprovalServiceServer
}

func NewApprovalService(repo approval.IApprovalRepository) proto.ApprovalServiceServer {
	return &ApprovalServiceImpl{
		_repo: repo,
	}
}

// Approve implements proto.ApprovalServiceServer.
func (a *ApprovalServiceImpl) Approve(_ context.Context, req *proto.ApprovalApproveRequest) (*proto.TxResult, error) {
	ia := a._repo.GetApproval(int(req.ApprovalId))
	if ia == nil {
		return a.errorV2(errors.New("审批不存在")), nil
	}
	iv := ia.GetApproval()
	if req.ApprovalUserId > 0 && iv.AssignUid != int(req.ApprovalUserId) {
		return a.errorV2(errors.New("没有审批权限")), nil
	}
	err := ia.Approve()
	return a.errorV2(err), nil
}

// Assign implements proto.ApprovalServiceServer.
func (a *ApprovalServiceImpl) Assign(_ context.Context, req *proto.ApprovalAssignRequest) (*proto.TxResult, error) {
	ia := a._repo.GetApproval(int(req.ApprovalId))
	if ia == nil {
		return a.errorV2(errors.New("审批不存在")), nil
	}
	err := ia.Assign(int(req.ApprovalUserId), req.ApprovalUsername)
	return a.errorV2(err), nil
}

// GetFlow implements proto.ApprovalServiceServer.
func (a *ApprovalServiceImpl) GetFlow(_ context.Context, req *proto.ApprovalFlowRequest) (*proto.SApprovalFlow, error) {
	fm := a._repo.FlowManager()
	flow := fm.GetFlow(int(req.FlowId))
	if flow == nil {
		return nil, errors.New("审批流程不存在")
	}
	nodes := make([]*proto.SApprovalNode, 0)
	for _, node := range flow.Nodes {
		nodes = append(nodes, &proto.SApprovalNode{
			NodeId:   int64(node.Id),
			NodeName: node.NodeName,
		})
	}
	return &proto.SApprovalFlow{
		FlowId:   int64(flow.Id),
		FlowName: flow.FlowName,
		FlowDesc: flow.FlowDesc,
		Nodes:    nodes,
	}, nil
}

// Reject implements proto.ApprovalServiceServer.
func (a *ApprovalServiceImpl) Reject(_ context.Context, req *proto.ApprovalRejectRequest) (*proto.TxResult, error) {
	ia := a._repo.GetApproval(int(req.ApprovalId))
	if ia == nil {
		return a.errorV2(errors.New("审批不存在")), nil
	}
	iv := ia.GetApproval()
	if req.ApprovalUserId > 0 && iv.AssignUid != int(req.ApprovalUserId) {
		return a.errorV2(errors.New("没有审批权限")), nil
	}
	err := ia.Reject(req.RejectReason)
	return a.errorV2(err), nil
}
