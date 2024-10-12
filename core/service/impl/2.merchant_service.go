/**
 * Copyright 2014 @ 56x.net.
 * name :
 * author : jarryliu
 * date : 2013-12-19 22:49
 * description :
 * history :
 */

package impl

import (
	"context"
	"errors"
	"fmt"
	"strings"

	de "github.com/ixre/go2o/core/domain/interface/domain"
	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/domain/interface/merchant"
	"github.com/ixre/go2o/core/domain/interface/merchant/staff"
	"github.com/ixre/go2o/core/domain/interface/merchant/wholesaler"
	"github.com/ixre/go2o/core/domain/interface/order"
	"github.com/ixre/go2o/core/domain/interface/wallet"
	"github.com/ixre/go2o/core/dto"
	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/go2o/core/query"
	"github.com/ixre/go2o/core/service/parser"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/domain/eventbus"
	"github.com/ixre/gof/types"
	context2 "golang.org/x/net/context"
)

var _ proto.MerchantServiceServer = new(merchantService)

type merchantService struct {
	_mchRepo    merchant.IMerchantRepo
	_memberRepo member.IMemberRepo
	_staffRepo  staff.IStaffRepo
	_query      *query.MerchantQuery
	_orderQuery *query.OrderQuery
	serviceUtil
	proto.UnimplementedMerchantServiceServer
}

// GetSettleConf implements proto.MerchantServiceServer.
func (m *merchantService) GetSettleConf(_ context.Context, req *proto.MerchantId) (*proto.SSettleConf, error) {
	im := m._mchRepo.GetMerchant(int(req.Value))
	if im == nil {
		return nil, errors.New("商户不存在")
	}
	conf := im.ConfManager()
	settle := conf.GetSettleConf()
	return &proto.SSettleConf{
		MchId:       int64(settle.MchId),
		MchName:     im.GetValue().MchName,
		OrderTxRate: float32(settle.OrderTxRate),
		OtherTxRate: float32(settle.OtherTxRate),
		SubMchNo:    settle.SubMchNo,
	}, nil

}

// SaveSettleConf implements proto.MerchantServiceServer.
func (m *merchantService) SaveSettleConf(_ context.Context, req *proto.SettleConfigSaveRequest) (*proto.TxResult, error) {
	im := m._mchRepo.GetMerchant(int(req.MchId))
	if im == nil {
		return m.errorV2(merchant.ErrNoSuchMerchant), nil
	}
	conf := im.ConfManager()
	err := conf.SaveSettleConf(&merchant.SettleConf{
		MchId:       int(req.MchId),
		OrderTxRate: float64(req.OrderTxRate),
		OtherTxRate: float64(req.OtherTxRate),
		SubMchNo:    req.SubMchNo,
	})
	return m.errorV2(err), nil
}

func NewMerchantService(r merchant.IMerchantRepo,
	memberRepo member.IMemberRepo,
	staffRepo staff.IStaffRepo,
	q *query.MerchantQuery, orderQuery *query.OrderQuery) proto.MerchantServiceServer {
	return &merchantService{
		_mchRepo:    r,
		_memberRepo: memberRepo,
		_staffRepo:  staffRepo,
		_query:      q,
		_orderQuery: orderQuery,
	}
}

// GetMerchantIdByMailAddress implements proto.MerchantServiceServer.
func (m *merchantService) GetMerchantIdByUsername(_ context.Context, mail *proto.String) (*proto.Int64, error) {
	mch := m._mchRepo.GetMerchantByUsername(mail.Value)
	if mch != nil {
		return &proto.Int64{Value: int64(mch.GetAggregateRootId())}, nil
	}
	return &proto.Int64{}, nil
}

// CreateMerchant 创建商户
func (m *merchantService) CreateMerchant(_ context.Context, r *proto.CreateMerchantRequest) (*proto.MerchantCreateResponse, error) {
	mch := r.Mch
	v := &merchant.Merchant{
		Username: r.Username,
		Password: domain.MerchantSha1Pwd(mch.Password, ""),
		MchName:  mch.MchName,
		IsSelf:   int16(r.IsSelf),
		MemberId: int(r.MemberId),
		Level:    0,
		Logo:     "",
		Province: 0,
		City:     0,
		District: 0,
	}
	im := m._mchRepo.CreateMerchant(v)
	err := im.SetValue(v)
	if err == nil && r.MemberId > 0 {
		err = im.BindMember(int(r.MemberId))
	}
	if err == nil {
		_, err = im.Save()
		if err == nil {
			// todo: 商城默认开通店铺，应单独提供方法开通店铺
			// o := shop.OnlineShop{
			// 	ShopName:   mch.MchName,
			// 	Logo:       mch.Logo,
			// 	Host:       "",
			// 	Alias:      "",
			// 	Telephone:  "",
			// 	Addr:       "",
			// 	ShopTitle:  "",
			// 	ShopNotice: "",
			// }
			// _, err = im.ShopManager().CreateOnlineShop(&o)
		}
	}
	rsp := &proto.MerchantCreateResponse{}
	if err == nil {
		rsp.MerchantId = int64(im.GetAggregateRootId())
	} else {
		rsp.ErrCode = 1
		rsp.ErrMsg = err.Error()
	}
	return rsp, nil
}

// SaveAuthenticate 提交商户认证信息
func (m *merchantService) SaveAuthenticate(_ context.Context, r *proto.SaveAuthenticateRequest) (*proto.TxResult, error) {
	mch := m._mchRepo.GetMerchant(int(r.MchId))
	if mch == nil {
		return m.errorV2(merchant.ErrNoSuchMerchant), nil
	}
	v := &merchant.Authenticate{
		OrgName:          r.OrgName,
		MchName:          r.MchName,
		LicenceNo:        r.OrgNo,
		LicencePic:       r.OrgPic,
		OrgAddress:       r.OrgAddress,
		Province:         int(r.Province),
		City:             int(r.City),
		District:         int(r.District),
		WorkCity:         int(r.WorkCity),
		QualificationPic: r.QualificationPic,
		PersonId:         r.PersonId,
		PersonName:       r.PersonName,
		PersonFrontPic:   r.PersonFrontPic,
		PersonBackPic:    r.PersonBackPic,
		PersonPhone:      r.PersonPhone,
		AuthorityPic:     r.AuthorityPic,
		BankName:         r.BankName,
		BankAccount:      r.BankAccount,
		BankAccountPic:   r.BankAccountPic,
		BankNo:           r.BankNo,
		ExtraData:        r.ExtraData,
		Version:          0,
	}
	_, err := mch.ProfileManager().SaveAuthenticate(v)
	return m.errorV2(err), nil
}

// ReviewAuthenticate 审核商户申请信息
func (m *merchantService) ReviewAuthenticate(_ context.Context, r *proto.MerchantReviewRequest) (*proto.TxResult, error) {
	mch := m._mchRepo.GetMerchant(int(r.MchId))
	if mch == nil {
		return m.errorV2(merchant.ErrNoSuchMerchant), nil
	}
	err := mch.ProfileManager().ReviewAuthenticate(r.Pass, r.Remark)
	return m.errorV2(err), nil
}

// ChangeMemberBind 更换会员绑定
func (m *merchantService) ChangeMemberBind(_ context2.Context, r *proto.ChangeMemberBindRequest) (*proto.Result, error) {
	im := m._mchRepo.GetMerchant(int(r.MerchantId))
	if im == nil {
		return m.error(merchant.ErrNoSuchMerchant), nil
	}
	mem := m._memberRepo.GetMemberByUser(r.Username)
	if mem == nil {
		return m.error(member.ErrNoSuchMember), nil
	}
	err := im.BindMember(int(mem.Id))
	if err != nil {
		return m.error(err), nil
	}
	return m.success(nil), nil
}

func (m *merchantService) GetMerchantIdByMember(_ context.Context, id *proto.MemberId) (*proto.Int64, error) {
	mch := m._mchRepo.GetManager().GetMerchantByMemberId(int(id.Value))
	if mch != nil {
		return &proto.Int64{Value: int64(mch.GetAggregateRootId())}, nil
	}
	return &proto.Int64{}, nil
}

// GetAccount 获取商户账户
func (m *merchantService) GetAccount(_ context.Context, id *proto.MerchantId) (*proto.SMerchantAccount, error) {
	v := m._mchRepo.GetAccount(int(id.Value))
	if v != nil {
		return m.parseAccountDto(v), nil
	}
	return nil, fmt.Errorf("no such merchant account")
}

// GetMerchantIdByHost 根据主机查询商户编号
func (m *merchantService) GetMerchantIdByHost(_ context.Context, host *proto.String) (*proto.Int64, error) {
	id := m._query.QueryMerchantIdByHost(host.Value)
	return &proto.Int64{Value: id}, nil
}

// GetMerchantMajorHost 获取商户的域名
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
	return nil, fmt.Errorf("no such sale conf")
}

// GetShopId 获取商户的店铺编号
func (m *merchantService) GetShopId(_ context.Context, id *proto.MerchantId) (*proto.Int64, error) {
	mch := m._mchRepo.GetMerchant(int(id.Value))
	shops := mch.ShopManager().GetShops()
	for _, v := range shops {
		return &proto.Int64{Value: int64(v.GetDomainId())}, nil
	}
	return &proto.Int64{}, nil
}

// UpdateLockStatus implements proto.MerchantServiceServer.
func (m *merchantService) UpdateLockStatus(_ context.Context, req *proto.MerchantLockStatusRequest) (*proto.TxResult, error) {
	mch := m._mchRepo.GetMerchant(int(req.MerchantId))
	if mch != nil {
		if req.Lock {
			err := mch.Lock()
			return m.errorV2(err), nil
		} else {
			err := mch.Unlock()
			return m.errorV2(err), nil
		}
	}
	return m.errorV2(merchant.ErrNoSuchMerchant), nil
}

// ChangePassword 修改密码
func (m *merchantService) ChangePassword(_ context.Context, r *proto.ModifyMerchantPasswordRequest) (*proto.TxResult, error) {
	mch := m._mchRepo.GetMerchant(int(r.MerchantId))
	var err error
	if mch == nil {
		err = merchant.ErrNoSuchMerchant
	} else {
		if l := len(r.OldPassword); l > 0 && l < 32 {
			err = de.ErrNotMD5Format
		} else if len(r.NewPassword) != 32 {
			err = de.ErrNotMD5Format
		} else {
			err = mch.ProfileManager().ChangePassword(r.NewPassword, r.OldPassword)
		}
	}
	return m.errorV2(err), nil
}

// GetApiInfo 获取API接口
func (m *merchantService) GetApiInfo(_ context.Context, id *proto.MerchantId) (*proto.SMerchantApiInfo, error) {
	mch := m._mchRepo.GetMerchant(int(id.Value))
	if mch != nil {
		v := mch.ApiManager().GetApiInfo()
		return m.parseApiDto(v), nil
	}
	return nil, fmt.Errorf("no such api info")
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
		err = acc.TransferToMember(int(r.Amount))
	}
	return m.error(err), nil
}

func (m *merchantService) GetMchBuyerGroup_(_ context.Context, id *proto.MerchantBuyerGroupId) (*proto.SMerchantBuyerGroup, error) {
	mch := m._mchRepo.GetMerchant(int(id.MerchantId))
	if mch != nil {
		v := mch.ConfManager().GetGroupByGroupId(int32(id.GroupId))
		return m.parseGroupDto(v), nil
	}
	return nil, fmt.Errorf("no such buyer group")
}

// 保存
func (m *merchantService) SaveMchBuyerGroup(_ context.Context, r *proto.SaveMerchantBuyerGroupRequest) (*proto.Result, error) {
	mch := m._mchRepo.GetMerchant(int(r.MerchantId))
	var err error
	if mch == nil {
		err = merchant.ErrNoSuchMerchant
	} else {
		v := m.parseGroup(r.Value)
		v.MerchantId = r.MerchantId
		//v.GroupId =
		_, err = mch.ConfManager().SaveMchBuyerGroup(v)
	}
	return m.result(err), nil
}

// 获取商户的(批发)买家分组
func (m *merchantService) GetBuyerGroups(_ context.Context, id *proto.MerchantId) (*proto.MerchantBuyerGroupListResponse, error) {
	mch := m._mchRepo.GetMerchant(int(id.Value))
	var arr []*proto.SMerchantBuyerGroup
	if mch != nil {
		list := mch.ConfManager().SelectBuyerGroup()
		arr = make([]*proto.SMerchantBuyerGroup, len(list))
		for i, v := range list {
			arr[i] = m.parseBuyerGroupDto(v)
		}
	}
	return &proto.MerchantBuyerGroupListResponse{
		Value: arr,
	}, nil
}

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

func (m *merchantService) GetTradeConf(_ context.Context, r *proto.TradeConfRequest) (*proto.STradeConf_, error) {
	mch := m._mchRepo.GetMerchant(int(r.MchId))
	if mch != nil {
		v := mch.ConfManager().GetTradeConf(int(r.TradeType))
		if v != nil {
			return m.parseTradeConfDto(v), nil
		}
	}
	return nil, fmt.Errorf("no such trade conf")
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
	if val.Password != pwd {
		//todo: 兼容旧密码
		if val.Password != domain.Sha1(pwd) {
			return 0, de.ErrPasswordNotMatch
		}
	}
	if (val.UserFlag & member.FlagLocked) == member.FlagLocked {
		return 0, member.ErrMemberLocked
	}
	return int64(val.Id), nil
}

// 登录，返回结果(Result_)和会员编号(Id);
// Result值为：-1:会员不存在; -2:账号密码不正确; -3:账号被停用
func (m *merchantService) testLogin(user string, pwd string) (_ merchant.IMerchantAggregateRoot, errCode int32, err error) {
	if user == "" || pwd == "" {
		return nil, 1, de.ErrPasswordNotMatch
	}
	if len(pwd) != 32 {
		return nil, 4, de.ErrNotMD5Format
	}
	//尝试作为独立的商户账号登陆
	mchList := m._mchRepo.FindList(nil, "username=?", user)
	if l := len(mchList); l == 0 {
		return nil, 5, merchant.ErrNoSuchMerchant
	}
	if len(mchList) > 1 {
		return nil, 6, errors.New("存在多个相同用户名的商户")
	}
	// if mch == nil {
	// 	// 使用会员身份登录
	// 	var id int64
	// 	id, err = m.testMemberLogin(user, domain.MemberSha1Pwd(pwd, ""))
	// 	if err != nil {
	// 		return nil, 2, err
	// 	}
	// 	mchId, _ := m.GetMerchantIdByMember(context.TODO(), &proto.MemberId{Value: id})
	// 	if mchId.Value > 0 {
	// 		mch = m._mchRepo.GetMerchant(int(mchId.Value))
	// 		return mch, 0, nil
	// 	}
	// 	return nil, 2, merchant.ErrNoSuchMerchant
	// }
	mch := m._mchRepo.CreateMerchant(mchList[0])
	mv := mch.GetValue()
	if pwd := domain.MerchantSha1Pwd(pwd, mch.GetValue().Salt); pwd != mv.Password {
		return nil, 1, de.ErrPasswordNotMatch
	}
	return mch, 0, nil
}

// CheckLogin 验证用户密码,并返回编号。可传入商户或会员的账号密码
func (m *merchantService) CheckLogin(_ context.Context, u *proto.MchUserPwdRequest) (*proto.MchLoginResponse, error) {
	user := strings.ToLower(strings.TrimSpace(u.Username))
	pwd := strings.TrimSpace(u.Password)
	mch, code, err := m.testLogin(user, pwd)
	if err != nil {
		return &proto.MchLoginResponse{
			ErrCode: code,
			ErrMsg:  err.Error(),
		}, nil
	}
	var shopId = 0
	shop := mch.ShopManager().GetOnlineShop()
	if shop != nil {
		shopId = shop.GetDomainId()
	}
	return &proto.MchLoginResponse{
		MerchantId: int64(mch.GetAggregateRootId()),
		ShopId:     int64(shopId),
	}, nil
}

func (m *merchantService) GetMerchant(_ context.Context, id *proto.Int64) (*proto.QMerchant, error) {
	mch := m._mchRepo.GetMerchant(int(id.Value))
	if mch != nil {
		c := mch.Complex()
		return m.parseMerchantDto(c), nil
	}
	return nil, merchant.ErrNoSuchMerchant
}

func (m *merchantService) SaveMerchant(_ context.Context, r *proto.SaveMerchantRequest) (*proto.Result, error) {
	mch := m._mchRepo.GetMerchant(int(r.MchId))
	if mch == nil {
		return m.error(merchant.ErrNoSuchMerchant), nil
	}
	v := mch.GetValue()
	d := r.Mch
	v.MchName = d.MchName
	v.Logo = d.Logo
	v.Level = int(d.Level)
	err := mch.SetValue(&v)
	if err == nil {
		_, err = mch.Save()
	}
	return m.result(err), nil
}

// Stat 获取商户的状态
func (m *merchantService) Stat(_ context.Context, mchId *proto.Int64) (r *proto.Result, err error) {
	mch := m._mchRepo.GetMerchant(int(mchId.Value))
	if mch == nil {
		err = merchant.ErrNoSuchMerchant
	} else {
		err = mch.Stat()
	}
	return m.result(err), nil
}

// SaveApiInfo 保存API信息
func (m *merchantService) SaveApiInfo(mchId int64, d *merchant.ApiInfo) error {
	mch := m._mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		return mch.ApiManager().SaveApiInfo(d)
	}
	return merchant.ErrNoSuchMerchant
}

// GetMemberLevels_ 获取所有会员等级
func (m *merchantService) GetMemberLevels_(mchId int64) []*merchant.MemberLevel {
	mch := m._mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		return mch.LevelManager().GetLevelSet()
	}
	return []*merchant.MemberLevel{}
}

// GetMemberLevelById_ 根据编号获取会员等级信息
func (m *merchantService) GetMemberLevelById_(mchId, id int32) *merchant.MemberLevel {
	mch := m._mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		return mch.LevelManager().GetLevelById(id)
	}
	return nil
}

// SaveMemberLevel_ 保存会员等级信息
func (m *merchantService) SaveMemberLevel_(mchId int64, v *merchant.MemberLevel) (int32, error) {
	mch := m._mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		return mch.LevelManager().SaveLevel(v)
	}
	return 0, merchant.ErrNoSuchMerchant
}

// DelMemberLevel_ 删除会员等级
func (m *merchantService) DelMemberLevel_(mchId, levelId int32) error {
	mch := m._mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		return mch.LevelManager().DeleteLevel(levelId)
	}
	return merchant.ErrNoSuchMerchant
}

// GetLevel_ 获取等级
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

// GetKeyMapsByKeyword_ 获取键值字典
func (m *merchantService) GetKeyMapsByKeyword_(mchId int64, keyword string) map[string]string {
	mch := m._mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		return mch.KvManager().GetsByChar(keyword)
	}
	return map[string]string{}
}

// SaveKeyMaps_ 保存键值字典
func (m *merchantService) SaveKeyMaps_(mchId int64, data map[string]string) error {
	mch := m._mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		return mch.KvManager().Sets(data)
	}
	return merchant.ErrNoSuchMerchant
}

// WithdrawToMemberAccount1 提到会员账户
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
//		MerchantId:      mchId,
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
//		ProcedureFee:       0,
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
//		MerchantId:      machId,
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

// GetStaff implements proto.MerchantServiceServer.
func (m *merchantService) GetStaff(_ context.Context, req *proto.StaffRequest) (*proto.SStaff, error) {
	staff := m._staffRepo.Get(int(req.Id))
	if staff == nil {
		return &proto.SStaff{}, nil
	}
	if req.MchId != 0 && req.MchId != int64(staff.MchId) {
		// 如果商户不匹配，则返回空
		return &proto.SStaff{}, nil
	}
	m.checkImInitialized(staff)
	return m.parseStaffDto(staff), nil
}

// GetStaffByMember implements proto.MerchantServiceServer.
func (m *merchantService) GetStaffByMember(_ context.Context, req *proto.StaffRequest) (*proto.SStaff, error) {
	staff := m._staffRepo.GetStaffByMemberId(int(req.Id))
	if staff == nil {
		return &proto.SStaff{}, nil
	}
	if req.MchId != 0 && req.MchId != int64(staff.MchId) {
		// 如果商户不匹配，则返回空
		return &proto.SStaff{}, nil
	}
	m.checkImInitialized(staff)
	return m.parseStaffDto(staff), nil
}

// checkImInitialized 检查员工IM是否初始化
func (m *merchantService) checkImInitialized(s *staff.Staff) {
	if s.ImInitialized == 0 {
		// 发布员工IM初始化事件
		eventbus.Publish(&staff.StaffRequireImInitEvent{
			Staff: *s,
		})
	}
}

// SaveStaff implements proto.MerchantServiceServer.
func (m *merchantService) SaveStaff(_ context.Context, r *proto.SaveStaffRequest) (*proto.Result, error) {
	staff := m._staffRepo.Get(int(r.Id))
	if staff == nil {
		return m.error(errors.New("no such staff")), nil
	}
	staff.Flag = int(r.Flag)
	staff.Gender = int(r.Gender)
	staff.Grade = int(r.Grade)
	staff.WorkStatus = int(r.WorkStatus)
	staff.Nickname = r.Nickname
	_, err := m._staffRepo.Save(staff)
	return m.result(err), nil
}

func (m *merchantService) parseStaffDto(src *staff.Staff) *proto.SStaff {
	return &proto.SStaff{
		Id:            int64(src.Id),
		MemberId:      int64(src.MemberId),
		StationId:     int32(src.StationId),
		MchId:         int64(src.MchId),
		Flag:          int32(src.Flag),
		Gender:        int32(src.Gender),
		Nickname:      src.Nickname,
		WorkStatus:    int32(src.WorkStatus),
		Grade:         int32(src.Grade),
		Status:        int32(src.Status),
		IsCertified:   int32(src.IsCertified),
		CertifiedName: src.CertifiedName,
		PremiumLevel:  int32(src.PremiumLevel),
		CreateTime:    int64(src.CreateTime),
		ImInitialized: int32(src.ImInitialized),
	}
}

func (m *merchantService) parseMerchantDto(src *merchant.ComplexMerchant) *proto.QMerchant {
	return &proto.QMerchant{
		MchId:         int64(src.Id),
		Username:      src.Username,
		MchName:       src.Name,
		MemberId:      src.MemberId,
		MailAddr:      src.Username,
		IsSelf:        src.SelfSales,
		Flag:          int32(src.Flag),
		Level:         src.Level,
		Province:      src.Province,
		City:          src.City,
		District:      src.District,
		Address:       src.Address,
		Logo:          src.Logo,
		Telephone:     src.Telephone,
		Status:        int32(src.Status),
		ExpiresTime:   src.ExpiresTime,
		LastLoginTime: src.LastLoginTime,
		CreateTime:    0,
		Authenticate:  &proto.QAuthenticate{},
	}
}

func (m *merchantService) parseTradeConf(conf *proto.STradeConf_) *merchant.TradeConf {
	return &merchant.TradeConf{
		//MerchantId:       conf.MerchantId,
		//TradeType:   int(conf.TradeType),
		//PlanId:      conf.PlanId,
		//Flag:        int(conf.Flag),
		//AmountBasis: int(conf.AmountBasis),
		//ProcedureFee:    int(conf.ProcedureFee),
		//TradeRate:   int(conf.TradeRate),
	}
}

func (m *merchantService) parseTradeConfDto(conf *merchant.TradeConf) *proto.STradeConf_ {
	return &proto.STradeConf_{
		//MerchantId:       conf.MerchantId,
		//TradeType:   int32(conf.TradeType),
		//PlanId:      conf.PlanId,
		//Flag:        int32(conf.Flag),
		//AmountBasis: int32(conf.AmountBasis),
		//ProcedureFee:    int32(conf.ProcedureFee),
		//TradeRate:   int32(conf.TradeRate),
	}
}

func (m *merchantService) parseAccountDto(v *merchant.Account) *proto.SMerchantAccount {
	return &proto.SMerchantAccount{
		Balance:           int64(v.Balance),
		FreezeAmount:      int64(v.FreezeAmount),
		AwaitAmount:       int64(v.AwaitAmount),
		PresentAmount:     int64(v.PresentAmount),
		SalesAmount:       int64(v.SalesAmount),
		RefundAmount:      int64(v.RefundAmount),
		WithdrawalAmount:  int64(v.WithdrawalAmount),
		InvoiceableAmount: int64(v.InvoiceableAmount),
		OfflineSales:      int64(v.OfflineSales),
		UpdateTime:        int64(v.UpdateTime),
	}
}

func (m *merchantService) parseSaleConf(v *proto.SMerchantSaleConf) *merchant.SaleConf {
	return &merchant.SaleConf{
		MchId:           int(v.MerchantId),
		FxSales:         types.ElseInt(v.FxSalesEnabled, 1, 0),
		CbPercent:       float64(v.CashBackPercent),
		CbTg1Percent:    float64(v.CashBackTg1Percent),
		CbTg2Percent:    float64(v.CashBackTg2Percent),
		CbMemberPercent: float64(v.CashBackMemberPercent),
		OaOpen:          types.ElseInt(v.AutoSetupOrder, 1, 0),
		OaTimeoutMinute: int(v.OrderTimeOutMinute),
		OaConfirmMinute: int(v.OrderConfirmAfterMinute),
		OaReceiveHour:   int(v.OrderTimeOutReceiveHour),
	}
}

func (m *merchantService) parseSaleConfDto(v merchant.SaleConf) *proto.SMerchantSaleConf {
	return &proto.SMerchantSaleConf{
		MerchantId:              int64(v.MchId),
		FxSalesEnabled:          v.FxSales == 1,
		CashBackPercent:         float64(v.CbPercent),
		CashBackTg1Percent:      float64(v.CbTg1Percent),
		CashBackTg2Percent:      float64(v.CbTg2Percent),
		CashBackMemberPercent:   float64(v.CbMemberPercent),
		AutoSetupOrder:          v.OaOpen == 1,
		OrderTimeOutMinute:      int32(v.OaTimeoutMinute),
		OrderConfirmAfterMinute: int32(v.OaConfirmMinute),
		OrderTimeOutReceiveHour: int32(v.OaReceiveHour),
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

func (m *merchantService) parseGroupDto(v *merchant.MchBuyerGroupSetting) *proto.SMerchantBuyerGroup {
	return &proto.SMerchantBuyerGroup{
		Id:              int64(v.ID),
		GroupId:         v.GroupId,
		Name:            v.Alias,
		EnableRetail:    v.EnableRetail == 1,
		EnableWholesale: v.EnableWholesale == 1,
		RebatePeriod:    v.RebatePeriod,
	}
}

func (m *merchantService) parseGroup(v *proto.SMerchantBuyerGroup) *merchant.MchBuyerGroupSetting {
	return &merchant.MchBuyerGroupSetting{
		ID:              int32(v.Id),
		Alias:           v.Name,
		GroupId:         v.GroupId,
		EnableRetail:    int32(types.ElseInt(v.EnableRetail, 1, 0)),
		EnableWholesale: int32(types.ElseInt(v.EnableWholesale, 1, 0)),
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
		State:          int32(v.Status),
		StatusText:     order.OrderStatus(v.Status).String(),
		CreateTime:     v.CreateTime,
		Items:          items,
		Data:           v.Data,
	}
}

func (m *merchantService) parseBuyerGroupDto(v *merchant.BuyerGroup) *proto.SMerchantBuyerGroup {
	return &proto.SMerchantBuyerGroup{
		Id:              v.GroupId,
		Name:            v.Name,
		GroupId:         v.GroupId,
		EnableRetail:    v.EnableRetail,
		EnableWholesale: v.EnableWholesale,
		RebatePeriod:    int32(v.RebatePeriod),
	}
}

// AdjustAccount implements proto.MerchantServiceServer.
func (m *merchantService) AdjustAccount(_ context.Context, req *proto.UserWalletAdjustRequest) (*proto.TxResult, error) {
	mch := m._mchRepo.GetMerchant(int(req.UserId))
	title := "系统冲正"
	// 人工冲正带[KF]字样
	if req.ManualAdjust {
		title = "[KF]系统冲正"
	}
	acc := mch.Account()
	err := acc.Adjust(title, int(req.Value), req.TransactionRemark, req.RelateUser)
	return m.errorV2(err), nil
}

// CarryToAccount 商户账户入账
func (m *merchantService) CarryToAccount(_ context.Context, req *proto.MerchantAccountCarrayRequest) (*proto.TxResult, error) {
	mch := m._mchRepo.GetMerchant(int(req.UserId))
	if mch == nil {
		return m.errorV2(member.ErrNoSuchMember), nil
	}
	acc := mch.Account()
	if acc == nil {
		return m.errorV2(member.ErrNoSuchMember), nil
	}
	id, err := acc.Carry(merchant.CarryParams{
		Freeze:            req.Freeze,
		OuterTxNo:         req.OuterTransactionNo,
		Amount:            int(req.Amount),
		TransactionFee:    int(req.TransactionFee),
		RefundAmount:      0,
		TransactionTitle:  req.TransactionTitle,
		TransactionRemark: req.TransactionRemark,
		OuterTxUid:        int(req.OuterTxUid),
		BillAmountType:    int(req.BillAmountType),
	})
	if err != nil {
		return m.errorV2(err), nil
	}
	return m.txResult(id, nil), nil
}

// CompleteTransaction implements proto.MerchantServiceServer.
func (m *merchantService) CompleteTransaction(_ context.Context, req *proto.FinishUserTransactionRequest) (*proto.TxResult, error) {
	mch := m._mchRepo.GetMerchant(int(req.UserId))
	if mch == nil {
		return m.errorV2(merchant.ErrNoSuchMerchant), nil
	}
	err := mch.Account().CompleteTransaction(
		int(req.TransactionId),
		req.OuterTransactionNo)
	return m.errorV2(err), nil
}

// RequestWithdraw implements proto.MerchantServiceServer.
func (m *merchantService) RequestWithdrawal(_ context.Context, req *proto.UserWithdrawRequest) (*proto.TxResult, error) {
	mch := m._mchRepo.GetMerchant(int(req.UserId))
	if mch == nil {
		return m.errorV2(merchant.ErrNoSuchMerchant), nil
	}
	title := ""
	kind := 0
	switch int(req.WithdrawalKind) {
	case int(proto.EUserWithdrawalKind_WithdrawToBankCard):
		title = "提现到银行卡"
		kind = wallet.KWithdrawToBankCard
	case int(proto.EUserWithdrawalKind_WithdrawToPayWallet):
		title = "充值到第三方账户"
		kind = wallet.KWithdrawToPayWallet
	case int(proto.EUserWithdrawalKind_WithdrawByExchange):
		title = "提现到余额"
		kind = wallet.KWithdrawExchange
	case int(proto.EUserWithdrawalKind_WithdrawCustom):
		title = "自定义提现"
		kind = wallet.KWithdrawCustom
	}
	acc := mch.Account()
	transactionId, txNo, err := acc.RequestWithdrawal(
		&wallet.WithdrawTransaction{
			Amount:           int(req.Amount),
			TransactionFee:   int(req.GetTransactionFee()),
			Kind:             kind,
			TransactionTitle: title,
			BankName:         "",
			AccountNo:        req.AccountNo,
			AccountName:      "",
		})
	if err != nil {
		return m.errorV2(err), nil
	}
	return m.txResult(
		int(transactionId),
		map[string]string{
			"transationNo": txNo,
		}), nil
}

// ReviewWithdrawal implements proto.MerchantServiceServer.
func (m *merchantService) ReviewWithdrawal(_ context.Context, req *proto.ReviewUserWithdrawalRequest) (*proto.TxResult, error) {
	mch := m._mchRepo.GetMerchant(int(req.UserId))
	if mch == nil {
		return m.errorV2(merchant.ErrNoSuchMerchant), nil
	}
	err := mch.Account().ReviewWithdrawal(
		int(req.TransactionId),
		req.Pass,
		req.TransactionRemark)
	return m.errorV2(err), nil
}

// Freeze implements proto.MerchantServiceServer.
func (m *merchantService) Freeze(_ context.Context, req *proto.UserWalletFreezeRequest) (*proto.TxResult, error) {
	mch := m._mchRepo.GetMerchant(int(req.UserId))
	if mch == nil {
		return m.errorV2(merchant.ErrNoSuchMerchant), nil
	}
	id, err := mch.Account().Freeze(
		wallet.TransactionData{
			TransactionTitle:  req.TransactionTitle,
			Amount:            int(req.Amount),
			TransactionFee:    0,
			OuterTxNo:         req.OuterTransactionNo,
			TransactionRemark: req.TransactionRemark,
			TransactionId:     int(req.TransactionId),
			OuterTxUid:        0,
		}, 0)
	if err != nil {
		return m.errorV2(err), nil
	}
	return m.txResult(id, nil), nil
}

// Unfreeze implements proto.MerchantServiceServer.
func (m *merchantService) Unfreeze(_ context.Context, req *proto.UserWalletUnfreezeRequest) (*proto.TxResult, error) {
	mch := m._mchRepo.GetMerchant(int(req.UserId))
	if mch == nil {
		return m.errorV2(merchant.ErrNoSuchMerchant), nil
	}
	err := mch.Account().Unfreeze(
		wallet.TransactionData{
			TransactionTitle:  req.TransactionTitle,
			Amount:            int(req.Amount),
			TransactionFee:    0,
			OuterTxNo:         req.OuterTransactionNo,
			TransactionRemark: req.TransactionRemark,
			TransactionId:     int(req.TransactionId),
			OuterTxUid:        0,
		}, req.IsRefundBalance, 0)
	if err != nil {
		return m.errorV2(err), nil
	}
	return m.txResult(0, nil), nil
}

// GetWalletTxLog 获取会员钱包交易记录
func (m *merchantService) GetWalletTxLog(_ context.Context, r *proto.UserWalletTxId) (*proto.UserWalletTxResponse, error) {
	mch := m._mchRepo.GetMerchant(int(r.UserId))
	v := mch.Account().GetWalletLog(r.TxId)
	if v == nil {
		return nil, errors.New("交易不存在")
	}
	return &proto.UserWalletTxResponse{
		TxId:               int64(v.Id),
		UserId:             r.UserId,
		OuterTransactionNo: v.OuterTxNo,
		Kind:               int32(v.Kind),
		TransactionTitle:   v.Subject,
		Amount:             int64(v.ChangeValue),
		TransactionFee:     int64(v.TransactionFee),
		ReviewStatus:       int32(v.ReviewStatus),
		TransactionRemark:  v.Remark,
		CreateTime:         int64(v.CreateTime),
		UpdateTime:         int64(v.UpdateTime),
		RelateUser:         int64(v.OprUid),
	}, nil
}

// TransferStaff 转移员工
func (m *merchantService) TransferStaff(_ context.Context, req *proto.TransferStaffRequest) (*proto.TxResult, error) {
	mch := m._mchRepo.GetMerchant(int(req.MchId))
	if mch == nil {
		return m.errorV2(merchant.ErrNoSuchMerchant), nil
	}
	txId, err := mch.EmployeeManager().RequestTransfer(int(req.StaffId), int(req.TransferMchId))
	if err != nil {
		return m.errorV2(err), nil
	}
	return m.txResult(txId, nil), nil
}

// AdjustBillAmount implements proto.MerchantServiceServer.
func (m *merchantService) ManualAdjustBillAmount(_ context.Context, req *proto.ManualAdjustMerchantBillAmountRequest) (*proto.TxResult, error) {
	mch := m._mchRepo.GetMerchant(int(req.MchId))
	if mch == nil {
		return m.errorV2(merchant.ErrNoSuchMerchant), nil
	}
	err := mch.Account().Adjust(req.Title, int(req.Amount), req.Remark, req.OprUid)
	if err == nil {
		err = mch.SaleManager().AdjustBillAmount(
			merchant.BillAmountType(req.BillAmountType),
			int(req.Amount),
			int(req.TxFee))
	}
	if err != nil {
		return m.errorV2(err), nil
	}
	return m.txResult(0, nil), nil
}

// GenerateBill implements proto.MerchantServiceServer.
func (m *merchantService) GenerateBill(_ context.Context, req *proto.GenerateMerchantBillRequest) (*proto.TxResult, error) {
	mch := m._mchRepo.GetMerchant(int(req.MchId))
	if mch == nil {
		return m.errorV2(merchant.ErrNoSuchMerchant), nil
	}
	manager := mch.SaleManager()
	if req.BillId <= 0 {
		// 默认生成当前月份的账单
		bill := manager.GetCurrentBill()
		req.BillId = int64(bill.Id)
	}
	err := manager.GenerateBill(int(req.BillId))
	if err != nil {
		return m.errorV2(err), nil
	}
	return m.txResult(0, nil), nil
}

// GetBill implements proto.MerchantServiceServer.
func (m *merchantService) GetBill(_ context.Context, req *proto.BillTimeRequest) (*proto.SMerchantBill, error) {
	mch := m._mchRepo.GetMerchant(int(req.MchId))
	if mch == nil {
		return nil, errors.New("商户不存在")
	}
	manager := mch.SaleManager()
	bill := manager.GetBillByTime(int(req.BillTime))
	if bill == nil {
		return nil, errors.New("账单不存在")
	}
	return m.parseMerchantBill(bill), nil
}

func (m *merchantService) parseMerchantBill(bill *merchant.MerchantBill) *proto.SMerchantBill {
	return &proto.SMerchantBill{
		Id:               int64(bill.Id),
		MchId:            int64(bill.MchId),
		BillTime:         int64(bill.BillTime),
		BillMonth:        bill.BillMonth,
		StartTime:        int64(bill.StartTime),
		EndTime:          int64(bill.EndTime),
		ShopOrderCount:   int32(bill.ShopOrderCount),
		StoreOrderCount:  int32(bill.StoreOrderCount),
		ShopTotalAmount:  int64(bill.ShopTotalAmount),
		StoreTotalAmount: int64(bill.StoreTotalAmount),
		OtherOrderCount:  int32(bill.OtherOrderCount),
		OtherTotalAmount: int64(bill.OtherTotalAmount),
		TotalTxFee:       int64(bill.TotalTxFee),
		Status:           int32(bill.Status),
		ReviewerId:       int64(bill.ReviewerId),
		ReviewerName:     bill.ReviewerName,
		ReviewTime:       int64(bill.ReviewTime),
		CreateTime:       int64(bill.CreateTime),
		BuildTime:        int64(bill.BuildTime),
		UpdateTime:       int64(bill.UpdateTime),
	}
}

// ReviewBill implements proto.MerchantServiceServer.
func (m *merchantService) ReviewBill(_ context.Context, req *proto.ReviewMerchantBillRequest) (*proto.TxResult, error) {
	mch := m._mchRepo.GetMerchant(int(req.MchId))
	if mch == nil {
		return m.errorV2(merchant.ErrNoSuchMerchant), nil
	}
	manager := mch.SaleManager()
	err := manager.ReviewBill(int(req.BillId), int(req.ReviewerId))
	if err != nil {
		return m.errorV2(err), nil
	}
	return m.txResult(int(req.BillId), nil), nil
}

// SettleBill implements proto.MerchantServiceServer.
func (m *merchantService) SettleBill(_ context.Context, req *proto.SettleMerchantBillRequest) (*proto.TxResult, error) {
	mch := m._mchRepo.GetMerchant(int(req.MchId))
	if mch == nil {
		return m.errorV2(merchant.ErrNoSuchMerchant), nil
	}
	manager := mch.SaleManager()
	err := manager.SettleBill(int(req.BillId))
	if err != nil {
		return m.errorV2(err), nil
	}
	return m.txResult(int(req.BillId), nil), nil
}

// RequestInvoice 商户申请发票
func (m *merchantService) RequestInvoice(_ context.Context, req *proto.MerchantRequestInvoiceRequest) (*proto.TxResult, error) {
	mch := m._mchRepo.GetMerchant(int(req.MchId))
	if mch == nil {
		return m.errorV2(merchant.ErrNoSuchMerchant), nil
	}
	account := mch.Account()
	txId, err := account.RequestInvoice(int(req.Amount), req.Remark)
	if err != nil {
		return m.errorV2(err), nil
	}
	return m.txResult(txId, nil), nil
}
