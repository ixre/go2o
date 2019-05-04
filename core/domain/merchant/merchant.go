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
	"github.com/ixre/gof/db/orm"
	"go2o/core/domain/interface/enum"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/domain/interface/merchant/user"
	"go2o/core/domain/interface/merchant/wholesaler"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/domain/interface/wallet"
	si "go2o/core/domain/merchant/shop"
	userImpl "go2o/core/domain/merchant/user"
	wsImpl "go2o/core/domain/merchant/wholesale"
	"go2o/core/domain/tmp"
	"go2o/core/infrastructure"
	"go2o/core/infrastructure/domain"
	"go2o/core/infrastructure/domain/util"
	"go2o/core/variable"
	"strings"
	"time"
)

var _ merchant.IMerchantManager = new(merchantManagerImpl)

type merchantManagerImpl struct {
	rep     merchant.IMerchantRepo
	valRepo valueobject.IValueRepo
}

func NewMerchantManager(rep merchant.IMerchantRepo,
	valRepo valueobject.IValueRepo) merchant.IMerchantManager {
	return &merchantManagerImpl{
		rep:     rep,
		valRepo: valRepo,
	}
}

// 创建会员申请商户密钥
func (m *merchantManagerImpl) CreateSignUpToken(memberId int64) string {
	return m.rep.CreateSignUpToken(memberId)
}

// 根据商户申请密钥获取会员编号
func (m *merchantManagerImpl) GetMemberFromSignUpToken(token string) int64 {
	return m.rep.GetMemberFromSignUpToken(token)
}

// 删除会员的商户申请资料
func (m *merchantManagerImpl) RemoveSignUp(memberId int64) error {
	_, err := tmp.Db().GetOrm().Delete(merchant.MchSignUp{}, "member_id= $1", memberId)
	return err
}
func (m *merchantManagerImpl) saveSignUpInfo(v *merchant.MchSignUp) (int32, error) {
	v.UpdateTime = time.Now().Unix()
	return orm.I32(orm.Save(tmp.Db().GetOrm(), v, int(v.Id)))
}

// 检查商户注册信息是否正确
func (m *merchantManagerImpl) checkSignUpInfo(v *merchant.MchSignUp) error {
	if v.MemberId <= 0 {
		return errors.New("会员不存在")
	}
	//todo: validate and check merchant name exists
	if v.MchName == "" {
		return merchant.ErrMissingMerchantName
	}
	if v.CompanyName == "" {
		return merchant.ErrMissingCompanyName
	}
	if v.CompanyNo == "" {
		return merchant.ErrMissingCompanyNo
	}
	if v.Address == "" {
		return merchant.ErrMissingAddress
	}
	if v.PersonName == "" {
		return merchant.ErrMissingPersonName
	}
	if v.PersonId == "" {
		return merchant.ErrMissingPersonId
	}
	if util.CheckChineseCardID(v.PersonId) != nil {
		return merchant.ErrPersonCardId
	}
	if v.Phone == "" {
		return merchant.ErrMissingPhone
	}
	if v.CompanyImage == "" {
		return merchant.ErrMissingCompanyImage
	}
	if v.PersonImage == "" {
		return merchant.ErrMissingPersonImage
	}
	return nil
}

// 提交商户注册信息
func (m *merchantManagerImpl) CommitSignUpInfo(v *merchant.MchSignUp) (int32, error) {
	err := m.checkSignUpInfo(v)
	if err != nil {
		return 0, err
	}
	v.Reviewed = enum.ReviewAwaiting
	v.SubmitTime = time.Now().Unix()
	return m.saveSignUpInfo(v)
}

// 审核商户注册信息
func (m *merchantManagerImpl) ReviewMchSignUp(id int32, pass bool, remark string) error {
	var err error
	v := m.GetSignUpInfo(id)
	if v == nil {
		return merchant.ErrNoSuchSignUpInfo
	}
	if pass {
		v.Reviewed = enum.ReviewPass
		v.Remark = ""
		if err = m.createNewMerchant(v); err != nil {
			return err
		}
	} else {
		v.Reviewed = enum.ReviewReject
		v.Remark = remark
		if strings.TrimSpace(v.Remark) == "" {
			return merchant.ErrRequireRejectRemark
		}
	}
	_, err = m.saveSignUpInfo(v)
	return err
}

// 创建新商户
func (m *merchantManagerImpl) createNewMerchant(v *merchant.MchSignUp) error {
	unix := time.Now().Unix()
	mchVal := &merchant.Merchant{
		MemberId: v.MemberId,
		// 商户名称
		Name: v.MchName,
		// 是否自营
		SelfSales: 0,
		// 商户等级
		Level: 1,
		// 标志
		Logo: "",
		// 公司名称
		CompanyName: "",
		// 省
		Province: v.Province,
		// 市
		City: v.City,
		// 区
		District: v.District,
		// 是否启用
		Enabled: 1,
		// 过期时间
		ExpiresTime: time.Now().Add(time.Hour * time.Duration(24*365)).Unix(),
		// 注册时间
		JoinTime: unix,
		// 更新时间
		UpdateTime: unix,
		// 登录时间
		LoginTime: 0,
		// 最后登录时间
		LastLoginTime: 0,
	}
	mch := m.rep.CreateMerchant(mchVal)

	err := mch.SetValue(mchVal)
	if err != nil {
		return err
	}
	mchId, err := mch.Save()
	if err == nil {
		names := m.valRepo.GetAreaNames([]int32{v.Province, v.City, v.District})
		location := strings.Join(names, "")
		ev := &merchant.EnterpriseInfo{
			MchId:        mchId,
			CompanyName:  v.CompanyName,
			CompanyNo:    v.CompanyNo,
			PersonName:   v.PersonName,
			PersonIdNo:   v.PersonId,
			PersonImage:  v.PersonImage,
			Tel:          v.Phone,
			Province:     v.Province,
			City:         v.City,
			District:     v.District,
			Location:     location,
			Address:      v.Address,
			CompanyImage: v.CompanyImage,
			AuthDoc:      v.AuthDoc,
			Reviewed:     v.Reviewed,
			ReviewTime:   unix,
			ReviewRemark: "",
			UpdateTime:   unix,
		}
		_, err = mch.ProfileManager().SaveEnterpriseInfo(ev)
		if err == nil {
			mch.ProfileManager().ReviewEnterpriseInfo(true, "")
		}
	}
	return err
}

// 获取商户申请信息
func (m *merchantManagerImpl) GetSignUpInfo(id int32) *merchant.MchSignUp {
	v := merchant.MchSignUp{}
	if tmp.Db().GetOrm().Get(id, &v) != nil {
		return nil
	}
	return &v
}

// 获取会员申请的商户信息
func (m *merchantManagerImpl) GetSignUpInfoByMemberId(memberId int64) *merchant.MchSignUp {
	v := merchant.MchSignUp{}
	if tmp.Db().GetOrm().GetBy(&v, "member_id= $1", memberId) != nil {
		return nil
	}
	return &v
}

// 获取会员关联的商户
func (m *merchantManagerImpl) GetMerchantByMemberId(memberId int64) merchant.IMerchant {
	v := merchant.Merchant{}
	if tmp.Db().GetOrm().GetBy(&v, "member_id= $1", memberId) == nil {
		return m.rep.CreateMerchant(&v)
	}
	return nil
}

var _ merchant.IMerchant = new(merchantImpl)

type merchantImpl struct {
	_value           *merchant.Merchant
	_account         merchant.IAccount
	_wholesaler      wholesaler.IWholesaler
	_host            string
	_rep             merchant.IMerchantRepo
	_wsRepo          wholesaler.IWholesaleRepo
	_itemRepo        item.IGoodsItemRepo
	_shopRepo        shop.IShopRepo
	_userRepo        user.IUserRepo
	_valRepo         valueobject.IValueRepo
	_memberRepo      member.IMemberRepo
	_userManager     user.IUserManager
	_confManager     merchant.IConfManager
	_saleManager     merchant.ISaleManager
	_levelManager    merchant.ILevelManager
	_kvManager       merchant.IKvManager
	_memberKvManager merchant.IKvManager
	//_mssManager      mss.IMssManager
	//_mssRepo          mss.IMssRepo
	_profileManager merchant.IProfileManager
	_apiManager     merchant.IApiManager
	_shopManager    shop.IShopManager
	_walletRepo     wallet.IWalletRepo
}

func NewMerchant(v *merchant.Merchant, rep merchant.IMerchantRepo,
	wsRepo wholesaler.IWholesaleRepo, itemRepo item.IGoodsItemRepo,
	shopRepo shop.IShopRepo, userRepo user.IUserRepo, memberRepo member.IMemberRepo,
	walletRepo wallet.IWalletRepo, valRepo valueobject.IValueRepo) merchant.IMerchant {
	mch := &merchantImpl{
		_value:      v,
		_rep:        rep,
		_wsRepo:     wsRepo,
		_itemRepo:   itemRepo,
		_shopRepo:   shopRepo,
		_userRepo:   userRepo,
		_valRepo:    valRepo,
		_memberRepo: memberRepo,
		_walletRepo: walletRepo,
	}
	return mch
}

func (m *merchantImpl) GetRepo() merchant.IMerchantRepo {
	return m._rep
}

func (m *merchantImpl) GetAggregateRootId() int32 {
	return m._value.ID
}

// 获取符合的商家信息
func (m *merchantImpl) Complex() *merchant.ComplexMerchant {
	src := m.GetValue()
	return &merchant.ComplexMerchant{
		Id:            src.ID,
		MemberId:      src.MemberId,
		Usr:           src.Usr,
		Pwd:           src.Pwd,
		Name:          src.Name,
		SelfSales:     src.SelfSales,
		Level:         src.Level,
		Logo:          src.Logo,
		CompanyName:   src.CompanyName,
		Province:      src.Province,
		City:          src.City,
		District:      src.District,
		Enabled:       src.Enabled,
		ExpiresTime:   src.ExpiresTime,
		JoinTime:      src.JoinTime,
		UpdateTime:    src.UpdateTime,
		LoginTime:     src.LoginTime,
		LastLoginTime: src.LastLoginTime,
	}
}

func (m *merchantImpl) GetValue() merchant.Merchant {
	return *m._value
}

func (m *merchantImpl) SetValue(v *merchant.Merchant) error {
	tv := m._value
	if v.ID == tv.ID {
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
		if len(v.CompanyName) != 0 {
			tv.CompanyName = v.CompanyName
		}
		tv.Pwd = v.Pwd
		tv.UpdateTime = time.Now().Unix()
	}
	return nil
}

// 保存
func (m *merchantImpl) Save() (int32, error) {
	id := m.GetAggregateRootId()
	if id > 0 {
		m.checkSelfSales()
		return m._rep.SaveMerchant(m._value)
	}
	return m.createMerchant()
}

// 自营检测,并返回是否需要保存
func (m *merchantImpl) checkSelfSales() bool {
	if m._value.SelfSales == 0 {
		//不为自营,但编号为1的商户
		if m.GetAggregateRootId() == 1 {
			m._value.SelfSales = 1
			m._value.Usr = "-"
			m._value.Pwd = "-"
			return true
		}
	} else if m.GetAggregateRootId() != 1 {
		//为自营,但ID不为0,异常数据
		m._value.SelfSales = 0
		m._value.Enabled = 0
		return true
	}
	return false
}

// 是否自营
func (m *merchantImpl) SelfSales() bool {
	return m._value.SelfSales == 1 ||
		m.GetAggregateRootId() == 1
}

// 获取商户的状态,判断 过期时间、判断是否停用。
// 过期时间通常按: 试合作期,即1个月, 后面每年延长一次。或者会员付费开通。
func (m *merchantImpl) Stat() error {
	if m._value == nil {
		return merchant.ErrNoSuchMerchant
	}
	if m._value.Enabled == 0 {
		return merchant.ErrMerchantDisabled
	}
	if m._value.ExpiresTime < time.Now().Unix() {
		return merchant.ErrMerchantExpires
	}
	return nil
}

// 设置商户启用或停用
func (m *merchantImpl) SetEnabled(enabled bool) error {
	if enabled {
		m._value.Enabled = 1
	} else {
		m._value.Enabled = 0
	}
	_, err := m.Save()
	return err
}

// 返回对应的会员编号
func (m *merchantImpl) Member() int64 {
	return m.GetValue().MemberId
}

// 获取商户账户
func (m *merchantImpl) Account() merchant.IAccount {
	if m._account == nil {
		v := m._rep.GetAccount(m.GetAggregateRootId())
		m._account = newAccountImpl(m, v, m._memberRepo, m._walletRepo)
	}
	return m._account
}

// 获取批发商实例
func (m *merchantImpl) Wholesaler() wholesaler.IWholesaler {
	if m._wholesaler == nil {
		mchId := m.GetAggregateRootId()
		v := m._wsRepo.GetWsWholesaler(mchId)
		if v == nil {
			v, _ = m.createWholesaler()
		}
		m._wholesaler = wsImpl.NewWholesaler(mchId, v, m._wsRepo, m._itemRepo)
	}
	return m._wholesaler
}

// 启用批发
func (m *merchantImpl) EnableWholesale() error {
	if m.Wholesaler() != nil {
		return errors.New("wholesale for merchant enabled!")
	}
	_, err := m.createWholesaler()
	return err
}

func (m *merchantImpl) createWholesaler() (*wholesaler.WsWholesaler, error) {
	v := &wholesaler.WsWholesaler{
		MchId:       m.GetAggregateRootId(),
		Rate:        1,
		ReviewState: enum.ReviewPass,
		//ReviewState: enum.ReviewAwaiting,
	}
	_, err := m._wsRepo.SaveWsWholesaler(v, true)
	return v, err
}

// 创建商户
func (m *merchantImpl) createMerchant() (int32, error) {
	if id := m.GetAggregateRootId(); id > 0 {
		return id, nil
	}

	id, err := m._rep.SaveMerchant(m._value)
	if err != nil {
		return id, err
	}

	//todo:事务

	// 初始化商户信息
	m._value.ID = id

	// 检测自营并保存
	if m.checkSelfSales() {
		m._rep.SaveMerchant(m._value)
	}

	//todo:  初始化商店

	// SiteConf
	//m._siteConf = &shop.ShopSiteConf{
	//	IndexTitle: "线上商店-" + v.Name,
	//	SubTitle:   "线上商店-" + v.Name,
	//	Logo:       v.Logo,
	//	State:      1,
	//	StateHtml:  "",
	//}
	//err = m._rep.SaveSiteConf(id, m._siteConf)
	//m._siteConf.VendorId = id

	// SaleConf
	//m._saleConf = &merchant.SaleConf{
	//	AutoSetupOrder:  1,
	//	IntegralBackNum: 0,
	//}
	//err = m._rep.SaveSaleConf(id, m._saleConf)
	//m._saleConf.VendorId = id

	// 创建API
	api := &merchant.ApiInfo{
		ApiId:     domain.NewApiId(int(id)),
		ApiSecret: domain.NewSecret(int(id)),
		WhiteList: "*",
		Enabled:   1,
	}
	err = m.ApiManager().SaveApiInfo(api)
	return id, err
}

// 获取商户的域名
func (m *merchantImpl) GetMajorHost() string {
	if len(m._host) == 0 {
		host := m._rep.GetMerchantMajorHost(m.GetAggregateRootId())
		if len(host) == 0 {
			host = fmt.Sprintf("%s.%s", m._value.Usr, infrastructure.GetApp().
				Config().GetString(variable.ServerDomain))
		}
		m._host = host
	}
	return m._host
}

// 返回用户服务
func (m *merchantImpl) UserManager() user.IUserManager {
	if m._userManager == nil {
		m._userManager = userImpl.NewUserManager(
			m.GetAggregateRootId(),
			m._userRepo)
	}
	return m._userManager
}

// 获取会员管理服务
func (m *merchantImpl) LevelManager() merchant.ILevelManager {
	if m._levelManager == nil {
		m._levelManager = NewLevelManager(m.GetAggregateRootId(), m._rep)
	}
	return m._levelManager
}

// 获取键值管理器
func (m *merchantImpl) KvManager() merchant.IKvManager {
	if m._kvManager == nil {
		m._kvManager = newKvManager(m, "kvset")
	}
	return m._kvManager
}

// 获取用户键值管理器
func (m *merchantImpl) MemberKvManager() merchant.IKvManager {
	if m._memberKvManager == nil {
		m._memberKvManager = newKvManager(m, "kvset_member")
	}
	return m._memberKvManager
}

// 消息系统管理器
//func (m *MerchantImpl) MssManager() mss.IMssManager {
//	if m._mssManager == nil {
//		m._mssManager = mssImpl.NewMssManager(m, m._mssRepo, m._rep)
//	}
//	return m._mssManager
//}

// 返回设置服务
func (m *merchantImpl) ConfManager() merchant.IConfManager {
	if m._confManager == nil {
		m._confManager = newConfigManagerImpl(m.GetAggregateRootId(),
			m._rep, m._memberRepo, m._valRepo)
	}
	return m._confManager
}

// 销售服务
func (m *merchantImpl) SaleManager() merchant.ISaleManager {
	if m._saleManager == nil {
		m._saleManager = newSaleManagerImpl(int(m.GetAggregateRootId()), m)
	}
	return m._saleManager
}

// 企业资料管理器
func (m *merchantImpl) ProfileManager() merchant.IProfileManager {
	if m._profileManager == nil {
		m._profileManager = newProfileManager(m, m._valRepo)
	}
	return m._profileManager
}

// API服务
func (m *merchantImpl) ApiManager() merchant.IApiManager {
	if m._apiManager == nil {
		m._apiManager = newApiManagerImpl(m)
	}
	return m._apiManager
}

// 商店服务
func (m *merchantImpl) ShopManager() shop.IShopManager {
	if m._shopManager == nil {
		m._shopManager = si.NewShopManagerImpl(m, m._shopRepo, m._valRepo)
	}
	return m._shopManager
}
