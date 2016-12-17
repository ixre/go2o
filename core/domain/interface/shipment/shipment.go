/**
 * Copyright 2015 @ z3q.net.
 * name : shipment
 * author : jarryliu
 * date : 2016-07-15 08:19
 * description :
 * history :
 */
package shipment

import "go2o/core/infrastructure/domain"

const (
	// 等待发货
	StatAwaitingShipment = 1 + iota
	// 已发货
	StatShipped
	// 已完成
	StatCompleted
)

var (
	ErrNotSetExpressTemplate *domain.DomainError = domain.NewDomainError(
		"err_not_set_express_tpl", "请设置运费模板")
	ErrMissingSpInfo *domain.DomainError = domain.NewDomainError(
		"err_shipment_missing_spinfo", "物流信息不完整、无法进行发货!",
	)
)

type (
	IShipmentOrder interface {
		// 获取聚合根编号
		GetAggregateRootId() int32
		// 获取值
		Value() ShipmentOrder
		// 发货
		Ship(spId int32, spOrderNo string) error
		// 发货完成
		Completed() error
		// 更新快递记录
		UpdateLog() error
	}

	IShipmentRepo interface {
		// 创建发货单
		CreateShipmentOrder(o *ShipmentOrder) IShipmentOrder

		// 获取发货单
		GetShipmentOrder(id int32) IShipmentOrder

		// 获取订单对应的发货单
		GetOrders(orderId int32) []IShipmentOrder

		// 保存发货单
		SaveShipmentOrder(o *ShipmentOrder) (int32, error)

		// 保存发货商品项
		SaveShipmentItem(v *Item) (int32, error)

		// 删除发货单
		DeleteShipmentOrder(id int32) error
	}

	ShipmentOrder struct {
		//  发货单编号
		Id int32 `db:"id" pk:"yes" auto:"yes"`
		// 订单编号
		OrderId int32 `db:"order_id"`
		// 快递服务商编号
		SpId int32 `db:"sp_id"`
		// 快递单号
		SpOrderNo string `db:"sp_order"`
		// 物流日志
		ExpressLog string `db:"exporess_log"`
		// 货物金额
		Amount float32 `db:"amount"`
		// 货物实际金额
		FinalAmount float32 `db:"final_amount"`
		// 发货时间
		ShipTime int64 `db:"ship_time"`
		// 状态
		Stat int `db:"state"`
		// 更新时间
		UpdateTime int64 `db:update_time"`
		// 配送项目
		Items []*Item `db:"-"`
	}

	Item struct {
		Id int32 `db:"id" auto:"yes" pk:"yes"`
		// 发货单编号
		OrderId int32 `db:"ship_order"`
		// 商品销售快照编号
		GoodsSnapId int32 `db:"snap_id"`
		// 数量
		Quantity int32 `db:"quantity"`
		// 商品金额
		Amount float32 `db:"amount"`
		// 商品实际金额
		FinalAmount float32 `db:"final_amount"`
	}
)
