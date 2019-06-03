/**
 * Copyright 2015 @ to2.net.
 * name : exchange
 * author : jarryliu
 * date : 2016-07-16 14:52
 * description :
 * history :
 */
package afterSales

type (
	// 换货单接口
	IExchangeOrder interface {
		// 将换货的商品重新发货
		ExchangeShip(spName string, spOrder string) error

		// 消费者延长收货时间
		LongReceive() error

		// 接收换货
		ExchangeReceive() error
	}

	// 换货单
	ExchangeOrder struct {
		// 编号
		Id int32 `db:"id" pk:"yes" auto:"no"`
		// 是否发货
		IsShipped int `db:"is_shipped"`
		// 快递名称
		ShipSpName string `db:"sp_name"`
		// 快递编号
		ShipSpOrder string `db:"sp_order"`
		// 发货时间
		ShipTime int64 `db:"ship_time"`
		// 是否收货
		IsReceived int `db:"is_received"`
		// 收货时间
		ReceiveTime int64 `db:"receive_time"`
	}
)
