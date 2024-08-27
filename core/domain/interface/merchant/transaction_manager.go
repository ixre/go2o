package merchant

const (
	// 账单状态:待生成
	BillStatusPending BillStatus = 0
	// 账单状态:已生成
	BillStatusGenerated BillStatus = 1
	// 账单状态:已复核
	BillStatusReviewed BillStatus = 2
	// 账单状态:已结算
	BillStatusSettled BillStatus = 3
)

type (
	// 账单金额类型
	BillAmountType int
	// 账单状态
	BillStatus int
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
	// GetCurrentBill 获取当前月份的账单
	GetCurrentBill() *MerchantBill
	// GetBillByTime 获取指定月份的账单
	GetBillByTime(billTime int) *MerchantBill
	// AdjustBillShopAmount 调整账单商城金额
	AdjustBillAmount(amountType BillAmountType, amount int, txFee int) error
	// GenerateBill 生成当前月份的账单
	GenerateBill() error
	// ReviewBill 审核账单
	ReviewBill(billId int, reviewerId int) error
	// SettleBill 结算账单
	SettleBill(billId int) error
}

// 账单结算事件
type BillSettledEvent struct {
	// 商户编号
	MchId int
	// 账单
	Bill *MerchantBill
}

// MerchantBill 商户月度账单
type MerchantBill struct {
	// 编号
	Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" auto:"yes" bson:"id"`
	// 商户编号
	MchId int `json:"mchId" db:"mch_id" gorm:"column:mch_id" bson:"mchId"`
	// 账单时间
	BillTime int `json:"billTime" db:"bill_time" gorm:"column:bill_time" bson:"billTime"`
	// 月份: 例:202408
	BillMonth string `json:"billMonth" db:"bill_month" gorm:"column:bill_month" bson:"billMonth"`
	// 账单开始时间
	StartTime int `json:"startTime" db:"start_time" gorm:"column:start_time" bson:"startTime"`
	// 账单结束时间
	EndTime int `json:"endTime" db:"end_time" gorm:"column:end_time" bson:"endTime"`
	// 商城订单数量
	ShopOrderCount int `json:"shopOrderCount" db:"shop_order_count" gorm:"column:shop_order_count" bson:"shopOrderCount"`
	// 线下订单数量
	StoreOrderCount int `json:"storeOrderCount" db:"store_order_count" gorm:"column:store_order_count" bson:"storeOrderCount"`
	// 商城总金额
	ShopTotalAmount int `json:"shopTotalAmount" db:"shop_total_amount" gorm:"column:shop_total_amount" bson:"shopTotalAmount"`
	// 线下总金额
	StoreTotalAmount int `json:"storeTotalAmount" db:"store_total_amount" gorm:"column:store_total_amount" bson:"storeTotalAmount"`
	// 其他订单总数量
	OtherOrderCount int `json:"otherOrderCount" db:"other_order_count" gorm:"column:other_order_count" bson:"otherOrderCount"`
	// 其他订单总金额
	OtherTotalAmount int `json:"otherTotalAmount" db:"other_total_amount" gorm:"column:other_total_amount" bson:"otherTotalAmount"`
	// 交易费
	TotalTxFee int `json:"totalTxFee" db:"total_tx_fee" gorm:"column:total_tx_fee" bson:"totalTxFee"`
	// 账单状态:  0: 待生成 1: 已生成 2: 已复核  3: 已结算
	Status int `json:"status" db:"status" gorm:"column:status" bson:"status"`
	// 审核人编号
	ReviewerId int `json:"reviewerId" db:"reviewer_id" gorm:"column:reviewer_id" bson:"reviewerId"`
	// 审核人名称
	ReviewerName string `json:"reviewerName" db:"reviewer_name" gorm:"column:reviewer_name" bson:"reviewerName"`
	// 审核时间
	ReviewTime int `json:"reviewTime" db:"review_time" gorm:"column:review_time" bson:"reviewTime"`
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
