package merchant

// 商户销售服务
type ISaleManager interface {
	// 计算交易费用,返回交易费及错误
	MathTradeFee(tradeType int, amount int64) (int64, error)
}
