package invoice

import (
	"reflect"

	"github.com/ixre/go2o/core/domain"
	"github.com/ixre/go2o/core/infrastructure/fw"
)

// 发票用户类型
type TenantType int
type IssueStatus int

const (
	// 系统
	TenantSystem TenantType = 0
	// 用户
	TenantUser TenantType = 1
	// 商户
	TenantMerchant TenantType = 2
)

const (
	// 发票状态: 待开票
	IssuePending = 1
	// 发票状态: 开票完成
	IssueSuccess = 2
	// 发票状态: 开票失败
	IssueFail = 3
	// 发票状态: 撤销
	IssueRevert = 4
)

type (
	// InvoiceUserAggregateRoot 发票用户聚合根
	InvoiceUserAggregateRoot interface {
		domain.IAggregateRoot
		// TeantType 获取租户类型
		TenantType() TenantType
		// TenantUserId 获取租户用户编号
		TenantUserId() int
		// Create 创建租户
		Create() error
		// GetInvoiceTitle 获取发票抬头
		GetInvoiceTitle(id int) *InvoiceTitle
		// CreateInvoiceTitle 保存发票抬头
		CreateInvoiceTitle(title *InvoiceTitle) error
		// CreateInvoice 创建发票
		RequestInvoice(data *InvoiceRequestData) (InvoiceDomain, error)
		// GetInvoice 获取发票
		GetInvoice(id int) InvoiceDomain
	}

	InvoiceDomain interface {
		domain.IDomain
		// GetValue 获取发票记录
		GetValue() *InvoiceRecord
		// GetItems 获取发票明细
		GetItems() []*InvoiceItem
		// Issue 开具发票,更新发票图片
		Issue(picture string) error
		// Issue 发票开具失败
		IssueFail(reason string) error
		// SendMail 发送发票到邮件中
		SendMail(mail string) error
		// Revert 撤销发票
		Revert(reason string) error
		// Save 保存发票
		Save() error
	}
)

// 发票申请数据
type InvoiceRequestData struct {
	// 关联单号
	OuterNo string `json:"outerNo"`
	// 开票人ID
	IssueTenantId int `json:"issueTenantId"`
	// 发票抬头
	TitleId int `json:"titleId"`
	// 接收邮箱
	ReceiveEmail string `json:"receiveEmail"`
	// 发票内容
	Subject string `json:"subject"`
	// 备注
	Remark string `json:"remark"`
	// 开票项目
	Items []*InvoiceItem `json:"items"`
}

var _ domain.IValueObject = new(InvoiceItem)

// IInvoiceTenantRepo 发票租户仓储接口
type IInvoiceTenantRepo interface {
	fw.Repository[InvoiceTenant]
	// Title 获取发票抬头仓储接口
	Title() IInvoiceTitleRepo
	// Records 获取发票记录仓储接口
	Records() IInvoiceRecordRepo
	// Items 获取发票项目仓储接口
	Items() IInvoiceItemRepo
	// GetTenant 获取租户
	GetTenant(id int) InvoiceUserAggregateRoot
	// CreateTenant 创建租户
	CreateTenant(v *InvoiceTenant) InvoiceUserAggregateRoot
}

// IInvoiceTitlesRepo 发票抬头仓储接口
type IInvoiceTitleRepo interface {
	fw.Repository[InvoiceTitle]
}

// IInvoiceRecordRepo 发票仓储接口
type IInvoiceRecordRepo interface {
	fw.Repository[InvoiceRecord]
}

// IInvoiceItemRepo 发票项目仓储接口
type IInvoiceItemRepo interface {
	fw.Repository[InvoiceItem]
}

// InvoiceTenant 发票租户
type InvoiceTenant struct {
	// 编号
	Id int `json:"id" db:"id" gorm:"column:id" bson:"id"`
	// 租户类型,1:会员  2:商户
	TenantType int `json:"tenantType" db:"tenant_type" gorm:"column:tenant_type" bson:"tenantType"`
	// 租户用户编号
	TenantUid int `json:"tenantUid" db:"tenant_uid" gorm:"column:tenant_uid" bson:"tenantUid"`
	// 创建时间
	CreateTime int `json:"createTime" db:"create_time" gorm:"column:create_time" bson:"createTime"`
}

func (i InvoiceTenant) TableName() string {
	return "invoice_tenant"
}

// InvoiceTitle 发票抬头
type InvoiceTitle struct {
	// 编号
	Id int `json:"id" db:"id" gorm:"column:id" bson:"id"`
	// 租户编号
	TenantId int `json:"tenantId" db:"tenant_id" gorm:"column:tenant_id" bson:"tenantId"`
	// 发票类型: 1:增值税普通发票 2:增值税专用发票 3:形式发票
	InvoiceType int `json:"invoiceType" db:"invoice_type" gorm:"column:invoice_type" bson:"invoiceType"`
	// 开具类型, 1: 个人 2:企业
	IssueType int `json:"issueType" db:"issue_type" gorm:"column:issue_type" bson:"issueType"`
	// HeaderName
	TitleName string `json:"titleName" db:"title_name" gorm:"column:title_name" bson:"titleName"`
	// 纳税人识别号
	TaxCode string `json:"taxCode" db:"tax_code" gorm:"column:tax_code" bson:"taxCode"`
	// 注册场所地址
	SignAddress string `json:"signAddress" db:"sign_address" gorm:"column:sign_address" bson:"signAddress"`
	// 注册固定电话
	SignTel string `json:"signTel" db:"sign_tel" gorm:"column:sign_tel" bson:"signTel"`
	// 基本户开户银行名
	BankName string `json:"bankName" db:"bank_name" gorm:"column:bank_name" bson:"bankName"`
	// 基本户开户账号
	BankAccount string `json:"bankAccount" db:"bank_account" gorm:"column:bank_account" bson:"bankAccount"`
	// 备注
	Remarks string `json:"remarks" db:"remarks" gorm:"column:remarks" bson:"remarks"`
	// 是否默认
	IsDefault int `json:"isDefault" db:"is_default" gorm:"column:is_default" bson:"isDefault"`
	// 创建时间
	CreateTime int `json:"createTime" db:"create_time" gorm:"column:create_time" bson:"createTime"`
}

func (i InvoiceTitle) TableName() string {
	return "invoice_title"
}

// InvoiceRecord 发票
type InvoiceRecord struct {
	// 编号
	Id int `json:"id" db:"id" gorm:"column:id" bson:"id"`
	// 发票代码
	InvoiceCode string `json:"invoiceCode" db:"invoice_code" gorm:"column:invoice_code" bson:"invoiceCode"`
	// 发票号码
	InvoiceNo string `json:"invoiceNo" db:"invoice_no" gorm:"column:invoice_no" bson:"invoiceNo"`
	// 租户编号
	TenantId int `json:"tenantId" db:"tenant_id" gorm:"column:tenant_id" bson:"tenantId"`
	// 开具租户编号
	IssueTenantId int `json:"issueTenantId" db:"issue_tenant_id" gorm:"column:issue_tenant_id" bson:"issueTenantId"`
	// 发票类型: 1:增值税普通发票 2:增值税专用发票 3:形式发票
	InvoiceType int `json:"invoiceType" db:"invoice_type" gorm:"column:invoice_type" bson:"invoiceType"`
	// 开具类型, 1: 个人 2:企业
	IssueType int `json:"issueType" db:"issue_type" gorm:"column:issue_type" bson:"issueType"`
	// 销售方名称
	SellerName string `json:"sellerName" db:"seller_name" gorm:"column:seller_name" bson:"sellerName"`
	// 销售方纳税人识别号
	SellerTaxCode string `json:"sellerTaxCode" db:"seller_tax_code" gorm:"column:seller_tax_code" bson:"sellerTaxCode"`
	// 买方名称
	PurchaserName string `json:"purchaserName" db:"purchaser_name" gorm:"column:purchaser_name" bson:"purchaserName"`
	// 买方纳税人识别号
	PurchaserTaxCode string `json:"purchaserTaxCode" db:"purchaser_tax_code" gorm:"column:purchaser_tax_code" bson:"purchaserTaxCode"`
	// 发票内容
	InvoiceSubject string `json:"invoiceSubject" db:"invoice_subject" gorm:"column:invoice_subject" bson:"invoiceSubject"`
	// 合计金额
	InvoiceAmount float64 `json:"invoiceAmount" db:"invoice_amount" gorm:"column:invoice_amount" bson:"invoiceAmount"`
	// 合计税额
	TaxAmount float64 `json:"taxAmount" db:"tax_amount" gorm:"column:tax_amount" bson:"taxAmount"`
	// 备注
	Remark string `json:"remark" db:"remark" gorm:"column:remark" bson:"remark"`
	// 开具备注/开票失败备注
	IssueRemark string `json:"issueRemark" db:"issue_remark" gorm:"column:issue_remark" bson:"issueRemark"`
	// 发票图片
	InvoicePic string `json:"invoicePic" db:"invoice_pic" gorm:"column:invoice_pic" bson:"invoicePic"`
	// 发票接收邮箱地址
	ReceiveEmail string `json:"receiveEmail" db:"receive_email" gorm:"column:receive_email" bson:"receiveEmail"`
	// 发票状态,1:待开票 2:开票完成 3:未通过
	InvoiceStatus int `json:"invoiceStatus" db:"invoice_status" gorm:"column:invoice_status" bson:"invoiceStatus"`
	// 开票时间
	InvoiceTime int `json:"invoiceTime" db:"invoice_time" gorm:"column:invoice_time" bson:"invoiceTime"`
	// 创建时间
	CreateTime int `json:"createTime" db:"create_time" gorm:"column:create_time" bson:"createTime"`
	// 更新时间
	UpdateTime int `json:"updateTime" db:"update_time" gorm:"column:update_time" bson:"updateTime"`
}

func (i InvoiceRecord) TableName() string {
	return "invoice_record"
}

// InvoiceItem 发票项目
type InvoiceItem struct {
	// 编号
	Id int `json:"id" db:"id" gorm:"column:id" bson:"id"`
	// 发票编号
	InvoiceId int `json:"invoiceId" db:"invoice_id" gorm:"column:invoice_id" bson:"invoiceId"`
	// 项目名称
	ItemName string `json:"itemName" db:"item_name" gorm:"column:item_name" bson:"itemName"`
	// 项目规格
	ItemSpec string `json:"itemSpec" db:"item_spec" gorm:"column:item_spec" bson:"itemSpec"`
	// 价格
	Price float64 `json:"price" db:"price" gorm:"column:price" bson:"price"`
	// 数量
	Quantity int `json:"quantity" db:"quantity" gorm:"column:quantity" bson:"quantity"`
	// 税率
	TaxRate float64 `json:"taxRate" db:"tax_rate" gorm:"column:tax_rate" bson:"taxRate"`
	// 计量单位
	Unit string `json:"unit" db:"unit" gorm:"column:unit" bson:"unit"`
	// 总金额
	Amount float64 `json:"amount" db:"amount" gorm:"column:amount" bson:"amount"`
	// 税额
	TaxAmount float64 `json:"taxAmount" db:"tax_amount" gorm:"column:tax_amount" bson:"taxAmount"`
	// 创建时间
	CreateTime int `json:"createTime" db:"create_time" gorm:"column:create_time" bson:"createTime"`
	// 更新时间
	UpdateTime int `json:"updateTime" db:"update_time" gorm:"column:update_time" bson:"updateTime"`
}

func (i InvoiceItem) TableName() string {
	return "invoice_item"
}

// Equal implements domain.IValueObject.
func (i *InvoiceItem) Equal(v interface{}) bool {
	return reflect.DeepEqual(i, v)
}
