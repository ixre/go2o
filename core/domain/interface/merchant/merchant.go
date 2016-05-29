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
	"go2o/core/domain/interface/merchant/mss"
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/domain/interface/merchant/user"
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

		// 返回对应的会员编号
		Member() int

		// 保存
		Save() (int, error)

		// 获取商户的域名
		GetMajorHost() string

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
		MssManager() mss.IMssManager
	}

	//合作商
	Merchant struct {
		Id int `db:"id" pk:"yes" auto:"yes"`
		// 关联的会员编号,作为结算账户
		MemberId int    `db:"member_id"`
		Usr      string `db:"usr"`
		Pwd      string `db:"pwd"`
		Name     string `db:"name"`
		// 商户等级
		Level    int   `db:"level"`
		Logo     string `db:"logo"`
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
)
