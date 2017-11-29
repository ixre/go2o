package repository

import (
	"database/sql"
	"github.com/jsix/gof/db"
	"github.com/jsix/gof/db/orm"
	"go2o/core/domain/interface/wallet"
	"log"
)

var _ wallet.IWalletRepo = new(WalletRepoImpl)

func NewWalletRepo(conn db.Connector) wallet.IWalletRepo {
	return &WalletRepoImpl{
		_orm: conn.GetOrm(),
	}
}

type WalletRepoImpl struct {
	_orm orm.Orm
}

func (w *WalletRepoImpl) CheckWalletUserMatch(userId int64, walletKind int, walletId int64) bool {
	l := w.GetWalletBy_("user_id=? AND wallet_kind=? AND id<>? LIMIT 1",
		userId, walletKind, walletId)
	return l == nil
}

func (w *WalletRepoImpl) GetLog(walletId int64, logId int64) *wallet.WalletLog {
	l := w.GetWalletLog_(logId)
	if l != nil && l.WalletId == walletId {
		return l
	}
	return nil
}

// Get WalletLog
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

// GetBy WalletLog
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

// Select WalletLog
func (w *WalletRepoImpl) SelectWalletLog_(where string, v ...interface{}) []*wallet.WalletLog {
	list := []*wallet.WalletLog{}
	err := w._orm.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:WalletLog")
	}
	return list
}

// Save WalletLog
func (w *WalletRepoImpl) SaveWalletLog_(v *wallet.WalletLog) (int, error) {
	id, err := orm.Save(w._orm, v, int(v.ID))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:WalletLog")
	}
	return id, err
}

// Delete WalletLog
func (w *WalletRepoImpl) DeleteWalletLog_(primary interface{}) error {
	err := w._orm.DeleteByPk(wallet.WalletLog{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:WalletLog")
	}
	return err
}

// Batch Delete WalletLog
func (w *WalletRepoImpl) BatchDeleteWalletLog_(where string, v ...interface{}) (int64, error) {
	r, err := w._orm.Delete(wallet.WalletLog{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:WalletLog")
	}
	return r, err
}

// Get Wallet
func (w *WalletRepoImpl) GetWallet_(primary interface{}) *wallet.Wallet {
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

// GetBy Wallet
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

// Select Wallet
func (w *WalletRepoImpl) SelectWallet_(where string, v ...interface{}) []*wallet.Wallet {
	list := []*wallet.Wallet{}
	err := w._orm.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:Wallet")
	}
	return list
}

// Save Wallet
func (w *WalletRepoImpl) SaveWallet_(v *wallet.Wallet) (int, error) {
	id, err := orm.Save(w._orm, v, int(v.ID))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:Wallet")
	}
	return id, err
}

// Delete Wallet
func (w *WalletRepoImpl) DeleteWallet_(primary interface{}) error {
	err := w._orm.DeleteByPk(wallet.Wallet{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:Wallet")
	}
	return err
}

// Batch Delete Wallet
func (w *WalletRepoImpl) BatchDeleteWallet_(where string, v ...interface{}) (int64, error) {
	r, err := w._orm.Delete(wallet.Wallet{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:Wallet")
	}
	return r, err
}
