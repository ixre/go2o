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
	mergeAdd("商户提现是否免费", MerchantTakeOutCashFree, "true", "")
	mergeAdd("收货提示信息", OrderReceiveAlertMessage, "确认收货后,款项将转给商户。请在收货前确保已经商品没有损坏和缺少!", "")
	mergeAdd("是否启用会员经验值功能", MemberExperienceEnabled, "true", "")
	mergeAdd("会员资料不完善提醒信息", MemberProfileNotCompletedMessage, "您的个人资料未完善,是否立即完善?", "")
	mergeAdd("会员未实名认证提示信息", MemberNotTrustedMessage, "您尚未实名认证!", "")
	mergeAdd("会员是否验证手机号码格式", MemberCheckPhoneFormat, "true", "")
	mergeAdd("注册后赠送积分数量", PresentIntegralNumOfRegister, "0", "")
	mergeAdd("会员邀请关系级数", MemberReferLayer, "3", "")
	mergeAdd("会员即时通讯是否必须", MemberImRequired, "false", "")

	// 会员提现
	mergeAdd("会员提现开关", MemberTakeOutOn, "true", "")
	mergeAdd("会员提现提示",MemberTakeOutMessage,"提现功能暂不可用","")
	mergeAdd("会员提现是否必须实名制认证", MemberTakeOutMustTrust, "true", "")
	mergeAdd("会员最低提现金额",MemberMinTakeOutAmount,"0.01","")
	mergeAdd("会员单笔最高提现金额",MemberMaxTakeOutAmount,"5000.00","")
	mergeAdd("会员提现手续费费率",MemberTakeOutCsn,"0.00","")
	mergeAdd("会员每日提现上限",MemberMaxTakeOutTimesOfDay,"0","")
	// 会员转账
	mergeAdd("会员转账开关",MemberTransferAccountsOn,"true","")
	mergeAdd("会员转账提示信息",MemberTransferAccountsMessage,"系统转账功能暂不可用","")
	mergeAdd("会员转账手续费费率",MemberTransferCsn,"0.00","")

	mergeAdd("会员默认个人签名", MemberDefaultPersonRemark, "什么也没留下", "")
	mergeAdd("商品默认图片", GoodsDefaultImage, "res/nopic.gif", "")
	mergeAdd("商品最低利润率,既(销售价-供货价)/销售价的比例", GoodsMinProfitRate, "0", "")
	mergeAdd("广告缓存时间（秒）", CacheAdMaxAge, "3600", "")
	return mergeData
}
