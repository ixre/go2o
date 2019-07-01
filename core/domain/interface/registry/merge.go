package registry

var mergeData = make([]*Registry, 0)

func mergeAdd(description string, key string, defaultValue string, options string) {
	r := &Registry{
		Key:          key,
		Value:        defaultValue,
		DefaultValue: defaultValue,
		Options:      options,
		UserDefine:   0,
		Description:  description,
	}
	mergeData = append(mergeData, r)
}

// 返回需要合并的注册表数据
func MergeRegistries() []*Registry {
	/** 域名 */
	mergeAdd("是否启用SSL连接", DomainEnabledSSL, "false", "")
	mergeAdd("零售门户前缀", DomainPrefixPortal, "www.", "")
	mergeAdd("批发门户域名前缀", DomainPrefixWholesalePortal, "whs.", "")
	mergeAdd("零售门户手机端域名前缀", DomainPrefixMobilePortal, "m.", "")
	mergeAdd("会员中心域名前缀", DomainPrefixMember, "u.", "")
	mergeAdd("商户系统域名前缀", DomainPrefixMerchant, "mch.", "")
	mergeAdd("通行证域名前缀", DomainPrefixPassport, "passport.", "")
	mergeAdd("通行证域名协议", DomainPassportProto, "http", "http或https")
	mergeAdd("API系统", DomainPrefixApi, "api.", "")
	mergeAdd("静态服务器前缀", DomainPrefixStatic, "static.", "")
	mergeAdd("图片服务器前缀", DomainPrefixImage, "img.", "")
	mergeAdd("批发中心移动端", DomainPrefixMobileWholesale, "mwhs.", "")
	mergeAdd("会员中心域名前缀(移动端)", DomainPrefixMobileMember, "mu.", "")
	mergeAdd("通行证域名前缀(移动端)", DomainPrefixMobilePassport, "mpp.", "")

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
	mergeAdd("系统是否挂起", SysSuspend, "false", "")
	mergeAdd("系统挂起提示消息", SysSuspendMessage, "系统正在升级维护，请稍后再试!", "")

	mergeAdd("收货提示信息", OrderReceiveAlertMessage, "确认收货后,款项将转给商户。请在收货前确保已经商品没有损坏和缺少!", "")

	/** 会员注册 */
	// 注册模式,1:普通注册 2:关闭注册 3:仅直接注册 4:仅邀请注册,等于member.RegisterMode
	mergeAdd("注册模式,1:普通注册 2:关闭注册 3:仅直接注册 4:仅邀请注册",MemberRegisterMode,"1","")
	mergeAdd("是否允许匿名注册", MemberRegisterAllowAnonymous ,"true","")
	mergeAdd("手机号码作为用户名", MemberRegisterPhoneAsUser,"false","")
	mergeAdd("是否需要填写手机", MemberRegisterNeedPhone,"false","")
	mergeAdd("必须绑定手机", MemberRegisterMustBindPhone,"false","")
	mergeAdd("是否需要填写即时通讯", MemberRegisterNeedIm,"false","")
	mergeAdd("注册提示", MemberRegisterNotice ,"","")
	mergeAdd("注册成功后跳转地址,默认登录页面", MemberRegisterReturnUrl,"/auth?uc=1","")
	mergeAdd("注册后赠送积分数量", MemberRegisterPresentIntegral, "0", "")

	mergeAdd("会员资料不完善提醒信息", MemberProfileNotCompletedMessage, "您的个人资料未完善,是否立即完善?", "")
	mergeAdd("会员未实名认证提示信息", MemberNotTrustedMessage, "您尚未实名认证!", "")
	mergeAdd("实名时是否需要先完善资", MemberRequireProfileOnTrusting, "false", "")
	mergeAdd("会员是否验证手机号码格式", MemberCheckPhoneFormat, "true", "")

	mergeAdd("会员邀请关系级数", MemberReferLayer, "3", "")
	mergeAdd("会员即时通讯是否必须", MemberImRequired, "false", "")

	// 会员提现
	mergeAdd("会员提现开关", MemberTakeOutOn, "true", "")
	mergeAdd("会员提现提示", MemberTakeOutMessage, "提现功能暂不可用", "")
	mergeAdd("会员提现是否必须实名制认证", MemberTakeOutMustTrust, "true", "")
	mergeAdd("会员最低提现金额", MemberMinTakeOutAmount, "0.01", "")
	mergeAdd("会员单笔最高提现金额", MemberMaxTakeOutAmount, "5000.00", "")
	mergeAdd("会员提现手续费费率", MemberTakeOutCsn, "0.00", "")
	mergeAdd("会员每日提现上限", MemberMaxTakeOutTimesOfDay, "0", "")
	// 会员转账
	mergeAdd("会员转账开关", MemberTransferAccountsOn, "true", "")
	mergeAdd("会员转账提示信息", MemberTransferAccountsMessage, "平台仅提供转账功能，请尽量当面交易以保证安全！", "")
	mergeAdd("会员转账手续费费率", MemberTransferCsn, "0.00", "")
	mergeAdd("活动账户转为赠送可提现奖金手续费费率", MemberFlowAccountConvertCsn, "0.20", "")
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
	// 商户订单
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
