package order

import (
	"errors"
	"go2o/core/domain/interface/enum"
	"go2o/core/domain/interface/express"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/order"
	"go2o/core/domain/interface/payment"
	"go2o/core/domain/interface/product"
	"go2o/core/domain/interface/promotion"
	"go2o/core/domain/interface/shipment"
	"go2o/core/domain/interface/valueobject"
	"time"
)

// 订单商品项(领域内部使用)
type orderItem struct {
	// 编号
	ID int32
	// 订单编号
	OrderId int64
	// 商品编号
	ItemId int32
	// 商品SKU编号
	SkuId int32
	// 快照编号
	SnapshotId int32
	// 数量
	Quantity int32
	// 退回数量(退货)
	ReturnQuantity int32
	// SKU描述
	//Sku string `db:"sku"`
	// 金额
	Amount float32
	// 最终金额, 可能会有优惠均摊抵扣的金额
	FinalAmount float32
	// 是否发货
	IsShipped int
	// 更新时间
	UpdateTime int64
	// 运营商编号
	VendorId int32
	// 商店编号
	ShopId int32
	// 重量,用于生成订单时存储数据
	Weight int32
	// 体积:毫升(ml)
	Bulk int32
	// 快递模板编号
	ExpressTplId int32
}

var _ order.IOrder = new(baseOrderImpl)

type baseOrderImpl struct {
	baseValue  *order.Order
	buyer      member.IMember
	repo       order.IOrderRepo
	memberRepo member.IMemberRepo
	manager    order.IOrderManager
	complex    *order.ComplexOrder
}

// 获取编号
func (o *baseOrderImpl) GetAggregateRootId() int64 {
	return o.baseValue.ID
}

// 订单类型
func (o *baseOrderImpl) Type() order.OrderType {
	return order.OrderType(o.baseValue.OrderType)
}

// 获取订单状态
func (o *baseOrderImpl) State() order.OrderState {
	return order.OrderState(o.baseValue.State)
}

// 获取购买的会员
func (o *baseOrderImpl) Buyer() member.IMember {
	if o.buyer == nil {
		o.buyer = o.memberRepo.GetMember(o.baseValue.BuyerId)
	}
	return o.buyer
}

// 提交订单。如遇拆单,需均摊优惠抵扣金额到商品
func (o *baseOrderImpl) Submit() error {
	if o.GetAggregateRootId() > 0 {
		return errors.New("订单不允许重复提交")
	}
	if o.baseValue.OrderNo == "" {
		o.baseValue.OrderNo = o.manager.GetFreeOrderNo(0)
	}
	if o.baseValue.State == 0 {
		o.baseValue.State = order.StatAwaitingPayment
	}
	o.baseValue.CreateTime = time.Now().Unix()
	return o.saveOrder()
}

// 获取订单号
func (o *baseOrderImpl) OrderNo() string {
	return o.baseValue.OrderNo
}

// 保存订单信息
func (o *baseOrderImpl) saveOrder() error {
	id, err := o.repo.SaveOrderList(o.baseValue)
	if err == nil {
		o.baseValue.ID = int64(id)
	}
	return err
}

// 设置并订单状态
func (o *baseOrderImpl) saveOrderState(state order.OrderState) {
	if o.baseValue.State != int32(order.StatBreak) {
		o.baseValue.State = int32(state)
		o.saveOrder()
	}
}

// 复合的订单信息
func (o *baseOrderImpl) Complex() *order.ComplexOrder {
	if o.complex == nil {
		o.complex = &order.ComplexOrder{
			Id:         o.GetAggregateRootId(),
			OrderType:  o.baseValue.OrderType,
			OrderNo:    o.OrderNo(),
			BuyerId:    o.baseValue.BuyerId,
			State:      o.baseValue.State,
			CreateTime: o.baseValue.CreateTime,
		}
	}
	return o.complex
}

// 生成支付单
func (o *baseOrderImpl) createPaymentOrder() *payment.PaymentOrder {
	orderId := o.GetAggregateRootId()
	if orderId <= 0 {
		panic("payment order must create after order submit!")
	}
	buyerId := o.Buyer().GetAggregateRootId()
	v := &payment.PaymentOrder{
		BuyUser:        buyerId,
		PaymentUser:    buyerId,
		VendorId:       0,
		OrderId:        int32(orderId),
		Type:           payment.TypeShopping,
		PaymentOptFlag: payment.OptPerm,
		PaymentSign:    enum.PaymentOnlinePay,
		CreateTime:     o.baseValue.CreateTime,
		TradeNo:        o.OrderNo(),
		State:          payment.StateAwaitingPayment,
	}
	v.FinalAmount = v.TotalFee - v.SubAmount - v.SystemDiscount -
		v.IntegralDiscount - v.BalanceDiscount
	return v
}

// 工厂方法生成订单
func FactoryNew(v *order.Order, manager order.IOrderManager,
	repo order.IOrderRepo, mchRepo merchant.IMerchantRepo,
	goodsRepo item.IGoodsItemRepo, productRepo product.IProductRepo,
	promRepo promotion.IPromotionRepo, memberRepo member.IMemberRepo,
	expressRepo express.IExpressRepo, shipRepo shipment.IShipmentRepo,
	payRepo payment.IPaymentRepo, valRepo valueobject.IValueRepo) order.IOrder {
	b := &baseOrderImpl{
		baseValue:  v,
		repo:       repo,
		memberRepo: memberRepo,
		manager:    manager,
	}
	t := order.OrderType(v.OrderType)
	switch t {
	case order.TRetail:
		return newNormalOrder(manager, b, repo, goodsRepo,
			productRepo, promRepo, expressRepo,
			payRepo, valRepo)
	case order.TWholesale:
		return newWholesaleOrder(b, repo, goodsRepo,
			expressRepo, payRepo, shipRepo, mchRepo, valRepo)
	}
	return nil
}
