package repos

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ixre/go2o/core/domain/interface/wallet"
	wi "github.com/ixre/go2o/core/domain/wallet"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
)

var _ wallet.IWalletRepo = new(WalletRepoImpl)

func NewWalletRepo(o orm.Orm) wallet.IWalletRepo {
	return &WalletRepoImpl{
		_orm:  o,
		_conn: o.Connector(),
	}
}

type WalletRepoImpl struct {
	_orm  orm.Orm
	_conn db.Connector
}

func (w *WalletRepoImpl) CreateWallet(userId int64,
	username string, walletType int,
	walletName string, flag int) wallet.IWallet {
	return w.createWallet1(&wallet.Wallet{
		UserId:     userId,
		UserName:   username,
		WalletType: walletType,
		WalletName: walletName,
		WalletFlag: flag,
	})
}

func (w *WalletRepoImpl) createWallet1(v *wallet.Wallet) wallet.IWallet {
	if v != nil {
		return wi.NewWallet(v, w)
	}
	return nil
}

func (w *WalletRepoImpl) GetWallet(walletId int64) wallet.IWallet {
	return w.createWallet1(w.getWallet_(walletId))
}

func (w *WalletRepoImpl) GetWalletByCode(code string) wallet.IWallet {
	l := w.GetWalletBy_("hash_code= $1", code)
	return w.createWallet1(l)
}

func (w *WalletRepoImpl) GetWalletByUserId(userId int64, walletType int) wallet.IWallet {
	l := w.GetWalletBy_("user_id= $1 AND wallet_type= $2 LIMIT 1", userId, walletType)
	return w.createWallet1(l)
}

func (w *WalletRepoImpl) CheckWalletUserMatch(userId int64, walletType int, walletId int64) bool {
	l := w.GetWalletBy_("user_id= $1 AND wallet_type= $2 AND id<> $3 LIMIT 1",
		userId, walletType, walletId)
	return l == nil
}

func (w *WalletRepoImpl) GetLog(walletId int64, logId int64) *wallet.WalletLog {
	l := w.GetWalletLog_(logId)
	if l != nil && l.WalletId == walletId {
		return l
	}
	return nil
}

func (w *WalletRepoImpl) PagingWalletLog(walletId int64, nodeId int, begin int, over int, where string, sort string) (total int, list []*wallet.WalletLog) {
	list = []*wallet.WalletLog{}
	table := "wal_wallet_log"
	err := w._conn.ExecScalar(fmt.Sprintf(`SELECT COUNT(1) FROM %s WHERE wallet_id= $1 %s`,
		table, where), &total, walletId)
	if total > 0 {
		if len(sort) > 0 {
			sort += ","
		}
		s := fmt.Sprintf(`SELECT * FROM %s WHERE wallet_id= $1 %s ORDER BY %s create_time DESC LIMIT $3 OFFSET $2`,
			table, where, sort)
		err = w._orm.SelectByQuery(&list, s, walletId, begin, over-begin)
	}
	if err != nil {
		log.Println("[ GO2O][ Repo][ Error]:", err.Error())
	}
	return total, list
}

// GetWalletLog_ Get WalletLog
func (w *WalletRepoImpl) GetWalletLog_(primary interface{}) *wallet.WalletLog {
	e := wallet.WalletLog{}
	err := w._orm.Get(primary, &e)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:WalletLog")
	}
	return nil
}

// GetWalletLogBy_ GetBy WalletLog
func (w *WalletRepoImpl) GetWalletLogBy_(where string, v ...interface{}) *wallet.WalletLog {
	e := wallet.WalletLog{}
	err := w._orm.GetBy(&e, where, v...)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:WalletLog")
	}
	return nil
}

// SelectWalletLog_ Select WalletLog
func (w *WalletRepoImpl) SelectWalletLog_(where string, v ...interface{}) []*wallet.WalletLog {
	var list []*wallet.WalletLog
	err := w._orm.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:WalletLog")
	}
	return list
}

// SaveWalletLog_ Save WalletLog
func (w *WalletRepoImpl) SaveWalletLog_(v *wallet.WalletLog) (int, error) {
	id, err := orm.Save(w._orm, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:WalletLog")
	}
	// if err == nil && v.Id > 0 {
	// 	eventbus.Publish(&events.AccountLogPushEvent{
	// 	})
	// }
	return id, err
}

// DeleteWalletLog_ Delete WalletLog
func (w *WalletRepoImpl) DeleteWalletLog_(primary interface{}) error {
	err := w._orm.DeleteByPk(wallet.WalletLog{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:WalletLog")
	}
	return err
}

// BatchDeleteWalletLog_ Batch Delete WalletLog
func (w *WalletRepoImpl) BatchDeleteWalletLog_(where string, v ...interface{}) (int64, error) {
	r, err := w._orm.Delete(wallet.WalletLog{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:WalletLog")
	}
	return r, err
}

// Get Wallet
func (w *WalletRepoImpl) getWallet_(primary interface{}) *wallet.Wallet {
	e := wallet.Wallet{}
	err := w._orm.Get(primary, &e)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:Wallet")
	}
	return nil
}

// GetWalletBy_ GetBy Wallet
func (w *WalletRepoImpl) GetWalletBy_(where string, v ...interface{}) *wallet.Wallet {
	e := wallet.Wallet{}
	err := w._orm.GetBy(&e, where, v...)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:Wallet")
	}
	return nil
}

// SelectWallet_ Select Wallet
func (w *WalletRepoImpl) SelectWallet_(where string, v ...interface{}) []*wallet.Wallet {
	var list []*wallet.Wallet
	err := w._orm.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:Wallet")
	}
	return list
}

// SaveWallet_ Save Wallet
func (w *WalletRepoImpl) SaveWallet_(v *wallet.Wallet) (int, error) {
	id, err := orm.Save(w._orm, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:Wallet")
	}
	return id, err
}

// DeleteWallet_ Delete Wallet
func (w *WalletRepoImpl) DeleteWallet_(primary interface{}) error {
	err := w._orm.DeleteByPk(wallet.Wallet{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:Wallet")
	}
	return err
}

// BatchDeleteWallet_ Batch Delete Wallet
func (w *WalletRepoImpl) BatchDeleteWallet_(where string, v ...interface{}) (int64, error) {
	r, err := w._orm.Delete(wallet.Wallet{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:Wallet")
	}
	return r, err
}
