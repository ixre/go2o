/**
 * Copyright 2015 @ S1N1 Team.
 * name : mss_rep
 * author : jarryliu
 * date : 2015-07-27 09:03
 * description :
 * history :
 */
package repository

import (
	"github.com/atnet/gof/db"
	"go2o/src/core/domain/interface/partner/mss"
)

var _ mss.IMssRep = new(MssRep)

type MssRep struct {
	_conn db.Connector
}

func NewMssRep(conn db.Connector) mss.IMssRep {
	return &MssRep{
		_conn: conn,
	}
}

// 获取邮箱模板
func (this *MssRep) GetMailTemplate(partnerId, id int) *mss.MailTemplate {
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
		this._conn.ExecScalar("SELECT MAX(id) FROM pt_mail_template WHERE partner_id=?", &v.Id, v.PartnerId)
	}
	return v.Id, err
}

// 获取所有的邮箱模版
func (this *MssRep) GetMailTemplates(partnerId int) []*mss.MailTemplate {
	var list = []*mss.MailTemplate{}
	this._conn.GetOrm().Select(&list, " partner_id=?", partnerId)
	return list
}

// 加入到发送对列
func (this *MssRep)  JoinMailTaskToQueen(v *mss.MailTask)error{
	_,_,err := this._conn.GetOrm().Save(nil,v)
	return err
}