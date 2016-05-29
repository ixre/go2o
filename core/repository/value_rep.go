/**
 * Copyright 2015 @ z3q.net.
 * name : value_rep
 * author : jarryliu
 * date : 2016-05-27 15:32
 * description :
 * history :
 */
package repository

import (
	"errors"
	"github.com/jsix/gof/db"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/valueobject"
)

var _ valueobject.IValueRep = new(valueRep)

var (
	// 默认注册权限设置
	defaultRegisterPerm = valueobject.RegisterPerm{
		RegisterMode:        member.RegisterModeNormal,
		AnonymousRegistered: true,
	}

	// 默认全局销售设置
	defaultGlobSaleConf = valueobject.GlobSaleConf{
		// 提现手续费费率
		ApplyCsn: 0.01,
		// 转账手续费费率
		TransCsn: 0.01,
		// 活动账户转为赠送可提现奖金手续费费率
		FlowConvertCsn: 0.05,
		// 赠送账户转换手续费费率
		PresentConvertCsn: 0.05,
		// 每一元返多少积分
		IntegralBackNum: 1,
		// 每单额外赠送
		IntegralBackExtra: 0,
		// 交易手续费类型
		TradeCsnType: valueobject.TradeCsnTypeByFee,
		// 按交易笔数收取手续费的金额
		TradeCsnFeeByOrder: 1, // 每笔订单最低收取1元
		// 按交易金额收取手续费的百分百
		TradeCsnPercentByFee: 0.01, // 1%收取

	}

	defaultGlobMchSaleConf = valueobject.GlobMerchantSaleConf{
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
		OrderTimeOutReceiveHour: 168, //c7天
	}
)

type valueRep struct {
	db.Connector
	_wxConf          *valueobject.WxApiConfig
	_rpConf          *valueobject.RegisterPerm
	_globSaleConf    *valueobject.GlobSaleConf
	_globMchSaleConf *valueobject.GlobMerchantSaleConf
}

func NewValueRep(conn db.Connector) valueobject.IValueRep {
	return &valueRep{
		Connector: conn,
	}
}

// 获取微信接口配置
func (this *valueRep) GetWxApiConfig() *valueobject.WxApiConfig {
	if this._wxConf == nil {
		this._wxConf = &valueobject.WxApiConfig{}
		unMarshalFromFile("conf/core/wx_api", this._wxConf)
	}
	return this._wxConf
}

// 保存微信接口配置
func (this *valueRep) SaveWxApiConfig(v *valueobject.WxApiConfig) error {
	if v != nil {
		//todo: 检查证书文件是否存在
		this._wxConf = v
		return marshalToFile("conf/core/wx_api", this._wxConf)
	}
	return errors.New("nil value")
}

// 获取注册权限
func (this *valueRep) GetRegisterPerm() *valueobject.RegisterPerm {
	if this._rpConf == nil {
		v := defaultRegisterPerm
		this._rpConf = &v
		unMarshalFromFile("conf/core/register_perm", this._rpConf)
	}
	return this._rpConf
}

// 保存注册权限
func (this *valueRep) SaveRegisterPerm(v *valueobject.RegisterPerm) error {
	if v != nil {
		this._rpConf = v
		return marshalToFile("conf/core/register_perm", this._rpConf)
	}
	return nil
}

// 获取全局系统销售设置
func (this *valueRep) GetGlobSaleConf() *valueobject.GlobSaleConf {
	if this._globSaleConf == nil {
		v := defaultGlobSaleConf
		this._globSaleConf = &v
		unMarshalFromFile("conf/core/sale_conf", this._globSaleConf)
	}
	return this._globSaleConf
}

// 保存全局系统销售设置
func (this *valueRep) SaveGlobSaleConf(v *valueobject.GlobSaleConf) error {
	if v != nil {
		this._globSaleConf = v
		return marshalToFile("conf/core/sale_conf", this._globSaleConf)
	}
	return nil
}

// 获取全局商户销售设置
func (this *valueRep) GetGlobMerchantSaleConf() *valueobject.GlobMerchantSaleConf {
	if this._globMchSaleConf == nil {
		v := defaultGlobMchSaleConf
		this._globMchSaleConf = &v
		unMarshalFromFile("conf/core/mch_sale_conf", this._globMchSaleConf)
	}
	return this._globMchSaleConf
}

// 保存全局商户销售设置
func (this *valueRep) SaveGlobMerchantSaleConf(v *valueobject.GlobMerchantSaleConf) error {
	if v != nil {
		this._globMchSaleConf = v
		return marshalToFile("conf/core/mch_sale_conf", this._globMchSaleConf)
	}
	return nil
}
