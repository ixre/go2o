/**
 * Copyright 2015 @ z3q.net.
 * name : shipment
 * author : jarryliu
 * date : 2016-07-15 08:19
 * description :
 * history :
 */
package shipment

type (
	IShipmentOrder interface {
		// 获取聚合根编号
		GetAggregateRootId() int
		// 获取值
		Value() ShipmentOrder
		// 发货
		Ship(spId int, spOrderNo string) error
		// 发货完成
		Shipped() error
		// 更新快递记录
		UpdateLog() error
	}

	IShipmentRep interface {
		// 创建发货单
		CreateShipmentOrder(o *ShipmentOrder) IShipmentOrder

		// 获取发货单
		GetShipmentOrder(id int) IShipmentOrder

		// 获取订单对应的发货单
		GetOrders(orderId int) []IShipmentOrder

		// 保存发货单
		SaveShipmentOrder(o *ShipmentOrder) (int, error)

		// 删除发货单
		DeleteShipmentOrder(id int) error
	}

	ShipmentOrder struct {
		//  发货单编号
		Id int `db:"id" pk:"yes" auto:"yes"`
		// 订单编号
		OrderId int `db:"order_id"`
		// 快递服务商编号
		SpId int `db:"sp_id"`
		// 快递单号
		SpOrderNo string `db:"sp_order"`
		// 物流日志
		ExpressLog string `db:"exporess_log"`
		// 货物金额
		Amount float32 `db:"amount"`
		// 货物实际金额
		FinalAmount float32 `db:"final_amount"`
		// 发货时间
		ShipTime int64 `db:"shop_time"`
		// 是否收货
		IsShipped int `db:"is_shipped"`
		// 更新时间
		UpdateTime int64 `db:update_time"`
		// 配送项目
		ShipmentItem []*Item `db:"-"`
	}

	Item struct {
		Id int `db:"id" auto:"yes" pk:"yes"`
		// 发货单编号
		OrderId int `db:"order_id"`
		// 商品销售快照编号
		GoodsSnapId int `db:"snap_id"`
		// 数量
		Quantity int `db:"quantity"`
		// 商品金额
		Amount float32 `db:"amount"`
		// 商品实际金额
		FinalAmount float32 `db:"final_amount"`
	}
)
