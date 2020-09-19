/**
 * Copyright 2014 @ to2.net.
 * name :
 * author : jarryliu
 * date : 2013-12-19 22:49
 * description :
 * history :
 */

package impl

import (
	"context"
	"github.com/ixre/gof/types"
	de "go2o/core/domain/interface/domain"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/domain/interface/merchant/wholesaler"
	"go2o/core/dto"
	"go2o/core/infrastructure/domain"
	"go2o/core/query"
	"go2o/core/service/proto"
	"strconv"
	"strings"
	"time"
)

var _ proto.MerchantServiceServer = new(merchantService)

type merchantService struct {
	_mchRepo    merchant.IMerchantRepo
	_memberRepo member.IMemberRepo
	_query      *query.MerchantQuery
	_orderQuery *query.OrderQuery
	serviceUtil
}

func (m *merchantService) GetAllTradeConf(_ context.Context, i *proto.Int64) (*proto.STradeConfListResponse, error) {
	panic("implement me")
}

func (m *merchantService) CreateMerchant(_ context.Context, r *proto.MerchantCreateRequest) (*proto.Result, error) {
	mch := r.Mch
	v := &merchant.Merchant{
		LoginUser:   mch.LoginUser,
		LoginPwd:    domain.MerchantSha1Pwd(mch.LoginPwd),
		Name:        mch.Name,
		SelfSales:   int16(mch.SelfSales),
		MemberId:    r.RelMemberId,
		Level:       0,
		Logo:        "",
		CompanyName: "",
		Province:    0,
		City:        0,
		District:    0,
	}
	im := m._mchRepo.CreateMerchant(v)
	err := im.SetValue(v)
	if err == nil {
		_, err = im.Save()
		if err == nil {
			o := shop.OnlineShop{
				ShopName:   mch.ShopName,
				Logo:       mch.ShopLogo,
				Host:       "",
				Alias:      "",
				Tel:        "",
				Addr:       "",
				ShopTitle:  "",
				ShopNotice: "",
			}
			_, err = im.ShopManager().CreateOnlineShop(&o)
		}
	}
	if err == nil {
		return m.success(map[string]string{
			"mch_id": strconv.Itoa(int(im.GetAggregateRootId())),
		}), nil
	}
	return m.result(err), nil
}

func (m *merchantService) GetTradeConf(_ context.Context, r *proto.TradeConfRequest) (*proto.STradeConf, error) {
	mch := m._mchRepo.GetMerchant(int(r.MchId))
	if mch != nil {
		v := mch.ConfManager().GetTradeConf(int(r.TradeType))
		if v != nil {
			return m.parseTradeConfDto(v), nil
		}
	}
	return nil, nil
}

func (m *merchantService) SaveTradeConf(_ context.Context, r *proto.TradeConfSaveRequest) (rsp *proto.Result, err error) {
	mch := m._mchRepo.GetMerchant(int(r.MchId))
	if mch == nil {
		err = merchant.ErrNoSuchMerchant
	} else {
		var dst []*merchant.TradeConf
		for _, v := range r.Arr {
			dst = append(dst, m.parseTradeConf(v))
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
	province int32, city int32, district int32) (int64, error) {
	unix := time.Now().Unix()
	v := &merchant.Merchant{
		MemberId: 0,
		// 用户
		LoginUser: user,
		// 密码
		LoginPwd: pwd,
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
		Province: int(province),
		// 市
		City: int(city),
		// 区
		District: int(district),
		// 是否启用
		Enabled: 1,
		// 过期时间
		ExpiresTime: time.Now().Add(time.Hour * time.Duration(24*365)).Unix(),
		// 注册时间
		CreateTime: unix,
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
			return 0, de.ErrCredential
		}
	}
	if val.State == member.StateStopped {
		return 0, member.ErrMemberLocked
	}
	return val.Id, nil
}

// 登录，返回结果(Result_)和会员编号(ID);
// Result值为：-1:会员不存在; -2:账号密码不正确; -3:账号被停用
func (m *merchantService) testLogin(user string, pwd string) (id int64, errCode int32, err error) {
	if user == "" || pwd == "" {
		return 0, 1, de.ErrCredential
	}
	if len(pwd) != 32 {
		return -1, 4, de.ErrNotMD5Format
	}
	//尝试作为独立的商户账号登陆
	mch := m._mchRepo.GetMerchantByLoginUser(user)
	if mch == nil {
		// 使用会员身份登录
		var id int64
		id, err = m.testMemberLogin(user, domain.MemberSha1Pwd(pwd))
		if err != nil {
			return 0, 2, err
		}
		if mch2 := m.GetMerchantByMemberId(id); mch2 != nil {
			return mch2.Id, 0, nil
		}
		return 0, 2, merchant.ErrNoSuchMerchant
	}
	mv := mch.GetValue()
	if pwd := domain.MerchantSha1Pwd(pwd); pwd != mv.LoginPwd {
		return 0, 1, de.ErrCredential
	}
	return mch.GetAggregateRootId(), 0, nil
}

// 验证用户密码,并返回编号。可传入商户或会员的账号密码
func (m *merchantService) CheckLogin(_ context.Context, u *proto.MchUserPwd) (*proto.Result, error) {
	user := strings.ToLower(strings.TrimSpace(u.User))
	pwd := strings.TrimSpace(u.Pwd)
	id, code, err := m.testLogin(user, pwd)
	if err != nil {
		return m.errorCodeResult(int(code), err), nil
	}
	return m.success(map[string]string{"mch_id": types.String(id)}), nil
}

// 获取企业信息,并返回是否为提交的信息
func (m *merchantService) GetEnterpriseInfo(mchId int64) (
	e *merchant.EnterpriseInfo) {
	mch := m._mchRepo.GetMerchant(int(mchId))
	return mch.ProfileManager().GetEnterpriseInfo()
}

// 保存企业信息
func (m *merchantService) SaveEnterpriseInfo(mchId int64,
	e *merchant.EnterpriseInfo) (int32, error) {
	mch := m._mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		return mch.ProfileManager().SaveEnterpriseInfo(e)
	}
	return 0, merchant.ErrNoSuchMerchant
}

// 审核企业信息
func (m *merchantService) ReviewEnterpriseInfo(mchId int64, pass bool,
	remark string) error {
	mch := m._mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		return mch.ProfileManager().ReviewEnterpriseInfo(pass, remark)
	}
	return merchant.ErrNoSuchMerchant
}

func (m *merchantService) GetMerchant(_ context.Context, id *proto.Int64) (*proto.SMerchant, error) {
	mch := m._mchRepo.GetMerchant(int(id.Value))
	if mch != nil {
		c := mch.Complex()
		return m.parseMerchantDto(c), nil
	}
	return nil, nil
}

//func (m *merchantService) GetMerchant(mchId int64) *merchant.Merchant {
//	mch := m._mchRepo.GetMerchant(int(mchId))
//	if mch != nil {
//		v := mch.GetValue()
//		return &v
//	}
//	return nil
//}

func (m *merchantService) GetAccount(mchId int) *merchant.Account {
	return m._mchRepo.GetAccount(mchId)
}

func (m *merchantService) SaveMerchant(mchId int64, v *merchant.Merchant) (int64, error) {
	var mch merchant.IMerchant
	var err error
	var isCreate bool
	v.Id = mchId
	if mchId > 0 {
		mch = m._mchRepo.GetMerchant(int(mchId))
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

func (m *merchantService) initializeMerchant(mchId int64) {

	// 初始化会员默认等级
	// m._mchRepo.GetMerchant(int(mchId))

	//conf := merchant.DefaultSaleConf
	//conf.VendorId = mch.GetAggregateRootId()
	// 保存销售设置
	//mch.ConfManager().SaveSaleConf(&conf)

	// 初始化销售标签
	//m._saleRepo.GetSale(mchId).LabelManager().InitSaleLabels()
}

// 获取商户的状态
func (m *merchantService) Stat(_ context.Context, mchId *proto.Int64) (r *proto.Result, err error) {
	mch := m._mchRepo.GetMerchant(int(mchId.Value))
	if mch == nil {
		err = merchant.ErrNoSuchMerchant
	} else {
		err = mch.Stat()
	}
	return m.result(err), nil
}

// 设置商户启用或停用
func (m *merchantService) SetEnabled(mchId int64, enabled bool) error {
	mch := m._mchRepo.GetMerchant(int(mchId))
	if mch == nil {
		return merchant.ErrNoSuchMerchant
	}
	return mch.SetEnabled(enabled)
}

// 根据主机查询商户编号
func (m *merchantService) GetMerchantIdByHost(host string) int64 {
	return m._query.QueryMerchantIdByHost(host)
}

// 获取商户的域名
func (m *merchantService) GetMerchantMajorHost(mchId int) string {
	mch := m._mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		return mch.GetMajorHost()
	}
	return ""
}

func (m *merchantService) SaveSaleConf(mchId int64, v *merchant.SaleConf) error {
	mch := m._mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		return mch.ConfManager().SaveSaleConf(v)
	}
	return merchant.ErrNoSuchMerchant
}

func (m *merchantService) GetSaleConf(mchId int64) *merchant.SaleConf {
	mch := m._mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		conf := mch.ConfManager().GetSaleConf()
		return &conf
	}
	return nil
}

func (m *merchantService) GetShopsOfMerchant(mchId int64) []*shop.Shop {
	mch := m._mchRepo.GetMerchant(int(mchId))
	shops := mch.ShopManager().GetShops()
	sv := make([]*shop.Shop, len(shops))
	for i, v := range shops {
		vv := v.GetValue()
		sv[i] = &vv
	}
	return sv
}

// 修改密码
func (m *merchantService) ModifyPassword(mchId int64, oldPwd, newPwd string) error {
	mch := m._mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		return mch.ProfileManager().ModifyPassword(newPwd, oldPwd)
	}
	return merchant.ErrNoSuchMerchant
}

func (m *merchantService) GetMerchantsId() []int32 {
	return m._mchRepo.GetMerchantsId()
}

// 保存API信息
func (m *merchantService) SaveApiInfo(mchId int64, d *merchant.ApiInfo) error {
	mch := m._mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		return mch.ApiManager().SaveApiInfo(d)
	}
	return merchant.ErrNoSuchMerchant
}

// 获取API接口
func (m *merchantService) GetApiInfo(mchId int) *merchant.ApiInfo {
	mch := m._mchRepo.GetMerchant(int(mchId))
	v := mch.ApiManager().GetApiInfo()
	return &v
}

// 启用/停用接口权限
func (m *merchantService) ApiPerm(mchId int64, enabled bool) error {
	mch := m._mchRepo.GetMerchant(int(mchId))
	if enabled {
		return mch.ApiManager().EnableApiPerm()
	}
	return mch.ApiManager().DisableApiPerm()
}

// 根据API ID获取MerchantId
func (m *merchantService) GetMerchantIdByApiId(apiId string) int64 {
	return m._mchRepo.GetMerchantIdByApiId(apiId)
}

// 获取所有会员等级
func (m *merchantService) GetMemberLevels(mchId int64) []*merchant.MemberLevel {
	mch := m._mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		return mch.LevelManager().GetLevelSet()
	}
	return []*merchant.MemberLevel{}
}

// 根据编号获取会员等级信息
func (m *merchantService) GetMemberLevelById(mchId, id int32) *merchant.MemberLevel {
	mch := m._mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		return mch.LevelManager().GetLevelById(id)
	}
	return nil
}

// 保存会员等级信息
func (m *merchantService) SaveMemberLevel(mchId int64, v *merchant.MemberLevel) (int32, error) {
	mch := m._mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		return mch.LevelManager().SaveLevel(v)
	}
	return 0, merchant.ErrNoSuchMerchant
}

// 删除会员等级
func (m *merchantService) DelMemberLevel(mchId, levelId int32) error {
	mch := m._mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		return mch.LevelManager().DeleteLevel(levelId)
	}
	return merchant.ErrNoSuchMerchant
}

// 获取等级
func (m *merchantService) GetLevel(mchId, level int32) *merchant.MemberLevel {
	mch := m._mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		return mch.LevelManager().GetLevelByValue(level)
	}
	return nil
}

// 获取下一个等级
func (m *merchantService) GetNextLevel(mchId, levelValue int32) *merchant.MemberLevel {
	mch := m._mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		return mch.LevelManager().GetNextLevel(levelValue)
	}
	return nil

}

// 获取键值字典
func (m *merchantService) GetKeyMapsByKeyword(mchId int64, keyword string) map[string]string {
	mch := m._mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		return mch.KvManager().GetsByChar(keyword)
	}
	return map[string]string{}
}

// 保存键值字典
func (m *merchantService) SaveKeyMaps(mchId int64, data map[string]string) error {
	mch := m._mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		return mch.KvManager().Sets(data)
	}
	return merchant.ErrNoSuchMerchant
}

// 查询分页订单
func (m *merchantService) PagedNormalOrderOfVendor(vendorId int64, begin, size int, pagination bool,
	where, orderBy string) (int, []*dto.PagedVendorOrder) {
	return m._orderQuery.PagedNormalOrderOfVendor(vendorId, begin, size, pagination, where, orderBy)
}

// 查询分页订单
func (m *merchantService) PagedWholesaleOrderOfVendor(vendorId int64, begin, size int, pagination bool,
	where, orderBy string) (int, []*dto.PagedVendorOrder) {
	return m._orderQuery.PagedWholesaleOrderOfVendor(vendorId, begin, size, pagination, where, orderBy)
}

// 查询分页订单
func (m *merchantService) PagedTradeOrderOfVendor(vendorId int64, begin, size int, pagination bool,
	where, orderBy string) (int32, []*proto.SComplexOrder) {
	return m._orderQuery.PagedTradeOrderOfVendor(vendorId, begin, size, pagination, where, orderBy)
}

// 提到会员账户
func (m *merchantService) TakeToMemberAccount(mchId int64, amount float32) error {
	mch := m._mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		acc := mch.Account()
		return acc.TransferToMember(amount)
	}
	return merchant.ErrNoSuchMerchant
}

// 提到会员账户
func (m *merchantService) TakeToMemberAccount1(mchId int64, amount float32) error {
	mch := m._mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		acc := mch.Account()
		return acc.TransferToMember1(amount)
	}
	return merchant.ErrNoSuchMerchant
}

// 账户充值
func (m *merchantService) ChargeAccount(mchId int64, kind int32, title,
	outerNo string, amount float64, relateUser int64) error {
	mch := m._mchRepo.GetMerchant(int(mchId))
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
//	v := &member.WalletAccountLog{
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
//	v := &member.WalletAccountLog{
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
//	_, err := m._memberRepo.SaveWalletAccountLog(v)
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
//		_, err1 := a._memberRepo.SaveWalletAccountLog(v)
//		return err1
//	}
//
//	return nil
//}
//

// 获取

// 同步批发商品
func (m *merchantService) SyncWholesaleItem(_ context.Context, vendorId *proto.Int64) (*proto.SyncWSItemsResponse, error) {
	mch := m._mchRepo.GetMerchant(int(vendorId.Value))
	var mp = map[string]int32{
		"add": 0, "del": 0,
	}
	if mch != nil {
		mp = mch.Wholesaler().SyncItems(true)
	}
	return &proto.SyncWSItemsResponse{Value: mp}, nil
}

func (m *merchantService) GetMchBuyerGroup_(mchId, id int64) *merchant.MchBuyerGroup {
	mch := m._mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		return mch.ConfManager().GetGroupByGroupId(int32(id))
	}
	return nil
}

// 保存
func (m *merchantService) SaveMchBuyerGroup_(mchId int64, v *merchant.MchBuyerGroup) (r *proto.Result, err error) {
	mch := m._mchRepo.GetMerchant(int(mchId))
	if mch == nil {
		err = merchant.ErrNoSuchMerchant
	} else {
		_, err = mch.ConfManager().SaveMchBuyerGroup(v)
	}
	return m.result(err), nil
}

// 获取买家分组
func (m *merchantService) GetBuyerGroups(mchId int64) []*merchant.BuyerGroup {
	mch := m._mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		return mch.ConfManager().SelectBuyerGroup()
	}
	return []*merchant.BuyerGroup{}
}

// 获取批发返点率
func (m *merchantService) GetRebateRate(mchId, groupId int64) []*wholesaler.WsRebateRate {
	mch := m._mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		return mch.Wholesaler().GetGroupRebateRate(int32(groupId))
	}
	return []*wholesaler.WsRebateRate{}
}

// 保存分组返点率
func (m *merchantService) SaveGroupRebateRate(mchId, groupId int64,
	arr []*wholesaler.WsRebateRate) error {
	mch := m._mchRepo.GetMerchant(int(mchId))
	if mch == nil {
		return merchant.ErrNoSuchMerchant
	}
	return mch.Wholesaler().SaveGroupRebateRate(int32(groupId), arr)
}

func (m *merchantService) parseMerchantDto(src *merchant.ComplexMerchant) *proto.SMerchant {
	return &proto.SMerchant{
		Id:            src.Id,
		MemberId:      src.MemberId,
		LoginUser:     src.Usr,
		LoginPwd:      src.Pwd,
		Name:          src.Name,
		SelfSales:     int32(src.SelfSales),
		Level:         src.Level,
		Logo:          src.Logo,
		CompanyName:   src.CompanyName,
		Province:      src.Province,
		City:          src.City,
		District:      src.District,
		Enabled:       src.Enabled,
		LastLoginTime: int32(src.LastLoginTime),
	}
}

func (m *merchantService) parseTradeConf(conf *proto.STradeConf) *merchant.TradeConf {
	return &merchant.TradeConf{
		//MchId:       conf.MchId,
		//TradeType:   int(conf.TradeType),
		//PlanId:      conf.PlanId,
		//Flag:        int(conf.Flag),
		//AmountBasis: int(conf.AmountBasis),
		//TradeFee:    int(conf.TradeFee),
		//TradeRate:   int(conf.TradeRate),
	}
}

func (m *merchantService) parseTradeConfDto(conf *merchant.TradeConf) *proto.STradeConf {
	return &proto.STradeConf{
		//MchId:       conf.MchId,
		//TradeType:   int32(conf.TradeType),
		//PlanId:      conf.PlanId,
		//Flag:        int32(conf.Flag),
		//AmountBasis: int32(conf.AmountBasis),
		//TradeFee:    int32(conf.TradeFee),
		//TradeRate:   int32(conf.TradeRate),
	}
}
