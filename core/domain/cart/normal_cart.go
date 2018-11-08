package cart

import (
	"errors"
	"go2o/core/domain/interface/cart"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/infrastructure/domain"
	"time"
)

var _ cart.ICart = new(cartImpl)

type cartImpl struct {
	value      *cart.NormalCart
	rep        cart.ICartRepo
	goodsRepo  item.IGoodsItemRepo
	memberRepo member.IMemberRepo
	summary    string
	shop       shop.IShop
	deliver    member.IDeliverAddress
	snapMap    map[int64]*item.Snapshot
}

func CreateCart(val *cart.NormalCart, rep cart.ICartRepo,
	memberRepo member.IMemberRepo, goodsRepo item.IGoodsItemRepo) cart.ICart {
	c := &cartImpl{
		value:      val,
		rep:        rep,
		memberRepo: memberRepo,
		goodsRepo:  goodsRepo,
	}
	return c.init()
}

// 创建新的购物车
func NewNormalCart(code string, rep cart.ICartRepo, memberRepo member.IMemberRepo,
	goodsRepo item.IGoodsItemRepo) cart.ICart {
	unix := time.Now().Unix()
	if code == "" {
		code = domain.GenerateCartCode(unix, time.Now().Nanosecond())
	}
	value := &cart.NormalCart{
		CartCode:   code,
		DeliverId:  0,
		PaymentOpt: 1,
		CreateTime: unix,
		UpdateTime: unix,
		Items:      []*cart.NormalCartItem{},
	}
	return CreateCart(value, rep, memberRepo, goodsRepo)
}

func (c *cartImpl) init() cart.ICart {
	// 获取购物车项
	if c.GetAggregateRootId() > 0 {
		if c.value.Items == nil {
			c.value.Items = c.rep.SelectNormalCartItem("cart_id=?",
				c.GetAggregateRootId())
		}
	}
	if c.value.Items == nil {
		c.value.Items = []*cart.NormalCartItem{}
	}
	// 初始化购物车的信息
	if c.value != nil && c.value.Items != nil {
		c.setAttachGoodsInfo(c.value.Items)
	}
	return c
}

// 购物车种类
func (c *cartImpl) Kind() cart.CartKind {
	return cart.KNormal
}

// 获取买家编号
func (c *cartImpl) BuyerId() int64 {
	return c.value.BuyerId
}

// 检查购物车(仅结算商品)
func (c *cartImpl) Check() error {
	if c.value == nil || len(c.value.Items) == 0 {
		return cart.ErrEmptyShoppingCart
	}
	for _, v := range c.value.Items {
		if v.Checked == 1 {
			it := c.goodsRepo.GetItem(v.ItemId)
			if it == nil {
				return item.ErrNoSuchItem // 没有商品
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
func (c *cartImpl) getSnapshotsMap(items []*cart.NormalCartItem) map[int64]*item.Snapshot {
	if c.snapMap == nil && items != nil {
		l := len(items)
		c.snapMap = make(map[int64]*item.Snapshot, l)
		if l > 0 {
			var ids = make([]int64, l)
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

func (c *cartImpl) getBuyerLevelId() int32 {
	if c.value.BuyerId > 0 {
		m := c.memberRepo.GetMember(c.value.BuyerId)
		if m != nil {
			return m.GetValue().Level
		}
	}
	return 0
}

func (c *cartImpl) setItemInfo(snap *item.GoodsItem, level int32) {
	// 设置会员价
	if level > 0 {
		gds := c.goodsRepo.CreateItem(snap)
		snap.Price = gds.GetPromotionPrice(level)
	}
}

// 设置附加的商品信息
func (c *cartImpl) setAttachGoodsInfo(items []*cart.NormalCartItem) {
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
		} else {
			iv := it.GetValue()
			sku = &item.Sku{
				ProductId:   iv.ProductId,
				ItemId:      iv.ID,
				Title:       iv.Title,
				Image:       iv.Image,
				SpecData:    "",
				SpecWord:    "",
				Code:        iv.Code,
				RetailPrice: iv.RetailPrice,
				Price:       iv.Price,
				Cost:        iv.Cost,
				Weight:      iv.Weight,
				Bulk:        iv.Bulk,
				Stock:       iv.StockNum,
				SaleNum:     iv.SaleNum,
			}
		}
		v.Sku = item.ParseSkuMedia(it.GetValue(), sku)
	}
}

// 获取聚合根编号
func (c *cartImpl) GetAggregateRootId() int32 {
	return c.value.Id
}

func (c *cartImpl) GetValue() cart.NormalCart {
	return *c.value
}

// 获取商品集合
func (c *cartImpl) Items() []*cart.NormalCartItem {
	return c.getItems()
}

func (c *cartImpl) getItems() []*cart.NormalCartItem {
	return c.value.Items
}

// 添加项
func (c *cartImpl) Put(itemId, skuId int64, num int32) error {
	_, err := c.put(itemId, skuId, num)
	return err
}

// 添加项
func (c *cartImpl) put(itemId, skuId int64, num int32) (*cart.NormalCartItem, error) {
	var err error
	if c.value.Items == nil {
		c.value.Items = []*cart.NormalCartItem{}
	}
	var sku *item.Sku
	it := c.goodsRepo.GetItem(itemId)
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
	// 判断商品SkuId
	if skuId > 0 {
		sku = it.GetSku(skuId)
		if sku == nil {
			return nil, item.ErrNoSuchSku
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

	v := &cart.NormalCartItem{
		CartId:   c.GetAggregateRootId(),
		VendorId: iv.VendorId,
		ShopId:   iv.ShopId,
		ItemId:   iv.ID,
		SkuId:    skuId,
		Quantity: num,
		Sku:      item.ParseSkuMedia(iv, sku),
		Checked:  1,
	}
	c.value.Items = append(c.value.Items, v)
	return v, err
}

// 更新商品数量，如数量为0，则删除
func (c *cartImpl) Update(itemId, skuId int64, quantity int32) error {
	return errors.New("not implement")
}

// 移出项
func (c *cartImpl) Remove(itemId, skuId int64, quantity int32) error {
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

// 获取项
func (c *cartImpl) GetItem(itemId, skuId int64) *cart.NormalCartItem {
	if c.value != nil && c.value.Items != nil {
		for _, v := range c.value.Items {
			if v.ItemId == itemId && v.SkuId == skuId {
				return v
			}
		}
	}
	return nil
}

// 获取购物车编码
func (c *cartImpl) Code() string {
	return c.value.CartCode
}

// 合并购物车，并返回新的购物车
func (c *cartImpl) Combine(ic cart.ICart) cart.ICart {
	if ic.Kind() != cart.KNormal {
		panic("only retail cart can be combine!")
	}
	if id := ic.GetAggregateRootId(); id != c.GetAggregateRootId() {
		rc := ic.(cart.INormalCart)
		for _, v := range rc.Items() {
			if item, err := c.put(v.ItemId,
				v.SkuId, v.Quantity); err == nil {
				if v.Checked == 1 {
					item.Checked = 1
				}
			}
		}
		ic.Destroy() //合并后,需销毁购物车
	}
	c.snapMap = nil //clean
	return c
}

// 设置购买会员收货地址
func (c *cartImpl) SetBuyerAddress(addressId int64) error {
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

func (c *cartImpl) setBuyerAddress(addressId int64) error {
	c.value.DeliverId = addressId
	_, err := c.Save()
	return err
}

// 标记商品结算
func (c *cartImpl) SignItemChecked(items []*cart.ItemPair) error {
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
func (c *cartImpl) SettlePersist(shopId, paymentOpt, deliverOpt int32, addressId int64) error {
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

	c.value.PaymentOpt = paymentOpt
	return nil
}

// 获取结算数据
func (c *cartImpl) GetSettleData() (s shop.IShop, d member.IDeliverAddress,
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
	return c.shop, c.deliver, c.value.PaymentOpt
}

// 保存购物车
func (c *cartImpl) Save() (int32, error) {
	c.value.UpdateTime = time.Now().Unix()
	id, err := c.rep.SaveShoppingCart(c.value)
	c.value.Id = id
	if c.value.Items != nil {
		for _, v := range c.value.Items {
			if v.Quantity <= 0 {
				c.rep.RemoveCartItem(v.Id)
			} else {
				v.CartId = c.GetAggregateRootId()
				v.Id, err = c.rep.SaveCartItem(v)
			}
		}
	}
	return id, err
}

// 获取勾选的商品
func (c *cartImpl) CheckedItems(checked map[int64][]int64) []*cart.ItemPair {
	items := []*cart.ItemPair{}
	if checked != nil {
		for _, v := range c.value.Items {
			arr, ok := checked[int64(v.ItemId)]
			if !ok {
				continue
			}
			for _, skuId := range arr {
				if skuId == int64(v.SkuId) {
					items = append(items, &cart.ItemPair{
						ItemId:   int64(v.ItemId),
						SkuId:    skuId,
						SellerId: v.VendorId,
						Quantity: v.Quantity,
					})
				}
			}
		}
	}
	return items
}

// 释放购物车,如果购物车的商品全部结算,则返回true
func (c *cartImpl) Release(_checked map[int64][]int64) bool {
	var checked []int
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
func (c *cartImpl) Destroy() (err error) {
	c.snapMap = nil //clean
	if err = c.rep.EmptyCartItems(c.GetAggregateRootId()); err == nil {
		return c.rep.DeleteCart(c.GetAggregateRootId())
	}
	return err
}
