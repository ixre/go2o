package repository

import (
	"database/sql"
	"github.com/jsix/gof/db"
	"github.com/jsix/gof/db/orm"
	"go2o/core/domain/interface/merchant/wholesaler"
	"log"
)

var _ wholesaler.IWholesaleRepo = new(wholesaleRepo)

type wholesaleRepo struct {
	_orm orm.Orm
}

// Create new WsWholesalerRepo
func NewWholesaleRepo(conn db.Connector) *wholesaleRepo {
	return &wholesaleRepo{
		_orm: conn.GetOrm(),
	}
}

// Get WsWholesaler
func (w *wholesaleRepo) GetWsWholesaler(primary interface{}) *wholesaler.WsWholesaler {
	e := wholesaler.WsWholesaler{}
	err := w._orm.Get(primary, &e)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:WsWholesaler")
	}
	return nil
}

// Save WsWholesaler
func (w *wholesaleRepo) SaveWsWholesaler(v *wholesaler.WsWholesaler, create bool) (int, error) {
	iid := int(v.MchId)
	if create {
		iid = 0
	}
	id, err := orm.Save(w._orm, v, iid)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:WsWholesaler")
	}
	return id, err
}

// Select WsRebateRate
func (w *wholesaleRepo) SelectWsRebateRate(where string, v ...interface{}) []*wholesaler.WsRebateRate {
	list := []*wholesaler.WsRebateRate{}
	err := w._orm.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:WsRebateRate")
	}
	return list
}

// Save WsRebateRate
func (w *wholesaleRepo) SaveWsRebateRate(v *wholesaler.WsRebateRate) (int, error) {
	id, err := orm.Save(w._orm, v, int(v.ID))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:WsRebateRate")
	}
	return id, err
}

// Batch Delete WsRebateRate
func (w *wholesaleRepo) BatchDeleteWsRebateRate(where string, v ...interface{}) (int64, error) {
	r, err := w._orm.Delete(wholesaler.WsRebateRate{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:WsRebateRate")
	}
	return r, err
}
