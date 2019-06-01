package registry

var (
	// 商户提现手续费
	MerchantTakeOutCashFree = KeyFormat("MerchantTakeOutCashFree")
	// 收货时的提示信息
	OrderReceiveAlertMessage = KeyFormat("OrderReceiveAlertMessage")
	// 是否启用会员经验值功能
	MemberExperienceEnabled = KeyFormat("MemberExperienceEnabled")
	// 会员资料不完善提醒信息
	MemberProfileNotCompletedMessage = KeyFormat("MemberProfileNotCompletedMessage")
	// 会员实名提醒信息
	MemberNotTrustedMessage = KeyFormat("MemberNotTrustedMessage")
	// 注册后赠送积分数量
	PresentIntegralNumOfRegister = KeyFormat("PresentIntegralNumOfRegister")
	// 会员邀请关系级数
	MemberReferLayer = KeyFormat("MemberReferLayer")
	// 会员即时通讯是否必须
	MemberImRequired = KeyFormat("MemberImRequired")
	// 会员是否验证手机号码格式
	MemberCheckPhoneFormat = KeyFormat("MemberCheckPhoneFormat")
	// 会员默认个人签名
	MemberDefaultPersonRemark = KeyFormat("MemberDefaultPersonRemark")

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
	MemberMaxTakeOutTimesOfDay =KeyFormat("MemberMaxTakeOutTimesOfDay")

	// 会员转账开关
	MemberTransferAccountsOn = KeyFormat("MemberTransferAccountsOn")
	// 会员转账提示信息
	MemberTransferAccountsMessage = KeyFormat("MemberTransferAccountsMessage")
	// 会员转账手续费费率
	MemberTransferCsn = KeyFormat("MemberTransferCsn")

	// 商品默认图片
	GoodsDefaultImage = KeyFormat("GoodsDefaultImage")
	// 商品最低利润率,既(销售价-供货价)/销售价的比例
	GoodsMinProfitRate = KeyFormat("GoodsMinProfitRate")
	// 广告缓存时间（秒）
	CacheAdMaxAge = KeyFormat("CacheAdMaxAge")
	// 敏感词,以|分割
	//ShopIncorrectAliasWords string = "ShopIncorrectAliasWords");
)
