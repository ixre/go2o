package shopping

import (
	"bytes"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/partner"
	"go2o/core/domain/interface/sale"
	"go2o/core/domain/interface/shopping"
	"go2o/core/infrastructure/domain"
	"strconv"
	"time"
)

type Cart struct {
	value       *shopping.ValueCart
	saleRep     sale.ISaleRep
	shoppingRep shopping.IShoppingRep
	partnerRep  partner.IPartnerRep
	memberRep   member.IMemberRep
	partnerId   int
	summary     string
	_shop       partner.IShop
	_deliver    member.IDeliver
}

func createCart(partnerRep partner.IPartnerRep, memberRep member.IMemberRep, saleRep sale.ISaleRep,
	shoppingRep shopping.IShoppingRep, partnerId int, val *shopping.ValueCart) shopping.ICart {
	return &Cart{
		value:       val,
		partnerId:   partnerId,
		partnerRep:  partnerRep,
		memberRep:   memberRep,
		shoppingRep: shoppingRep,
		saleRep:     saleRep,
	}
}

//todo: partnerId 应去掉，可能在多个商家买东西
func newCart(partnerRep partner.IPartnerRep, memberRep member.IMemberRep, saleRep sale.ISaleRep,
	shoppingRep shopping.IShoppingRep, partnerId int, buyerId int) shopping.ICart {
	unix := time.Now().Unix()
	cartKey := domain.GenerateCartKey(unix, time.Now().Nanosecond())
	value := &shopping.ValueCart{
		CartKey:    cartKey,
		BuyerId:    buyerId,
		OrderNo:    "",
		IsBought:   0,
		ShopId:     0,
		DeliverId:  0,
		PaymentOpt: 1,
		DeliverOpt: 1,
		CreateTime: unix,
		UpdateTime: unix,
	}

	return &Cart{
		value:       value,
		partnerRep:  partnerRep,
		memberRep:   memberRep,
		partnerId:   partnerId,
		shoppingRep: shoppingRep,
		saleRep:     saleRep,
	}
}

func (this *Cart) GetDomainId() int {
	return this.value.Id
}

func (this *Cart) GetValue() shopping.ValueCart {
	return *this.value
}

// 添加项
func (this *Cart) AddItem(goodsId, num int) *shopping.ValueCartItem {
	if this.value.Items == nil {
		this.value.Items = []*shopping.ValueCartItem{}
	}

	// 添加数量
	for _, v := range this.value.Items {
		if v.GoodsId == goodsId {
			v.Num = v.Num + num
			return v
		}
	}

	// 添加项
	pro := this.saleRep.GetValueGoods(this.partnerId, goodsId)
	if pro != nil {
		v := &shopping.ValueCartItem{
			CartId:     this.GetDomainId(),
			GoodsId:    goodsId,
			Num:        num,
			Name:       pro.Name,
			GoodsNo:    pro.GoodsNo,
			SmallTitle: pro.SmallTitle,
			Image:      pro.Image,
			Price:      pro.Price,
			SalePrice:  pro.SalePrice,
		}
		this.value.Items = append(this.value.Items, v)
		return v
	}
	return nil
}

// 移出项
func (this *Cart) RemoveItem(goodsId, num int) error {
	if this.value.Items == nil {
		return shopping.ErrEmptyShoppingCart
	}

	// 删除数量
	for i, v := range this.value.Items {
		if v.GoodsId == goodsId {
			if newNum := v.Num - num; newNum <= 0 {
				// 移出购物车
				this.value.Items = append(this.value.Items[:i],
					this.value.Items[i+1:]...)
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
	if this.value.BuyerId > 0 {
		return shopping.ErrCartBuyerBinded
	}
	this.value.BuyerId = buyerId
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
		pt, err = this.partnerRep.GetPartner(this.partnerId)
		if err != nil {
			return err
		}
		shop = pt.GetShop(shopId)
		if shop == nil {
			return partner.ErrNoSuchShop
		}
		this._shop = shop
		this.value.ShopId = shopId
	}

	if this.value.BuyerId > 0 && deliverId > 0 {
		var m member.IMember
		m, err = this.memberRep.GetMember(this.value.BuyerId)
		if err != nil {
			return err
		}
		deliver = m.GetDeliver(deliverId)
		if deliver == nil {
			return member.ErrInvalidSession
		}
		this._deliver = deliver
		this.value.DeliverId = deliverId
	}

	this.value.PaymentOpt = paymentOpt
	this.value.DeliverOpt = deliverOpt
	return nil
}

// 获取结算数据
func (this *Cart) GetSettleData() (s partner.IShop, d member.IDeliver, paymentOpt, deliverOpt int) {
	var err error
	if this.value.ShopId > 0 && this._shop == nil {
		var pt partner.IPartner
		pt, err = this.partnerRep.GetPartner(this.partnerId)
		if err == nil {
			this._shop = pt.GetShop(this.value.ShopId)
		}
	}
	if this.value.DeliverId > 0 && this._deliver == nil {
		var m member.IMember
		m, err = this.memberRep.GetMember(this.value.BuyerId)
		if err == nil {
			this._deliver = m.GetDeliver(this.value.DeliverId)
		}
	}
	return this._shop, this._deliver, this.value.PaymentOpt, this.value.DeliverOpt
}

// 保存购物车
func (this *Cart) Save() (int, error) {
	rep := this.shoppingRep
	this.value.UpdateTime = time.Now().Unix()
	id, err := rep.SaveShoppingCart(this.value)
	this.value.Id = id

	if this.value.Items != nil {
		for _, v := range this.value.Items {
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

// 绑定订单
func (this *Cart) BindOrder(orderNo string) error {
	if this.GetDomainId() <= 0 || len(this.value.OrderNo) != 0 {
		return shopping.ErrDisallowBindForCart
	}
	this.value.OrderNo = orderNo
	this.value.IsBought = 1
	_, err := this.Save()
	return err
}

// 获取总览信息
func (this *Cart) GetSummary() string {
	if len(this.summary) != 0 {
		return this.summary
	}
	buf := bytes.NewBufferString("")
	length := len(this.value.Items)
	var pro *sale.ValueGoods
	for i, v := range this.value.Items {
		pro = this.saleRep.GetValueGoods(this.partnerId, v.GoodsId)
		if pro != nil {
			buf.WriteString(pro.Name)
			if len(pro.SmallTitle) != 0 {
				buf.WriteString("(" + pro.SmallTitle + ")")
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
	for _, v := range this.value.Items {
		qua = float32(v.Num)
		totalFee = totalFee + v.Price*qua
		orderFee = orderFee + v.SalePrice*qua
	}
	return totalFee, orderFee
}
