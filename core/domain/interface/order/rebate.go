package order

type (
	// RebateItem 订单返利详情
	RebateItem struct {
		// 编号
		Id int `db:"id" pk:"yes" auto:"yes" json:"id"`
		// 返利单Id
		DebateId int64 `db:"debate_id" json:"debateId"`
		// 商品编号
		ItemId int64 `db:"item_id" json:"itemId"`
		// 商品名称
		ItemName string `db:"item_name" json:"itemName"`
		// 商品图片
		ItemImage string `db:"item_image" json:"itemImage"`
		// 商品金额
		ItemAmount int64 `db:"item_amount" json:"itemAmount"`
		// 返利金额
		RebateAmount int64 `db:"rebate_amount" json:"rebateAmount"`
	}

	// AffiliateRebate 订单返利
	AffiliateRebate struct {
		// 编号
		Id int64 `db:"id" pk:"yes" auto:"yes" json:"id"`
		// 返利方案Id
		PlanId int `db:"plan_id" json:"planId"`
		// 成交人Id
		TraderId int64 `db:"trader_id" json:"traderId"`
		// 分享码
		AffiliateCode string `db:"affiliate_code" json:"affiliateCode"`
		// 订单号
		OrderNo string `db:"order_no" json:"orderNo"`
		// 订单标题
		OrderSubject string `db:"order_subject" json:"orderSubject"`
		// 订单金额
		OrderAmount int64 `db:"order_amount" json:"orderAmount"`
		// 返利金额
		RebaseAmount int64 `db:"rebase_amount" json:"rebaseAmount"`
		// 返利状态
		Status int `db:"status" json:"status"`
		// 创建时间
		CreateTime int64 `db:"create_time" json:"createTime"`
		// 更新时间
		UpdateTime int64 `db:"update_time" json:"updateTime"`
	}
)
