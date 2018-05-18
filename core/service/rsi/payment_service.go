package rsi

/**
 * Copyright 2015 @ z3q.net.
 * name : payment_service.go
 * author : jarryliu
 * date : 2016-07-03 13:24
 * description :
 * history :
 */
import (
	"context"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/order"
	"go2o/core/domain/interface/payment"
	"go2o/core/module"
	"go2o/core/service/thrift/parser"
	"go2o/gen-code/thrift/define"
)

var _ define.PaymentService = new(paymentService)

type paymentService struct {
	repo       payment.IPaymentRepo
	orderRepo  order.IOrderRepo
	memberRepo member.IMemberRepo
}

func NewPaymentService(rep payment.IPaymentRepo, orderRepo order.IOrderRepo,
	memberRepo member.IMemberRepo) *paymentService {
	return &paymentService{
		repo:       rep,
		orderRepo:  orderRepo,
		memberRepo: memberRepo,
	}
}

// 根据编号获取支付单
func (p *paymentService) GetPaymentOrderById(ctx context.Context, id int32) (*define.SPaymentOrder, error) {
	po := p.repo.GetPaymentOrderById(int(id))
	if po != nil {
		v := po.Get()
		return parser.PaymentOrderDto(&v), nil
	}
	return nil, nil
}

// 根据交易号获取支付单编号
func (p *paymentService) GetPaymentOrderId(ctx context.Context, tradeNo string) (int32, error) {
	po := p.repo.GetPaymentOrder(tradeNo)
	if po != nil {
		return int32(po.GetAggregateRootId()), nil
	}
	return 0, nil
}

// 根据支付单号获取支付单
func (p *paymentService) GetPaymentOrder(ctx context.Context, paymentNo string) (*define.SPaymentOrder, error) {
	if po := p.repo.GetPaymentOrder(paymentNo); po != nil {
		v := po.Get()
		return parser.PaymentOrderDto(&v), nil
	}
	return nil, nil
}

// 创建支付单
func (p *paymentService) SubmitPaymentOrder(ctx context.Context, s *define.SPaymentOrder) (*define.Result_, error) {
	v := parser.PaymentOrder(s)
	o := p.repo.CreatePaymentOrder(v)
	err := o.Submit()
	return parser.Result(nil, err), nil
}

// 调整支付单金额
func (p *paymentService) AdjustOrder(ctx context.Context, paymentNo string, amount float64) (*define.Result_, error) {
	var err error
	o := p.repo.GetPaymentOrder(paymentNo)
	if o == nil {
		err = payment.ErrNoSuchPaymentOrder
	} else {
		err = o.Adjust(int(amount * 100))
	}
	return parser.Result(0, err), nil
}

// 积分抵扣支付单
func (p *paymentService) DiscountByIntegral(ctx context.Context, orderId int32,
	integral int64, ignoreOut bool) (r *define.Result_, err error) {
	var amount int
	o := p.repo.GetPaymentOrderById(int(orderId))
	if o == nil {
		err = payment.ErrNoSuchPaymentOrder
	} else {
		amount, err = o.IntegralDiscount(int(integral), ignoreOut)
	}
	return parser.Result(amount, err), nil
}

// 余额抵扣
func (p *paymentService) DiscountByBalance(ctx context.Context, orderId int32, remark string) (*define.Result_, error) {
	var err error
	o := p.repo.GetPaymentOrderById(int(orderId))
	if o == nil {
		err = payment.ErrNoSuchPaymentOrder
	} else {
		err = o.BalanceDiscount(remark)
	}
	return parser.Result(0, err), nil
}

// 钱包账户支付
func (p *paymentService) PaymentByWallet(ctx context.Context,
	tradeNo string, mergePay bool, remark string) (r *define.Result_, err error) {
	// 单个支付订单支付
	if !mergePay {
		ip := p.repo.GetPaymentOrder(tradeNo)
		if ip == nil {
			err = payment.ErrNoSuchPaymentOrder
		} else {
			err = ip.PaymentByWallet(remark)
		}
		return parser.Result(0, err), nil
	}
	// 合并支付单
	arr := p.repo.GetMergePayOrders(tradeNo)
	payUid := arr[0].Get().PayUid
	finalFee := 0
	for _, v := range arr {
		finalFee += v.Get().FinalFee
	}
	acc := p.memberRepo.GetAccount(payUid)
	if int(acc.Balance*100) < finalFee {
		err = member.ErrAccountBalanceNotEnough
	} else {
		for _, v := range arr {
			if err = v.PaymentByWallet(remark); err != nil {
				break
			}
		}
	}
	return parser.Result(nil, err), nil
}

// 余额钱包混合支付，优先扣除余额。
func (p *paymentService) HybridPayment(ctx context.Context, orderId int32, remark string) (r *define.Result_, err error) {
	o := p.repo.GetPaymentOrderById(int(orderId))
	if o == nil {
		err = payment.ErrNoSuchPaymentOrder
	} else {
		err = o.HybridPayment(remark)
	}
	return parser.Result(0, err), nil
}

// 完成支付单支付，并传入支付方式及外部订单号
func (p *paymentService) FinishPayment(ctx context.Context, tradeNo string, spName string,
	outerNo string) (r *define.Result_, err error) {
	o := p.repo.GetPaymentOrder(tradeNo)
	if o == nil {
		err = payment.ErrNoSuchPaymentOrder
	} else {
		err = o.PaymentFinish(spName, outerNo)
	}
	return parser.Result(nil, err), nil
}

// 支付网关
func (p *paymentService) GatewayV1(ctx context.Context, action string, userId int64, data map[string]string) (r *define.Result_, err error) {
	mod := module.Get(module.M_PAY).(*module.PaymentModule)
	// 获取令牌
	if action == "get_token" {
		token := mod.CreateToken(userId)
		return parser.Result(token, nil), nil
	}
	// 提交支付请求
	if action == "submit" {
		err = mod.Submit(userId, data)
	}
	// 验证支付
	if action == "payment" {
		err = mod.CheckAndPayment(userId, data)
	}
	return parser.Result(nil, err), nil
}

// 获取支付预交易数据
func (p *paymentService) GetPaymentOrderInfo(ctx context.Context,
	tradeNo string, mergePay bool) (*define.SPrepareTradeData, error) {
	var arr []payment.IPaymentOrder
	if mergePay {
		arr = p.repo.GetMergePayOrders(tradeNo)
	} else {
		ip := p.repo.GetPaymentOrder(tradeNo)
		if ip != nil {
			arr = []payment.IPaymentOrder{ip}
		}
	}
	return p.getMergePaymentOrdersInfo(tradeNo, arr, !true)
}

// 获取合并支付的支付单的支付数据
func (p *paymentService) getMergePaymentOrdersInfo(tradeNo string,
	tradeOrders []payment.IPaymentOrder, checkPay bool) (*define.SPrepareTradeData, error) {
	d := &define.SPrepareTradeData{ErrCode: 1, TradeOrders: []*define.SPaymentOrderData{}}
	if len(tradeOrders) == 0 {
		d.ErrMsg = "无效的支付订单"
		return d, nil
	}
	d.TradeState = payment.StateAwaitingPayment // 待支付
	for _, ip := range tradeOrders {
		// 检查支付状态
		if checkPay {
			if err := ip.CheckPaymentState(); err != nil {
				d.ErrMsg = err.Error()
				return d, nil
			}
		}
		iv := ip.Get()
		so := &define.SPaymentOrderData{
			OrderNo:      iv.OutOrderNo,
			Subject:      iv.Subject,
			TradeType:    iv.TradeType,
			State:        int32(iv.State),
			ProcedureFee: int32(iv.ProcedureFee),
			FinalFee:     int32(iv.FinalFee),
		}
		// 更新支付状态
		if so.State != payment.StateAwaitingPayment {
			d.TradeState = so.State
		}
		// 更新支付标志
		if i := int32(iv.PaymentFlag); d.PayFlag != i {
			d.PayFlag = i
		}
		// 更新支付金额
		d.TradeOrders = append(d.TradeOrders, so)
		d.ProcedureFee += int32(so.ProcedureFee) // 手续费
		d.FinalFee += int32(so.FinalFee)         // 最终金额
		d.TotalAmount += int32(iv.TotalAmount)   // 累计金额
	}
	d.ErrCode = 0
	d.TradeNo = tradeNo // 交易单号
	return d, nil
}
