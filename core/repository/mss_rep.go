/**
 * Copyright 2015 @ z3q.net.
 * name : mss_rep
 * author : jarryliu
 * date : 2015-07-27 09:03
 * description :
 * history :
 */
package repository

import (
	"fmt"
	"github.com/jsix/gof/db"
	"github.com/jsix/gof/util"
	"go2o/core"
	"go2o/core/domain/interface/mss"
	"go2o/core/domain/interface/mss/notify"
	"go2o/core/domain/interface/valueobject"
	mssImpl "go2o/core/domain/mss"
	notifyImpl "go2o/core/domain/mss/notify"
	"go2o/core/variable"
)

var _ mss.IMssRep = new(MssRep)

type MssRep struct {
	_conn         db.Connector
	_sysManger    mss.IMessageManager
	_notifyManger notify.INotifyManager
	_notifyRep    notify.INotifyRep
	_valRep       valueobject.IValueRep
	_globMss      mss.IUserMessageManager
}

func NewMssRep(conn db.Connector, notifyRep notify.INotifyRep,
	valRep valueobject.IValueRep) mss.IMssRep {
	return &MssRep{
		_conn:      conn,
		_notifyRep: notifyRep,
		_valRep:    valRep,
	}
}

// 系统消息服务
func (this *MssRep) MessageManager() mss.IMessageManager {
	if this._sysManger == nil {
		this._sysManger = mssImpl.NewMessageManager(this)
	}
	return this._sysManger
}

// 通知服务
func (this *MssRep) NotifyManager() notify.INotifyManager {
	if this._notifyManger == nil {
		this._notifyManger = notifyImpl.NewNotifyManager(
			this._notifyRep, this._valRep)
	}
	return this._notifyManger
}

func (this *MssRep) GetProvider() mss.IUserMessageManager {
	if this._globMss == nil {
		this._globMss = mssImpl.NewMssManager(0, this)
	}
	return this._globMss
}

// 获取短信配置
func (this *MssRep) GetConfig(userId int) *mss.Config {
	conf := mss.Config{}
	filePath := "conf/core/mss_conf"
	if userId != 0 {
		filePath = fmt.Sprintf("conf/mch/%d/mss_conf", userId)
	}
	globFile := util.NewGobFile(filePath)
	handleError(globFile.Unmarshal(&conf))
	return &conf
}

// 保存消息设置
func (this *MssRep) SaveConfig(userId int, conf *mss.Config) error {
	filePath := "conf/core/mss_conf"
	if userId != 0 {
		filePath = fmt.Sprintf("conf/mch/%d/mss_conf", userId)
	}
	globFile := util.NewGobFile(filePath)
	return globFile.Save(conf)
}

// 获取邮箱模板
func (this *MssRep) GetMailTemplate(merchantId, id int) *mss.MailTemplate {
	var e mss.MailTemplate
	if err := this._conn.GetOrm().Get(id, &e); err == nil {
		return &e
	}
	return nil
}

// 保存邮箱模版
func (this *MssRep) SaveMailTemplate(v *mss.MailTemplate) (int, error) {
	var err error
	var orm = this._conn.GetOrm()
	if v.Id > 0 {
		_, _, err = orm.Save(v.Id, v)
	} else {
		_, _, err = orm.Save(nil, v)
		this._conn.ExecScalar("SELECT MAX(id) FROM pt_mail_template WHERE merchant_id=?", &v.Id, v.MerchantId)
	}
	return v.Id, err
}

// 获取所有的邮箱模版
func (this *MssRep) GetMailTemplates(merchantId int) []*mss.MailTemplate {
	var list = []*mss.MailTemplate{}
	this._conn.GetOrm().Select(&list, "merchant_id=?", merchantId)
	return list
}

// 删除邮件模板
func (this *MssRep) DeleteMailTemplate(merchantId, id int) error {
	_, err := this._conn.GetOrm().Delete(mss.MailTemplate{}, "merchant_id=? AND id=?", merchantId, id)
	return err
}

// 加入到发送对列
func (this *MssRep) JoinMailTaskToQueen(v *mss.MailTask) error {
	var err error
	if v.Id > 0 {
		_, _, err = this._conn.GetOrm().Save(v.Id, v)
	} else {
		_, _, err = this._conn.GetOrm().Save(nil, v)
		if err == nil {
			err = this._conn.ExecScalar("SELECT max(id) FROM pt_mail_queue", &v.Id)
		}
	}

	if err == nil {
		rc := core.GetRedisConn()
		defer rc.Close()
		rc.Do("RPUSH", variable.KvNewMailTask, v.Id) // push to queue
	}
	return err
}

// 保存消息
func (this *MssRep) SaveMessage(v *mss.Message) (int, error) {
	var err error
	if v.Id > 0 {
		_, _, err = this._conn.GetOrm().Save(v.Id, v)
	} else {
		var id int64
		_, id, err = this._conn.GetOrm().Save(nil, v)
		v.Id = int(id)
	}
	return v.Id, err
}

// 获取消息
func (this *MssRep) GetMessage(id int) *mss.Message {
	e := mss.Message{}
	if this._conn.GetOrm().Get(id, &e) == nil {
		return &e
	}
	return nil
}

// 保存用户消息关联
func (this *MssRep) SaveUserMsg(v *mss.To) (int, error) {
	var err error
	if v.Id > 0 {
		_, _, err = this._conn.GetOrm().Save(v.Id, v)
	} else {
		var id int64
		_, id, err = this._conn.GetOrm().Save(nil, v)
		v.Id = int(id)
	}
	return v.Id, err
}

// 保存消息内容
func (this *MssRep) SaveMsgContent(v *mss.Content) (int, error) {
	var err error
	if v.Id > 0 {
		_, _, err = this._conn.GetOrm().Save(v.Id, v)
	} else {
		var id int64
		_, id, err = this._conn.GetOrm().Save(nil, v)
		v.Id = int(id)
	}
	return v.Id, err
}

// 获取消息内容
func (this *MssRep) GetMessageContent(msgId int) *mss.Content {
	e := mss.Content{}
	if this._conn.GetOrm().GetBy(&e, "msg_id=?", msgId) == nil {
		return &e
	}
	return nil
}

// 获取消息目标
func (this *MssRep) GetMessageTo(msgId, toUserId, toRole int) *mss.To {
	e := mss.To{}
	if this._conn.GetOrm().GetByQuery(&e, `SELECT * FROM msg_to t INNER JOIN msg_content c ON c.id = t.id
WHERE msg_id=? AND to_id =? AND to_role=?`, msgId, toUserId, toRole) == nil {
		return &e
	}
	return nil
}
