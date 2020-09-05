/**
 * Copyright 2015 @ to2.net.
 * name : shipment
 * author : jarryliu
 * date : 2016-07-15 08:19
 * description :
 * history :
 */
package shipment

import (
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/infrastructure/domain"
)

const (
	// 等待发货
	StatAwaitingShipment = 1 + iota
	// 已发货
	StatShipped
	// 已完成
	StatCompleted
)

var (
	ErrNotSetExpressTemplate *domain.DomainError = domain.NewError(
		"err_not_set_express_tpl", "请设置运费模板")
	ErrMissingSpInfo *domain.DomainError = domain.NewError(
		"err_shipment_missing_spinfo", "物流信息不完整、无法进行发货!",
	)
)

type (
	IShipmentOrder interface {
		// 获取聚合根编号
		GetAggregateRootId() int64
		// 获取值
		Value() ShipmentOrder
		// 发货
		Ship(spId int32, spOrderNo string) error
		// 发货完成
		Completed() error
		// 更新快递记录
		UpdateLog() error
		// 智能选择门店
		SmartChoiceShop(address string) (shop.IShop, error)
	}

	IShipmentRepo interface {
		// 创建发货单
		CreateShipmentOrder(o *ShipmentOrder) IShipmentOrder
		// 获取发货单
		GetShipmentOrder(id int64) IShipmentOrder
		// 获取订单对应的发货单
		GetShipOrders(orderId int64, sub bool) []IShipmentOrder
		// 保存发货单
		SaveShipmentOrder(o *ShipmentOrder) (int, error)
		// 保存发货商品项
		SaveShipmentItem(v *Item) (int, error)
		// 删除发货单
		DeleteShipmentOrder(id int64) error
	}

	// 发货单
	ShipmentOrder struct {
		// 编号
		ID int64 `db:"id" pk:"yes" auto:"yes"`
		// 订单编号
		OrderId int64 `db:"order_id"`
		// 子订单编号
		SubOrderId int64 `db:"sub_orderid"`
		// 快递SP编号
		SpId int32 `db:"sp_id"`
		// 快递SP单号
		SpOrder string `db:"sp_order"`
		// 物流日志
		ShipmentLog string `db:"shipment_log"`
		// 运费
		Amount float64 `db:"amount"`
		// 实际运费
		FinalAmount float64 `db:"final_amount"`
		// 发货时间
		ShipTime int64 `db:"ship_time"`
		// 状态
		State int `db:"state"`
		// 更新时间
		UpdateTime int64 `db:"update_time"`
		// 配送项目
		Items []*Item `db:"-"`
	}
	// 发货单详情
	Item struct {
		// 编号
		ID int64 `db:"id" pk:"yes" auto:"yes"`
		// 发货单编号
		ShipOrder int64 `db:"ship_order"`
		// 商品交易快照编号
		SnapshotId int64 `db:"snapshot_id"`
		// 商品数量
		Quantity int32 `db:"quantity"`
		// 运费
		Amount float64 `db:"amount"`
		// 实际运费
		FinalAmount float64 `db:"final_amount"`
	}
)
