package events

import (
	"github.com/ixre/go2o/core/domain/interface/order"
	"github.com/ixre/go2o/core/domain/interface/wallet"
)

// 钱包更新事件
type WalletLogClickhouseUpdateEvent struct {
	Data *wallet.WalletLog
}

// 订单分销事件
type OrderAffiliteRebateEvent struct {
	// 订单号
	OrderNo string
	// 订单金额
	OrderAmount int64
	// 分销商品
	AffiliteItems []*order.SubOrderItem
}
