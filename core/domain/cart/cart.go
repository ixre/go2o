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
	"go2o/core/domain/interface/sale/item"
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

func (c *cartImpl) init() cart.ICart {
	// 初始化购物车的信息
	if c._value != nil && c._value.Items != nil {
		c.setAttachGoodsInfo(c._value.Items)
	}
	return c
}

// 检查购物车(仅结算商品)
func (c *cartImpl) Check() error {
	if c._value == nil || len(c._value.Items) == 0 {
		return cart.ErrEmptyShoppingCart
	}
	for _, v := range c._value.Items {
		if v.Checked == 1 {
			snap := c._goodsRep.GetLatestSnapshot(v.SkuId)
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
func (c *cartImpl) getSnapshotsMap(items []*cart.CartItem) map[int]*goods.Snapshot {
	if c._snapMap == nil {
		if items != nil {
			l := len(items)
			c._snapMap = make(map[int]*goods.Snapshot, l)
			if l > 0 {
				var ids []int = make([]int, l)
				for i, v := range items {
					ids[i] = v.SkuId
				}
				snapList := c._goodsRep.GetSnapshots(ids)
				for _, v := range snapList {
					v2 := v
					c._snapMap[v.SkuId] = &v2
				}
			}
		}
	}
	return c._snapMap
}

func (c *cartImpl) getBuyerLevelId() int {
	if c._value.BuyerId > 0 {
		m := c._memberRep.GetMember(c._value.BuyerId)
		if m != nil {
			return m.GetValue().Level
		}
	}
	return -1
}

func (c *cartImpl) setGoodsInfo(snap *goods.Snapshot, level int) {
	// 设置会员价
	if level > 0 {
		gds := c._goodsRep.GetGoodsBySKuId(snap.SkuId).(sale.IGoods)
		snap.SalePrice = gds.GetPromotionPrice(level)
	}
}

// 设置附加的商品信息
func (c *cartImpl) setAttachGoodsInfo(items []*cart.CartItem) {
	list := c.getSnapshotsMap(items)
	if list == nil {
		return
	}
	var level int
	for _, v := range items {
		gv, ok := list[v.SkuId]
		//  会员价
		if gv.LevelSales == 1 && level != -1 {
			if level == 0 {
				level = c.getBuyerLevelId()
			}
			c.setGoodsInfo(gv, level)
		}
		// 设置购物车项的数据
		if ok {
			v.Snapshot = gv
			v.Name = gv.GoodsTitle
			v.Price = gv.Price
			v.GoodsNo = gv.GoodsNo
			v.Image = gv.Image
			v.SalePrice = gv.SalePrice
		}
	}
}

// 获取聚合根编号
func (c *cartImpl) GetAggregateRootId() int {
	return c._value.Id
}

func (c *cartImpl) GetValue() cart.ValueCart {
	return *c._value
}

// 获取购物车中的商品
func (c *cartImpl) GetCartGoods() []sale.IGoods {
	//todo: IMPL
	//var gs []sale.IGoods = make([]sale.IGoods, len(c._value.Items))
	//for i, v := range c._value.Items {
	//    gs[i] = c._goodsRep.getGoods
	//}
	//return gs
	return []sale.IGoods{}
}

// 获取商品编号与购物车项的集合
func (c *cartImpl) Items() map[int]*cart.CartItem {
	list := make(map[int]*cart.CartItem)
	for _, v := range c._value.Items {
		list[v.SkuId] = v
	}
	return list
}

// 添加项
func (c *cartImpl) AddItem(vendorId int, shopId int, skuId int,
	num int, checked bool) (*cart.CartItem, error) {
	var err error
	if c._value.Items == nil {
		c._value.Items = []*cart.CartItem{}
	}
	snap := c._goodsRep.GetLatestSnapshot(skuId)
	if snap == nil {
		return nil, goods.ErrNoSuchGoods // 没有商品
	}
	if snap.ShelveState != item.ShelvesOn {
		return nil, goods.ErrNotOnShelves //未上架
	}
	if snap.StockNum == 0 {
		return nil, sale.ErrFullOfStock // 已经卖完了
	}
	// 添加数量
	for _, v := range c._value.Items {
		if v.SkuId == skuId {
			if v.Quantity+num > snap.StockNum {
				return v, sale.ErrOutOfStock // 库存不足
			}
			v.Quantity += num
			if checked {
				v.Checked = 1
			}
			return v, err
		}
	}

	c._snapMap = nil //clean

	// 设置商品的相关信息
	c.setGoodsInfo(snap, c.getBuyerLevelId())

	v := &cart.CartItem{
		CartId:     c.GetAggregateRootId(),
		VendorId:   vendorId,
		ShopId:     shopId,
		Snapshot:   snap,
		SnapshotId: snap.SkuId,
		SkuId:      skuId,
		Quantity:   num,
		Name:       snap.GoodsTitle,
		GoodsNo:    snap.GoodsNo,
		Image:      snap.Image,
		Price:      snap.Price,
		SalePrice:  snap.SalePrice,
	}
	if checked {
		v.Checked = 1
	}
	c._value.Items = append(c._value.Items, v)
	return v, err
}

// 移出项
func (c *cartImpl) RemoveItem(goodsId, num int) error {
	if c._value.Items == nil {
		return cart.ErrEmptyShoppingCart
	}

	// 删除数量
	for _, v := range c._value.Items {
		if v.SkuId == goodsId {
			if newNum := v.Quantity - num; newNum <= 0 {
				// 移出购物车
				//c.value.Items = append(c.value.Items[:i],c.value.Items[i+1:]...)
				v.Quantity = 0
			} else {
				v.Quantity = newNum
			}
			break
		}
	}

	c._snapMap = nil //clean

	return nil
}

// 获取购物车的KEY
func (c *cartImpl) Key() string {
	return c._value.CartKey
}

/*
func (c *cartImpl) combineBuyerCart() cart.ICart {

    var hasOutCart = len(cartKey) != 0
    var hasBuyer = c._value.BuyerId > 0

    var memCart cart.ICart = nil // 消费者的购物车
    var outCart cart.ICart = c // 当前购物车

    if hasBuyer {
        // 如果没有传递cartKey ，或者传递的cart和会员绑定的购物车相同，直接返回
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
        outCart, _ = c.GetCartByKey(cartKey)
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

// 合并购物车，并返回新的购物车
func (c *cartImpl) Combine(ic cart.ICart) cart.ICart {
	if ic.GetAggregateRootId() != c.GetAggregateRootId() {
		for _, v := range ic.GetValue().Items {
			if item, err := c.AddItem(v.VendorId, v.ShopId,
				v.SkuId, v.Quantity, v.Checked == 1); err == nil {
				if v.Checked == 1 {
					item.Checked = 1
				}
			}
		}
		ic.Destroy() //合并后,需销毁购物车
	}
	c._snapMap = nil //clean
	return c
}

// 设置购买会员
func (c *cartImpl) SetBuyer(buyerId int) error {
	if c._value.BuyerId > 0 {
		return cart.ErrCartBuyerBind
	}
	c._value.BuyerId = buyerId
	memCart := c._rep.GetMemberCurrentCart(buyerId)
	if memCart != nil && memCart.Key() != c.Key() {
		c.Combine(memCart)
	}
	_, err := c.Save()
	return err
}

// 设置购买会员收货地址
func (c *cartImpl) SetBuyerAddress(addressId int) error {
	if c._value.BuyerId < 0 {
		return cart.ErrCartNoBuyer
	}
	m := c._memberRep.GetMember(c._value.BuyerId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	addr := m.Profile().GetDeliver(addressId)
	if addr == nil {
		return member.ErrNoSuchDeliverAddress
	}
	return c.setBuyerAddress(addressId)
}

func (c *cartImpl) setBuyerAddress(addressId int) error {
	c._value.DeliverId = addressId
	_, err := c.Save()
	return err
}

// 标记商品结算
func (c *cartImpl) SignItemChecked(skuArr []int) error {
	mp := c.Items()
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
	err := c.Check()
	if err == nil {
		_, err = c.Save()
	}
	return err
}

// 结算数据持久化
func (c *cartImpl) SettlePersist(shopId, paymentOpt, deliverOpt, deliverId int) error {
	//var shop shop.IShop
	var deliver member.IDeliverAddress
	var err error

	if shopId > 0 {
		//var mch merchant.IMerchant
		//mch, err = c._partnerRep.GetMerchant(c._merchantId)
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

	if c._value.BuyerId > 0 && deliverId > 0 {
		var m member.IMember = c._memberRep.GetMember(c._value.BuyerId)
		if m == nil {
			return member.ErrNoSuchMember
		}
		deliver = m.Profile().GetDeliver(deliverId)
		if deliver == nil {
			return member.ErrInvalidSession
		}
		c._deliver = deliver
		c._value.DeliverId = deliverId
	}

	c._value.PaymentOpt = paymentOpt
	c._value.DeliverOpt = deliverOpt
	return nil
}

// 获取结算数据
func (c *cartImpl) GetSettleData() (s shop.IShop, d member.IDeliverAddress,
	paymentOpt, deliverOpt int) {
	//var err error
	if c._value.ShopId > 0 && c._shop == nil {
		//var pt merchant.IMerchant
		//pt, err = c._partnerRep.GetMerchant(c._merchantId)
		//if err == nil {
		//	c._shop = pt.ShopManager().GetShop(c._value.ShopId)
		//}
		//todo: not implement
	}
	if c._deliver == nil {
		pm := c._memberRep.GetMember(c._value.BuyerId).Profile()
		if c._value.DeliverId > 0 {
			c._deliver = pm.GetDeliver(c._value.DeliverId)
		} else {
			c._deliver = pm.GetDefaultAddress()
			if c._deliver != nil {
				c.setBuyerAddress(c._deliver.GetDomainId())
			}
		}
	}
	return c._shop, c._deliver, c._value.PaymentOpt, c._value.DeliverOpt
}

// 保存购物车
func (c *cartImpl) Save() (int, error) {
	c._value.UpdateTime = time.Now().Unix()
	id, err := c._rep.SaveShoppingCart(c._value)
	c._value.Id = id
	if c._value.Items != nil {
		for _, v := range c._value.Items {
			if v.Quantity <= 0 {
				c._rep.RemoveCartItem(v.Id)
			} else {
				v.Id, err = c._rep.SaveCartItem(v)
			}
		}
	}
	return id, err
}

// 释放购物车,如果购物车的商品全部结算,则返回true
func (c *cartImpl) Release() bool {
	checked := []int{}
	for i, v := range c._value.Items {
		if v.Checked == 1 {
			checked = append(checked, i)
		}
	}
	// 如果为部分结算,则移除商品并返回false
	if len(checked) < len(c._value.Items) {
		for _, i := range checked {
			v := c._value.Items[i]
			c.RemoveItem(v.SkuId, v.Quantity)
		}
		c.Save()
		return false
	}
	return true
}

// 销毁购物车
func (c *cartImpl) Destroy() (err error) {
	c._snapMap = nil //clean
	if err = c._rep.EmptyCartItems(c.GetAggregateRootId()); err == nil {
		return c._rep.DeleteCart(c.GetAggregateRootId())
	}
	return err
}

// 获取总览信息
func (c *cartImpl) GetSummary() string {
	if len(c._summary) != 0 {
		return c._summary
	}
	buf := bytes.NewBufferString("")

	list := c.getSnapshotsMap(c._value.Items)
	if list != nil {
		length := len(list)
		for i, v := range c._value.Items {
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
func (c *cartImpl) GetJsonItems() []byte {
	var goods []*order.OrderGoods = make([]*order.OrderGoods, len(c._value.Items))
	for i, v := range c._value.Items {
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
func (c *cartImpl) GetFee() (totalFee float32, orderFee float32) {
	var qua float32
	for _, v := range c._value.Items {
		if v.Checked == 1 {
			qua = float32(v.Quantity)
			totalFee += v.Price * qua
			orderFee += v.SalePrice * qua
		}
	}
	return totalFee, orderFee
}
