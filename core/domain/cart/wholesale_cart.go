package cart

import (
	"go2o/core/domain/interface/cart"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/merchant/shop"
)

var _ cart.ICart = new(wholesaleCartImpl)
var _ cart.IWholesaleCart = new(wholesaleCartImpl)

type wholesaleCartImpl struct {
	value      *cart.WsCart
	rep        cart.ICartRepo
	goodsRepo  item.IGoodsItemRepo
	memberRepo member.IMemberRepo
	summary    string
	shop       shop.IShop
	deliver    member.IDeliverAddress
	snapMap    map[int32]*item.Snapshot
}

func CreateWholesaleCart(val *cart.WsCart, rep cart.ICartRepo,
	memberRepo member.IMemberRepo, goodsRepo item.IGoodsItemRepo) cart.ICart {
	return (&wholesaleCartImpl{
		value:      val,
		rep:        rep,
		memberRepo: memberRepo,
		goodsRepo:  goodsRepo,
	}).init()
}

func (c *wholesaleCartImpl) init() cart.ICart {
	// 初始化购物车的信息
	if c.value != nil && c.value.Items != nil {
		c.setAttachGoodsInfo(c.value.Items)
	}
	return c
}

// 购物车种类
func (c *wholesaleCartImpl) Kind() cart.CartKind {
	return cart.KWholesale
}

// 获取买家编号
func (c *wholesaleCartImpl) BuyerId() int32 {
	return c.value.BuyerId
}

// 检查购物车(仅结算商品)
func (c *wholesaleCartImpl) Check() error {
	if c.value == nil || len(c.value.Items) == 0 {
		return cart.ErrEmptyShoppingCart
	}
	for _, v := range c.value.Items {
		if v.Checked == 1 {
			it := c.goodsRepo.GetItem(v.ItemId)
			if it == nil {
				return item.ErrNoSuchGoods // 没有商品
			}
			stock := it.GetValue().StockNum
			if v.SkuId > 0 {
				if sku := it.GetSku(v.SkuId); sku != nil {
					stock = sku.Stock
				}
			}
			if stock == 0 {
				return item.ErrFullOfStock // 已经卖完了
			}
			if stock < v.Quantity {
				return item.ErrOutOfStock // 超出库存
			}
		}
	}
	return nil
}

// 获取商品的快照列表
func (c *wholesaleCartImpl) getSnapshotsMap(items []*cart.WsCartItem) map[int32]*item.Snapshot {
	if c.snapMap == nil && items != nil {
		l := len(items)
		c.snapMap = make(map[int32]*item.Snapshot, l)
		if l > 0 {
			var ids []int32 = make([]int32, l)
			for i, v := range items {
				ids[i] = v.ItemId
			}
			snapList := c.goodsRepo.GetSnapshots(ids)
			for _, v := range snapList {
				v2 := v
				c.snapMap[v.ItemId] = &v2
			}
		}
	}
	return c.snapMap
}

func (c *wholesaleCartImpl) getBuyerLevelId() int32 {
	if c.value.BuyerId > 0 {
		m := c.memberRepo.GetMember(c.value.BuyerId)
		if m != nil {
			return m.GetValue().Level
		}
	}
	return 0
}

func (c *wholesaleCartImpl) setItemInfo(snap *item.GoodsItem, level int32) {
	// 设置会员价
	if level > 0 {
		gds := c.goodsRepo.CreateItem(snap)
		snap.Price = gds.GetPromotionPrice(level)
	}
}

// 设置附加的商品信息
func (c *wholesaleCartImpl) setAttachGoodsInfo(items []*cart.WsCartItem) {
	list := c.getSnapshotsMap(items)
	if list == nil {
		return
	}
	var sku *item.Sku
	for _, v := range items {
		it := c.goodsRepo.GetItem(v.ItemId)
		if it == nil {
			continue
		}
		if v.SkuId > 0 {
			sku = it.GetSku(v.SkuId)
		}
		v.Sku = item.ParseSkuMedia(it.GetValue(), sku)

		//  会员价
		//var level int32
		//if gv.LevelSales == 1 && level != -1 {
		//    if level == 0 {
		//        level = c.getBuyerLevelId()
		//    }
		//    //todo: ???
		//    //c.setGoodsInfo(gv, level)
		//}
		//
	}
}

// 获取聚合根编号
func (c *wholesaleCartImpl) GetAggregateRootId() int32 {
	return c.value.ID
}

func (c *wholesaleCartImpl) GetValue() cart.WsCart {
	return *c.value
}

// 获取购物车中的商品
func (c *wholesaleCartImpl) GetCartGoods() []item.IGoodsItem {
	//todo: IMPL
	//var gs []item.IGoods = make([]item.IGoods, len(c._value.Items))
	//for i, v := range c._value.Items {
	//    gs[i] = c._goodsRepo.getGoods
	//}
	//return gs
	return []item.IGoodsItem{}
}

// 获取商品编号与购物车项的集合
func (c *wholesaleCartImpl) Items() map[int32]*cart.WsCartItem {
	list := make(map[int32]*cart.WsCartItem)
	for _, v := range c.value.Items {
		list[v.SkuId] = v
	}
	return list
}

func (c *wholesaleCartImpl) getItems() []*cart.WsCartItem {
	return c.value.Items
}

// 添加项
func (c *wholesaleCartImpl) Put(itemId, skuId int32, num int32) error {
	_, err := c.put(itemId, skuId, num)
	return err
}

// 添加项
func (c *wholesaleCartImpl) put(itemId, skuId int32, num int32) (*cart.WsCartItem, error) {
	var err error
	if c.value.Items == nil {
		c.value.Items = nil //[]*cart.WsCartItem{}
	}

	var sku *item.Sku
	it := c.goodsRepo.GetItem(itemId)
	if it == nil {
		return nil, item.ErrNoSuchGoods // 没有商品
	}
	iv := it.GetValue()
	// 库存,如有SKU，则使用SKU的库存
	stock := iv.StockNum
	// 判断是否上架

	if iv.ShelveState != item.ShelvesOn {
		return nil, item.ErrNotOnShelves //未上架
	}
	// 判断商品SkuId
	if skuId > 0 {
		sku = it.GetSku(skuId)
		if sku == nil {
			return nil, item.ErrNoSuchItemSku
		}
		stock = sku.Stock
	} else if iv.SkuNum > 0 {
		return nil, cart.ErrItemNoSku
	}
	// 检查是否已经卖完了
	if stock == 0 {
		return nil, item.ErrFullOfStock
	}

	// 添加数量
	for _, v := range c.value.Items {
		if v.ItemId == itemId && v.SkuId == skuId {
			if v.Quantity+num > stock {
				return v, item.ErrOutOfStock // 库存不足
			}
			v.Quantity += num
			return v, err
		}
	}

	c.snapMap = nil

	// 设置商品的相关信息
	c.setItemInfo(iv, c.getBuyerLevelId())

	v := &cart.WsCartItem{
		CartId:   c.GetAggregateRootId(),
		VendorId: iv.VendorId,
		ShopId:   iv.ShopId,
		ItemId:   iv.Id,
		SkuId:    skuId,
		Quantity: num,
		Sku:      item.ParseSkuMedia(iv, sku),
		Checked:  1,
	}
	c.value.Items = append(c.value.Items, v)
	return v, err
}

// 移出项
func (c *wholesaleCartImpl) Remove(itemId, skuId, num int32) error {
	if c.value.Items == nil {
		return cart.ErrEmptyShoppingCart
	}
	// 删除数量
	for _, v := range c.value.Items {
		if v.ItemId == itemId && v.SkuId == skuId {
			if newNum := v.Quantity - num; newNum <= 0 {
				v.Quantity = 0
			} else {
				v.Quantity = newNum
			}
			break
		}
	}
	c.snapMap = nil //clean

	return nil
}

// 获取购物车的KEY
func (c *wholesaleCartImpl) Key() string {
	return c.value.Code
}

/*
func (c *wholesaleCartImpl) combineBuyerCart() cart.ICart {

    var hasOutCart = len(cartCode) != 0
    var hasBuyer = c._value.BuyerId > 0

    var memCart cart.ICart = nil // 消费者的购物车
    var outCart cart.ICart = c // 当前购物车

    if hasBuyer {
        // 如果没有传递cartCode ，或者传递的cart和会员绑定的购物车相同，直接返回
        if memCart = c._rep.GetMemberCurrentCart(c._value.BuyerId);
            memCart != nil {
            if memCart.Key() == outCart.Key() {
                return memCart
            }
        } else {
            memCart = c.NewCart()
        }
    }

    if hasOutCart {
        outCart, _ = c.GetCartByKey(cartCode)
    }

    // 合并购物车
    if outCart != nil && hasBuyer {
        if buyerId := outCart.GetValue().BuyerId; buyerId <= 0 || buyerId == c._buyerId {
            memCart, _ = memCart.Combine(outCart)
            outCart.Destroy()
            memCart.Save()
        }
    }

    if memCart != nil {
        return memCart
    }

    if outCart != nil {
        return outCart
    }

    return c.NewCart()

    //	if !hasOutCart {
    //		if c == nil {
    //			// 新的购物车不存在，直接返回会员的购物车
    //			if mc != nil {
    //				return mc
    //			}
    //		} else {
    //			cv := c.GetValue()
    //			//合并购物车
    //			if cv.BuyerId <= 0 {
    //				// 设置购买者
    //				if hasBuyer {
    //					c.SetBuyer(buyerId)
    //				}
    //			} else if mc != nil && cv.BuyerId == buyerId {
    //				// 合并购物车
    //				nc, err := mc.Combine(c)
    //				if err == nil {
    //					nc.Save()
    //					return nc
    //				}
    //				return mc
    //			}
    //
    //			// 如果没有购买，则返回
    //			return c
    //		}
    //	}

    // 返回一个新的购物车
    //	return c.NewCart(buyerId)
}
*/

// 设置购买会员
func (c *wholesaleCartImpl) SetBuyer(buyerId int32) error {
	if c.value.BuyerId > 0 {
		return cart.ErrCartBuyerBind
	}
	c.value.BuyerId = buyerId
	_, err := c.Save()
	return err
}

// 设置购买会员收货地址
func (c *wholesaleCartImpl) SetBuyerAddress(addressId int32) error {
	if c.value.BuyerId < 0 {
		return cart.ErrCartNoBuyer
	}
	m := c.memberRepo.GetMember(c.value.BuyerId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	addr := m.Profile().GetAddress(addressId)
	if addr == nil {
		return member.ErrNoSuchAddress
	}
	return c.setBuyerAddress(addressId)
}

func (c *wholesaleCartImpl) setBuyerAddress(addressId int32) error {
	c.value.DeliverId = addressId
	_, err := c.Save()
	return err
}

// 标记商品结算
func (c *wholesaleCartImpl) SignItemChecked(items []*cart.CartItem) error {
	mp := c.getItems()
	for _, item := range mp {
		item.Checked = 0
		for _, v := range items {
			if v.SkuId == item.SkuId && v.ItemId == item.ItemId {
				item.Checked = 1
				break
			}
		}
	}
	err := c.Check()
	if err == nil {
		_, err = c.Save()
	}
	return err
}

// 结算数据持久化
func (c *wholesaleCartImpl) SettlePersist(shopId, paymentOpt, deliverOpt, addressId int32) error {
	//var shop shop.IShop
	var deliver member.IDeliverAddress
	var err error

	if shopId > 0 {
		//var mch merchant.IMerchant
		//mch, err = c._partnerRepo.GetMerchant(c._mchId)
		//if err != nil {
		//	return err
		//}
		//shop = mch.ShopManager().GetShop(shopId)
		//if shop == nil {
		//	return merchant.ErrNoSuchShop
		//}
		//c._shop = shop
		//c._value.ShopId = shopId

		//todo: not implement
		return err
	}

	if c.value.BuyerId > 0 && addressId > 0 {
		m := c.memberRepo.GetMember(c.value.BuyerId)
		if m == nil {
			return member.ErrNoSuchMember
		}
		deliver = m.Profile().GetAddress(addressId)
		if deliver == nil {
			return member.ErrInvalidSession
		}
		c.deliver = deliver
		c.value.DeliverId = addressId
	}

	return nil
}

// 获取结算数据
func (c *wholesaleCartImpl) GetSettleData() (s shop.IShop, d member.IDeliverAddress,
	paymentOpt int32) {

	if c.deliver == nil {
		pm := c.memberRepo.GetMember(c.value.BuyerId).Profile()
		if c.value.DeliverId > 0 {
			c.deliver = pm.GetAddress(c.value.DeliverId)
		} else {
			c.deliver = pm.GetDefaultAddress()
			if c.deliver != nil {
				c.setBuyerAddress(c.deliver.GetDomainId())
			}
		}
	}
	return c.shop, c.deliver, -1
}

// 保存购物车
func (c *wholesaleCartImpl) Save() (int32, error) {
	panic("not impl")
	//c.value.UpdateTime = time.Now().Unix()
	//id, err := c.rep.SaveShoppingCart(c.value)
	//c.value.ID = id
	//if c.value.Items != nil {
	//    for _, v := range c.value.Items {
	//        if v.Quantity <= 0 {
	//            c.rep.RemoveCartItem(v.ID)
	//        } else {
	//            v.ID, err = c.rep.SaveCartItem(v)
	//        }
	//    }
	//}
	//return id, err
}

// 释放购物车,如果购物车的商品全部结算,则返回true
func (c *wholesaleCartImpl) Release() bool {
	checked := []int{}
	for i, v := range c.value.Items {
		if v.Checked == 1 {
			checked = append(checked, i)
		}
	}
	// 如果为部分结算,则移除商品并返回false
	if len(checked) < len(c.value.Items) {
		for _, i := range checked {
			v := c.value.Items[i]
			c.Remove(v.ItemId, v.SkuId, v.Quantity)
		}
		c.Save()
		return false
	}
	return true
}

// 销毁购物车
func (c *wholesaleCartImpl) Destroy() (err error) {
	c.snapMap = nil //clean
	if err = c.rep.EmptyCartItems(c.GetAggregateRootId()); err == nil {
		return c.rep.DeleteCart(c.GetAggregateRootId())
	}
	return err
}

// 获取订单金额,返回totalFee为总额，
// orderFee为实际订单的金额(扣去促销优惠等后的金额)
func (c *wholesaleCartImpl) GetFee() (totalFee float32, orderFee float32) {
	var qua float32
	for _, v := range c.value.Items {
		if v.Checked == 1 {
			qua = float32(v.Quantity)
			totalFee += v.Sku.RetailPrice * qua
			orderFee += v.Sku.Price * qua
		}
	}
	return totalFee, orderFee
}
