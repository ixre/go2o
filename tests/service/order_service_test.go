package service

import (
	"context"
	"strings"
	"testing"

	"github.com/ixre/go2o/core/domain/interface/order"
	"github.com/ixre/go2o/core/service/impl"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/types/typeconv"
)

// 测试提交普通订单
func TestSubmitNormalOrder(t *testing.T) {
	var memberId int64 = 1
	ret, err := impl.OrderService.SubmitOrder(
		context.TODO(),
		&proto.SubmitOrderRequest{
			BuyerId:       memberId,
			OrderType:     int32(order.TRetail),
			AddressId:     1,
			Subject:       "",
			CouponCode:    "",
			BalanceDeduct: false,
		})
	if err != nil {
		t.Error(err)
	}
	if ret.ErrCode > 0 {
		t.Error(ret.ErrMsg)
	}
}

func TestSubmitOrderSubjectPostgresInsert(t *testing.T) {
	s := "指间陶艺 精品宜兴紫砂茶宠 名家全手工小号座镇貔貅茶玩雕塑摆件精品 威震八方(貔貅) 公款"
	s2 := strings.Replace(s, " ", "", -1)
	s3 := s2[:15] + "..."
	t.Log("----", s3)
}

// 测试获取子订单
func TestGetSubOrder(t *testing.T) {
	// -- 更新状态
	// update sale_sub_order set status = 1 WHERE order_no IN('1230322007642433','1230322001642486')
	// -- 更新deductAmount
	// update pay_order set deduct_amount = deduct_amount+1000,final_amount = final_amount-1000 where id=670
	// -- 删除已生成的支付单
	// delete FROM pay_order where out_order_no IN('1230322007642433','1230322001642486')
	orderNo := "1230324001307478"
	ret, _ := impl.OrderService.GetOrder(context.TODO(), &proto.OrderRequest{
		OrderNo:    orderNo,
		WithDetail: true,
	})
	t.Log(typeconv.MustJson(ret))
}

// 测试拆分支付单
func TestBreakPaymentOrder(t *testing.T){
	orderNo := "1230322000642437"
	ret, _ := impl.OrderService.BreakPaymentOrder(context.TODO(), &proto.BreakPaymentRequest{
		PaymentOrderNo:    orderNo,
	})
	t.Log(typeconv.MustJson(ret))	
}