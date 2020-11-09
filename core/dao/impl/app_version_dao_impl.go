package impl

import (
	"database/sql"
	"github.com/ixre/gof/db/orm"
	"go2o/core/dao"
	"go2o/core/dao/model"
	"log"
)

var _ dao.IAppProdDao = new(appProdDaoImpl)
var appProdDaoImplMapped = false

type appProdDaoImpl struct {
	_orm orm.Orm
}


// Create new AppProdDao
func NewAppProdDao(o orm.Orm) dao.IAppProdDao {
	if !appProdDaoImplMapped {
		_ = o.Mapping(model.AppProd{}, "app_prod")
		_ = o.Mapping(model.AppVersion{},"app_version")
		appProdDaoImplMapped = true
	}
	return &appProdDaoImpl{
		_orm: o,
	}
}

// Get APP产品
func (t *appProdDaoImpl) Get(primary interface{}) *model.AppProd {
	e := model.AppProd{}
	err := t._orm.Get(primary, &e)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:AppProd")
	}
	return nil
}

// GetBy APP产品
func (t *appProdDaoImpl) GetBy(where string, v ...interface{}) *model.AppProd {
	e := model.AppProd{}
	err := t._orm.GetBy(&e, where, v...)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:AppProd")
	}
	return nil
}

// Count APP版本 by condition
func (t *appProdDaoImpl) Count(where string, v ...interface{}) (int, error) {
	return t._orm.Count(model.AppProd{}, where, v...)
}

// Select APP产品
func (t *appProdDaoImpl) Select(where string, v ...interface{}) []*model.AppProd {
	list := make([]*model.AppProd, 0)
	err := t._orm.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:AppProd")
	}
	return list
}

// Save APP产品
func (t *appProdDaoImpl) Save(v *model.AppProd) (int, error) {
	id, err := orm.Save(t._orm, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:AppProd")
	}
	return id, err
}

// Delete APP产品
func (t *appProdDaoImpl) Delete(primary interface{}) error {
	err := t._orm.DeleteByPk(model.AppProd{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:AppProd")
	}
	return err
}

// Batch Delete APP产品
func (t *appProdDaoImpl) BatchDelete(where string, v ...interface{}) (int64, error) {
	r, err := t._orm.Delete(model.AppProd{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:AppProd")
	}
	return r, err
}


// Get APP版本
func (t *appProdDaoImpl) GetVersion(primary interface{}) *model.AppVersion {
	e := model.AppVersion{}
	err := t._orm.Get(primary, &e)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:AppVersion")
	}
	return nil
}

// GetBy APP版本
func (t *appProdDaoImpl) GetVersionBy(where string, v ...interface{}) *model.AppVersion {
	e := model.AppVersion{}
	err := t._orm.GetBy(&e, where, v...)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:AppVersion")
	}
	return nil
}

// Select APP版本
func (t *appProdDaoImpl) SelectVersion(where string, v ...interface{}) []*model.AppVersion {
	list := make([]*model.AppVersion, 0)
	err := t._orm.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:AppVersion")
	}
	return list
}

// Save APP版本
func (t *appProdDaoImpl) SaveVersion(v *model.AppVersion) (int, error) {
	id, err := orm.Save(t._orm, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:AppVersion")
	}
	return id, err
}

// Delete APP版本
func (t *appProdDaoImpl) DeleteVersion(primary interface{}) error {
	err := t._orm.DeleteByPk(model.AppVersion{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:AppVersion")
	}
	return err
}

// Batch Delete APP版本
func (t *appProdDaoImpl) BatchDeleteVersion(where string, v ...interface{}) (int64, error) {
	r, err := t._orm.Delete(model.AppVersion{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:AppVersion")
	}
	return r, err
}
