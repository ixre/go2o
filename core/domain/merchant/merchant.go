/**
 * Copyright 2014 @ 56x.net.
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
	"time"

	"github.com/ixre/go2o/core/domain/interface/approval"
	"github.com/ixre/go2o/core/domain/interface/domain/enum"
	"github.com/ixre/go2o/core/domain/interface/invoice"
	"github.com/ixre/go2o/core/domain/interface/item"
	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/domain/interface/merchant"
	"github.com/ixre/go2o/core/domain/interface/merchant/shop"
	"github.com/ixre/go2o/core/domain/interface/merchant/staff"
	"github.com/ixre/go2o/core/domain/interface/merchant/user"
	"github.com/ixre/go2o/core/domain/interface/merchant/wholesaler"
	rbac "github.com/ixre/go2o/core/domain/interface/rabc"
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/domain/interface/sys"
	"github.com/ixre/go2o/core/domain/interface/valueobject"
	"github.com/ixre/go2o/core/domain/interface/wallet"
	si "github.com/ixre/go2o/core/domain/merchant/shop"
	staffImpl "github.com/ixre/go2o/core/domain/merchant/staff"
	userImpl "github.com/ixre/go2o/core/domain/merchant/user"
	wsImpl "github.com/ixre/go2o/core/domain/merchant/wholesale"
	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/go2o/core/initial/provide"
	"github.com/ixre/go2o/core/variable"
	"github.com/ixre/gof/storage"
)

var _ merchant.IMerchantAggregateRoot = new(merchantImpl)

type merchantImpl struct {
	_storage         storage.Interface
	_value           *merchant.Merchant
	_account         merchant.IAccount
	_wholesaler      wholesaler.IWholesaler
	_host            string
	_repo            merchant.IMerchantRepo
	_wsRepo          wholesaler.IWholesaleRepo
	_itemRepo        item.IItemRepo
	_shopRepo        shop.IShopRepo
	_userRepo        user.IUserRepo
	_staffRepo       staff.IStaffRepo
	_valRepo         valueobject.IValueRepo
	_memberRepo      member.IMemberRepo
	_sysRepo         sys.ISystemRepo
	_userManager     user.IUserManager
	_confManager     merchant.IConfManager
	_saleManager     merchant.IMerchantTransactionManager
	_levelManager    merchant.ILevelManager
	_kvManager       merchant.IKvManager
	_memberKvManager merchant.IKvManager
	//_mssManager      mss.IMssManager
	//_mssRepo          mss.IMssRepo
	_profileManager  merchant.IProfileManager
	_apiManager      merchant.IApiManager
	_shopManager     shop.IShopManager
	_employeeManager staff.IStaffManager
	_walletRepo      wallet.IWalletRepo
	_registryRepo    registry.IRegistryRepo
	_approvalRepo    approval.IApprovalRepository
	_invoiceRepo     invoice.IInvoiceRepo
	// 之前绑定的会员编号
	_lastBindMemberId int
	_rbacRepo         rbac.IRbacRepository
}

// EmployeeManager implements merchant.IMerchant.
func (m *merchantImpl) EmployeeManager() staff.IStaffManager {
	if m._employeeManager == nil {
		m._employeeManager = staffImpl.NewStaffManager(m,
			m._staffRepo,
			m._memberRepo,
			m._sysRepo,
			m._repo,
			m._approvalRepo,
			m._storage)
	}
	return m._employeeManager
}

func NewMerchant(v *merchant.Merchant,
	storage storage.Interface,
	rep merchant.IMerchantRepo,
	wsRepo wholesaler.IWholesaleRepo, itemRepo item.IItemRepo,
	shopRepo shop.IShopRepo, userRepo user.IUserRepo,
	employeeRepo staff.IStaffRepo,
	memberRepo member.IMemberRepo,
	sysRepo sys.ISystemRepo,
	walletRepo wallet.IWalletRepo, valRepo valueobject.IValueRepo,
	registryRepo registry.IRegistryRepo,
	invoiceRepo invoice.IInvoiceRepo,
	approvalRepo approval.IApprovalRepository,
	rbacRepo rbac.IRbacRepository,
) merchant.IMerchantAggregateRoot {
	mch := &merchantImpl{
		_storage:      storage,
		_value:        v,
		_repo:         rep,
		_wsRepo:       wsRepo,
		_itemRepo:     itemRepo,
		_shopRepo:     shopRepo,
		_userRepo:     userRepo,
		_staffRepo:    employeeRepo,
		_valRepo:      valRepo,
		_memberRepo:   memberRepo,
		_sysRepo:      sysRepo,
		_walletRepo:   walletRepo,
		_registryRepo: registryRepo,
		_approvalRepo: approvalRepo,
		_invoiceRepo:  invoiceRepo,
		_rbacRepo:     rbacRepo,
	}
	return mch
}

func (m *merchantImpl) GetRepo() merchant.IMerchantRepo {
	return m._repo
}

func (m *merchantImpl) GetAggregateRootId() int {
	return m._value.Id
}

// Complex 获取符合的商家信息
func (m *merchantImpl) Complex() *merchant.ComplexMerchant {
	src := m.GetValue()
	return &merchant.ComplexMerchant{
		Id:        int32(src.Id),
		MemberId:  int64(src.MemberId),
		Username:  src.Username,
		Pwd:       src.Password,
		Name:      src.MchName,
		SelfSales: int32(src.IsSelf),
		Level:     int32(src.Level),
		Flag:      src.Flag,
		Address:   src.Address,
		//Logo:          src.Logo,
		//CompanyName:   src.CompanyName,
		Province:  int32(src.Province),
		City:      int32(src.City),
		District:  int32(src.District),
		Telephone: src.Tel,
		Status:    int(src.Status),
		// Enabled:       int32(src.Enabled),
		// ExpiresTime:   src.ExpiresTime,
		// JoinTime:      src.CreateTime,
		// UpdateTime:    src.UpdateTime,
		// LoginTime:     src.LoginTime,
		// LastLoginTime: src.LastLoginTime,
	}
}

func (m *merchantImpl) GetValue() merchant.Merchant {
	return *m._value
}

func (m *merchantImpl) SetValue(v *merchant.Merchant) error {
	if err := m.check(v); err != nil {
		return err
	}
	tv := m._value
	if v.Id == tv.Id {
		tv.MchName = v.MchName
		tv.Province = v.Province
		tv.City = v.City
		tv.District = v.District
		tv.Address = v.Address
	}

	if m.GetAggregateRootId() <= 0 {
		m._value.MemberId = v.MemberId
	}
	if len(tv.Username) == 0 {
		tv.Username = v.Username
	}
	if v.LastLoginTime > 0 {
		tv.LastLoginTime = v.LastLoginTime
	}
	if len(v.Logo) != 0 {
		tv.Logo = v.Logo
	}
	tv.IsSelf = v.IsSelf
	tv.Status = v.Status
	tv.ExpiresTime = v.ExpiresTime
	return nil
}

// 检查商户注册信息是否正确
func (m *merchantImpl) check(v *merchant.Merchant) error {
	//todo: validate and check merchant name exists
	if len(v.MchName) > 0 {
		//todo: 检查商户名称是否存在
		//return merchant.ErrMissingMerchantName
	}
	// if v.CompanyName == "" {
	// 	return merchant.ErrMissingCompanyName
	// }
	// if v.CompanyNo == "" {
	// 	return merchant.ErrMissingCompanyNo
	// }
	// if v.Address == "" {
	// 	return merchant.ErrMissingAddress
	// }
	// if v.PersonName == "" {
	// 	return merchant.ErrMissingPersonName
	// }
	// if v.PersonId == "" {
	// 	return merchant.ErrMissingPersonId
	// }
	// if util.CheckChineseCardID(v.PersonId) != nil {
	// 	return merchant.ErrPersonCardId
	// }
	// if v.Phone == "" {
	// 	return merchant.ErrMissingPhone
	// }
	// if v.CompanyImage == "" {
	// 	return merchant.ErrMissingCompanyImage
	// }
	// if v.PersonImage == "" {
	// 	return merchant.ErrMissingPersonImage
	// }
	return nil
}

// GrantFlag 标志赋值, 如果flag小于零, 则异或运算
func (m *merchantImpl) GrantFlag(flag int) error {
	v, err := domain.GrantFlag(m._value.Flag, flag)
	if err == nil {
		m._value.Flag = v
	}
	return err
}

// Lock implements merchant.IMerchantAggregateRoot.
func (m *merchantImpl) Lock() error {
	f := m._value.Flag
	err := m.GrantFlag(merchant.FLocked)
	if err == nil && f != m._value.Flag {
		_, err = m.Save()
	}
	return err
}

// Unlock implements merchant.IMerchantAggregateRoot.
func (m *merchantImpl) Unlock() error {
	f := m._value.Flag
	err := m.GrantFlag(-merchant.FLocked)
	if err == nil && f != m._value.Flag {
		_, err = m.Save()
	}
	return err
}

// ContainFlag implements merchant.IMerchant.
func (m *merchantImpl) ContainFlag(flag int) bool {
	return domain.TestFlag(m._value.Flag, flag)
}

// 绑定会员
func (m *merchantImpl) BindMember(memberId int) error {
	if m._value.MemberId == memberId {
		return merchant.ErrMemberBindExists
	}
	exist := m._repo.CheckMemberBind(memberId, m.GetAggregateRootId())
	if exist {
		return merchant.ErrBindAnotherMerchant
	}
	m._lastBindMemberId = m._value.MemberId
	m._value.MemberId = memberId
	if m.GetAggregateRootId() > 0 {
		err := m.applyBindMember()
		if err == nil {
			_, err = m.Save()
		}
		return err
	}
	return nil
}

// 绑定会员
func (m *merchantImpl) applyBindMember() error {
	// 解绑
	if m._lastBindMemberId > 0 {
		origin := m._memberRepo.GetMember(int64(m._lastBindMemberId))
		if origin != nil {
			_ = origin.GrantFlag(-member.FlagSeller)
		}
	}
	// 添加商户标志
	im := m._memberRepo.GetMember(int64(m._value.MemberId))
	if im == nil {
		return member.ErrNoSuchMember
	}
	err := im.GrantFlag(member.FlagSeller)
	if err == nil {
		m._lastBindMemberId = m._value.MemberId
	}
	return err
}

// Save 保存
func (m *merchantImpl) Save() (int64, error) {
	id := m.GetAggregateRootId()
	if id > 0 {
		id, err := m._repo.SaveMerchant(m._value)
		return int64(id), err
	}
	return m.createMerchant()
}

// SelfSales 是否自营
func (m *merchantImpl) SelfSales() bool {
	return m._value.IsSelf == 1 || m.GetAggregateRootId() == 1
}

// Stat 获取商户的状态,判断 过期时间、判断是否停用。
// 过期时间通常按: 试合作期,即1个月, 后面每年延长一次。或者会员付费开通。
func (m *merchantImpl) Stat() error {
	if m._value == nil {
		return merchant.ErrNoSuchMerchant
	}
	if m._value.Status == 0 {
		return merchant.ErrMerchantDisabled
	}
	if m._value.ExpiresTime > 0 && m._value.ExpiresTime < int(time.Now().Unix()) {
		return merchant.ErrMerchantExpires
	}
	return nil
}

// Member 返回对应的会员编号
func (m *merchantImpl) Member() int64 {
	return int64(m.GetValue().MemberId)
}

// Account 获取商户账户
func (m *merchantImpl) Account() merchant.IAccount {
	if m._account == nil {
		v := m._repo.GetAccount(int(m.GetAggregateRootId()))
		m._account = newAccountImpl(m, v, m._memberRepo, m._walletRepo, m._invoiceRepo)
	}
	return m._account
}

// Wholesaler 获取批发商实例
func (m *merchantImpl) Wholesaler() wholesaler.IWholesaler {
	if m._wholesaler == nil {
		mchId := m.GetAggregateRootId()
		v := m._wsRepo.GetWsWholesaler(mchId)
		if v == nil {
			v, _ = m.createWholesaler()
		}
		m._wholesaler = wsImpl.NewWholesaler(int64(mchId), v, m._wsRepo, m._itemRepo)
	}
	return m._wholesaler
}

// EnableWholesale 启用批发
func (m *merchantImpl) EnableWholesale() error {
	if m.Wholesaler() != nil {
		return errors.New("wholesale for merchant enabled!")
	}
	_, err := m.createWholesaler()
	return err
}

func (m *merchantImpl) createWholesaler() (*wholesaler.WsWholesaler, error) {
	v := &wholesaler.WsWholesaler{
		MchId:        int64(m.GetAggregateRootId()),
		Rate:         1,
		ReviewStatus: enum.ReviewApproved,
		//ReviewStatus: enum.ReviewPending,
	}
	_, err := m._wsRepo.SaveWsWholesaler(v, true)
	return v, err
}

// 创建商户
func (m *merchantImpl) createMerchant() (int64, error) {
	// 验证商户用户是否占用
	if len(m._value.Username) == 0 {
		return 0, merchant.ErrMissingMerchantUser
	}
	if m._repo.CheckUserExists(m._value.Username, 0) {
		return 0, merchant.ErrMerchantUserExists
	}

	unix := time.Now().Unix()
	m._value.CreateTime = int(unix)
	// 状态: 待开通
	m._value.Status = 0
	m._value.Flag = 0

	// 设置邮箱
	m._value.MailAddr = m._value.Username

	id, err := m._repo.SaveMerchant(m._value)
	if err != nil {
		return int64(id), err
	}
	m._value.Id = id
	// 绑定会员
	if m._value.MemberId > 0 {
		err = m.applyBindMember()
	}
	// 初始化认证信息
	auth := &merchant.Authenticate{
		Id:               0,
		MchId:            0,
		OrgName:          "",
		LicenceNo:        "",
		LicencePic:       "",
		WorkCity:         0,
		QualificationPic: "",
		PersonId:         "",
		PersonName:       "",
		PersonFrontPic:   "",
		PersonBackPic:    "",
		PersonPhone:      "",
		AuthorityPic:     "",
		BankName:         "",
		BankAccount:      "",
		BankNo:           "",
		ExtraData:        "",
		ReviewTime:       0,
		ReviewStatus:     0,
		ReviewRemark:     "",
		UpdateTime:       int(unix),
	}
	auth.Id, _ = m._repo.SaveAuthenticate(auth)

	// 创建API
	api := &merchant.ApiInfo{
		ApiId:     domain.NewApiId(int(id)),
		ApiSecret: domain.NewSecret(int(id)),
		WhiteList: "*",
		Enabled:   1,
	}
	err = m.ApiManager().SaveApiInfo(api)
	return int64(id), err
}

// GetMajorHost 获取商户的域名
func (m *merchantImpl) GetMajorHost() string {
	if len(m._host) == 0 {
		host := m._repo.GetMerchantMajorHost(int(m.GetAggregateRootId()))
		if len(host) == 0 {
			cfg := provide.GetApp().Config()
			host = fmt.Sprintf("%s.%s", m._value.Username, cfg.GetString(variable.ServerDomain))
		}
		m._host = host
	}
	return m._host
}

// UserManager 返回用户服务
func (m *merchantImpl) UserManager() user.IUserManager {
	if m._userManager == nil {
		m._userManager = userImpl.NewUserManager(int64(m.GetAggregateRootId()), m._userRepo)
	}
	return m._userManager
}

// LevelManager 获取会员管理服务
func (m *merchantImpl) LevelManager() merchant.ILevelManager {
	if m._levelManager == nil {
		m._levelManager = NewLevelManager(int64(m.GetAggregateRootId()), m._repo)
	}
	return m._levelManager
}

// KvManager 获取键值管理器
func (m *merchantImpl) KvManager() merchant.IKvManager {
	if m._kvManager == nil {
		m._kvManager = newKvManager(m, "kvset")
	}
	return m._kvManager
}

// MemberKvManager 获取用户键值管理器
func (m *merchantImpl) MemberKvManager() merchant.IKvManager {
	if m._memberKvManager == nil {
		m._memberKvManager = newKvManager(m, "kvset_member")
	}
	return m._memberKvManager
}

// 消息系统管理器
//func (m *MerchantImpl) MssManager() mss.IMssManager {
//	if m._mssManager == nil {
//		m._mssManager = mssImpl.NewMssManager(m, m._mssRepo, m._repo)
//	}
//	return m._mssManager
//}

// ConfManager 返回设置服务
func (m *merchantImpl) ConfManager() merchant.IConfManager {
	if m._confManager == nil {
		m._confManager = newConfigManagerImpl(int(m.GetAggregateRootId()),
			m._repo, m._memberRepo, m._valRepo)
	}
	return m._confManager
}

// TransactionManager 销售服务
func (m *merchantImpl) TransactionManager() merchant.IMerchantTransactionManager {
	if m._saleManager == nil {
		m._saleManager = newTransactionManagerImpl(int(m.GetAggregateRootId()),
			m,
			m._repo,
			m._invoiceRepo,
			m._walletRepo,
			m._rbacRepo,
			m._registryRepo)
	}
	return m._saleManager
}

// ProfileManager 企业资料管理器
func (m *merchantImpl) ProfileManager() merchant.IProfileManager {
	if m._profileManager == nil {
		m._profileManager = newProfileManager(m, m._valRepo, m._invoiceRepo)
	}
	return m._profileManager
}

// ApiManager API服务
func (m *merchantImpl) ApiManager() merchant.IApiManager {
	if m._apiManager == nil {
		m._apiManager = newApiManagerImpl(m)
	}
	return m._apiManager
}

// ShopManager 商店服务
func (m *merchantImpl) ShopManager() shop.IShopManager {
	if m._shopManager == nil {
		m._shopManager = si.NewShopManagerImpl(m, m._shopRepo, m._valRepo, m._registryRepo)
	}
	return m._shopManager
}
