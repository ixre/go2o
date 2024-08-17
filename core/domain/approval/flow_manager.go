/**
 * Copyright (C) 2007-2024 fze.NET, All rights reserved.
 *
 * name: flow_manager.go
 * author: jarrysix (jarrysix#gmail.com)
 * date: 2024-08-17 10:35:49
 * description: 工作流管理器
 * history:
 */

package approval

import "github.com/ixre/go2o/core/domain/interface/approval"

var _ approval.IFlowManager = new(flowManagerImpl)

type flowManagerImpl struct {
}

func NewFlowManager() approval.IFlowManager {
	return &flowManagerImpl{}
}

// CreateFlow implements approval.IFlowManager.
func (f *flowManagerImpl) CreateFlow(name string, desc string, nodes []*approval.ApprovalFlow) (int, error) {
	panic("unimplemented")
}

// GetFlow implements approval.IFlowManager.
func (f *flowManagerImpl) GetFlow(id int) *approval.ApprovalFlow {
	if id == 101 {
		return &staffTransferApprovalFlow
	}
	return nil
}
