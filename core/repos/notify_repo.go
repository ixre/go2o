/**
 * Copyright 2015 @ 56x.net.
 * name : notify_repo.go
 * author : jarryliu
 * date : 2016-07-06 18:45
 * description :
 * history :
 */
package repos

import (
	"database/sql"
	"log"

	mss "github.com/ixre/go2o/core/domain/interface/message"
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/infrastructure/fw"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/util"
)

var _ mss.INotifyRepo = new(notifyRepoImpl)

type notifyRepoImpl struct {
	_conn        db.Connector
	_orm         orm.Orm
	_itemGob     *util.GobFile
	_notifyItems map[string]*mss.NotifyItem
	registryRepo registry.IRegistryRepo
	manager      mss.INotifyManager
	templateRepo fw.Repository[mss.NotifyTemplate]
}

func NewNotifyRepo(o orm.Orm, on fw.ORM, registryRepo registry.IRegistryRepo) mss.INotifyRepo {
	s := &notifyRepoImpl{
		_conn:        o.Connector(),
		_orm:         o,
		registryRepo: registryRepo,
		_itemGob:     util.NewGobFile("conf/core/mss_notify"),
	}
	s.templateRepo = &fw.BaseRepository[mss.NotifyTemplate]{
		ORM: on,
	}
	return s
}

// TemplateRepo 获取通知模板仓储
func (n *notifyRepoImpl) TemplateRepo() fw.Repository[mss.NotifyTemplate] {
	return n.templateRepo
}

func (this *notifyRepoImpl) getNotifyItemMap() map[string]*mss.NotifyItem {
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
func (this *notifyRepoImpl) GetAllNotifyItem() []mss.NotifyItem {
	var list []mss.NotifyItem
	for _, v := range mss.DefaultNotifyItems {
		v2 := this.getNotifyItemMap()[v.Key]
		if v2 != nil {
			list = append(list, *v2)
		}
	}
	return list
}

// 获取通知项
func (this *notifyRepoImpl) GetNotifyItem(key string) *mss.NotifyItem {
	return this.getNotifyItemMap()[key]
}

// 保存通知项
func (this *notifyRepoImpl) SaveNotifyItem(v *mss.NotifyItem) error {
	this._notifyItems[v.Key] = v
	return this._itemGob.Save(this._notifyItems)
}

// SelectNotifyTemplate Select 系统通知模板
func (s *notifyRepoImpl) GetAllNotifyTemplate() []*mss.NotifyTemplate {
	list := make([]*mss.NotifyTemplate, 0)
	err := s._orm.Select(&list, "")
	if err != nil && err != sql.ErrNoRows {
		log.Printf("[ Orm][ Error]: %s; Entity:NotifyTemplate\n", err.Error())
	}
	return list
}

// SaveNotifyTemplate Save 系统通知模板
func (s *notifyRepoImpl) SaveNotifyTemplate(v *mss.NotifyTemplate) (*mss.NotifyTemplate, error) {
	id, err := orm.Save(s._orm, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Printf("[ Orm][ Error]: %s; Entity:NotifyTemplate\n", err.Error())
	}
	v.Id = id
	return v, err
}

// DeleteNotifyTemplate Delete 系统通知模板
func (s *notifyRepoImpl) DeleteNotifyTemplate(primary interface{}) error {
	err := s._orm.DeleteByPk(mss.NotifyTemplate{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("[ Orm][ Error]: %s; Entity:NotifyTemplate\n", err.Error())
	}
	return err
}
