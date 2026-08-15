/**
 * Copyright 2015 @ 56x.net.
 * name : member_conf
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package merchant

type (
	IConfManager interface {
		// 获取销售配置
		GetSaleConf() SaleConf
		// 保存销售配置
		SaveSaleConf(*SaleConf) error
		// 使用系统的配置并保存
		UseGlobSaleConf() error
		// 保存客户分组
		SaveMchBuyerGroup(v *MchBuyerGroupSetting) (int32, error)
		// 获取商户的全部客户分组
		SelectBuyerGroup() []*BuyerGroup
		// 根据分组编号获取分组设置
		GetGroupByGroupId(groupId int32) *MchBuyerGroupSetting
		// 获取所有的交易设置
		GetAllTradeConf_() []*TradeConf
		// 根据交易类型获取交易设置
		GetTradeConf(tradeType int) *TradeConf
		// 保存交易设置
		SaveTradeConf([]*TradeConf) error
		// 获取结算设置
		GetSettleConf() *SettleConf
		// 保存结算设置
		SaveSettleConf(*SettleConf) error
	}

	// 商户客户分组设置
	MchBuyerGroupSetting struct {
		//编号
		ID int32 `db:"id" pk:"yes" auto:"yes"`
		//商家编号
		MerchantId int64 `db:"mch_id"`
		//客户分组编号
		GroupId int64 `db:"group_id"`
		//分组别名
		Alias string `db:"alias"`
		// 是否启用零售
		EnableRetail int32 `db:"enable_retail"`
		// 是否启用批发
		EnableWholesale int32 `db:"enable_wholesale"`
		// 批发返点周期
		RebatePeriod int32 `db:"rebate_period"`
	}

	// 全局客户分组
	BuyerGroup struct {
		//编号
		GroupId int64
		//分组别名
		Name string
		// 是否启用零售
		EnableRetail bool
		// 是否启用批发
		EnableWholesale bool
		// 批发返点周期
		RebatePeriod int
	}
)

// MchSaleConf 商户销售设置
type SaleConf struct {
	// MchId
	MchId int `json:"mchId" db:"mch_id" gorm:"column:mch_id" pk:"yes" auto:"yes" bson:"mchId"`
	// 是否启用分销,0:不启用, 1:启用
	FxSales int `json:"fxSales" db:"fx_sales" gorm:"column:fx_sales" bson:"fxSales"`
	// 反现比例,0则不返现
	CbPercent float64 `json:"cbPercent" db:"cb_percent" gorm:"column:cb_percent" bson:"cbPercent"`
	// 一级比例
	CbTg1Percent float64 `json:"cbTg1Percent" db:"cb_tg1_percent" gorm:"column:cb_tg1_percent" bson:"cbTg1Percent"`
	// 二级比例
	CbTg2Percent float64 `json:"cbTg2Percent" db:"cb_tg2_percent" gorm:"column:cb_tg2_percent" bson:"cbTg2Percent"`
	// 会员比例
	CbMemberPercent float64 `json:"cbMemberPercent" db:"cb_member_percent" gorm:"column:cb_member_percent" bson:"cbMemberPercent"`
	// 开启自动设置订单
	OaOpen int `json:"oaOpen" db:"oa_open" gorm:"column:oa_open" bson:"oaOpen"`
	// 订单超时取消（分钟）
	OaTimeoutMinute int `json:"oaTimeoutMinute" db:"oa_timeout_minute" gorm:"column:oa_timeout_minute" bson:"oaTimeoutMinute"`
	// 订单自动确认（分钟）
	OaConfirmMinute int `json:"oaConfirmMinute" db:"oa_confirm_minute" gorm:"column:oa_confirm_minute" bson:"oaConfirmMinute"`
	// 超时自动收货（小时）
	OaReceiveHour int `json:"oaReceiveHour" db:"oa_receive_hour" gorm:"column:oa_receive_hour" bson:"oaReceiveHour"`
	// UpdateTime
	UpdateTime int `json:"updateTime" db:"update_time" gorm:"column:update_time" bson:"updateTime"`
	//IntegralBackNum         int     `db:"ib_num"`                         // 每一元返多少积分
	//IntegralBackExtra       int     `db:"ib_extra"`                       // 每单额外赠送
	//TakeOutCsn                float32 `db:"apply_csn"`                      // 提现手续费费率
	//TransferCsn                float32 `db:"trans_csn"`                      // 转账手续费费率
	//FlowConvertCsn          float32 `db:"flow_convert_csn"`               // 活动账户转为赠送可提现奖金手续费费率
	//PresentConvertCsn       float32 `db:"present_convert_csn"`            // 钱包账户转换手续费费率

}

func (m SaleConf) TableName() string {
	return "mch_sale_conf"
}

// MchSettleConf 商户结算设置
type SettleConf struct {
	// 编号
	Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" auto:"yes" bson:"id"`
	// 商户编号
	MchId int `json:"mchId" db:"mch_id" gorm:"column:mch_id" bson:"mchId"`
	// 订单交易费率
	OrderTxRate float64 `json:"orderTxRate" db:"order_tx_rate" gorm:"column:order_tx_rate" bson:"orderTxRate"`
	// 其他服务手续费比例
	OtherTxRate float64 `json:"otherTxRate" db:"other_tx_rate" gorm:"column:other_tx_rate" bson:"otherTxRate"`
	// 结算子商户号
	SubMchNo string `json:"subMchNo" db:"sub_mch_no" gorm:"column:sub_mch_no" bson:"subMchNo"`
	// 创建时间
	UpdateTime int `json:"updateTime" db:"update_time" gorm:"column:update_time" bson:"updateTime"`
}

func (m SettleConf) TableName() string {
	return "mch_settle_conf"
}
