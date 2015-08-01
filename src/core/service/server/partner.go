/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package server

import (
	"encoding/json"
	"github.com/atnet/gof/net/jsv"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/sale"
	"go2o/src/core/domain/interface/shopping"
	"go2o/src/core/dto"
	"go2o/src/core/infrastructure/domain"
	"go2o/src/core/service/dps"
	"strconv"
	"time"
)

// 合作商户的接口
type Partner struct{}

func (this *Partner) GetPartner(m *jsv.Args, r *jsv.Result) error {
	return nil
	_, err, e := VerifyPartner(m)
	if err != nil {
		return err
	}
	r.Result = true
	r.Data = e
	return nil
}

func (this *Partner) GetSiteConf(m *jsv.Args, r *jsv.Result) error {
	partnerId, err, _ := VerifyPartner(m)
	if err != nil {
		return err
	}

	siteConf := dps.PartnerService.GetSiteConf(partnerId)
	r.Result = true
	r.Data = siteConf
	return nil
}

func (this *Partner) GetHost(m *jsv.Args, r *jsv.Result) error {
	partnerId, err := strconv.Atoi((*m)["partner_id"].(string))
	if err != nil {
		return err
	}
	host := dps.PartnerService.GetPartnerMajorHost(partnerId)
	r.Data = host
	r.Result = true
	return nil
}

func (this *Partner) GetShops(m *jsv.Args, r *jsv.Result) error {
	partnerId, err, _ := VerifyPartner(m)

	if err != nil {
		return err
	}

	shops := dps.PartnerService.GetShopsOfPartner(partnerId)
	r.Result = true
	r.Data = shops
	return nil
}

func (this *Partner) Category(m *jsv.Args, r *jsv.Result) error {
	partnerId, err, _ := VerifyPartner(m)
	if err != nil {
		return err
	}

	categories := dps.SaleService.GetCategories(partnerId)
	r.Result = true
	r.Data = categories
	return nil
}

func (this *Partner) GetItems(m *jsv.Args, r *jsv.Result) error {
	partnerId, err, _ := VerifyPartner(m)
	if err != nil {
		return err
	}
	cid, _ := strconv.Atoi((*m)["cid"].(string))
	num, _ := strconv.Atoi((*m)["num"].(string))

	_, items := dps.SaleService.GetPagedOnShelvesGoods(partnerId, cid, 0, num)
	r.Result = true
	r.Data = items

	return nil
}

func (this *Partner) RegisterMember(m *jsv.Args, r *jsv.Result) error {

	var err error

	e := member.ValueMember{}

	if err = jsv.UnmarshalMap((*m)["json"], &e); err != nil {
		return err
	}
	var cardId string
	var tgId int
	var partnerId int

	cardId = (*m)["card_id"].(string)
	tgId, _ = strconv.Atoi((*m)["tg_id"].(string))
	partnerId, _ = strconv.Atoi((*m)["partner_id"].(string))

	//如果卡片ID为空时，自动生成
	if cardId == "" {
		cardId = time.Now().Format("200601021504")
	}
	e.RegFrom = "WEB"

	e.Pwd = domain.MemberSha1Pwd(e.Pwd)
	id, err := dps.MemberService.SaveMember(&e)

	if err == nil {
		dps.MemberService.SaveRelation(id, cardId, tgId, partnerId)
		r.Result = true
		return nil
	}
	return err
}

func (this *Partner) GetShoppingCart(m *jsv.Args, r *dto.ShoppingCart) error {
	partnerId, _ := strconv.Atoi((*m)["partner_id"].(string))
	memberId, _ := strconv.Atoi((*m)["member_id"].(string))
	var cartKey string = (*m)["cart_key"].(string)
	cart := dps.ShoppingService.GetShoppingCart(partnerId, memberId, cartKey)
	*r = *cart
	return nil
}

func (this *Partner) GetCartSettle(m *jsv.Args, r *dto.SettleMeta) error {
	partnerId, _ := strconv.Atoi((*m)["partner_id"].(string))
	memberId, _ := strconv.Atoi((*m)["member_id"].(string))
	var cartKey string = (*m)["cart_key"].(string)
	settle := dps.ShoppingService.GetCartSettle(partnerId, memberId, cartKey)
	*r = *settle
	return nil
}

func (this *Partner) BuildOrder(m *jsv.Args, r *jsv.Result) error {
	partnerId, err, _ := VerifyPartner(m)
	if err != nil {
		return err
	}

	memberId, err := strconv.Atoi((*m)["member_id"].(string))
	couponCode := (*m)["coupon_code"].(string)
	if err != nil {
		return err
	}

	data, err := dps.ShoppingService.BuildOrder(partnerId, memberId, "", couponCode)
	if err != nil {
		return err
	}

	js, _ := json.Marshal(data)

	r.Result = true
	r.Data = string(js)
	return nil
}

func (this *Partner) GetOrderByNo(m *jsv.Args, r *shopping.ValueOrder) error {
	partnerId, err, _ := VerifyPartner(m)
	if err != nil {
		return err
	}
	order := dps.ShoppingService.GetOrderByNo(partnerId, (*m)["order_no"].(string))
	if order != nil {
		*r = *order
	}
	return nil
}

func (this *Partner) CheckUsrExist(m *jsv.Args, r *jsv.Result) error {
	_, err, _ := VerifyPartner(m)
	if err != nil {
		return err
	}
	r.Result = true
	r.Data = dps.MemberService.CheckUsr((*m)["usr"].(string), 0)
	return nil
}

func (this *Partner) AddCartItem(m *jsv.Args, item *dto.CartItem) error {
	partnerId, _ := strconv.Atoi((*m)["partner_id"].(string))
	memberId, _ := strconv.Atoi((*m)["member_id"].(string))
	cartKey := (*m)["cart_key"].(string)
	goodsId, _ := strconv.Atoi((*m)["goods_id"].(string))
	num, _ := strconv.Atoi((*m)["num"].(string))

	v, err := dps.ShoppingService.AddCartItem(partnerId,
		memberId, cartKey, goodsId, num)
	if v == nil {
		return sale.ErrNoSuchGoods
	}

	*item = *v
	return err
}

func (this *Partner) SubCartItem(m *jsv.Args, r *jsv.Result) error {
	partnerId, _ := strconv.Atoi((*m)["partner_id"].(string))
	memberId, _ := strconv.Atoi((*m)["member_id"].(string))
	cartKey := (*m)["cart_key"].(string)
	goodsId, _ := strconv.Atoi((*m)["goods_id"].(string))
	num, _ := strconv.Atoi((*m)["num"].(string))
	err := dps.ShoppingService.SubCartItem(partnerId,
		memberId, cartKey, goodsId, num)
	//r.Result = err == nil
	return err
}

// 订单持久
func (this *Partner) OrderPersist(m *jsv.Args, r *jsv.Result) error {
	partnerId, _ := strconv.Atoi((*m)["partner_id"].(string))
	memberId, _ := strconv.Atoi((*m)["member_id"].(string))
	deliverId, _ := strconv.Atoi((*m)["deliver_id"].(string))
	paymentOpt, _ := strconv.Atoi((*m)["payment_opt"].(string))
	deliverOpt, _ := strconv.Atoi((*m)["deliver_opt"].(string))
	shopId, _ := strconv.Atoi((*m)["shop_id"].(string))
	return dps.ShoppingService.PrepareSettlePersist(partnerId, memberId, shopId, paymentOpt, deliverOpt, deliverId)
}

// 需要传递配送地址
func (this *Partner) SubmitOrder(m *jsv.Args, r *jsv.Result) error {
	partnerId, err, _ := VerifyPartner(m)
	if err != nil {
		return err
	}
	memberId, _ := strconv.Atoi((*m)["member_id"].(string))
	couponCode := (*m)["coupon_code"].(string)

	orderNo, err := dps.ShoppingService.SubmitOrder(partnerId, memberId, couponCode, true)
	if err != nil {
		return err
	} else {
		r.Result = true
		r.Data = orderNo
	}
	return nil
}
