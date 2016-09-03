/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-19 22:49
 * description :
 * history :
 */

package dps

import (
	"errors"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/domain/interface/sale"
	"go2o/core/dto"
	"go2o/core/query"
	"log"
	"strings"
	"time"
)

type merchantService struct {
	_mchRep     merchant.IMerchantRep
	_saleRep    sale.ISaleRep
	_query      *query.MerchantQuery
	_orderQuery *query.OrderQuery
}

func NewMerchantService(r merchant.IMerchantRep, saleRep sale.ISaleRep,
	q *query.MerchantQuery, orderQuery *query.OrderQuery) *merchantService {
	return &merchantService{
		_mchRep:     r,
		_query:      q,
		_saleRep:    saleRep,
		_orderQuery: orderQuery,
	}
}

// 创建会员申请商户密钥
func (m *merchantService) CreateSignUpToken(memberId int) string {
	return m._mchRep.CreateSignUpToken(memberId)
}

// 根据商户申请密钥获取会员编号
func (m *merchantService) GetMemberFromSignUpToken(token string) int {
	return m._mchRep.GetMemberFromSignUpToken(token)
}

// 获取会员商户申请信息
func (m *merchantService) GetMchSignUpInfoByMemberId(memberId int) *merchant.MchSignUp {
	return m._mchRep.GetManager().GetSignUpInfoByMemberId(memberId)
}

// 获取商户申请信息
func (m *merchantService) GetSignUp(id int) *merchant.MchSignUp {
	return m._mchRep.GetManager().GetSignUpInfo(id)
}

// 审核商户申请信息
func (m *merchantService) ReviewSignUp(id int, pass bool, remark string) error {
	return m._mchRep.GetManager().ReviewMchSignUp(id, pass, remark)
}

// 商户注册
func (m *merchantService) SignUp(usr, pwd, companyName string,
	province int, city int, district int) (int, error) {
	unix := time.Now().Unix()
	v := &merchant.Merchant{
		MemberId: 0,
		// 用户
		Usr: usr,
		// 密码
		Pwd: pwd,
		// 商户名称
		Name: companyName,
		// 是否自营
		SelfSales: 0,
		// 商户等级
		Level: 1,
		// 标志
		Logo: "",
		// 省
		Province: province,
		// 市
		City: city,
		// 区
		District: district,
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
	mch, err := m._mchRep.CreateMerchant(v)
	if err == nil {
		err = mch.SetValue(v)
		if err == nil {
			return mch.Save()
		}
	}
	return -1, err
}

// 提交注册信息
func (m *merchantService) SignUpPost(e *merchant.MchSignUp) (int, error) {
	return m._mchRep.GetManager().CommitSignUpInfo(e)
}

func (m *merchantService) GetMerchantByMemberId(memberId int) *merchant.Merchant {
	mch := m._mchRep.GetManager().GetMerchantByMemberId(memberId)
	if mch != nil {
		v := mch.GetValue()
		return &v
	}
	return nil
}

// 删除会员的商户申请资料
func (m *merchantService) RemoveMerchantSignUp(memberId int) error {
	return m._mchRep.GetManager().RemoveSignUp(memberId)
}

// 验证用户密码并返回编号
func (m *merchantService) Verify(usr, pwd string) int {
	usr = strings.ToLower(strings.TrimSpace(usr))
	return m._query.Verify(usr, pwd)
}

// 获取企业信息
func (m *merchantService) GetReviewedEnterpriseInfo(mchId int) *merchant.EnterpriseInfo {
	mch, err := m._mchRep.GetMerchant(mchId)
	if err != nil {
		handleError(err)
		return &merchant.EnterpriseInfo{}
	}
	return mch.ProfileManager().GetReviewedEnterpriseInfo()
}

// 获取企业信息,并返回是否为提交的信息
func (m *merchantService) GetReviewingEnterpriseInfo(mchId int) (
	e *merchant.EnterpriseInfo, isPost bool) {
	mch, err := m._mchRep.GetMerchant(mchId)
	if err != nil {
		handleError(err)
		return nil, false
	}
	mg := mch.ProfileManager()
	e = mg.GetReviewingEnterpriseInfo()
	if e != nil {
		return e, true
	}
	e = mg.GetReviewedEnterpriseInfo()
	if e != nil {
		v := *e
		v.IsHandled = 0
		v.Reviewed = 0
		return &v, false
	}
	return nil, false

}

// 保存企业信息
func (m *merchantService) SaveEnterpriseInfo(mchId int,
	e *merchant.EnterpriseInfo) (int, error) {
	mch, err := m._mchRep.GetMerchant(mchId)
	if err != nil {
		return -1, err
	}
	return mch.ProfileManager().SaveEnterpriseInfo(e)
}

// 审核企业信息
func (m *merchantService) ReviewEnterpriseInfo(mchId int, pass bool,
	remark string) error {
	mch, err := m._mchRep.GetMerchant(mchId)
	if err != nil {
		return err
	}
	return mch.ProfileManager().ReviewEnterpriseInfo(pass, remark)
}

func (m *merchantService) GetMerchant(merchantId int) (*merchant.Merchant, error) {
	mch, err := m._mchRep.GetMerchant(merchantId)
	if mch != nil {
		v := mch.GetValue()
		return &v, err
	}
	return nil, err
}

func (m *merchantService) GetAccount(mchId int) *merchant.Account {
	return m._mchRep.GetAccount(mchId)
}

func (m *merchantService) SaveMerchant(merchantId int, v *merchant.Merchant) (int, error) {
	var mch merchant.IMerchant
	var err error
	var isCreate bool

	v.Id = merchantId

	if merchantId > 0 {
		mch, _ = m._mchRep.GetMerchant(merchantId)
		if mch == nil {
			err = errors.New("no such partner")
		} else {
			err = mch.SetValue(v)
		}
	} else {
		isCreate = true
		mch, err = m._mchRep.CreateMerchant(v)
	}

	if err != nil {
		return 0, err
	}

	merchantId, err = mch.Save()

	if isCreate {
		m.initializeMerchant(merchantId)
	}

	return merchantId, err
}

func (m *merchantService) initializeMerchant(merchantId int) {

	// 初始化会员默认等级
	// m._mchRep.GetMerchant(merchantId)

	//conf := merchant.DefaultSaleConf
	//conf.MerchantId = mch.GetAggregateRootId()
	// 保存销售设置
	//mch.ConfManager().SaveSaleConf(&conf)

	// 初始化销售标签
	m._saleRep.GetSale(merchantId).LabelManager().InitSaleLabels()
}

// 获取商户的状态
func (m *merchantService) Stat(merchantId int) error {
	mch, err := m._mchRep.GetMerchant(merchantId)
	if err != nil {
		return err
	}
	return mch.Stat()
}

// 根据主机查询商户编号
func (m *merchantService) GetMerchantIdByHost(host string) int {
	return m._query.QueryMerchantIdByHost(host)
}

// 获取商户的域名
func (m *merchantService) GetMerchantMajorHost(merchantId int) string {
	mch, err := m._mchRep.GetMerchant(merchantId)
	if err != nil {
		log.Println("[ Merchant][ Service]-", err.Error())
	}
	return mch.GetMajorHost()
}

func (m *merchantService) SaveSaleConf(merchantId int, v *merchant.SaleConf) error {
	mch, _ := m._mchRep.GetMerchant(merchantId)
	return mch.ConfManager().SaveSaleConf(v)
}

func (m *merchantService) GetSaleConf(merchantId int) *merchant.SaleConf {
	mch, err := m._mchRep.GetMerchant(merchantId)
	if err != nil {
		log.Println("[ Merchant][ Service]-", err.Error(), ",ID:", merchantId)
	}
	conf := mch.ConfManager().GetSaleConf()
	return &conf
}

func (m *merchantService) GetShopsOfMerchant(merchantId int) []*shop.Shop {
	mch, err := m._mchRep.GetMerchant(merchantId)
	if err != nil {
		log.Println("[ Merchant][ Service]-", err.Error())
	}
	shops := mch.ShopManager().GetShops()
	sv := make([]*shop.Shop, len(shops))
	for i, v := range shops {
		vv := v.GetValue()
		sv[i] = &vv
	}
	return sv
}

// 获取商城
func (m *merchantService) GetOnlineShops(mchId int) []*shop.Shop {
	mch, _ := m._mchRep.GetMerchant(mchId)
	shops := mch.ShopManager().GetShops()
	sv := []*shop.Shop{}
	for _, v := range shops {
		if v.Type() == shop.TypeOnlineShop {
			vv := v.GetValue()
			sv = append(sv, &vv)
		}
	}
	return sv
}

// 修改密码
func (m *merchantService) ModifyPassword(mchId int, oldPwd, newPwd string) error {
	mch, _ := m._mchRep.GetMerchant(mchId)
	if mch != nil {
		return mch.ProfileManager().ModifyPassword(newPwd, oldPwd)
	}
	return merchant.ErrNoSuchMerchant
}

func (m *merchantService) GetMerchantsId() []int {
	return m._mchRep.GetMerchantsId()
}

// 保存API信息
func (m *merchantService) SaveApiInfo(merchantId int, d *merchant.ApiInfo) error {
	mch, _ := m._mchRep.GetMerchant(merchantId)
	return mch.ApiManager().SaveApiInfo(d)
}

// 获取API接口
func (m *merchantService) GetApiInfo(merchantId int) *merchant.ApiInfo {
	mch, _ := m._mchRep.GetMerchant(merchantId)
	v := mch.ApiManager().GetApiInfo()
	return &v
}

// 启用/停用接口权限
func (m *merchantService) ApiPerm(merchantId int, enabled bool) error {
	mch, _ := m._mchRep.GetMerchant(merchantId)
	if enabled {
		return mch.ApiManager().EnableApiPerm()
	}
	return mch.ApiManager().DisableApiPerm()
}

// 根据API ID获取MerchantId
func (m *merchantService) GetMerchantIdByApiId(apiId string) int {
	return m._mchRep.GetMerchantIdByApiId(apiId)
}

// 获取所有会员等级
func (m *merchantService) GetMemberLevels(merchantId int) []*merchant.MemberLevel {
	mch, _ := m._mchRep.GetMerchant(merchantId)
	return mch.LevelManager().GetLevelSet()
}

// 根据编号获取会员等级信息
func (m *merchantService) GetMemberLevelById(merchantId, id int) *merchant.MemberLevel {
	mch, _ := m._mchRep.GetMerchant(merchantId)
	return mch.LevelManager().GetLevelById(id)
}

// 保存会员等级信息
func (m *merchantService) SaveMemberLevel(merchantId int, v *merchant.MemberLevel) (int, error) {
	mch, _ := m._mchRep.GetMerchant(merchantId)
	return mch.LevelManager().SaveLevel(v)
}

// 删除会员等级
func (m *merchantService) DelMemberLevel(merchantId, levelId int) error {
	mch, _ := m._mchRep.GetMerchant(merchantId)
	return mch.LevelManager().DeleteLevel(levelId)
}

// 获取等级
func (m *merchantService) GetLevel(merchantId, level int) *merchant.MemberLevel {
	mch, _ := m._mchRep.GetMerchant(merchantId)
	return mch.LevelManager().GetLevelByValue(level)
}

// 获取下一个等级
func (m *merchantService) GetNextLevel(merchantId, levelValue int) *merchant.MemberLevel {
	mch, _ := m._mchRep.GetMerchant(merchantId)
	return mch.LevelManager().GetNextLevel(levelValue)
}

// 获取键值字典
func (m *merchantService) GetKeyMapsByKeyword(merchantId int, keyword string) map[string]string {
	mch, _ := m._mchRep.GetMerchant(merchantId)
	return mch.KvManager().GetsByChar(keyword)
}

// 保存键值字典
func (m *merchantService) SaveKeyMaps(merchantId int, data map[string]string) error {
	mch, err := m._mchRep.GetMerchant(merchantId)
	if err != nil {
		return err
	}
	return mch.KvManager().Sets(data)
}

// 查询分页订单
func (m *merchantService) PagedOrdersOfVendor(vendorId, begin, size int, pagination bool,
	where, orderBy string) (int, []*dto.PagedVendorOrder) {
	return m._orderQuery.PagedOrdersOfVendor(vendorId, begin, size, pagination, where, orderBy)
}

// 提到会员账户
func (m *merchantService) TakeToMemberAccount(mchId int, amount float32) error {
	mch, err := m._mchRep.GetMerchant(mchId)
	if err == nil {
		acc := mch.Account()
		err = acc.TransferToMember(amount)
	}
	return err
}
