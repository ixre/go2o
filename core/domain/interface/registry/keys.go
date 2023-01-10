package registry

var (

	/** 域名 */
	// 访问协议,https或http
	HttpProtocols = KeyFormat("HttpProtocols")
	// 根域名
	Domain = KeyFormat("Domain")
	// 控制面板前缀
	DomainPrefixDashboard = KeyFormat("DomainPrefixDashboard")
	// 零售门户前缀
	DomainPrefixPortal = KeyFormat("DomainPrefixPortal")
	// 批发门户域名前缀
	DomainPrefixWholesalePortal = KeyFormat("DomainPrefixWholesalePortal")
	// 零售门户手机端域名前缀
	DomainPrefixMobilePortal = KeyFormat("DomainPrefixMobilePortal")
	// 会员中心域名前缀
	DomainPrefixMember = KeyFormat("DomainPrefixMember")
	// 商户系统域名前缀
	DomainPrefixMerchant = KeyFormat("DomainPrefixMerchant")
	// 通行证域名前缀
	DomainPrefixPassport = KeyFormat("DomainPrefixPassport")
	// 通行证域名协议,默认为http,可以使用https安全加密
	DomainPassportProto = KeyFormat("DomainPassportProto")
	// API前缀
	DomainPrefixApi = KeyFormat("DomainPrefixApi")
	// HAPI前缀
	DomainPrefixHApi = KeyFormat("DomainPrefixHApi")
	// 文件服务器前缀
	DomainFileServerPrefix = KeyFormat("DomainFileServerPrefix")
	// 静态服务器前缀
	DomainPrefixStatic = KeyFormat("DomainPrefixStatic")
	// 图片服务器前缀
	DomainPrefixImage = KeyFormat("DomainPrefixImage")
	// 批发中心移动端
	DomainPrefixMobileWholesale = KeyFormat("DomainPrefixMobileWholesale")
	// 会员中心域名前缀(移动端)
	DomainPrefixMobileMember = KeyFormat("DomainPrefixMobileMember")
	// 通行证域名前缀(移动端)
	DomainPrefixMobilePassport = KeyFormat("DomainPrefixMobilePassport")

	/* 管理面板 */
	// 面板钩子显示名称
	BoardHookDisplayName = KeyFormat("board_hook_display_name")
	// 面板链接钩子访问密钥
	BoardHookToken = KeyFormat("board_hook_token")
	// 面板链接钩子URL地址
	BoardHookURL = KeyFormat("board_hook_url")

	/* API设置 */
	// 接口需要的最低版本
	ApiRequireVersion = KeyFormat("api_require_version")

	/* 平台 */
	// 平台名称
	PlatformName = KeyFormat("PlatformName")
	// 客服电话
	PlatformServiceTel = KeyFormat("PlatformServiceTel")
	// Logo标志
	PlatformLogo = KeyFormat("PlatformLogo")
	// 反色标志
	PlatformInverseColorLogo = KeyFormat("PlatformInverseColorLogo")
	// 零售门户标志
	PlatformRetailSiteLogo = KeyFormat("PlatformRetailSiteLogo")
	// 批发门户标志
	PlatformWholesaleSiteLogo = KeyFormat("PlatformWholesaleSiteLogo")

	/** 系统 */
	// 启用商户店铺商品分类
	EnableMchGoodsCategory = KeyFormat("EnableMchGoodsCategory")
	// 启用商户页面分类
	EnableMchPageCategory = KeyFormat("EnableMchPageCategory")
	// 开启调试模式
	EnableDebugMode = KeyFormat("EnableDebugMode")
	// 系统是否挂起
	SysSuspend = KeyFormat("SysSuspend")
	// 系统挂起提示消息
	SysSuspendMessage = KeyFormat("SysSuspendMessage")
	// 接口JWT密钥
	SysJWTSecret = KeyFormat("sys_jwt_secret")
	// 超级管理员登录密钥
	SysSuperLoginToken = KeyFormat("sys_super_login_token")
	/** 短信 */
	// 推送短信发送事件,如果开启,则短信通过系统外部发送
	SmsPushSendEvent = KeyFormat("SmsPushSendEvent")
	// 默认短信服务商
	SmsDefaultProvider = KeyFormat("SmsDefaultProvider")
	// 用户注册短信模板ID
	SmsRegisterTemplateId = KeyFormat("SmsRegisterTemplateId")
	// 用户验证码短信模板ID
	SmsMemberCheckTemplateId = KeyFormat("SmsMemberCheckTemplateId")
	// 短信接收间隔
	SmsSendDuration = KeyFormat("SmsSendDuration")

	/** 账户 */
	// 余额账户
	AccountBalanceAlias = KeyFormat("AccountBalanceAlias")
	// 积分账户
	AccountIntegralAlias = KeyFormat("AccountIntegralAlias")
	// 钱包账户
	AccountWalletAlias = KeyFormat("AccountWalletAlias")
	// 流动金账户
	AccountFlowAlias = KeyFormat("AccountFlowAlias")
	// 增利金账户
	AccountGrowthAlias = KeyFormat("AccountGrowthAlias")
	// 收货时的提示信息
	OrderReceiveAlertMessage = KeyFormat("OrderReceiveAlertMessage")

	/** 会员注册 */
	// 注册模式,1:普通注册 2:关闭注册 3:仅直接注册 4:仅邀请注册,等于member.RegisterMode
	MemberRegisterMode = KeyFormat("MemberRegisterMode")
	// 是否允许匿名注册
	MemberRegisterAllowAnonymous = KeyFormat("MemberRegisterAllowAnonymous")
	// 手机号码作为用户名
	MemberRegisterPhoneAsUser = KeyFormat("MemberRegisterPhoneAsUser")
	// 是否需要填写手机
	MemberRegisterNeedPhone = KeyFormat("MemberRegisterNeedPhone")
	// 必须绑定手机
	MemberRegisterMustBindPhone = KeyFormat("MemberRegisterMustBindPhone")
	// 是否需要填写即时通讯
	MemberRegisterNeedIm = KeyFormat("MemberRegisterNeedIm")
	// 注册提示
	MemberRegisterNotice = KeyFormat("MemberRegisterNotice")
	// 注册后赠送积分数量
	MemberRegisterPresentIntegral = KeyFormat("MemberRegisterPresentIntegral")
	// 邀请注册成功后跳转地址
	MemberInviteRegisterReturnUrl = KeyFormat("MemberInviteRegisterReturnUrl")
	// 邀请注册开启桥接页面,如跳转到注册前先显示一个页面
	MemberInviteEnableBridge = KeyFormat("MemberInviteEnableBridge")

	// 会员资料不完善提醒信息
	MemberProfileNotCompletedMessage = KeyFormat("MemberProfileNotCompletedMessage")
	// 会员实名提醒信息
	MemberNotTrustedMessage = KeyFormat("MemberNotTrustedMessage")
	// 实名时是否需要先完善资料
	MemberRequireProfileOnTrusting = KeyFormat("MemberRequireProfileOnTrusting")
	// 会员邀请关系级数
	MemberReferLayer = KeyFormat("MemberReferLayer")
	// 会员即时通讯是否必须
	MemberImRequired = KeyFormat("MemberImRequired")
	// 会员是否验证手机号码格式
	MemberCheckPhoneFormat = KeyFormat("MemberCheckPhoneFormat")
	// 会员默认个人签名
	MemberDefaultPersonRemark = KeyFormat("MemberDefaultPersonRemark")
	// 会员实名是否需要证件照片
	MemberTrustRequireCardImage = KeyFormat("MemberTrustRequireCardImage")
	// 会员实名是否需要人相图片
	MemberTrustRequirePeopleImage = KeyFormat("MemberTrustRequirePeopleImage")

	/*会员提现*/

	// MemberWithdrawEnabled  会员提现开关
	MemberWithdrawEnabled = KeyFormat("MemberWithdrawEnabled")
	// MemberWithdrawMessage 会员提现提示
	MemberWithdrawMessage = KeyFormat("MemberWithdrawMessage")
	// MemberWithdrawalMustVerification 会员提现是否必须实名制认证
	MemberWithdrawalMustVerification = KeyFormat("MemberWithdrawalMustVerification")
	// MemberWithdrawMinAmount 会员最低提现金额
	MemberWithdrawMinAmount = KeyFormat("MemberWithdrawMinAmount")
	// MemberWithdrawMaxAmount 会员单笔最高提现金额
	MemberWithdrawMaxAmount = KeyFormat("MemberWithdrawMaxAmount")
	// MemberWithdrawProcedureRate 会员提现手续费费率
	MemberWithdrawProcedureRate = KeyFormat("MemberWithdrawProcedureRate")
	// MemberWithdrawMaxTimeOfDay 会员每日提现上限
	MemberWithdrawMaxTimeOfDay = KeyFormat("MemberWithdrawMaxTimeOfDay")

	/*会员转账*/
	// MemberAccountTransferEnabled 会员转账开关
	MemberAccountTransferEnabled = KeyFormat("MemberAccountTransferEnabled")
	// MemberAccountTransferMessage 会员转账提示信息
	MemberAccountTransferMessage = KeyFormat("MemberAccountTransferMessage")
	// MemberAccountTransferProcedureRate 会员转账手续费费率
	MemberAccountTransferProcedureRate = KeyFormat("MemberAccountTransferProcedureRate")

	// 活动账户转为赠送可提现奖金手续费费率
	MemberFlowAccountConvertCsn = KeyFormat("MemberFlowAccountConvertCsn")
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
	// 兑换1元所需要的积分个数
	IntegralExchangeQuantity = KeyFormat("IntegralExchangeQuantity")
	// 抵扣1元所需要的积分个数
	IntegralDiscountQuantity = KeyFormat("IntegralDiscountQuantity")

	/** 订单 */
	// 是否启用订单返利
	OrderEnableAffiliateRebate = KeyFormat("OrderEnableAffiliateRebate")
	// 全局订单返利比例
	OrderGlobalAffiliateRebateRate = KeyFormat("OrderGlobalAffiteRebateRate")
	// 推送分销事件
	OrderPushAffiliateEvent = KeyFormat("OrderPushAffiliateEvent")
	//  /* 商户订单 */
	// MchOrderSettleMode 商户订单结算模式
	MchOrderSettleMode = KeyFormat("MchOrderSettleMode")
	// MchOrderSettleRate 商户订单结算比例
	MchOrderSettleRate = KeyFormat("MchOrderSettleRate")
	// MchSingleOrderServiceFee 商户订单每单服务费(按单结算)
	MchSingleOrderServiceFee = KeyFormat("MchSingleOrderServiceFee")
	// MchMonthFreeOrders 商户订单每月免服务费订单数
	MchMonthFreeOrders = KeyFormat("MchMonthFreeOrders")
	// 是否必须认证后才可以继续操作
	MchMustBeTrust = KeyFormat("MchMustBeTrust")
	// MchOrderRequireTicket 商户交易单是否需上传发票
	MchOrderRequireTicket = KeyFormat("MchOrderRequireTicket")

	/* 商户 */
	// 商户提现手续费
	MerchantTakeOutCashFree = KeyFormat("MerchantTakeOutCashFree")
	// 商户提现手续费费率
	MerchantTakeOutCsn = KeyFormat("MerchantTakeOutCsn")
	// 商户提现
	MerchantMinTakeOutAmount = KeyFormat("MerchantMinTakeOutAmount")

	// 商品默认图片
	GoodsDefaultImage = KeyFormat("GoodsDefaultImage")
	// 商品最低利润率,既(销售价-供货价)/销售价的比例
	GoodsMinProfitRate = KeyFormat("GoodsMinProfitRate")
	// 广告缓存时间（秒）
	CacheAdMaxAge = KeyFormat("CacheAdMaxAge")
	// 敏感词,以|分割
	//ShopIncorrectAliasWords string = "ShopIncorrectAliasWords");
)
