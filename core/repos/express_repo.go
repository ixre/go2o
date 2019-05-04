/**
 * Copyright 2015 @ z3q.net.
 * name : express_rep
 * author : jarryliu
 * date : 2016-07-05 18:33
 * description :
 * history :
 */
package repos

import (
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	expImpl "go2o/core/domain/express"
	"go2o/core/domain/interface/express"
	"go2o/core/domain/interface/valueobject"
	"sync"
)

type expressRepo struct {
	db.Connector
	*expImpl.ExpressRepBase
	_valRepo       valueobject.IValueRepo
	ProvidersCache []*express.ExpressProvider
	mux            sync.Mutex
}

func NewExpressRepo(conn db.Connector, valRepo valueobject.IValueRepo) express.IExpressRepo {
	return &expressRepo{
		Connector: conn,
		_valRepo:  valRepo,
	}
}

// 获取所有快递公司
func (er *expressRepo) GetExpressProviders() []*express.ExpressProvider {
	mux.Lock()
	if er.ProvidersCache == nil {
		er.ProvidersCache = []*express.ExpressProvider{}
		err := er.GetOrm().Select(&er.ProvidersCache, "")
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
func (er *expressRepo) GetExpressProvider(id int32) *express.ExpressProvider {
	for _, v := range er.GetExpressProviders() {
		if v.Id == id {
			return v
		}
	}
	return nil
}

// 保存快递公司
func (er *expressRepo) SaveExpressProvider(v *express.ExpressProvider) (int32, error) {
	er.ProvidersCache = nil
	return orm.I32(orm.Save(er.GetOrm(), v, int(v.Id)))
}

// 获取用户的快递
func (er *expressRepo) GetUserExpress(userId int32) express.IUserExpress {
	return expImpl.NewUserExpress(userId, er, er._valRepo)
}

// 获取用户的快递模板
func (er *expressRepo) GetUserAllTemplate(userId int32) []*express.ExpressTemplate {
	var list []*express.ExpressTemplate
	er.GetOrm().Select(&list, "user_id= $1", userId)
	return list
}

// 删除快递模板
func (er *expressRepo) DeleteExpressTemplate(userId int32, templateId int32) error {
	_, err := er.GetOrm().Delete(express.ExpressTemplate{},
		"id= $1 AND user_id= $2", templateId, userId)
	return err
}

// 保存快递模板
func (er *expressRepo) SaveExpressTemplate(v *express.ExpressTemplate) (int32, error) {
	return orm.I32(orm.Save(er.GetOrm(), v, int(v.Id)))
}

// 获取模板的所有地区设置
func (er *expressRepo) GetExpressTemplateAllAreaSet(templateId int32) []express.ExpressAreaTemplate {
	var list []express.ExpressAreaTemplate
	er.GetOrm().Select(&list, "template_id= $1", templateId)
	return list
}

// 保存模板的地区设置
func (er *expressRepo) SaveExpressTemplateAreaSet(v *express.ExpressAreaTemplate) (int32, error) {
	return orm.I32(orm.Save(er.GetOrm(), v, int(v.Id)))
}

// 删除模板的地区设置
func (er *expressRepo) DeleteAreaExpressTemplate(templateId int32, areaSetId int32) error {
	_, err := er.Connector.GetOrm().Delete(express.ExpressAreaTemplate{},
		"id= $1 AND template_id = $2", areaSetId, templateId)
	return err
}
