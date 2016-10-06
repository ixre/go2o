/**
 * Copyright 2015 @ z3q.net.
 * name : msg_manager
 * author : jarryliu
 * date : 2015-07-26 21:30
 * description :
 * history :
 */
package mss

import (
	"go2o/core/domain/interface/mss/notify"
	"go2o/core/infrastructure/domain"
)

//todo: waiting refactor

var (
	ErrNotSupportMessageType *domain.DomainError = domain.NewDomainError(
		"err_not_support_message_type", "不支持的消息类型")

	ErrNotEnabled *domain.DomainError = domain.NewDomainError(
		"err_template_not_enabled", "模板未启用")

	ErrTemplateUsed *domain.DomainError = domain.NewDomainError(
		"err_template_used", "模板被使用，无法删除")

	ErrMessageUpdate *domain.DomainError = domain.NewDomainError(
		"err_message_update", "消息不需要更新")

	ErrMessageNotSave *domain.DomainError = domain.NewDomainError(
		"err_message_not_save", "请在消息发送前保存")

	ErrUnknownMessageUseFor *domain.DomainError = domain.NewDomainError(
		"err_unknown_message_use_for", "未知的消息用途")

	ErrMessageAllUser *domain.DomainError = domain.NewDomainError(
		"err_message_all_user", "消息为全员消息,指定了多余的参数",
	)

	ErrMessageToRole *domain.DomainError = domain.NewDomainError(
		"err_message_to_role", "消息为用户类型消息,指定了多余的用户",
	)

	ErrUnknownRole *domain.DomainError = domain.NewDomainError(
		"err_unknown_role", "未知的用户类型")

	ErrNoSuchReceiveUser *domain.DomainError = domain.NewDomainError(
		"err_no_such_receive_user", "消息没有指定接收用户")
)

type (
	// 系统管理
	IMessageManager interface {
		// 创建消息对象
		CreateMessage(msg *Message, content interface{}) IMessage

		// 创建用于会员通知的消息对象
		CreateMemberNotifyMessage(memberId int, msgType int,
			content interface{}) IMessage

		// 获取消息
		GetMessage(id int) IMessage

		// 获取聊天会话编号
		GetChatSessionId(senderRole, senderId, toRole, toId int) int

		// 创建聊天会话
		CreateChatSession(senderRole, senderId, toRole, toId int) (Message, error)
	}

	// Message manager,主要用于管理用户的模板
	IUserMessageManager interface {
		// 获取聚合根编号
		GetAggregateRootId() int

		// 获取配置
		GetConfig() Config

		// 保存消息设置
		SaveConfig(conf *Config) error

		// 获取邮箱模板
		GetMailTemplate(int) *MailTemplate

		// 保存邮箱模版
		SaveMailTemplate(*MailTemplate) (int, error)

		// 获取所有的邮箱模版
		GetMailTemplates() []*MailTemplate

		// 删除邮件模板
		DeleteMailTemplate(int) error
	}

	IMssRep interface {
		// 获取消息提供者
		GetProvider() IUserMessageManager

		// 系统消息服务
		MessageManager() IMessageManager

		// 通知服务
		NotifyManager() notify.INotifyManager

		// 获取消息设置
		GetConfig(userId int) *Config

		// 保存消息设置
		SaveConfig(userId int, conf *Config) error

		// 获取邮箱模板
		GetMailTemplate(userId, id int) *MailTemplate

		// 保存邮箱模版
		SaveMailTemplate(*MailTemplate) (int, error)

		// 获取所有的邮箱模版
		GetMailTemplates(userId int) []*MailTemplate

		// 删除邮件模板
		DeleteMailTemplate(userId, id int) error

		// 加入到发送对列
		JoinMailTaskToQueen(*MailTask) error

		// 保存消息
		SaveMessage(msg *Message) (int, error)

		// 获取消息
		GetMessage(id int) *Message

		// 保存用户消息关联
		SaveUserMsg(to *To) (int, error)

		// 保存消息内容
		SaveMsgContent(co *Content) (int, error)

		// 获取消息内容
		GetMessageContent(msgId int) *Content

		// 获取消息目标
		GetMessageTo(msgId, toUserId, toRole int) *To
	}

	// 系统消息发送配置
	//todo: 过时的
	Config struct {
		//注册完成
		RegisterNotifyEnabled bool
		// 通知类型
		RegisterNotifyType int
		// 通知模板选择
		RegisterNotifyTpl int
		// 注册通知的标签数据
		RegisterNotifyTagData string

		// 资料完成
		ProfileCompleteNotifyEnabled bool
		// 通知类型
		ProfileCompleteNotifyType int
		// 通知模板选择
		ProfileCompleteNotifyTpl int
	}
)
