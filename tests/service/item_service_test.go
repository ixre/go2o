package service

import (
	"context"
	"testing"

	"github.com/ixre/go2o/core/service/impl"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/types/typeconv"
)

func TestPagingShopGoodsRequest(t *testing.T) {
	goods, err := impl.ItemService.GetShopPagedOnShelvesGoods(context.TODO(), &proto.PagingShopGoodsRequest{
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
