package rsi

import (
	"go2o/core/domain/interface/wallet"
	"go2o/core/service/thrift/parser"
	"go2o/gen-code/thrift/define"
)

var _ define.WalletService = new(walletServiceImpl)

func NewWalletService(repo wallet.IWalletRepo) define.WalletService {
	return &walletServiceImpl{
		_repo: repo,
	}
}

type walletServiceImpl struct {
	_repo wallet.IWalletRepo
}

func (w *walletServiceImpl) CreateWallet(userId int64, walletType int32, flag int32, remark string) (r *define.Result_, err error) {
	v := &wallet.Wallet{
		UserId:     userId,
		WalletType: int(walletType),
		WalletFlag: int(flag),
		Remark:     remark,
	}
	iw := w._repo.CreateWallet(v)
	return parser.Result(iw.Save()), nil
}

func (w *walletServiceImpl) GetWalletId(userId int64, walletType int32) (r int32, err error) {
	panic("implement me")
}

func (w *walletServiceImpl) Adjust(walletId int64, value int32, title string, outerNo string, opuId int32, opuName string) (r *define.Result_, err error) {
	panic("implement me")
}

func (w *walletServiceImpl) Discount(walletId int64, value int32, title string, outerNo string, opuId int32, opuName string, must bool) (r *define.Result_, err error) {
	panic("implement me")
}

func (w *walletServiceImpl) Freeze(walletId int64, value int32, title string, outerNo string, opuId int32, opuName string) (r *define.Result_, err error) {
	panic("implement me")
}

func (w *walletServiceImpl) Unfreeze(walletId int64, value int32, title string, outerNo string, opuId int32, opuName string) (r *define.Result_, err error) {
	panic("implement me")
}

func (w *walletServiceImpl) Charge(walletId int64, value int32, by int32, title string, outerNo string, opuId int32, opuName string) (r *define.Result_, err error) {
	panic("implement me")
}

func (w *walletServiceImpl) Transfer(walletId int64, toWalletId int64, value int32, tradeFee int32, remark string) (r *define.Result_, err error) {
	panic("implement me")
}

func (w *walletServiceImpl) RequestTakeOut(walletId int64, value int32, tradeFee int32, kind int32, title string) (r *define.Result_, err error) {
	panic("implement me")
}

func (w *walletServiceImpl) ReviewTakeOut(walletId int64, takeId int64, reviewPass bool, remark string) (r *define.Result_, err error) {
	panic("implement me")
}

func (w *walletServiceImpl) FinishTakeOut(walletId int64, takeId int64, outerNo string) (r *define.Result_, err error) {
	panic("implement me")
}
