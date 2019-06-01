package registry

var (
	// 收货时的提示信息
	OrderReceiveAlertMessage = KeyFormat("OrderReceiveAlertMessage")

	// 会员资料不完善提醒信息
	MemberProfileNotCompletedMessage = KeyFormat("MemberProfileNotCompletedMessage")
	// 会员实名提醒信息
	MemberNotTrustedMessage = KeyFormat("MemberNotTrustedMessage")
	// 注册后赠送积分数量
	MemberPresentIntegralNumOfRegister = KeyFormat("MemberPresentIntegralNumOfRegister")
	// 会员邀请关系级数
	MemberReferLayer = KeyFormat("MemberReferLayer")
	// 会员即时通讯是否必须
	MemberImRequired = KeyFormat("MemberImRequired")
	// 会员是否验证手机号码格式
	MemberCheckPhoneFormat = KeyFormat("MemberCheckPhoneFormat")
	// 会员默认个人签名
	MemberDefaultPersonRemark = KeyFormat("MemberDefaultPersonRemark")
	/*会员提现*/
	// 会员提现开关
	MemberTakeOutOn = KeyFormat("MemberTakeOutOn")
	// 会员提现提示
	MemberTakeOutMessage = KeyFormat("MemberTakeOutMessage")
	// 会员提现是否必须实名制认证
	MemberTakeOutMustTrust = KeyFormat("MemberTakeOutMustTrust")
	// 会员最低提现金额
	MemberMinTakeOutAmount = KeyFormat("MemberMinTakeOutAmount")
	// 会员单笔最高提现金额
	MemberMaxTakeOutAmount = KeyFormat("MemberMaxTakeOutAmount")
	// 会员提现手续费费率
	MemberTakeOutCsn = KeyFormat("MemberTakeOutCsn")
	// 会员每日提现上限
	MemberMaxTakeOutTimesOfDay = KeyFormat("MemberMaxTakeOutTimesOfDay")

	/*会员转账*/
	// 会员转账开关
	MemberTransferAccountsOn = KeyFormat("MemberTransferAccountsOn")
	// 会员转账提示信息
	MemberTransferAccountsMessage = KeyFormat("MemberTransferAccountsMessage")
	// 会员转账手续费费率
	MemberTransferCsn = KeyFormat("MemberTransferCsn")

	/* 经验值 */
	// 是否启用会员经验值功能
	ExperienceEnabled = KeyFormat("ExperienceEnabled")
	// 会员普通消费1元产生的经验比例
	ExperienceRateByOrder = KeyFormat("ExperienceRateByOrder")
	// 会员线下消费1元产生的经验比例
	ExperienceRateByTradeOrder = KeyFormat("ExperienceRateByTradeOrder")
	// 会员批发消费1元产生的经验比例
	ExperienceRateByWholesaleOrder = KeyFormat("ExperienceRateByWholesaleOrder")

	/* 积分 */
	// 会员普通消费1元产生的积分比例
	IntegralRateByConsumption = KeyFormat("IntegralRateByConsumption")
	// 会员线下消费1元产生的积分比例
	IntegralRateByTradeOrder = KeyFormat("IntegralRateByTradeOrder")
	// 会员批发消费1元产生的积分比例
	IntegralRateByWholesaleOrder = KeyFormat("IntegralRateByWholesaleOrder")

	/* 商户订单 */
	// 商户订单结算模式
	MchOrderSettleMode = KeyFormat("MchOrderSettleMode")
	// 商户订单结算比例
	MchOrderSettleRate = KeyFormat("MchOrderSettleRate")
	// 商户订单每单服务费(按单结算)
	MchSingleOrderServiceFee = KeyFormat("MchSingleOrderServiceFee")
	// 商户订单每月免服务费订单数
	MchMonthFreeOrders = KeyFormat("MchMonthFreeOrders")

	/* 商户 */
	// 商户提现手续费
	MerchantTakeOutCashFree = KeyFormat("MerchantTakeOutCashFree")
	// 商户提现手续费费率
	MerchantTakeOutCsn = KeyFormat("MerchantTakeOutCsn")

	// 商品默认图片
	GoodsDefaultImage = KeyFormat("GoodsDefaultImage")
	// 商品最低利润率,既(销售价-供货价)/销售价的比例
	GoodsMinProfitRate = KeyFormat("GoodsMinProfitRate")
	// 广告缓存时间（秒）
	CacheAdMaxAge = KeyFormat("CacheAdMaxAge")
	// 敏感词,以|分割
	//ShopIncorrectAliasWords string = "ShopIncorrectAliasWords");
)
