package rsi

import (
	"context"
	"go2o/core/domain/interface/wallet"
	"go2o/core/service/auto_gen/rpc/ttype"
	"go2o/core/service/auto_gen/rpc/wallet_service"
	"go2o/core/service/thrift/parser"
)

var _ wallet_service.WalletService = new(walletServiceImpl)

func NewWalletService(repo wallet.IWalletRepo) wallet_service.WalletService {
	return &walletServiceImpl{
		_repo: repo,
	}
}

type walletServiceImpl struct {
	_repo wallet.IWalletRepo
	serviceUtil
}

func (w *walletServiceImpl) CreateWallet(ctx context.Context, userId int64, walletType int32, flag int32, remark string) (*ttype.Result_, error) {
	v := &wallet.Wallet{
		UserId:     userId,
		WalletType: int(walletType),
		WalletFlag: int(flag),
		Remark:     remark,
	}
	iw := w._repo.CreateWallet(v)
	_, err := iw.Save()
	return w.result(err), nil
}

func (w *walletServiceImpl) GetWalletId(ctx context.Context, userId int64, walletType int32) (r int64, err error) {
	iw := w._repo.GetWalletByUserId(userId, int(walletType))
	if iw != nil {
		return iw.GetAggregateRootId(), nil
	}
	return 0, nil
}

func (w *walletServiceImpl) GetWallet(ctx context.Context, walletId int64) (r *wallet_service.SWallet, err error) {
	iw := w._repo.GetWallet(walletId)
	if iw != nil {
		return w.parseWallet(iw.Get()), nil
	}
	return nil, nil
}

func (w *walletServiceImpl) GetWalletLog(ctx context.Context, walletId int64, id int64) (r *wallet_service.SWalletLog, err error) {
	iw := w._repo.GetWallet(walletId)
	if iw != nil {
		if l := iw.GetLog(id); l.ID > 0 {
			return w.parseWalletLog(l), nil
		}
	}
	return nil, nil
}
func (w *walletServiceImpl) Adjust(ctx context.Context, walletId int64, value int32, title string, outerNo string, opuId int32, opuName string) (r *ttype.Result_, err error) {
	iw := w._repo.GetWallet(walletId)
	if iw == nil {
		err = wallet.ErrNoSuchWalletAccount
	} else {
		err = iw.Adjust(int(value), title, outerNo, int(opuId), opuName)
	}
	return w.result(err), nil
}

func (w *walletServiceImpl) Discount(ctx context.Context, walletId int64, value int32, title string, outerNo string, must bool) (r *ttype.Result_, err error) {
	iw := w._repo.GetWallet(walletId)
	if iw == nil {
		err = wallet.ErrNoSuchWalletAccount
	} else {
		err = iw.Discount(int(value), title, outerNo, must)
	}
	return w.result(err), nil
}

func (w *walletServiceImpl) Freeze(ctx context.Context, walletId int64, value int32, title string, outerNo string, opuId int32, opuName string) (r *ttype.Result_, err error) {
	iw := w._repo.GetWallet(walletId)
	if iw == nil {
		err = wallet.ErrNoSuchWalletAccount
	} else {
		err = iw.Freeze(int(value), title, outerNo, int(opuId), opuName)
	}
	return w.result(err), nil
}

func (w *walletServiceImpl) Unfreeze(ctx context.Context, walletId int64, value int32, title string, outerNo string, opuId int32, opuName string) (r *ttype.Result_, err error) {
	iw := w._repo.GetWallet(walletId)
	if iw == nil {
		err = wallet.ErrNoSuchWalletAccount
	} else {
		err = iw.Unfreeze(int(value), title, outerNo, int(opuId), opuName)
	}
	return w.result(err), nil
}

func (w *walletServiceImpl) Charge(ctx context.Context, walletId int64, value int32, by int32, title string, outerNo string, opuId int32, opuName string) (r *ttype.Result_, err error) {
	iw := w._repo.GetWallet(walletId)
	if iw == nil {
		err = wallet.ErrNoSuchWalletAccount
	} else {
		err = iw.Charge(int(value), int(by), title, outerNo, int(opuId), opuName)
	}
	return w.result(err), nil
}

func (w *walletServiceImpl) Transfer(ctx context.Context, walletId int64, toWalletId int64, value int32, tradeFee int32, remark string) (r *ttype.Result_, err error) {
	iw := w._repo.GetWallet(walletId)
	if iw == nil {
		err = wallet.ErrNoSuchWalletAccount
	} else {
		title := "钱包转账"
		toTitle := "钱包收款"
		//todo: title
		err = iw.Transfer(toWalletId, int(value), int(tradeFee), title, toTitle, remark)
	}
	return w.result(err), nil
}

func (w *walletServiceImpl) RequestTakeOut(ctx context.Context, walletId int64, value int32, tradeFee int32, kind int32, title string) (r *ttype.Result_, err error) {
	iw := w._repo.GetWallet(walletId)
	if iw == nil {
		err = wallet.ErrNoSuchWalletAccount
	} else {
		_, tradeNo, err1 := iw.RequestTakeOut(int(value), int(tradeFee), int(kind), title)
		if err1 != nil {
			err = err1
		} else {
			return w.success(map[string]string{
				"TradeNo": tradeNo,
			}), nil
		}
	}
	return w.result(err), nil
}

func (w *walletServiceImpl) ReviewTakeOut(ctx context.Context, walletId int64, takeId int64, reviewPass bool, remark string, opuId int32, opuName string) (r *ttype.Result_, err error) {
	iw := w._repo.GetWallet(walletId)
	if iw == nil {
		err = wallet.ErrNoSuchWalletAccount
	} else {
		err = iw.ReviewTakeOut(takeId, reviewPass, remark, int(opuId), opuName)
	}
	return w.result(err), nil
}

func (w *walletServiceImpl) FinishTakeOut(ctx context.Context, walletId int64, takeId int64, outerNo string) (r *ttype.Result_, err error) {
	iw := w._repo.GetWallet(walletId)
	if iw == nil {
		err = wallet.ErrNoSuchWalletAccount
	} else {
		err = iw.FinishTakeOut(takeId, outerNo)
	}
	return w.result(err), nil
}

func (w *walletServiceImpl) PagingWalletLog(ctx context.Context, walletId int64, params *ttype.SPagingParams) (r *ttype.SPagingResult_, err error) {
	iw := w._repo.GetWallet(walletId)
	if iw == nil {
		return parser.PagingResult(0, nil, wallet.ErrNoSuchWalletAccount), nil
	}
	sortBy := params.OrderField
	if params.OrderDesc {
		sortBy += " DESC"
	}
	total, list := iw.PagingLog(int(params.Begin), int(params.Over), params.Opt, sortBy)
	return parser.PagingResult(total, list, err), nil
}

func (w *walletServiceImpl) parseWallet(v wallet.Wallet) *wallet_service.SWallet {
	return &wallet_service.SWallet{
		ID:             v.ID,
		HashCode:       v.HashCode,
		NodeId:         int32(v.NodeId),
		UserId:         v.UserId,
		WalletType:     int32(v.WalletType),
		WalletFlag:     int32(v.WalletFlag),
		Balance:        int32(v.Balance),
		PresentBalance: int32(v.PresentBalance),
		AdjustAmount:   int32(v.AdjustAmount),
		FreezeAmount:   int32(v.FreezeAmount),
		LatestAmount:   int32(v.LatestAmount),
		ExpiredAmount:  int32(v.ExpiredAmount),
		TotalCharge:    int32(v.TotalCharge),
		TotalPresent:   int32(v.TotalPresent),
		TotalPay:       int32(v.TotalPay),
		State:          int32(v.State),
		Remark:         v.Remark,
		CreateTime:     v.CreateTime,
		UpdateTime:     v.UpdateTime,
	}
}
func (w *walletServiceImpl) parseWalletLog(l wallet.WalletLog) *wallet_service.SWalletLog {
	return &wallet_service.SWalletLog{
		ID:           l.ID,
		WalletId:     l.WalletId,
		Kind:         int32(l.Kind),
		Title:        l.Title,
		OuterChan:    l.OuterChan,
		OuterNo:      l.OuterNo,
		Value:        int32(l.Value),
		Balance:      int32(l.Balance),
		TradeFee:     int32(l.TradeFee),
		OperatorId:   int32(l.OperatorId),
		OperatorName: l.OperatorName,
		Remark:       l.Remark,
		ReviewState:  int32(l.ReviewState),
		ReviewRemark: l.ReviewRemark,
		ReviewTime:   l.ReviewTime,
		CreateTime:   l.CreateTime,
		UpdateTime:   l.UpdateTime,
	}
}
