/**
 * Copyright 2015 @ z3q.net.
 * name : default
 * author : jarryliu
 * date : 2016-07-23 11:21
 * description :
 * history :
 */
package repository

import (
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/valueobject"
)

var (
	DefaultRegistry = valueobject.Registry{
		AlertMessageForOrderReceive: "确认收货后,款项将转给商户。请在收货前确保已经商品没有损坏和缺少!",
		// 会员中心首页模板文件名称
		UCenterIndexTplFile:              "index.html",
		MemberProfileNotCompletedMessage: "您的个人资料未完善,是否立即完善?",
		MemberNotTrustedMessage:          "您尚未实名认证!",
		// 注册后赠送积分数量
		PresentIntegralNumOfRegister: 0,
		Extend: map[string]string{},
	}

	// 默认平台设置
	//todo: 默认值
	defaultPlatformConf = valueobject.PlatformConf{
		Name:                    "GO2O",
		Logo:                    "https://raw.githubusercontent.com/jsix/go2o/master/docs/mark.gif",
		MchGoodsCategory:        false,
		MchPageCategory:         false,
		EnabledMemberExperience: !true,
	}

	// 默认注册权限设置
	defaultRegisterPerm = valueobject.RegisterPerm{
		RegisterMode:        member.RegisterModeNormal,
		NeedPhone:           false,
		MustBinPhone:        false,
		NeedIm:              false,
		AnonymousRegistered: true,
		CallBackUrl:         "/auth?uc=1", //默认进入会员中心
	}

	// 默认全局销售设置
	defaultGlobNumberConf = valueobject.GlobNumberConf{
		// 兑换1元所需要的积分
		IntegralExchangeRate: 1000,
		// 消费1元产生的经验
		ExperienceRateByOrder: 1,
		// 消费1元产生的积分
		IntegralRateByOrder: 1,
		// 每单额外赠送
		IntegralBackExtra: 0,
		// 提现手续费费率
		ApplyCsn: 0.01,
		// 转账手续费费率
		TransCsn: 0.01,
		// 活动账户转为赠送可提现奖金手续费费率
		FlowConvertCsn: 0.05,
		// 赠送账户转换手续费费率
		PresentConvertCsn: 0.05,
		// 交易手续费类型
		TradeCsnType: valueobject.TradeCsnTypeByFee,
		// 按交易笔数收取手续费的金额
		TradeCsnFeeByOrder: 1, // 每笔订单最低收取1元
		// 按交易金额收取手续费的百分百
		TradeCsnPercentByFee: 0.01, // 1%收取
	}

	defaultGlobMchSaleConf = valueobject.GlobMchSaleConf{
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
		OrderTimeOutMinute: 720, // 12小时
		// 订单自动确认时间
		OrderConfirmAfterMinute: 10,
		// 订单超时自动收货
		OrderTimeOutReceiveHour: 168, //7天
	}

	// 默认短信接口设置
	defaultSmsConf = map[int]*valueobject.SmsApiPerm{
		valueobject.SmsAli:     &valueobject.SmsApiPerm{Default: true},
		valueobject.SmsNetEasy: &valueobject.SmsApiPerm{},
		valueobject.Sms253Com:  &valueobject.SmsApiPerm{},
	}
)
