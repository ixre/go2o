/**
 * Copyright 2015 @ to2.net.
 * name : mss_rep
 * author : jarryliu
 * date : 2015-07-27 09:03
 * description :
 * history :
 */
package repos

import (
	"database/sql"
	"fmt"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/util"
	"go2o/core"
	"go2o/core/domain/interface/mss"
	"go2o/core/domain/interface/mss/notify"
	"go2o/core/domain/interface/registry"
	"go2o/core/domain/interface/valueobject"
	mssImpl "go2o/core/domain/mss"
	notifyImpl "go2o/core/domain/mss/notify"
	"go2o/core/variable"
)

var _ mss.IMssRepo = new(mssRepo)

type mssRepo struct {
	_conn         db.Connector
	_sysManger    mss.IMessageManager
	_notifyManger notify.INotifyManager
	_notifyRepo   notify.INotifyRepo
	_valRepo      valueobject.IValueRepo
	registryRepo  registry.IRegistryRepo
	_globMss      mss.IUserMessageManager
}

func NewMssRepo(conn db.Connector, notifyRepo notify.INotifyRepo,
	registryRepo  registry.IRegistryRepo,
	valRepo valueobject.IValueRepo) mss.IMssRepo {
	return &mssRepo{
		_conn:       conn,
		_notifyRepo: notifyRepo,
		registryRepo:registryRepo,
		_valRepo:    valRepo,
	}
}

// 系统消息服务
func (m *mssRepo) MessageManager() mss.IMessageManager {
	if m._sysManger == nil {
		m._sysManger = mssImpl.NewMessageManager(m)
	}
	return m._sysManger
}

// 通知服务
func (m *mssRepo) NotifyManager() notify.INotifyManager {
	if m._notifyManger == nil {
		m._notifyManger = notifyImpl.NewNotifyManager(
			m._notifyRepo, m.registryRepo)
	}
	return m._notifyManger
}

func (m *mssRepo) GetProvider() mss.IUserMessageManager {
	if m._globMss == nil {
		m._globMss = mssImpl.NewMssManager(0, m)
	}
	return m._globMss
}

// 获取短信配置
func (m *mssRepo) GetConfig(userId int32) *mss.Config {
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
func (m *mssRepo) SaveConfig(userId int32, conf *mss.Config) error {
	filePath := "conf/core/mss_conf"
	if userId != 0 {
		filePath = fmt.Sprintf("conf/mch/%d/mss_conf", userId)
	}
	globFile := util.NewGobFile(filePath)
	return globFile.Save(conf)
}

// 获取邮箱模板
func (m *mssRepo) GetMailTemplate(mchId, id int32) *mss.MailTemplate {
	var e mss.MailTemplate
	if err := m._conn.GetOrm().Get(id, &e); err == nil {
		return &e
	}
	return nil
}

// 保存邮箱模版
func (m *mssRepo) SaveMailTemplate(v *mss.MailTemplate) (int32, error) {
	return orm.I32(orm.Save(m._conn.GetOrm(), v, int(v.Id)))
}

// 获取所有的邮箱模版
func (m *mssRepo) GetMailTemplates(mchId int32) []*mss.MailTemplate {
	var list = []*mss.MailTemplate{}
	m._conn.GetOrm().Select(&list, "merchant_id= $1", mchId)
	return list
}

// 删除邮件模板
func (m *mssRepo) DeleteMailTemplate(mchId, id int32) error {
	_, err := m._conn.GetOrm().Delete(mss.MailTemplate{},
		"merchant_id= $1 AND id= $2", mchId, id)
	return err
}

// 加入到发送对列
func (m *mssRepo) JoinMailTaskToQueen(v *mss.MailTask) error {
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
func (m *mssRepo) SaveMessage(v *mss.Message) (int32, error) {
	return orm.I32(orm.Save(m._conn.GetOrm(), v, int(v.Id)))
}

// 获取消息
func (m *mssRepo) GetMessage(id int32) *mss.Message {
	e := mss.Message{}
	if m._conn.GetOrm().Get(id, &e) == nil {
		e.To = []mss.User{}
		m._conn.Query(`SELECT to_id,to_role FROM msg_to WHERE msg_id= $1`, func(rs *sql.Rows) {
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
func (m *mssRepo) SaveUserMsg(v *mss.To) (int32, error) {
	return orm.I32(orm.Save(m._conn.GetOrm(), v, int(v.Id)))
}

// 保存消息内容
func (m *mssRepo) SaveMsgContent(v *mss.Content) (int32, error) {
	return orm.I32(orm.Save(m._conn.GetOrm(), v, int(v.Id)))
}

// 获取消息内容
func (m *mssRepo) GetMessageContent(msgId int32) *mss.Content {
	e := mss.Content{}
	if m._conn.GetOrm().GetBy(&e, "msg_id= $1", msgId) == nil {
		return &e
	}
	return nil
}

// 获取消息目标
func (m *mssRepo) GetMessageTo(msgId int32, toUserId int32, toRole int) *mss.To {
	e := mss.To{}
	if m._conn.GetOrm().GetByQuery(&e, `SELECT * FROM msg_to
		WHERE msg_id= $1 AND to_id = $2 AND to_role= $3`,
		msgId, toUserId, toRole) == nil {
		return &e
	}
	return nil
}
