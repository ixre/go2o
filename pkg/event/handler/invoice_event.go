/**
 * Copyright (C) 2007-2024 fze.NET, All rights reserved.
 *
 * name: invoice_event.go
 * author: jarrysix (jarrysix@gmail.com)
 * date: 2024-09-22 16:26:14
 * description: 发票事件处理
 * history:
 */

package handler

type InvoiceEventHandler struct {
}

func NewInvoiceEventHandler() *InvoiceEventHandler {
	return &InvoiceEventHandler{}
}

// HandleInvoiceRevertEvent 处理发票撤销事件
func (i *InvoiceEventHandler) HandleInvoiceRevertEvent(event interface{}) {
	// e := event.(*invoice.InvoiceRevertEvent)
	// 撤销发票应在具体的实现中订阅事件并处理, 这里不做任何处理
	// 通常订阅事件的场景为： 撤销除商城订单以外的自定义订单的发票, 需将订单标记为未开票
}
