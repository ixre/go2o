package impl

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : cart_service.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-09-27 11:04
 * description :
 * history :
 */

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/ixre/gof/util"
	"go2o/core/domain/interface/cart"
	proItem "go2o/core/domain/interface/item"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/dto"
	"go2o/core/service/parser"
	"go2o/core/service/proto"
	"strconv"
	"strings"
)

var _ proto.CartServiceServer = new(cartServiceImpl)

type cartServiceImpl struct {
	itemRepo proItem.IGoodsItemRepo
	cartRepo cart.ICartRepo
	mchRepo  merchant.IMerchantRepo
	shopRepo shop.IShopRepo
	serviceUtil
}

func NewCartService(cartRepo cart.ICartRepo,
	goodsRepo proItem.IGoodsItemRepo,
	mchRepo merchant.IMerchantRepo,
	shopRepo shop.IShopRepo) *cartServiceImpl {
	return &cartServiceImpl{
		cartRepo: cartRepo,
		itemRepo: goodsRepo,
		mchRepo:  mchRepo,
		shopRepo: shopRepo,
	}
}

/*---------------- 批发购物车 ----------------*/

// 批发购物车接口
func (s *cartServiceImpl) WholesaleCartV1(_ context.Context, r *proto.WsCartRequest) (*proto.Result, error) {
	//todo: check member
	c := s.cartRepo.GetMyCart(r.MemberId, cart.KWholesale)
	if r.Data == nil {
		r.Data = map[string]string{}
	}
	switch r.Action {
	case "GET":
		return s.wsGetCart(c, r.Data)
	case "MINI":
		return s.wsGetSimpleCart(c, r.Data)
	case "PUT":
		return s.wsPutItem(c, r.Data)
	case "UPDATE":
		return s.wsUpdateItem(c, r.Data)
	case "CHECK":
		return s.wsCheckCart(c, r.Data)
	}
	return s.result(errors.New("unknown action")), nil
}

// 转换勾选字典,数据如：{"1":["10","11"],"2":["20","21"]}
func (s *cartServiceImpl) parseCheckedMap(data string) (m map[int64][]int64) {
	if data != "" && data != "{}" {
		src := map[string][]string{}
		err := json.Unmarshal([]byte(data), &src)
		if err == nil {
			m = map[int64][]int64{}
			for k, v := range src {
				itemId, _ := strconv.Atoi(k)
				var skuList []int64
				for _, v2 := range v {
					skuId, _ := strconv.Atoi(v2)
					skuList = append(skuList, int64(skuId))
				}
				m[int64(itemId)] = skuList
			}
			return m
		}
	}
	return nil
}

// 获取可结算的购物车
func (s *cartServiceImpl) wsGetCart(c cart.ICart, data map[string]string) (*proto.Result, error) {
	//统计checked
	checked := s.parseCheckedMap(data["checked"])
	checkout := data["checkout"] == "true"
	v := c.(cart.IWholesaleCart).JdoData(checkout, checked)
	if v != nil {
		for _, v2 := range v.Seller {
			mch := s.mchRepo.GetMerchant(int(v2.SellerId))
			if mch != nil {
				v2.Data["SellerName"] = mch.GetValue().CompanyName
			}
		}
	}
	return s.success(nil), nil
}

// 获取简易的购物车
func (s *cartServiceImpl) wsGetSimpleCart(c cart.ICart, data map[string]string) (*proto.Result, error) {
	size, err := strconv.Atoi(data["size"])
	if err != nil {
		size = 5
	}
	qd := c.(cart.IWholesaleCart).QuickJdoData(size)
	return s.success(map[string]string{"JsonData": qd}), nil
}

// 转换提交到购物车的数据(PUT和UPDATE), 数据如：91:1;92:1
func (s *cartServiceImpl) wsParseCartPostedData(skuData string) (arr []*cart.ItemPair) {
	arr = []*cart.ItemPair{}
	splitArr := strings.Split(skuData, ";")
	for _, str := range splitArr {
		i := strings.Index(str, ":")
		if i == -1 {
			continue
		}
		skuId, err := util.I64Err(strconv.Atoi(str[:i]))
		quantity, err1 := util.I32Err(strconv.Atoi(str[i+1:]))
		if err == nil && err1 == nil {
			arr = append(arr, &cart.ItemPair{
				SkuId:    skuId,
				Quantity: quantity,
			})
		}
	}
	return arr
}

// 生成结算提交的数据(立即购买),skuData参考数据：skuId:num;skuId2;num2
func (s *cartServiceImpl) createCheckedData(itemId int64, arr []*cart.ItemPair) string {
	buf := bytes.NewBufferString("{\"")
	buf.WriteString(strconv.Itoa(int(itemId)))
	buf.WriteString("\":[")
	for i, v := range arr {
		if i != 0 {
			buf.WriteString(",")
		}
		buf.WriteString("\"")
		buf.WriteString(strconv.Itoa(int(v.SkuId)))
		buf.WriteString("\"")
	}
	buf.WriteString("]}")
	return buf.String()
}

// 放入商品，data["Data"]
func (s *cartServiceImpl) wsPutItem(c cart.ICart, data map[string]string) (*proto.Result, error) {
	aId := c.GetAggregateRootId()
	itemId, err := util.I64Err(strconv.Atoi(data["ItemId"]))
	arr := s.wsParseCartPostedData(data["Data"])
	for _, v := range arr {
		err = c.Put(itemId, v.SkuId, v.Quantity)
		if err != nil {
			break
		}
	}
	if err == nil {
		_, err = c.Save()
		if err == nil {
			mp := make(map[string]string)
			mp["cartId"] = strconv.Itoa(int(aId))
			mp["checked"] = s.createCheckedData(itemId, arr)
			return s.success(mp), nil
		}
	}
	return s.result(err), nil
}

func (s *cartServiceImpl) wsUpdateItem(c cart.ICart, data map[string]string) (*proto.Result, error) {
	itemId, err := util.I64Err(strconv.Atoi(data["ItemId"]))
	arr := s.wsParseCartPostedData(data["Data"])
	for _, v := range arr {
		err = c.Update(itemId, v.SkuId, v.Quantity)
		if err != nil {
			break
		}
	}
	if err == nil {
		_, err = c.Save()
	}
	return s.result(err), nil
}

// 勾选购物车，格式如：1:2;1:5
func (s *cartServiceImpl) wsCheckCart(c cart.ICart, data map[string]string) (*proto.Result, error) {
	checked := data["Checked"]
	var arr []*cart.ItemPair
	splitArr := strings.Split(checked, ";")
	for _, str := range splitArr {
		i := strings.Index(str, ":")
		if i == -1 {
			continue
		}
		itemId, err := util.I64Err(strconv.Atoi(str[:i]))
		skuId, err1 := util.I64Err(strconv.Atoi(str[i+1:]))
		if err == nil && err1 == nil {
			arr = append(arr, &cart.ItemPair{
				ItemId: itemId,
				SkuId:  skuId,
			})
		}
	}
	err := c.SignItemChecked(arr)
	return s.result(err), nil
}

/*---------------- 普通购物车 ----------------*/

// 普通购物车接口
func (s *cartServiceImpl) NormalCartV1(_ context.Context, r *proto.NormalCartRequest) (*proto.Result, error) {
	//todo: check member
	c := s.cartRepo.GetMyCart(r.MemberId, cart.KWholesale)
	if r.Data == nil {
		r.Data = map[string]string{}
	}
	switch r.Action {
	case "GET":
		return s.wsGetCart(c, r.Data)
	case "MINI":
		return s.wsGetSimpleCart(c, r.Data)
	case "PUT":
		return s.wsPutItem(c, r.Data)
	case "UPDATE":
		return s.wsUpdateItem(c, r.Data)
	case "CHECK":
		return s.wsCheckCart(c, r.Data)
	}
	return s.result(errors.New("unknown action")), nil
}

//  获取购物车
func (s *cartServiceImpl) getShoppingCart(buyerId int64, code string) cart.ICart {
	var c cart.ICart
	var cc cart.ICart
	if len(code) > 0 {
		cc = s.cartRepo.GetShoppingCartByKey(code)
	}
	// 如果传入会员编号，则合并购物车
	if buyerId > 0 {
		c = s.cartRepo.GetMyCart(buyerId, cart.KNormal)
		if cc != nil {
			rc := c.(cart.INormalCart)
			rc.Combine(cc)
			c.Save()
		}
		return c
	}
	// 如果只传入code,且购物车存在，直接返回。
	if cc != nil {
		return cc
	}
	// 不存在，则新建购物车
	c = s.cartRepo.NewNormalCart(code)
	//_, err := c.Save()
	//domain.HandleError(err, "service")
	return c
}

// 获取购物车,当购物车编号不存在时,将返回一个新的购物车
func (s *cartServiceImpl) GetShoppingCart(memberId int64,
	cartCode string) *proto.SShoppingCart {
	c := s.getShoppingCart(memberId, cartCode)
	return s.parseCart(c)
}

// 转换购物车数据
func (s *cartServiceImpl) parseCart(c cart.ICart) *proto.SShoppingCart {
	dto := parser.ParseToDtoCart(c)
	for _, v := range dto.Shops {
		is := s.shopRepo.GetOnlineShop(int(v.ShopId))
		if is != nil {
			v.ShopName = is.ShopName
		} else {
			for _, it := range v.Items {
				c.Remove(it.ItemId, it.SkuId, it.Quantity)
			}
		}
	}
	return dto
}

// 放入购物车
func (s *cartServiceImpl) PutInCart(memberId int64, code string,
	itemId, skuId int64, quantity int32) (*proto.SShoppingCartItem, error) {
	c := s.getShoppingCart(memberId, code)
	if c == nil {
		return nil, cart.ErrNoSuchCart
	}
	err := c.Put(itemId, skuId, quantity)
	if err == nil {
		if _, err = c.Save(); err == nil {
			rc := c.(cart.INormalCart)
			item := rc.GetItem(itemId, skuId)
			return parser.ParseCartItem(item), err
		}
	}
	return nil, err
}
func (s *cartServiceImpl) SubCartItem(memberId int64, code string,
	itemId, skuId int64, quantity int32) error {
	c := s.getShoppingCart(memberId, code)
	if c == nil {
		return cart.ErrNoSuchCart
	}
	err := c.Remove(itemId, skuId, quantity)
	if err == nil {
		_, err = c.Save()
	}
	return err
}

// 勾选商品结算
func (s *cartServiceImpl) CartCheckSign(memberId int64,
	cartCode string, arr []*proto.SShoppingCartItem) error {
	c := s.getShoppingCart(memberId, cartCode)
	items := make([]*cart.ItemPair, len(arr))
	for i, v := range arr {
		items[i] = &cart.ItemPair{
			ItemId:  v.ItemId,
			SkuId:   v.SkuId,
			Checked: 1,
		}
	}
	err := c.SignItemChecked(items)
	if err == nil {
		_, err = c.Save()
	}
	return err
}

func (s *cartServiceImpl) SetBuyerAddress(buyerId int64, cartCode string, addressId int64) error {
	cart := s.getShoppingCart(buyerId, cartCode)
	return cart.SetBuyerAddress(addressId)
}


func (s *cartServiceImpl) GetCartSettle(memberId int64,
	cartCode string) *dto.SettleMeta {
	cart := s.getShoppingCart(memberId, cartCode)
	sp, deliver, payOpt := cart.GetSettleData()
	st := new(dto.SettleMeta)
	st.PaymentOpt = payOpt
	if sp != nil {
		v := sp.GetValue()
		ols := sp.(shop.IOnlineShop)
		st.Shop = &dto.SettleShopMeta{
			Id:   v.Id,
			Name: v.Name,
			Tel:  ols.GetShopValue().Tel,
		}
	}

	if deliver != nil {
		v := deliver.GetValue()
		st.Deliver = &dto.SettleDeliverMeta{
			Id:             v.ID,
			ConsigneeName:  v.ConsigneeName,
			ConsigneePhone: v.ConsigneePhone,
			Address:        strings.Replace(v.Area, " ", "", -1) + v.DetailAddress,
		}
	}

	return st
}
