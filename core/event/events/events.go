package events

import (
	"github.com/ixre/go2o/core/domain/interface/mss"
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
	// 买家编号
	BuyerId int64
	// 订单金额
	OrderAmount int64
	// 分销商品
	AffiliateItems []*order.SubOrderItem
}

type SendSmsEvent struct {
	// 短信服务商
	Provider string
	// 发送短信场景
	Scene mss.MessageScene
	// 手机号
	Phone string
	// 数据
	Data []string
	// 短信内容
	Template string
	// 接口地址
	ApiConf *sms.SmsApi
}
