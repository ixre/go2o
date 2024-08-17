/**
 * Copyright (C) 2007-2024 fze.NET, All rights reserved.
 *
 * name: approval.go
 * author: jarrysix (jarrysix#gmail.com)
 * date: 2024-08-16 21:48:42
 * description: 审批聚合根
 * history:
 */

package approval

import (
	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/go2o/core/infrastructure/fw"
)

var (
	// 待审批
	FinalAwaitStatus = 0
	// 审批通过
	FinalPassStatus = 1
	// 审批拒绝
	FinalRejectStatus = 2
)

var (
	// 开始节点
	NodeTypeStart = 1
	// 结束节点
	NodeTypeEnd = 2
	// 其他节点
	NodeTypeOther = 3
)

type (
	// IApprovalAggregateRoot 审批聚合根
	IApprovalAggregateRoot interface {
		domain.IAggregateRoot
		// Save 保存审批
		Save() error
		// FlowId 获取工作流编号
		FlowId() int
		// 获取审批
		GetApproval() *Approval
		// 通过
		Approve() error
		// 拒绝
		Reject(remark string) error
		// 是否最终状态
		IsFinal() bool
		// 获取工作流
		Flow() *ApprovalFlow
		// 处理节点,在审批后会进行自动调用, 实现审批业务需重写
		Process(nodeKey string, tx *ApprovalLog) error
		// 分配审批人,当节点审批后切换到下个节点, 需分配审批人
		Assign(uid int, name string) error
	}

	// IFlowManager 工作流管理器
	IFlowManager interface {
		// GetFlow 获取工作流
		GetFlow(id int) *ApprovalFlow
		// CreateFlow 创建工作流
		CreateFlow(name, desc string, nodes []*ApprovalFlow) (int, error)
	}

	// IApprovalRepository 审批表仓储
	IApprovalRepository interface {
		fw.Repository[Approval]
		// GetLogRepo 获取审核日志仓储
		GetLogRepo() fw.Repository[ApprovalLog]
		// FlowManager 工作流管理器
		FlowManager() IFlowManager
		// GetCurrentNodeLog 获取当前节点审核日志
		GetCurrentNodeLog(approvalId int) *ApprovalLog
	}
)

// Approval 审批表
type Approval struct {
	// 编号
	Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" auto:"yes" bson:"id"`
	// 审批流水号
	ApprovalNo string `json:"approvalNo" db:"approval_no" gorm:"column:approval_no" bson:"approvalNo"`
	// 工作流编号
	FlowId int `json:"flowId" db:"flow_id" gorm:"column:flow_id" bson:"flowId"`
	// 当前节点编号
	NodeId int `json:"nodeId" db:"node_id" gorm:"column:node_id" bson:"nodeId"`
	// 审批人
	AssignUid int `json:"assignUid" db:"assign_uid" gorm:"column:assign_uid" bson:"assignUid"`
	// 审批人名称
	AssignName string `json:"assignName" db:"assign_name" gorm:"column:assign_name" bson:"assignName"`
	// 最终状态,  0: 审核中  1: 已通过  2:不通过
	FinalStatus int `json:"finalStatus" db:"final_status" gorm:"column:final_status" bson:"finalStatus"`
	// 创建时间
	CreateTime int `json:"createTime" db:"create_time" gorm:"column:create_time" bson:"createTime"`
	// 更新时间
	UpdateTime int `json:"updateTime" db:"update_time" gorm:"column:update_time" bson:"updateTime"`
}

func (a Approval) TableName() string {
	return "approval"
}

// ApprovalLog 审核日志
type ApprovalLog struct {
	// 编号
	Id int `json:"id" db:"id" gorm:"column:id" bson:"id"`
	// 审批编号
	ApprovalId int `json:"approvalId" db:"approval_id" gorm:"column:approval_id" bson:"approvalId"`
	// 节点编号
	NodeId int `json:"nodeId" db:"node_id" gorm:"column:node_id" bson:"nodeId"`
	// 节点名称
	NodeName string `json:"nodeName" db:"node_name" gorm:"column:node_name" bson:"nodeName"`
	// 审批人编号
	AssignUid int `json:"assignUid" db:"assign_uid" gorm:"column:assign_uid" bson:"assignUid"`
	// 审批人名称
	AssignName string `json:"assignName" db:"assign_name" gorm:"column:assign_name" bson:"assignName"`
	// 审核状态
	ApprovalStatus int `json:"approvalStatus" db:"approval_status" gorm:"column:approval_status" bson:"approvalStatus"`
	// 审核备注
	ApprovalRemark string `json:"approvalRemark" db:"approval_remark" gorm:"column:approval_remark" bson:"approvalRemark"`
	// 审核时间
	ApprovalTime int `json:"approvalTime" db:"approval_time" gorm:"column:approval_time" bson:"approvalTime"`
	// 创建时间
	CreateTime int `json:"createTime" db:"create_time" gorm:"column:create_time" bson:"createTime"`
}

func (a ApprovalLog) TableName() string {
	return "approval_log"
}

// ApprovalFlow ApprovalFlow
type ApprovalFlow struct {
	// Id
	Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" auto:"yes" bson:"id"`
	// 工作流名称
	FlowName string `json:"flowName" db:"flow_name" gorm:"column:flow_name" bson:"flowName"`
	// 工作流描述
	FlowDesc string `json:"flowDesc" db:"flow_desc" gorm:"column:flow_desc" bson:"flowDesc"`
	// 交易号前缀
	TxPrefix string `json:"txPrefix" db:"tx_prefix" gorm:"column:tx_prefix" bson:"txPrefix"`
	// 节点
	Nodes []*ApprovalFlowNode `json:"nodes" db:"-" gorm:"-:all" bson:"nodes"`
}

func (a ApprovalFlow) TableName() string {
	return "approval_flow"
}

// ApprovalFlowNode ApprovalFlowNode
type ApprovalFlowNode struct {
	// Id
	Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" auto:"yes" bson:"id"`
	// 工作流编号
	FlowId int `json:"flowId" db:"flow_id" gorm:"column:flow_id" bson:"flowId"`
	// 节点KEY
	NodeKey string `json:"nodeKey" db:"node_key" gorm:"column:node_key" bson:"nodeKey"`
	// 节点类型 1:起始节点   2: 结束节点   3: 其他节点
	NodeType int `json:"nodeType" db:"node_type" gorm:"column:node_type" bson:"nodeType"`
	// 节点名称
	NodeName string `json:"nodeName" db:"node_name" gorm:"column:node_name" bson:"nodeName"`
	// 节点描述
	NodeDesc string `json:"nodeDesc" db:"node_desc" gorm:"column:node_desc" bson:"nodeDesc"`
}

func (a ApprovalFlowNode) TableName() string {
	return "approval_flow_node"
}
