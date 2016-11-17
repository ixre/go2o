/**
 * Copyright 2015 @ z3q.net.
 * name : express_rep
 * author : jarryliu
 * date : 2016-07-05 18:33
 * description :
 * history :
 */
package repository

import (
	"github.com/jsix/gof/db"
	"github.com/jsix/gof/db/orm"
	expImpl "go2o/core/domain/express"
	"go2o/core/domain/interface/express"
	"go2o/core/domain/interface/valueobject"
	"sync"
)

type expressRep struct {
	db.Connector
	*expImpl.ExpressRepBase
	_valRep        valueobject.IValueRep
	ProvidersCache []*express.ExpressProvider
	mux            sync.Mutex
}

func NewExpressRep(conn db.Connector, valRep valueobject.IValueRep) express.IExpressRep {
	return &expressRep{
		Connector: conn,
		_valRep:   valRep,
	}
}

// 获取所有快递公司
func (this *expressRep) GetExpressProviders() []*express.ExpressProvider {
	mux.Lock()
	if this.ProvidersCache == nil {
		this.ProvidersCache = []*express.ExpressProvider{}
		err := this.GetOrm().Select(&this.ProvidersCache, "")
		if err != nil {
			panic(err)
		}
		if len(this.ProvidersCache) == 0 {
			this.ProvidersCache = this.SaveDefaultExpressProviders(this)
		}
	}
	mux.Unlock()
	return this.ProvidersCache
}

// 获取快递公司
func (this *expressRep) GetExpressProvider(id int64) *express.ExpressProvider {
	for _, v := range this.GetExpressProviders() {
		if v.Id == id {
			return v
		}
	}
	return nil
}

// 保存快递公司
func (this *expressRep) SaveExpressProvider(v *express.ExpressProvider) (int64, error) {
	this.ProvidersCache = nil
	return orm.Save(this.GetOrm(), v, v.Id)
}

// 获取用户的快递
func (this *expressRep) GetUserExpress(userId int64) express.IUserExpress {
	return expImpl.NewUserExpress(userId, this, this._valRep)
}

// 获取用户的快递模板
func (this *expressRep) GetUserAllTemplate(userId int64) []*express.ExpressTemplate {
	list := []*express.ExpressTemplate{}
	this.GetOrm().Select(&list, "user_id=?", userId)
	return list
}

// 删除快递模板
func (this *expressRep) DeleteExpressTemplate(userId int64, templateId int64) error {
	_, err := this.GetOrm().Delete(express.ExpressTemplate{},
		"id=? AND user_id=?", templateId, userId)
	return err
}

// 保存快递模板
func (this *expressRep) SaveExpressTemplate(v *express.ExpressTemplate) (int64, error) {
	return orm.Save(this.GetOrm(), v, v.Id)
}

// 获取模板的所有地区设置
func (this *expressRep) GetExpressTemplateAllAreaSet(templateId int64) []express.ExpressAreaTemplate {
	list := []express.ExpressAreaTemplate{}
	this.GetOrm().Select(&list, "template_id=?", templateId)
	return list
}

// 保存模板的地区设置
func (this *expressRep) SaveExpressTemplateAreaSet(v *express.ExpressAreaTemplate) (int64, error) {
	return orm.Save(this.GetOrm(), v, v.Id)
}

// 删除模板的地区设置
func (this *expressRep) DeleteAreaExpressTemplate(templateId int64, areaSetId int64) error {
	_, err := this.Connector.GetOrm().Delete(express.ExpressAreaTemplate{},
		"id= ? AND template_id=?", areaSetId, templateId)
	return err
}
