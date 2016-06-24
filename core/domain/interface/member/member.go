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
		ProfileManager() IProfileManager

		// 会员收藏服务
		FavoriteManager() IFavoriteManager

		// 获取值
		GetValue() Member

		// 邀请管理
		Invitation() IInvitationManager

		// 设置值
		SetValue(*Member) error

		// 获取账户
		GetAccount() IAccount

		// 锁定会员
		Lock() error

		// 解锁会员
		Unlock() error

		//　保存积分记录
		SaveIntegralLog(*IntegralLog) error

		// 增加经验值
		AddExp(exp int) error

		// 获取等级
		GetLevel() *Level

		//　增加积分
		// todo:merchantId 不需要
		AddIntegral(merchantId int, backType int, integral int, log string) error

		// 获取关联的会员
		GetRelation() *MemberRelation

		// 更新会员绑定
		SaveRelation(r *MemberRelation) error

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

		// 修改密码,旧密码可为空
		ModifyPassword(newPwd, oldPwd string) error

		// 修改交易密码，旧密码可为空
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
		CreateDeliver(*DeliverAddress) (IDeliver, error)

		// 获取配送地址
		GetDeliverAddress() []IDeliver

		// 获取配送地址
		GetDeliver(int) IDeliver

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
		// 状态
		State int `db:"state"`
		// 最后登陆时间
		LastLoginTime int64 `db:"last_login_time"`
		// 更新时间
		UpdateTime int64 `db:"update_time"`
		// 动态令牌，用于登陆或API调用
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

	// 实名认证信息
	TrustedInfo struct {
		//会员编号
		MemberId int `db:"member_id" pk:"yes"`
		//真实姓名
		RealName string `db:"real_name"`
		//身份证号码
		BodyNumber string `db:"body_number"`
		//认证图片、身份证、人与身份证的图像等
		TrustImage string `db:"trust_image"`
		//是否处理
		IsHandle int `db:"is_handle"`
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
)

func (this BankInfo) Right() bool {
	return len(this.Name) > 0 && len(this.Account) > 0 &&
		len(this.AccountName) > 0
}

func (this BankInfo) Locked() bool {
	return this.IsLocked == BankLocked
}
