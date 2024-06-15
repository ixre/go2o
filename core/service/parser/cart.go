package parser

import (
	"github.com/ixre/go2o/core/domain/interface/cart"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/math"
	"github.com/ixre/gof/types"
)

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : cart.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-09-08 08:29
 * description :
 * history :
 */

func ParseCartItem(item *cart.NormalCartItem) *proto.SShoppingCartItem {
	i := &proto.SShoppingCartItem{
		ItemId:   item.ItemId,
		SkuId:    item.SkuId,
		Quantity: item.Quantity,
		Checked:  item.Checked == 1,
		ShopId:   item.ShopId,
	}
	if item.Sku != nil {
		i.Image = item.Sku.Image
		i.OriginPrice = math.Round(float64(item.Sku.OriginPrice), 2)
		i.Price = math.Round(float64(item.Sku.Price), 2)
		i.SpecWord = item.Sku.SpecWord
		if i.Title == "" {
			i.Title = item.Sku.Title
		}
		i.Code = item.Sku.ItemCode
		i.StockText = types.StringCond(item.Sku.Stock > 0,
			"有货", "无货")
	}
	return i
}

func ParseToDtoCart(ic cart.ICartAggregateRoot) *proto.SShoppingCart {
	c := &proto.SShoppingCart{}
	if ic.Kind() != cart.KNormal {
		panic("购物车类型非零售")
	}
	rc := ic.(cart.INormalCart)
	v := rc.Value()

	c.CartId = ic.GetAggregateRootId()
	c.CartCode = v.CartCode
	c.Sellers = []*proto.SShoppingCartGroup{}

	items := rc.Items()
	if len(items) > 0 {
		mp := make(map[int64]*proto.SShoppingCartGroup, 0) //保存运营商到map
		for _, v := range items {
			vendor, ok := mp[v.ShopId]
			if !ok {
				vendor = &proto.SShoppingCartGroup{
					SellerId: v.VendorId,
					ShopId:   v.ShopId,
					Items:    []*proto.SShoppingCartItem{},
				}
				mp[v.ShopId] = vendor
				c.Sellers = append(c.Sellers, vendor)
			}
			if v.Checked == 1 {
				vendor.Checked = true
			}
			vendor.Items = append(vendor.Items, ParseCartItem(v))
		}
	}
	return c
}

func ParsePrepareOrderGroups(ic cart.ICartAggregateRoot) []*proto.SPrepareOrderGroup {
	arr := make([]*proto.SPrepareOrderGroup, 0)
	if ic.Kind() != cart.KNormal {
		panic("购物车类型非零售")
	}
	rc := ic.(cart.INormalCart)
	items := rc.Items()
	if items != nil && len(items) > 0 {
		mp := make(map[int64]*proto.SPrepareOrderGroup, 0) //保存运营商到map
		for _, v := range items {
			if v.Checked != 1 {
				continue
			}
			vendor, ok := mp[v.ShopId]
			if !ok {
				vendor = &proto.SPrepareOrderGroup{
					SellerId: v.VendorId,
					ShopId:   v.ShopId,
					Items:    []*proto.SPrepareOrderItem{},
				}
				mp[v.ShopId] = vendor
				arr = append(arr, vendor)
			}
			vendor.Items = append(vendor.Items, parsePrepareOrderItem(v))
		}
	}
	return arr
}

func parsePrepareOrderItem(item *cart.NormalCartItem) *proto.SPrepareOrderItem {
	i := &proto.SPrepareOrderItem{
		ItemId:   item.ItemId,
		SkuId:    item.SkuId,
		Quantity: item.Quantity,
	}
	if item.Sku != nil {
		i.Image = item.Sku.Image
		i.Price = math.Round(float64(item.Sku.Price), 2)
		i.SpecWord = item.Sku.SpecWord
		if i.Title == "" {
			i.Title = item.Sku.Title
		}
	}
	return i
}
