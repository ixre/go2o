/**
 * Copyright 2015 @ z3q.net.
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
	}

	// 销售设置为商户填写,同时可以恢复默认
	SaleConf struct {
		// 合作商编号
		MerchantId int64 `db:"mch_id" auto:"no" pk:"yes"`
		// 是否启用分销模式
		FxSalesEnabled int `db:"fx_sales"`
		// 返现比例,0则不返现
		CashBackPercent float32 `db:"cb_percent"`
		// 一级比例
		CashBackTg1Percent float32 `db:"cb_tg1_percent"`
		// 二级比例
		CashBackTg2Percent float32 `db:"cb_tg2_percent"`
		// 会员比例
		CashBackMemberPercent float32 `db:"cb_member_percent"`
		// 自动设置订单
		AutoSetupOrder int `db:"oa_open"`
		// 订单超时分钟数
		OrderTimeOutMinute int `db:"oa_timeout_minute"`
		// 订单自动确认时间
		OrderConfirmAfterMinute int `db:"oa_confirm_minute"`
		// 订单超时自动收货
		OrderTimeOutReceiveHour int `db:"oa_receive_hour"`

		//IntegralBackNum         int     `db:"ib_num"`                         // 每一元返多少积分
		//IntegralBackExtra       int     `db:"ib_extra"`                       // 每单额外赠送
		//TakeOutCsn                float32 `db:"apply_csn"`                      // 提现手续费费率
		//TransferCsn                float32 `db:"trans_csn"`                      // 转账手续费费率
		//FlowConvertCsn          float32 `db:"flow_convert_csn"`               // 活动账户转为赠送可提现奖金手续费费率
		//PresentConvertCsn       float32 `db:"present_convert_csn"`            // 赠送账户转换手续费费率
	}
)
