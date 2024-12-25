/**
 * Copyright (C) 2007-2024 fze.NET, All rights reserved.
 *
 * name: workflow.go
 * author: jarrysix (jarrysix#gmail.com)
 * date: 2024-08-17 10:34:38
 * description: 定义审批工作流
 * history:
 */

package approval

import "github.com/ixre/go2o/core/domain/interface/approval"

var staffTransferApprovalFlow = approval.ApprovalFlow{
	Id:       101,
	FlowName: "员工转商户",
	FlowDesc: "商户员工转移到其他商户",
	TxPrefix: "MSTM-",
	Nodes: []*approval.ApprovalFlowNode{
		{
			Id:       1,
			NodeKey:  "aggree",
			NodeName: "商户同意",
			NodeDesc: "原商户同意转出",
			NodeType: 1,
		},
		{
			Id:       2,
			NodeKey:  "finish",
			NodeName: "新商户同意",
			NodeDesc: "新商户同意转入",
			NodeType: 2,
		},
	},
}
