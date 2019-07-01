/**
 * Copyright 2015 @ to2.net.
 * name : default
 * author : jarryliu
 * date : 2016-07-23 11:21
 * description :
 * history :
 */
package repos

import (
	"go2o/core/domain/interface/valueobject"
)

var (
	DefaultRegistry = valueobject.Registry_{
		// 商户提现是否免费
		MerchantTakeOutCashFree: true,
		// 收货提示信息
		OrderReceiveAlertMessage: "确认收货后,款项将转给商户。请在收货前确保已经商品没有损坏和缺少!",
		// 是否启用会员经验值功能
		MemberExperienceEnabled: true,
		//会员资料不完善提醒信息
		MemberProfileNotCompletedMessage: "您的个人资料未完善,是否立即完善?",
		// 会员未实名认证提示信息
		MemberNotTrustedMessage: "您尚未实名认证!",
		// 会员是否验证手机号码格式
		MemberCheckPhoneFormat: true,
		// 注册后赠送积分数量
		PresentIntegralNumOfRegister: 0,
		MemberReferLayer:             3,
		// 会员即时通讯是否必须
		MemberImRequired: false,
		// 会员提现开关
		MemberTakeOutOn: true,
		// 提现是否必须实名制认证
		TakeOutMustTrust: true,

		// 会员默认个人签名
		MemberDefaultPersonRemark: "什么也没留下",

		// 商品默认图片
		GoodsDefaultImage: "res/nopic.gif",
		// 商品最低利润率,既(销售价-供货价)/销售价的比例
		GoodsMinProfitRate: 0,
		// 广告缓存时间（秒）
		CacheAdMaxAge: 3600,

		// 商铺别名敏感词,以|分割
		ShopIncorrectAliasWords: "shop|master|o2o|super|www|sys|system|mall|mch|system|passport|api|image|static|img",
		//RegistryData: map[string]string{
		//	"PlatformName":      "GO2O",
		//	"Logo":              "https://raw.githubusercontent.com/jsix/go2o/master/docs/mark.gif",
		//	"RetailSiteLogo":    "https://raw.githubusercontent.com/jsix/go2o/master/docs/mark.gif",
		//	"WholesaleSiteLogo": "https://raw.githubusercontent.com/jsix/go2o/master/docs/mark.gif",
		//	"Telephone":         "021-88888888",
		//	// 会员转账开关
		//	//valueobject.RKMemberTransferAccountsOn: "true",
		//	// 会员转账提示信息
		//	valueobject.RKMemberTransferAccountsMessage: "平台仅提供转账功能，请尽量当面交易以保证安全！",
		//},
	}
	systemIncorrectWords = `系统|官方|shop|www|政府|mall|mch|商户|客服|system|`

	// 移动应用配置
	DefaultMoAppConf = valueobject.MoAppConf{
		// 应用名称
		AppName: DefaultRegistry.RegistryData["PlatformName"],
		// APP图标地址
		AppIcon: DefaultRegistry.RegistryData["RetailSiteLogo"],
		// 描述
		Description: "移动应用正在开发中",
		// 模板文件
		ShowTplPath: "app.html",
		// 安卓APP版本
		AndroidVersion: "1.0",
		// 安卓APP版发布地址
		AndroidReleaseUrl: "",
		// 苹果APP版本
		IosVersion: "1.0",
		// 苹果APP版发布地址
		IosReleaseUrl: "",
		// 微软APP版本
		WpVersion: "1.0",
		// 微软APP版发布地址
		WpReleaseUrl: "",
	}


	// 默认全局销售设置
	//DefaultGlobNumberConf = valueobject.GlobNumberConf{
	// 兑换1元所需要的积分
	//IntegralExchangeRate: 100,
	// 抵扣1元所需要的积分
	//IntegralDiscountRate: 100,
	// 消费1元产生的经验
	//ExperienceRateByOrder: 1,
	// 消费1元产生的积分
	//IntegralRateByConsumption: 1,
	// 每单额外赠送
	//IntegralBackExtra: 0,
	// 提现手续费费率
	//TakeOutCsn: 0.01,
	// 转账手续费费率
	//TransferCsn: 0.01,
	// 活动账户转为赠送可提现奖金手续费费率
	//FlowConvertCsn: 0.05,
	// 钱包账户转换手续费费率
	//PresentConvertCsn: 0.05,
	// 交易手续费类型
	//TradeCsnType: valueobject.TradeCsnTypeByFee,
	// 按交易笔数收取手续费的金额
	//TradeCsnFeeByOrder: 1, // 每笔订单最低收取1元
	// 按交易金额收取手续费的百分百
	//TradeCsnPercentByFee: 0.01, // 1%收取
	//MinTakeOutAmount:     0.01,
	// 单笔最高提现金额
	//MaxTakeOutAmount: 5000,
	//}

	DefaultGlobMchSaleConf = valueobject.GlobMchSaleConf{
		// 商户订单结算模式
		//MchOrderSettleMode: enum.MchModeSettleByRate,
		// 商户订单结算比例
		//MchOrderSettleRate: 1,
		// 商户交易单是否需上传发票
		//TradeOrderRequireTicket: false,
		// 是否启用分销模式
		FxSalesEnabled: false,
		// 返现比例,0则不返现
		CashBackPercent: 0.1,
		// 一级比例
		CashBackTg1Percent: 0.5,
		// 二级比例
		CashBackTg2Percent: 0.3,
		// 会员比例
		CashBackMemberPercent: 0.2,

		// 自动设置订单
		AutoSetupOrder: 1,
		// 订单超时分钟数
		OrderTimeOutMinute: 1440, // 24小时
		// 订单自动确认时间
		OrderConfirmAfterMinute: 10,
		// 订单超时自动收货
		OrderTimeOutReceiveHour: 168, //7天
	}

	// 默认短信接口设置
	//defaultSmsConf = map[int]*valueobject.SmsApiPerm{
	//	valueobject.SmsHttp:   {Default: true},
	//	valueobject.SmsAli:    {},
	//	valueobject.Sms253Com: {},
	//}
)
