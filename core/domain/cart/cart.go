package cart

import (
	"bytes"
	"encoding/json"
	"go2o/core/domain/interface/cart"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/domain/interface/sale"
	"go2o/core/domain/interface/sale/goods"
	"go2o/core/domain/interface/sale/item"
	"go2o/core/domain/interface/shopping"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/infrastructure/domain"
	"strconv"
	"time"
)

type cartImpl struct {
	_value     *cart.ValueCart
	_rep       cart.ICartRep
	_goodsRep  goods.IGoodsRep
	_memberRep member.IMemberRep
	_summary   string
	_shop      shop.IShop
	_deliver   member.IDeliverAddress
}

func CreateCart(val *cart.ValueCart, rep cart.ICartRep,
	memberRep member.IMemberRep, goodsRep goods.IGoodsRep) cart.ICart {
	return (&cartImpl{
		_value:     val,
		_rep:       rep,
		_memberRep: memberRep,
		_goodsRep:  goodsRep,
	}).init()
}

//todo: merchantId 应去掉，可能在多个商户买东西
// 创建新的购物车
func NewCart(buyerId int, rep cart.ICartRep, memberRep member.IMemberRep,
	goodsRep goods.IGoodsRep) cart.ICart {
	unix := time.Now().Unix()
	cartKey := domain.GenerateCartKey(unix, time.Now().Nanosecond())
	value := &cart.ValueCart{
		CartKey:    cartKey,
		BuyerId:    buyerId,
		ShopId:     0,
		DeliverId:  0,
		PaymentOpt: 1,
		DeliverOpt: 1,
		CreateTime: unix,
		UpdateTime: unix,
		Items:      nil,
	}
	return CreateCart(value, rep, memberRep, goodsRep)
}

func (this *cartImpl) init() cart.ICart {
	// 初始化购物车的信息
	if this._value != nil && this._value.Items != nil {
		this.setAttachGoodsInfo(this._value.Items)
	}
	return this
}

// 检查购物车
func (this *cartImpl) Check() error {
	if this._value == nil || len(this._value.Items) == 0 {
		return cart.ErrEmptyShoppingCart
	}

	for _, v := range this._value.Items {
		snap := this._goodsRep.GetLatestSnapshot(v.GoodsId)
		if snap == nil {
			return goods.ErrNoSuchGoods // 没有商品
		}
		gs := this._goodsRep.GetValueGoodsById(snap.SkuId)
		stockNum := gs.StockNum
		if stockNum == 0 {
			return sale.ErrFullOfStock // 已经卖完了
		}
		if stockNum < v.Quantity {
			return sale.ErrOutOfStock // 超出库存
		}
	}
	return nil
}

// 设置附加的商品信息
func (this *cartImpl) setAttachGoodsInfo(items []*cart.CartItem) {
	if items != nil {
		l := len(items)
		if l == 0 {
			return
		}
		var ids []int = make([]int, l)
		for i, v := range items {
			ids[i] = v.GoodsId
		}

		// 设置附加的值
		goodsList, err := this._goodsRep.GetGoodsByIds(ids...)
		if err == nil {
			var goodsMap = make(map[int]*valueobject.Goods, len(goodsList))
			for _, v := range goodsList {
				goodsMap[v.GoodsId] = v
			}

			var level int
			var gds sale.IGoods
			var sl sale.ISale

			//  更新登陆后的优惠价
			if this._value.BuyerId > 0 {

				//todo: impl
				/*
				   sl = this._saleRep.GetSale(this._merchantId)
				   m := this._memberRep.GetMember(this._value.BuyerId)
				   if m != nil {
				       level = m.GetValue().Level
				   }*/
			}

			for _, v := range items {
				gv, ok := goodsMap[v.GoodsId]
				if level > 0 {
					gds = sl.GoodsManager().CreateGoodsByItem(
						sl.ItemManager().CreateItem(item.ParseToPartialValueItem(gv)),
						goods.ParseToValueGoods(gv),
					)
					if p := gds.GetPromotionPrice(level); p < gv.SalePrice {
						gv.SalePrice = p
					}
				}
				if ok {
					v.Name = gv.Name
					v.Price = gv.Price
					v.GoodsNo = gv.GoodsNo
					v.Image = gv.Image
					v.SalePrice = gv.SalePrice
				}
			}
		}
	}
}

// 获取聚合根编号
func (this *cartImpl) GetAggregateRootId() int {
	return this._value.Id
}

func (this *cartImpl) GetValue() cart.ValueCart {
	return *this._value
}

// 获取购物车中的商品
func (this *cartImpl) GetCartGoods() []sale.IGoods {
	//todo: not implement
	/*
	   sl := this._saleRep.GetSale(this._merchantId)
	   var gs []sale.IGoods = make([]sale.IGoods, len(this._value.Items))
	   for i, v := range this._value.Items {
	       gs[i] = sl.GoodsManager().GetGoods(v.GoodsId)
	   }
	   return gs
	*/
	return []sale.IGoods{}
}

// 获取商品编号与购物车项的集合
func (this *cartImpl) Items() map[int]*cart.CartItem {
	list := make(map[int]*cart.CartItem)
	for _, v := range this._value.Items {
		list[v.GoodsId] = v
	}
	return list
}

// 添加项
func (this *cartImpl) AddItem(mchId int, shopId int, goodsId int,
	num int) (*cart.CartItem, error) {
	var err error
	if this._value.Items == nil {
		this._value.Items = []*cart.CartItem{}
	}
	snap := this._goodsRep.GetLatestSnapshot(goodsId)
	if snap == nil {
		return nil, goods.ErrNoSuchGoods // 没有商品
	}

	//todo: 是否上架, 下架后应将当前销售快照删除,或快照包含是否上架的信息

	//if !gds.GetItem().IsOnShelves() {
	//    return nil, goods.ErrNotOnShelves //未上架
	//}
	gs := this._goodsRep.GetValueGoodsById(snap.SkuId)
	stockNum := gs.StockNum
	if stockNum == 0 {
		return nil, sale.ErrFullOfStock // 已经卖完了
	}

	// 添加数量
	for _, v := range this._value.Items {
		if v.GoodsId == goodsId {
			if v.Quantity+num > stockNum {
				return v, sale.ErrOutOfStock // 库存不足
			}
			v.Quantity += num
			return v, err
		}
	}

	v := &cart.CartItem{
		CartId:     this.GetAggregateRootId(),
		VendorId:   mchId,
		ShopId:     shopId,
		SnapshotId: snap.SkuId,
		GoodsId:    goodsId,
		Quantity:   num,
		Name:       snap.GoodsTitle,
		GoodsNo:    snap.GoodsNo,
		Image:      snap.Image,
		Price:      snap.Price,
		//todo: 优惠价
		// SalePrice:  gv.PromPrice, // 使用优惠价
	}
	this._value.Items = append(this._value.Items, v)
	return v, err
}

// 移出项
func (this *cartImpl) RemoveItem(goodsId, num int) error {
	if this._value.Items == nil {
		return cart.ErrEmptyShoppingCart
	}

	// 删除数量
	for _, v := range this._value.Items {
		if v.GoodsId == goodsId {
			if newNum := v.Quantity - num; newNum <= 0 {
				// 移出购物车
				//this.value.Items = append(this.value.Items[:i],this.value.Items[i+1:]...)
				v.Quantity = 0
			} else {
				v.Quantity = newNum
			}
			break
		}
	}
	return nil
}

// 获取购物车的KEY
func (this *cartImpl) Key() string {
	return this._value.CartKey
}

/*
func (this *cartImpl) combineBuyerCart() cart.ICart {

    var hasOutCart = len(cartKey) != 0
    var hasBuyer = this._value.BuyerId > 0

    var memCart cart.ICart = nil // 消费者的购物车
    var outCart cart.ICart = this // 当前购物车

    if hasBuyer {
        // 如果没有传递cartKey ，或者传递的cart和会员绑定的购物车相同，直接返回
        if memCart = this._rep.GetMemberCurrentCart(this._value.BuyerId);
            memCart != nil {
            if memCart.Key() == outCart.Key() {
                return memCart
            }
        } else {
            memCart = this.NewCart()
        }
    }

    if hasOutCart {
        outCart, _ = this.GetCartByKey(cartKey)
    }

    // 合并购物车
    if outCart != nil && hasBuyer {
        if buyerId := outCart.GetValue().BuyerId; buyerId <= 0 || buyerId == this._buyerId {
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

    return this.NewCart()

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
    //	return this.NewCart(buyerId)
}
*/

// 合并购物车，并返回新的购物车
func (this *cartImpl) Combine(c cart.ICart) (cart.ICart, error) {
	if c.GetAggregateRootId() != this.GetAggregateRootId() {
		for _, v := range c.GetValue().Items {
			this.AddItem(v.VendorId, v.ShopId, v.GoodsId, v.Quantity)
		}
	}
	return this, nil
}

// 设置购买会员
func (this *cartImpl) SetBuyer(buyerId int) error {
	if this._value.BuyerId > 0 {
		return cart.ErrCartBuyerBinded
	}
	this._value.BuyerId = buyerId
	memCart := this._rep.GetMemberCurrentCart(buyerId)
	if memCart != nil && memCart.Key() != this.Key() {
		this.Combine(memCart)
	}
	_, err := this.Save()
	return err
}

// 结算数据持久化
func (this *cartImpl) SettlePersist(shopId, paymentOpt, deliverOpt, deliverId int) error {
	//var shop shop.IShop
	var deliver member.IDeliverAddress
	var err error

	if shopId > 0 {
		//var mch merchant.IMerchant
		//mch, err = this._partnerRep.GetMerchant(this._merchantId)
		//if err != nil {
		//	return err
		//}
		//shop = mch.ShopManager().GetShop(shopId)
		//if shop == nil {
		//	return merchant.ErrNoSuchShop
		//}
		//this._shop = shop
		//this._value.ShopId = shopId

		//todo: not implement
		return err
	}

	if this._value.BuyerId > 0 && deliverId > 0 {
		var m member.IMember = this._memberRep.GetMember(this._value.BuyerId)
		if m == nil {
			return member.ErrNoSuchMember
		}
		deliver = m.Profile().GetDeliver(deliverId)
		if deliver == nil {
			return member.ErrInvalidSession
		}
		this._deliver = deliver
		this._value.DeliverId = deliverId
	}

	this._value.PaymentOpt = paymentOpt
	this._value.DeliverOpt = deliverOpt
	return nil
}

// 获取结算数据
func (this *cartImpl) GetSettleData() (s shop.IShop, d member.IDeliverAddress, paymentOpt, deliverOpt int) {
	//var err error
	if this._value.ShopId > 0 && this._shop == nil {
		//var pt merchant.IMerchant
		//pt, err = this._partnerRep.GetMerchant(this._merchantId)
		//if err == nil {
		//	this._shop = pt.ShopManager().GetShop(this._value.ShopId)
		//}
		//todo: not implement
	}
	if this._value.DeliverId > 0 && this._deliver == nil {
		var m member.IMember
		m = this._memberRep.GetMember(this._value.BuyerId)
		if m != nil {
			this._deliver = m.Profile().GetDeliver(this._value.DeliverId)
		}
	}
	return this._shop, this._deliver, this._value.PaymentOpt, this._value.DeliverOpt
}

// 保存购物车
func (this *cartImpl) Save() (int, error) {
	this._value.UpdateTime = time.Now().Unix()
	id, err := this._rep.SaveShoppingCart(this._value)
	this._value.Id = id

	if this._value.Items != nil {
		for _, v := range this._value.Items {
			if v.Quantity <= 0 {
				this._rep.RemoveCartItem(v.Id)
			} else {
				i, err := this._rep.SaveCartItem(v)
				if err != nil {
					v.Id = i
				}
			}
		}
	}

	return id, err
}

// 销毁购物车
func (this *cartImpl) Destroy() (err error) {
	if err = this._rep.EmptyCartItems(this.GetAggregateRootId()); err == nil {
		return this._rep.DeleteCart(this.GetAggregateRootId())
	}
	return err
}

// 获取总览信息
func (this *cartImpl) GetSummary() string {
	if len(this._summary) != 0 {
		return this._summary
	}
	buf := bytes.NewBufferString("")
	length := len(this._value.Items)

	var snap *goods.GoodsSnapshot
	for i, v := range this._value.Items {

		snap = this._goodsRep.GetSaleSnapshot(v.SnapshotId)
		if snap != nil {
			buf.WriteString(snap.GoodsName)
			if len(snap.SmallTitle) != 0 {
				buf.WriteString("(" + snap.SmallTitle + ")")
			}
			buf.WriteString("*" + strconv.Itoa(v.Quantity))
			if i < length-1 {
				buf.WriteString("\n")
			}
		}
	}
	return buf.String()
}

// 获取Json格式的商品数据
func (this *cartImpl) GetJsonItems() []byte {
	var goods []*shopping.OrderGoods = make([]*shopping.OrderGoods, len(this._value.Items))
	for i, v := range this._value.Items {
		goods[i] = &shopping.OrderGoods{
			GoodsId:    v.GoodsId,
			GoodsImage: v.Image,
			Quantity:   v.Quantity,
			Name:       v.Name,
		}
	}
	d, _ := json.Marshal(goods)
	return d
}

// 获取订单金额,返回totalFee为总额，
// orderFee为实际订单的金额(扣去促销优惠等后的金额)
func (this *cartImpl) GetFee() (totalFee float32, orderFee float32) {
	var qua float32
	for _, v := range this._value.Items {
		qua = float32(v.Quantity)
		totalFee = totalFee + v.Price*qua
		orderFee = orderFee + v.SalePrice*qua
	}
	return totalFee, orderFee
}
