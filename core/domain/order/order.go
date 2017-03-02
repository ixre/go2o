package order

import (
	"go2o/core/domain/interface/express"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/order"
	"go2o/core/domain/interface/payment"
	"go2o/core/domain/interface/product"
	"go2o/core/domain/interface/promotion"
	"go2o/core/domain/interface/valueobject"
	"time"
)

var _ order.IOrder = new(baseOrderImpl)

type baseOrderImpl struct {
	baseValue  *order.Order
	buyer      member.IMember
	repo       order.IOrderRepo
	memberRepo member.IMemberRepo
	manager    order.IOrderManager
}

// 获取编号
func (o *baseOrderImpl) GetAggregateRootId() int64 {
	return o.baseValue.ID
}

// 订单类型
func (o *baseOrderImpl) Type() order.OrderType {
	return order.OrderType(o.baseValue.OrderType)
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
	if o.GetAggregateRootId() <= 0 {
		o.baseValue.OrderNo = o.manager.GetFreeOrderNo(0)
		o.baseValue.CreateTime = time.Now().Unix()
		o.baseValue.State = order.StatAwaitingPayment
	}
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

// 工厂方法生成订单
func FactoryNew(v *order.Order, manager order.IOrderManager,
	repo order.IOrderRepo, mchRepo merchant.IMerchantRepo,
	goodsRepo item.IGoodsItemRepo, productRepo product.IProductRepo,
	promRepo promotion.IPromotionRepo, memberRepo member.IMemberRepo,
	expressRepo express.IExpressRepo, payRepo payment.IPaymentRepo,
	valRepo valueobject.IValueRepo) order.IOrder {
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
	}
	return nil
}
