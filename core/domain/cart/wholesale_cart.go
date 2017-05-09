package cart

import (
	"encoding/json"
	"github.com/jsix/gof/util"
	"go2o/core/domain/interface/cart"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/infrastructure/format"
	"strconv"
	"time"
)

var _ cart.ICart = new(wholesaleCartImpl)
var _ cart.IWholesaleCart = new(wholesaleCartImpl)

type wholesaleCartImpl struct {
	value      *cart.WsCart
	rep        cart.ICartRepo
	itemRepo   item.IGoodsItemRepo
	memberRepo member.IMemberRepo
	mchRepo    merchant.IMerchantRepo
	summary    string
	shop       shop.IShop
	deliver    member.IDeliverAddress
	snapMap    map[int32]*item.Snapshot
}

func CreateWholesaleCart(val *cart.WsCart, rep cart.ICartRepo,
	memberRepo member.IMemberRepo, mchRepo merchant.IMerchantRepo,
	itemRepo item.IGoodsItemRepo) cart.ICart {
	return (&wholesaleCartImpl{
		value:      val,
		rep:        rep,
		memberRepo: memberRepo,
		itemRepo:   itemRepo,
	}).init()
}

func (c *wholesaleCartImpl) init() cart.ICart {
	// 获取购物车项
	if c.GetAggregateRootId() > 0 {
		if c.value.Items == nil {
			c.value.Items = c.rep.SelectWsCartItem("cart_id=?",
				c.GetAggregateRootId())
		}
	} else {
		c.value.Items = []*cart.WsCartItem{}
	}
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
func (c *wholesaleCartImpl) BuyerId() int64 {
	return c.value.BuyerId
}

// 检查购物车(仅结算商品)
func (c *wholesaleCartImpl) Check() error {
	if c.value == nil || len(c.value.Items) == 0 {
		return cart.ErrEmptyShoppingCart
	}
	for _, v := range c.value.Items {
		if v.Checked == 1 {
			it := c.itemRepo.GetItem(v.ItemId)
			if it == nil {
				return item.ErrNoSuchItem // 没有商品
			}
			// 验证批发权限
			wsIt := it.Wholesale()
			if wsIt == nil || !wsIt.CanWholesale() {
				return item.ErrItemWholesaleOff
			}
			// 验证库存
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

// 添加项
func (c *wholesaleCartImpl) put(itemId, skuId int32, num int32) (*cart.WsCartItem, error) {
	var err error
	if c.value.Items == nil {
		c.value.Items = []*cart.WsCartItem{}
	}
	var sku *item.Sku
	it := c.itemRepo.GetItem(itemId)
	if it == nil {
		return nil, item.ErrNoSuchItem // 没有商品
	}
	iv := it.GetValue()
	// 库存,如有SKU，则使用SKU的库存
	stock := iv.StockNum
	// 判断是否上架
	if iv.ShelveState != item.ShelvesOn {
		return nil, item.ErrNotOnShelves //未上架
	}
	// 验证批发权限
	wsIt := it.Wholesale()
	if wsIt == nil || !wsIt.CanWholesale() {
		return nil, item.ErrItemWholesaleOff
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
			snapList := c.itemRepo.GetSnapshots(ids)
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
		gds := c.itemRepo.CreateItem(snap)
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
		it := c.itemRepo.GetItem(v.ItemId)
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

// 移出项
func (c *wholesaleCartImpl) Remove(itemId, skuId, quantity int32) error {
	if c.value.Items == nil {
		return cart.ErrEmptyShoppingCart
	}
	// 删除数量
	for _, v := range c.value.Items {
		if v.ItemId == itemId && v.SkuId == skuId {
			if newNum := v.Quantity - quantity; newNum <= 0 {
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

// 获取购物车编码
func (c *wholesaleCartImpl) Code() string {
	return c.value.Code
}

// 设置购买会员收货地址
func (c *wholesaleCartImpl) SetBuyerAddress(addressId int64) error {
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

func (c *wholesaleCartImpl) setBuyerAddress(addressId int64) error {
	c.value.DeliverId = addressId
	_, err := c.Save()
	return err
}

// 标记商品结算
func (c *wholesaleCartImpl) SignItemChecked(items []*cart.ItemPair) error {
	mp := c.getItems()
	// 遍历购物车商品，默认不结算。
	for _, item := range mp {
		item.Checked = 0
		// 如果传入结算商品信息，则标记购物车项结算状态
		for _, v := range items {
			if v.SkuId == item.SkuId && v.ItemId == item.ItemId {
				item.Checked = v.Checked
				break
			}
		}
	}
	return c.Check()
}

// 结算数据持久化
func (c *wholesaleCartImpl) SettlePersist(shopId, paymentOpt, deliverOpt int32,
	addressId int64) error {
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
	c.value.UpdateTime = time.Now().Unix()
	id, err := util.I32Err(c.rep.SaveWsCart(c.value))
	c.value.ID = id
	if c.value.Items != nil {
		for _, v := range c.value.Items {
			if v.Quantity <= 0 {
				//c.rep.RemoveCartItem(v.ID)
				c.rep.BatchDeleteWsCartItem("id=?", v.ID)
			} else {
				v.CartId = c.GetAggregateRootId()
				v.ID, err = util.I32Err(c.rep.SaveWsCartItem(v))
			}
		}
	}
	return id, err
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

// 获取购物车商品Jdo数据
func (c *wholesaleCartImpl) getJdoItemData(list []*cart.WsCartItem,
	itemId int32) cart.WCartItemJdo {
	it := c.itemRepo.GetItem(itemId)
	v := it.GetValue()
	itw := it.Wholesale()
	itJdo := cart.WCartItemJdo{
		ItemId:    int64(itemId),
		ItemName:  v.Title,
		ItemImage: v.Image,
		SkuList:   []cart.WCartSkuJdo{},
		Data:      map[string]string{},
	}
	skuSignMap := make(map[int64]bool)
	for _, v := range list {
		if v.ItemId != itemId || skuSignMap[int64(v.SkuId)] {
			continue
		}
		skuV := it.GetSku(v.SkuId)
		skuJdo := cart.WCartSkuJdo{
			SkuId:            int64(v.SkuId),
			SkuImage:         skuV.Image,
			SpecWord:         skuV.SpecWord,
			Quantity:         v.Quantity,
			Price:            0,
			DiscountPrice:    0,
			CanSalesQuantity: skuV.Stock,
			JData:            "{}",
		}
		mp := map[string]interface{}{}
		mp["canSalesQuantity"] = skuV.Stock
		c.setJdoSkuData(itw, &skuJdo, mp)

		itJdo.SkuList = append(itJdo.SkuList, skuJdo)
		skuSignMap[skuJdo.SkuId] = true
	}
	return itJdo
}

func (c *wholesaleCartImpl) setJdoSkuData(itw item.IWholesaleItem,
	sku *cart.WCartSkuJdo, mp map[string]interface{}) {
	prArr := itw.GetSkuPrice(int32(sku.SkuId))

	var min float64
	priceRange := [][]string{}
	for _, v := range prArr {
		if min == 0 {
			min = v.WholesalePrice
		}
		if v.WholesalePrice < min {
			min = v.WholesalePrice
		}
		priceRange = append(priceRange, []string{
			strconv.Itoa(int(v.RequireQuantity)),
			format.DecimalToString(v.WholesalePrice),
		})
	}
	sku.Price = min
	sku.DiscountPrice = min
	mp["priceRange"] = priceRange
	data, _ := json.Marshal(mp)
	sku.JData = string(data)

}

// Jdo数据
func (c *wholesaleCartImpl) JdoData() *cart.WCartJdo {
	var jdo cart.WCartJdo = []cart.WCartSellerJdo{}
	venMap := make(map[int32]int)
	itSignMap := make(map[int64]bool)
	for _, v := range c.value.Items {
		// 如果已处理过商品，则跳过
		if v.VendorId <= 0 || itSignMap[int64(v.ItemId)] {
			continue
		}
		vi, ok := venMap[v.VendorId]
		//初始化VendorJdo
		if !ok {
			vJdo := cart.WCartSellerJdo{
				SellerId: v.VendorId,
				Items:    []cart.WCartItemJdo{},
				Data:     map[string]string{},
			}
			jdo = append(jdo, vJdo)
			vi = len(jdo) - 1
			venMap[v.VendorId] = vi
		}
		// 设置商品信息
		itJdo := c.getJdoItemData(c.value.Items, v.ItemId)
		jdo[vi].Items = append(jdo[vi].Items, itJdo)
		itSignMap[int64(v.ItemId)] = true
	}
	return &jdo
}
