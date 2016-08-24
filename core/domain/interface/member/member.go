/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-09 09:49
 * description :
 * history :
 */

package member

const (
	StateStopped = 0 //已停用
	StateOk      = 1 //正常
	BankNoLock   = 0
	BankLocked   = 1
)

const (
	// 收藏店铺
	FavTypeShop = iota + 1
	// 收藏商品
	FavTypeGoods
)

type (
	IMember interface {
		// 获取聚合根编号
		GetAggregateRootId() int

		// 会员资料服务
		Profile() IProfileManager

		// 会员收藏服务
		Favorite() IFavoriteManager

		// 礼品卡服务
		GiftCard() IGiftCardManager

		// 邀请管理
		Invitation() IInvitationManager

		// 获取值
		GetValue() Member

		// 设置值
		SetValue(*Member) error

		// 获取账户
		GetAccount() IAccount

		// 发送验证码,传入操作及消息类型,并返回验证码,及错误
		SendCheckCode(operation string, mssType int) (string, error)

		// 对比验证码
		CompareCode(code string) error

		// 锁定会员
		Lock() error

		// 解锁会员
		Unlock() error

		// 获取关联的会员
		GetRelation() *Relation

		// 更新会员绑定
		SaveRelation(r *Relation) error

		// 更改邀请人
		ChangeReferees(memberId int) error

		// 增加经验值
		AddExp(exp int) error

		// 获取等级
		GetLevel() *Level

		// 更换用户名
		ChangeUsr(string) error

		// 保存
		Save() (int, error)
	}

	// 会员资料服务
	IProfileManager interface {
		// 获取资料
		GetProfile() Profile

		// 保存资料
		SaveProfile(v *Profile) error

		// 资料是否完善
		ProfileCompleted() bool

		// 修改密码,旧密码可为空; 传入原始密码。
		ModifyPassword(newPwd, oldPwd string) error

		// 修改交易密码，旧密码可为空; 传入原始密码。
		ModifyTradePassword(newPwd, oldPwd string) error

		// 获取提现银行信息
		GetBank() BankInfo

		// 保存提现银行信息,保存后将锁定
		SaveBank(*BankInfo) error

		// 解锁提现银行卡信息
		UnlockBank() error

		// 实名认证信息
		GetTrustedInfo() TrustedInfo

		// 保存实名认证信息
		SaveTrustedInfo(v *TrustedInfo) error

		// 审核实名认证,若重复审核将返回错误
		ReviewTrustedInfo(pass bool, remark string) error

		// 创建配送地址
		CreateDeliver(*DeliverAddress) IDeliverAddress

		// 获取配送地址
		GetDeliverAddress() []IDeliverAddress

		// 获取配送地址
		GetDeliver(int) IDeliverAddress

		// 获取默认收货地址
		GetDefaultAddress() IDeliverAddress

		// 删除配送地址
		DeleteDeliver(int) error
	}

	// 收藏服务
	IFavoriteManager interface {
		// 收藏
		Favorite(favType, referId int) error

		// 是否已收藏
		Favored(favType, referId int) bool

		// 取消收藏
		Cancel(favType, referId int) error

		// 收藏商品
		FavoriteGoods(goodsId int) error

		// 取消收藏商品
		CancelGoodsFavorite(goodsId int) error

		// 收藏店铺
		FavoriteShop(shopId int) error

		// 取消收藏店铺
		CancelShopFavorite(shopId int) error

		// 商品是否已收藏
		GoodsFavored(goodsId int) bool

		// 商店是否已收藏
		ShopFavored(shopId int) bool
	}

	Member struct {
		// 编号
		Id int `db:"id" auto:"yes" pk:"yes"`
		// 用户名
		Usr string `db:"usr"`
		// 密码
		Pwd string `db:"Pwd"`
		// 交易密码
		TradePwd string `db:"trade_pwd"`
		// 经验值
		Exp int `db:"exp"`
		// 等级
		Level int `db:"level"`
		// 邀请码
		InvitationCode string `db:"invitation_code"`
		// 注册来源
		RegFrom string `db:"reg_from"`
		// 注册IP
		RegIp string `db:"reg_ip"`
		// 注册时间
		RegTime int64 `db:"reg_time"`
		// 校验码
		CheckCode string `db:"check_code"`
		// 校验码过期时间
		CheckExpires int64 `db:"check_expires"`
		// 状态
		State int `db:"state"`
		// 最后登录时间
		LastLoginTime int64 `db:"last_login_time"`
		// 更新时间
		UpdateTime int64 `db:"update_time"`
		// 动态令牌，用于登录或API调用
		DynamicToken string `db:"-"`
		// 超时时间
		TimeoutTime int64 `db:"-"`
	}

	// 会员资料
	Profile struct {
		//会员编号
		MemberId int `db:"member_id" pk:"yes" auto:"no"`
		//姓名
		Name string `db:"name"`
		//头像
		Avatar string `db:"avatar"`
		//性别
		Sex int `db:"sex"`
		//生日
		BirthDay string `db:"birthday"`
		//电话
		Phone string `db:"phone"`
		//地址
		Address string `db:"address"`
		//即时通讯
		Im string `db:"im"`
		//电子邮件
		Email string `db:"email"`
		// 省
		Province int `db:"province"`
		// 市
		City int `db:"city"`
		// 区
		District int `db:"district"`
		//备注
		Remark string `db:"remark"`

		// 扩展1
		Ext1 string `db:"ext_1"`
		// 扩展2
		Ext2 string `db:"ext_2"`
		// 扩展3
		Ext3 string `db:"ext_3"`
		// 扩展4
		Ext4 string `db:"ext_4"`
		// 扩展5
		Ext5 string `db:"ext_4"`
		// 扩展6
		Ext6 string `db:"ext_4"`
		//更新时间
		UpdateTime int64 `db:"update_time"`
	}

	//会员关联表
	Relation struct {
		MemberId int `db:"member_id" pk:"yes"`
		//会员卡号
		CardId string `db:"card_id"`
		//推荐人（会员）
		RefereesId int `db:"invi_member_id"`
		// 会员关系字符串
		ReferStr string `db:"refer_str"`
		//注册关联商户编号
		RegisterMerchantId int `db:"reg_merchant_id"`
	}

	// 实名认证信息
	TrustedInfo struct {
		//会员编号
		MemberId int `db:"member_id" pk:"yes"`
		//真实姓名
		RealName string `db:"real_name"`
		//身份证号码
		CardId string `db:"card_id"`
		//认证图片、身份证、人与身份证的图像等
		TrustImage string `db:"trust_image"`
		//是否审核通过
		Reviewed int `db:"reviewed"`
		//审核时间
		ReviewTime int64 `db:"review_time"`
		//审核备注
		Remark string `db:"remark"`
		//更新时间
		UpdateTime int64 `db:"update_time"`
	}

	// 银行卡信息,因为重要且非频繁更新的数据
	// 所以需要用IsLocked来标记是否锁定
	BankInfo struct {
		//会员编号
		MemberId int `db:"member_id" pk:"yes"`
		//名称
		Name string `db:"name"`
		//账号
		Account string `db:"account"`
		//账户名
		AccountName string `db:"account_name"`
		//支行网点
		Network string `db:"network"`
		//状态
		State int `db:"state"`
		//是否锁定
		IsLocked int `db:"is_locked"`
		//更新时间
		UpdateTime int64 `db:"update_time"`
	}

	// 收藏
	Favorite struct {
		// 编号
		Id int `db:"id"`
		// 会员编号
		MemberId int `db:"member_id"`
		// 收藏类型
		FavType int `db:"fav_type"`
		// 引用编号
		ReferId int `db:"refer_id"`
		// 收藏时间
		UpdateTime int64 `db:"update_time"`
	}

	// 收货地址
	IDeliverAddress interface {
		GetDomainId() int
		GetValue() DeliverAddress
		SetValue(*DeliverAddress) error
		Save() (int, error)
	}

	// 收货地址
	DeliverAddress struct {
		//编号
		Id int `db:"id" pk:"yes" auto:"yes"`
		//会员编号
		MemberId int `db:"member_id"`
		//收货人
		RealName string `db:"real_name"`
		//电话
		Phone string `db:"phone"`
		//省
		Province int `db:"province"`
		//市
		City int `db:"city"`
		//区
		District int `db:"district"`
		//地区(省市区连接)
		Area string `db:"area"`
		//地址
		Address string `db:"address"`
		//是否默认
		IsDefault int `db:"is_default"`
	}
)

func (this BankInfo) Right() bool {
	return len(this.Name) > 0 && len(this.Account) > 0 &&
		len(this.AccountName) > 0
}

func (this BankInfo) Locked() bool {
	return this.IsLocked == BankLocked
}
