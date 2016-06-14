/**
 * Copyright 2015 @ z3q.net.
 * name : MssManager
 * author : jarryliu
 * date : 2015-07-26 23:08
 * description :
 * history :
 */
package mss

import (
    "go2o/core/domain/interface/mss"
    "time"
)

var _ mss.IMessageProvider = new(messageProviderImpl)

var _ mss.ISystemManager = new(messageManagerImpl)

type messageManagerImpl struct {
    _rep mss.IMssRep
}

func NewMessageManager(rep mss.IMssRep) mss.ISystemManager {
    return &messageManagerImpl{
        _rep:rep,
    }
}


// 获取所有的通知项
func (this *messageManagerImpl) GetAllNotifyItem()[]mss.NotifyItem{
    return this._rep.GetAllNotifyItem()
}

// 获取通知项配置
func (this *messageManagerImpl) GetNotifyItem(key string) mss.NotifyItem {
    return *this._rep.GetNotifyItem(key)
}
// 保存通知项设置
func (this *messageManagerImpl) SaveNotifyItem(item *mss.NotifyItem) error {
    v := this._rep.GetNotifyItem(item.Key)
    if v == nil {
        return mss.ErrNoSuchNotifyItem
    }
    v.Content = item.Content
    v.TplId = item.TplId
    v.NotifyBy = item.NotifyBy
    return this._rep.SaveNotifyItem(v)
}

type messageProviderImpl struct {
    _appUserId     int
    _mssRep        mss.IMssRep
    _mailTemplates []*mss.MailTemplate
    _config        *mss.Config
}

func NewMssManager(appUserId int, rep mss.IMssRep) mss.IMessageProvider {
    return &messageProviderImpl{
        _appUserId: appUserId,
        _mssRep:    rep,
    }
}

// 获取聚合根编号
func (this *messageProviderImpl) GetAggregateRootId() int {
    return this._appUserId
}

// 获取配置
func (this *messageProviderImpl) GetConfig() mss.Config {
    if this._config == nil {
        this._config = this._mssRep.GetConfig(this._appUserId)
    }
    return *this._config
}

// 保存消息设置
func (this *messageProviderImpl) SaveConfig(conf *mss.Config) error {
    err := this._mssRep.SaveConfig(this._appUserId, conf)
    if err == nil {
        this._config = nil
    }
    return err
}



// 创建消息模版对象
func (this *messageProviderImpl) CreateMessage(msg *mss.Message) (mss.IMessage) {
    return newMailTemplate(msg,this._mssRep)


    //todo: other message type
    //var err error
    //switch v.(type) {
    //case *mss.MailTemplate:
    //    tpl := v.(*mss.MailTemplate)
    //    if tpl.Enabled == 0 {
    //        err = mss.ErrNotEnabled
    //    }
    //    return newMailTemplate(msg,this._mssRep), err
    //}
    //return nil, mss.ErrNotSupportMessageType
}

// 发送消息
func (this *messageProviderImpl) Send(msg mss.IMessage, msgContent interface{},
    data mss.MessageData) error {
    if(msg.GetDomainId() ==0){
        msg.Save()
    }
    return msg.Send(msgContent,data)
}

// 获取邮箱模板
func (this *messageProviderImpl) GetMailTemplate(id int) *mss.MailTemplate {
    return this._mssRep.GetMailTemplate(this._appUserId, id)
}

// 保存邮箱模版
func (this *messageProviderImpl) SaveMailTemplate(v *mss.MailTemplate) (
int, error) {
    v.MerchantId = this._appUserId
    v.UpdateTime = time.Now().Unix()
    if v.CreateTime == 0 {
        v.CreateTime = v.UpdateTime
    }
    return this._mssRep.SaveMailTemplate(v)
}

// 删除邮件模板
func (this *messageProviderImpl) DeleteMailTemplate(id int) error {
    //merchantId := this._partner.GetAggregateRootId()
    //if this._partnerRep.CheckKvContainValue(merchantId, "kvset", strconv.Itoa(id), "mail") > 0 {
    //	return mss.ErrTemplateUsed
    //}
    return this._mssRep.DeleteMailTemplate(this._appUserId, id)
}

// 获取所有的邮箱模版
func (this *messageProviderImpl) GetMailTemplates() []*mss.MailTemplate {
    if this._mailTemplates == nil {
        this._mailTemplates = this._mssRep.GetMailTemplates(this._appUserId)
    }
    return this._mailTemplates
}
