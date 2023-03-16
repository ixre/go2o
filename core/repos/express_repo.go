/**
 * Copyright 2015 @ 56x.net.
 * name : express_rep
 * author : jarryliu
 * date : 2016-07-05 18:33
 * description :
 * history :
 */
package repos

import (
	"sync"

	expImpl "github.com/ixre/go2o/core/domain/express"
	"github.com/ixre/go2o/core/domain/interface/express"
	"github.com/ixre/go2o/core/domain/interface/valueobject"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
)

type expressRepo struct {
	db.Connector
	*expImpl.ExpressRepBase
	_valRepo       valueobject.IValueRepo
	ProvidersCache []*express.Provider
	mux            sync.Mutex
	o              orm.Orm
}

func NewExpressRepo(o orm.Orm, valRepo valueobject.IValueRepo) express.IExpressRepo {
	return &expressRepo{
		Connector: o.Connector(),
		o:         o,
		_valRepo:  valRepo,
	}
}

// 获取所有快递公司
func (er *expressRepo) GetExpressProviders() []*express.Provider {
	mux.Lock()
	if er.ProvidersCache == nil {
		er.ProvidersCache = []*express.Provider{}
		err := er.o.Select(&er.ProvidersCache, "")
		if err != nil {
			panic(err)
		}
		if len(er.ProvidersCache) == 0 {
			er.ProvidersCache = er.SaveDefaultExpressProviders(er)
		}
	}
	mux.Unlock()
	return er.ProvidersCache
}

// 获取快递公司
func (er *expressRepo) GetExpressProvider(id int32) *express.Provider {
	for _, v := range er.GetExpressProviders() {
		if v.Id == id {
			return v
		}
	}
	return nil
}

// 保存快递公司
func (er *expressRepo) SaveExpressProvider(v *express.Provider) (int32, error) {
	er.ProvidersCache = nil
	return orm.I32(orm.Save(er.o, v, int(v.Id)))
}

// 获取用户的快递
func (er *expressRepo) GetUserExpress(userId int) express.IUserExpress {
	return expImpl.NewUserExpress(userId, er, er._valRepo)
}

// 获取用户的快递模板
func (er *expressRepo) GetUserAllTemplate(userId int) []*express.ExpressTemplate {
	var list []*express.ExpressTemplate
	er.o.Select(&list, "vendor_id= $1", userId)
	return list
}

// 删除快递模板
func (er *expressRepo) DeleteExpressTemplate(userId int, templateId int) error {
	_, err := er.o.Delete(express.ExpressTemplate{},
		"id= $1 AND vendor_id= $2", templateId, userId)
	return err
}

// 保存快递模板
func (er *expressRepo) SaveExpressTemplate(v *express.ExpressTemplate) (int, error) {
	return orm.Save(er.o, v, int(v.Id))
}

// 获取模板的所有地区设置
func (er *expressRepo) GetExpressTemplateAllAreaSet(templateId int) []express.RegionExpressTemplate {
	var list []express.RegionExpressTemplate
	er.o.Select(&list, "template_id= $1", templateId)
	return list
}

// 保存模板的地区设置
func (er *expressRepo) SaveExpressTemplateAreaSet(v *express.RegionExpressTemplate) (int, error) {
	return orm.Save(er.o, v, int(v.Id))
}

// 删除模板的地区设置
func (er *expressRepo) DeleteAreaExpressTemplate(templateId int, areaSetId int) error {
	_, err := er.o.Delete(express.RegionExpressTemplate{},
		"id= $1 AND template_id = $2", areaSetId, templateId)
	return err
}
