package order

type (
	// DistributionItem 订单返利详情
	DistributionItem struct {
		// 编号
		Id int64 `db:"id" pk:"yes" auto:"yes" json:"id"`
		// 分销单Id
		DistributeId int64 `db:"distribute_id" json:"distributeId"`
		// 商品编号
		ItemId int64 `db:"item_id" json:"itemId"`
		// 商品名称
		ItemName string `db:"item_name" json:"itemName"`
		// 商品图片
		ItemImage string `db:"item_image" json:"itemImage"`
		// 商品金额
		ItemAmount int64 `db:"item_amount" json:"itemAmount"`
		// 分销金额
		DistributeAmount int64 `db:"distribute_amount" json:"distributeAmount"`
	}

	// AffiliateDistribution 订单分销
	AffiliateDistribution struct {
		// 编号
		Id int64 `db:"id" pk:"yes" auto:"yes" json:"id"`
		// 返利方案Id
		PlanId int `db:"plan_id" json:"planId"`
		// 买家
		BuyerId int64 `db:"buyer_id" json:"buyerId"`
		// 返利所有人编号
		OwnerId int64 `db:"owner_id" json:"ownerId"`
		// 标志
		Flag int16 `db:"flag" json:"flag"`
		// 是否已读
		IsRead int16 `db:"is_read" json:"isRead"`
		// 分享码
		AffiliateCode string `db:"affiliate_code" json:"affiliateCode"`
		// 订单号
		OrderNo string `db:"order_no" json:"orderNo"`
		// 订单标题
		OrderSubject string `db:"order_subject" json:"orderSubject"`
		// 订单金额
		OrderAmount int64 `db:"order_amount" json:"orderAmount"`
		// 分销奖励金额
		DistributeAmount int64 `db:"distribute_amount" json:"distributeAmount"`
		// 返利状态
		Status int `db:"status" json:"status"`
		// 创建时间
		CreateTime int64 `db:"create_time" json:"createTime"`
		// 更新时间
		UpdateTime int64 `db:"update_time" json:"updateTime"`
	}
)
