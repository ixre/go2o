/**
 * Copyright 2014 @ 56x.net.
 * name :
 * author : jarryliu
 * date : 2013-12-12 16:53
 * description :
 * history :
 */

package merchant

import (
	"github.com/ixre/go2o/core/domain/interface/merchant/shop"
	"github.com/ixre/go2o/core/domain/interface/merchant/staff"
	"github.com/ixre/go2o/core/domain/interface/merchant/user"
	"github.com/ixre/go2o/core/domain/interface/merchant/wholesaler"
	"github.com/ixre/go2o/core/domain/interface/wallet"
	"github.com/ixre/go2o/core/infrastructure/domain"
)

const (
	KindAccountCharge           = 1
	KindAccountCarry            = 2
	KindAccountPresent          = 3
	KindAccountTakePayment      = 4
	KindAccountTransferToMember = 5

	// KindＭachTakeOutToBankCard 商户提现
	KindＭachTakeOutToBankCard = 100
	// KindＭachTakOutRefund 商户提现失败返还给会员
	KindＭachTakOutRefund = 101
)

const (
	// FSelfSales 是否自营
	FSelfSales = 1 << iota
	// FLocked 停用
	FLocked
	// 已认证
	FAuthenticated
)

type (
	// IMerchantAggregateRoot 商户接口
	//todo: 实现商户等级,商户的品牌
	IMerchantAggregateRoot interface {
		domain.IAggregateRoot
		GetValue() Merchant
		// ContainFlag 判断是否包含标志
		ContainFlag(f int) bool
		// GrantFlag 标志赋值, 如果flag小于零, 则异或运算(去除)
		GrantFlag(flag int) error
		// Complex 获取符合的商家信息
		Complex() *ComplexMerchant
		// SetValue 设置值
		SetValue(*Merchant) error
		// Stat 获取商户的状态,判断 过期时间、判断是否停用。
		// 过期时间通常按: 试合作期,即1个月, 后面每年延长一次。或者会员付费开通。
		Stat() error
		// Lock 锁定
		Lock() error
		// Unlock 解锁
		Unlock() error
		// SelfSales 是否自营
		SelfSales() bool
		// Member 返回对应的会员编号
		Member() int64
		// Save 保存
		Save() (int64, error)
		// GetMajorHost 获取商户的域名
		GetMajorHost() string
		// BindMember 绑定会员号
		BindMember(memberId int) error
		// Account 获取商户账户
		Account() IAccount
		// EnableWholesale 启用批发
		EnableWholesale() error
		// Wholesaler 获取批发商实例
		Wholesaler() wholesaler.IWholesaler
		// UserManager 返回用户服务
		UserManager() user.IUserManager
		// ConfManager 返回设置服务
		ConfManager() IConfManager
		// TransactionManager 销售服务
		TransactionManager() IMerchantTransactionManager
		// LevelManager 获取会员等级服务
		LevelManager() ILevelManager
		// KvManager 获取键值管理器
		KvManager() IKvManager
		// ProfileManager 企业资料服务
		ProfileManager() IProfileManager
		// ApiManager API服务
		ApiManager() IApiManager
		// ShopManager 商店服务
		ShopManager() shop.IShopManager
		// MemberKvManager 获取会员键值管理器
		MemberKvManager() IKvManager
		// EmpManager 员工服务
		EmployeeManager() staff.IStaffManager
		// 消息系统管理器
		//MssManager() mss.IMssManager
	}

	// IAccount 账户
	IAccount interface {
		domain.IDomain
		// GetValue 获取账户值
		GetValue() *Account
		// Save 保存
		Save() error
		// GetBalanceLog 根据编号获取余额变动信息
		GetBalanceLog(id int) *BalanceLog
		// GetBalanceLogByOuterNo 根据号码获取余额变动信息
		//GetBalanceLogByOuterNo(outerNo string) *BalanceLog
		// SaveBalanceLog 保存余额变动信息
		SaveBalanceLog(*BalanceLog) (int, error)

		// GetWalletLog 获取钱包账户日志
		GetWalletLog(txId int64) *wallet.WalletLog

		// Carry 订单结账(商户结算),返回交易流水编号和错误
		Carry(p CarryParams) (txId int, err error)

		// Consume 消耗商户支出，例如广告费、提现等
		Consume(transactionTitle string, amount int, outerTxNo string, transactionRemark string) error

		// Freeze 账户冻结
		Freeze(d wallet.TransactionData, relateUser int64) (int, error)

		// Unfreeze 账户解冻, isRefundBalance 是否退回余额
		Unfreeze(d wallet.TransactionData, isRefundBalance bool, relateUser int64) error

		// Adjust 客服调整
		Adjust(title string, amount int, remark string, relateUser int64) error

		// RequestWithdrawal 申请提现(只支持钱包),drawType：提现方式,返回info_id,交易号 及错误
		RequestWithdrawal(w *wallet.WithdrawTransaction) (int, string, error)

		// ReviewWithdrawal 提现审核
		ReviewWithdrawal(transactionId int, pass bool, reason string) error

		// CompleteTransaction 完成交易(打款),outerTransactionNo为外部交易号
		CompleteTransaction(transactionId int, outerTransactionNo string) error

		// TransferToMember todo: 以下需要重构或移除
		// 转到会员账户
		TransferToMember(amount int) error

		// TransferToMember1 商户积分转会员积分
		TransferToMember1(amount float32) error

		// RequestInvoice 申请发票,返回发票申请ID和错误
		RequestInvoice(amount int, remark string) (int, error)
	}

	// 订单参数
	CarryParams struct {
		// 是否先冻结
		Freeze bool
		// 外部订单号,非订单添加前缀，如:XT:100000
		OuterTxNo string
		// 订单金额(含交易费)
		Amount int
		// 交易费
		TransactionFee int
		// 退款金额
		RefundAmount int
		// 交易描述,如：订单结算
		TransactionTitle string
		// 交易备注,如：洗衣液
		TransactionRemark string
		// 关联的外部用户编号,可为空
		OuterTxUid int
	}
	IMerchantManager interface {
		// GetMerchantByMemberId 获取会员关联的商户
		GetMerchantByMemberId(memberId int) IMerchantAggregateRoot
	}

	// 商户
	ComplexMerchant struct {
		Id int32
		// 关联的会员编号,作为结算账户
		MemberId int64
		// 用户
		Username string
		// 密码
		Pwd string
		// 商户名称
		Name string
		// 是否自营
		SelfSales int32
		Flag      int
		// 商户等级
		Level int32
		// 标志
		Logo    string
		Address string
		// 电话
		Telephone string
		// 公司名称
		CompanyName string
		// 省
		Province int32
		// 市
		City int32
		// 区
		District int32
		// 是否启用
		Enabled int32
		// 过期时间
		ExpiresTime int64
		// 注册时间
		JoinTime int64
		Status   int
		// 更新时间
		UpdateTime int64
		// 登录时间
		LoginTime int64
		// 最后登录时间
		LastLoginTime int64
	}

	// MchMerchant 商户
	Merchant struct {
		// 编号
		Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" auto:"yes" bson:"id"`
		// 会员编号
		MemberId int `json:"memberId" db:"member_id" gorm:"column:member_id" bson:"memberId"`
		// 登录用户
		Username string `json:"username" db:"username" gorm:"column:username" bson:"username"`
		// 登录密码
		Password string `json:"password" db:"password" gorm:"column:password" bson:"password"`
		// 邮箱地址
		MailAddr string `json:"mailAddr" db:"mail_addr" gorm:"column:mail_addr" bson:"mailAddr"`
		// 加密盐
		Salt string `json:"salt" db:"salt" gorm:"column:salt" bson:"salt"`
		// 名称
		MchName string `json:"mchName" db:"mch_name" gorm:"column:mch_name" bson:"mchName"`
		// 全称
		FullName string `json:"fullName" db:"full_name" gorm:"column:full_name" bson:"fullName"`
		// 是否自营
		IsSelf int16 `json:"isSelf" db:"is_self" gorm:"column:is_self" bson:"isSelf"`
		// 标志
		Flag int `json:"flag" db:"flag" gorm:"column:flag" bson:"flag"`
		// 商户等级
		Level int `json:"level" db:"level" gorm:"column:level" bson:"level"`
		// 所在省
		Province int `json:"province" db:"province" gorm:"column:province" bson:"province"`
		// 所在市
		City int `json:"city" db:"city" gorm:"column:city" bson:"city"`
		// 所在区
		District int `json:"district" db:"district" gorm:"column:district" bson:"district"`
		// 公司地址
		Address string `json:"address" db:"address" gorm:"column:address" bson:"address"`
		// 标志
		Logo string `json:"logo" db:"logo" gorm:"column:logo" bson:"logo"`
		// 公司电话
		Tel string `json:"tel" db:"tel" gorm:"column:tel" bson:"tel"`
		// 状态: 0:未审核 1:已开通  2:停用  3: 关闭
		Status int16 `json:"status" db:"status" gorm:"column:status" bson:"status"`
		// 过期时间
		ExpiresTime int `json:"expiresTime" db:"expires_time" gorm:"column:expires_time" bson:"expiresTime"`
		// 最后登录时间
		LastLoginTime int `json:"lastLoginTime" db:"last_login_time" gorm:"column:last_login_time" bson:"lastLoginTime"`
		// 创建时间
		CreateTime int `json:"createTime" db:"create_time" gorm:"column:create_time" bson:"createTime"`
	}

	// 商户余额日志
	BalanceLog struct {
		// 编号
		Id int32 `db:"id" pk:"yes" auto:"yes"`
		// 商户编号
		MchId int64 `db:"mch_id"`
		// 日志类型
		Kind int `db:"kind"`
		// 标题
		Title string `db:"title"`
		// 外部订单号
		OuterNo string `db:"outer_no"`
		// 金额
		Amount int64 `db:"amount"`
		// 手续费
		CsnAmount int64 `db:"csn_amount"`
		// 状态
		State int `db:"state"`
		// 创建时间
		CreateTime int64 `db:"create_time"`
		// 更新时间
		UpdateTime int64 `db:"update_time"`
	}

	// 商户每日报表
	MchDayChart struct {
		// 编号
		Id int32 `db:"id" pk:"yes" auto:"yes"`
		// 商户编号
		MchId int64 `db:"mch_id"`
		// 新增订单数量
		OrderNumber int `db:"order_number"`
		// 订单额
		OrderAmount float32 `db:"order_amount"`
		// 购物会员数
		BuyerNumber int `db:"buyer_number"`
		// 支付单数量
		PaidNumber int `db:"paid_number"`
		// 支付总金额
		PaidAmount float32 `db:"paid_amount"`
		// 完成订单数
		CompleteOrders int `db:"complete_orders"`
		// 入帐金额
		InAmount float32 `db:"in_amount"`
		// 线下订单数量
		OfflineOrders int `db:"offline_orders"`
		// 线下订单金额
		OfflineAmount float32 `db:"offline_amount"`
		// 日期
		Date int64 `db:"date"`
		// 日期字符串
		DateStr string `db:"date_str"`
		// 更新时间
		UpdateTime int64 `db:"update_time"`
	}
)

func (m Merchant) TableName() string {
	return "mch_merchant"
}

// MchAccount 商户账户
type Account struct {
	// 商户编号
	MchId int `json:"mchId" db:"mch_id" gorm:"column:mch_id" pk:"yes" auto:"yes" bson:"mchId"`
	// 余额
	Balance int `json:"balance" db:"balance" gorm:"column:balance" bson:"balance"`
	// 冻结金额
	FreezeAmount int `json:"freezeAmount" db:"freeze_amount" gorm:"column:freeze_amount" bson:"freezeAmount"`
	// 待入账金额
	AwaitAmount int `json:"awaitAmount" db:"await_amount" gorm:"column:await_amount" bson:"awaitAmount"`
	// 平台赠送金额
	PresentAmount int `json:"presentAmount" db:"present_amount" gorm:"column:present_amount" bson:"presentAmount"`
	// 累计销售总额
	SalesAmount int `json:"salesAmount" db:"sales_amount" gorm:"column:sales_amount" bson:"salesAmount"`
	// 累计退款金额
	RefundAmount int `json:"refundAmount" db:"refund_amount" gorm:"column:refund_amount" bson:"refundAmount"`
	// 已提取金额
	WithdrawalAmount int `json:"takeAmount" db:"take_amount" gorm:"column:take_amount" bson:"takeAmount"`
	// 线下销售金额
	OfflineSales int `json:"offlineSales" db:"offline_sales" gorm:"column:offline_sales" bson:"offlineSales"`
	// 可开票金额
	InvoiceableAmount int `json:"invoiceableAmount" db:"invoiceable_amount" gorm:"column:invoiceable_amount" bson:"invoiceableAmount"`
	// 更新时间
	UpdateTime int `json:"updateTime" db:"update_time" gorm:"column:update_time" bson:"updateTime"`
}

func (m Account) TableName() string {
	return "mch_account"
}
