/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-12 16:55
 * description :
 * history :
 */

package merchant

import (
	"errors"
	"fmt"
	"go2o/src/core/domain/interface/enum"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/merchant"
	"go2o/src/core/domain/interface/merchant/mss"
	"go2o/src/core/domain/interface/merchant/user"
	mssImpl "go2o/src/core/domain/merchant/mss"
	userImpl "go2o/src/core/domain/merchant/user"
	"go2o/src/core/infrastructure"
	"go2o/src/core/infrastructure/domain"
	"go2o/src/core/variable"
	"time"
)

var _ merchant.IMerchant = new(Merchant)

type Merchant struct {
	_value    *merchant.MerchantValue
	_saleConf *merchant.SaleConf
	_siteConf *merchant.SiteConf
	_apiInfo  *merchant.ApiInfo
	_shops    []merchant.IShop
	_host     string

	_rep       merchant.IMerchantRep
	_userRep   user.IUserRep
	_memberRep member.IMemberRep

	_userManager     user.IUserManager
	_confManager     merchant.IConfManager
	_levelManager    merchant.ILevelManager
	_kvManager       merchant.IKvManager
	_memberKvManager merchant.IKvManager
	_mssManager      mss.IMssManager
	_mssRep          mss.IMssRep
}

func NewMerchant(v *merchant.MerchantValue, rep merchant.IMerchantRep, userRep user.IUserRep,
	memberRep member.IMemberRep, mssRep mss.IMssRep) (merchant.IMerchant, error) {

	var err error

	if v == nil {
		err = errors.New("101:no such partner")
		return nil, err
	}
	if time.Now().Unix() > v.ExpiresTime {
		err = errors.New("103: partner is expires")
	}

	return &Merchant{
		_value:     v,
		_rep:       rep,
		_userRep:   userRep,
		_memberRep: memberRep,
		_mssRep:    mssRep,
	}, err
}

func (this *Merchant) clearShopCache() {
	this._shops = nil
}

func (this *Merchant) GetAggregateRootId() int {
	return this._value.Id
}
func (this *Merchant) GetValue() merchant.MerchantValue {
	return *this._value
}

func (this *Merchant) SetValue(v *merchant.MerchantValue) error {
	tv := this._value
	if v.Id == tv.Id {
		tv.Name = v.Name
		tv.Address = v.Address
		if v.LastLoginTime > 0 {
			tv.LastLoginTime = v.LastLoginTime
		}
		if v.LoginTime > 0 {
			tv.LoginTime = v.LoginTime
		}

		if len(v.Logo) != 0 {
			tv.Logo = v.Logo
		}
		tv.Phone = v.Phone
		tv.Pwd = v.Pwd
		tv.UpdateTime = time.Now().Unix()

	}
	return nil
}

// 保存
func (this *Merchant) Save() (int, error) {
	var id int = this.GetAggregateRootId()
	if id > 0 {
		return this._rep.SaveMerchant(this._value)
	}

	return this.createMerchant()
}

// 创建商户
func (this *Merchant) createMerchant() (int, error) {
	if id := this.GetAggregateRootId(); id > 0 {
		return id, nil
	}

	v := this._value
	id, err := this._rep.SaveMerchant(v)
	if err != nil {
		return id, err
	}

	//todo:事务

	// 初始化商户信息
	this._value.Id = id

	// SiteConf
	this._siteConf = &merchant.SiteConf{
		IndexTitle: "线上商店-" + v.Name,
		SubTitle:   "线上商店-" + v.Name,
		Logo:       v.Logo,
		State:      1,
		StateHtml:  "",
	}
	err = this._rep.SaveSiteConf(id, this._siteConf)
	this._siteConf.MerchantId = id

	// SaleConf
	this._saleConf = &merchant.SaleConf{
		AutoSetupOrder:  1,
		IntegralBackNum: 0,
	}
	err = this._rep.SaveSaleConf(id, this._saleConf)
	this._saleConf.MerchantId = id

	// 创建API
	this._apiInfo = &merchant.ApiInfo{
		ApiId:     domain.NewApiId(id),
		ApiSecret: domain.NewSecret(id),
		WhiteList: "*",
		Enabled:   1,
	}
	err = this._rep.SaveApiInfo(id, this._apiInfo)

	return id, err
}

// 获取商户的域名
func (this *Merchant) GetMajorHost() string {
	if len(this._host) == 0 {
		host := this._rep.GetMerchantMajorHost(this.GetAggregateRootId())
		if len(host) == 0 {
			host = fmt.Sprintf("%s.%s", this._value.Usr, infrastructure.GetApp().
				Config().GetString(variable.ServerDomain))
		}
		this._host = host
	}
	return this._host
}

// 获取销售配置
func (this *Merchant) GetSaleConf() merchant.SaleConf {
	if this._saleConf == nil {
		//10%分成
		//0.2,         #上级
		//0.1,         #上上级
		//0.8          #消费者自己
		this._saleConf = this._rep.GetSaleConf(
			this.GetAggregateRootId())

		this.verifySaleConf(this._saleConf)
	}
	return *this._saleConf
}

// 保存销售配置
func (this *Merchant) SaveSaleConf(v *merchant.SaleConf) error {

	this.GetSaleConf()

	if v.RegisterMode == merchant.ModeRegisterClosed ||
		v.RegisterMode == merchant.ModeRegisterNormal ||
		v.RegisterMode == merchant.ModeRegisterMustInvitation ||
		v.RegisterMode == merchant.ModeRegisterMustRedirect {
		this._saleConf.RegisterMode = v.RegisterMode
	} else {
		return merchant.ErrRegisterMode
	}

	if v.FlowConvertCsn < 0 || v.PresentConvertCsn < 0 ||
		v.ApplyCsn < 0 || v.TransCsn < 0 ||
		v.FlowConvertCsn > 1 || v.PresentConvertCsn > 1 ||
		v.ApplyCsn > 1 || v.TransCsn > 1 {
		return merchant.ErrSalesPercent
	}

	this.verifySaleConf(v)

	this._saleConf = v
	this._saleConf.MerchantId = this._value.Id

	return this._rep.SaveSaleConf(this.GetAggregateRootId(), this._saleConf)
}

// 注册权限验证,如果没有权限注册,返回错误
func (this *Merchant) RegisterPerm(isInvitation bool) error {
	conf := this.GetSaleConf()
	if conf.RegisterMode == merchant.ModeRegisterClosed {
		return merchant.ErrRegOff
	}
	if conf.RegisterMode == merchant.ModeRegisterMustInvitation && !isInvitation {
		return merchant.ErrRegMustInvitation
	}
	if conf.RegisterMode == merchant.ModeRegisterMustRedirect && isInvitation {
		return merchant.ErrRegOffInvitation
	}
	return nil
}

// 验证销售设置
func (this *Merchant) verifySaleConf(v *merchant.SaleConf) {
	if v.OrderTimeOutMinute <= 0 {
		v.OrderTimeOutMinute = 1440 // 一天
	}

	if v.OrderConfirmAfterMinute <= 0 {
		v.OrderConfirmAfterMinute = 60 // 一小时后自动确认
	}

	if v.OrderTimeOutReceiveHour <= 0 {
		v.OrderTimeOutReceiveHour = 7 * 24 // 7天后自动确认
	}
}

// 获取站点配置
func (this *Merchant) GetSiteConf() merchant.SiteConf {
	if this._siteConf == nil {
		this._siteConf = this._rep.GetSiteConf(this.GetAggregateRootId())
	}
	return *this._siteConf
}

// 保存站点配置
func (this *Merchant) SaveSiteConf(v *merchant.SiteConf) error {
	this._siteConf = v
	this._siteConf.MerchantId = this._value.Id
	return this._rep.SaveSiteConf(this.GetAggregateRootId(), this._siteConf)
}

// 获取API信息
func (this *Merchant) GetApiInfo() merchant.ApiInfo {
	if this._apiInfo == nil {
		this._apiInfo = this._rep.GetApiInfo(this.GetAggregateRootId())
	}
	return *this._apiInfo
}

// 保存API信息
func (this *Merchant) SaveApiInfo(v *merchant.ApiInfo) error {
	this._apiInfo = v
	this._apiInfo.MerchantId = this._value.Id
	return this._rep.SaveApiInfo(this.GetAggregateRootId(), this._apiInfo)
}

// 新建商店
func (this *Merchant) CreateShop(v *merchant.ValueShop) merchant.IShop {
	v.CreateTime = time.Now().Unix()
	v.MerchantId = this.GetAggregateRootId()
	return newShop(this, v, this._rep)
}

// 获取所有商店
func (this *Merchant) GetShops() []merchant.IShop {
	if this._shops == nil {
		shops := this._rep.GetShopsOfMerchant(this.GetAggregateRootId())
		this._shops = make([]merchant.IShop, len(shops))
		for i, v := range shops {
			this._shops[i] = this.CreateShop(v)
		}
	}

	return this._shops
}

// 获取营业中的商店
func (this *Merchant) GetBusinessInShops() []merchant.IShop {
	var list []merchant.IShop = make([]merchant.IShop, 0)
	for _, v := range this._shops {
		if v.GetValue().State == enum.ShopBusinessIn {
			list = append(list, v)
		}
	}
	return list
}

// 获取商店
func (this *Merchant) GetShop(shopId int) merchant.IShop {
	//	v := this.rep.GetValueShop(this.GetAggregateRootId(), shopId)
	//	if v == nil {
	//		return nil
	//	}
	//	return this.CreateShop(v)
	shops := this.GetShops()

	for _, v := range shops {
		if v.GetValue().Id == shopId {
			return v
		}
	}
	return nil
}

// 删除门店
func (this *Merchant) DeleteShop(shopId int) error {
	//todo : 检测订单数量
	return this._rep.DeleteShop(this.GetAggregateRootId(), shopId)
}

// 返回用户服务
func (this *Merchant) UserManager() user.IUserManager {
	if this._userManager == nil {
		this._userManager = userImpl.NewUserManager(
			this.GetAggregateRootId(),
			this._userRep)
	}
	return this._userManager
}

// 返回设置服务
func (this *Merchant) ConfManager() merchant.IConfManager {
	if this._confManager == nil {
		this._confManager = &ConfManager{
			_merchantId: this.GetAggregateRootId(),
			_rep:        this._rep,
		}
	}
	return this._confManager
}

// 获取会员管理服务
func (this *Merchant) LevelManager() merchant.ILevelManager {
	if this._levelManager == nil {
		this._levelManager = NewLevelManager(this.GetAggregateRootId(), this._memberRep)
	}
	return this._levelManager
}

// 获取键值管理器
func (this *Merchant) KvManager() merchant.IKvManager {
	if this._kvManager == nil {
		this._kvManager = newKvManager(this, "kvset")
	}
	return this._kvManager
}

// 获取用户键值管理器
func (this *Merchant) MemberKvManager() merchant.IKvManager {
	if this._memberKvManager == nil {
		this._memberKvManager = newKvManager(this, "kvset_member")
	}
	return this._memberKvManager
}

// 消息系统管理器
func (this *Merchant) MssManager() mss.IMssManager {
	if this._mssManager == nil {
		this._mssManager = mssImpl.NewMssManager(this, this._mssRep, this._rep)
	}
	return this._mssManager
}
