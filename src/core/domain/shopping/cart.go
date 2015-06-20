package shopping

import (
	"bytes"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/domain/interface/sale"
	"go2o/src/core/domain/interface/shopping"
	"go2o/src/core/infrastructure/domain"
	"strconv"
	"time"
	"go2o/src/core/domain/interface/valueobject"
)

type Cart struct {
	_value       *shopping.ValueCart
	_saleRep     sale.ISaleRep
	_shoppingRep shopping.IShoppingRep
	_partnerRep  partner.IPartnerRep
	_memberRep   member.IMemberRep
	_partnerId   int
	_summary     string
	_shop        partner.IShop
	_deliver     member.IDeliver
}

func createCart(partnerRep partner.IPartnerRep, memberRep member.IMemberRep, saleRep sale.ISaleRep,
	shoppingRep shopping.IShoppingRep, partnerId int, val *shopping.ValueCart) shopping.ICart {
	return (&Cart{
		_value:       val,
		_partnerId:   partnerId,
		_partnerRep:  partnerRep,
		_memberRep:   memberRep,
		_shoppingRep: shoppingRep,
		_saleRep:     saleRep,
	}).init()
}

//todo: partnerId 应去掉，可能在多个商家买东西
func newCart(partnerRep partner.IPartnerRep, memberRep member.IMemberRep, saleRep sale.ISaleRep,
	shoppingRep shopping.IShoppingRep, partnerId int, buyerId int) shopping.ICart {
	unix := time.Now().Unix()
	cartKey := domain.GenerateCartKey(unix, time.Now().Nanosecond())
	value := &shopping.ValueCart{
		CartKey:    cartKey,
		BuyerId:    buyerId,
		ShopId:     0,
		DeliverId:  0,
		PaymentOpt: 1,
		DeliverOpt: 1,
		CreateTime: unix,
		UpdateTime: unix,
		Items:nil,
	}

	return (&Cart{
		_value:       value,
		_partnerRep:  partnerRep,
		_memberRep:   memberRep,
		_partnerId:   partnerId,
		_shoppingRep: shoppingRep,
		_saleRep:     saleRep,
	}).init()
}

func (this *Cart) init()shopping.ICart{
	// 初始化购物车的信息
	if this._value != nil && this._value.Items != nil {
		this.setAttachGoodsInfo(this._value.Items)
	}
	return this
}


func (this *Cart) setAttachGoodsInfo(items []*shopping.ValueCartItem) {
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
		goods, err := this._saleRep.GetGoodsByIds(ids...)
		if err == nil {
			var goodsMap = make(map[int]*valueobject.Goods, len(goods))
			for _, v := range goods {
				goodsMap[v.GoodsId] = v
			}

			for _, v := range items {
				gv, ok := goodsMap[v.GoodsId]
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

func (this *Cart) GetDomainId() int {
	return this._value.Id
}

func (this *Cart) GetValue() shopping.ValueCart {
	return *this._value
}

// 添加项
func (this *Cart) AddItem(goodsId, num int) *shopping.ValueCartItem {
	if this._value.Items == nil {
		this._value.Items = []*shopping.ValueCartItem{}
	}

	// 添加数量
	for _, v := range this._value.Items {
		if v.GoodsId == goodsId {
			v.Num = v.Num + num
			return v
		}
	}

	sl := this._saleRep.GetSale(this._partnerId)
	goods := sl.GetGoods(goodsId);
	gv := goods.GetPackedValue()
	snap := goods.GetLatestSnapshot()

	if goods != nil {
		v := &shopping.ValueCartItem{
			CartId:     this.GetDomainId(),
			SnapshotId: snap.Id,
			GoodsId:    goodsId,
			Num:        num,
			Name:       gv.Name,
			GoodsNo:    gv.GoodsNo,
			Image:      gv.Image,
			Price:      gv.Price,
			SalePrice:  gv.SalePrice,
		}
		this._value.Items = append(this._value.Items, v)
		return v
	}
	return nil
}

// 移出项
func (this *Cart) RemoveItem(goodsId, num int) error {
	if this._value.Items == nil {
		return shopping.ErrEmptyShoppingCart
	}

	// 删除数量
	for _, v := range this._value.Items {
		if v.GoodsId == goodsId {
			if newNum := v.Num - num; newNum <= 0 {
				// 移出购物车
				//this.value.Items = append(this.value.Items[:i],this.value.Items[i+1:]...)
				v.Num = 0
			} else {
				v.Num = newNum
			}
			break
		}
	}
	return nil
}

// 合并购物车，并返回新的购物车
func (this *Cart) Combine(c shopping.ICart) (shopping.ICart, error) {
	if c.GetDomainId() != this.GetDomainId() {
		for _, v := range c.GetValue().Items {
			this.AddItem(v.GoodsId, v.Num)
		}
	}
	return this, nil
}

// 设置购买会员
func (this *Cart) SetBuyer(buyerId int) error {
	if this._value.BuyerId > 0 {
		return shopping.ErrCartBuyerBinded
	}
	this._value.BuyerId = buyerId
	_, err := this.Save()
	return err
}

// 结算数据持久化
func (this *Cart) SettlePersist(shopId, paymentOpt, deliverOpt, deliverId int) error {
	var shop partner.IShop
	var deliver member.IDeliver
	var err error

	if shopId > 0 {
		var pt partner.IPartner
		pt, err = this._partnerRep.GetPartner(this._partnerId)
		if err != nil {
			return err
		}
		shop = pt.GetShop(shopId)
		if shop == nil {
			return partner.ErrNoSuchShop
		}
		this._shop = shop
		this._value.ShopId = shopId
	}

	if this._value.BuyerId > 0 && deliverId > 0 {
		var m member.IMember
		m, err = this._memberRep.GetMember(this._value.BuyerId)
		if err != nil {
			return err
		}
		deliver = m.GetDeliver(deliverId)
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
func (this *Cart) GetSettleData() (s partner.IShop, d member.IDeliver, paymentOpt, deliverOpt int) {
	var err error
	if this._value.ShopId > 0 && this._shop == nil {
		var pt partner.IPartner
		pt, err = this._partnerRep.GetPartner(this._partnerId)
		if err == nil {
			this._shop = pt.GetShop(this._value.ShopId)
		}
	}
	if this._value.DeliverId > 0 && this._deliver == nil {
		var m member.IMember
		m, err = this._memberRep.GetMember(this._value.BuyerId)
		if err == nil {
			this._deliver = m.GetDeliver(this._value.DeliverId)
		}
	}
	return this._shop, this._deliver, this._value.PaymentOpt, this._value.DeliverOpt
}

// 保存购物车
func (this *Cart) Save() (int, error) {
	rep := this._shoppingRep
	this._value.UpdateTime = time.Now().Unix()
	id, err := rep.SaveShoppingCart(this._value)
	this._value.Id = id

	if this._value.Items != nil {
		for _, v := range this._value.Items {
			if v.Num <= 0 {
				rep.RemoveCartItem(v.Id)
			} else {
				i, err := rep.SaveCartItem(v)
				if err != nil {
					v.Id = i
				}
			}
		}
	}

	return id, err
}

// 销毁购物车
func (this *Cart) Destroy() (err error) {
	if err = this._shoppingRep.EmptyCartItems(this.GetDomainId()); err == nil {
		return this._shoppingRep.DeleteCart(this.GetDomainId())
	}
	return err
}

// 绑定订单
//func (this *Cart) BindOrder(orderNo string) error {
//	if this.GetDomainId() <= 0 || len(this._value.OrderNo) != 0 {
//		return shopping.ErrDisallowBindForCart
//	}
//	this._value.OrderNo = orderNo
//	this._value.IsBought = 1
//	_, err := this.Save()
//	return err
//}

// 获取总览信息
func (this *Cart) GetSummary() string {
	if len(this._summary) != 0 {
		return this._summary
	}
	buf := bytes.NewBufferString("")
	length := len(this._value.Items)

	var snap *sale.GoodsSnapshot
	for i, v := range this._value.Items {

		snap = this._saleRep.GetGoodsSnapshot(v.SnapshotId)
		if snap != nil {
			buf.WriteString(snap.GoodsName)
			if len(snap.SmallTitle) != 0 {
				buf.WriteString("(" + snap.SmallTitle + ")")
			}
			buf.WriteString("*" + strconv.Itoa(v.Num))
			if i < length-1 {
				buf.WriteString("\n")
			}
		}
	}
	return buf.String()
}

// 获取订单金额,返回totalFee为总额，
// orderFee为实际订单的金额(扣去促销优惠等后的金额)
func (this *Cart) GetFee() (totalFee float32, orderFee float32) {
	var qua float32
	for _, v := range this._value.Items {
		qua = float32(v.Num)
		totalFee = totalFee + v.Price*qua
		orderFee = orderFee + v.SalePrice*qua
	}
	return totalFee, orderFee
}
