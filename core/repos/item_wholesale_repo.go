package repos

import (
	"database/sql"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/util"
	"go2o/core/domain/interface/item"
	"log"
)

var _ item.IItemWholesaleRepo = new(itemWholesaleRepo)

type itemWholesaleRepo struct {
	_orm orm.Orm
}

func NewItemWholesaleRepo(conn db.Connector) item.IItemWholesaleRepo {
	return &itemWholesaleRepo{
		_orm: conn.GetOrm(),
	}
}

// Get WsItem
func (w *itemWholesaleRepo) GetWsItem(primary interface{}) *item.WsItem {
	e := item.WsItem{}
	err := w._orm.Get(primary, &e)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:WsItem")
	}
	return nil
}

// Select WsItem
func (w *itemWholesaleRepo) SelectWsItem(where string, v ...interface{}) []*item.WsItem {
	list := []*item.WsItem{}
	err := w._orm.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:WsItem")
	}
	return list
}

// Save WsItem
func (w *itemWholesaleRepo) SaveWsItem(v *item.WsItem, create bool) (int, error) {
	iid := util.BoolExt.TInt(create, 0, int(v.ItemId))
	id, err := orm.Save(w._orm, v, iid)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:WsItem")
	}
	return id, err
}

// Delete WsItem
func (w *itemWholesaleRepo) DeleteWsItem(primary interface{}) error {
	err := w._orm.DeleteByPk(item.WsItem{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:WsItem")
	}
	return err
}

// Batch Delete WsItem
func (w *itemWholesaleRepo) BatchDeleteWsItem(where string, v ...interface{}) (int64, error) {
	r, err := w._orm.Delete(item.WsItem{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:WsItem")
	}
	return r, err
}

// Get WsItemDiscount
func (w *itemWholesaleRepo) GetWsItemDiscount(primary interface{}) *item.WsItemDiscount {
	e := item.WsItemDiscount{}
	err := w._orm.Get(primary, &e)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:WsItemDiscount")
	}
	return nil
}

// Select WsItemDiscount
func (w *itemWholesaleRepo) SelectWsItemDiscount(where string, v ...interface{}) []*item.WsItemDiscount {
	list := []*item.WsItemDiscount{}
	err := w._orm.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:WsItemDiscount")
	}
	return list
}

// Save WsItemDiscount
func (w *itemWholesaleRepo) SaveWsItemDiscount(v *item.WsItemDiscount) (int, error) {
	id, err := orm.Save(w._orm, v, int(v.ID))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:WsItemDiscount")
	}
	return id, err
}

// Delete WsItemDiscount
func (w *itemWholesaleRepo) DeleteWsItemDiscount(primary interface{}) error {
	err := w._orm.DeleteByPk(item.WsItemDiscount{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:WsItemDiscount")
	}
	return err
}

// Batch Delete WsItemDiscount
func (w *itemWholesaleRepo) BatchDeleteWsItemDiscount(where string, v ...interface{}) (int64, error) {
	r, err := w._orm.Delete(item.WsItemDiscount{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:WsItemDiscount")
	}
	return r, err
}

// Get WsSkuPrice
func (w *itemWholesaleRepo) GetWsSkuPrice(primary interface{}) *item.WsSkuPrice {
	e := item.WsSkuPrice{}
	err := w._orm.Get(primary, &e)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:WsSkuPrice")
	}
	return nil
}

// Select WsSkuPrice
func (w *itemWholesaleRepo) SelectWsSkuPrice(where string, v ...interface{}) []*item.WsSkuPrice {
	list := []*item.WsSkuPrice{}
	err := w._orm.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:WsSkuPrice")
	}
	return list
}

// Save WsSkuPrice
func (w *itemWholesaleRepo) SaveWsSkuPrice(v *item.WsSkuPrice) (int, error) {
	id, err := orm.Save(w._orm, v, int(v.ID))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:WsSkuPrice")
	}
	return id, err
}

// Delete WsSkuPrice
func (w *itemWholesaleRepo) DeleteWsSkuPrice(primary interface{}) error {
	err := w._orm.DeleteByPk(item.WsSkuPrice{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:WsSkuPrice")
	}
	return err
}

// Batch Delete WsSkuPrice
func (w *itemWholesaleRepo) BatchDeleteWsSkuPrice(where string, v ...interface{}) (int64, error) {
	r, err := w._orm.Delete(item.WsSkuPrice{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:WsSkuPrice")
	}
	return r, err
}
