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

	mergeAdd("收货提示信息", OrderReceiveAlertMessage, "确认收货后,款项将转给商户。请在收货前确保已经商品没有损坏和缺少!", "")

	mergeAdd("会员资料不完善提醒信息", MemberProfileNotCompletedMessage, "您的个人资料未完善,是否立即完善?", "")
	mergeAdd("会员未实名认证提示信息", MemberNotTrustedMessage, "您尚未实名认证!", "")
	mergeAdd("会员是否验证手机号码格式", MemberCheckPhoneFormat, "true", "")
	mergeAdd("注册后赠送积分数量", MemberPresentIntegralNumOfRegister, "0", "")
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
	mergeAdd("会员转账提示信息", MemberTransferAccountsMessage, "系统转账功能暂不可用", "")
	mergeAdd("会员转账手续费费率", MemberTransferCsn, "0.00", "")
	// 经验值
	mergeAdd("是否启用会员经验值功能", ExperienceEnabled, "true", "")
	mergeAdd("会员普通消费1元产生的经验比例", ExperienceRateByOrder, "1.00", "")
	mergeAdd("会员线下消费1元产生的经验比例", ExperienceRateByTradeOrder, "0.00", "")
	mergeAdd("会员批发消费1元产生的经验比例", ExperienceRateByWholesaleOrder, "0.00", "")

	// 积分
	mergeAdd("会员普通消费1元产生的积分比例", IntegralRateByConsumption, "1.00", "")
	mergeAdd("会员线下消费1元产生的积分比例", IntegralRateByTradeOrder, "0.00", "")
	mergeAdd("会员批发消费1元产生的积分比例", IntegralRateByWholesaleOrder, "0.00", "")
	// 商户订单
	mergeAdd("商户订单结算模式", MchOrderSettleMode, "1", "1:按供货价,2:按销售额,3:按单")
	mergeAdd("商户订单结算比例", MchOrderSettleRate, "0.05", "")
	mergeAdd("商户订单每单服务费(按单结算)", MchSingleOrderServiceFee, "1.00", "")
	mergeAdd("商户订单每月免服务费订单数", MchMonthFreeOrders, "0", "")
	// 商户
	mergeAdd("商户提现是否免费", MerchantTakeOutCashFree, "true", "")
	mergeAdd("商户提现手续费费率", MerchantTakeOutCsn, "0.00", "")

	mergeAdd("会员默认个人签名", MemberDefaultPersonRemark, "什么也没留下", "")
	mergeAdd("商品默认图片", GoodsDefaultImage, "res/nopic.gif", "")
	mergeAdd("商品最低利润率,既(销售价-供货价)/销售价的比例", GoodsMinProfitRate, "0", "")
	mergeAdd("广告缓存时间（秒）", CacheAdMaxAge, "3600", "")
	return mergeData
}
