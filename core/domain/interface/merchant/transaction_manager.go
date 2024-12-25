/**
 * Copyright (C) 2007-2024 fze.NET, All rights reserved.
 *
 * name: transaction_manager.go
 * author: jarrysix (jarrysix#gmail.com)
 * date: 2024-08-20 22:39:44
 * description: 商户交易服务
 * history:
 */

package merchant

import "github.com/ixre/go2o/core/infrastructure/domain"

const (
	// 账单类型:日度账单
	BillTypeDaily = 1
	// 账单类型:月度账单
	BillTypeMonthly = 2
)

const (
	// 账单状态:待生成
	BillStatusPending BillStatus = 1
	// 账单状态:待核对
	BillStatusWaitConfirm BillStatus = 2
	// 账单状态:待复核
	BillStatusWaitReview BillStatus = 3
	// 账单状态:已复核
	BillStatusReviewed BillStatus = 4
)

type (
	// 账单金额类型
	BillAmountType int
	// 账单状态
	BillStatus int
)

const (
	// 未结算
	BillNoSettlement = 0
	// 待结算
	BillWaitSettlement = 1
	// 已结算
	BillSettlemented = 2
)

const (
	// 无需结算
	BillResultNone = 0
	// 结算在途
	BillResultPending = 1
	// 结算失败
	BillResultFailed = 2
	// 结算到帐
	BillResultSuccess = 3
)

const (
	// 账单金额类型:商城
	BillAmountTypeShop = 0
	// 账单金额类型:线下
	BillAmountTypeStore = 1
	// 账单金额类型:其他
	BillAmountTypeOther = 2
)

// 商户交易服务
type IMerchantTransactionManager interface {
	// 计算交易费用,返回交易费及错误
	MathTransactionFee(tradeType int, amount int) (int, error)
	// GetSettlementType 获取系统设置结算周期
	GetSettlementPeriod() int
	// GetCurrentDailyBill 获取当前账单
	GetCurrentDailyBill() IBillDomain
	// GetBill 获取指定账单
	GetBill(billId int) IBillDomain
	// GetBillByTime 获取指定月份的账单
	GetBillByTime(billTime int) IBillDomain
	// GenerateMonthlyBill 生成月度账单,如果为日结，则账单会自动复核。
	// 月度账单的生成时间需保证上月的日账单已复核完毕，
	// 建议为手动生成每月账单,如未生成，则每月3日定时生成。
	GenerateMonthlyBill(year, month int) error
}

// IBillDomain 账单领域
type IBillDomain interface {
	domain.IDomain
	// 获取值
	Value() *MerchantBill
	// Update 更新账单金额
	UpdateAmount() error
	// UpdateBillAmount 调整账单金额,如果amount为负数,则表示退款
	UpdateBillAmount(amount int, txFee int) error
	// Generate 生成账单
	Generate() error
	// Confirm 确认账单,按日结算的账单不需要商户确认
	Confirm() error
	// Review 审核账单
	Review(reviewerId int, remark string) error
	// Settle 结算账单
	Settle() error
	// UpdateSettleInfo 更新结算信息,settleTxNo 结算单号, message 错误信息
	UpdateSettleInfo(spCode string, settleTxNo string, settleResult int, message string) error
}

// MerchantBillSettleEvent 账单结算事件
type MerchantBillSettleEvent struct {
	// 商户编号
	MchId int
	// 商户名称
	MchName string
	// 结算子商户号
	SubMerchantNo string
	// 账单
	Bill *MerchantBill
	// 结算备注
	SettleRemark string
}

// MerchantBill 商户账单
type MerchantBill struct {
	// 编号
	Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" auto:"yes" bson:"id"`
	// 商户编号
	MchId int `json:"mchId" db:"mch_id" gorm:"column:mch_id" bson:"mchId"`
	// 账单类型, 1: 日账单  2: 月度账单
	BillType int `json:"billType" db:"bill_type" gorm:"column:bill_type" bson:"billType"`
	// 账单时间
	BillTime int `json:"billTime" db:"bill_time" gorm:"column:bill_time" bson:"billTime"`
	// 月份: 例:202408
	BillMonth string `json:"billMonth" db:"bill_month" gorm:"column:bill_month" bson:"billMonth"`
	// 账单开始时间
	StartTime int `json:"startTime" db:"start_time" gorm:"column:start_time" bson:"startTime"`
	// 账单结束时间
	EndTime int `json:"endTime" db:"end_time" gorm:"column:end_time" bson:"endTime"`
	// 交易笔数
	TxCount int `json:"txCount" db:"tx_count" gorm:"column:tx_count" bson:"txCount"`
	// 交易总金额
	TxAmount int `json:"txAmount" db:"tx_amount" gorm:"column:tx_amount" bson:"txAmount"`
	// 交易手续费
	TxFee int `json:"txFee" db:"tx_fee" gorm:"column:tx_fee" bson:"txFee"`
	// 交易退款金额
	RefundAmount int `json:"refundAmount" db:"refund_amount" gorm:"column:refund_amount" bson:"refundAmount"`
	// 实际账单金额
	FinalAmount int `json:"finalAmount" db:"final_amount" gorm:"column:final_amount" bson:"finalAmount"`
	// 账单状态:  0: 待生成 1: 待确认   2: 待复核 3: 待结算  4: 已结算
	Status int `json:"status" db:"status" gorm:"column:status" bson:"status"`
	// 审核人编号
	ReviewerId int `json:"reviewerId" db:"reviewer_id" gorm:"column:reviewer_id" bson:"reviewerId"`
	// 审核人名称
	ReviewerName string `json:"reviewerName" db:"reviewer_name" gorm:"column:reviewer_name" bson:"reviewerName"`
	// 审核备注
	ReviewRemark string `json:"reviewRemark" db:"review_remark" gorm:"column:review_remark" bson:"reviewRemark"`
	// 审核时间
	ReviewTime int `json:"reviewTime" db:"review_time" gorm:"column:review_time" bson:"reviewTime"`
	// 账单备注
	BillRemark string `json:"billRemark" db:"bill_remark" gorm:"column:bill_remark" bson:"billRemark"`
	// 用户备注
	UserRemark string `json:"userRemark" db:"user_remark" gorm:"column:user_remark" bson:"userRemark"`
	// 结算状态 0: 无需结算 1: 待结算 2: 已结算
	SettleStatus int `json:"settleStatus" db:"settle_status" gorm:"column:settle_status" bson:"settleStatus"`
	// 结算通道编码
	SettleSpCode string `json:"settleSpCode" db:"settle_sp_code" gorm:"column:settle_sp_code" bson:"settleSpCode"`
	// 结算单号
	SettleTxNo string `json:"settleTxNo" db:"settle_tx_no" gorm:"column:settle_tx_no" bson:"settleTxNo"`
	// 结算结果 0: 无需结算 1: 结算在途  2: 结算失败  3: 结算到帐
	SettleResult int `json:"settleResult" db:"settle_result" gorm:"column:settle_result" bson:"settleResult"`
	// 结算备注
	SettleRemark string `json:"settleRemark" db:"settle_remark" gorm:"column:settle_remark" bson:"settleRemark"`
	// 创建时间
	CreateTime int `json:"createTime" db:"create_time" gorm:"column:create_time" bson:"createTime"`
	// 账单生成时间
	BuildTime int `json:"buildTime" db:"build_time" gorm:"column:build_time" bson:"buildTime"`
	// 更新时间
	UpdateTime int `json:"updateTime" db:"update_time" gorm:"column:update_time" bson:"updateTime"`
}

func (m MerchantBill) TableName() string {
	return "mch_bill"
}
