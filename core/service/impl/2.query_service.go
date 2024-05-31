package impl

import (
	"context"
	"strings"

	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/domain/interface/order"
	"github.com/ixre/go2o/core/domain/interface/product"
	"github.com/ixre/go2o/core/domain/interface/valueobject"
	"github.com/ixre/go2o/core/dto"
	"github.com/ixre/go2o/core/infrastructure/format"
	"github.com/ixre/go2o/core/query"
	"github.com/ixre/go2o/core/service/parser"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/go2o/core/variable"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/storage"
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
	itemQuery       *query.ItemQuery
	statisticsQuery *query.StatisticsQuery
	catRepo         product.ICategoryRepo
	proto.UnimplementedQueryServiceServer
}

func NewQueryService(o orm.Orm, s storage.Interface,
	catRepo product.ICategoryRepo) proto.QueryServiceServer {
	shopQuery := query.NewShopQuery(o, s)
	return &queryService{
		shopQuery:       shopQuery,
		itemQuery:       query.NewItemQuery(o),
		memberQuery:     query.NewMemberQuery(o),
		orderQuery:      query.NewOrderQuery(o),
		catRepo:         catRepo,
		statisticsQuery: query.NewStatisticsQuery(o, s),
	}
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
		AwaitShipmentOrders: ret[int32(order.StatAwaitingPickup)] + ret[int32(order.StatAwaitingShipment)],
		AwaitReceiveOrders:  ret[int32(order.StatShipped)],
		CompletedOrders:     ret[int32(order.StatCompleted)],
	}, nil
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
			v.Logo = format.GetFileFullUrl(v.Logo)
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
	n, list := q.orderQuery.PagingTradeOrderOfBuyer(
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
			SpecWord:       it.SpecWord,
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
		v.Avatar = format.GetFileFullUrl(v.Avatar)
		rsp.Value[i] = &proto.MemberListSingle{
			MemberId:      int64(v.MemberId),
			Username:      v.Usr,
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
			Username: v.User,
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
func (q *queryService) PagingMemberAccountLog(_ context.Context, r *proto.PagingAccountLogRequest) (*proto.MemberAccountPagingLogResponse, error) {
	var total int
	var rows []*proto.SMemberAccountLog
	switch member.AccountType(r.AccountType) {
	case member.AccountIntegral:
		total, rows = q.memberQuery.PagedIntegralAccountLog(
			r.MemberId, r.ValueFilter, r.Params.Begin,
			r.Params.End, r.Params.SortBy)
	case member.AccountBalance:
		total, rows = q.memberQuery.PagedBalanceAccountLog(
			r.MemberId, r.ValueFilter, int(r.Params.Begin),
			int(r.Params.End), r.Params.Where,
			r.Params.SortBy)
	case member.AccountWallet:
		total, rows = q.memberQuery.PagedWalletAccountLog(
			r.MemberId, r.ValueFilter, int(r.Params.Begin),
			int(r.Params.End), r.Params.Where,
			r.Params.SortBy)
	}
	rs := &proto.MemberAccountPagingLogResponse{
		Total: int32(total),
		Data:  rows,
	}
	return rs, nil
}

// PagingOnShelvesGoods 获取分页上架的商品
func (q *queryService) PagingOnShelvesGoods(_ context.Context, r *proto.PagingShopGoodsRequest) (*proto.PagingShopGoodsResponse, error) {
	ret := &proto.PagingShopGoodsResponse{
		Total: 0,
		Data:  make([]*proto.SGoods, 0),
	}
	var ids []int
	if r.CategoryId > 0 {
		cat := q.catRepo.GlobCatService().GetCategory(int(r.CategoryId))
		if cat == nil {
			return ret, nil
		}
		ids = cat.GetChildes()
		ids = append(ids, int(r.CategoryId))

	}
	if len(strings.TrimSpace(r.Params.SortBy)) == 0 {
		r.Params.SortBy = "item_snapshot.update_time DESC"
	}
	var total int
	var list []*valueobject.Goods
	switch r.ItemType {
	case proto.EItemSalesType_IT_NORMAL:
		total, list = q.itemQuery.GetPagingOnShelvesGoods(
			r.ShopId, ids, int(r.Flag),
			int(r.Params.Begin),
			int(r.Params.End),
			r.Keyword,
			r.Params.Where,
			r.Params.SortBy)
	case proto.EItemSalesType_IT_WHOLESALE:
		// total, list = q.getPagedOnShelvesItemForWholesale(
		// 	int32(r.CategoryId),
		// 	int32(r.Params.Begin),
		// 	int32(r.Params.End),
		// 	where,
		// 	r.Params.SortBy)
	}

	ret.Total = int64(total)
	for _, v := range list {
		v.Image = format.GetGoodsImageUrl(v.Image)
		ret.Data = append(ret.Data, q.parseGoods(v))
	}
	return ret, nil
}

func (q *queryService) parseGoods(v *valueobject.Goods) *proto.SGoods {
	return &proto.SGoods{
		ItemId:      v.ItemId,
		ProductId:   v.ProductId,
		VendorId:    int64(v.VendorId),
		ShopId:      int64(v.ShopId),
		CategoryId:  v.CategoryId,
		Title:       v.Title,
		ShortTitle:  v.ShortTitle,
		GoodsNo:     v.GoodsNo,
		Image:       v.Image,
		OriginPrice: v.OriginPrice,
		Price:       v.Price,
		PromPrice:   v.PromPrice,
		PriceRange:  v.PriceRange,
		SkuId:       v.SkuId,
		ItemFlag:    int32(v.ItemFlag),
		StockNum:    v.StockNum,
		SaleNum:     v.SaleNum,
	}
}

// QueryItemSalesHistory 查询商品销售记录
func (q *queryService) QueryItemSalesHistory(_ context.Context, req *proto.QueryItemSalesHistoryRequest) (*proto.QueryItemSalesHistoryResponse, error) {
	list := q.itemQuery.QueryItemSalesHistory(req.ItemId, int(req.Size), req.Random)
	ret := &proto.QueryItemSalesHistoryResponse{
		Value: []*proto.SItemSalesHistory{},
	}
	for _, v := range list {
		dst := &proto.SItemSalesHistory{
			BuyerUserCode:   v.BuyerUserCode,
			BuyerName:       v.BuyerName,
			BuyerPortrait:   v.BuyerPortrait,
			BuyTime:         v.BuyTime,
			IsFinishPayment: v.OrderState > order.StatAwaitingPayment,
		}
		if req.MaskBuyer {
			dst.BuyerName = format.MaskNickname(dst.BuyerName)
		}
		ret.Value = append(ret.Value, dst)
	}
	return ret, nil
}

// SearchItem 搜索商品
func (q *queryService) SearchItem(_ context.Context, req *proto.SearchItemRequest) (*proto.SearchItemResponse, error) {
	list := q.itemQuery.SearchItem(int(req.ShopId), req.Keyword, int(req.Size))
	ret := &proto.SearchItemResponse{
		Value: []*proto.SSearchItemResult{},
	}
	for _, v := range list {
		dst := &proto.SSearchItemResult{
			ItemId:     v.ItemId,
			ItemFlag:   int32(v.ItemFlag),
			Code:       v.Code,
			SellerId:   v.SellerId,
			Title:      v.Title,
			Image:      v.Image,
			PriceRange: v.PriceRange,
			StockNum:   v.StockNum,
		}
		ret.Value = append(ret.Value, dst)
	}
	return ret, nil
}

func (q *queryService) getPagedOnShelvesItemForWholesale(catId int32, start,
	end int32, where, sortBy string) (int32, []*proto.SUnifiedViewItem) {

	total, list := q.itemQuery.GetPagingOnShelvesItemForWholesale(catId,
		start, end, where, sortBy)
	arr := make([]*proto.SUnifiedViewItem, len(list))
	for j, v := range list {
		v.Image = format.GetGoodsImageUrl(v.Image)
		dto := parser.ItemDtoV2(*v)
		q.attachWholesaleItemDataV2(dto)
		arr[j] = dto
	}
	return total, arr
}

// 附加批发商品的信息
func (q *queryService) attachWholesaleItemDataV2(dto *proto.SUnifiedViewItem) {
	dto.Data = make(map[string]string)
	// vendor := q.mchRepo.GetMerchant(int(dto.VendorId))
	// if vendor != nil {
	// 	vv := vendor.GetValue()
	// 	pStr := q.valueRepo.GetAreaName(int32(vv.Province))
	// 	cStr := q.valueRepo.GetAreaName(int32(vv.City))
	// 	dto.Data["VendorName"] = vv.CompanyName
	// 	dto.Data["ShipArea"] = pStr + cStr
	// 	// 认证信息
	// 	ei := vendor.ProfileManager().GetEnterpriseInfo()
	// 	if ei != nil && ei.Reviewed == enum.ReviewPass {
	// 		dto.Data["Authorized"] = "true"
	// 	} else {
	// 		dto.Data["Authorized"] = "false"
	// 	}
	// 	// 品牌
	// 	b := q.promRepo.BrandService().Get(int(dto.BrandId))
	// 	if b != nil {
	// 		dto.Data["BrandName"] = b.Name
	// 		dto.Data["BrandImage"] = b.Image
	// 		dto.Data["BrandId"] = strconv.Itoa(int(b.Id))
	// 	}
	// }
}
