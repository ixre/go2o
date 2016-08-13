/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-12 16:53
 * description :
 * history :
 */

package merchant

import (
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/domain/interface/merchant/user"
)

const (
	KindAccountSettleOrder      = 1
	KindAccountTransferToMember = 2
)

type (
	// 商户接口
	//todo: 实现商户等级,商户的品牌
	IMerchant interface {
		GetAggregateRootId() int

		GetValue() Merchant

		SetValue(*Merchant) error

		// 获取商户的状态,判断 过期时间、判断是否停用。
		// 过期时间通常按: 试合作期,即1个月, 后面每年延长一次。或者会员付费开通。
		Stat() error

		// 是否自营
		SelfSales() bool

		// 返回对应的会员编号
		Member() int

		// 保存
		Save() (int, error)

		// 获取商户的域名
		GetMajorHost() string

		// 获取商户账户
		Account() IAccount

		// 返回用户服务
		UserManager() user.IUserManager

		// 返回设置服务
		ConfManager() IConfManager

		// 获取会员等级服务
		LevelManager() ILevelManager

		// 获取键值管理器
		KvManager() IKvManager

		// 企业资料服务
		ProfileManager() IProfileManager

		// API服务
		ApiManager() IApiManager

		// 商店服务
		ShopManager() shop.IShopManager

		// 获取会员键值管理器
		MemberKvManager() IKvManager

		// 消息系统管理器
		//MssManager() mss.IMssManager
	}

	// 账户
	IAccount interface {
		// 获取领域对象编号
		GetDomainId() int

		// 获取账户值
		GetValue() *Account

		// 保存
		Save() error

		// 根据编号获取余额变动信息
		GetBalanceLog(id int) *BalanceLog

		// 根据号码获取余额变动信息
		GetBalanceLogByOuterNo(outerNo string) *BalanceLog

		// 保存余额变动信息
		SaveBalanceLog(*BalanceLog) (int, error)

		// 订单结账
		SettleOrder(orderNo string, amount float32, csn float32, refundAmount float32, remark string) error

		// 提现
		//todo:???

		// 转到会员账户
		TransferToMember(amount float32) error

		// 赠送
		Present(amount float32, remark string) error
	}

	//合作商
	Merchant struct {
		Id int `db:"id" pk:"yes" auto:"yes"`
		// 关联的会员编号,作为结算账户
		MemberId int `db:"member_id"`

		Usr string `db:"usr"`
		Pwd string `db:"pwd"`

		// 商户名称
		Name string `db:"name"`

		// 是否自营
		SelfSales int `db:"self_sales"`

		// 商户等级
		Level int    `db:"level"`
		Logo  string `db:"logo"`
		// 省
		Province int `db:"province"`
		// 市
		City int `db:"city"`
		// 区
		District int `db:"district"`
		// 是否启用
		Enabled int `db:"enabled"`

		ExpiresTime   int64 `db:"expires_time"`
		JoinTime      int64 `db:"join_time"`
		UpdateTime    int64 `db:"update_time"`
		LoginTime     int64 `db:"login_time"`
		LastLoginTime int64 `db:"last_login_time"`
	}

	// 商户账户表
	Account struct {
		// 商户编号
		MchId int `db:"mch_id" pk:"yes"`
		// 余额
		Balance float32 `db:"balance"`
		// 冻结金额
		FreezeAmount float32 `db:"freeze_amount"`
		// 待入账金额
		AwaitAmount float32 `db:"await_amount"`
		// 平台赠送金额
		PresentAmount float32 `db:"present_amount"`
		// 累计销售总额
		SalesAmount float32 `db:"sales_amount"`
		// 累计退款金额
		RefundAmount float32 `db:"refund_amount"`
		// 已提取金额
		TakeAmount float32 `db:"take_amount"`
		// 线下销售金额
		OfflineSales float32 `db:"offline_sales"`
		// 更新时间
		UpdateTime int64 `db:"update_time"`
	}

	// 商户余额日志
	BalanceLog struct {
		// 编号
		Id int `db:"id" pk:"yes" auto:"yes"`
		// 商户编号
		MchId int `db:"Mchid"`
		// 日志类型
		Kind int `db:"Kind"`
		// 标题
		Title string `db:"Title"`
		// 外部订单号
		OuterNo string `db:"Outerno"`
		// 金额
		Amount float32 `db:"Amount"`
		// 手续费
		CsnAmount float32 `db:"csn_amount"`
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
		Id int `db:"id" pk:"yes" auto:"yes"`
		// 商户编号
		MchId int `db:"mch_id"`
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
		Date int `db:"date"`
		// 日期字符串
		DateStr string `db:"date_str"`
		// 更新时间
		UpdateTime int64 `db:"update_time"`
	}
)
