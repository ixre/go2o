package service

import (
	"context"
	"testing"

	"github.com/ixre/go2o/core/domain/interface/item"
	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/inject"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/types/typeconv"
)

func TestPagingMemberAccountLog(t *testing.T) {
	s := inject.GetQueryService()
	ret, _ := s.PagingMemberAccountLog(context.TODO(), &proto.PagingAccountLogRequest{
		MemberId:    723,
		ValueFilter: 0,
		AccountType: int32(member.AccountWallet),
		Params: &proto.SPagingParams{
			Begin:  0,
			End:    10,
			SortBy: "create_time DESC,id DESC",
		},
	})
	t.Log(typeconv.MustJson(ret))
}

func TestPagingIntegralLog(t *testing.T) {
	params := &proto.SPagingParams{
		Parameters: nil,
		SortBy:     "",
		Begin:      0,
		End:        10,
	}
	r, _ := inject.GetQueryService().PagingMemberAccountLog(context.TODO(),
		&proto.PagingAccountLogRequest{
			MemberId:    1,
			AccountType: int32(member.AccountIntegral),
			Params:      params,
		})
	t.Log(typeconv.MustJson(r))
}

func TestPagingBalanceLog(t *testing.T) {
	memberId := 702
	params := &proto.SPagingParams{
		Parameters: nil,
		SortBy:     "",
		Begin:      0,
		End:        10,
	}
	r, _ := inject.GetQueryService().PagingMemberAccountLog(context.TODO(),
		&proto.PagingAccountLogRequest{
			MemberId:    int64(memberId),
			AccountType: int32(member.AccountBalance),
			Params:      params,
		})
	t.Log(typeconv.MustJson(r))
}

func TestQueryPagingFlagGoods(t *testing.T) {
	params := &proto.SPagingParams{
		Parameters: nil,
		SortBy:     "",
		Begin:      0,
		End:        10,
	}
	r, _ := inject.GetQueryService().PagingOnShelvesGoods(context.TODO(),
		&proto.PagingShopGoodsRequest{
			ShopId:     0,
			CategoryId: 0,
			Flag:       item.FlagNewOnShelve,
			Params:     params,
		})
	t.Log(typeconv.MustJson(r))
}

func TestPagingShopGoodsRequest(t *testing.T) {
	goods, err := inject.GetQueryService().PagingOnShelvesGoods(context.TODO(), &proto.PagingShopGoodsRequest{
		ShopId:     0,
		CategoryId: 2185,
		Params: &proto.SPagingParams{
			Begin:  0,
			End:    20,
			Where:  "",
			SortBy: "item_info.sale_num DESC",
		},
	})

	if err != nil {
		t.Error(err)
	} else {
		t.Log(len(goods.Data), typeconv.MustJson(goods.Data))
	}
}

func TestMemberStatifics(t *testing.T) {
	var memberId int64 = 729
	mp, _ := inject.GetQueryService().MemberStatistics(context.TODO(), &proto.MemberStatisticsRequest{
		MemberId: memberId,
	})
	t.Log("未支付订单数", mp.AwaitPaymentOrders)
}

func TestQuerySearchItem(t *testing.T) {
	list, err := inject.GetQueryService().SearchItem(context.TODO(), &proto.SearchItemRequest{
		ShopId:  0,
		Keyword: "1",
		//CategoryId: 0,
		//Begin:      0,
		Size: 10,
	})

	if err != nil {
		t.Error(err)
	} else {
		t.Log(len(list.Value), typeconv.MustJson(list.Value))
	}
}
