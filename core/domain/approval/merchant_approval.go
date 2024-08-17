/**
 * Copyright (C) 2007-2024 fze.NET, All rights reserved.
 *
 * name: merchant_approval.go
 * author: jarrysix (jarrysix#gmail.com)
 * date: 2024-08-17 10:35:26
 * description: 员工转商户审批业务
 * history:
 */

package approval

import (
	"github.com/ixre/go2o/core/domain/interface/approval"
)

var _ approval.IApprovalAggregateRoot = new(staffTransferApprovalImpl)

type staffTransferApprovalImpl struct {
	*ApprovalImpl
}

func NewStaffTransferApproval(value *approval.Approval, repo approval.IApprovalRepository) approval.IApprovalAggregateRoot {
	s := &staffTransferApprovalImpl{}
	s._value = value
	s._repo = repo
	return s
}

func (s *staffTransferApprovalImpl) Process(nodeKey string, tx *approval.ApprovalLog) error {
	if nodeKey == "aggree" {
		// 原商户同意
	}
	if nodeKey == "finish" {
		//
	}
	return nil
}
