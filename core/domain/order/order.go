package order

import (
	"errors"
	"strings"
	"time"

	"github.com/ixre/go2o/core/domain/interface/cart"
	"github.com/ixre/go2o/core/domain/interface/express"
	"github.com/ixre/go2o/core/domain/interface/item"
	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/domain/interface/merchant"
	"github.com/ixre/go2o/core/domain/interface/merchant/shop"
	"github.com/ixre/go2o/core/domain/interface/order"
	"github.com/ixre/go2o/core/domain/interface/payment"
	"github.com/ixre/go2o/core/domain/interface/product"
	"github.com/ixre/go2o/core/domain/interface/promotion"
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/domain/interface/shipment"
	"github.com/ixre/go2o/core/domain/interface/valueobject"
	"github.com/ixre/go2o/core/infrastructure/format"
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
	Amount int64
	// 最终金额, 可能会有优惠均摊抵扣的金额
	FinalAmount int64
	// 是否发货
	IsShipped int
	// 更新时间
	UpdateTime int64
	// 运营商编号
	VendorId int64
	// 商店编号
	ShopId int64
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
	itemRepo   item.IItemRepo
	manager    order.IOrderManager
	complex    *order.ComplexOrder
}

// GetAggregateRootId 获取编号
func (o *baseOrderImpl) GetAggregateRootId() int64 {
	return o.baseValue.Id
}

// Type 订单类型
func (o *baseOrderImpl) Type() order.OrderType {
	return order.OrderType(o.baseValue.OrderType)
}

// State 获取订单状态
func (o *baseOrderImpl) State() order.OrderStatus {
	return order.OrderStatus(o.baseValue.Status)
}

// Buyer 获取购买的会员
func (o *baseOrderImpl) Buyer() member.IMember {
	if o.buyer == nil {
		o.buyer = o.memberRepo.GetMember(o.baseValue.BuyerId)
	}
	return o.buyer
}

// SetShipmentAddress 设置配送地址
func (o *baseOrderImpl) SetShipmentAddress(addressId int64) error {
	if addressId <= 0 {
		return order.ErrNoSuchAddress
	}
	buyer := o.Buyer()
	if buyer == nil {
		return member.ErrNoSuchMember
	}
	addr := buyer.Profile().GetAddress(addressId)
	if addr == nil {
		return order.ErrNoSuchAddress
	}
	d := addr.GetValue()
	o.baseValue.ShippingAddress = strings.Replace(d.Area, " ", "", -1) + d.DetailAddress
	o.baseValue.ConsigneeName = d.ConsigneeName
	o.baseValue.ConsigneePhone = d.ConsigneePhone
	return nil
}

// GetPaymentOrder implements order.IOrder
func (*baseOrderImpl) GetPaymentOrder() payment.IPaymentOrder {
	panic("unimplemented")
}

// Submit 提交订单。如遇拆单,需均摊优惠抵扣金额到商品
func (o *baseOrderImpl) Submit() error {
	if o.GetAggregateRootId() > 0 {
		return errors.New("订单不允许重复提交")
	}
	if o.baseValue.OrderNo == "" {
		o.baseValue.OrderNo = o.manager.GetFreeOrderNo(0)
	}
	if o.baseValue.Status == 0 {
		o.baseValue.Status = order.StatAwaitingPayment
	}
	m := o.Buyer().GetValue()
	o.baseValue.BuyerUser = m.User
	o.baseValue.CreateTime = time.Now().Unix()
	return o.saveOrder()
}

// BuildCart 通过订单创建购物车
func (o *baseOrderImpl) BuildCart() cart.ICart {
	//todo: 实现批发等订单的构造购物车
	panic("implement in sub class")
}

// OrderNo 获取订单号
func (o *baseOrderImpl) OrderNo() string {
	return o.baseValue.OrderNo
}

// 保存订单信息
func (o *baseOrderImpl) saveOrder() error {
	id, err := o.repo.SaveOrder(o.baseValue)
	if err == nil {
		o.baseValue.Id = int64(id)
	}
	return err
}

// 设置并订单状态
func (o *baseOrderImpl) saveOrderState(state order.OrderStatus) {
	if state == order.StatBreak {
		o.baseValue.IsBreak = 1
	}
	if o.baseValue.Status != int(state) {
		o.baseValue.Status = int(state)
		_ = o.saveOrder()
	}
}

// 绑定商品信息
func (o *baseOrderImpl) bindItemInfo(i *order.ComplexItem) {
	unitPrice := float64(i.FinalAmount) / float64(i.Quantity) / 100
	i.Data["UnitPrice"] = format.FormatFloat64(unitPrice)
	it := o.itemRepo.GetSalesSnapshot(i.SnapshotId)
	i.Data["ItemImage"] = it.Image
	i.Data["ItemName"] = it.GoodsTitle
	i.Data["SpecWord"] = it.Sku
	//todo: ??  SKU货号，供打出订单
}

// Complex 复合的订单信息
func (o *baseOrderImpl) Complex() *order.ComplexOrder {
	if o.complex == nil {
		o.complex = &order.ComplexOrder{
			OrderId:        o.GetAggregateRootId(),
			OrderType:      int32(o.baseValue.OrderType),
			OrderNo:        o.OrderNo(),
			BuyerId:        o.baseValue.BuyerId,
			BuyerUser:      o.baseValue.BuyerUser,
			Subject:        o.baseValue.Subject,
			ItemCount:      o.baseValue.ItemCount,
			ItemAmount:     o.baseValue.ItemAmount,
			DiscountAmount: o.baseValue.DiscountAmount,
			ExpressFee:     o.baseValue.ExpressFee,
			PackageFee:     o.baseValue.PackageFee,
			FinalAmount:    o.baseValue.FinalAmount,
			IsBreak:        int32(o.baseValue.IsBreak),
			Status:         int32(o.baseValue.Status),
			StateText:      "",
			CreateTime:     o.baseValue.CreateTime,
			UpdateTime:     o.baseValue.UpdateTime,
			Data:           make(map[string]string),
			Details:        []*order.ComplexOrderDetails{},
			Consignee: &order.ComplexConsignee{
				ConsigneeName:   o.baseValue.ConsigneeName,
				ConsigneePhone:  o.baseValue.ConsigneePhone,
				ShippingAddress: o.baseValue.ShippingAddress,
			},
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
		Id:             0,
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
		PayFlag:        payment.PAllFlag,
		ExtraData:      "",
		TradeChannel:   0,
		OutTradeSp:     "",
		OutTradeNo:     "",
		State:          payment.StateAwaitingPayment,
		SubmitTime:     unix,
		ExpiresTime:    0,
		PaidTime:       0,
		UpdateTime:     unix,
		TradeMethods:   make([]*payment.TradeMethodData, 0),
	}
	return v2
}

// FactoryOrder 工厂方法生成订单
func FactoryOrder(v *order.Order, manager order.IOrderManager,
	repo order.IOrderRepo, mchRepo merchant.IMerchantRepo,
	itemRepo item.IItemRepo, productRepo product.IProductRepo,
	promRepo promotion.IPromotionRepo, memberRepo member.IMemberRepo,
	expressRepo express.IExpressRepo, shipRepo shipment.IShipmentRepo,
	payRepo payment.IPaymentRepo, cartRepo cart.ICartRepo,
	shopRepo shop.IShopRepo, valRepo valueobject.IValueRepo,
	registryRepo registry.IRegistryRepo) order.IOrder {
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
			payRepo, cartRepo, shopRepo, registryRepo, valRepo)
	case order.TWholesale:
		return newWholesaleOrder(b, repo, itemRepo,
			expressRepo, payRepo, shipRepo, mchRepo, shopRepo, valRepo, registryRepo)
	case order.TTrade:
		return newTradeOrder(b, payRepo, mchRepo, valRepo, registryRepo)
	}
	return nil
}
