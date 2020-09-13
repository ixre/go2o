package impl

import (
	"context"
	"github.com/ixre/gof"
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
type queryService struct{
	shopQuery *query.ShopQuery
}

func NewQueryService()*queryService{
	ctx := gof.CurrentApp
	shopQuery:= query.NewShopQuery(ctx)
	return &queryService{
		shopQuery: shopQuery,
	}
}

// 获取分页店铺数据
func (q *queryService) PagingShops(_ context.Context, r *proto.QueryPagingShopRequest) (*proto.QueryPagingShopsResponse, error) {
	n, rows := q.shopQuery.PagedOnBusinessOnlineShops(
		int(r.Params.Begin),
		int(r.Params.Over),
		"", r.Params.SortBy)
	ret := &proto.QueryPagingShopsResponse{
		Total: int32(n),
		Value: make([]*proto.QueryPagingShop, n),
	}
	if len(rows) > 0 {
		for _, v := range rows {
			v.Logo = format.GetResUrl(v.Logo)
			if v.Host == "" {
				v.Host = v.Alias + "." + variable.Domain
			}
			ret.Value = append(ret.Value,&proto.QueryPagingShop{
				Id:                   v.Id,
				Name:                 v.Name,
				Alias:                v.Alias,
				Host:                 v.Host,
				Logo:                 v.Logo,
				CreateTime:           v.CreateTime,
			})
		}
	}
	return ret,nil
}

