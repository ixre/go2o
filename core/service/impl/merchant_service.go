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
	"go2o/core/domain/interface/order"
	"go2o/core/dto"
	"go2o/core/infrastructure/domain"
	"go2o/core/query"
	"go2o/core/service/parser"
	"go2o/core/service/proto"
	"strconv"
	"strings"
)

var _ proto.MerchantServiceServer = new(merchantService)

type merchantService struct {
	_mchRepo    merchant.IMerchantRepo
	_memberRepo member.IMemberRepo
	_query      *query.MerchantQuery
	_orderQuery *query.OrderQuery
	serviceUtil
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
func (m *merchantService) CreateSignUpToken(_ context.Context, id *proto.MemberId) (*proto.String, error) {
	s := m._mchRepo.CreateSignUpToken(id.Value)
	return &proto.String{Value: s}, nil
}

// 根据商户申请密钥获取会员编号
func (m *merchantService) GetMemberFromSignUpToken(_ context.Context, s *proto.String) (*proto.Int64, error) {
	i := m._mchRepo.GetMemberFromSignUpToken(s.Value)
	return &proto.Int64{Value: i}, nil
}

// 提交注册信息
func (m *merchantService) SignUp(_ context.Context, up *proto.SMchSignUp) (*proto.Result, error) {
	im := m._mchRepo.GetManager()
	_, err := im.CommitSignUpInfo(m.parseMchSignUp(up))
	return m.error(err), nil
}

// 获取会员商户申请信息
func (m *merchantService) GetMchSignUpId(_ context.Context, id *proto.MemberId) (*proto.Int64, error) {
	v := m._mchRepo.GetManager().GetSignUpInfoByMemberId(id.Value)
	if v != nil {
		return &proto.Int64{Value: int64(v.Id)}, nil
	}
	return &proto.Int64{}, nil
}

// 获取商户申请信息
func (m *merchantService) GetSignUp(_ context.Context, id *proto.Int64) (*proto.SMchSignUp, error) {
	im := m._mchRepo.GetManager()
	v := im.GetSignUpInfo(int32(id.Value))
	if v != nil {
		return m.parseMchSIgnUpDto(v), nil
	}
	return nil, nil
}

// 审核商户申请信息
func (m *merchantService) ReviewSignUp(_ context.Context, r *proto.MchReviewRequest) (*proto.Result, error) {
	im := m._mchRepo.GetManager()
	err := im.ReviewMchSignUp(int32(r.MerchantId), r.Pass, r.Remark)
	return m.error(err), nil
}

// 删除会员的商户申请资料
func (m *merchantService) RemoveMerchantSignUp(_ context.Context, id *proto.MemberId) (*proto.Result, error) {
	err := m._mchRepo.GetManager().RemoveSignUp(id.Value)
	return m.error(err), nil
}

func (m *merchantService) GetMerchantIdByMember(_ context.Context, id *proto.MemberId) (*proto.Int64, error) {
	mch := m._mchRepo.GetManager().GetMerchantByMemberId(id.Value)
	if mch != nil {
		return &proto.Int64{Value: mch.GetAggregateRootId()}, nil
	}
	return &proto.Int64{}, nil
}

// 获取企业信息,并返回是否为提交的信息
func (m *merchantService) GetEnterpriseInfo(_ context.Context, id *proto.MerchantId) (*proto.SEnterpriseInfo, error) {
	mch := m._mchRepo.GetMerchant(int(id.Value))
	if mch != nil {
		v := mch.ProfileManager().GetEnterpriseInfo()
		if v != nil {
			return m.parseEnterpriseInfoDto(v), nil
		}
	}
	return nil, nil
}

// 保存企业信息
func (m *merchantService) SaveEnterpriseInfo(_ context.Context, r *proto.SaveEnterpriseRequest) (*proto.Result, error) {
	mch := m._mchRepo.GetMerchant(int(r.MerchantId))
	var err error
	if mch == nil {
		err = merchant.ErrNoSuchMerchant
	} else {
		_, err = mch.ProfileManager().SaveEnterpriseInfo(
			m.parseEnterpriseInfo(r.Value))
	}
	return m.error(err), nil
}

// 审核企业信息
func (m *merchantService) ReviewEnterpriseInfo(_ context.Context, r *proto.MchReviewRequest) (*proto.Result, error) {
	mch := m._mchRepo.GetMerchant(int(r.MerchantId))
	var err error
	if mch == nil {
		err = merchant.ErrNoSuchMerchant
	} else {
		im := mch.ProfileManager()
		err = im.ReviewEnterpriseInfo(r.Pass, r.Remark)
	}
	return m.error(err), nil
}

// 获取商户账户
func (m *merchantService) GetAccount(_ context.Context, id *proto.MerchantId) (*proto.SMerchantAccount, error) {
	v := m._mchRepo.GetAccount(int(id.Value))
	if v != nil {
		return m.parseAccountDto(v), nil
	}
	return nil, nil
}

// 设置商户启用或停用
func (m *merchantService) SetEnabled(_ context.Context, r *proto.MerchantDisableRequest) (*proto.Result, error) {
	mch := m._mchRepo.GetMerchant(int(r.MerchantId))
	var err error
	if mch == nil {
		err = merchant.ErrNoSuchMerchant
	} else {
		err = mch.SetEnabled(r.Enabled)
	}
	return m.error(err), nil
}

// 根据主机查询商户编号
func (m *merchantService) GetMerchantIdByHost(_ context.Context, host *proto.String) (*proto.Int64, error) {
	id := m._query.QueryMerchantIdByHost(host.Value)
	return &proto.Int64{Value: id}, nil
}

// 获取商户的域名
func (m *merchantService) GetMerchantMajorHost(_ context.Context, id *proto.MerchantId) (*proto.String, error) {
	mch := m._mchRepo.GetMerchant(int(id.Value))
	if mch != nil {
		return &proto.String{
			Value: mch.GetMajorHost(),
		}, nil
	}
	return &proto.String{}, nil
}

func (m *merchantService) SaveSaleConf(_ context.Context, r *proto.SaveMerchantSaleConfRequest) (*proto.Result, error) {
	mch := m._mchRepo.GetMerchant(int(r.MerchantId))
	var err error
	if mch == nil {
		err = merchant.ErrNoSuchMerchant
	} else {
		err = mch.ConfManager().SaveSaleConf(m.parseSaleConf(r.Value))
	}
	return m.error(err), nil
}

func (m *merchantService) GetSaleConf(_ context.Context, id *proto.MerchantId) (*proto.SMerchantSaleConf, error) {
	mch := m._mchRepo.GetMerchant(int(id.Value))
	if mch != nil {
		conf := mch.ConfManager().GetSaleConf()
		return m.parseSaleConfDto(conf), nil
	}
	return nil, nil
}

// 获取商户的店铺编号
func (m *merchantService) GetShopId(_ context.Context, id *proto.MerchantId) (*proto.Int64, error) {
	mch := m._mchRepo.GetMerchant(int(id.Value))
	shops := mch.ShopManager().GetShops()
	for _, v := range shops {
		return &proto.Int64{Value: int64(v.GetDomainId())}, nil
	}
	return &proto.Int64{}, nil
}

// 修改密码
func (m *merchantService) ModifyPassword(_ context.Context, r *proto.ModifyMerchantPasswordRequest) (*proto.Result, error) {
	mch := m._mchRepo.GetMerchant(int(r.MerchantId))
	var err error
	if mch == nil {
		err = merchant.ErrNoSuchMerchant
	} else {
		if l := len(r.Origin); l > 0 && l < 32 {
			err = de.ErrNotMD5Format
		} else if len(r.Password) != 32 {
			err = de.ErrNotMD5Format
		} else {
			err = mch.ProfileManager().ModifyPassword(r.Origin, r.Password)
		}
	}
	return m.error(err), nil
}

// 获取API接口
func (m *merchantService) GetApiInfo(_ context.Context, id *proto.MerchantId) (*proto.SMerchantApiInfo, error) {
	mch := m._mchRepo.GetMerchant(int(id.Value))
	if mch != nil {
		v := mch.ApiManager().GetApiInfo()
		return m.parseApiDto(v), nil
	}
	return nil, nil
}

// 启用/停用接口权限
func (m *merchantService) ToggleApiPerm(_ context.Context, r *proto.MerchantApiPermRequest) (*proto.Result, error) {
	mch := m._mchRepo.GetMerchant(int(r.MerchantId))
	im := mch.ApiManager()
	var err error
	if r.Enabled {
		err = im.EnableApiPerm()
	} else {
		err = im.DisableApiPerm()
	}
	return m.error(err), nil
}

// 根据API ID获取MerchantId
func (m *merchantService) GetMerchantIdByApiId(_ context.Context, s *proto.String) (*proto.Int64, error) {
	i := m._mchRepo.GetMerchantIdByApiId(s.Value)
	return &proto.Int64{Value: i}, nil
}

// 查询分页订单
func (m *merchantService) PagedNormalOrderOfVendor(_ context.Context, r *proto.MerchantOrderRequest) (*proto.PagingMerchantOrderListResponse, error) {
	total, list := m._orderQuery.PagedNormalOrderOfVendor(
		r.MerchantId,
		int(r.Params.Begin),
		int(r.Params.End-r.Params.Begin),
		r.Pagination,
		r.Params.Where,
		r.Params.SortBy)
	ret := &proto.PagingMerchantOrderListResponse{
		Total: int64(total),
		Data:  make([]*proto.SMerchantOrder, len(list)),
	}
	for i, v := range list {
		ret.Data[i] = m.parseOrder(v)
	}
	return ret, nil
}

// 查询分页订单
func (m *merchantService) PagedWholesaleOrderOfVendor(_ context.Context, r *proto.MerchantOrderRequest) (*proto.PagingMerchantOrderListResponse, error) {
	total, list := m._orderQuery.PagedWholesaleOrderOfVendor(
		r.MerchantId,
		int(r.Params.Begin),
		int(r.Params.End-r.Params.Begin),
		r.Pagination,
		r.Params.Where,
		r.Params.SortBy)
	ret := &proto.PagingMerchantOrderListResponse{
		Total: int64(total),
		Data:  make([]*proto.SMerchantOrder, len(list)),
	}
	for i, v := range list {
		ret.Data[i] = m.parseOrder(v)
	}
	return ret, nil
}

func (m *merchantService) PagedTradeOrderOfVendor(_ context.Context, r *proto.MerchantOrderRequest) (*proto.PagingMerchantOrderListResponse, error) {
	total, list := m._orderQuery.PagedTradeOrderOfVendor(
		r.MerchantId,
		int(r.Params.Begin),
		int(r.Params.End-r.Params.Begin),
		r.Pagination,
		r.Params.Where,
		r.Params.SortBy)
	ret := &proto.PagingMerchantOrderListResponse{
		Total: int64(total),
		Data:  make([]*proto.SMerchantOrder, len(list)),
	}
	for i, v := range list {
		ret.Data[i] = m.parseOrder(v)
	}
	return ret, nil
}

// 提到会员账户
func (m *merchantService) WithdrawToMemberAccount(_ context.Context, r *proto.WithdrawToMemberAccountRequest) (*proto.Result, error) {
	mch := m._mchRepo.GetMerchant(int(r.MerchantId))
	var err error
	if mch == nil {
		err = merchant.ErrNoSuchMerchant
	} else {
		acc := mch.Account()
		err = acc.TransferToMember(float32(r.Amount))
	}
	return m.error(err), nil
}

// 账户充值
func (m *merchantService) ChargeAccount(_ context.Context, r *proto.MerchantChargeRequest) (*proto.Result, error) {
	mch := m._mchRepo.GetMerchant(int(r.MerchantId))
	var err error
	if mch == nil {
		err = merchant.ErrNoSuchMerchant
	} else {
		err = mch.Account().Charge(r.Kind, r.Amount, r.Title, r.OuterNo, r.RelateUserId)
	}
	return m.error(err), nil
}

func (m *merchantService) GetMchBuyerGroup_(_ context.Context, id *proto.MerchantBuyerGroupId) (*proto.SMerchantBuyerGroup, error) {
	mch := m._mchRepo.GetMerchant(int(id.MerchantId))
	if mch != nil {
		v := mch.ConfManager().GetGroupByGroupId(int32(id.GroupId))
		return m.parseGroupDto(v), nil
	}
	return nil, nil
}

// 保存
func (m *merchantService) SaveMchBuyerGroup_(_ context.Context, r *proto.SaveMerchantBuyerGroupRequest) (*proto.Result, error) {
	mch := m._mchRepo.GetMerchant(int(r.MerchantId))
	var err error
	if mch == nil {
		err = merchant.ErrNoSuchMerchant
	} else {
		v := m.parseGroup(r.Value)
		v.MchId = r.MerchantId
		//v.GroupId =
		_, err = mch.ConfManager().SaveMchBuyerGroup(v)
	}
	return m.result(err), nil
}

//todo: mchBuyerGroup去调还是BuyerGroup去调

//// 获取买家分组
//func (m *merchantService) GetBuyerGroups(_ context.Context, id *proto.MerchantId) (*proto.MerchantBuyerGroupListResponse, error) {
//	mch := m._mchRepo.GetMerchant(int(id.Value))
//	if mch != nil {
//		list := mch.ConfManager().SelectBuyerGroup()
//		arr := make([]*proto.SMerchantBuyerGroup,len(list))
//		for i,v := range list{
//			arr[i] = m.parseGroupDto(v)
//		}
//	}
//	return []*merchant.BuyerGroup{}
//}

// 获取批发返点率
func (m *merchantService) GetRebateRate(_ context.Context, id *proto.MerchantBuyerGroupId) (*proto.WholesaleRebateRateListResponse, error) {
	mch := m._mchRepo.GetMerchant(int(id.MerchantId))
	ret := &proto.WholesaleRebateRateListResponse{
		Value: make([]*proto.SWholesaleRebateRate, 0),
	}
	if mch != nil {
		arr := mch.Wholesaler().GetGroupRebateRate(int32(id.GroupId))
		for _, v := range arr {
			ret.Value = append(ret.Value, m.parseRebateRateDto(v))
		}
	}
	return ret, nil
}

// 保存分组返点率
func (m *merchantService) SaveGroupRebateRate(_ context.Context, r *proto.SaveWholesaleRebateRateRequest) (*proto.Result, error) {
	mch := m._mchRepo.GetMerchant(int(r.MerchantId))
	var err error
	if mch == nil {
		err = merchant.ErrNoSuchMerchant
	} else {
		arr := make([]*wholesaler.WsRebateRate, len(r.Value))
		for i, v := range r.Value {
			arr[i] = m.parseRebateRate(v)
		}
		err = mch.Wholesaler().SaveGroupRebateRate(int32(r.GroupId), arr)
	}
	return m.error(err), nil
}

func (m *merchantService) GetAllTradeConf_(_ context.Context, i *proto.Int64) (*proto.STradeConfListResponse, error) {
	return &proto.STradeConfListResponse{
		Value: make([]*proto.STradeConf_, 0),
	}, nil
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

func (m *merchantService) GetTradeConf(_ context.Context, r *proto.TradeConfRequest) (*proto.STradeConf_, error) {
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

// 登录，返回结果(Result_)和会员编号(Id);
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

// 登录，返回结果(Result_)和会员编号(Id);
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
		mchId, _ := m.GetMerchantIdByMember(context.TODO(), &proto.MemberId{Value: id})
		if mchId.Value > 0 {
			return mchId.Value, 0, nil
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

func (m *merchantService) GetMerchant(_ context.Context, id *proto.Int64) (*proto.SMerchant, error) {
	mch := m._mchRepo.GetMerchant(int(id.Value))
	if mch != nil {
		c := mch.Complex()
		return m.parseMerchantDto(c), nil
	}
	return nil, nil
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

// 保存API信息
func (m *merchantService) SaveApiInfo(mchId int64, d *merchant.ApiInfo) error {
	mch := m._mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		return mch.ApiManager().SaveApiInfo(d)
	}
	return merchant.ErrNoSuchMerchant
}

// 获取所有会员等级
func (m *merchantService) GetMemberLevels_(mchId int64) []*merchant.MemberLevel {
	mch := m._mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		return mch.LevelManager().GetLevelSet()
	}
	return []*merchant.MemberLevel{}
}

// 根据编号获取会员等级信息
func (m *merchantService) GetMemberLevelById_(mchId, id int32) *merchant.MemberLevel {
	mch := m._mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		return mch.LevelManager().GetLevelById(id)
	}
	return nil
}

// 保存会员等级信息
func (m *merchantService) SaveMemberLevel_(mchId int64, v *merchant.MemberLevel) (int32, error) {
	mch := m._mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		return mch.LevelManager().SaveLevel(v)
	}
	return 0, merchant.ErrNoSuchMerchant
}

// 删除会员等级
func (m *merchantService) DelMemberLevel_(mchId, levelId int32) error {
	mch := m._mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		return mch.LevelManager().DeleteLevel(levelId)
	}
	return merchant.ErrNoSuchMerchant
}

// 获取等级
func (m *merchantService) GetLevel_(mchId, level int32) *merchant.MemberLevel {
	mch := m._mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		return mch.LevelManager().GetLevelByValue(level)
	}
	return nil
}

// 获取下一个等级
func (m *merchantService) GetNextLevel_(mchId, levelValue int32) *merchant.MemberLevel {
	mch := m._mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		return mch.LevelManager().GetNextLevel(levelValue)
	}
	return nil

}

// 获取键值字典
func (m *merchantService) GetKeyMapsByKeyword_(mchId int64, keyword string) map[string]string {
	mch := m._mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		return mch.KvManager().GetsByChar(keyword)
	}
	return map[string]string{}
}

// 保存键值字典
func (m *merchantService) SaveKeyMaps_(mchId int64, data map[string]string) error {
	mch := m._mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		return mch.KvManager().Sets(data)
	}
	return merchant.ErrNoSuchMerchant
}

// 查询分页订单

// 提到会员账户
func (m *merchantService) WithdrawToMemberAccount1(mchId int64, amount float32) error {
	mch := m._mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		acc := mch.Account()
		return acc.TransferToMember1(amount)
	}
	return merchant.ErrNoSuchMerchant
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
//		err := a.ChargeMachAccountByKind(memberId, mach.Id,
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

func (m *merchantService) parseMerchantDto(src *merchant.ComplexMerchant) *proto.SMerchant {
	return &proto.SMerchant{
		Id:            src.Id,
		MemberId:      src.MemberId,
		LoginUser:     src.Usr,
		LoginPwd:      src.Pwd,
		Name:          src.Name,
		SelfSales:     src.SelfSales,
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

func (m *merchantService) parseTradeConf(conf *proto.STradeConf_) *merchant.TradeConf {
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

func (m *merchantService) parseTradeConfDto(conf *merchant.TradeConf) *proto.STradeConf_ {
	return &proto.STradeConf_{
		//MchId:       conf.MchId,
		//TradeType:   int32(conf.TradeType),
		//PlanId:      conf.PlanId,
		//Flag:        int32(conf.Flag),
		//AmountBasis: int32(conf.AmountBasis),
		//TradeFee:    int32(conf.TradeFee),
		//TradeRate:   int32(conf.TradeRate),
	}
}

func (m *merchantService) parseMchSignUp(v *proto.SMchSignUp) *merchant.MchSignUp {
	return &merchant.MchSignUp{
		Id:           int32(v.Id),
		SignNo:       v.SignNo,
		MemberId:     v.MemberId,
		Usr:          v.User,
		Pwd:          v.Pwd,
		MchName:      v.MchName,
		Province:     v.Province,
		City:         v.City,
		District:     v.District,
		Address:      v.Address,
		ShopName:     v.ShopName,
		CompanyName:  v.CompanyName,
		CompanyNo:    v.CompanyNo,
		PersonName:   v.PersonName,
		PersonId:     v.PersonId,
		PersonImage:  v.PersonImage,
		Phone:        v.Phone,
		CompanyImage: v.CompanyImage,
		AuthDoc:      v.AuthDoc,
		Remark:       v.Remark,
		Reviewed:     v.ReviewState,
		SubmitTime:   v.SubmitTime,
	}
}

func (m *merchantService) parseMchSIgnUpDto(v *merchant.MchSignUp) *proto.SMchSignUp {
	return &proto.SMchSignUp{
		Id:           int64(v.Id),
		SignNo:       v.SignNo,
		MemberId:     v.MemberId,
		User:         v.Usr,
		Pwd:          v.Pwd,
		MchName:      v.MchName,
		Province:     v.Province,
		City:         v.City,
		District:     v.District,
		Address:      v.Address,
		ShopName:     v.ShopName,
		CompanyName:  v.CompanyName,
		CompanyNo:    v.CompanyNo,
		PersonName:   v.PersonName,
		PersonId:     v.PersonId,
		PersonImage:  v.PersonImage,
		Phone:        v.Phone,
		CompanyImage: v.CompanyImage,
		AuthDoc:      v.AuthDoc,
		Remark:       v.Remark,
		ReviewState:  v.Reviewed,
		SubmitTime:   v.SubmitTime,
	}
}

func (m *merchantService) parseEnterpriseInfoDto(v *merchant.EnterpriseInfo) *proto.SEnterpriseInfo {
	return &proto.SEnterpriseInfo{
		Id:           int64(v.ID),
		MerchantId:   v.MchId,
		CompanyName:  v.CompanyName,
		CompanyNo:    v.CompanyNo,
		PersonName:   v.PersonName,
		PersonIdNo:   v.PersonIdNo,
		PersonImage:  v.PersonImage,
		Telephone:    v.Tel,
		Province:     v.Province,
		City:         v.City,
		District:     v.District,
		Location:     v.Location,
		Address:      v.Address,
		CompanyImage: v.CompanyImage,
		AuthDoc:      v.AuthDoc,
		ReviewState:  v.Reviewed,
		ReviewTime:   v.ReviewTime,
		ReviewRemark: v.ReviewRemark,
		UpdateTime:   v.UpdateTime,
	}
}

func (m *merchantService) parseEnterpriseInfo(v *proto.SEnterpriseInfo) *merchant.EnterpriseInfo {
	return &merchant.EnterpriseInfo{
		ID:           int32(v.Id),
		MchId:        v.MerchantId,
		CompanyName:  v.CompanyName,
		CompanyNo:    v.CompanyNo,
		PersonName:   v.PersonName,
		PersonIdNo:   v.PersonIdNo,
		PersonImage:  v.PersonImage,
		Tel:          v.Telephone,
		Province:     v.Province,
		City:         v.City,
		District:     v.District,
		Location:     v.Location,
		Address:      v.Address,
		CompanyImage: v.CompanyImage,
		AuthDoc:      v.AuthDoc,
		Reviewed:     int32(v.ReviewState),
		ReviewTime:   v.ReviewTime,
		ReviewRemark: v.ReviewRemark,
		UpdateTime:   v.UpdateTime,
	}
}

func (m *merchantService) parseAccountDto(v *merchant.Account) *proto.SMerchantAccount {
	return &proto.SMerchantAccount{
		Balance:       float64(v.Balance),
		FreezeAmount:  float64(v.FreezeAmount),
		AwaitAmount:   float64(v.AwaitAmount),
		PresentAmount: float64(v.PresentAmount),
		SalesAmount:   float64(v.SalesAmount),
		RefundAmount:  float64(v.RefundAmount),
		TakeAmount:    float64(v.TakeAmount),
		OfflineSales:  float64(v.OfflineSales),
		UpdateTime:    v.UpdateTime,
	}
}

func (m *merchantService) parseSaleConf(v *proto.SMerchantSaleConf) *merchant.SaleConf {
	return &merchant.SaleConf{
		MerchantId:              v.MerchantId,
		FxSalesEnabled:          types.IntCond(v.FxSalesEnabled, 1, 0),
		CashBackPercent:         float32(v.CashBackPercent),
		CashBackTg1Percent:      float32(v.CashBackTg1Percent),
		CashBackTg2Percent:      float32(v.CashBackTg2Percent),
		CashBackMemberPercent:   float32(v.CashBackMemberPercent),
		AutoSetupOrder:          types.IntCond(v.AutoSetupOrder, 1, 0),
		OrderTimeOutMinute:      int(v.OrderTimeOutMinute),
		OrderConfirmAfterMinute: int(v.OrderConfirmAfterMinute),
		OrderTimeOutReceiveHour: int(v.OrderTimeOutReceiveHour),
	}
}

func (m *merchantService) parseSaleConfDto(v merchant.SaleConf) *proto.SMerchantSaleConf {
	return &proto.SMerchantSaleConf{
		MerchantId:              v.MerchantId,
		FxSalesEnabled:          v.FxSalesEnabled == 1,
		CashBackPercent:         float64(v.CashBackPercent),
		CashBackTg1Percent:      float64(v.CashBackTg1Percent),
		CashBackTg2Percent:      float64(v.CashBackTg2Percent),
		CashBackMemberPercent:   float64(v.CashBackMemberPercent),
		AutoSetupOrder:          v.AutoSetupOrder == 1,
		OrderTimeOutMinute:      int32(v.OrderTimeOutMinute),
		OrderConfirmAfterMinute: int32(v.OrderConfirmAfterMinute),
		OrderTimeOutReceiveHour: int32(v.OrderTimeOutReceiveHour),
	}
}

func (m *merchantService) parseApiDto(v merchant.ApiInfo) *proto.SMerchantApiInfo {
	arr := strings.Split(v.WhiteList, ",")
	if len(v.WhiteList) == 0 {
		arr = []string{}
	}
	return &proto.SMerchantApiInfo{
		ApiId:     v.ApiId,
		ApiSecret: v.ApiSecret,
		WhiteList: arr,
		Enabled:   v.Enabled == 1,
	}
}

func (m *merchantService) parseGroupDto(v *merchant.MchBuyerGroup) *proto.SMerchantBuyerGroup {
	return &proto.SMerchantBuyerGroup{
		Id:              int64(v.ID),
		Alias:           v.Alias,
		EnableRetail:    v.EnableRetail == 1,
		EnableWholesale: v.EnableWholesale == 1,
		RebatePeriod:    v.RebatePeriod,
	}
}

func (m *merchantService) parseGroup(v *proto.SMerchantBuyerGroup) *merchant.MchBuyerGroup {
	return &merchant.MchBuyerGroup{
		ID:              int32(v.Id),
		Alias:           v.Alias,
		EnableRetail:    int32(types.IntCond(v.EnableRetail, 1, 0)),
		EnableWholesale: int32(types.IntCond(v.EnableWholesale, 1, 0)),
		RebatePeriod:    v.RebatePeriod,
	}
}

func (m *merchantService) parseRebateRateDto(v *wholesaler.WsRebateRate) *proto.SWholesaleRebateRate {
	return &proto.SWholesaleRebateRate{
		Id:            int64(v.ID),
		WsId:          int64(v.WsId),
		BuyerGroupId:  int64(v.BuyerGid),
		RequireAmount: v.RequireAmount,
		RebateRate:    v.RebateRate,
	}
}

func (m *merchantService) parseRebateRate(v *proto.SWholesaleRebateRate) *wholesaler.WsRebateRate {
	return &wholesaler.WsRebateRate{
		ID:            int32(v.Id),
		WsId:          int32(v.WsId),
		BuyerGid:      int32(v.BuyerGroupId),
		RequireAmount: v.RequireAmount,
		RebateRate:    v.RebateRate,
	}
}

func (m *merchantService) parseOrder(v *dto.PagedVendorOrder) *proto.SMerchantOrder {
	items := make([]*proto.SOrderItem, 0)
	if v.Items != nil {
		for _, v := range v.Items {
			items = append(items, parser.ParseOrderItem(v))
		}
	}
	return &proto.SMerchantOrder{
		OrderId:        v.Id,
		OrderNo:        v.OrderNo,
		ParentNo:       v.ParentNo,
		BuyerId:        int64(v.BuyerId),
		BuyerName:      v.BuyerName,
		ItemAmount:     float64(v.ItemAmount),
		DiscountAmount: float64(v.DiscountAmount),
		ExpressFee:     float64(v.ExpressFee),
		PackageFee:     float64(v.PackageFee),
		IsPaid:         v.IsPaid,
		FinalAmount:    float64(v.FinalAmount),
		State:          int32(v.State),
		StateText:      order.OrderState(v.State).String(),
		CreateTime:     v.CreateTime,
		Items:          items,
		Data:           v.Data,
	}
}
