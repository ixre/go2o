package merchant

const (
	TFlagNormal = 1 << iota
	// 免费
	TFlagFree
	// 关毕交易权限
	TFlagNoPerm
)

const (
	// 普通订单
	TKNormalOrder = 1
	// 批发订单
	TKWholesaleOrder = 2
	// 交易订单
	TKTradeOrder = 3
)

// 商户交易设置
type TradeConf struct {
	// 编号
	ID int64 `db:"id" pk:"yes" auto:"yes"`
	// 商户编号
	MchId int64 `db:"mch_id"`
	// 交易类型
	TradeType int `db:"trade_type"`
	// 交易方案，根据方案来自动调整比例
	PlanId int64 `db:"plan_id"`
	// 交易标志
	Flag int `db:"flag"`
	// 交易手续费依据,1:未设置 2:按金额 3:按比例
	AmountBasis int `db:"amount_basis"`
	// 交易费，按单笔收取
	TradeFee int `db:"trade_fee"`
	// 交易手续费比例
	TradeRate int `db:"trade_rate"`
	// 更新时间
	UpdateTime int64 `db:"update_time"`
}
