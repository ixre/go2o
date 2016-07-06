/**
 * Copyright 2015 @ z3q.net.
 * name : value_rep.go
 * author : jarryliu
 * date : 2016-05-27 15:28
 * description :
 * history :
 */
package valueobject

var (
	TradeCsnTypeByOrder = 1 //按订单笔数收取手续费
	TradeCsnTypeByFee   = 2 //按交易金额收取手续费
)

var (
	SmsAli     = 1 //阿里大鱼
	SmsNetEasy = 2 //网易
	SmsIndex   = []int{
		SmsAli,
		SmsNetEasy,
	}
	SmsTextMap = map[int]string{
		SmsAli:     "阿里大鱼",
		SmsNetEasy: "网易",
	}
)

type (
	// 平台设置
	PlatformConf struct {
		// 平台名称
		Name string
		// 标志
		Logo string
		// 允许商户创建商品分类
		MchGoodsCategory bool
		// 允许商户创建页面分类
		MchPageCategory bool
	}

	// 微信API设置
	//todo: ??? 应在应用层
	WxApiConfig struct {
		/**===== 微信公众平台设置 =====**/

		//APP ID
		AppId string
		//APP 密钥
		AppSecret string
		//通信密钥
		MpToken string
		//通信AES KEY
		MpAesKey string
		//原始ID
		OriId string

		/**===== 用于微信支付 =====**/

		//商户编号
		MchId string
		//商户接口密钥
		MchApiKey string
		//微信支付的证书路径(上传)
		MchCertPath string
		//微信支付的证书公钥路径(上传)
		MchCertKeyPath string
		//是否启用红包功能
		RedPackEnabled bool
		//红包金额限制
		RedPackAmountLimit float32
		//红包每日数量限制
		RedPackDayTimeLimit int

		//MchPayNotifyPath string //微信支付异步通知的路径
	}

	// 注册权限设置
	RegisterPerm struct {
		// 注册模式,等于member.RegisterMode
		RegisterMode int
		// 是否允许匿名注册
		AnonymousRegistered bool
		// 注册提示
		Notice string
		// 用户条款内容
		Licence string
		// 注册回调页
		CallBackUrl string
	}

	// 全局数值设置
	GlobNumberConf struct {
		// 兑换1元所需要的积分
		IntegralExchangeRate int
		// 提现手续费费率
		ApplyCsn float32
		// 转账手续费费率
		TransCsn float32
		// 活动账户转为赠送可提现奖金手续费费率
		FlowConvertCsn float32
		// 赠送账户转换手续费费率
		PresentConvertCsn float32
		// 每一元返多少积分
		IntegralBackNum int
		// 每单额外赠送
		IntegralBackExtra int
		// 交易手续费类型
		TradeCsnType int
		// 按交易笔数收取手续费的金额
		TradeCsnFeeByOrder float32
		// 按交易金额收取手续费的百分百
		TradeCsnPercentByFee float32
		// 最低提现金额
		MinApplyAmount float32
		// 单笔最高提现金额
		MaxApplyAmount float32
	}

	// 全局商户销售设置
	GlobMchSaleConf struct {
		// 是否启用分销模式
		FxSalesEnabled bool
		// 返现比例,0则不返现
		CashBackPercent float32
		// 一级比例
		CashBackTg1Percent float32
		// 二级比例
		CashBackTg2Percent float32
		// 会员比例
		CashBackMemberPercent float32

		// 自动设置订单
		AutoSetupOrder int
		// 订单超时分钟数
		OrderTimeOutMinute int
		// 订单自动确认时间
		OrderConfirmAfterMinute int
		// 订单超时自动收货
		OrderTimeOutReceiveHour int
	}

	IValueRep interface {
		// 获取微信接口配置
		GetWxApiConfig() WxApiConfig

		// 保存微信接口配置
		SaveWxApiConfig(v *WxApiConfig) error

		// 获取注册权限
		GetRegisterPerm() RegisterPerm

		// 保存注册权限
		SaveRegisterPerm(v *RegisterPerm) error

		// 获取全局系统数值设置
		GetGlobNumberConf() GlobNumberConf

		// 保存全局系统数值设置
		SaveGlobNumberConf(v *GlobNumberConf) error

		// 获取平台设置
		GetPlatformConf() PlatformConf

		// 保存平台设置
		SavePlatformConf(v *PlatformConf) error

		// 获取全局商户销售设置
		GetGlobMchSaleConf() GlobMchSaleConf

		// 保存全局商户销售设置
		SaveGlobMchSaleConf(v *GlobMchSaleConf) error

		// 获取短信设置
		GetSmsApiSet() SmsApiSet

		// 保存短信API
		SaveSmsApiPerm(provider int, s *SmsApiPerm) error

		// 获取默认的短信API,返回API提供商编号及API信息
		GetDefaultSmsApiPerm() (int, *SmsApiPerm)

		// 获取下级区域
		GetChildAreas(id int) []*Area

		// 获取地区名称
		GetAreaNames(id []int) []string
	}

	// 短信接口
	SmsApiPerm struct {
		ApiKey    string //接口编号
		ApiSecret string //接口密钥
		Default   bool   //是否默认的接口使用
	}
	// 短信接口设置
	SmsApiSet map[int]*SmsApiPerm

	//http://www.stats.gov.cn/tjsj/tjbz/xzqhdm/
	// 区域,中国行政区划
	Area struct {
		Code   int    `db:"code" json:"code"`
		Parent int    `db:"parent" json:"parent"`
		Name   string `db:"name" json:"name"`
	}
)
