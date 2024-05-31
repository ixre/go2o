package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/ixre/go2o/core/domain/interface/order"
	"github.com/ixre/go2o/core/inject"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/types/typeconv"
)

func Test_merchantService_PagedNormalOrderOfVendor(t *testing.T) {
	ret, _ := inject.GetMerchantService().PagedNormalOrderOfVendor(context.TODO(), &proto.MerchantOrderRequest{
		MerchantId: 1,
		Pagination: false,
		Params: &proto.SPagingParams{
			Begin: 0,
			End:   10,
			Where: fmt.Sprintf("o.status = %d", order.StatAwaitingPayment),
		},
	})
	t.Log(typeconv.MustJson(ret))
}
