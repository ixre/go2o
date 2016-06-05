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
    "fmt"
    "go2o/core/domain/interface/merchant"
    "go2o/core/domain/interface/merchant/mss"
    "go2o/core/domain/interface/merchant/shop"
    "go2o/core/domain/interface/merchant/user"
    "go2o/core/domain/interface/valueobject"
    mssImpl "go2o/core/domain/merchant/mss"
    si "go2o/core/domain/merchant/shop"
    userImpl "go2o/core/domain/merchant/user"
    "go2o/core/infrastructure"
    "go2o/core/infrastructure/domain"
    "go2o/core/variable"
    "time"
)

var _ merchant.IMerchant = new(MerchantImpl)

type MerchantImpl struct {
    _value           *merchant.Merchant
    _host            string
    _rep             merchant.IMerchantRep
    _shopRep         shop.IShopRep
    _userRep         user.IUserRep
    _valRep          valueobject.IValueRep
    _userManager     user.IUserManager
    _confManager     merchant.IConfManager
    _levelManager    merchant.ILevelManager
    _kvManager       merchant.IKvManager
    _memberKvManager merchant.IKvManager
    _mssManager      mss.IMssManager
    _mssRep          mss.IMssRep
    _profileManager  merchant.IProfileManager
    _apiManager      merchant.IApiManager
    _shopManager     shop.IShopManager
}

func NewMerchant(v *merchant.Merchant, rep merchant.IMerchantRep,
shopRep shop.IShopRep, userRep user.IUserRep,
mssRep mss.IMssRep, valRep valueobject.IValueRep) (merchant.IMerchant, error) {
    mch := &MerchantImpl{
        _value:   v,
        _rep:     rep,
        _shopRep: shopRep,
        _userRep: userRep,
        _mssRep:  mssRep,
        _valRep:  valRep,
    }
    return mch, mch.Stat()
}

func (this *MerchantImpl) GetRep() merchant.IMerchantRep {
    return this._rep
}

func (this *MerchantImpl) GetAggregateRootId() int {
    return this._value.Id
}
func (this *MerchantImpl) GetValue() merchant.Merchant {
    return *this._value
}

func (this *MerchantImpl) SetValue(v *merchant.Merchant) error {
    tv := this._value
    if v.Id == tv.Id {
        tv.Name = v.Name
        tv.Province = v.Province
        tv.City = v.City
        tv.District = v.District
        if v.LastLoginTime > 0 {
            tv.LastLoginTime = v.LastLoginTime
        }
        if v.LoginTime > 0 {
            tv.LoginTime = v.LoginTime
        }

        if len(v.Logo) != 0 {
            tv.Logo = v.Logo
        }
        tv.Pwd = v.Pwd
        tv.UpdateTime = time.Now().Unix()
    }
    return nil
}

// 保存
func (this *MerchantImpl) Save() (int, error) {
    var id int = this.GetAggregateRootId()

    if id > 0 {
        this.checkSelfSales()
        return this._rep.SaveMerchant(this._value)
    }

    return this.createMerchant()
}

// 自营检测,并返回是否需要保存
func (this *MerchantImpl) checkSelfSales() bool {
    if this._value.SelfSales == 0 { //不为自营,但编号为1的商户
        if this.GetAggregateRootId() == 1 {
            this._value.SelfSales = 1
            this._value.Usr = "-"
            this._value.Pwd = "-"
            return true
        }
    } else if this.GetAggregateRootId() != 1 { //为自营,但ID不为0,异常数据
        this._value.SelfSales = 0
        this._value.Enabled = 0
        return true
    }
    return false
}

// 是否自营
func (this *MerchantImpl) SelfSales() bool {
    return this._value.SelfSales == 1 ||
    this.GetAggregateRootId() == 1
}

// 获取商户的状态,判断 过期时间、判断是否停用。
// 过期时间通常按: 试合作期,即1个月, 后面每年延长一次。或者会员付费开通。
func (this *MerchantImpl) Stat() error {
    if this._value == nil {
        return merchant.ErrNoSuchMerchant
    }
    if this._value.Enabled == 0 {
        //log.Println("[MERCHANT][ IMPL] - ",this._value)
        return merchant.ErrMerchantDisabled
    }
    if this._value.ExpiresTime < time.Now().Unix() {
        return merchant.ErrMerchantExpires
    }
    return nil
}

// 返回对应的会员编号
func (this *MerchantImpl) Member() int {
    return this.GetValue().MemberId
}

// 创建商户
func (this *MerchantImpl) createMerchant() (int, error) {
    if id := this.GetAggregateRootId(); id > 0 {
        return id, nil
    }

    id, err := this._rep.SaveMerchant(this._value)
    if err != nil {
        return id, err
    }

    //todo:事务

    // 初始化商户信息
    this._value.Id = id

    // 检测自营并保存
    if this.checkSelfSales() {
        this._rep.SaveMerchant(this._value)
    }

    //todo:  初始化商店

    // SiteConf
    //this._siteConf = &shop.ShopSiteConf{
    //	IndexTitle: "线上商店-" + v.Name,
    //	SubTitle:   "线上商店-" + v.Name,
    //	Logo:       v.Logo,
    //	State:      1,
    //	StateHtml:  "",
    //}
    //err = this._rep.SaveSiteConf(id, this._siteConf)
    //this._siteConf.MerchantId = id

    // SaleConf
    //this._saleConf = &merchant.SaleConf{
    //	AutoSetupOrder:  1,
    //	IntegralBackNum: 0,
    //}
    //err = this._rep.SaveSaleConf(id, this._saleConf)
    //this._saleConf.MerchantId = id

    // 创建API
    api := &merchant.ApiInfo{
        ApiId:     domain.NewApiId(id),
        ApiSecret: domain.NewSecret(id),
        WhiteList: "*",
        Enabled:   1,
    }
    err = this.ApiManager().SaveApiInfo(api)
    return id, err
}

// 获取商户的域名
func (this *MerchantImpl) GetMajorHost() string {
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

// 返回用户服务
func (this *MerchantImpl) UserManager() user.IUserManager {
    if this._userManager == nil {
        this._userManager = userImpl.NewUserManager(
            this.GetAggregateRootId(),
            this._userRep)
    }
    return this._userManager
}

// 获取会员管理服务
func (this *MerchantImpl) LevelManager() merchant.ILevelManager {
    return nil
    /*
       if this._levelManager == nil {
           this._levelManager = NewLevelManager(this.GetAggregateRootId(), this._memberRep)
       }
       return this._levelManager*/
}

// 获取键值管理器
func (this *MerchantImpl) KvManager() merchant.IKvManager {
    if this._kvManager == nil {
        this._kvManager = newKvManager(this, "kvset")
    }
    return this._kvManager
}

// 获取用户键值管理器
func (this *MerchantImpl) MemberKvManager() merchant.IKvManager {
    if this._memberKvManager == nil {
        this._memberKvManager = newKvManager(this, "kvset_member")
    }
    return this._memberKvManager
}

// 消息系统管理器
func (this *MerchantImpl) MssManager() mss.IMssManager {
    if this._mssManager == nil {
        this._mssManager = mssImpl.NewMssManager(this, this._mssRep, this._rep)
    }
    return this._mssManager
}

// 返回设置服务
func (this *MerchantImpl) ConfManager() merchant.IConfManager {
    if this._confManager == nil {
        this._confManager = newConfigManagerImpl(this, this._rep, this._valRep)
    }
    return this._confManager
}

// 企业资料管理器
func (this *MerchantImpl) ProfileManager() merchant.IProfileManager {
    if this._profileManager == nil {
        this._profileManager = newProfileManager(this)
    }
    return this._profileManager
}

// API服务
func (this *MerchantImpl) ApiManager() merchant.IApiManager {
    if this._apiManager == nil {
        this._apiManager = newApiManagerImpl(this)
    }
    return this._apiManager
}

// 商店服务
func (this *MerchantImpl) ShopManager() shop.IShopManager {
    if this._shopManager == nil {
        this._shopManager = si.NewShopManagerImpl(this, this._shopRep)
    }
    return this._shopManager
}
