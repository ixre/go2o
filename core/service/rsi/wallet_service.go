package rsi

import (
	"go2o/core/domain/interface/wallet"
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
