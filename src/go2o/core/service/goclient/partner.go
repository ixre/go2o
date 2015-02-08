/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package goclient

import (
	"fmt"
	"github.com/atnet/gof/net/jsv"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/partner"
	"go2o/core/domain/interface/sale"
	"go2o/core/domain/interface/shopping"
	"go2o/core/dto"
	"strconv"
)

type partnerClient struct {
	conn *jsv.TCPConn
}

func (this *partnerClient) GetPartner(partnerId int, secret string) (a *partner.ValuePartner, err error) {
	var result jsv.Result
	err = this.conn.WriteAndDecode([]byte(fmt.Sprintf(
		`{"partner_id":"%s","secret":"%s"}>>Partner.GetPartner`,
		strconv.Itoa(partnerId), secret)), &result, 256)
	if err != nil {
		return nil, err
	}
	a = &partner.ValuePartner{}
	err = jsv.UnmarshalMap(result.Data, &a)
	return a, err
}

func (this *partnerClient) Category(partnerId int, secret string) (a []sale.ValueCategory, err error) {
	var result jsv.Result
	err = this.conn.WriteAndDecode([]byte(fmt.Sprintf(
		`{"partner_id":"%s","secret":"%s"}>>Partner.Category`,
		strconv.Itoa(partnerId), secret)), &result, 2048)
	if err != nil {
		return nil, err
	}
	a = []sale.ValueCategory{}
	err = jsv.UnmarshalMap(result.Data, &a)
	return a, err
}

func (this *partnerClient) GetShops(partnerId int, secret string) (a []partner.ValueShop, err error) {
	var result jsv.Result
	err = this.conn.WriteAndDecode([]byte(fmt.Sprintf(
		`{"partner_id":"%s","secret":"%s"}>>Partner.GetShops`,
		strconv.Itoa(partnerId), secret)), &result, 2048)
	if err != nil {
		return nil, err
	}
	a = []partner.ValueShop{}
	err = jsv.UnmarshalMap(result.Data, &a)
	return a, err
}

func (this *partnerClient) GetItems(partnerId int, secret string, categoryId int, getNum int) (
	a []*dto.ListGoods, err error) {
	var result jsv.Result
	err = this.conn.WriteAndDecode([]byte(fmt.Sprintf(
		`{"partner_id":"%s","secret":"%s","cid":"%d","num":"%d"}>>Partner.GetItems`,
		strconv.Itoa(partnerId), secret, categoryId, getNum)), &result, 2048)
	if err != nil {
		return nil, err
	}
	a = []*dto.ListGoods{}
	err = jsv.UnmarshalMap(result.Data, &a)
	return a, err
}

func (this *partnerClient) GetHost(partnerId int, secret string) (host string, err error) {
	var result jsv.Result
	err = this.conn.WriteAndDecode([]byte(fmt.Sprintf(
		`{"partner_id":"%s","secret":"%s"}>>Partner.GetHost`,
		strconv.Itoa(partnerId), secret)), &result, 2048)
	if err != nil {
		return "", err
	}
	return result.Data.(string), nil
}

func (this *partnerClient) GetSiteConf(partnerId int, secret string) (c *partner.SiteConf, err error) {
	var result jsv.Result
	err = this.conn.WriteAndDecode([]byte(fmt.Sprintf(
		`{"partner_id":"%s","secret":"%s"}>>Partner.GetSiteConf`,
		strconv.Itoa(partnerId), secret)), &result, 2048)
	if err != nil {
		return nil, err
	}
	c = new(partner.SiteConf)
	err = jsv.UnmarshalMap(result.Data, &c)
	return c, nil
}

//根据订单号获取订单
func (this *partnerClient) BuildOrder(partnerId int, secret string, memberId int,
	cartData string, couponCode string) (string, error) {
	var result jsv.Result
	err := this.conn.WriteAndDecode([]byte(fmt.Sprintf(
		`{"partner_id":"%s","secret":"%s","member_id":"%d","cart":"%s","coupon_code":"%s"}>>Partner.BuildOrder`,
		strconv.Itoa(partnerId), secret, memberId, cartData,
		couponCode)), &result, 2048)
	if err != nil {
		return "{}", err
	}

	if result.Data == nil {
		return "{}", err
	}
	return result.Data.(string), err
}

// 提交订单，并返回订单号
func (this *partnerClient) SubmitOrder(partnerId int, secret string, memberId int,
	shopId int, paymentMethod int, deliverAddrId int, cart string, couponCode string,
	note string) (orderNo string, err error) {
	var result jsv.Result
	err = this.conn.WriteAndDecode([]byte(fmt.Sprintf(
		`{"partner_id":"%d","secret":"%s","member_id":"%d","shop_id":"%d",`+
			`"pay_method":"%d","addr_id":"%d","cart":"%s","coupon_code":"%s"`+
			`,"note":"%s"}>>Partner.SubmitOrder`,
		partnerId,
		secret,
		memberId,
		shopId,
		paymentMethod,
		deliverAddrId,
		cart,
		couponCode,
		note)), &result, 256)
	if err != nil {
		return "", err
	}
	return result.Data.(string), err
}

//根据订单号获取订单
func (this *partnerClient) GetOrderByNo(partnerId int, secret string, order_no string) (*shopping.ValueOrder, error) {
	var result *shopping.ValueOrder = new(shopping.ValueOrder)
	err := this.conn.WriteAndDecode([]byte(fmt.Sprintf(
		`{"partner_id":"%d","secret":"%s","order_no":"%s"}>>Partner.GetOrderByNo`,
		partnerId, secret, order_no)), result, 2048)
	return result, err
}

func (this *partnerClient) UserIsExist(partnerId int, secret string, usr string) bool {
	var result jsv.Result
	err := this.conn.WriteAndDecode([]byte(fmt.Sprintf(
		`{"partner_id":"%s","secret":"%s","usr":"%s"}>>Partner.CheckUsrExist`,
		strconv.Itoa(partnerId), secret, usr)), &result, 72)
	if err != nil {
		return true
	}
	return result.Data.(bool)
}

//注册会员
func (this *partnerClient) RegisterMember(m *member.ValueMember, ptId, tgId int, cardId string) (
	b bool, err error) {
	var result jsv.Result
	err = this.conn.WriteAndDecode([]byte(fmt.Sprintf(
		`{"partner_id":"%d","tg_id":"%d","card_id":"%s","json":%s}>>Partner.RegisterMember`,
		ptId, tgId, cardId, jsv.MarshalString(m))), &result, -1)
	if err != nil {
		return false, err
	}
	return result.Result, err
}

func (this *partnerClient) GetShoppingCart(partnerId int, memberId int, cartKey string) (
	a *dto.ShoppingCart) {
	var result jsv.Result
	err := this.conn.WriteAndDecode([]byte(fmt.Sprintf(
		`{"partner_id":"%d","member_id":"%d","cart_key":"%s"}>>Partner.GetShoppingCart`,
		partnerId, memberId, cartKey)), &result, 1024)
	if err != nil {
		return nil
	}
	a = &dto.ShoppingCart{}
	err = jsv.UnmarshalMap(result.Data, &a)
	return a
}

func (this *partnerClient) AddCartItem(partnerId int, memberId int, cartKey string, goodsId, num int) (*dto.CartItem, error) {
	var result dto.CartItem
	err := this.conn.WriteAndDecode([]byte(fmt.Sprintf(
		`{"partner_id":"%d","member_id":"%d","cart_key":"%s","goods_id":"%d","num":"%d"}>>Partner.AddCartItem`,
		partnerId, memberId, cartKey, goodsId, num)), &result, 1024)
	if err != nil {
		return nil, err
	}
	return &result, err
}
