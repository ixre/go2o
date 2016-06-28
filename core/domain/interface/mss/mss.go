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
	"errors"
	"go2o/core/infrastructure/domain"
)

//todo: waiting refactor

const (
	TypeSiteMessage = 1 + iota
	TypeEmailMessage
	TypePhoneMessage
)

var (
	ErrNotSupportMessageType *domain.DomainError = domain.NewDomainError(
		"err_not_support_message_type", "不支持的消息类型")

	ErrNotEnabled *domain.DomainError = domain.NewDomainError(
		"err_template_not_enabled", "模板未启用")

	ErrTemplateUsed *domain.DomainError = domain.NewDomainError(
		"err_template_used", "模板被使用，无法删除")

	ErrNoSuchNotifyItem *domain.DomainError = domain.NewDomainError(
		"err_no_such_notify_item", "通知项不存在")

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

	// 类型字典
	NotifyTypeMap = map[int]string{
		TypeSiteMessage:  "站内信",
		TypeEmailMessage: "邮件",
		TypePhoneMessage: "短信",
	}

	// 类型顺序
	NotifyTypeIndex = []int{
		TypeSiteMessage,
		TypeEmailMessage,
		TypePhoneMessage,
	}

	// 默认通知项
	DefaultNotifyItems = NotifyItemSet{
		&NotifyItem{
			Key:      "注册通知",
			TplId:    -1,
			NotifyBy: TypeSiteMessage,
			Content:  "您好,恭喜您已注册成功{platform}的会员!",
			Tags: map[string]string{
				"platform": "平台名称",
			},
		},
		&NotifyItem{
			Key:        "验证手机",
			TplId:      -1,
			ReadonlyBy: true,
			NotifyBy:   TypePhoneMessage,
			Content:    "您正在进行{operation},本次验证码为{code},有效期为{minutes}分种,[{platform}]。",
			Tags: map[string]string{
				"operation": "操作,如找回密码,重置手机等",
				"code":      "验证码",
				"minutes":   "有效时间",
				"platform":  "平台名称",
			},
		},
	}
)

//可通过外部添加
func RegisterNotifyItem(key string, item *NotifyItem) {
	for _, v := range DefaultNotifyItems {
		if v.Key == key {
			panic(errors.New("通知项" + key + "已存在!"))
		}
	}
	DefaultNotifyItems = append(DefaultNotifyItems, item)
}

type (
	// 系统管理
	IMessageManager interface {
		// 获取所有的通知项
		GetAllNotifyItem() []NotifyItem
		// 获取通知项配置
		GetNotifyItem(key string) NotifyItem
		// 保存通知项设置
		SaveNotifyItem(item *NotifyItem) error

		// 创建消息对象
		CreateMessage(msg *Message, content interface{}) IMessage

		// 获取消息
		GetMessage(id int) IMessage
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
		GetManager() IMessageManager

		// 获取消息设置
		GetConfig(userId int) *Config

		// 保存消息设置
		SaveConfig(userId int, conf *Config) error

		// 获取所有的通知项
		GetAllNotifyItem() []NotifyItem

		// 获取通知项
		GetNotifyItem(key string) *NotifyItem

		// 保存通知项
		SaveNotifyItem(v *NotifyItem) error

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

	// 通知项
	NotifyItem struct {
		Key string
		// 发送方式
		NotifyBy int
		// 不允许修改发送方式
		ReadonlyBy bool
		TplId      int
		Content    string
		Tags       map[string]string
	}

	// 通知项集合
	NotifyItemSet []*NotifyItem

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
