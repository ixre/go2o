/**
 * Copyright 2015 @ 56x.net.
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
	"log"

	mss "github.com/ixre/go2o/core/domain/interface/message"
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/domain/interface/valueobject"
	mssImpl "github.com/ixre/go2o/core/domain/message"
	notifyImpl "github.com/ixre/go2o/core/domain/message/notify"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/util"
)

var _ mss.IMessageRepo = new(messageRepoImpl)

type messageRepoImpl struct {
	_conn         db.Connector
	_sysManger    mss.IMessageManager
	_notifyManger mss.INotifyManager
	_notifyRepo   mss.INotifyRepo
	_valRepo      valueobject.IValueRepo
	registryRepo  registry.IRegistryRepo
	_globMss      mss.IUserMessageManager
	_orm          orm.Orm
}

var messageRepoMapped = false

func NewMssRepo(o orm.Orm, notifyRepo mss.INotifyRepo,
	registryRepo registry.IRegistryRepo,
	valRepo valueobject.IValueRepo) mss.IMessageRepo {
	if !messageRepoMapped {
		_ = o.Mapping(mss.NotifyTemplate{}, "sys_notify_template")
		messageRepoMapped = true
	}
	return &messageRepoImpl{
		_conn:        o.Connector(),
		_orm:         o,
		_notifyRepo:  notifyRepo,
		registryRepo: registryRepo,
		_valRepo:     valRepo,
	}
}

// 系统消息服务
func (m *messageRepoImpl) MessageManager() mss.IMessageManager {
	if m._sysManger == nil {
		m._sysManger = mssImpl.NewMessageManager(m)
	}
	return m._sysManger
}

// 通知服务
func (m *messageRepoImpl) NotifyManager() mss.INotifyManager {
	if m._notifyManger == nil {
		m._notifyManger = notifyImpl.NewNotifyManager(
			m._notifyRepo, m, m.registryRepo)
	}
	return m._notifyManger
}

func (m *messageRepoImpl) GetProvider() mss.IUserMessageManager {
	if m._globMss == nil {
		m._globMss = mssImpl.NewMssManager(0, m)
	}
	return m._globMss
}

// SelectNotifyTemplate Select 系统通知模板
func (s *messageRepoImpl) GetAllNotifyTemplate() []*mss.NotifyTemplate {
	list := make([]*mss.NotifyTemplate, 0)
	err := s._orm.Select(&list, "")
	if err != nil && err != sql.ErrNoRows {
		log.Printf("[ Orm][ Error]: %s; Entity:NotifyTemplate\n", err.Error())
	}
	return list
}

// SaveNotifyTemplate Save 系统通知模板
func (s *messageRepoImpl) SaveNotifyTemplate(v *mss.NotifyTemplate) (int, error) {
	id, err := orm.Save(s._orm, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Printf("[ Orm][ Error]: %s; Entity:NotifyTemplate\n", err.Error())
	}
	return id, err
}

// DeleteNotifyTemplate Delete 系统通知模板
func (s *messageRepoImpl) DeleteNotifyTemplate(primary interface{}) error {
	err := s._orm.DeleteByPk(mss.NotifyTemplate{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("[ Orm][ Error]: %s; Entity:NotifyTemplate\n", err.Error())
	}
	return err
}

// 获取短信配置
func (m *messageRepoImpl) GetConfig(userId int64) *mss.Config {
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
func (m *messageRepoImpl) SaveConfig(userId int64, conf *mss.Config) error {
	filePath := "conf/core/mss_conf"
	if userId != 0 {
		filePath = fmt.Sprintf("conf/mch/%d/mss_conf", userId)
	}
	globFile := util.NewGobFile(filePath)
	return globFile.Save(conf)
}

// 获取邮箱模板
func (m *messageRepoImpl) GetMailTemplate(mchId int64, id int32) *mss.MailTemplate {
	var e mss.MailTemplate
	if err := m._orm.Get(id, &e); err == nil {
		return &e
	}
	return nil
}

// 保存邮箱模版
func (m *messageRepoImpl) SaveMailTemplate(v *mss.MailTemplate) (int32, error) {
	return orm.I32(orm.Save(m._orm, v, int(v.Id)))
}

// 获取所有的邮箱模版
func (m *messageRepoImpl) GetMailTemplates(mchId int64) []*mss.MailTemplate {
	var list = []*mss.MailTemplate{}
	m._orm.Select(&list, "merchant_id= $1", mchId)
	return list
}

// 删除邮件模板
func (m *messageRepoImpl) DeleteMailTemplate(mchId, id int64) error {
	_, err := m._orm.Delete(mss.MailTemplate{},
		"merchant_id= $1 AND id= $2", mchId, id)
	return err
}

// 加入到发送对列
func (m *messageRepoImpl) JoinMailTaskToQueen(v *mss.MailTask) error {
	var err error
	if v.Id > 0 {
		_, _, err = m._orm.Save(v.Id, v)
	} else {
		_, _, err = m._orm.Save(nil, v)
		if err == nil {
			err = m._conn.ExecScalar("SELECT max(id) FROM pt_mail_queue", &v.Id)
		}
	}

	if err == nil {
		//rc := core.GetRedisConn()
		//defer rc.Close()
		//rc.Do("RPUSH", variable.KvNewMailTask, v.Id) // push to queue
	}
	return err
}

// 保存消息
func (m *messageRepoImpl) SaveMessage(v *mss.Message) (int32, error) {
	return orm.I32(orm.Save(m._orm, v, int(v.Id)))
}

// 获取消息
func (m *messageRepoImpl) GetMessage(id int32) *mss.Message {
	e := mss.Message{}
	if m._orm.Get(id, &e) == nil {
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
func (m *messageRepoImpl) SaveUserMsg(v *mss.To) (int32, error) {
	return orm.I32(orm.Save(m._orm, v, int(v.Id)))
}

// 保存消息内容
func (m *messageRepoImpl) SaveMsgContent(v *mss.Content) (int32, error) {
	return orm.I32(orm.Save(m._orm, v, int(v.Id)))
}

// 获取消息内容
func (m *messageRepoImpl) GetMessageContent(msgId int32) *mss.Content {
	e := mss.Content{}
	if m._orm.GetBy(&e, "msg_id= $1", msgId) == nil {
		return &e
	}
	return nil
}

// 获取消息目标
func (m *messageRepoImpl) GetMessageTo(msgId int32, toUserId int32, toRole int) *mss.To {
	e := mss.To{}
	if m._orm.GetByQuery(&e, `SELECT * FROM msg_to
		WHERE msg_id= $1 AND to_id = $2 AND to_role= $3`,
		msgId, toUserId, toRole) == nil {
		return &e
	}
	return nil
}
