package cart

import (
	"bytes"
	"encoding/json"
	"go2o/core/domain/interface/cart"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/domain/interface/order"
	"go2o/core/domain/interface/sale"
	"go2o/core/domain/interface/sale/goods"
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
	_snapMap   map[int]*goods.Snapshot
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

// 检查购物车(仅结算商品)
func (this *cartImpl) Check() error {
	if this._value == nil || len(this._value.Items) == 0 {
		return cart.ErrEmptyShoppingCart
	}
	for _, v := range this._value.Items {
		if v.Checked == 1 {
			snap := this._goodsRep.GetLatestSnapshot(v.SkuId)
			if snap == nil {
				return goods.ErrNoSuchGoods // 没有商品
			}
			if snap.StockNum == 0 {
				return sale.ErrFullOfStock // 已经卖完了
			}
			if snap.StockNum < v.Quantity {
				return sale.ErrOutOfStock // 超出库存
			}
		}
	}
	return nil
}

// 获取商品的快招列表
func (this *cartImpl) getSnapshotsMap(items []*cart.CartItem) map[int]*goods.Snapshot {
	if this._snapMap == nil {
		if items != nil {
			l := len(items)
			this._snapMap = make(map[int]*goods.Snapshot, l)
			if l > 0 {
				var ids []int = make([]int, l)
				for i, v := range items {
					ids[i] = v.SkuId
				}
				snapList := this._goodsRep.GetSnapshots(ids)
				for _, v := range snapList {
					v2 := v
					this._snapMap[v.SkuId] = &v2
				}
			}
		}
	}
	return this._snapMap
}

func (this *cartImpl) getBuyerLevelId() int {
	if this._value.BuyerId > 0 {
		m := this._memberRep.GetMember(this._value.BuyerId)
		if m != nil {
			return m.GetValue().Level
		}
	}
	return -1
}

func (this *cartImpl) setGoodsInfo(snap *goods.Snapshot, level int) {
	// 设置会员价
	if level > 0 {
		gds := this._goodsRep.GetGoodsBySKuId(snap.SkuId).(sale.IGoods)
		snap.SalePrice = gds.GetPromotionPrice(level)
	}
}

// 设置附加的商品信息
func (this *cartImpl) setAttachGoodsInfo(items []*cart.CartItem) {
	list := this.getSnapshotsMap(items)
	if list == nil {
		return
	}
	var level int
	for _, v := range items {
		gv, ok := list[v.SkuId]
		//  会员价
		if gv.LevelSales == 1 && level != -1 {
			if level == 0 {
				level = this.getBuyerLevelId()
			}
			this.setGoodsInfo(gv, level)
		}
		// 设置购物车项的数据
		if ok {
			v.Name = gv.GoodsTitle
			v.Price = gv.Price
			v.GoodsNo = gv.GoodsNo
			v.Image = gv.Image
			v.SalePrice = gv.SalePrice
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
	//todo: IMPL
	//var gs []sale.IGoods = make([]sale.IGoods, len(this._value.Items))
	//for i, v := range this._value.Items {
	//    gs[i] = this._goodsRep.getGoods
	//}
	//return gs
	return []sale.IGoods{}
}

// 获取商品编号与购物车项的集合
func (this *cartImpl) Items() map[int]*cart.CartItem {
	list := make(map[int]*cart.CartItem)
	for _, v := range this._value.Items {
		list[v.SkuId] = v
	}
	return list
}

// 添加项
func (this *cartImpl) AddItem(vendorId int, shopId int, skuId int,
	num int) (*cart.CartItem, error) {
	var err error
	if this._value.Items == nil {
		this._value.Items = []*cart.CartItem{}
	}
	snap := this._goodsRep.GetLatestSnapshot(skuId)
	if snap == nil {
		return nil, goods.ErrNoSuchGoods // 没有商品
	}
	if snap.OnShelves != 1 {
		return nil, goods.ErrNotOnShelves //未上架
	}
	if snap.StockNum == 0 {
		return nil, sale.ErrFullOfStock // 已经卖完了
	}
	// 添加数量
	for _, v := range this._value.Items {
		if v.SkuId == skuId {
			if v.Quantity+num > snap.StockNum {
				return v, sale.ErrOutOfStock // 库存不足
			}
			v.Quantity += num
			return v, err
		}
	}

	this._snapMap = nil //clean

	// 设置商品的相关信息
	this.setGoodsInfo(snap, this.getBuyerLevelId())

	v := &cart.CartItem{
		CartId:     this.GetAggregateRootId(),
		VendorId:   vendorId,
		ShopId:     shopId,
		SnapshotId: snap.SkuId,
		SkuId:      skuId,
		Quantity:   num,
		Name:       snap.GoodsTitle,
		GoodsNo:    snap.GoodsNo,
		Image:      snap.Image,
		Price:      snap.Price,
		SalePrice:  snap.SalePrice,
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
		if v.SkuId == goodsId {
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

	this._snapMap = nil //clean

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
func (this *cartImpl) Combine(c cart.ICart) cart.ICart {
	if c.GetAggregateRootId() != this.GetAggregateRootId() {
		for _, v := range c.GetValue().Items {
			if item, err := this.AddItem(v.VendorId, v.ShopId,
				v.SkuId, v.Quantity); err == nil {
				if v.Checked == 1 {
					item.Checked = 1
				}
			}
		}
		c.Destroy() //合并后,需销毁购物车
	}
	this._snapMap = nil //clean
	return this
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

// 标记商品结算
func (this *cartImpl) SignItemChecked(skuArr []int) error {
	mp := this.Items()
	arrMap := make(map[int]int, len(skuArr))
	for _, v := range skuArr {
		arrMap[v] = 0
	}
	for skuId, item := range mp {
		if _, ok := arrMap[skuId]; ok {
			item.Checked = 1
		} else {
			item.Checked = 0
		}
	}
	err := this.Check()
	if err == nil {
		_, err = this.Save()
	}
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
func (this *cartImpl) GetSettleData() (s shop.IShop, d member.IDeliverAddress,
	paymentOpt, deliverOpt int) {
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
				v.Id, err = this._rep.SaveCartItem(v)
			}
		}
	}
	return id, err
}

// 释放购物车,如果购物车的商品全部结算,则返回true
func (this *cartImpl) Release() bool {
	checked := []int{}
	for i, v := range this._value.Items {
		if v.Checked == 1 {
			checked = append(checked, i)
		}
	}
	// 如果为部分结算,则移除商品并返回false
	if len(checked) < len(this._value.Items) {
		for _, i := range checked {
			v := this._value.Items[i]
			this.RemoveItem(v.SkuId, v.Quantity)
		}
		this.Save()
		return false
	}
	return true
}

// 销毁购物车
func (this *cartImpl) Destroy() (err error) {
	this._snapMap = nil //clean
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

	list := this.getSnapshotsMap(this._value.Items)
	if list != nil {
		length := len(list)
		for i, v := range this._value.Items {
			snap := list[v.SkuId]
			if snap != nil {
				buf.WriteString(snap.GoodsTitle)
				if len(snap.SmallTitle) != 0 {
					buf.WriteString("(" + snap.SmallTitle + ")")
				}
				buf.WriteString("*" + strconv.Itoa(v.Quantity))
				if i < length-1 {
					buf.WriteString("\n")
				}
			}
		}
	}
	return buf.String()
}

// 获取Json格式的商品数据
func (this *cartImpl) GetJsonItems() []byte {
	var goods []*order.OrderGoods = make([]*order.OrderGoods, len(this._value.Items))
	for i, v := range this._value.Items {
		goods[i] = &order.OrderGoods{
			GoodsId:    v.SkuId,
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
		if v.Checked == 1 {
			qua = float32(v.Quantity)
			totalFee += v.Price * qua
			orderFee += v.SalePrice * qua
		}
	}
	return totalFee, orderFee
}
