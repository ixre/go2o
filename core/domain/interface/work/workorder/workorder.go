/**
 * Copyright (C) 2009-2024 56X.NET, All rights reserved.
 *
 * name : model_gen.go
 * author : jarrysix
 * date : 2024/07/20 16:22:23
 * description :
 * history :
 */
package workorder

import (
	"github.com/ixre/go2o/core/domain"
	"github.com/ixre/go2o/core/infrastructure/fw"
)

const (
	StatusPending    = 1 // 待处理
	StatusProcessing = 2 // 处理中
	StatusFinished   = 3 // 已完结
)

const (
	FlagUserClosed = 1 // 用户关闭
)

const (
	ClassSuggest = 1 // 建议
	ClassAppeal  = 2 // 申诉
)

type (
	// 工单聚合根
	IWorkorderAggregateRoot interface {
		domain.IAggregateRoot
		// 获取工单数
		Value() *Workorder
		// 提交工单
		Submit() error
		// 分配客服
		AllocateAgentId(userId int) error
		// 完结
		Finish() error
		// 用户关闭工单
		Close() error
		// 评价
		Apprise(isUsefully bool, rank int, apprise string) error
		// 提交回复
		SubmitComment(content string, isReplay bool, refCommentId int) error
	}
)

// IWorkorderRepo 工单仓储
type IWorkorderRepo interface {
	fw.Repository[Workorder]
	// 评论仓储
	CommentRepo() fw.Repository[WorkorderComment]
	// 创建工单
	CreateWorkorder(value *Workorder) IWorkorderAggregateRoot
	// 获取工单
	GetWorkorder(id int) IWorkorderAggregateRoot
}

// Workorder Workorder
type Workorder struct {
	// 编号
	Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" auto:"yes" bson:"id"`
	// 工单号
	OrderNo string `json:"orderNo" db:"order_no" gorm:"column:order_no" bson:"orderNo"`
	// 会员编号
	MemberId int `json:"memberId" db:"member_id" gorm:"column:member_id" bson:"memberId"`
	// 反馈类型, 1: 建议 2:申诉
	ClassId int `json:"classId" db:"class_id" gorm:"column:class_id" bson:"classId"`
	// 关联商户
	MchId int `json:"mchId" db:"mch_id" gorm:"column:mch_id" bson:"mchId"`
	// 标志, 1:用户关闭
	Flag int `json:"flag" db:"flag" gorm:"column:flag" bson:"flag"`
	// 关联业务, 如:CHARGE:2014050060
	Wip string `json:"wip" db:"wip" gorm:"column:wip" bson:"wip"`
	// Subject
	Subject string `json:"subject" db:"subject" gorm:"column:subject" bson:"subject"`
	// 投诉内容
	Content string `json:"content" db:"content" gorm:"column:content" bson:"content"`
	// 是否开放评论
	IsOpened int `json:"isOpened" db:"is_opened" gorm:"column:is_opened" bson:"isOpened"`
	// 诉求描述
	HopeDesc string `json:"hopeDesc" db:"hope_desc" gorm:"column:hope_desc" bson:"hopeDesc"`
	// 图片
	FirstPhoto string `json:"firstPhoto" db:"first_photo" gorm:"column:first_photo" bson:"firstPhoto"`
	// 图片列表
	PhotoList string `json:"photoList" db:"photo_list" gorm:"column:photo_list" bson:"photoList"`
	// 状态,1:待处理 2:处理中 3:已完结
	Status int `json:"status" db:"status" gorm:"column:status" bson:"status"`
	// 分配的客服编号
	AllocateAid int `json:"allocateAid" db:"allocate_aid" gorm:"column:allocate_aid" bson:"allocateAid"`
	// 服务评分
	ServiceRank int `json:"serviceRank" db:"service_rank" gorm:"column:service_rank" bson:"serviceRank"`
	// 服务评价
	ServiceApprise string `json:"serviceApprise" db:"service_apprise" gorm:"column:service_apprise" bson:"serviceApprise"`
	// 是否有用 0:未评价 1:是 2:否
	IsUsefully int `json:"isUsefully" db:"is_usefully" gorm:"column:is_usefully" bson:"isUsefully"`
	// 创建时间
	CreateTime int `json:"createTime" db:"create_time" gorm:"column:create_time" bson:"createTime"`
	// 更新时间
	UpdateTime int `json:"updateTime" db:"update_time" gorm:"column:update_time" bson:"updateTime"`
}

func (w Workorder) TableName() string {
	return "workorder"
}

// WorkorderComment 工单详情
type WorkorderComment struct {
	// 编号
	Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" auto:"yes" bson:"id"`
	// 案件编号
	OrderId int `json:"orderId" db:"order_id" gorm:"column:order_id" bson:"orderId"`
	// 是否为回复信息,0:用户信息 1: 回复信息
	IsReplay int `json:"isReplay" db:"is_replay" gorm:"column:is_replay" bson:"isReplay"`
	// Content
	Content string `json:"content" db:"content" gorm:"column:content" bson:"content"`
	// 是否撤回 0:否 1:是
	IsRevert int `json:"isRevert" db:"is_revert" gorm:"column:is_revert" bson:"isRevert"`
	// 引用评论编号
	RefCid int `json:"refCid" db:"ref_cid" gorm:"column:ref_cid" bson:"refCid"`
	// 创建时间
	CreateTime int `json:"createTime" db:"create_time" gorm:"column:create_time" bson:"createTime"`
}

func (w WorkorderComment) TableName() string {
	return "workorder_comment"
}
