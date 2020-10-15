package impl

import (
	"context"
	"github.com/ixre/gof"
	"go2o/core/dto"
	"go2o/core/infrastructure/format"
	"go2o/core/query"
	"go2o/core/service/proto"
	"go2o/core/variable"
)

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : query_service.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-09-13 13:09
 * description :
 * history :
 */

var _ proto.QueryServiceServer = new(queryService)

type queryService struct {
	shopQuery   *query.ShopQuery
	orderQuery  *query.OrderQuery
	memberQuery *query.MemberQuery
}

func NewQueryService() *queryService {
	ctx := gof.CurrentApp
	shopQuery := query.NewShopQuery(ctx)
	return &queryService{
		shopQuery:   shopQuery,
		memberQuery: query.NewMemberQuery(ctx.Db()),
		orderQuery:  query.NewOrderQuery(ctx.Db()),
	}
}

// 获取分页店铺数据
func (q *queryService) PagingShops(_ context.Context, r *proto.QueryPagingShopRequest) (*proto.QueryPagingShopsResponse, error) {
	n, rows := q.shopQuery.PagedOnBusinessOnlineShops(
		int(r.Params.Begin),
		int(r.Params.End),
		"", r.Params.SortBy)
	ret := &proto.QueryPagingShopsResponse{
		Total: int64(n),
		Value: make([]*proto.QueryPagingShop, n),
	}
	if len(rows) > 0 {
		for _, v := range rows {
			v.Logo = format.GetResUrl(v.Logo)
			if v.Host == "" {
				v.Host = v.Alias + "." + variable.Domain
			}
			ret.Value = append(ret.Value, &proto.QueryPagingShop{
				Id:         v.Id,
				Name:       v.Name,
				Alias:      v.Alias,
				Host:       v.Host,
				Logo:       v.Logo,
				CreateTime: v.CreateTime,
			})
		}
	}
	return ret, nil
}

// 查询分页普通订单
func (q *queryService) MemberNormalOrders(_ context.Context, r *proto.MemberOrderPagingRequest) (*proto.MemberOrderPagingResponse, error) {
	n, list := q.orderQuery.QueryPagerOrder(
		r.MemberId,
		r.Params.Begin,
		r.Params.End,
		true,
		r.Params.Where,
		r.Params.SortBy)
	ret := &proto.MemberOrderPagingResponse{
		Total: int64(n),
		Value: make([]*proto.PagedMemberSubOrder, n),
	}
	for i, v := range list {
		ret.Value[i] = q.parseOrder(v)
	}
	return ret, nil
}

// 查询分页批发订单
func (q *queryService) QueryWholesaleOrders(_ context.Context, r *proto.MemberOrderPagingRequest) (*proto.MemberOrderPagingResponse, error) {
	n, list := q.orderQuery.PagedWholesaleOrderOfBuyer(
		r.MemberId,
		r.Params.Begin,
		r.Params.End,
		true,
		r.Params.Where,
		r.Params.SortBy)
	ret := &proto.MemberOrderPagingResponse{
		Total: int64(n),
		Value: make([]*proto.PagedMemberSubOrder, n),
	}
	for i, v := range list {
		ret.Value[i] = q.parseOrder(v)
	}
	return ret, nil
}

// 查询分页交易/服务类订单
func (q *queryService) QueryTradeOrders(_ context.Context, r *proto.MemberOrderPagingRequest) (*proto.MemberOrderPagingResponse, error) {
	n, list := q.orderQuery.PagedTradeOrderOfBuyer(
		r.MemberId,
		r.Params.Begin,
		r.Params.End,
		true,
		r.Params.Where,
		r.Params.SortBy)
	ret := &proto.MemberOrderPagingResponse{
		Total: int64(n),
		Value: make([]*proto.PagedMemberSubOrder, n),
	}
	for i, v := range list {
		ret.Value[i] = q.parseTradeOrder(v)
	}
	return ret, nil
}

func (q *queryService) parseOrder(src *dto.PagedMemberSubOrder) *proto.PagedMemberSubOrder {
	dst := &proto.PagedMemberSubOrder{
		OrderId:        src.Id,
		OrderNo:        src.OrderNo,
		ParentNo:       src.ParentNo,
		VendorId:       src.VendorId,
		ShopId:         src.ShopId,
		ShopName:       src.ShopName,
		ItemAmount:     float64(src.ItemAmount),
		DiscountAmount: float64(src.DiscountAmount),
		ExpressFee:     float64(src.ExpressFee),
		PackageFee:     float64(src.PackageFee),
		IsPaid:         src.IsPaid,
		FinalAmount:    float64(src.FinalAmount),
		State:          int32(src.State),
		StateText:      src.StateText,
		CreateTime:     src.CreateTime,
		Items:          make([]*proto.SOrderItem, 0),
	}
	for _, v := range src.Items {
		dst.Items = append(dst.Items, q.parseOrderItem(v))
	}
	return dst
}

func (q *queryService) parseOrderItem(v *dto.OrderItem) *proto.SOrderItem {
	return &proto.SOrderItem{
		Id:             int64(v.Id),
		SnapshotId:     int64(v.SnapshotId),
		SkuId:          int64(v.SkuId),
		ItemId:         int64(v.ItemId),
		ItemTitle:      v.GoodsTitle,
		Image:          v.Image,
		Price:          float64(v.Price),
		FinalPrice:     float64(v.FinalPrice),
		Quantity:       int32(v.Quantity),
		ReturnQuantity: int32(v.ReturnQuantity),
		Amount:         float64(v.Amount),
		FinalAmount:    float64(v.FinalAmount),
		IsShipped:      v.IsShipped == 1,
		Data:           nil,
	}
}

func (q *queryService) parseTradeOrder(src *proto.SSingleOrder) *proto.PagedMemberSubOrder {
	return &proto.PagedMemberSubOrder{
		OrderId:  src.OrderId,
		OrderNo:  src.OrderNo,
		VendorId: src.SellerId,
		ShopId:   src.ShopId,
		//ShopName:       src.,
		ItemAmount:     src.ItemAmount,
		DiscountAmount: src.DiscountAmount,
		ExpressFee:     src.ExpressFee,
		PackageFee:     src.PackageFee,
		//IsPaid:         src.IsPaid,
		FinalAmount: src.FinalAmount,
		State:       int32(src.State),
		//StateText:      src.StateText,
		CreateTime: src.SubmitTime,
		Items:      make([]*proto.SOrderItem, 0),
	}
}

func (q *queryService) QueryMemberList(_ context.Context, r *proto.MemberListRequest) (*proto.MemberListResponse, error) {
	list := q.memberQuery.QueryMemberList(r.IdList)
	var rsp = &proto.MemberListResponse{
		Value: make([]*proto.MemberListSingle, len(list)),
	}
	for i, v := range list {
		v.Avatar = format.GetResUrl(v.Avatar)
		rsp.Value[i] = &proto.MemberListSingle{
			MemberId:      int64(v.MemberId),
			User:          v.Usr,
			NickName:      v.Name,
			Avatar:        v.Avatar,
			Level:         v.Level,
			Integral:      v.Integral,
			Balance:       float64(v.Balance),
			WalletBalance: float64(v.WalletBalance),
		}
	}
	return rsp, nil
}

// 根据用户或手机筛选会员
func (q *queryService) SearchMembers(_ context.Context, r *proto.MemberSearchRequest) (*proto.MemberListResponse, error) {
	list := q.memberQuery.FilterMemberByUserOrPhone(r.Keyword)
	ret := &proto.MemberListResponse{
		Value: make([]*proto.MemberListSingle, len(list)),
	}
	for i, v := range list {
		ret.Value[i] = &proto.MemberListSingle{
			MemberId: int64(v.Id),
			User:     v.User,
			NickName: v.Name,
			Avatar:   v.Avatar,
		}
	}
	return ret, nil
}

// 获取分页店铺收藏

func (q *queryService) QueryMemberFavoriteShops(_ context.Context, r *proto.FavoriteQueryRequest) (*proto.PagingShopFavoriteResponse, error) {
	total, rows := q.memberQuery.PagedShopFav(r.MemberId, int(r.Begin), int(r.End), r.Where)
	ret := &proto.PagingShopFavoriteResponse{
		Total: int64(total),
		Data:  make([]*proto.SPagingShopFavorite, len(rows)),
	}
	for i, v := range rows {
		ret.Data[i] = &proto.SPagingShopFavorite{
			Id:         int64(v.Id),
			ShopId:     int64(v.ShopId),
			ShopName:   v.ShopName,
			Logo:       v.Logo,
			UpdateTime: v.UpdateTime,
		}
	}
	return ret, nil
}

// 获取分页商品收藏
func (q *queryService) QueryMemberFavoriteGoods(_ context.Context, r *proto.FavoriteQueryRequest) (*proto.PagingGoodsFavoriteResponse, error) {
	total, rows := q.memberQuery.PagedGoodsFav(r.MemberId, int(r.Begin), int(r.End), r.Where)
	ret := &proto.PagingGoodsFavoriteResponse{
		Total: int64(total),
		Data:  make([]*proto.SPagingGoodsFavorite, len(rows)),
	}
	for i, v := range rows {
		ret.Data[i] = &proto.SPagingGoodsFavorite{
			Id:         int64(v.Id),
			SkuId:      int64(v.SkuId),
			GoodsName:  v.GoodsName,
			Image:      v.Image,
			OnShelves:  v.OnShelves == 1,
			StockNum:   int32(v.StockNum),
			SalePrice:  v.SalePrice,
			UpdateTime: v.UpdateTime,
		}
	}
	return ret, nil
}
