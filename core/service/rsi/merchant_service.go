/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-19 22:49
 * description :
 * history :
 */

package rsi

import (
	"context"
	"github.com/ixre/gof/util"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/domain/interface/merchant/wholesaler"
	"go2o/core/dto"
	"go2o/core/infrastructure/domain"
	"go2o/core/query"
	"go2o/core/service/auto_gen/rpc/mch_service"
	"go2o/core/service/auto_gen/rpc/order_service"
	"go2o/core/service/auto_gen/rpc/ttype"
	"go2o/core/service/thrift/parser"
	"strings"
	"time"
)

var _ mch_service.MerchantService = new(merchantService)

type merchantService struct {
	_mchRepo    merchant.IMerchantRepo
	_memberRepo member.IMemberRepo
	_query      *query.MerchantQuery
	_orderQuery *query.OrderQuery
	serviceUtil
}

func (m *merchantService) GetAllTradeConf(ctx context.Context, mchId int32) (r []*mch_service.STradeConf, err error) {
	mch := m._mchRepo.GetMerchant(mchId)
	var arr []*mch_service.STradeConf
	if mch != nil {
		for _, v := range mch.ConfManager().GetAllTradeConf() {
			arr = append(arr, parser.TradeConfDto(v))
		}
	}
	return arr, nil
}

func (m *merchantService) GetTradeConf(ctx context.Context, mchId int32, tradeType int32) (r *mch_service.STradeConf, err error) {
	mch := m._mchRepo.GetMerchant(mchId)
	if mch != nil {
		v := mch.ConfManager().GetTradeConf(int(tradeType))
		if v != nil {
			return parser.TradeConfDto(v), nil
		}
	}
	return nil, nil
}

func (m *merchantService) SaveTradeConf(ctx context.Context, mchId int32, arr []*mch_service.STradeConf) (r *ttype.Result_, err error) {
	mch := m._mchRepo.GetMerchant(mchId)
	if mch == nil {
		err = merchant.ErrNoSuchMerchant
	} else {
		var dst []*merchant.TradeConf
		for _, v := range arr {
			dst = append(dst, parser.TradeConf(v))
		}
		err = mch.ConfManager().SaveTradeConf(dst)
	}
	return m.result(err), nil
}

func NewMerchantService(r merchant.IMerchantRepo, memberRepo member.IMemberRepo,
	q *query.MerchantQuery, orderQuery *query.OrderQuery) *merchantService {
	return &merchantService{
		_mchRepo:    r,
		_memberRepo: memberRepo,
		_query:      q,
		_orderQuery: orderQuery,
	}
}

// 创建会员申请商户密钥
func (m *merchantService) CreateSignUpToken(memberId int64) string {
	return m._mchRepo.CreateSignUpToken(memberId)
}

// 根据商户申请密钥获取会员编号
func (m *merchantService) GetMemberFromSignUpToken(token string) int64 {
	return m._mchRepo.GetMemberFromSignUpToken(token)
}

// 获取会员商户申请信息
func (m *merchantService) GetMchSignUpInfoByMemberId(memberId int64) *merchant.MchSignUp {
	return m._mchRepo.GetManager().GetSignUpInfoByMemberId(memberId)
}

// 获取商户申请信息
func (m *merchantService) GetSignUp(id int32) *merchant.MchSignUp {
	return m._mchRepo.GetManager().GetSignUpInfo(id)
}

// 审核商户申请信息
func (m *merchantService) ReviewSignUp(id int32, pass bool, remark string) error {
	return m._mchRepo.GetManager().ReviewMchSignUp(id, pass, remark)
}

// 商户注册
func (m *merchantService) SignUp(user, pwd, companyName string,
	province int32, city int32, district int32) (int32, error) {
	unix := time.Now().Unix()
	v := &merchant.Merchant{
		MemberId: 0,
		// 用户
		Usr: user,
		// 密码
		Pwd: pwd,
		// 商户名称
		Name: companyName,
		// 是否自营
		SelfSales: 0,
		// 商户等级
		Level: 1,
		// 标志
		Logo:        "",
		CompanyName: companyName,
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
	mch := m._mchRepo.CreateMerchant(v)
	err := mch.SetValue(v)
	if err == nil {
		return mch.Save()
	}
	return -1, err
}

// 提交注册信息
func (m *merchantService) SignUpPost(e *merchant.MchSignUp) (int32, error) {
	return m._mchRepo.GetManager().CommitSignUpInfo(e)
}

func (m *merchantService) GetMerchantByMemberId(memberId int64) *merchant.Merchant {
	mch := m._mchRepo.GetManager().GetMerchantByMemberId(memberId)
	if mch != nil {
		v := mch.GetValue()
		return &v
	}
	return nil
}

// 删除会员的商户申请资料
func (m *merchantService) RemoveMerchantSignUp(memberId int64) error {
	return m._mchRepo.GetManager().RemoveSignUp(memberId)
}

// 登录，返回结果(Result_)和会员编号(ID);
// Result值为：-1:会员不存在; -2:账号密码不正确; -3:账号被停用
func (m *merchantService) testMemberLogin(user string, pwd string) (id int64, err error) {
	user = strings.ToLower(strings.TrimSpace(user))
	val := m._memberRepo.GetMemberByUser(user)
	if val == nil {
		val = m._memberRepo.GetMemberValueByPhone(user)
	}
	if val == nil {
		return 0, member.ErrNoSuchMember
	}
	if val.Pwd != pwd {
		//todo: 兼容旧密码
		if val.Pwd != domain.Sha1(pwd) {
			return 0, member.ErrCredential
		}
	}
	if val.State == member.StateStopped {
		return 0, member.ErrMemberLocked
	}
	return val.Id, nil
}

// 验证用户密码,并返回编号。可传入商户或会员的账号密码
func (m *merchantService) CheckLogin(ctx context.Context, user, oriPwd string) (r *ttype.Result_, err error) {
	user = strings.ToLower(strings.TrimSpace(user))
	oriPwd = strings.TrimSpace(oriPwd)
	var mchId int32
	if user == "" || oriPwd == "" {
		return m.error(member.ErrCredential), nil
	}
	//尝试作为独立的商户账号登陆
	encPwd := domain.MerchantSha1Pwd(user, oriPwd)
	mchId = m._query.Verify(user, encPwd)
	if mchId <= 0 {
		// 使用会员身份登录
		var id int64
		mEncPwd := domain.MemberSha1Pwd(domain.Md5(oriPwd))
		id, err = m.testMemberLogin(user, mEncPwd)
		if err == nil {
			mch := m.GetMerchantByMemberId(id)
			if mch != nil {
				mchId = mch.ID
			}
		}
	}
	if mchId < 0 && err == nil {
		err = merchant.ErrNoSuchMerchant
	}
	if err != nil {
		return m.error(err), nil
	}
	return m.success(map[string]string{"mch_id": util.Str(mchId)}), nil
}

// 获取企业信息,并返回是否为提交的信息
func (m *merchantService) GetEnterpriseInfo(mchId int32) (
	e *merchant.EnterpriseInfo) {
	mch := m._mchRepo.GetMerchant(mchId)
	return mch.ProfileManager().GetEnterpriseInfo()
}

// 保存企业信息
func (m *merchantService) SaveEnterpriseInfo(mchId int32,
	e *merchant.EnterpriseInfo) (int32, error) {
	mch := m._mchRepo.GetMerchant(mchId)
	if mch != nil {
		return mch.ProfileManager().SaveEnterpriseInfo(e)
	}
	return 0, merchant.ErrNoSuchMerchant
}

// 审核企业信息
func (m *merchantService) ReviewEnterpriseInfo(mchId int32, pass bool,
	remark string) error {
	mch := m._mchRepo.GetMerchant(mchId)
	if mch != nil {
		return mch.ProfileManager().ReviewEnterpriseInfo(pass, remark)
	}
	return merchant.ErrNoSuchMerchant
}

func (m *merchantService) Complex(ctx context.Context, mchId int32) (*mch_service.SComplexMerchant, error) {
	mch := m._mchRepo.GetMerchant(mchId)
	if mch != nil {
		c := mch.Complex()
		return parser.MerchantDto(c), nil
	}
	return nil, nil
}

func (m *merchantService) GetMerchant(mchId int32) *merchant.Merchant {
	mch := m._mchRepo.GetMerchant(mchId)
	if mch != nil {
		v := mch.GetValue()
		return &v
	}
	return nil
}

func (m *merchantService) GetAccount(mchId int32) *merchant.Account {
	return m._mchRepo.GetAccount(mchId)
}

func (m *merchantService) SaveMerchant(mchId int32, v *merchant.Merchant) (int32, error) {
	var mch merchant.IMerchant
	var err error
	var isCreate bool
	v.ID = mchId

	if mchId > 0 {
		mch = m._mchRepo.GetMerchant(mchId)
	} else {
		isCreate = true
		mch = m._mchRepo.CreateMerchant(v)
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

func (m *merchantService) initializeMerchant(mchId int32) {

	// 初始化会员默认等级
	// m._mchRepo.GetMerchant(mchId)

	//conf := merchant.DefaultSaleConf
	//conf.VendorId = mch.GetAggregateRootId()
	// 保存销售设置
	//mch.ConfManager().SaveSaleConf(&conf)

	// 初始化销售标签
	//m._saleRepo.GetSale(mchId).LabelManager().InitSaleLabels()
}

// 获取商户的状态
func (m *merchantService) Stat(ctx context.Context, mchId int32) (r *ttype.Result_, err error) {
	mch := m._mchRepo.GetMerchant(mchId)
	if mch == nil {
		err = merchant.ErrNoSuchMerchant
	} else {
		err = mch.Stat()
	}
	return m.result(err), nil
}

// 设置商户启用或停用
func (m *merchantService) SetEnabled(mchId int32, enabled bool) error {
	mch := m._mchRepo.GetMerchant(mchId)
	if mch == nil {
		return merchant.ErrNoSuchMerchant
	}
	return mch.SetEnabled(enabled)
}

// 根据主机查询商户编号
func (m *merchantService) GetMerchantIdByHost(host string) int32 {
	return m._query.QueryMerchantIdByHost(host)
}

// 获取商户的域名
func (m *merchantService) GetMerchantMajorHost(mchId int32) string {
	mch := m._mchRepo.GetMerchant(mchId)
	if mch != nil {
		return mch.GetMajorHost()
	}
	return ""
}

func (m *merchantService) SaveSaleConf(mchId int32, v *merchant.SaleConf) error {
	mch := m._mchRepo.GetMerchant(mchId)
	if mch != nil {
		return mch.ConfManager().SaveSaleConf(v)
	}
	return merchant.ErrNoSuchMerchant
}

func (m *merchantService) GetSaleConf(mchId int32) *merchant.SaleConf {
	mch := m._mchRepo.GetMerchant(mchId)
	if mch != nil {
		conf := mch.ConfManager().GetSaleConf()
		return &conf
	}
	return nil
}

func (m *merchantService) GetShopsOfMerchant(mchId int32) []*shop.Shop {
	mch := m._mchRepo.GetMerchant(mchId)
	shops := mch.ShopManager().GetShops()
	sv := make([]*shop.Shop, len(shops))
	for i, v := range shops {
		vv := v.GetValue()
		sv[i] = &vv
	}
	return sv
}

// 修改密码
func (m *merchantService) ModifyPassword(mchId int32, oldPwd, newPwd string) error {
	mch := m._mchRepo.GetMerchant(mchId)
	if mch != nil {
		return mch.ProfileManager().ModifyPassword(newPwd, oldPwd)
	}
	return merchant.ErrNoSuchMerchant
}

func (m *merchantService) GetMerchantsId() []int32 {
	return m._mchRepo.GetMerchantsId()
}

// 保存API信息
func (m *merchantService) SaveApiInfo(mchId int32, d *merchant.ApiInfo) error {
	mch := m._mchRepo.GetMerchant(mchId)
	if mch != nil {
		return mch.ApiManager().SaveApiInfo(d)
	}
	return merchant.ErrNoSuchMerchant
}

// 获取API接口
func (m *merchantService) GetApiInfo(mchId int32) *merchant.ApiInfo {
	mch := m._mchRepo.GetMerchant(mchId)
	v := mch.ApiManager().GetApiInfo()
	return &v
}

// 启用/停用接口权限
func (m *merchantService) ApiPerm(mchId int32, enabled bool) error {
	mch := m._mchRepo.GetMerchant(mchId)
	if enabled {
		return mch.ApiManager().EnableApiPerm()
	}
	return mch.ApiManager().DisableApiPerm()
}

// 根据API ID获取MerchantId
func (m *merchantService) GetMerchantIdByApiId(apiId string) int32 {
	return m._mchRepo.GetMerchantIdByApiId(apiId)
}

// 获取所有会员等级
func (m *merchantService) GetMemberLevels(mchId int32) []*merchant.MemberLevel {
	mch := m._mchRepo.GetMerchant(mchId)
	if mch != nil {
		return mch.LevelManager().GetLevelSet()
	}
	return []*merchant.MemberLevel{}
}

// 根据编号获取会员等级信息
func (m *merchantService) GetMemberLevelById(mchId, id int32) *merchant.MemberLevel {
	mch := m._mchRepo.GetMerchant(mchId)
	if mch != nil {
		return mch.LevelManager().GetLevelById(id)
	}
	return nil
}

// 保存会员等级信息
func (m *merchantService) SaveMemberLevel(mchId int32, v *merchant.MemberLevel) (int32, error) {
	mch := m._mchRepo.GetMerchant(mchId)
	if mch != nil {
		return mch.LevelManager().SaveLevel(v)
	}
	return 0, merchant.ErrNoSuchMerchant
}

// 删除会员等级
func (m *merchantService) DelMemberLevel(mchId, levelId int32) error {
	mch := m._mchRepo.GetMerchant(mchId)
	if mch != nil {
		return mch.LevelManager().DeleteLevel(levelId)
	}
	return merchant.ErrNoSuchMerchant
}

// 获取等级
func (m *merchantService) GetLevel(mchId, level int32) *merchant.MemberLevel {
	mch := m._mchRepo.GetMerchant(mchId)
	if mch != nil {
		return mch.LevelManager().GetLevelByValue(level)
	}
	return nil
}

// 获取下一个等级
func (m *merchantService) GetNextLevel(mchId, levelValue int32) *merchant.MemberLevel {
	mch := m._mchRepo.GetMerchant(mchId)
	if mch != nil {
		return mch.LevelManager().GetNextLevel(levelValue)
	}
	return nil

}

// 获取键值字典
func (m *merchantService) GetKeyMapsByKeyword(mchId int32, keyword string) map[string]string {
	mch := m._mchRepo.GetMerchant(mchId)
	if mch != nil {
		return mch.KvManager().GetsByChar(keyword)
	}
	return map[string]string{}
}

// 保存键值字典
func (m *merchantService) SaveKeyMaps(mchId int32, data map[string]string) error {
	mch := m._mchRepo.GetMerchant(mchId)
	if mch != nil {
		return mch.KvManager().Sets(data)
	}
	return merchant.ErrNoSuchMerchant
}

// 查询分页订单
func (m *merchantService) PagedNormalOrderOfVendor(vendorId int32, begin, size int, pagination bool,
	where, orderBy string) (int, []*dto.PagedVendorOrder) {
	return m._orderQuery.PagedNormalOrderOfVendor(vendorId, begin, size, pagination, where, orderBy)
}

// 查询分页订单
func (m *merchantService) PagedWholesaleOrderOfVendor(vendorId int32, begin, size int, pagination bool,
	where, orderBy string) (int, []*dto.PagedVendorOrder) {
	return m._orderQuery.PagedWholesaleOrderOfVendor(vendorId, begin, size, pagination, where, orderBy)
}

// 查询分页订单
func (m *merchantService) PagedTradeOrderOfVendor(vendorId int32, begin, size int, pagination bool,
	where, orderBy string) (int32, []*order_service.SComplexOrder) {
	return m._orderQuery.PagedTradeOrderOfVendor(vendorId, begin, size, pagination, where, orderBy)
}

// 提到会员账户
func (m *merchantService) TakeToMemberAccount(mchId int32, amount float32) error {
	mch := m._mchRepo.GetMerchant(mchId)
	if mch != nil {
		acc := mch.Account()
		return acc.TransferToMember(amount)
	}
	return merchant.ErrNoSuchMerchant
}

// 提到会员账户
func (m *merchantService) TakeToMemberAccount1(mchId int32, amount float32) error {
	mch := m._mchRepo.GetMerchant(mchId)
	if mch != nil {
		acc := mch.Account()
		return acc.TransferToMember1(amount)
	}
	return merchant.ErrNoSuchMerchant
}

// 账户充值
func (m *merchantService) ChargeAccount(mchId int32, kind int32, title,
	outerNo string, amount float64, relateUser int64) error {
	mch := m._mchRepo.GetMerchant(mchId)
	if mch == nil {
		return merchant.ErrNoSuchMerchant
	}
	return mch.Account().Charge(kind, amount, title, outerNo, relateUser)
}

//
////商户利润修改
//func (m *merchantService) UpdateMechOfflineRate(id int32, rate float32, return_rate float32) error {
//	return m._mchRepo.UpdateMechOfflineRate(id, rate, return_rate)
//}
//
////获取当前商家的利润
//func (m *merchantService) GetOfflineRate(id int32) (float32, float32, error) {
//	return m._mchRepo.GetOfflineRate(id)
//}
//
////修改当前账户信息
//func (m *merchantService) TakeOutBankCard(mchId  int32, amount float32) error {
//	account := m.GetAccount(mchId)
//	account.Balance = account.Balance - amount
//	err := m._mchRepo.UpdateAccount(account)
//	return err
//}
//
////添加商户提取日志
//func (m *merchantService) TakeOutBankCardLog(memberId  int32, mchId  int32, amount float32) {
//	o := &merchant.BalanceLog{
//		MchId:      mchId,
//		Kind:       100,
//		Title:      "商户提现",
//		OuterNo:    "00002",
//		Amount:     amount * (-1),
//		ProcedureFee:  0,
//		State:      1,
//		CreateTime: time.Now().Unix(),
//		UpdateTime: time.Now().Unix(),
//	}
//	m._mchRepo.SaveMachBlanceLog(o)
//
//	v := &member.MWalletLog{
//		MemberId:     memberId,
//		Kind: merchant.KindＭachTakeOutToBankCard,
//		OuterNo:      "00000000",
//		Title:        "商户提现到银行卡",
//		Amount:       amount * (-1),
//		CsnFee:       0,
//		State:        1,
//		CreateTime:   time.Now().Unix(),
//		UpdateTime:   time.Now().Unix(),
//	}
//	m._mchRepo.SavePresionBlanceLog(v)
//}
//
//func (m *merchantService) UpdateMachAccount(account *merchant.Account) {
//	m._mchRepo.UpdateAccount(account)
//}
//func (m *merchantService) SaveMachBlanceLog(v *merchant.BalanceLog) {
//	m._mchRepo.SaveMachBlanceLog(v)
//}
//
//// 充值到商户账户
//func (m *merchantService) ChargeMachAccountByKind(memberId, machId int32,
//	kind int, title string, outerNo string, amount float32, relateUser int) error {
//	if amount <= 0 || math.IsNaN(float64(amount)) {
//		return member.ErrIncorrectAmount
//	}
//	unix := time.Now().Unix()
//	v := &member.MWalletLog{
//		MemberId:     memberId,
//		Kind: kind,
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
//		ProcedureFee:  0,
//		State:      1,
//		CreateTime: time.Now().Unix(),
//		UpdateTime: time.Now().Unix(),
//	}
//	m._mchRepo.SaveMachBlanceLog(o)
//	_, err := m._memberRepo.SavePresentLog(v)
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
//func (a *merchantService) ConfirmApplyCash(memberId  int32, infoId int32,
//	pass bool, remark string) error {
//	m := a._memberRepo.GetMember(memberId)
//	if m == nil {
//		return member.ErrNoSuchMember
//	}
//	v := a._memberRepo.GetWalletLog(infoId)
//	if v.Kind != merchant.KindＭachTakeOutToBankCard {
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
//		err := a.ChargeMachAccountByKind(memberId, mach.ID,
//			merchant.KindＭachTakOutRefund,
//			"商户提现退回", v.OuterNo, (-v.Amount),
//			member.DefaultRelateUser)
//		if err != nil {
//			return err
//		}
//		v.UpdateTime = time.Now().Unix()
//		_, err1 := a._memberRepo.SavePresentLog(v)
//		return err1
//	}
//
//	return nil
//}
//>>>>>>> echo3

// 获取

// 同步批发商品
func (m *merchantService) SyncWholesaleItem(ctx context.Context, vendorId int32) (map[string]int32, error) {
	mch := m._mchRepo.GetMerchant(vendorId)
	if mch != nil {
		return mch.Wholesaler().SyncItems(true), nil
	}
	return map[string]int32{
		"add": 0, "del": 0,
	}, nil
}

func (m *merchantService) GetMchBuyerGroup_(mchId, id int32) *merchant.MchBuyerGroup {
	mch := m._mchRepo.GetMerchant(mchId)
	if mch != nil {
		return mch.ConfManager().GetGroupByGroupId(id)
	}
	return nil
}

// 保存
func (m *merchantService) SaveMchBuyerGroup_(mchId int32, v *merchant.MchBuyerGroup) (r *ttype.Result_, err error) {
	mch := m._mchRepo.GetMerchant(mchId)
	if mch == nil {
		err = merchant.ErrNoSuchMerchant
	} else {
		_, err = mch.ConfManager().SaveMchBuyerGroup(v)
	}
	return m.result(err), nil
}

// 获取买家分组
func (m *merchantService) GetBuyerGroups(mchId int32) []*merchant.BuyerGroup {
	mch := m._mchRepo.GetMerchant(mchId)
	if mch != nil {
		return mch.ConfManager().SelectBuyerGroup()
	}
	return []*merchant.BuyerGroup{}
}

// 获取批发返点率
func (m *merchantService) GetRebateRate(mchId, groupId int32) []*wholesaler.WsRebateRate {
	mch := m._mchRepo.GetMerchant(mchId)
	if mch != nil {
		return mch.Wholesaler().GetGroupRebateRate(groupId)
	}
	return []*wholesaler.WsRebateRate{}
}

// 保存分组返点率
func (m *merchantService) SaveGroupRebateRate(mchId, groupId int32,
	arr []*wholesaler.WsRebateRate) error {
	mch := m._mchRepo.GetMerchant(mchId)
	if mch == nil {
		return merchant.ErrNoSuchMerchant
	}
	return mch.Wholesaler().SaveGroupRebateRate(groupId, arr)
}
