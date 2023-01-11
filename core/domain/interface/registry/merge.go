package registry

var mergeData = make([]*Registry, 0)

func mergeAdd(description string, key string, defaultValue string, options string) {
	r := &Registry{
		Key:          key,
		Value:        defaultValue,
		DefaultValue: defaultValue,
		Options:      options,
		Flag:         FlagUserDefine,
		Description:  description,
	}
	mergeData = append(mergeData, r)
}

// 返回需要合并的注册表数据
func MergeRegistries() []*Registry {
	/** 域名 */
	mergeAdd("访问协议", HttpProtocols, "http", "")
	mergeAdd("域名", Domain, "yourdomain.com", "")
	mergeAdd("管理面板前缀", DomainPrefixDashboard, "board.", "")
	mergeAdd("零售门户前缀", DomainPrefixPortal, "www.", "")
	mergeAdd("批发门户域名前缀", DomainPrefixWholesalePortal, "whs.", "")
	mergeAdd("零售门户手机端域名前缀", DomainPrefixMobilePortal, "m.", "")
	mergeAdd("会员中心域名前缀", DomainPrefixMember, "u.", "")
	mergeAdd("商户系统域名前缀", DomainPrefixMerchant, "mch.", "")
	mergeAdd("通行证域名前缀", DomainPrefixPassport, "passport.", "")
	mergeAdd("通行证域名协议", DomainPassportProto, "http", "http或https")
	mergeAdd("API前缀", DomainPrefixApi, "api.", "")
	mergeAdd("HAPI前缀", DomainPrefixHApi, "hapi.", "")
	mergeAdd("文件服务器前缀", DomainFileServerPrefix, "/files/", "")
	mergeAdd("静态服务器前缀", DomainPrefixStatic, "static.", "")
	mergeAdd("图片服务器前缀", DomainPrefixImage, "img.", "")
	mergeAdd("批发中心移动端", DomainPrefixMobileWholesale, "mwhs.", "")
	mergeAdd("会员中心域名前缀(移动端)", DomainPrefixMobileMember, "mu.", "")
	mergeAdd("通行证域名前缀(移动端)", DomainPrefixMobilePassport, "mpp.", "")

	/* 管理面板 */
	mergeAdd("面板钩子显示名称", BoardHookDisplayName, "Hooks", "")
	mergeAdd("面板链接钩子访问密钥", BoardHookToken, "", "")
	mergeAdd("面板链接钩子URL地址", BoardHookURL, "", "")
	/* API设置 */
	mergeAdd("接口需要的最低版本", ApiRequireVersion, "0.0.1", "")

	/* 账户 */
	mergeAdd("余额账户", AccountBalanceAlias, "余额", "")
	mergeAdd("积分账户", AccountIntegralAlias, "积分", "")
	mergeAdd("钱包账户", AccountWalletAlias, "钱包", "")
	mergeAdd("流动金账户", AccountFlowAlias, "流动金", "")
	mergeAdd("增利金账户", AccountGrowthAlias, "增利金", "")

	/* 平台 */
	mergeAdd("平台名称", PlatformName, "GO2O商城系统", "")
	mergeAdd("客服客服电话", PlatformServiceTel, "+86-021-66666666", "")
	mergeAdd("Logo标志", PlatformLogo, "//raw.githubusercontent.com/jsix/go2o/master/docs/mark.gif", "")
	mergeAdd("反色标志", PlatformInverseColorLogo, "//raw.githubusercontent.com/jsix/go2o/master/docs/mark.gif", "")
	mergeAdd("零售门户标志", PlatformRetailSiteLogo, "//raw.githubusercontent.com/jsix/go2o/master/docs/mark.gif", "")
	mergeAdd("批发门户标志", PlatformWholesaleSiteLogo, "//raw.githubusercontent.com/jsix/go2o/master/docs/mark.gif", "")

	/** 系统 */
	mergeAdd("启用商户店铺商品分类", EnableMchGoodsCategory, "false", "")
	mergeAdd("启用商户页面分类", EnableMchPageCategory, "false", "")
	mergeAdd("开启调试模式", EnableDebugMode, "false", "")
	mergeAdd("系统是否挂起", SysSuspend, "false", "")
	mergeAdd("系统挂起提示消息", SysSuspendMessage, "系统正在升级维护，请稍后再试!", "")
	mergeAdd("接口JWT密钥", SysJWTSecret, "", "")
	mergeAdd("超级管理员登录密钥", SysSuperLoginToken, "", "")
	mergeAdd("收货提示信息", OrderReceiveAlertMessage, "确认收货后,款项将转给商户。请在收货前确保已经商品没有损坏和缺少!", "")

	/** 短信 */
	mergeAdd("默认短信服务商", SmsDefaultProvider, "", "")
	mergeAdd("用户注册短信模板ID", SmsRegisterTemplateId, "", "")
	mergeAdd("用户验证码短信模板ID", SmsMemberCheckTemplateId, "", "")
	mergeAdd("短信接收间隔,默认(2s)", SmsSendDuration, "2000", "")
	// 注册模式,1:普通注册 2:关闭注册 3:仅直接注册 4:仅邀请注册,等于member.RegisterMode
	mergeAdd("注册模式,1:普通注册 2:关闭注册 3:仅直接注册 4:仅邀请注册", MemberRegisterMode, "1", "")
	mergeAdd("是否允许匿名注册", MemberRegisterAllowAnonymous, "true", "")
	mergeAdd("手机号码作为用户名", MemberRegisterPhoneAsUser, "false", "")
	mergeAdd("是否需要填写手机", MemberRegisterNeedPhone, "false", "")
	mergeAdd("必须绑定手机", MemberRegisterMustBindPhone, "false", "")
	mergeAdd("是否需要填写即时通讯", MemberRegisterNeedIm, "false", "")
	mergeAdd("注册提示", MemberRegisterNotice, "", "")
	mergeAdd("邀请注册成功后跳转地址", MemberInviteRegisterReturnUrl, "", "")
	mergeAdd("注册后赠送积分数量", MemberRegisterPresentIntegral, "0", "")
	mergeAdd("邀请注册开启桥接页面,跳转到注册页前显示一个自定义页面", MemberInviteEnableBridge, "false", "")

	mergeAdd("会员资料不完善提醒信息", MemberProfileNotCompletedMessage, "您的个人资料未完善,是否立即完善?", "")
	mergeAdd("会员未实名认证提示信息", MemberNotTrustedMessage, "您尚未实名认证!", "")
	mergeAdd("实名时是否需要先完善资", MemberRequireProfileOnTrusting, "false", "")
	mergeAdd("会员是否验证手机号码格式", MemberCheckPhoneFormat, "true", "")
	mergeAdd("会员邀请关系级数", MemberReferLayer, "3", "")
	mergeAdd("会员即时通讯是否必须", MemberImRequired, "false", "")
	mergeAdd("会员实名是否需要证件照片", MemberTrustRequireCardImage, "false", "")
	mergeAdd("会员实名是否需要人相图片", MemberTrustRequirePeopleImage, "true", "")

	// 会员提现
	mergeAdd("启用会员提现", MemberWithdrawEnabled, "true", "")
	mergeAdd("会员提现提示", MemberWithdrawMessage, "提现功能暂不可用", "")
	mergeAdd("会员提现是否必须实名制认证", MemberWithdrawalMustVerification, "true", "")
	mergeAdd("会员最低提现金额", MemberWithdrawMinAmount, "0.01", "")
	mergeAdd("会员单笔最高提现金额", MemberWithdrawMaxAmount, "5000.00", "")
	mergeAdd("会员提现手续费费率", MemberWithdrawProcedureRate, "0.00", "")
	mergeAdd("会员每日提现上限", MemberWithdrawMaxTimeOfDay, "0", "")
	// 会员转账
	mergeAdd("启用会员转账", MemberAccountTransferEnabled, "true", "")
	mergeAdd("会员转账提示信息", MemberAccountTransferMessage, "平台仅提供转账功能，请尽量当面交易以保证安全！", "")
	mergeAdd("会员转账手续费费率", MemberAccountTransferProcedureRate, "0.00", "")
	mergeAdd("活动账户转为赠送可提现奖金手续费费率", MemberFlowAccountConvertCsn, "0.20", "")
	// 会员信息推送
	mergeAdd("是否启用会员账户信息消息推送", MemberAccountPushEnabled, "0", "0:关闭,1:启用")
	mergeAdd("是否启用会员邀请关系变更消息推送", MemberRelationPushEnabled, "0", "0:关闭 2:启用")

	// 经验值
	mergeAdd("是否启用会员经验值功能", ExperienceEnabled, "true", "")
	mergeAdd("会员普通消费1元产生的经验比例", ExperienceRateByOrder, "1.00", "")
	mergeAdd("会员线下消费1元产生的经验比例", ExperienceRateByTradeOrder, "0.00", "")
	mergeAdd("会员批发消费1元产生的经验比例", ExperienceRateByWholesaleOrder, "0.00", "")

	// 积分
	mergeAdd("会员普通消费1元产生的积分比例", IntegralRateByConsumption, "1.00", "")
	mergeAdd("会员线下消费1元产生的积分比例", IntegralRateByTradeOrder, "0.00", "")
	mergeAdd("会员批发消费1元产生的积分比例", IntegralRateByWholesaleOrder, "0.00", "")
	mergeAdd("兑换1元所需要的积分(0不兑换)", IntegralExchangeQuantity, "1000", "")
	mergeAdd("抵扣1元所需要的积分(0不抵扣)", IntegralDiscountQuantity, "1000", "")

	// 订单
	mergeAdd("是否启用订单返利", OrderEnableAffiliateRebate, "0", "")
	mergeAdd("全局订单返利比例", OrderGlobalAffiliateRebateRate, "0", "")
	mergeAdd("推送分销事件", OrderPushAffiliateEvent, "0", "0:不推送(内部处理),1:仅推送(内部处理),2:推送并处理(外部处理分销)")
	// 商户订单
	mergeAdd("是否必须认证后才可上传商品", MchMustBeTrust, "true", "")
	mergeAdd("商户订单结算模式", MchOrderSettleMode, "1", "1:按供货价,2:按销售额,3:按单")
	mergeAdd("商户订单结算比例", MchOrderSettleRate, "0.05", "")
	mergeAdd("商户订单每单服务费(按单结算)", MchSingleOrderServiceFee, "1.00", "")
	mergeAdd("商户订单每月免服务费订单数", MchMonthFreeOrders, "0", "")
	mergeAdd("商户交易单是否需上传发票", MchOrderRequireTicket, "false", "")
	// 商户
	mergeAdd("商户提现是否免费", MerchantTakeOutCashFree, "true", "")
	mergeAdd("商户提现手续费费率", MerchantTakeOutCsn, "0.00", "")
	mergeAdd("商户提现最低金额", MerchantMinTakeOutAmount, "1", "")

	mergeAdd("会员默认个人签名", MemberDefaultPersonRemark, "什么也没留下", "")
	mergeAdd("商品默认图片", GoodsDefaultImage, "res/nopic.gif", "")
	mergeAdd("商品最低利润率,既(销售价-供货价)/销售价的比例", GoodsMinProfitRate, "0", "")
	mergeAdd("广告缓存时间（秒）", CacheAdMaxAge, "3600", "")
	return mergeData
}
