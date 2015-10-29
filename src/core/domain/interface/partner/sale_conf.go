/**
 * Copyright 2014 @ ops.
 * name :
 * author : jarryliu
 * date : 2013-11-11 19:51
 * description :
 * history :
 */

package partner

type SaleConf struct {
	PartnerId             int     `db:"partner_id" auto:"no" pk:"yes"` // 合作商编号
	CashBackPercent       float32 `db:"cb_percent"`                    // 返现比例,0则不返现
	CashBackTg1Percent    float32 `db:"cb_tg1_percent"`                // 一级比例
	CashBackTg2Percent    float32 `db:"cb_tg2_percent"`                // 二级比例
	CashBackMemberPercent float32 `db:"cb_member_percent"`             // 会员比例
	IntegralBackNum       int     `db:"ib_num"`                        // 每一元返多少积分
	IntegralBackExtra     int     `db:"ib_extra"`                      // 每单额外赠送
	AutoSetupOrder        int     `db:"auto_setup_order"`              // 自动设置订单
	ApplyCsn              float32 `db:"apply_csn"`                     // 提现手续费费率
	TransCsn              float32 `db:"trans_csn"`                     // 转账手续费费率
	FlowConvertCsn        float32 `db:"flow_convert_csn"`              // 活动账户转为赠送可提现奖金手续费费率
	PresentConvertCsn     float32 `db:"present_convert_csn"`           // 赠送账户转换手续费费率
}
