/**
 * Copyright (C) 2007-2024 fze.NET, All rights reserved.
 *
 * name: merchant_event.go
 * author: jarrysix (jarrysix#gmail.com)
 * date: 2024-09-13 00:15:38
 * description: 商户事件处理
 * history:
 */

package handler

type MerchantEventHandler struct {
}

func NewMerchantEventHandler() *MerchantEventHandler {
	return &MerchantEventHandler{}
}

// HandleStaffRequireImInitEvent 处理员工IM初始化事件
func (m *MerchantEventHandler) HandleStaffRequireImInitEvent(event interface{}) {
	// e := event.(*staff.StaffRequireImInitEvent)
	// 初始化员工IM应在具体的实现中订阅事件并处理, 这里只是为了展示，不做任何处理
}

// HandleMerchantBillSettleEvent 处理商户账单结算事件
func (m *MerchantEventHandler) HandleMerchantBillSettleEvent(event interface{}) {
	// e := event.(*merchant.MerchantBillSettleEvent)
	// 结算账单应在具体的实现中订阅事件并处理, 这里只是为了展示，不做任何处理
}

// HandleStaffTransferApprovedEvent 处理员工转移审批通过事件
func (m *MerchantEventHandler) HandleStaffTransferApprovedEvent(event interface{}) {
	//e := event.(*staff.StaffTransferApprovedEvent)
	// 员工转移审批通过事件应在具体的实现中订阅事件并处理, 这里只是为了展示，不做任何处理
}
