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
	"database/sql"
	"fmt"
	"github.com/jsix/gof/db"
	"github.com/jsix/gof/db/orm"
	"github.com/jsix/gof/util"
	"go2o/core"
	"go2o/core/domain/interface/mss"
	"go2o/core/domain/interface/mss/notify"
	"go2o/core/domain/interface/valueobject"
	mssImpl "go2o/core/domain/mss"
	notifyImpl "go2o/core/domain/mss/notify"
	"go2o/core/variable"
)

var _ mss.IMssRep = new(mssRep)

type mssRep struct {
	_conn         db.Connector
	_sysManger    mss.IMessageManager
	_notifyManger notify.INotifyManager
	_notifyRep    notify.INotifyRep
	_valRep       valueobject.IValueRep
	_globMss      mss.IUserMessageManager
}

func NewMssRep(conn db.Connector, notifyRep notify.INotifyRep,
	valRep valueobject.IValueRep) mss.IMssRep {
	return &mssRep{
		_conn:      conn,
		_notifyRep: notifyRep,
		_valRep:    valRep,
	}
}

// 系统消息服务
func (m *mssRep) MessageManager() mss.IMessageManager {
	if m._sysManger == nil {
		m._sysManger = mssImpl.NewMessageManager(m)
	}
	return m._sysManger
}

// 通知服务
func (m *mssRep) NotifyManager() notify.INotifyManager {
	if m._notifyManger == nil {
		m._notifyManger = notifyImpl.NewNotifyManager(
			m._notifyRep, m._valRep)
	}
	return m._notifyManger
}

func (m *mssRep) GetProvider() mss.IUserMessageManager {
	if m._globMss == nil {
		m._globMss = mssImpl.NewMssManager(0, m)
	}
	return m._globMss
}

// 获取短信配置
func (m *mssRep) GetConfig(userId int64) *mss.Config {
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
func (m *mssRep) SaveConfig(userId int, conf *mss.Config) error {
	filePath := "conf/core/mss_conf"
	if userId != 0 {
		filePath = fmt.Sprintf("conf/mch/%d/mss_conf", userId)
	}
	globFile := util.NewGobFile(filePath)
	return globFile.Save(conf)
}

// 获取邮箱模板
func (m *mssRep) GetMailTemplate(merchantId, id int) *mss.MailTemplate {
	var e mss.MailTemplate
	if err := m._conn.GetOrm().Get(id, &e); err == nil {
		return &e
	}
	return nil
}

// 保存邮箱模版
func (m *mssRep) SaveMailTemplate(v *mss.MailTemplate) (int, error) {
	var err error
	var orm = m._conn.GetOrm()
	if v.Id > 0 {
		_, _, err = orm.Save(v.Id, v)
	} else {
		_, _, err = orm.Save(nil, v)
		m._conn.ExecScalar("SELECT MAX(id) FROM pt_mail_template WHERE merchant_id=?", &v.Id, v.MerchantId)
	}
	return v.Id, err
}

// 获取所有的邮箱模版
func (m *mssRep) GetMailTemplates(merchantId int) []*mss.MailTemplate {
	var list = []*mss.MailTemplate{}
	m._conn.GetOrm().Select(&list, "merchant_id=?", merchantId)
	return list
}

// 删除邮件模板
func (m *mssRep) DeleteMailTemplate(merchantId, id int) error {
	_, err := m._conn.GetOrm().Delete(mss.MailTemplate{}, "merchant_id=? AND id=?", merchantId, id)
	return err
}

// 加入到发送对列
func (m *mssRep) JoinMailTaskToQueen(v *mss.MailTask) error {
	var err error
	if v.Id > 0 {
		_, _, err = m._conn.GetOrm().Save(v.Id, v)
	} else {
		_, _, err = m._conn.GetOrm().Save(nil, v)
		if err == nil {
			err = m._conn.ExecScalar("SELECT max(id) FROM pt_mail_queue", &v.Id)
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
func (m *mssRep) SaveMessage(v *mss.Message) (int64, error) {
	var err error
	if v.Id > 0 {
		_, _, err = m._conn.GetOrm().Save(v.Id, v)
	} else {
		var id int64
		_, id, err = m._conn.GetOrm().Save(nil, v)
		v.Id = int(id)
	}
	return v.Id, err
}

// 获取消息
func (m *mssRep) GetMessage(id int64) *mss.Message {
	e := mss.Message{}
	if m._conn.GetOrm().Get(id, &e) == nil {
		e.To = []mss.User{}
		m._conn.Query(`SELECT to_id,to_role FROM msg_to WHERE msg_id=?`, func(rs *sql.Rows) {
			for rs.Next() {
				u := mss.User{}
				rs.Scan(&u.Id, &u.Role)
				e.To = append(e.To, u)
			}
		}, id)
		return &e
	}
	return nil
}

// 保存用户消息关联
func (m *mssRep) SaveUserMsg(v *mss.To) (int64, error) {
	return orm.Save(m._conn.GetOrm(), v, v.Id)
}

// 保存消息内容
func (m *mssRep) SaveMsgContent(v *mss.Content) (int64, error) {
	var err error
	if v.Id > 0 {
		_, _, err = m._conn.GetOrm().Save(v.Id, v)
	} else {
		var id int64
		_, id, err = m._conn.GetOrm().Save(nil, v)
		v.Id = int(id)
	}
	return v.Id, err
}

// 获取消息内容
func (m *mssRep) GetMessageContent(msgId int64) *mss.Content {
	e := mss.Content{}
	if m._conn.GetOrm().GetBy(&e, "msg_id=?", msgId) == nil {
		return &e
	}
	return nil
}

// 获取消息目标
func (m *mssRep) GetMessageTo(msgId int64, toUserId int64, toRole int) *mss.To {
	e := mss.To{}
	if m._conn.GetOrm().GetByQuery(&e, `SELECT * FROM msg_to
		WHERE msg_id=? AND to_id =? AND to_role=?`, msgId, toUserId, toRole) == nil {
		return &e
	}
	return nil
}
