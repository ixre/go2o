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
	"github.com/jsix/gof/db"
	"go2o/core"
	"go2o/core/domain/interface/mss"
	mssImpl "go2o/core/domain/mss"
	"go2o/core/variable"
)

var _ mss.IMssRep = new(MssRep)

type MssRep struct {
	_conn    db.Connector
	_globMss mss.IMessageProvider
}

func NewMssRep(conn db.Connector) mss.IMssRep {
	return &MssRep{
		_conn: conn,
	}
}

func (this *MssRep) GetManager() mss.IMessageProvider {
	if this._globMss == nil {
		this._globMss = mssImpl.NewMssManager(0, this)
	}
	return this._globMss
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
