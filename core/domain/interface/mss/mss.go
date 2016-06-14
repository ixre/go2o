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
)

//todo: waiting refactor

const (
    NotifyByMessage = 1 + iota
    NotifyByEMail
    NotifyByPhoneMessage
)

var (
    // 类型字典
    NotifyTypeMap = map[int]string{
        NotifyByMessage:      "站内信",
        NotifyByEMail:        "邮件",
        NotifyByPhoneMessage: "短信",
    }

    // 类型顺序
    NotifyTypeIndex = []int{
        NotifyByMessage,
        NotifyByEMail,
        NotifyByPhoneMessage,
    }

    // 默认通知项
    DefaultNotifyItems = NotifyItemSet{
        &NotifyItem{
            Key:"register_ok",
            TplId:-1,
            NotifyBy:NotifyByMessage,
            Content:"您好,恭喜您已注册成功{platform}的会员!",
            Tags:map[string]string{
                "platform":"平台名称",
            },
        },
        &NotifyItem{
            Key:"valid_phone",
            TplId:-1,
            NotifyBy:NotifyByMessage,
            Content:"您正在进行{operation},本次验证码为{code},有效期为{minutes}分种,[{platform}]。",
            Tags:map[string]string{
                "operation":"操作,如找回密码,重置手机等",
                "code":"验证码",
                "minutes":"有效时间",
                "platform":"平台名称",
            },
        },
    }
)

//可通过外部添加
func RegisterNotifyItem(key string, item *NotifyItem) {
    for _,v := range DefaultNotifyItems{
        if v.Key == key{
            panic(errors.New("通知项" + key + "已存在!"))
        }
    }
    DefaultNotifyItems = append(DefaultNotifyItems,item)
}

type (
    // 系统管理
    ISystemManager interface {
        // 获取所有的通知项
        GetAllNotifyItem()[]NotifyItem
        // 获取通知项配置
        GetNotifyItem(key string)NotifyItem
        // 保存通知项设置
        SaveNotifyItem(item *NotifyItem) error
    }

    // Message manager
    IMessageProvider interface {
        // 获取聚合根编号
        GetAggregateRootId() int

        // 获取配置
        GetConfig()Config

        // 保存消息设置
        SaveConfig(conf *Config) error

        // 发送消息
        Send(tpl IMsgTemplate, d MsgData, to []string) error

        // 获取邮箱模板
        GetMailTemplate(int) *MailTemplate

        // 保存邮箱模版
        SaveMailTemplate(*MailTemplate) (int, error)

        // 获取所有的邮箱模版
        GetMailTemplates() []*MailTemplate

        // 删除邮件模板
        DeleteMailTemplate(int) error

        // 创建消息模版对象
        CreateMsgTemplate(v interface{}) (IMsgTemplate, error)
    }

    IMssRep interface {
        // 获取消息提供者
        GetProvider() IMessageProvider

        // 系统消息服务
        GetManager() ISystemManager

        // 获取消息设置
        GetConfig(userId int) *Config

        // 保存消息设置
        SaveConfig(userId int, conf *Config) error

        // 获取所有的通知项
        GetAllNotifyItem()[]NotifyItem

        // 获取通知项
        GetNotifyItem(key string)*NotifyItem

        // 保存通知项
        SaveNotifyItem(v *NotifyItem)error

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
    }


    // 通知项
    NotifyItem struct {
        Key string
        NotifyBy int
        TplId    int
        Content  string
        Tags     map[string]string
    }

    // 通知项集合
    NotifyItemSet []*NotifyItem

    // 系统消息发送配置
    //todo: 过时的
    Config struct {
        //注册完成
        RegisterNotifyEnabled        bool
        // 通知类型
        RegisterNotifyType           int
        // 通知模板选择
        RegisterNotifyTpl            int
        // 注册通知的标签数据
        RegisterNotifyTagData        string

        // 资料完成
        ProfileCompleteNotifyEnabled bool
        // 通知类型
        ProfileCompleteNotifyType    int
        // 通知模板选择
        ProfileCompleteNotifyTpl     int
    }
)
