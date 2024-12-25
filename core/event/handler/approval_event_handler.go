/**
 * Copyright (C) 2007-2024 fze.NET, All rights reserved.
 *
 * name: approval_event_handler.go
 * author: jarrysix (jarrysix#gmail.com)
 * date: 2024-08-17 14:40:08
 * description: 处理审批事件
 * history:
 */

package handler

import (
	"fmt"

	"github.com/ixre/go2o/core/domain/interface/approval"
	"github.com/ixre/go2o/core/domain/interface/merchant"
	"github.com/ixre/go2o/core/infrastructure/logger"
)

func (h *EventHandler) OnApprovalProcess(data interface{}) {
	event := data.(*approval.ApprovalProcessEvent)
	ap := event.Approval
	if ap.FlowId() == approval.FlowStaffTransfer {
		processStaffTransferEvent(h, event)
		return
	}
	logger.Error("approval process event not support: %d", event.Approval.FlowId())
}

// 处理员工转移事件
func processStaffTransferEvent(h *EventHandler, event *approval.ApprovalProcessEvent) {
	trans := h._staffRepo.TransferRepo().Get(event.Approval.GetApproval().BizId)
	if trans == nil {
		logger.Error("staff transfer not found: %d", event.Approval.GetApproval().BizId)
		panic(fmt.Errorf("staff transfer not found: %d", event.Approval.GetApproval().BizId))
	}
	ic := h._mchRepo.CreateMerchant(&merchant.Merchant{
		Id: trans.OriginMchId,
	})
	err := ic.EmployeeManager().TransferApproval(trans, event)
	if err != nil {
		panic(err)
	}
}
