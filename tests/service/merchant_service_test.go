package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/ixre/go2o/pkg/inject"
	"github.com/ixre/go2o/pkg/interface/domain/order"
	"github.com/ixre/go2o/pkg/service/proto"
	"github.com/ixre/gof/typeconv"
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
