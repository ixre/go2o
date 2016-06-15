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
	mssImpl "go2o/core/domain/mss"
	"go2o/core/infrastructure/domain"
	"go2o/core/variable"
)

var _ mss.IMssRep = new(MssRep)

type MssRep struct {
	_conn        db.Connector
	_globMss     mss.IUserMessageManager
	_notifyItems map[string]*mss.NotifyItem
	_itemGob     *util.GobFile
	_sysManger   mss.IMessageManager
}

func NewMssRep(conn db.Connector) mss.IMssRep {
	return &MssRep{
		_conn:    conn,
		_itemGob: util.NewGobFile("conf/core/mss_notify"),
	}
}

// 系统消息服务
func (this *MssRep) GetManager() mss.IMessageManager {
	if this._sysManger == nil {
		this._sysManger = mssImpl.NewMessageManager(this)
	}
	return this._sysManger
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
	domain.HandleError(globFile.Unmarshal(&conf))
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

func (this *MssRep) getNotifyItemMap() map[string]*mss.NotifyItem {
	if this._notifyItems == nil {
		this._notifyItems = map[string]*mss.NotifyItem{}
		err := this._itemGob.Unmarshal(&this._notifyItems)
		//拷贝系统默认的配置
		if err != nil {
			for _, v := range mss.DefaultNotifyItems {
				vv := *v
				this._notifyItems[v.Key] = &vv
			}
		}
	}
	return this._notifyItems
}

// 获取所有的通知项
func (this *MssRep) GetAllNotifyItem() []mss.NotifyItem {
	list := []mss.NotifyItem{}
	for _, v := range mss.DefaultNotifyItems {
		v2 := this.getNotifyItemMap()[v.Key]
		list = append(list, *v2)
	}
	return list
}

// 获取通知项
func (this *MssRep) GetNotifyItem(key string) *mss.NotifyItem {

	return this.getNotifyItemMap()[key]
}

// 保存通知项
func (this *MssRep) SaveNotifyItem(v *mss.NotifyItem) error {
	this._notifyItems[v.Key] = v
	return this._itemGob.Save(this._notifyItems)
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
	//todo:
	msg := mss.Message{}
	return &msg
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
