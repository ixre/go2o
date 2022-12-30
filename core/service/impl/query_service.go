package impl

import (
	"context"

	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/domain/interface/order"
	"github.com/ixre/go2o/core/dto"
	"github.com/ixre/go2o/core/infrastructure/format"
	"github.com/ixre/go2o/core/query"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/go2o/core/variable"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/storage"
	"github.com/ixre/gof/types/typeconv"
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
	shopQuery       *query.ShopQuery
	orderQuery      *query.OrderQuery
	memberQuery     *query.MemberQuery
	statisticsQuery *query.StatisticsQuery
	proto.UnimplementedQueryServiceServer
}

// SummaryStatistics implements proto.QueryServiceServer
func (q *queryService) SummaryStatistics(context.Context, *proto.SummaryStatisticsRequest) (*proto.SummaryStatisticsResponse, error) {
	s := q.statisticsQuery.QuerySummary()
	return &proto.SummaryStatisticsResponse{
		TotalMembers:                s.TotalMembers,
		TodayJoinMembers:            s.TodayJoinMembers,
		TodayLoginMembers:           s.TodayLoginMembers,
		TodayCreateOrders:           s.TodayCreateOrders,
		AwaitShipmentOrders:         s.AwaitShipmentOrders,
		AwaitReviewWithdrawRequests: s.AwaitReviewWithdrawRequests,
	}, nil
}

// MemberStatistics 获取会员的订单状态及其数量
func (q *queryService) MemberStatistics(_ context.Context, req *proto.MemberStatisticsRequest) (*proto.MemberStatisticsResponse, error) {
	ret := make(map[int32]int32, 0)
	for k, v := range q.memberQuery.OrdersQuantity(req.MemberId) {
		ret[int32(k)] = int32(v)
	}
	return &proto.MemberStatisticsResponse{
		AwaitPaymentOrders:  ret[int32(order.StatAwaitingPayment)],
		AwaitShipmentOrders: ret[int32(order.StatAwaitingShipment)],
		AwaitReceiveOrders:  ret[int32(order.StatShipped)],
		CompletedOrders:     ret[int32(order.StatCompleted)],
	}, nil
}

func NewQueryService(o orm.Orm, s storage.Interface) *queryService {
	shopQuery := query.NewShopQuery(o, s)
	return &queryService{
		shopQuery:       shopQuery,
		memberQuery:     query.NewMemberQuery(o),
		orderQuery:      query.NewOrderQuery(o),
		statisticsQuery: query.NewStatisticsQuery(o, s),
	}
}

// PagingShops 获取分页店铺数据
func (q *queryService) PagingShops(_ context.Context, r *proto.QueryPagingShopRequest) (*proto.QueryPagingShopsResponse, error) {
	n, rows := q.shopQuery.PagedOnBusinessOnlineShops(
		int(r.Params.Begin),
		int(r.Params.End),
		"", r.Params.SortBy)
	ret := &proto.QueryPagingShopsResponse{
		Total: int64(n),
		Value: make([]*proto.QueryPagingShop, len(rows)),
	}
	if len(rows) > 0 {
		for i, v := range rows {
			v.Logo = format.GetResUrl(v.Logo)
			if v.Host == "" {
				v.Host = v.Alias + "." + variable.Domain
			}
			ret.Value[i] = &proto.QueryPagingShop{
				Id:         v.Id,
				Name:       v.Name,
				Alias:      v.Alias,
				Host:       v.Host,
				Logo:       v.Logo,
				CreateTime: v.CreateTime,
			}
		}
	}
	return ret, nil
}

// MemberNormalOrders 查询分页普通订单
func (q *queryService) MemberNormalOrders(_ context.Context, r *proto.MemberOrderPagingRequest) (*proto.MemberOrderPagingResponse, error) {
	n, list := q.orderQuery.QueryPagingNormalOrder(
		r.MemberId,
		r.Params.Begin,
		r.Params.End,
		true,
		r.Params.Where,
		r.Params.SortBy)
	ret := &proto.MemberOrderPagingResponse{
		Total: int64(n),
		Value: make([]*proto.SMemberPagingOrder, len(list)),
	}
	for i, v := range list {
		ret.Value[i] = q.parseOrder(v)
	}
	return ret, nil
}

// QueryWholesaleOrders 查询分页批发订单
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
		Value: make([]*proto.SMemberPagingOrder, len(list)),
	}
	for i, v := range list {
		ret.Value[i] = q.parseOrder(v)
	}
	return ret, nil
}

// QueryTradeOrders 查询分页交易/服务类订单
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
		Value: make([]*proto.SMemberPagingOrder, len(list)),
	}
	for i, v := range list {
		ret.Value[i] = q.parseTradeOrder(v)
	}
	return ret, nil
}

func (q *queryService) parseOrder(src *dto.MemberPagingOrderDto) *proto.SMemberPagingOrder {
	dst := &proto.SMemberPagingOrder{
		OrderNo:        src.OrderNo,
		BuyerId:        src.BuyerId,
		BuyerUser:      src.BuyerUser,
		ShopId:         src.ShopId,
		ShopName:       src.ShopName,
		ItemCount:      int64(src.ItemCount),
		ItemAmount:     src.ItemAmount,
		DiscountAmount: src.DiscountAmount,
		DeductAmount:   src.DeductAmount,
		ExpressFee:     src.ExpressFee,
		PackageFee:     src.PackageFee,
		FinalAmount:    src.FinalAmount,
		Items:          []*proto.SOrderItem{},
		Status:         int32(src.Status),
		StatusText:     src.StatusText,
		CreateTime:     src.CreateTime,
	}

	for _, it := range src.Items {
		dst.Items = append(dst.Items, &proto.SOrderItem{
			Id:             int64(it.Id),
			SnapshotId:     int64(it.SnapshotId),
			SkuId:          int64(it.SkuId),
			ItemId:         int64(it.ItemId),
			ItemTitle:      it.ItemTitle,
			Image:          it.Image,
			Price:          it.Price,
			FinalPrice:     it.FinalPrice,
			Quantity:       int32(it.Quantity),
			ReturnQuantity: int32(it.ReturnQuantity),
			ItemAmount:     int64(it.FinalPrice * int64(it.Quantity)),
			FinalAmount:    it.FinalAmount,
			IsShipped:      it.IsShipped == 1,
			Data:           map[string]string{},
		})
	}
	return dst
}

func (q *queryService) parseOrderItem(v *dto.OrderItem) *proto.SOrderItem {
	return &proto.SOrderItem{
		Id:             int64(v.Id),
		SnapshotId:     int64(v.SnapshotId),
		SkuId:          int64(v.SkuId),
		ItemId:         int64(v.ItemId),
		ItemTitle:      v.ItemTitle,
		Image:          v.Image,
		Price:          v.Price,
		FinalPrice:     v.FinalPrice,
		Quantity:       int32(v.Quantity),
		ReturnQuantity: int32(v.ReturnQuantity),
		ItemAmount:     v.Amount,
		FinalAmount:    v.FinalAmount,
		IsShipped:      v.IsShipped == 1,
		Data:           nil,
	}
}

func (q *queryService) parseTradeOrder(src *proto.SSingleOrder) *proto.SMemberPagingOrder {
	return &proto.SMemberPagingOrder{
		OrderNo: src.OrderNo,
		//ShopName:       src.,
		ItemAmount:     src.ItemAmount,
		DiscountAmount: src.DiscountAmount,
		ExpressFee:     src.ExpressFee,
		PackageFee:     src.PackageFee,
		//IsPaid:         src.IsPaid,
		FinalAmount: src.FinalAmount,
		Status:      int32(src.Status),
		//StateText:      src.StateText,
		CreateTime: src.SubmitTime,
		//Items:      make([]*proto.SOrderItem, 0),
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
			Nickname:      v.Name,
			Portrait:      v.Avatar,
			Level:         v.Level,
			Integral:      v.Integral,
			Balance:       v.Balance,
			WalletBalance: v.WalletBalance,
		}
	}
	return rsp, nil
}

// SearchMembers 根据用户或手机筛选会员
func (q *queryService) SearchMembers(_ context.Context, r *proto.MemberSearchRequest) (*proto.MemberListResponse, error) {
	list := q.memberQuery.FilterMemberByUserOrPhone(r.Keyword)
	ret := &proto.MemberListResponse{
		Value: make([]*proto.MemberListSingle, len(list)),
	}
	for i, v := range list {
		ret.Value[i] = &proto.MemberListSingle{
			MemberId: int64(v.Id),
			User:     v.User,
			Nickname: v.Name,
			Portrait: v.Avatar,
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

// QueryMemberFavoriteGoods 获取分页商品收藏
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

// 获取钱包账户分页记录
func (q *queryService) PagingMemberAccountLog(_ context.Context, r *proto.PagingAccountInfoRequest) (*proto.SPagingResult, error) {
	var total int
	var rows []map[string]interface{}
	switch member.AccountType(r.AccountType) {
	case member.AccountIntegral:
		total, rows = q.memberQuery.PagedIntegralAccountLog(
			r.MemberId, r.Params.Begin,
			r.Params.End, r.Params.SortBy)
	case member.AccountBalance:
		total, rows = q.memberQuery.PagedBalanceAccountLog(
			r.MemberId, int(r.Params.Begin),
			int(r.Params.End), r.Params.Where,
			r.Params.SortBy)
	case member.AccountWallet:
		total, rows = q.memberQuery.PagedWalletAccountLog(
			r.MemberId, int(r.Params.Begin),
			int(r.Params.End), r.Params.Where,
			r.Params.Where)
	}
	rs := &proto.SPagingResult{
		ErrCode: 0,
		ErrMsg:  "",
		Count:   int32(total),
		Data:    typeconv.MustJson(rows),
	}
	return rs, nil
}
