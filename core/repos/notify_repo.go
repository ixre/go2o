/**
 * Copyright 2015 @ to2.net.
 * name : notify_repo.go
 * author : jarryliu
 * date : 2016-07-06 18:45
 * description :
 * history :
 */
package repos

import (
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/util"
	"go2o/core/domain/interface/mss/notify"
	"go2o/core/domain/interface/registry"
	notify2 "go2o/core/domain/mss/notify"
)

var _ notify.INotifyRepo = new(notifyRepImpl)

type notifyRepImpl struct {
	_conn        db.Connector
	_itemGob     *util.GobFile
	_notifyItems map[string]*notify.NotifyItem
	registryRepo registry.IRegistryRepo
	manager notify.INotifyManager
}

func (n *notifyRepImpl) Manager() notify.INotifyManager {
	if(n.manager == nil){
		n.manager = notify2.NewNotifyManager(n,n.registryRepo)
	}
	panic("implement me")
}

func NewNotifyRepo(conn db.Connector,registryRepo registry.IRegistryRepo) notify.INotifyRepo {
	return &notifyRepImpl{
		_conn:    conn,
		registryRepo:registryRepo,
		_itemGob: util.NewGobFile("conf/core/mss_notify"),
	}
}

func (this *notifyRepImpl) getNotifyItemMap() map[string]*notify.NotifyItem {
	if this._notifyItems == nil {
		this._notifyItems = map[string]*notify.NotifyItem{}
		err := this._itemGob.Unmarshal(&this._notifyItems)
		//拷贝系统默认的配置
		if err != nil {
			for _, v := range notify.DefaultNotifyItems {
				vv := *v
				this._notifyItems[v.Key] = &vv
			}
		}
	}
	return this._notifyItems
}

// 获取所有的通知项
func (this *notifyRepImpl) GetAllNotifyItem() []notify.NotifyItem {
	list := []notify.NotifyItem{}
	for _, v := range notify.DefaultNotifyItems {
		v2 := this.getNotifyItemMap()[v.Key]
		if v2 != nil {
			list = append(list, *v2)
		}
	}
	return list
}

// 获取通知项
func (this *notifyRepImpl) GetNotifyItem(key string) *notify.NotifyItem {
	return this.getNotifyItemMap()[key]
}

// 保存通知项
func (this *notifyRepImpl) SaveNotifyItem(v *notify.NotifyItem) error {
	this._notifyItems[v.Key] = v
	return this._itemGob.Save(this._notifyItems)
}
