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
	"github.com/ixre/gof/math"
	"github.com/ixre/gof/util"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/infrastructure/domain"
	"go2o/core/infrastructure/format"
	"go2o/core/service/auto_gen/rpc/ttype"
)

var (
	ErrNoSuchCart *domain.DomainError = domain.NewError(
		"empty_shopping_no_such_cart", "购物车无法使用")

	ErrKindNotMatch *domain.DomainError = domain.NewError(
		"err_cart_kind_not_match", "购物车类型不匹配")

	ErrEmptyShoppingCart *domain.DomainError = domain.NewError(
		"empty_shopping_cart", "购物车没有商品")

	ErrNoChecked *domain.DomainError = domain.NewError(
		"empty_shopping_cart_no_checked", "购物车没有结算的商品")

	ErrCartNoBuyer *domain.DomainError = domain.NewError(
		"err_cart_no_buyer", "购物车未绑定")

	ErrItemNoSku *domain.DomainError = domain.NewError(
		"err_cart_item_no_sku", "请选择商品规格")
)

const (
	// 普通(B2C)购物车
	KNormal CartKind = 1
	// 零售(B2C-线下)购物车
	KRetail CartKind = 2
	// 批发(B2B)购物车
	KWholesale CartKind = 3
)

const (
	// 手工创建
	FlagManualCreate = 1
	// 是否可改价
	FlagPriceEditable = 2
)

type (
	// 购物车类型
	CartKind int
	// 购物车
	ICart interface {
		// 获取聚合根编号
		GetAggregateRootId() int32
		// 购物车种类
		Kind() CartKind
		// 克隆
		Clone() ICart
		// 获取购物车编码
		Code() string
		// 获取买家编号
		BuyerId() int64
		// 预先准备购物车
		Prepare() error
		// 标记商品结算
		SignItemChecked(items []*ItemPair) error
		// 获取勾选的商品,checked:为商品与商品SKU数据
		CheckedItems(checked map[int64][]int64) []*ItemPair

		// 添加商品到购物车,如商品没有SKU,则skuId传入0
		// todo: 这里有问题、如果是线下店的购物车,如何实现?
		// 暂时以店铺区分,2017-02-28考虑单独的购物车或子系统
		Put(itemId, skuId int64, quantity int32) error
		// 更新商品数量，如数量为0，则删除
		Update(itemId, skuId int64, quantity int32) error
		// 移出项
		Remove(itemId, skuId int64, quantity int32) error
		// 保存购物车
		Save() (int32, error)
		// 释放购物车,如果购物车的商品全部结算,则返回true
		Release(checked map[int64][]int64) bool
		// 销毁购物车
		Destroy() error

		// 结算数据持久化
		SettlePersist(shopId, paymentOpt, deliverOpt int32, addressId int64) error
		// 获取结算数据
		GetSettleData() (s shop.IShop, d member.IDeliverAddress, paymentOpt int32)

		// 设置购买会员收货地址
		SetBuyerAddress(addressId int64) error
	}

	//商品普通购物车,未登陆时以code标识，登陆后以买家编号标识
	INormalCart interface {
		// 获取购物车值
		Value() NormalCart
		// 获取商品集合
		Items() []*NormalCartItem
		// 合并购物车，并返回新的购物车
		Combine(ICart) ICart
		// 获取项
		GetItem(itemId, skuId int64) *NormalCartItem
	}

	//商品批发购物车
	IWholesaleCart interface {
		// 获取购物车值
		GetValue() WsCart
		// 获取商品集合
		Items() []*WsCartItem
		// Jdo数据
		JdoData(checkout bool, checked map[int64][]int64) *WCartJdo
		// 简单Jdo数据,max为最多数量
		QuickJdoData(max int) string
	}

	// 根据数据获取购物车,
	// 如果member的cart与key不一致，则合并购物车；
	// 如果会员没有购物车，则绑定为key的购物车
	// 如果都没有，则创建一个购物车
	ICartRepo interface {
		// 获取买家的购物车
		GetMyCart(buyerId int64, k CartKind) ICart
		// 创建一个购物车
		NewNormalCart(code string) ICart
		// 创建一个普通购物车
		CreateNormalCart(r *NormalCart) ICart
		// 获取购物车
		GetNormalCart(id int32) ICart

		// 获取购物车
		GetShoppingCartByKey(key string) ICart
		// 获取购物车
		GetShoppingCart(key string) *NormalCart
		// 获取最新的购物车
		GetLatestCart(buyerId int64) *NormalCart
		// 保存购物车
		SaveShoppingCart(*NormalCart) (int32, error)
		// 移出购物车项
		RemoveCartItem(id int32) error
		// 保存购物车项
		SaveCartItem(*NormalCartItem) (int32, error)
		// 清空购物车项
		EmptyCartItems(cartId int32) error
		// 删除购物车
		DeleteCart(cartId int32) error

		// Select SaleCartItem
		SelectNormalCartItem(where string, v ...interface{}) []*NormalCartItem
		// Save SaleCart
		SaveNormalCart(v *NormalCart) (int, error)
		// Delete SaleCart
		DeleteNormalCart(primary interface{}) error

		// Save WsCart
		SaveWsCart(v *WsCart) (int, error)
		// Delete WsCart
		DeleteWsCart(primary interface{}) error
		// Select WsCartItem
		SelectWsCartItem(where string, v ...interface{}) []*WsCartItem
		// Save WsCartItem
		SaveWsCartItem(v *WsCartItem) (int, error)
		// Batch Delete WsCartItem
		BatchDeleteWsCartItem(where string, v ...interface{}) (int64, error)
	}

	// 购物车商品
	ItemPair struct {
		// 商品编号
		ItemId int64
		// SKU编号
		SkuId int64
		// 卖家编号
		SellerId int32
		// 数量
		Quantity int32
		// 是否勾选结算
		Checked int32
	}

	// 购物车
	NormalCart struct {
		Id         int32  `db:"id" pk:"yes" auto:"yes"`
		CartCode   string `db:"code"`
		BuyerId    int64  `db:"buyer_id"`
		PaymentOpt int32  `db:"payment_opt"`
		//todo: del???
		DeliverId  int64             `db:"deliver_id"`
		CreateTime int64             `db:"create_time"`
		UpdateTime int64             `db:"update_time"`
		Items      []*NormalCartItem `db:"-"`
	}

	// 购物车项
	NormalCartItem struct {
		// 编号
		Id int32 `db:"id" pk:"yes" auto:"yes"`
		// 购物车编号
		CartId int32 `db:"cart_id"`
		// 运营商编号
		VendorId int32 `db:"vendor_id"`
		// 店铺编号
		ShopId int32 `db:"shop_id"`
		// 商品编号
		ItemId int64 `db:"item_id"`
		// SKU编号
		SkuId int64 `db:"sku_id"`
		// 数量
		Quantity int32 `db:"quantity"`
		// 是否勾选结算
		Checked int32 `db:"checked"`
		// 订单依赖的SKU媒介
		Sku *item.SkuMedia `db:"-"`
	}

	// 商品批发购物车
	WsCart struct {
		// 编号
		ID int32 `db:"id" pk:"yes" auto:"yes"`
		// 购物车编码
		Code string `db:"code"`
		// 买家编号
		BuyerId int64 `db:"buyer_id"`
		// 送货地址
		DeliverId int64 `db:"deliver_id"`
		// 创建时间
		CreateTime int64 `db:"create_time"`
		// 修改时间
		UpdateTime int64 `db:"update_time"`
		// 商品项
		Items []*WsCartItem `db:"-"`
	}

	// 批发购物车商品项
	WsCartItem struct {
		// 编号
		ID int32 `db:"id" pk:"yes" auto:"yes"`
		// 购物车编号
		CartId int32 `db:"cart_id"`
		// 运营商编号
		SellerId int32 `db:"vendor_id"`
		// 店铺编号
		ShopId int32 `db:"shop_id"`
		// 商品编号
		ItemId int64 `db:"item_id"`
		// SKU编号
		SkuId int64 `db:"sku_id"`
		// 数量
		Quantity int32 `db:"quantity"`
		// 订单依赖的SKU媒介
		Sku *item.SkuMedia `db:"-"`
	}

	// 批发购物车JSON数据对象
	WCartJdo struct {
		Seller []WCartSellerJdo
		Data   map[string]string
	}

	// 批发购物车卖家JSON数据对象
	WCartSellerJdo struct {
		// 运营商编号
		SellerId int32
		// 购物车商品
		Item []WCartItemJdo
		// 其他数据
		Data map[string]string
	}

	// 批发购物车商品JSON数据对象
	WCartItemJdo struct {
		// 商品编号
		ItemId int64
		// 商品标题
		ItemName string
		// 商品图片
		ItemImage string
		// SKU列表
		Sku []WCartSkuJdo
		// 其他数据
		Data map[string]string
	}

	// 批发购物车规格JSON数据对象
	WCartSkuJdo struct {
		// SKU编号
		SkuId int64
		// SKU编码
		SkuCode string
		// SKU图片
		SkuImage string
		// 规格文本
		SpecWord string
		// 数量
		Quantity int32
		// 价格
		Price float64
		// 折扣价
		DiscountPrice float64
		// 可售数量
		CanSalesQuantity int32
		// 数据JSON表示
		JData string
	}
)

func ParseCartItem(item *NormalCartItem) *ttype.SShoppingCartItem {
	i := &ttype.SShoppingCartItem{
		ItemId:   item.ItemId,
		SkuId:    item.SkuId,
		Quantity: item.Quantity,
		Checked:  item.Checked == 1,
		ShopId:   item.ShopId,
	}
	if item.Sku != nil {
		i.Image = format.GetGoodsImageUrl(item.Sku.Image)
		i.RetailPrice = math.Round(float64(item.Sku.RetailPrice), 2)
		i.Price = math.Round(float64(item.Sku.Price), 2)
		i.SpecWord = item.Sku.SpecWord
		if i.Title == "" {
			i.Title = item.Sku.Title
		}
		i.Code = item.Sku.ItemCode
		i.StockText = util.BoolExt.TString(item.Sku.Stock > 0,
			"有货", "无货")
	}
	return i
}

func ParseToDtoCart(c ICart) *ttype.SShoppingCart {
	cart := &ttype.SShoppingCart{}
	if c.Kind() != KNormal {
		panic("购物车类型非零售")
	}
	rc := c.(INormalCart)
	v := rc.Value()

	cart.CartId = c.GetAggregateRootId()
	cart.Code = v.CartCode
	cart.Shops = []*ttype.SShoppingCartGroup{}

	items := rc.Items()
	if items != nil && len(items) > 0 {
		mp := make(map[int32]*ttype.SShoppingCartGroup, 0) //保存运营商到map
		for _, v := range items {
			vendor, ok := mp[v.ShopId]
			if !ok {
				vendor = &ttype.SShoppingCartGroup{
					VendorId: v.VendorId,
					ShopId:   v.ShopId,
					Items:    []*ttype.SShoppingCartItem{},
				}
				mp[v.ShopId] = vendor
				cart.Shops = append(cart.Shops, vendor)
			}
			if v.Checked == 1 {
				vendor.Checked = true
			}
			vendor.Items = append(vendor.Items, ParseCartItem(v))
			//cart.TotalNum += v.Quantity
		}
	}

	return cart
}
