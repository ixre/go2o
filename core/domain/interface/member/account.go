/**
 * Copyright 2015 @ z3q.net.
 * name : account
 * author : jarryliu
 * date : 2015-07-24 08:48
 * description :
 * history :
 */
package member

const (
	// 余额账户
	AccountBalance = 1
	// 积分账户
	AccountIntegral = 2
	// 钱包账户
	AccountWallet = 3
	// 流通金账户
	AccountFlow = 4
)

const (
	// 用户充值
	ChargeByUser int32 = 1
	// 系统自动充值
	ChargeBySystem int32 = 2
	// 客服充值
	ChargeByService int32 = 3
	// 退款充值
	ChargeByRefund int32 = 4
)

const (
	// 自定义的业务类型
	KindMine int32 = 30
	// 会员充值
	KindBalanceCharge int32 = 1
	// 系统充值
	KindBalanceSystemCharge int32 = 2
	// 支付抵扣
	KindBalanceDiscount int32 = 3
	// 退款
	KindBalanceRefund int32 = 4
	// 转入
	KindBalanceTransferIn int32 = 5
	// 转出
	KindBalanceTransferOut int32 = 6
	// 失效
	KindBalanceExpired int32 = 7
	// 冻结
	KindBalanceFreeze int32 = 8
	// 解冻
	KindBalanceUnfreeze int32 = 9
	// 客服调整
	KindAdjust int32 = 10

	// 客服充值
	KindBalanceServiceCharge int32 = 15
	// 客服扣减
	KindBalanceServiceDiscount int32 = 16
)

const (
	// 赠送金额
	KindWalletAdd int32 = 1
	// 抵扣奖金
	KindWalletDiscount int32 = 2
	// 转入
	KindWalletTransferIn int32 = 5
	// 转出
	KindWalletTransferOut int32 = 6
	// 失效
	KindWalletExpired int32 = 7
	// 冻结
	KindWalletFreeze int32 = 8
	// 解冻
	KindWalletUnfreeze int32 = 9
	// 客服调整
	KindWalletAdjust int32 = 10
	// 提现到余额
	KindWalletTakeOutToBalance int32 = 11
	// 提现到银行卡(人工提现)
	KindWalletTakeOutToBankCard int32 = 12
	// 提现到第三方
	KindWalletTakeOutToThirdPart int32 = 13
	// 提现退还到银行卡
	KindWalletTakeOutRefund int32 = 14
	// 支付单退款
	KindWalletPaymentRefund int32 = 15

	// 客服赠送
	KindWalletServiceAdd int32 = 21
	// 客服扣减
	KindWalletServiceDiscount int32 = 22
)

const (
	KindGrow = 7 // 增利

	//KindCommission = 9 // 手续费

	// 赠送
	//KindBalancePresent = 3

	// 流通账户
	KindBalanceFlow int32 = 4 // 账户流通

	// 提现
	//KindBalanceApplyCash = 11
	// 转账
	KindBalanceTransfer int32 = 12
	StatusOK                  = 1
)

const (
	// 赠送
	TypeIntegralPresent = 1
	// 积分抵扣
	TypeIntegralDiscount = 2
	// 积分冻结
	TypeIntegralFreeze = 3
	// 积分解冻
	TypeIntegralUnfreeze = 4
	// 购物赠送
	TypeIntegralShoppingPresent = 5
	// 支付抵扣
	TypeIntegralPaymentDiscount = 6
)

type (
	IAccount interface {
		// 获取领域对象编号
		GetDomainId() int64

		// 获取账户值
		GetValue() *Account

		// 保存
		Save() (int64, error)

		// 设置优先(默认)支付方式, account 为账户类型
		SetPriorityPay(account int, enabled bool) error

		// 根据编号获取余额变动信息
		GetBalanceInfo(id int32) *BalanceInfo

		// 根据号码获取余额变动信息
		// GetBalanceInfoByNo(no string) *BalanceInfo

		// 保存余额变动信息
		SaveBalanceInfo(*BalanceInfo) (int32, error)

		// 获取钱包账户日志
		GetWalletLog(id int32) *MWalletLog

		// 充值,客服操作时,需提供操作人(relateUser)
		//ChargeForBalance(chargeType int32, title string, outerNo string, amount float32, relateUser int64) error

		// 扣减余额
		DiscountBalance(title string, outerNo string, amount float32, relateUser int64) error

		// 冻结余额
		Freeze(title string, outerNo string, amount float32, relateUser int64) error

		// 解冻金额
		Unfreeze(title string, outerNo string, amount float32, relateUser int64) error

		// 扣减奖金,mustLargeZero是否必须大于0, 赠送金额存在扣为负数的情况
		DiscountWallet(title string, outerNo string, amount float32,
			relateUser int64, mustLargeZero bool) error

		// 冻结赠送金额
		FreezeWallet(title string, outerNo string, amount float32, relateUser int64) error

		// 解冻赠送金额
		UnfreezeWallet(title string, outerNo string, amount float32, relateUser int64) error

		// 支付单抵扣消费,tradeNo为支付单单号
		PaymentDiscount(tradeNo string, amount float32, remark string) error

		//　增加积分
		AddIntegral(iType int, outerNo string, value int64, remark string) error

		// 积分抵扣
		IntegralDiscount(logType int, outerNo string, value int64, remark string) error

		// 冻结积分,当new为true不扣除积分,反之扣除积分
		FreezesIntegral(value int64, new bool, remark string) error

		// 解冻积分
		UnfreezesIntegral(value int64, remark string) error

		// 退款
		RequestBackBalance(backType int, title string, amount float32) error

		// 完成退款
		FinishBackBalance(id int32, tradeNo string) error

		// 申请提现,applyType：提现方式,返回info_id,交易号 及错误
		RequestTakeOut(takeKind int32, title string, amount float32, commission float32) (int32, string, error)

		// 确认提现
		ConfirmTakeOut(id int32, pass bool, remark string) error

		// 完成提现
		FinishTakeOut(id int32, tradeNo string) error

		// 将冻结金额标记为失效
		FreezeExpired(accountKind int, amount float32, remark string) error

		// 转账
		TransferAccount(accountKind int, toMember int64, amount float32,
			csnRate float32, remark string) error

		// 接收转账
		ReceiveTransfer(accountKind int, fromMember int64, tradeNo string,
			amount float32, remark string) error

		// 退款
		Refund(accountKind int, kind int32, title string, outerNo string, amount float32, relateUser int64) error

		// 充值
		Charge(account int32, kind int32, title, outerNo string, amount float32, relateUser int64) error

		// 客服调整
		Adjust(account int, title string, amount float32, remark string, relateUser int64) error

		// 赠送金额,客服操作时,需提供操作人(relateUser)
		//ChargeForPresent(title string, outerNo string, amount float32, relateUser int64) error

		// 赠送金额(指定业务类型)
		//ChargePresentByKind(kind int32, title string, outerNo string, amount float32, relateUser int64) error

		// 转账余额到其他账户
		TransferBalance(kind int32, amount float32, tradeNo string, toTitle, fromTitle string) error

		// 转账返利账户,kind为转账类型，如 KindBalanceTransfer等
		// commission手续费
		TransferWallet(kind int32, amount float32, commission float32, tradeNo string,
			toTitle string, fromTitle string) error

		// 转账活动账户,kind为转账类型，如 KindBalanceTransfer等
		// commission手续费
		TransferFlow(kind int32, amount float32, commission float32, tradeNo string,
			toTitle string, fromTitle string) error

		// 将活动金转给其他人
		TransferFlowTo(memberId int64, kind int32, amount float32, commission float32,
			tradeNo string, toTitle string, fromTitle string) error
	}

	// 余额变动信息
	BalanceInfo struct {
		Id       int64  `db:"id" auto:"yes" pk:"yes"`
		MemberId int64  `db:"member_id"`
		TradeNo  string `db:"trade_no"`
		Kind     int32  `db:"kind"`
		Type     int    `db:"type"`
		Title    string `db:"title"`
		// 金额
		Amount float32 `db:"amount"`
		// 手续费
		CsnAmount float32 `db:"csn_amount"`
		// 引用编号
		RefId      int64 `db:"ref_id"`
		State      int   `db:"state"`
		CreateTime int64 `db:"create_time"`
		UpdateTime int64 `db:"update_time"`
	}

	// 账户值对象
	Account struct {
		// 会员编号
		MemberId int64 `db:"member_id" pk:"yes"`
		// 积分
		Integral int64 `db:"integral"`
		// 不可用积分
		FreezeIntegral int64 `db:"freeze_integral"`
		// 余额
		Balance float32 `db:"balance"`
		// 不可用余额
		FreezeBalance float32 `db:"freeze_balance"`
		// 失效的账户余额
		ExpiredBalance float32 `db:"expired_balance"`
		//奖金账户余额
		WalletBalance float32 `db:"wallet_balance"`
		//冻结赠送金额
		FreezeWallet float32 `db:"freeze_wallet"`
		//失效的赠送金额
		ExpiredPresent float32 `db:"expired_wallet"`
		//总赠送金额
		TotalPresentFee float32 `db:"total_wallet_amount"`
		//流动账户余额
		FlowBalance float32 `db:"flow_balance"`
		//当前理财账户余额
		GrowBalance float32 `db:"grow_balance"`
		//理财总投资金额,不含收益
		GrowAmount float32 `db:"grow_amount"`
		//当前收益金额
		GrowEarnings float32 `db:"grow_earnings"`
		//累积收益金额
		GrowTotalEarnings float32 `db:"grow_total_earnings"`
		//总消费金额
		TotalExpense float32 `db:"total_expense"`
		//总充值金额
		TotalCharge float32 `db:"total_charge"`
		//总支付额
		TotalPay float32 `db:"total_pay"`
		// 优先(默认)支付选项
		PriorityPay int `db:"priority_pay"`
		//更新时间
		UpdateTime int64 `db:"update_time"`
	}

	// 积分记录
	IntegralLog struct {
		// 编号
		Id int64 `db:"id" pk:"yes" auto:"yes"`
		// 会员编号
		MemberId int64 `db:"member_id"`
		// 类型
		Type int `db:"type"`
		// 关联的编号
		OuterNo string `db:"outer_no"`
		// 积分值
		Value int64 `db:"value"`
		// 备注
		Remark string `db:"remark"`
		// 创建时间
		CreateTime int64 `db:"create_time"`
	}

	// 余额日志
	BalanceLog struct {
		Id       int64  `db:"id" auto:"yes" pk:"yes"`
		MemberId int64  `db:"member_id"`
		OuterNo  string `db:"outer_no"`
		// 业务类型
		BusinessKind int32 `db:"kind"`

		Title string `db:"title"`
		// 金额
		Amount float32 `db:"amount"`
		// 手续费
		CsnFee float32 `db:"csn_fee"`
		// 关联操作人,仅在客服操作时,记录操作人
		RelateUser int64 `db:"rel_user"`
		// 状态
		State int32 `db:"state"`
		// 备注
		Remark string `db:"remark"`
		// 创建时间
		CreateTime int64 `db:"create_time"`
		// 更新时间
		UpdateTime int64 `db:"update_time"`
	}

	// 钱包账户日志
	MWalletLog struct {
		ID int64 `db:"id" auto:"yes" pk:"yes"`
		// 会员编号
		MemberId int64 `db:"member_id"`
		// 外部单号
		OuterNo string `db:"outer_no"`
		// 业务类型
		BusinessKind int32 `db:"kind"`
		// 标题
		Title string `db:"title"`
		// 金额
		Amount float32 `db:"amount"`
		// 手续费
		CsnFee float32 `db:"csn_fee"`
		// 关联操作人,仅在客服操作时,记录操作人
		RelateUser int64 `db:"rel_user"`
		// 状态
		State int32 `db:"state"`
		// 备注
		Remark string `db:"remark"`
		// 创建时间
		CreateTime int64 `db:"create_time"`
		// 更新时间
		UpdateTime int64 `db:"update_time"`
	}
)
