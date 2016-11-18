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
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/domain/interface/sale"
	"go2o/core/dto"
	"go2o/core/query"
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
func (m *merchantService) CreateSignUpToken(memberId int64) string {
	return m._mchRep.CreateSignUpToken(memberId)
}

// 根据商户申请密钥获取会员编号
func (m *merchantService) GetMemberFromSignUpToken(token string) int64 {
	return m._mchRep.GetMemberFromSignUpToken(token)
}

// 获取会员商户申请信息
func (m *merchantService) GetMchSignUpInfoByMemberId(memberId int64) *merchant.MchSignUp {
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
	mch := m._mchRep.CreateMerchant(v)
	err := mch.SetValue(v)
	if err == nil {
		return mch.Save()
	}
	return -1, err
}

// 提交注册信息
func (m *merchantService) SignUpPost(e *merchant.MchSignUp) (int64, error) {
	return m._mchRep.GetManager().CommitSignUpInfo(e)
}

func (m *merchantService) GetMerchantByMemberId(memberId int64) *merchant.Merchant {
	mch := m._mchRep.GetManager().GetMerchantByMemberId(memberId)
	if mch != nil {
		v := mch.GetValue()
		return &v
	}
	return nil
}

// 删除会员的商户申请资料
func (m *merchantService) RemoveMerchantSignUp(memberId int64) error {
	return m._mchRep.GetManager().RemoveSignUp(memberId)
}

// 验证用户密码并返回编号
func (m *merchantService) Verify(usr, pwd string) (int, error) {
	usr = strings.ToLower(strings.TrimSpace(usr))
	pwd = strings.TrimSpace(pwd)
	if usr == "" || pwd == "" {
		return 0, member.ErrCredential
	}
	mchId := m._query.Verify(usr, pwd)
	if mchId <= 0 {
		return mchId, merchant.ErrNoSuchMerchant
	}
	mch := m._mchRep.GetMerchant(mchId)
	return mchId, mch.Stat()
}

// 获取企业信息
func (m *merchantService) GetReviewedEnterpriseInfo(mchId int) *merchant.EnterpriseInfo {
	mch := m._mchRep.GetMerchant(mchId)
	if mch != nil {
		return mch.ProfileManager().GetReviewedEnterpriseInfo()
	}
	return nil
}

// 获取企业信息,并返回是否为提交的信息
func (m *merchantService) GetReviewingEnterpriseInfo(mchId int) (
	e *merchant.EnterpriseInfo, isPost bool) {
	mch := m._mchRep.GetMerchant(mchId)
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
	e *merchant.EnterpriseInfo) (int64, error) {
	mch := m._mchRep.GetMerchant(mchId)
	if mch != nil {
		return mch.ProfileManager().SaveEnterpriseInfo(e)
	}
	return 0, merchant.ErrNoSuchMerchant
}

// 审核企业信息
func (m *merchantService) ReviewEnterpriseInfo(mchId int, pass bool,
	remark string) error {
	mch := m._mchRep.GetMerchant(mchId)
	if mch != nil {
		return mch.ProfileManager().ReviewEnterpriseInfo(pass, remark)
	}
	return merchant.ErrNoSuchMerchant
}

func (m *merchantService) GetMerchant(mchId int) *merchant.Merchant {
	mch := m._mchRep.GetMerchant(mchId)
	if mch != nil {
		v := mch.GetValue()
		return &v
	}
	return nil
}

func (m *merchantService) GetAccount(mchId int64) *merchant.Account {
	return m._mchRep.GetAccount(mchId)
}

func (m *merchantService) SaveMerchant(mchId int, v *merchant.Merchant) (int64, error) {
	var mch merchant.IMerchant
	var err error
	var isCreate bool
	v.Id = mchId

	if mchId > 0 {
		mch = m._mchRep.GetMerchant(mchId)
	} else {
		isCreate = true
		mch = m._mchRep.CreateMerchant(v)
	}
	if mch == nil {
		return 0, merchant.ErrNoSuchMerchant
	}
	err = mch.SetValue(v)
	if err == nil {
		mchId, err = mch.Save()
		if isCreate {
			m.initializeMerchant(mchId)
		}
	}
	return mchId, err
}

func (m *merchantService) initializeMerchant(mchId int) {

	// 初始化会员默认等级
	// m._mchRep.GetMerchant(mchId)

	//conf := merchant.DefaultSaleConf
	//conf.MerchantId = mch.GetAggregateRootId()
	// 保存销售设置
	//mch.ConfManager().SaveSaleConf(&conf)

	// 初始化销售标签
	m._saleRep.GetSale(mchId).LabelManager().InitSaleLabels()
}

// 获取商户的状态
func (m *merchantService) Stat(mchId int) error {
	mch := m._mchRep.GetMerchant(mchId)
	if mch != nil {
		return mch.Stat()
	}
	return merchant.ErrNoSuchMerchant
}

// 设置商户启用或停用
func (m *merchantService) SetEnabled(mchId int, enabled bool) error {
	mch := m._mchRep.GetMerchant(mchId)
	if mch == nil {
		return merchant.ErrNoSuchMerchant
	}
	return mch.SetEnabled(enabled)
}

// 根据主机查询商户编号
func (m *merchantService) GetMerchantIdByHost(host string) int {
	return m._query.QueryMerchantIdByHost(host)
}

// 获取商户的域名
func (m *merchantService) GetMerchantMajorHost(mchId int64) string {
	mch := m._mchRep.GetMerchant(mchId)
	if mch != nil {
		return mch.GetMajorHost()
	}
	return ""
}

func (m *merchantService) SaveSaleConf(mchId int, v *merchant.SaleConf) error {
	mch := m._mchRep.GetMerchant(mchId)
	if mch != nil {
		return mch.ConfManager().SaveSaleConf(v)
	}
	return merchant.ErrNoSuchMerchant
}

func (m *merchantService) GetSaleConf(mchId int) *merchant.SaleConf {
	mch := m._mchRep.GetMerchant(mchId)
	if mch != nil {
		conf := mch.ConfManager().GetSaleConf()
		return &conf
	}
	return nil
}

func (m *merchantService) GetShopsOfMerchant(mchId int) []*shop.Shop {
	mch := m._mchRep.GetMerchant(mchId)
	shops := mch.ShopManager().GetShops()
	sv := make([]*shop.Shop, len(shops))
	for i, v := range shops {
		vv := v.GetValue()
		sv[i] = &vv
	}
	return sv
}

// 获取商城
func (m *merchantService) GetOnlineShops1(mchId int) []*shop.Shop {
	mch := m._mchRep.GetMerchant(mchId)
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

// 获取线上店铺
func (m *merchantService) GetOnlineShopOfVendor(mchId int) *shop.ShopDto {
	mch := m._mchRep.GetMerchant(mchId)
	if mch != nil {
		return mch.ShopManager().GetOnlineShop().Data()
	}
	return nil
}

// 修改密码
func (m *merchantService) ModifyPassword(mchId int, oldPwd, newPwd string) error {
	mch := m._mchRep.GetMerchant(mchId)
	if mch != nil {
		return mch.ProfileManager().ModifyPassword(newPwd, oldPwd)
	}
	return merchant.ErrNoSuchMerchant
}

func (m *merchantService) GetMerchantsId() []int64 {
	return m._mchRep.GetMerchantsId()
}

// 保存API信息
func (m *merchantService) SaveApiInfo(mchId int, d *merchant.ApiInfo) error {
	mch := m._mchRep.GetMerchant(mchId)
	if mch != nil {
		return mch.ApiManager().SaveApiInfo(d)
	}
	return merchant.ErrNoSuchMerchant
}

// 获取API接口
func (m *merchantService) GetApiInfo(mchId int) *merchant.ApiInfo {
	mch := m._mchRep.GetMerchant(mchId)
	v := mch.ApiManager().GetApiInfo()
	return &v
}

// 启用/停用接口权限
func (m *merchantService) ApiPerm(mchId int64, enabled bool) error {
	mch := m._mchRep.GetMerchant(mchId)
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
func (m *merchantService) GetMemberLevels(mchId int64) []*merchant.MemberLevel {
	mch := m._mchRep.GetMerchant(mchId)
	if mch != nil {
		return mch.LevelManager().GetLevelSet()
	}
	return []*merchant.MemberLevel{}
}

// 根据编号获取会员等级信息
func (m *merchantService) GetMemberLevelById(mchId, id int) *merchant.MemberLevel {
	mch := m._mchRep.GetMerchant(mchId)
	if mch != nil {
		return mch.LevelManager().GetLevelById(id)
	}
	return nil
}

// 保存会员等级信息
func (m *merchantService) SaveMemberLevel(mchId int64, v *merchant.MemberLevel) (int64, error) {
	mch := m._mchRep.GetMerchant(mchId)
	if mch != nil {
		return mch.LevelManager().SaveLevel(v)
	}
	return 0, merchant.ErrNoSuchMerchant
}

// 删除会员等级
func (m *merchantService) DelMemberLevel(mchId, levelId int) error {
	mch := m._mchRep.GetMerchant(mchId)
	if mch != nil {
		return mch.LevelManager().DeleteLevel(levelId)
	}
	return merchant.ErrNoSuchMerchant
}

// 获取等级
func (m *merchantService) GetLevel(mchId, level int) *merchant.MemberLevel {
	mch := m._mchRep.GetMerchant(mchId)
	if mch != nil {
		return mch.LevelManager().GetLevelByValue(level)
	}
	return nil
}

// 获取下一个等级
func (m *merchantService) GetNextLevel(mchId, levelValue int) *merchant.MemberLevel {
	mch := m._mchRep.GetMerchant(mchId)
	if mch != nil {
		return mch.LevelManager().GetNextLevel(levelValue)
	}
	return nil

}

// 获取键值字典
func (m *merchantService) GetKeyMapsByKeyword(mchId int64, keyword string) map[string]string {
	mch := m._mchRep.GetMerchant(mchId)
	if mch != nil {
		return mch.KvManager().GetsByChar(keyword)
	}
	return map[string]string{}
}

// 保存键值字典
func (m *merchantService) SaveKeyMaps(mchId int64, data map[string]string) error {
	mch := m._mchRep.GetMerchant(mchId)
	if mch != nil {
		return mch.KvManager().Sets(data)
	}
	return merchant.ErrNoSuchMerchant
}

// 查询分页订单
func (m *merchantService) PagedOrdersOfVendor(vendorId, begin, size int, pagination bool,
	where, orderBy string) (int, []*dto.PagedVendorOrder) {
	return m._orderQuery.PagedOrdersOfVendor(vendorId, begin, size, pagination, where, orderBy)
}

// 提到会员账户
func (m *merchantService) TakeToMemberAccount(mchId int, amount float32) error {
	mch := m._mchRep.GetMerchant(mchId)
	if mch != nil {
		acc := mch.Account()
		return acc.TransferToMember(amount)
	}
	return merchant.ErrNoSuchMerchant
}

// 提到会员账户
func (m *merchantService) TakeToMemberAccount1(mchId int, amount float32) error {
	mch := m._mchRep.GetMerchant(mchId)
	if mch != nil {
		acc := mch.Account()
		return acc.TransferToMember1(amount)
	}
	return merchant.ErrNoSuchMerchant
}

//
////商户利润修改
//func (m *merchantService) UpdateMechOfflineRate(id int, rate float32, return_rate float32) error {
//	return m._mchRep.UpdateMechOfflineRate(id, rate, return_rate)
//}
//
////获取当前商家的利润
//func (m *merchantService) GetOfflineRate(id int) (float32, float32, error) {
//	return m._mchRep.GetOfflineRate(id)
//}
//
////修改当前账户信息
//func (m *merchantService) TakeOutBankCard(mchId int, amount float32) error {
//	account := m.GetAccount(mchId)
//	account.Balance = account.Balance - amount
//	err := m._mchRep.UpdateAccount(account)
//	return err
//}
//
////添加商户提取日志
//func (m *merchantService) TakeOutBankCardLog(memberId int64, mchId int, amount float32) {
//	o := &merchant.BalanceLog{
//		MchId:      mchId,
//		Kind:       100,
//		Title:      "商户提现",
//		OuterNo:    "00002",
//		Amount:     amount * (-1),
//		CsnAmount:  0,
//		State:      1,
//		CreateTime: time.Now().Unix(),
//		UpdateTime: time.Now().Unix(),
//	}
//	m._mchRep.SaveMachBlanceLog(o)
//
//	v := &member.PresentLog{
//		MemberId:     memberId,
//		BusinessKind: merchant.KindＭachTakeOutToBankCard,
//		OuterNo:      "00000000",
//		Title:        "商户提现到银行卡",
//		Amount:       amount * (-1),
//		CsnFee:       0,
//		State:        1,
//		CreateTime:   time.Now().Unix(),
//		UpdateTime:   time.Now().Unix(),
//	}
//	m._mchRep.SavePresionBlanceLog(v)
//}
//
//func (m *merchantService) UpdateMachAccount(account *merchant.Account) {
//	m._mchRep.UpdateAccount(account)
//}
//func (m *merchantService) SaveMachBlanceLog(v *merchant.BalanceLog) {
//	m._mchRep.SaveMachBlanceLog(v)
//}
//
//// 充值到商户账户
//func (m *merchantService) ChargeMachAccountByKind(memberId, machId int,
//	kind int, title string, outerNo string, amount float32, relateUser int) error {
//	if amount <= 0 || math.IsNaN(float64(amount)) {
//		return member.ErrIncorrectAmount
//	}
//	unix := time.Now().Unix()
//	v := &member.PresentLog{
//		MemberId:     memberId,
//		BusinessKind: kind,
//		Title:        title,
//		OuterNo:      outerNo,
//		Amount:       amount,
//		State:        1,
//		RelateUser:   relateUser,
//		CreateTime:   unix,
//		UpdateTime:   unix,
//	}
//
//	o := &merchant.BalanceLog{
//		MchId:      machId,
//		Kind:       kind,
//		Title:      title,
//		OuterNo:    "00002",
//		Amount:     amount,
//		CsnAmount:  0,
//		State:      1,
//		CreateTime: time.Now().Unix(),
//		UpdateTime: time.Now().Unix(),
//	}
//	m._mchRep.SaveMachBlanceLog(o)
//	_, err := m._memberRep.SavePresentLog(v)
//	if err == nil {
//		machAcc := m.GetAccount(machId)
//		machAcc.Balance = machAcc.Balance + amount
//		machAcc.UpdateTime = unix
//		m.UpdateMachAccount(machAcc)
//	}
//	return err
//}
//
//// 确认提现
//func (a *merchantService) ConfirmApplyCash(memberId int64, infoId int,
//	pass bool, remark string) error {
//	m := a._memberRep.GetMember(memberId)
//	if m == nil {
//		return member.ErrNoSuchMember
//	}
//	v := a._memberRep.GetPresentLog(infoId)
//	if v.BusinessKind != merchant.KindＭachTakeOutToBankCard {
//		return errors.New("非商户提现")
//	}
//	if pass {
//		v.State = enum.ReviewPass
//	} else {
//		if v.State == enum.ReviewReject {
//			return dm.ErrState
//		}
//		v.Remark += "失败:" + remark
//		v.State = enum.ReviewReject
//		mach := a.GetMerchantByMemberId(v.MemberId)
//		err := a.ChargeMachAccountByKind(memberId, mach.Id,
//			merchant.KindＭachTakOutRefund,
//			"商户提现退回", v.OuterNo, (-v.Amount),
//			member.DefaultRelateUser)
//		if err != nil {
//			return err
//		}
//		v.UpdateTime = time.Now().Unix()
//		_, err1 := a._memberRep.SavePresentLog(v)
//		return err1
//	}
//
//	return nil
//}
//>>>>>>> echo3
