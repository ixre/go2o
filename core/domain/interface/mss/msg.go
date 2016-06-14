/**
 * Copyright 2015 @ z3q.net.
 * name : msg_template
 * author : jarryliu
 * date : 2015-07-26 21:57
 * description :
 * history :
 */
package mss

var (
    RoleSystem = 0
    RoleMember = 1
    RoleMerchant = 2
)

type(
    IMessage interface {
        // 应用数据
        //Parse(MessageData) string
        // 加入到发送对列
        //JoinQueen(to []string) error

        // 获取领域编号
        GetDomainId()int

        // 保存
        Save()(int,error)

        // 发送
        Send(msg interface{}, data MessageData) error
    }

    // 消息值
    ValueMessage string

    // 简讯
    ValuePhoneMessage  string

    // 邮件消息
    ValueMailMessage struct {
        // 主题
        Subject string `json:"subject"`
        // 内容
        Body    string `json:"body"`
    }

    // 站内信
    ValueSiteMessage struct {
        // 主题
        Subject string  `json:"subject"`
        // 信息内容
        Message string `json:"message"`
    }

    // 消息,优先级为: AllUser ->  ToRole  ->  To
    Message struct {
        // 消息编号
        Id         int  `db:"id" pk:"yes" auto:"yes"`
        // 发送人角色
        SenderRole   int `db:"sender_role"`
        // 发送人类型
        SenderId int `db:"sender_int"`
        // 发送的目标
        To         []User `db:"-"`
        // 发送的用户角色
        ToRole     int  `db:"to_role"`
        // 全系统接收
        AllUser    int   `db:"all_user"`
        // 消息类型
        Type      int `db:"msg_type"`
        // 是否只能阅读
        Readonly   int  `db:"read_only"`
    }

    User  struct {
        Id   int
        Role int
    }
)


// 消息数据
type MessageData map[string]string


// 邮件模版
type MailTemplate struct {
    // 编号
    Id         int `db:"id" pk:"yes" auto:"yes"`
    // 商户编号
    MerchantId int `db:"merchant_id"`
    // 名称
    Name       string `db:"name"`
    // 主题
    Subject    string `db:"subject"`
    // 内容
    Body       string `db:"body"`

    // 是否启用
    Enabled    int `db:"enabled"`

    // 创建时间
    CreateTime int64 `db:"create_time"`
    // 更新时间
    UpdateTime int64 `db:"update_time"`
}

type MailTask struct {
    // 编号
    Id         int `db:"id" pk:"yes" auto:"yes"`
    // 任务编号,无任务为0
    TaskId     int `db:"task_id"`
    // 商户编号
    MerchantId int `db:"merchant_id"`
    // 发送至
    SendTo     string `db:"send_to"`
    // 主题
    Subject    string `db:"subject"`
    // 内容
    Body       string `db:"body"`
    // 是否发送(0,1)
    IsSend     int `db:"is_send"`
    // 是否失败(0,1)
    IsFailed   int `db:"is_failed"`
    // 创建时间
    CreateTime int64 `db:"create_time"`
    // 发送时间
    SendTime   int64 `db:"update_time`
}