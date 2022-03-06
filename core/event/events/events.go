package events

import "github.com/ixre/go2o/core/domain/interface/wallet"

type WalletLogClickhouseWriteEvent struct {
	Data *wallet.WalletLog
}
