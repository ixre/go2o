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
)

var _ order.IOrder = new(baseOrderImpl)

type baseOrderImpl struct {
	baseValue *order.Order
	repo      order.IOrderRepo
}

// 获取编号
func (o *baseOrderImpl) GetAggregateRootId() int64 {
	return o.baseValue.ID
}

// 订单类型
func (o *baseOrderImpl) Type() order.OrderType {
	return order.OrderType(o.baseValue.OrderType)
}

// 提交订单。如遇拆单,需均摊优惠抵扣金额到商品
func (o *baseOrderImpl) Submit() error {
	panic("not implement!")
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
		baseValue: v,
		repo:      repo,
	}
	t := order.OrderType(v.OrderType)
	switch t {
	case order.TRetail:
		return newOrder(manager, b, mchRepo, repo, goodsRepo,
			productRepo, promRepo, memberRepo, expressRepo,
			payRepo, valRepo)
	}
	return nil
}
