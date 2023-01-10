package events

import (
	"github.com/ixre/go2o/core/domain/interface/order"
	"github.com/ixre/go2o/core/domain/interface/wallet"
	"github.com/ixre/go2o/core/infrastructure/tool/sms"
)

// 钱包更新事件
type WalletLogClickhouseUpdateEvent struct {
	Data *wallet.WalletLog
}

// 订单分销事件
type OrderAffiliateRebateEvent struct {
	// 订单号
	OrderNo string
	// 订单金额
	OrderAmount int64
	// 分销商品
	AffiliateItems []*order.SubOrderItem
}

type SendSmsEvent struct {
	Provider string
	Phone    string
	Content  string
	Sene  string 
	Data     map[string]interface{}
	ApiConf  *sms.SmsApi
}
