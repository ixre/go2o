package repos

import (
	"database/sql"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	"go2o/core/domain/interface/merchant/wholesaler"
	"log"
)

var _ wholesaler.IWholesaleRepo = new(wholesaleRepo)

type wholesaleRepo struct {
	_orm  orm.Orm
	_conn db.Connector
}

// Create new WsWholesalerRepo
func NewWholesaleRepo(conn db.Connector) *wholesaleRepo {
	return &wholesaleRepo{
		_orm:  conn.GetOrm(),
		_conn: conn,
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

// 同步商品
func (w *wholesaleRepo) SyncItems(vendorId int32, shelve, review int32) (add int, del int) {

	del, _, err2 := w._conn.Exec(`DELETE FROM ws_item WHERE
    vendor_id= $1 AND item_id NOT IN (SELECT id FROM item_info
    WHERE vendor_id= $2)`, vendorId, vendorId)
	if err2 != nil {
		log.Println("wholesale item sync fail:", err2.Error())
	}
	return add, del
}

// 获取待同步商品
func (w *wholesaleRepo) GetAwaitSyncItems(vendorId int32) []int {
	add := []int{}
	i := 0
	err := w._conn.Query(`SELECT id FROM item_info WHERE
		vendor_id = $1 AND id NOT IN (SELECT item_id FROM
		 ws_item WHERE vendor_id= $2)`, func(rows *sql.Rows) {
		for rows.Next() {
			rows.Scan(&i)
			add = append(add, i)
		}
	}, vendorId, vendorId)
	if err != nil && err != sql.ErrNoRows {
		log.Println("wholesale get awit item fail:", err.Error())
	}
	return add
}

// Select WsRebateRate
func (w *wholesaleRepo) SelectWsRebateRate(where string, v ...interface{}) []*wholesaler.WsRebateRate {
	var list []*wholesaler.WsRebateRate
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
