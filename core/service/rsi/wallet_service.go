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
	iw := w._repo.GetWalletByUserId(userId, int(walletType))
	if iw != nil {
		return int32(iw.GetAggregateRootId()), nil
	}
	return 0, nil
}

func (w *walletServiceImpl) Adjust(walletId int64, value int32, title string, outerNo string, opuId int32, opuName string) (r *define.Result_, err error) {
	iw := w._repo.GetWallet(walletId)
	if iw == nil {
		err = wallet.ErrNoSuchWalletAccount
	} else {
		err = iw.Adjust(int(value), title, outerNo, int(opuId), opuName)
	}
	return parser.Result(nil, err), nil
}

func (w *walletServiceImpl) Discount(walletId int64, value int32, title string, outerNo string, opuId int32, opuName string, must bool) (r *define.Result_, err error) {
	iw := w._repo.GetWallet(walletId)
	if iw == nil {
		err = wallet.ErrNoSuchWalletAccount
	} else {
		err = iw.Discount(int(value), title, outerNo, int(opuId), opuName, must)
	}
	return parser.Result(nil, err), nil
}

func (w *walletServiceImpl) Freeze(walletId int64, value int32, title string, outerNo string, opuId int32, opuName string) (r *define.Result_, err error) {
	iw := w._repo.GetWallet(walletId)
	if iw == nil {
		err = wallet.ErrNoSuchWalletAccount
	} else {
		err = iw.Freeze(int(value), title, outerNo, int(opuId), opuName)
	}
	return parser.Result(nil, err), nil
}

func (w *walletServiceImpl) Unfreeze(walletId int64, value int32, title string, outerNo string, opuId int32, opuName string) (r *define.Result_, err error) {
	iw := w._repo.GetWallet(walletId)
	if iw == nil {
		err = wallet.ErrNoSuchWalletAccount
	} else {
		err = iw.Unfreeze(int(value), title, outerNo, int(opuId), opuName)
	}
	return parser.Result(nil, err), nil
}

func (w *walletServiceImpl) Charge(walletId int64, value int32, by int32, title string, outerNo string, opuId int32, opuName string) (r *define.Result_, err error) {
	iw := w._repo.GetWallet(walletId)
	if iw == nil {
		err = wallet.ErrNoSuchWalletAccount
	} else {
		err = iw.Charge(int(value), int(by), title, outerNo, int(opuId), opuName)
	}
	return parser.Result(nil, err), nil
}

func (w *walletServiceImpl) Transfer(walletId int64, toWalletId int64, value int32, tradeFee int32, remark string) (r *define.Result_, err error) {
	iw := w._repo.GetWallet(walletId)
	if iw == nil {
		err = wallet.ErrNoSuchWalletAccount
	} else {
		title := "钱包转账"
		toTitle := "钱包收款"
		//todo: title
		err = iw.Transfer(toWalletId, int(value), int(tradeFee), title, toTitle, remark)
	}
	return parser.Result(nil, err), nil
}

func (w *walletServiceImpl) RequestTakeOut(walletId int64, value int32, tradeFee int32, kind int32, title string) (r *define.Result_, err error) {
	iw := w._repo.GetWallet(walletId)
	if iw == nil {
		err = wallet.ErrNoSuchWalletAccount
	} else {
		_, tradeNo, err1 := iw.RequestTakeOut(int(value), int(tradeFee), int(kind), title)
		if err1 != nil {
			err = err1
		} else {
			return parser.Result(tradeNo, nil), nil
		}
	}
	return parser.Result(nil, err), nil
}

func (w *walletServiceImpl) ReviewTakeOut(walletId int64, takeId int64, reviewPass bool, remark string, opuId int32, opuName string) (r *define.Result_, err error) {
	iw := w._repo.GetWallet(walletId)
	if iw == nil {
		err = wallet.ErrNoSuchWalletAccount
	} else {
		err = iw.ReviewTakeOut(takeId, reviewPass, remark, int(opuId), opuName)
	}
	return parser.Result(nil, err), nil
}

func (w *walletServiceImpl) FinishTakeOut(walletId int64, takeId int64, outerNo string) (r *define.Result_, err error) {
	iw := w._repo.GetWallet(walletId)
	if iw == nil {
		err = wallet.ErrNoSuchWalletAccount
	} else {
		err = iw.FinishTakeOut(takeId, outerNo)
	}
	return parser.Result(nil, err), nil
}

func (w *walletServiceImpl) PagingWalletLog(walletId int64, params *define.PagingParams) (r *define.PagingResult_, err error) {
	iw := w._repo.GetWallet(walletId)
	if iw == nil {
		return parser.PagingResult(0, nil, wallet.ErrNoSuchWalletAccount), nil
	}
	sortBy := params.OrderField
	if params.OrderDesc {
		sortBy += " DESC"
	}
	total, list := iw.PagingLog(int(params.Begin), int(params.Over),params.Opt, sortBy)
	return parser.PagingResult(total, list, err), nil
}
