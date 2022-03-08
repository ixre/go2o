package events

import "github.com/ixre/go2o/core/domain/interface/wallet"

type WalletLogClickhouseUpdateEvent struct {
	Data *wallet.WalletLog
}
