package wholesaler

/* 批发商 */

type (
	// 批发商
	IWholesaler interface {
		// 获取领域编号
		GetDomainId() int32
		// 获取值
		Value() *WsWholesaler
		// 审核批发商
		Review(pass bool, reason string) error
		// 停止批发权限
		Abort() error
		// 保存
		Save() (int32, error)
		// 保存批发返点率
		SaveRebateRate(*WsRebateRate) (int32, error)
		// 获取批发返点率
		RebateRates() []*WsRebateRate
	}

	IWholesaleRepo interface {
		// Get WsWholesaler
		GetWsWholesaler(primary interface{}) *WsWholesaler
		// Save WsWholesaler
		SaveWsWholesaler(v *WsWholesaler, create bool) (int, error)
		// Select WsRebateRate
		SelectWsRebateRate(where string, v ...interface{}) []*WsRebateRate
		// Save WsRebateRate
		SaveWsRebateRate(v *WsRebateRate) (int, error)
		// Batch Delete WsRebateRate
		BatchDeleteWsRebateRate(where string, v ...interface{}) (int64, error)
	}
	// 批发商
	WsWholesaler struct {
		// 供货商编号等于供货商（等同与商户编号)
		MchId int32 `db:"mch_id" pk:"yes" auto:"yes"`
		// 批发商评级
		Rate int `db:"rate"`
		// 批发商审核状态
		ReviewState int32 `db:"review_state"`
	}
	// 批发客户分组返点比例设置
	WsRebateRate struct {
		// 编号
		ID int32 `db:"id" pk:"yes" auto:"yes"`
		// 批发商编号
		WsId int32 `db:"ws_id"`
		// 客户分组编号
		BuyerGid int32 `db:"buyer_gid"`
		// 下限金额
		RequireAmount int32 `db:"require_amount"`
		// 返点率
		RebateRate float64 `db:"rebate_rate"`
	}
)
