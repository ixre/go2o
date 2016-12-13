/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-08 10:23
 * description :
 * history :
 */

package cart

import (
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/domain/interface/sale"
	"go2o/core/dto"
	"go2o/core/infrastructure/domain"
	"go2o/core/infrastructure/format"
)

var (
	ErrEmptyShoppingCart *domain.DomainError = domain.NewDomainError(
		"empty_shopping_cart", "购物车没有商品")

	ErrCartNoBuyer *domain.DomainError = domain.NewDomainError(
		"err_cart_no_buyer", "购物车未绑定")

	ErrCartBuyerBind *domain.DomainError = domain.NewDomainError(
		"err_cart_buyer_binded", "购物车已绑定")

	ErrDisallowBindForCart *domain.DomainError = domain.NewDomainError(
		"err_cart_disallow_bind", "无法为购物车绑定订单")
)

type (
	ICart interface {
		// 获取聚合根编号
		GetAggregateRootId() int32

		// 获取购物车的KEY
		Key() string

		// 获取购物车值
		GetValue() ValueCart

		// 标记商品结算
		SignItemChecked(skuArr []int32) error

		// 检查购物车(仅结算商品)
		Check() error

		// 获取商品编号与购物车项的集合
		Items() map[int32]*CartItem

		// 获取购物车中的商品
		GetCartGoods() []sale.IGoods

		// 结算数据持久化
		SettlePersist(shopId, paymentOpt, deliverOpt, deliverId int32) error

		// 获取结算数据
		GetSettleData() (s shop.IShop, d member.IDeliverAddress, paymentOpt, deliverOpt int32)

		// 设置购买会员
		SetBuyer(buyerId int32) error

		// 设置购买会员收货地址
		SetBuyerAddress(addressId int32) error

		// 添加项,需传递商户编号、店铺编号
		// todo: 这里有问题、如果是线下店的购物车,如何实现?
		AddItem(vendorId, shopId, skuId int32, num int, checked bool) (*CartItem, error)

		// 移出项
		RemoveItem(skuId int32, num int) error

		// 合并购物车，并返回新的购物车
		Combine(ICart) ICart

		// 保存购物车
		Save() (int32, error)

		// 释放购物车,如果购物车的商品全部结算,则返回true
		Release() bool

		// 销毁购物车
		Destroy() error

		// 获取汇总信息
		GetSummary() string

		// 获取Json格式的商品数据
		GetJsonItems() []byte

		// 获取金额
		GetFee() (totalFee float32, orderFee float32)
	}

	// 根据数据获取购物车,
	// 如果member的cart与key不一致，则合并购物车；
	// 如果会员没有购物车，则绑定为key的购物车
	// 如果都没有，则创建一个购物车
	ICartRepo interface {
		// 创建购物车对象
		CreateCart(v *ValueCart) ICart

		// 创建一个购物车
		NewCart() ICart

		// 获取购物车
		GetShoppingCartByKey(key string) ICart

		// 获取会员没有结算的购物车
		GetMemberCurrentCart(buyerId int32) ICart

		// 获取购物车
		GetShoppingCart(key string) *ValueCart

		// 获取最新的购物车
		GetLatestCart(buyerId int32) *ValueCart

		// 保存购物车
		SaveShoppingCart(*ValueCart) (int32, error)

		// 移出购物车项
		RemoveCartItem(id int32) error

		// 保存购物车项
		SaveCartItem(*CartItem) (int32, error)

		// 清空购物车项
		EmptyCartItems(cartId int32) error

		// 删除购物车
		DeleteCart(cartId int32) error
	}

	//todo:  shopId应去掉,同时应存储邮费等信息
	ValueCart struct {
		Id      int32  `db:"id" pk:"yes" auto:"yes"`
		CartKey string `db:"cart_key"`
		BuyerId int32  `db:"buyer_id"`
		//OrderNo    string           `db:"order_no"`
		//IsBought   int              `db:"is_bought"`
		PaymentOpt int32       `db:"payment_opt"`
		DeliverOpt int32       `db:"deliver_opt"`
		DeliverId  int32       `db:"deliver_id"`
		ShopId     int32       `db:"shop_id"`
		CreateTime int64       `db:"create_time"`
		UpdateTime int64       `db:"update_time"`
		Items      []*CartItem `db:"-"`
	}

	// 购物车项
	CartItem struct {
		Id         int32          `db:"id" pk:"yes" auto:"yes"`
		CartId     int32          `db:"cart_id"`
		VendorId   int32          `db:"vendor_id"`
		ShopId     int32          `db:"shop_id"`
		SkuId      int32          `db:"goods_id"`
		SnapshotId int32          `db:"snap_id"`
		Quantity   int            `db:"quantity"`
		Checked    int            `db:"checked" json:"checked"` // 是否结算
		Snapshot   *item.Snapshot `db:"-"`

		Price      float32 `db:"-"`
		SalePrice  float32 `db:"-"`
		Sku        string  `db:"-"`
		Name       string  `db:"-"`
		GoodsNo    string  `db:"-"`
		SmallTitle string  `db:"-"`
		Image      string  `db:"-"`
	}
)

func ParseCartItem(item *CartItem) *dto.CartItem {
	i := &dto.CartItem{
		GoodsId:    item.SkuId,
		GoodsName:  item.Name,
		GoodsNo:    item.GoodsNo,
		SmallTitle: item.SmallTitle,
		GoodsImage: format.GetGoodsImageUrl(item.Image),
		Quantity:   item.Quantity,
		Price:      item.Price,
		SalePrice:  item.SalePrice,
	}
	if item.Checked == 1 {
		i.Checked = true
	}
	return i
}

func ParseToDtoCart(c ICart) *dto.ShoppingCart {
	cart := &dto.ShoppingCart{}
	v := c.GetValue()
	cart.Id = c.GetAggregateRootId()
	cart.BuyerId = v.BuyerId
	cart.CartKey = v.CartKey
	cart.UpdateTime = v.UpdateTime
	t, f := c.GetFee()
	cart.TotalFee = t
	cart.OrderFee = f
	cart.Summary = c.GetSummary()
	cart.Vendors = []*dto.CartVendorGroup{}

	if v.Items != nil {
		if l := len(v.Items); l > 0 {
			mp := make(map[int32]*dto.CartVendorGroup, 0) //保存运营商到map
			for _, v := range v.Items {
				vendor, ok := mp[v.ShopId]
				if !ok {
					vendor = &dto.CartVendorGroup{
						VendorId: v.VendorId,
						ShopId:   v.ShopId,
						Items:    []*dto.CartItem{},
					}
					mp[v.ShopId] = vendor
					cart.Vendors = append(cart.Vendors, vendor)
				}
				if v.Checked == 1 {
					vendor.CheckedNum += 1
				}
				vendor.Items = append(vendor.Items, ParseCartItem(v))
				cart.TotalNum += v.Quantity
			}
		}
	}

	return cart
}
