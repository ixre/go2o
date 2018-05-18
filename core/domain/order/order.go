package order

import (
	"errors"
	"go2o/core/domain/interface/cart"
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
	"go2o/core/infrastructure/format"
	"time"
)

// 订单商品项(领域内部使用)
type orderItem struct {
	// 编号
	ID int32
	// 订单编号
	OrderId int64
	// 商品编号
	ItemId int64
	// 商品SKU编号
	SkuId int64
	// 快照编号
	SnapshotId int64
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
	itemRepo   item.IGoodsItemRepo
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

// 通过订单创建购物车
func (o *baseOrderImpl) BuildCart() cart.ICart {
	//todo: 实现批发等订单的构造购物车
	panic("implement in sub class")
}

// 获取订单号
func (o *baseOrderImpl) OrderNo() string {
	return o.baseValue.OrderNo
}

// 保存订单信息
func (o *baseOrderImpl) saveOrder() error {
	id, err := o.repo.SaveOrder(o.baseValue)
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

// 绑定商品信息
func (o *baseOrderImpl) bindItemInfo(i *order.ComplexItem) {
	unitPrice := i.FinalAmount / float64(i.Quantity)
	i.Data["UnitPrice"] = format.DecimalToString(unitPrice)
	it := o.itemRepo.GetSalesSnapshot(i.SnapshotId)
	i.Data["ItemImage"] = it.Image
	i.Data["ItemName"] = it.GoodsTitle
	i.Data["SpecWord"] = it.Sku
	//todo: ??  SKU货号，供打出订单
}

// 复合的订单信息
func (o *baseOrderImpl) Complex() *order.ComplexOrder {
	if o.complex == nil {
		o.complex = &order.ComplexOrder{
			OrderId:    o.GetAggregateRootId(),
			OrderType:  o.baseValue.OrderType,
			OrderNo:    o.OrderNo(),
			BuyerId:    o.baseValue.BuyerId,
			State:      o.baseValue.State,
			CreateTime: o.baseValue.CreateTime,
			Data:       make(map[string]string),
		}
	}
	return o.complex
}

// 生成支付单
func (o *baseOrderImpl) createPaymentOrder() *payment.Order {
	orderId := o.GetAggregateRootId()
	if orderId <= 0 {
		panic("payment order must create after order submit!")
	}
	buyerId := o.Buyer().GetAggregateRootId()
	unix := time.Now().Unix()
	v2 := &payment.Order{
		ID:             0,
		SellerId:       0,
		TradeType:      "",
		TradeNo:        o.OrderNo(),
		SubOrder:       0,
		OrderType:      int(order.TRetail),
		OutOrderNo:     o.OrderNo(),
		Subject:        "支付订单",
		BuyerId:        buyerId,
		PayUid:         buyerId,
		TotalAmount:    0,
		DiscountAmount: 0,
		DeductAmount:   0,
		AdjustAmount:   0,
		FinalFee:       0,
		PaymentFlag:    payment.PAllFlag,
		ExtraData:      "",
		TradeChannel:   0,
		OutTradeSp:     "",
		OutTradeNo:     "",
		State:          payment.StateAwaitingPayment,
		SubmitTime:     unix,
		ExpiresTime:    0,
		PaidTime:       0,
		UpdateTime:     unix,
		TradeChannels:  make([]*payment.TradeChan, 0),
	}
	return v2
}

// 工厂方法生成订单
func FactoryOrder(v *order.Order, manager order.IOrderManager,
	repo order.IOrderRepo, mchRepo merchant.IMerchantRepo,
	itemRepo item.IGoodsItemRepo, productRepo product.IProductRepo,
	promRepo promotion.IPromotionRepo, memberRepo member.IMemberRepo,
	expressRepo express.IExpressRepo, shipRepo shipment.IShipmentRepo,
	payRepo payment.IPaymentRepo, cartRepo cart.ICartRepo,
	valRepo valueobject.IValueRepo) order.IOrder {
	b := &baseOrderImpl{
		baseValue:  v,
		repo:       repo,
		itemRepo:   itemRepo,
		memberRepo: memberRepo,
		manager:    manager,
	}
	t := order.OrderType(v.OrderType)
	switch t {
	case order.TRetail:
		return newNormalOrder(manager, b, repo, itemRepo,
			productRepo, promRepo, expressRepo,
			payRepo, cartRepo, valRepo)
	case order.TWholesale:
		return newWholesaleOrder(b, repo, itemRepo,
			expressRepo, payRepo, shipRepo, mchRepo, valRepo)
	case order.TTrade:
		return newTradeOrder(b, payRepo, mchRepo, valRepo)
	}
	return nil
}
