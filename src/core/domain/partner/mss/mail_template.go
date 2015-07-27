/**
 * Copyright 2015 @ S1N1 Team.
 * name : mail_template
 * author : jarryliu
 * date : 2015-07-27 09:19
 * description :
 * history :
 */
package mss

import (
	"go2o/src/core/domain/interface/partner/mss"
	mssIns "go2o/src/core/infrastructure/mss"
	"time"
)

var _ mss.IMsgTemplate = new(mailTemplate)

type mailTemplate struct {
	_rep       mss.IMssRep
	_partnerId int
	_tpl       *mss.MailTemplate
	_data      mss.MsgData
}

func newMailTemplate(partnerId int, rep mss.IMssRep, tpl *mss.MailTemplate) mss.IMsgTemplate {
	return &mailTemplate{
		_rep:       rep,
		_partnerId: partnerId,
		_tpl:       tpl,
	}
}

// 应用数据
func (this *mailTemplate) ApplyData(d mss.MsgData) {
	this._data = d
}

// 加入到发送对列
func (this *mailTemplate) JoinQueen(to []string) error {
	unix := time.Now().Unix()
	for _, t := range to {
		task := &mss.MailTask{
			PartnerId:  this._partnerId,
			Subject:    mssIns.Transplate(this._tpl.Subject, this._data),
			Body:       mssIns.Transplate(this._tpl.Body, this._data),
			To:         t,
			CreateTime: unix,
		}
		this._rep.JoinMailTaskToQueen(task)
	}
	return nil
}
