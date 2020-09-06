package impl

/**
 * Copyright 2014 @ to2.net.
 * name :
 * author : jarryliu
 * date : 2020-09-05 20:14
 * description :
 * history :
 */
import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ixre/gof"
	"github.com/ixre/gof/math"
	"github.com/ixre/gof/storage"
	"github.com/ixre/gof/util"
	de "go2o/core/domain/interface/domain"
	"go2o/core/domain/interface/domain/enum"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/mss/notify"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/dto"
	"go2o/core/infrastructure/domain"
	"go2o/core/infrastructure/format"
	"go2o/core/module"
	"go2o/core/query"
	"go2o/core/service/proto"
	"go2o/core/variable"
	"strconv"
	"strings"
	"time"
)

var _ proto.MemberServiceServer = new(memberService)

type memberService struct {
	repo       member.IMemberRepo
	mchService *merchantService
	query      *query.MemberQuery
	orderQuery *query.OrderQuery
	valRepo    valueobject.IValueRepo
	serviceUtil
	sto storage.Interface
}

// 交换会员编号
func (s *memberService) SwapMemberId(_ context.Context, r *proto.SwapMemberRequest) (*proto.Int64, error) {
	var memberId int64
	switch r.Cred {
	default:
	case proto.ECredentials_User:
		memberId = s.repo.GetMemberIdByUser(r.Value)
	case proto.ECredentials_Code:
		memberId = int64(s.repo.GetMemberIdByCode(r.Value))
	case proto.ECredentials_Phone:
		memberId = s.repo.GetMemberIdByPhone(r.Value)
	case proto.ECredentials_Email:
		memberId = s.repo.GetMemberIdByEmail(r.Value)
	case proto.ECredentials_InviteCode:
		memberId = s.repo.GetMemberIdByInviteCode(r.Value)
	}
	return &proto.Int64{Value: memberId}, nil
}

func NewMemberService(mchService *merchantService, repo member.IMemberRepo,
	q *query.MemberQuery, oq *query.OrderQuery, valRepo valueobject.IValueRepo) *memberService {
	s := &memberService{
		repo:       repo,
		query:      q,
		mchService: mchService,
		orderQuery: oq,
		valRepo:    valRepo,
	}
	return s
	//return s.init()
}

func (s *memberService) init() *memberService {
	db := gof.CurrentApp.Db()
	var list []*member.Member
	db.GetOrm().Select(&list, "")
	for _, v := range list {
		im := s.repo.CreateMember(v)
		if rl := im.GetRelation(); rl != nil {
			im.BindInviter(rl.InviterId, true)
		}
		//if len(v.InviteCode) < 6 {
		//	im := s.repo.CreateMember(v)
		//	v.InviteCode = s.generateInviteCode()
		//	im.SetValue(v)
		//	im.Save()
		//}
	}
	return s
}

// 根据会员编号获取会员
func (s *memberService) getMemberValue(memberId int64) *member.Member {
	if memberId > 0 {
		v := s.repo.GetMember(memberId)
		if v != nil {
			vv := v.GetValue()
			return &vv
		}
	}
	return nil
}

// 根据会员编号获取会员
func (s *memberService) GetMember(_ context.Context, id *proto.Int64) (*proto.SMember, error) {
	v := s.getMemberValue(id.Value)
	if v != nil {
		return s.parseMemberDto(v), nil
	}
	return nil, nil
}

// 根据用户名获取会员
func (s *memberService) GetMemberByUser(_ context.Context, user *proto.String) (*proto.SMember, error) {
	v := s.repo.GetMemberByUser(user.Value)
	if v != nil {
		return s.parseMemberDto(v), nil
	}
	return nil, nil
}

// 获取资料
func (s *memberService) GetProfile(_ context.Context, i *proto.Int64) (*proto.SProfile, error) {
	m := s.repo.GetMember(i.Value)
	if m != nil {
		v := m.Profile().GetProfile()
		return s.parseMemberProfile(&v), nil
	}
	return nil, nil
}

// 保存资料
func (s *memberService) SaveProfile(v *proto.SProfile) error {
	if v.MemberId > 0 {
		v2 := s.parseMemberProfile2(v)
		m := s.repo.GetMember(v.MemberId)
		if m == nil {
			return member.ErrNoSuchMember
		}
		return m.Profile().SaveProfile(v2)
	}
	return nil
}

// 升级为高级会员
func (s *memberService) Premium(_ context.Context, r *proto.PremiumRequest) (*proto.Result, error) {
	m := s.repo.GetMember(r.MemberId)
	if m == nil {
		return s.result(member.ErrNoSuchMember), nil
	}
	err := m.Premium(int(r.V), r.Expires)
	return s.result(err), nil
}

// 检查会员的会话Token是否正确
func (s *memberService) CheckToken(_ context.Context, r *proto.CheckTokenRequest) (*proto.Bool, error) {
	md := module.Get(module.MM).(*module.MemberModule)
	return &proto.Bool{
		Value: md.CheckToken(r.MemberId, r.Token),
	}, nil
}

// 获取会员的会员Token,reset表示是否重置会员的token
func (s *memberService) GetToken(_ context.Context, r *proto.GetTokenRequest) (*proto.String, error) {
	pubToken := ""
	md := module.Get(module.MM).(*module.MemberModule)
	if !r.Reset_ {
		pubToken = md.GetToken(r.MemberId)
	}
	if r.Reset_ || (pubToken == "" && r.MemberId > 0) {
		m := s.getMemberValue(r.MemberId)
		if m != nil {
			return &proto.String{Value: md.ResetToken(r.MemberId, m.Pwd)}, nil
		}
	}
	return &proto.String{Value: pubToken}, nil
}

// 移除会员的Token
func (s *memberService) RemoveToken(_ context.Context, id *proto.Int64) (*proto.Empty, error) {
	md := module.Get(module.MM).(*module.MemberModule)
	md.RemoveToken(id.Value)
	return &proto.Empty{}, nil
}

// 更改手机号码，不验证手机格式
func (s *memberService) ChangePhone(_ context.Context, r *proto.ChangePhoneRequest) (*proto.Result, error) {
	err := s.changePhone(r.MemberId, r.Phone)
	return s.result(err), nil
}

// 更改邀请人
func (s *memberService) ChangeInviterId(_ context.Context, r *proto.ChangeInviterRequest) (*proto.Result, error) {
	im := s.repo.GetMember(r.MemberId)
	if im == nil {
		return s.result(member.ErrNoSuchMember), nil
	}
	err := im.BindInviter(r.InviterId, true)
	return s.result(err), nil
}

// 是否已收藏
func (s *memberService) Favored(memberId int64, favType int, referId int32) bool {
	return s.repo.CreateMemberById(memberId).
		Favorite().Favored(favType, referId)
}

// 取消收藏
func (s *memberService) Cancel(memberId int64, favType int, referId int32) error {
	return s.repo.CreateMemberById(memberId).
		Favorite().Cancel(favType, referId)
}

// 收藏商品
func (s *memberService) FavoriteGoods(memberId int64, goodsId int32) error {
	return s.repo.CreateMemberById(memberId).
		Favorite().Favorite(member.FavTypeGoods, goodsId)
}

// 取消收藏商品
func (s *memberService) CancelGoodsFavorite(memberId int64, goodsId int32) error {
	return s.repo.CreateMemberById(memberId).
		Favorite().Cancel(member.FavTypeGoods, goodsId)
}

// 收藏店铺
func (s *memberService) FavoriteShop(memberId int64, shopId int32) error {
	return s.repo.CreateMemberById(memberId).
		Favorite().Favorite(member.FavTypeShop, shopId)
}

// 取消收藏店铺
func (s *memberService) CancelShopFavorite(memberId int64, shopId int32) error {
	return s.repo.CreateMemberById(memberId).
		Favorite().Cancel(member.FavTypeShop, shopId)
}

// 商品是否已收藏
func (s *memberService) GoodsFavored(memberId int64, goodsId int32) bool {
	return s.repo.CreateMemberById(memberId).
		Favorite().Favored(member.FavTypeGoods, goodsId)
}

// 商店是否已收藏
func (s *memberService) ShopFavored(memberId int64, shopId int32) bool {
	return s.repo.CreateMemberById(memberId).
		Favorite().Favored(member.FavTypeShop, shopId)
}

// 获取会员的订单状态及其数量
func (s *memberService) OrdersQuantity(_ context.Context, id *proto.Int64) (*proto.OrderQuantityMapResponse, error) {
	ret := make(map[int32]int32, 0)
	for k, v := range s.query.OrdersQuantity(id.Value) {
		ret[int32(k)] = int32(v)
	}
	return &proto.OrderQuantityMapResponse{Data: ret}, nil
}

/**================ 会员等级 ==================**/
// 获取所有会员等级
func (s *memberService) GetMemberLevels() []*member.Level {
	return s.repo.GetManager().LevelManager().GetLevelSet()
}

// 等级列表
func (s *memberService) MemberLevelList(_ context.Context, empty *proto.Empty) (*proto.SMemberLevelListResponse, error) {
	var arr []*proto.SMemberLevel
	list := s.repo.GetManager().LevelManager().GetLevelSet()
	for _, v := range list {
		arr = append(arr, s.parseLevelDto(v))
	}
	return &proto.SMemberLevelListResponse{List: arr}, nil
}

// 根据编号获取会员等级信息
func (s *memberService) GetMemberLevel(_ context.Context, i *proto.Int32) (*proto.SMemberLevel, error) {
	lv := s.repo.GetManager().LevelManager().GetLevelById(int(i.Value))
	if lv != nil {
		return s.parseLevelDto(lv), nil
	}
	return nil, nil
}

// 保存会员等级信息
func (s *memberService) SaveMemberLevel(_ context.Context, level *proto.SMemberLevel) (*proto.Result, error) {
	lv := &member.Level{
		ID:            int(level.ID),
		Name:          level.Name,
		RequireExp:    int(level.RequireExp),
		ProgramSignal: level.ProgramSignal,
		Enabled:       int(level.Enabled),
		IsOfficial:    int(level.IsOfficial),
		AllowUpgrade:  int(level.AllowUpgrade),
	}
	_, err := s.repo.GetManager().LevelManager().SaveLevel(lv)
	if err != nil {
		return s.error(err), nil
	}
	return s.success(nil), nil
}

// 根据SIGN获取等级
func (s *memberService) GetLevelBySign(_ context.Context, sign *proto.String) (*proto.SMemberLevel, error) {
	lv := s.repo.GetManager().LevelManager().GetLevelByProgramSign(sign.Value)
	if lv != nil {
		return s.parseLevelDto(lv), nil
	}
	return nil, nil
}

// 根据可编程字符获取会员等级
func (s *memberService) GetLevelByProgramSign(sign string) *member.Level {
	return s.repo.GetManager().LevelManager().GetLevelByProgramSign(sign)
}

// 删除会员等级
func (s *memberService) DelMemberLevel(levelId int32) error {
	return s.repo.GetManager().LevelManager().DeleteLevel(int(levelId))
}

// 获取下一个等级
func (s *memberService) GetNextLevel(levelId int32) *member.Level {
	return s.repo.GetManager().LevelManager().GetNextLevelById(int(levelId))
}

// 获取启用中的最大等级,用于判断是否可以升级
func (s *memberService) GetHighestLevel() member.Level {
	lv := s.repo.GetManager().LevelManager().GetHighestLevel()
	if lv != nil {
		return *lv
	}
	return member.Level{}
}

func (s *memberService) GetWalletLog(memberId int64, logId int32) *member.WalletAccountLog {
	m := s.repo.GetMember(memberId)
	return m.GetAccount().GetWalletLog(logId)
}

func (s *memberService) getMember(memberId int64) (
	member.IMember, error) {
	if memberId <= 0 {
		return nil, member.ErrNoSuchMember
	}
	m := s.repo.GetMember(memberId)
	if m == nil {
		return m, member.ErrNoSuchMember
	}
	return m, nil
}

//todo: remove
func (s *memberService) GetMemberIdByInviteCode(code string) int64 {
	return s.repo.GetMemberIdByInviteCode(code)
}

// 根据信息获取会员编号
func (s *memberService) GetMemberIdByBasis(str string, basic int) int64 {
	switch basic {
	default:
	case notify.TypePhoneMessage:
		return s.repo.GetMemberIdByPhone(str)
	case notify.TypeEmailMessage:
		return s.repo.GetMemberIdByEmail(str)
	}
	return -1
}

// 发送会员验证码消息, 并返回验证码, 验证码通过data.code获取
func (s *memberService) SendCode(_ context.Context, r *proto.SendCodeRequest) (*proto.Result, error) {
	m := s.repo.GetMember(r.MemberId)
	if m == nil {
		return s.error(member.ErrNoSuchMember), nil
	}
	msgType, _ := strconv.Atoi(r.MsgType)
	code, err := m.SendCheckCode(r.Operation, msgType)
	if err != nil {
		return s.error(err), nil
	}
	return s.success(map[string]string{"code": code}), nil
}

// 比较验证码是否正确
func (s *memberService) CompareCode(_ context.Context, r *proto.CompareCodeRequest) (*proto.Result, error) {
	m := s.repo.GetMember(r.MemberId)
	if m == nil {
		return s.error(member.ErrNoSuchMember), nil
	}
	if err := m.CompareCode(r.Code); err != nil {
		return s.error(err), nil
	}
	return s.success(nil), nil
}

// 更改会员用户名
func (s *memberService) ChangeUser(_ context.Context, r *proto.ChangeUserRequest) (*proto.Result, error) {
	var err error
	m := s.repo.GetMember(int64(r.MemberId))
	if m == nil {
		err = member.ErrNoSuchMember
	} else {
		if err = m.ChangeUser(r.Usr); err == nil {
			return s.success(nil), nil
		}
	}
	return s.result(err), nil
}

// 获取会员等级信息
func (s *memberService) MemberLevelInfo(_ context.Context, id *proto.Int64) (*proto.SMemberLevelInfo, error) {
	level := &proto.SMemberLevelInfo{Level: -1}
	im := s.repo.GetMember(id.Value)
	if im != nil {
		v := im.GetValue()
		level.Exp = int32(v.Exp)
		level.Level = int32(v.Level)
		lv := im.GetLevel()
		level.LevelName = lv.Name
		level.ProgramSignal = lv.ProgramSignal
		nextLv := s.repo.GetManager().LevelManager().GetNextLevelById(lv.ID)
		if nextLv == nil {
			level.NextLevel = -1
		} else {
			level.NextLevel = int32(nextLv.ID)
			level.NextLevelName = nextLv.Name
			level.NextProgramSignal = nextLv.ProgramSignal
			level.RequireExp = int32(nextLv.RequireExp - v.Exp)
		}
	}
	return level, nil
}

// 更改会员等级
func (s *memberService) UpdateLevel(_ context.Context, r *proto.UpdateLevelRequest) (*proto.Result, error) {
	m := s.repo.GetMember(r.MemberId)
	var err error
	if m == nil {
		err = member.ErrNoSuchMember
	} else {
		err = m.ChangeLevel(int(r.Level), int(r.PaymentOrderId), r.Review)
	}
	return s.result(err), nil
}

// 上传会员头像
func (s *memberService) ChangeAvatar(memberId int64, avatar string) error {
	m := s.repo.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.Profile().ChangeAvatar(avatar)
}

// 保存用户
func (s *memberService) SaveMember(v *proto.SMember) (int64, error) {
	if v.Id > 0 {
		return s.updateMember(v)
	}
	return -1, errors.New("Create member use \"RegisterMember\" method.")
}

func (s *memberService) updateMember(v *proto.SMember) (int64, error) {
	m := s.repo.GetMember(v.Id)
	if m == nil {
		return -1, member.ErrNoSuchMember
	}
	mv := s.parseMember(v)
	if err := m.SetValue(mv); err != nil {
		return m.GetAggregateRootId(), err
	}
	return m.Save()
}

// 注册会员
func (s *memberService) RegisterMemberV2(_ context.Context, r *proto.RegisterMemberRequest) (*proto.Result, error) {
	if len(r.Pwd) != 32 {
		return s.error(de.ErrNotMD5Format), nil
	}
	v := &member.Member{
		User:     r.User,
		Pwd:      domain.Sha1Pwd(r.Pwd),
		Name:     r.Name,
		RealName: "",
		Avatar:   r.Avatar,
		Phone:    r.Phone,
		Email:    r.Email,
		RegFrom:  r.Extend["reg_from"],
		RegIp:    r.Extend["reg_ip"],
		Flag:     int(r.Flag),
	}
	// 验证邀请码
	inviteCode := r.Extend["invite_code"]
	inviterId, err := s.repo.GetManager().CheckInviteRegister(inviteCode)
	if err != nil {
		return s.error(err), nil
	}
	// 创建会员
	m := s.repo.CreateMember(v)
	id, err := m.Save()

	if err == nil {
		// 保存关联信息
		err = m.BindInviter(inviterId, true)
		//m := s.repo.CreateMember(v) //创建会员
		//id, err := m.Save()
		//if err == nil {
		//	pro.Sex = 1
		//	pro.MemberId = id
		//	//todo: 如果注册失败，则删除。应使用SQL-TRANSFER
		//	if err = m.Profile().SaveProfile(pro); err != nil {
		//		s.repo.DeleteMember(id)
		//}
		return s.success(map[string]string{
			"member_id": util.Str(id),
		}), nil
	}
	return s.error(err), nil
}

func (s *memberService) GetRelation(memberId int64) *member.InviteRelation {
	return s.repo.GetRelation(memberId)
}

// 激活会员
func (s *memberService) Active(_ context.Context, i *proto.Int64) (*proto.Result, error) {
	m := s.repo.GetMember(i.Value)
	if m == nil {
		return s.error(member.ErrNoSuchMember), nil
	}
	if err := m.Active(); err != nil {
		return s.error(err), nil
	}
	return s.success(nil), nil
}

// 锁定/解锁会员
func (s *memberService) Lock(_ context.Context, r *proto.LockRequest) (*proto.Result, error) {
	m := s.repo.GetMember(r.MemberId)
	if m == nil {
		return s.error(member.ErrNoSuchMember), nil
	}
	if err := m.Lock(int(r.Minutes), r.Remark); err != nil {
		return s.error(err), nil
	}
	return s.success(nil), nil
}

// 解锁会员
func (s *memberService) Unlock(_ context.Context, i *proto.Int64) (*proto.Result, error) {
	m := s.repo.GetMember(i.Value)
	if m == nil {
		return s.error(member.ErrNoSuchMember), nil
	}
	if err := m.Unlock(); err != nil {
		return s.error(err), nil
	}
	return s.success(nil), nil
}

// 判断资料是否完善
func (s *memberService) ProfileCompleted(memberId int64) bool {
	m := s.repo.GetMember(memberId)
	if m != nil {
		return m.Profile().ProfileCompleted()
	}
	return false
}

// 判断资料是否完善
func (s *memberService) CheckProfileComplete(_ context.Context, id *proto.Int64) (*proto.Result, error) {
	m := s.repo.GetMember(id.Value)
	var err error
	if m == nil {
		err = member.ErrNoSuchMember
	} else {
		err = m.Profile().CheckProfileComplete()
		if err != nil {
			switch err.Error() {
			case "phone":
				err = errors.New("未完善手机")
			case "birthday":
				err = errors.New("未完善生日")
			case "address":
				err = errors.New("未完善地址")
			case "im":
				err = errors.New("未完善" + variable.AliasMemberIM)
			}
		}
	}
	return s.result(err), nil
}

// 更改密码
func (s *memberService) ModifyPwd(_ context.Context, r *proto.ModifyPwdRequest) (*proto.Result, error) {
	m := s.repo.GetMember(r.MemberId)
	if m == nil {
		return s.error(member.ErrNoSuchMember), nil
	}
	pwd := r.Pwd
	old := r.Old
	if l := len(r.Pwd); l != 32 {
		return s.error(de.ErrNotMD5Format), nil
	} else {
		pwd = domain.MemberSha1Pwd(pwd)
	}
	if l := len(old); l > 0 && l != 32 {
		return s.error(de.ErrNotMD5Format), nil
	} else {
		old = domain.MemberSha1Pwd(old)
	}
	err := m.Profile().ModifyPassword(pwd, old)
	if err != nil {
		return s.error(err), nil
	}
	return s.success(nil), nil
}

// 更改交易密码
func (s *memberService) ModifyTradePwd(_ context.Context, r *proto.ModifyPwdRequest) (*proto.Result, error) {
	m := s.repo.GetMember(r.MemberId)
	if m == nil {
		return s.error(member.ErrNoSuchMember), nil
	}
	pwd, old := r.Pwd, r.Old
	if l := len(pwd); l != 32 {
		return s.error(de.ErrNotMD5Format), nil
	} else {
		pwd = domain.TradePwd(pwd)
	}
	if l := len(old); l > 0 && l != 32 {
		return s.error(de.ErrNotMD5Format), nil
	} else {
		old = domain.TradePwd(old)
	}
	err := m.Profile().ModifyTradePassword(pwd, old)
	if err != nil {
		return s.error(err), nil
	}
	return s.success(nil), nil
}

// 登录，返回结果(Result_)和会员编号(ID);
// Result值为：-1:会员不存在; -2:账号密码不正确; -3:账号被停用
func (s *memberService) testLogin(user string, pwd string) (id int64, errCode int32, err error) {
	user = strings.ToLower(user)
	memberId := s.repo.GetMemberIdByUser(user)
	if len(pwd) != 32 {
		return -1, 4, de.ErrNotMD5Format
	}
	if memberId <= 0 {
		//todo: 界面加上使用手机号码登陆
		//val = m.repo.GetMemberValueByPhone(user)
		return 0, 2, member.ErrNoSuchMember
	}
	val := s.repo.GetMember(memberId).GetValue()
	if val.Pwd != domain.Sha1Pwd(pwd) {
		return 0, 1, de.ErrCredential
	}
	if val.Flag&member.FlagLocked == member.FlagLocked {
		return 0, 3, member.ErrMemberLocked
	}

	return val.Id, 0, nil
}

// 登录，返回结果(Result_)和会员编号(ID);
// Result值为：-1:会员不存在; -2:账号密码不正确; -3:账号被停用
func (s *memberService) CheckLogin(_ context.Context, r *proto.LoginRequest) (*proto.Result, error) {
	id, code, err := s.testLogin(r.User, r.Pwd)
	if err != nil {
		r := s.error(err)
		r.ErrCode = code
		return r, nil
	}
	var memberCode = ""
	if r.Update {
		m := s.repo.GetMember(id)
		memberCode = m.GetValue().Code
		go m.UpdateLoginTime()
	}
	mp := map[string]string{
		"id":   strconv.Itoa(int(id)),
		"code": memberCode,
	}
	return s.success(mp), nil
}

// 检查交易密码
func (s *memberService) CheckTradePwd(_ context.Context, r *proto.CheckTradePwdRequest) (*proto.Result, error) {
	m := s.repo.GetMember(r.MemberId)
	if m == nil {
		return s.result(member.ErrNoSuchMember), nil
	}
	mv := m.GetValue()
	if mv.TradePwd == "" {
		return s.error(member.ErrNotSetTradePwd), nil
	}
	if len(r.TradePwd) != 32 {
		return s.error(de.ErrNotMD5Format), nil
	}
	if encPwd := domain.TradePwd(r.TradePwd); mv.TradePwd != encPwd {
		return s.error(member.ErrIncorrectTradePwd), nil
	}
	return s.success(nil), nil
}

// 检查与现有用户不同的用户是否存在,如存在则返回错误
//func (s *memberService) CheckUser(user string, memberId int64) error {
//	if len(user) < 6 {
//		return member.ErrUserLength
//	}
//	if s.repo.CheckUserExist(user, memberId) {
//		return member.ErrUserExist
//	}
//	return nil
//}

// 获取会员账户
func (s *memberService) GetAccount(_ context.Context, id *proto.Int64) (*proto.SAccount, error) {
	m := s.repo.CreateMember(&member.Member{Id: id.Value})
	acc := m.GetAccount()
	if acc != nil {
		return s.parseAccountDto(acc.GetValue()), nil
	}
	return nil, nil
}

// 获取上级邀请人会员编号数组
func (s *memberService) InviterArray(_ context.Context, r *proto.DepthRequest) (*proto.InviterIdListResponse, error) {
	m := s.repo.CreateMember(&member.Member{Id: r.MemberId})
	var arr []int64
	if m != nil {
		arr = m.Invitation().InviterArray(r.MemberId, int(r.Depth))
	}
	return &proto.InviterIdListResponse{
		List: arr,
	}, nil
}

// 按条件获取荐指定等级会员的数量
func (s *memberService) InviteMembersQuantity(_ context.Context, r *proto.DepthRequest) (*proto.Int32, error) {
	where := ""
	memberId := r.MemberId
	switch r.Depth {
	case 1:
		where = fmt.Sprintf(" inviter_id = %d", memberId)
	case 2:
		where = fmt.Sprintf(" inviter_id = %d OR inviter_d2 = %d", memberId, memberId)
	case 3:
		where = fmt.Sprintf(" inviter_id = %d OR inviter_d2 = %d OR inviter_d3 = %d", memberId, memberId, memberId)
	}
	if len(where) == 0 {
		return &proto.Int32{Value: 0}, nil
	}
	q := s.query.InviteMembersQuantity(memberId, where)
	return &proto.Int32{Value: int32(q)}, nil
}

// 按条件获取荐指定等级会员的数量
func (s *memberService) QueryInviteQuantity(_ context.Context, r *proto.InviteQuantityRequest) (*proto.Int64, error) {
	where := ""
	if r.Data != nil && len(r.Data) > 0 {
		where = s.parseGetInviterDataParams(r.Data)
	}
	return &proto.Int64{
		Value: int64(s.query.GetInviteQuantity(r.MemberId, where)),
	}, nil
}

// 按条件获取荐指定等级会员的列表
func (s *memberService) QueryInviteArray(_ context.Context, r *proto.InviteQuantityRequest) (*proto.MemberIdListResponse, error) {
	where := ""
	if r.Data != nil && len(r.Data) > 0 {
		where = s.parseGetInviterDataParams(r.Data)
	}
	return &proto.MemberIdListResponse{
		List: s.query.GetInviteArray(r.MemberId, where),
	}, nil
}

func (s *memberService) parseGetInviterDataParams(data map[string]string) string {
	buf := bytes.NewBufferString("")
	begin := data["begin"]
	end := data["end"]
	level := data["level"]
	operate := data["operate"]
	trust := data["trust"]

	if begin != "" && end != "" {
		buf.WriteString(" AND reg_time BETWEEN ")
		buf.WriteString(begin)
		buf.WriteString(" AND ")
		buf.WriteString(end)
	} else if begin != "" {
		buf.WriteString(" AND reg_time >= ")
		buf.WriteString(begin)
	} else if end != "" {
		buf.WriteString(" AND reg_time <= ")
		buf.WriteString(end)
	}

	if level != "" {
		if operate == "" {
			operate = ">="
		}
		buf.WriteString(" AND level ")
		buf.WriteString(operate)
		buf.WriteString(level)
	}

	if trust != "" {
		buf.WriteString(" AND review_state ")
		if trust == "true" {
			buf.WriteString(" = ")
		} else {
			buf.WriteString(" <> ")
		}
		trustOk := strconv.Itoa(int(enum.ReviewPass))
		buf.WriteString(trustOk)
	}
	return buf.String()
}

func (s *memberService) GetBank(memberId int64) *member.BankInfo {
	m := s.repo.CreateMember(&member.Member{Id: memberId})
	b := m.Profile().GetBank()
	return &b
}

func (s *memberService) SaveBankInfo(v *member.BankInfo) error {
	m := s.repo.CreateMember(&member.Member{Id: v.MemberId})
	return m.Profile().SaveBank(v)
}

// 解锁银行卡信息
func (s *memberService) UnlockBankInfo(memberId int64) error {
	m := s.repo.CreateMember(&member.Member{Id: memberId})
	return m.Profile().UnlockBank()
}

// 获取收款码
func (s *memberService) ReceiptsCodes(_ context.Context, id *proto.Int64) (*proto.SReceiptsCodeListResponse, error) {
	m := s.repo.GetMember(id.Value)
	if m == nil {
		return &proto.SReceiptsCodeListResponse{
			List: make([]*proto.SReceiptsCode, 0),
		}, nil
	}
	arr := m.Profile().ReceiptsCodes()
	list := make([]*proto.SReceiptsCode, len(arr))
	for i, v := range arr {
		list[i] = &proto.SReceiptsCode{
			Identity:  v.Identity,
			Name:      v.Name,
			AccountId: v.AccountId,
			CodeUrl:   v.CodeUrl,
			State:     int32(v.State),
		}
	}
	return &proto.SReceiptsCodeListResponse{List: list}, nil
}

// 保存收款码
func (s *memberService) SaveReceiptsCode(_ context.Context, r *proto.ReceiptsCodeSaveRequest) (*proto.Result, error) {
	m := s.repo.GetMember(r.MemberId)
	if m == nil {
		return s.error(member.ErrNoSuchMember), nil
	}
	v := &member.ReceiptsCode{
		Identity:  r.Code.Identity,
		Name:      r.Code.Name,
		AccountId: r.Code.AccountId,
		CodeUrl:   r.Code.CodeUrl,
		State:     int(r.Code.State),
	}
	if err := m.Profile().SaveReceiptsCode(v); err != nil {
		return s.error(err), nil
	}
	return s.success(nil), nil
}

// 获取银行卡
func (s *memberService) BankCards(_ context.Context, id *proto.Int64) (*proto.BankCardListResponse, error) {
	m := s.repo.CreateMember(&member.Member{Id: id.Value})
	b := m.Profile().GetBank()
	arr := make([]*proto.SBankcard, 0)
	arr = append(arr, &proto.SBankcard{
		Id:          -1,
		BankName:    b.BankName,
		AccountName: b.AccountName,
		Account:     b.Account,
		Network:     b.Network,
		State:       int32(b.State),
	})
	return &proto.BankCardListResponse{
		List: arr,
	}, nil
}

// 保存银行卡
func (s *memberService) SaveBankcard(_ context.Context, r *proto.BankCardSaveRequest) (*proto.Result, error) {
	m := s.repo.CreateMember(&member.Member{Id: r.MemberId})
	var v = &member.BankInfo{
		BankName:    r.Card.BankName,
		AccountName: r.Card.AccountName,
		Account:     r.Card.Account,
		Network:     r.Card.Network,
		State:       int(r.Card.State),
	}
	if err := m.Profile().SaveBank(v); err != nil {
		return s.error(err), nil
	}
	return s.success(nil), nil
}

// 实名认证信息
func (s *memberService) GetTrustInfo(_ context.Context, i *proto.Int64) (*proto.STrustedInfo, error) {
	t := member.TrustedInfo{}
	m := s.repo.GetMember(i.Value)
	if m != nil {
		t = m.Profile().GetTrustedInfo()
	}
	return &proto.STrustedInfo{
		RealName:         t.RealName,
		CountryCode:      t.CountryCode,
		CardType:         int32(t.CardType),
		CardId:           t.CardId,
		CardImage:        t.CardImage,
		CardReverseImage: t.CardReverseImage,
		TrustImage:       t.TrustImage,
		ManualReview:     int32(t.ManualReview),
		ReviewState:      int32(t.ReviewState),
		ReviewTime:       t.ReviewTime,
		Remark:           t.Remark,
	}, nil
}

// 保存实名认证信息
func (s *memberService) SubmitTrustInfo(_ context.Context, r *proto.SubmitTrustInfoRequest) (result *proto.Result, err error) {
	m := s.repo.GetMember(r.MemberId)
	if m == nil {
		err = member.ErrNoSuchMember
	} else {
		err = m.Profile().SaveTrustedInfo(&member.TrustedInfo{
			MemberId:         r.MemberId,
			RealName:         r.Info.RealName,
			CountryCode:      r.Info.CountryCode,
			CardType:         int(r.Info.CardType),
			CardId:           r.Info.CardId,
			CardImage:        r.Info.CardImage,
			CardReverseImage: r.Info.CardReverseImage,
			TrustImage:       r.Info.TrustImage,
		})
	}
	if err != nil {
		return s.error(err), nil
	}
	return s.success(nil), nil
}

// 审核实名认证,若重复审核将返回错误
func (s *memberService) ReviewTrustedInfo(_ context.Context, r *proto.ReviewTrustInfoRequest) (*proto.Result, error) {
	m := s.repo.GetMember(r.MemberId)
	err := m.Profile().ReviewTrustedInfo(r.ReviewPass, r.Remark)
	if err != nil {
		return s.error(err), nil
	}
	return s.success(nil), nil
}

// 获取分页店铺收藏
func (s *memberService) PagedShopFav(memberId int64, begin, end int,
	where string) (int, []*dto.PagedShopFav) {
	return s.query.PagedShopFav(memberId, begin, end, where)
}

// 获取分页店铺收藏
func (s *memberService) PagedGoodsFav(memberId int64, begin, end int,
	where string) (int, []*dto.PagedGoodsFav) {
	return s.query.PagedGoodsFav(memberId, begin, end, where)
}

// 获取钱包账户分页记录
func (s *memberService) PagingAccountLog(_ context.Context, r *proto.PagingAccountInfoRequest) (*proto.SPagingResult, error) {
	var total int
	var rows []map[string]interface{}
	switch r.AccountType {
	case member.AccountIntegral:
		total, rows = s.query.PagedIntegralAccountLog(r.MemberId, r.Params.Begin, r.Params.Over, r.Params.SortBy)
	case member.AccountBalance:
		total, rows = s.query.PagedBalanceAccountLog(r.MemberId, int(r.Params.Begin), int(r.Params.Over), "", "")
	case member.AccountWallet:
		total, rows = s.query.PagedWalletAccountLog(r.MemberId, int(r.Params.Begin), int(r.Params.Over), "", "")
	}
	rs := &proto.SPagingResult{
		ErrCode: 0,
		ErrMsg:  "",
		Count:   int32(total),
		Data:    s.json(rows),
	}
	return rs, nil
}

// 获取余额账户分页记录
func (s *memberService) PagedBalanceAccountLog(memberId int64, begin, end int,
	where, orderBy string) (int, []map[string]interface{}) {
	return s.query.PagedBalanceAccountLog(memberId, begin, end, where, orderBy)
}

// 获取钱包账户分页记录
func (s *memberService) PagedWalletAccountLog(memberId int64, begin, end int,
	where, orderBy string) (int, []map[string]interface{}) {
	return s.query.PagedWalletAccountLog(memberId, begin, end, where, orderBy)
}

// 查询分页普通订单
func (s *memberService) QueryNormalOrder(memberId int64, begin, size int, pagination bool,
	where, orderBy string) (num int, rows []*dto.PagedMemberSubOrder) {
	return s.orderQuery.QueryPagerOrder(memberId, begin, size, pagination, where, orderBy)
}

// 查询分页批发订单
func (s *memberService) QueryWholesaleOrder(memberId int64, begin, size int, pagination bool,
	where, orderBy string) (num int, rows []*dto.PagedMemberSubOrder) {
	return s.orderQuery.PagedWholesaleOrderOfBuyer(memberId, begin, size, pagination, where, orderBy)
}

// 查询分页订单
func (s *memberService) PagedTradeOrder(buyerId int64, begin, size int, pagination bool,
	where, orderBy string) (int, []*proto.SComplexOrder) {
	return s.orderQuery.PagedTradeOrderOfBuyer(buyerId, begin, size, pagination, where, orderBy)
}

/*********** 收货地址 ***********/

// 获取会员的收货地址
func (s *memberService) GetAddressList(_ context.Context, id *proto.Int64) (*proto.AddressListResponse, error) {
	src := s.repo.GetDeliverAddress(id.Value)
	var arr []*proto.SAddress
	for _, v := range src {
		arr = append(arr, s.parseAddressDto(v))
	}
	return &proto.AddressListResponse{List: arr}, nil
}

//获取配送地址
func (s *memberService) GetAddress(_ context.Context, r *proto.GetAddressRequest) (*proto.SAddress, error) {
	m := s.repo.CreateMember(&member.Member{Id: r.MemberId})
	pro := m.Profile()
	var addr member.IDeliverAddress
	if r.AddrId > 0 {
		addr = pro.GetAddress(r.AddrId)
	} else {
		addr = pro.GetDefaultAddress()
	}
	if addr != nil {
		v := addr.GetValue()
		d := s.parseAddressDto(&v)
		d.Area = s.valRepo.GetAreaString(
			v.Province, v.City, v.District)
		return d, nil
	}
	return nil, nil
}

//保存配送地址
func (s *memberService) SaveAddress(memberId int64, src *proto.SAddress) (int64, error) {
	e := s.parseAddress(src)
	m := s.repo.CreateMember(&member.Member{Id: memberId})
	var v member.IDeliverAddress
	var err error
	if e.ID > 0 {
		v = m.Profile().GetAddress(e.ID)
		err = v.SetValue(e)
	} else {
		v = m.Profile().CreateDeliver(e)
	}
	if err != nil {
		return -1, err
	}
	return v.Save()
}

//删除配送地址
func (s *memberService) DeleteAddress(memberId int64, deliverId int64) error {
	m := s.repo.CreateMember(&member.Member{Id: memberId})
	return m.Profile().DeleteAddress(deliverId)
}

//设置余额优先支付
func (s *memberService) BalancePriorityPay(memberId int64, enabled bool) error {
	m := s.repo.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().SetPriorityPay(member.AccountBalance, enabled)
}

//判断会员是否由指定会员邀请推荐的
func (s *memberService) IsInvitation(memberId int64, invitationMemberId int64) bool {
	m := s.repo.CreateMember(&member.Member{Id: memberId})
	return m.Invitation().InvitationBy(invitationMemberId)
}

// 获取我邀请的会员及会员邀请的人数
func (s *memberService) GetMyPagedInvitationMembers(memberId int64,
	begin, end int) (total int, rows []*dto.InvitationMember) {
	iv := s.repo.CreateMember(&member.Member{Id: memberId}).Invitation()
	total, rows = iv.GetInvitationMembers(begin, end)
	if l := len(rows); l > 0 {
		arr := make([]int32, l)
		for i := 0; i < l; i++ {
			arr[i] = rows[i].MemberId
		}
		num := iv.GetSubInvitationNum(arr)
		for i := 0; i < l; i++ {
			rows[i].InvitationNum = num[rows[i].MemberId]
			rows[i].Avatar = format.GetResUrl(rows[i].Avatar)
		}
	}
	return total, rows
}

// 获取会员最后更新时间
func (s *memberService) GetMemberLatestUpdateTime(memberId int64) int64 {
	return s.repo.GetMemberLatestUpdateTime(memberId)
}

func (s *memberService) GetMemberList(ids []int64) []*dto.MemberSummary {
	list := s.query.GetMemberList(ids)
	for _, v := range list {
		v.Avatar = format.GetResUrl(v.Avatar)
	}
	return list
}

// 标志赋值, 如果flag小于零, 则异或运算
func (s *memberService) GrantFlag(_ context.Context, r *proto.GrantFlagRequest) (*proto.Result, error) {
	m := s.repo.GetMember(r.MemberId)
	if m == nil {
		return s.error(member.ErrNoSuchMember), nil
	}
	if err := m.GrantFlag(int(r.Flag)); err != nil {
		return s.error(err), nil
	}
	return s.success(nil), nil
}

// 获取会员汇总信息
func (s *memberService) Complex(_ context.Context, id *proto.Int64) (*proto.SComplexMember, error) {
	m := s.repo.GetMember(id.Value)
	if m != nil {
		x := m.Complex()
		return s.parseComplexMemberDto(x), nil
	}
	return nil, nil
}

// 冻结积分,当new为true不扣除积分,反之扣除积分
func (s *memberService) FreezesIntegral(memberId int64, title string, value int64,
	new bool) error {
	m := s.repo.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().FreezesIntegral(title, int(value), new, 0)
}

// 解冻积分
func (s *memberService) UnfreezesIntegral(memberId int64, title string, value int64) error {
	m := s.repo.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().UnfreezesIntegral(title, int(value))
}

// 充值,account为账户类型,kind为业务类型
func (s *memberService) AccountCharge(_ context.Context, r *proto.AccountChangeRequest) (*proto.Result, error) {
	var err error
	m := s.repo.CreateMember(&member.Member{Id: r.MemberId})
	acc := m.GetAccount()
	if acc == nil {
		err = member.ErrNoSuchMember
	} else {
		err = acc.Charge(r.Account, r.Title, int(r.Amount), r.OuterNo, r.Remark)
	}
	return s.result(err), nil
}

// 账户抵扣
func (s *memberService) AccountDiscount(_ context.Context, r *proto.AccountChangeRequest) (*proto.Result, error) {
	m, err := s.getMember(r.MemberId)
	if err == nil {
		acc := m.GetAccount()
		err = acc.Discount(int(r.Account), r.Title, int(r.Amount), r.OuterNo, r.Remark)
	}
	return s.result(err), nil
}

// 账户消耗
func (s *memberService) AccountConsume(_ context.Context, r *proto.AccountChangeRequest) (*proto.Result, error) {
	m, err := s.getMember(r.MemberId)
	if err == nil {
		acc := m.GetAccount()
		err = acc.Consume(int(r.Account), r.Title, int(r.Amount), r.OuterNo, r.Remark)
	}
	return s.result(err), nil
}

// 账户消耗
func (s *memberService) AccountRefund(_ context.Context, r *proto.AccountChangeRequest) (*proto.Result, error) {
	m, err := s.getMember(r.MemberId)
	if err == nil {
		acc := m.GetAccount()
		err = acc.Refund(int(r.Account), r.Title, int(r.Account), r.OuterNo, r.Remark)
	}
	return s.result(err), nil
}

// 调整账户
func (s *memberService) AccountAdjust(_ context.Context, r *proto.AccountAdjustRequest) (*proto.Result, error) {
	m, err := s.getMember(r.MemberId)
	if err == nil {
		tit := "[KF]系统冲正"
		if r.Value > 0 {
			tit = "[KF]系统充值"
		}
		acc := m.GetAccount()
		err = acc.Adjust(int(r.Account), tit, int(r.Value), r.Remark, r.RelateUser)
	}
	return s.result(err), nil
}

// !银行四要素认证
func (s *memberService) B4EAuth(_ context.Context, r *proto.B4EAuthRequest) (*proto.Result, error) {
	mod := module.Get(module.B4E).(*module.Bank4E)
	if r.Action == "get" {
		data := mod.GetBasicInfo(r.MemberId)
		d, err := json.Marshal(data)
		if err != nil {
			return s.error(err), nil
		}
		return s.success(map[string]string{"data": string(d)}), nil
	}
	if r.Action == "update" {
		err := mod.UpdateInfo(r.MemberId,
			r.Data["real_name"],
			r.Data["id_card"],
			r.Data["phone"],
			r.Data["bank_account"])
		return s.result(err), nil
	}
	return s.error(errors.New("未知操作")), nil
}

// 验证交易密码
func (s *memberService) VerifyTradePwd(memberId int64, tradePwd string) (bool, error) {
	im, err := s.getMember(memberId)
	if err == nil {
		m := im.GetValue()
		if len(m.TradePwd) == 0 {
			return false, member.ErrNotSetTradePwd
		}
		if m.TradePwd != tradePwd {
			return false, member.ErrIncorrectTradePwd
		}
		return true, err
	}
	return false, err
}

// 提现并返回提现编号,交易号以及错误信息
func (s *memberService) SubmitTakeOutRequest(memberId int64, takeKind int32,
	applyAmount float32, commission float32) (int32, string, error) {
	m, err := s.getMember(memberId)
	if err != nil {
		return 0, "", err
	}

	acc := m.GetAccount()
	var title string
	switch int(takeKind) {
	case member.KindWalletTakeOutToBankCard:
		title = "提现到银行卡"
	case member.KindWalletTakeOutToBalance:
		title = "充值账户"
	case member.KindWalletTakeOutToThirdPart:
		title = "充值到第三方账户"
	}
	return acc.RequestTakeOut(int(takeKind), title, applyAmount, commission)
}

// 获取最近的提现描述
func (s *memberService) GetLatestApplyCashText(memberId int64) string {
	var latestInfo string
	latestApplyInfo := s.query.GetLatestWalletLogByKind(memberId,
		member.KindWalletTakeOutToBankCard)
	if latestApplyInfo != nil {
		var sText string
		switch latestApplyInfo.ReviewState {
		case enum.ReviewAwaiting:
			sText = "已申请"
		case enum.ReviewPass:
			sText = "已审核,等待打款"
		case enum.ReviewReject:
			sText = "被退回"
		case enum.ReviewConfirm:
			sText = "已完成"
		}
		if latestApplyInfo.Amount < 0 {
			latestApplyInfo.Amount = -latestApplyInfo.Amount
		}
		latestInfo = fmt.Sprintf(`<b>最近提现：</b>%s&nbsp;申请提现%s ，状态：<span class="status">%s</span>。`,
			time.Unix(latestApplyInfo.CreateTime, 0).Format("2006-01-02 15:04"),
			format.FormatFloat(latestApplyInfo.Amount),
			sText)
	}
	return latestInfo
}

// 确认提现
func (s *memberService) ConfirmTakeOutRequest(memberId int64,
	infoId int32, pass bool, remark string) error {
	m, err := s.getMember(memberId)
	if err == nil {
		err = m.GetAccount().ConfirmTakeOut(infoId, pass, remark)
	}
	return err
}

// 完成提现
func (s *memberService) FinishTakeOutRequest(memberId int64, id int32, tradeNo string) error {
	m, err := s.getMember(memberId)
	if err != nil {
		return err
	}
	return m.GetAccount().FinishTakeOut(id, tradeNo)
}

// 冻结余额
func (s *memberService) Freeze(memberId int64, title string,
	tradeNo string, amount float32, referId int64) error {
	m := s.repo.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().Freeze(title, tradeNo, amount, referId)
}

// 解冻金额
func (s *memberService) Unfreeze(memberId int64, title string,
	tradeNo string, amount float32, referId int64) error {
	m := s.repo.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().Unfreeze(title, tradeNo, amount, referId)
}

// 冻结赠送金额
func (s *memberService) FreezeWallet(memberId int64, title string,
	tradeNo string, amount float32, referId int64) error {
	m := s.repo.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().FreezeWallet(title, tradeNo, amount, referId)
}

// 解冻赠送金额
func (s *memberService) UnfreezeWallet(memberId int64, title string,
	tradeNo string, amount float32, referId int64) error {
	m := s.repo.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().UnfreezeWallet(title, tradeNo, amount, referId)
}

// 将冻结金额标记为失效
func (s *memberService) FreezeExpired(memberId int64, accountKind int, amount float32,
	remark string) error {
	m := s.repo.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().FreezeExpired(accountKind, amount, remark)
}

// 转账余额到其他账户
func (s *memberService) TransferAccount(accountKind int, fromMember int64,
	toMember int64, amount float32, csnRate float32, remark string) error {
	m := s.repo.GetMember(fromMember)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().TransferAccount(accountKind, toMember,
		amount, csnRate, remark)
}

// 转账余额到其他账户
func (s *memberService) TransferBalance(memberId int64, kind int32, amount float32, tradeNo string,
	toTitle, fromTitle string) error {
	m := s.repo.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().TransferBalance(int(kind), amount, tradeNo, toTitle, fromTitle)
}

// 转账活动账户,kind为转账类型，如 KindBalanceTransfer等
// commission手续费
func (s *memberService) TransferFlow(memberId int64, kind int32, amount float32,
	commission float32, tradeNo string, toTitle string, fromTitle string) error {
	m := s.repo.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().TransferFlow(int(kind), amount, commission, tradeNo,
		toTitle, fromTitle)
}

// 将活动金转给其他人
func (s *memberService) TransferFlowTo(memberId int64, toMemberId int64, kind int32,
	amount float32, commission float32, tradeNo string, toTitle string,
	fromTitle string) error {
	m := s.repo.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().TransferFlowTo(toMemberId, int(kind), amount,
		commission, tradeNo, toTitle, fromTitle)
}

// 根据用户或手机筛选会员
func (s *memberService) FilterMemberByUserOrPhone(key string) []*dto.SimpleMember {
	return s.query.FilterMemberByUserOrPhone(key)
}

// 根据用户名货手机获取会员
func (s *memberService) GetMemberByUserOrPhone(key string) *dto.SimpleMember {
	return s.query.GetMemberByUserOrPhone(key)
}

// 根据手机获取会员编号
func (s *memberService) GetMemberIdByPhone1(phone string) int64 {
	return s.query.GetMemberIdByPhone(phone)
}

// 会员推广排名
func (s *memberService) GetMemberInviRank(mchId int32, allTeam bool,
	levelComp string, level int, startTime int64, endTime int64,
	num int) []*dto.RankMember {
	return s.query.GetMemberInviRank(mchId, allTeam, levelComp, level, startTime, endTime, num)
}

//********* 促销  **********//

// 可用的优惠券分页数据
func (s *memberService) PagedAvailableCoupon(memberId int, start, end int) (int, []*dto.SimpleCoupon) {
	return s.repo.CreateMemberById(int64(memberId)).
		GiftCard().PagedAvailableCoupon(start, end)
}

// 已使用的优惠券
func (s *memberService) PagedAllCoupon(memberId int, start, end int) (int, []*dto.SimpleCoupon) {
	return s.repo.CreateMemberById(int64(memberId)).
		GiftCard().PagedAllCoupon(start, end)
}

// 过期的优惠券
func (s *memberService) PagedExpiresCoupon(memberId int, start, end int) (int, []*dto.SimpleCoupon) {
	return s.repo.CreateMemberById(int64(memberId)).
		GiftCard().PagedExpiresCoupon(start, end)
}

// 更改手机号
func (s *memberService) changePhone(memberId int64, phone string) error {
	m := s.repo.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.Profile().ChangePhone(phone)
}

func (s *memberService) parseLevelDto(src *member.Level) *proto.SMemberLevel {
	return &proto.SMemberLevel{
		ID:            int32(src.ID),
		Name:          src.Name,
		RequireExp:    int32(src.RequireExp),
		ProgramSignal: src.ProgramSignal,
		Enabled:       int32(src.Enabled),
		IsOfficial:    int32(src.IsOfficial),
	}
}

func (s *memberService) parseMemberDto(src *member.Member) *proto.SMember {
	return &proto.SMember{
		Id:             src.Id,
		User:           src.User,
		Code:           src.Code,
		Pwd:            src.Pwd,
		TradePwd:       src.TradePwd,
		Exp:            int64(src.Exp),
		Level:          int32(src.Level),
		PremiumUser:    int32(src.PremiumUser),
		PremiumExpires: src.PremiumExpires,
		InviteCode:     src.InviteCode,
		RegIp:          src.RegIp,
		RegFrom:        src.RegFrom,
		State:          int32(src.State),
		Flag:           int32(src.Flag),
		Avatar:         src.Avatar,
		Phone:          src.Phone,
		Email:          src.Email,
		Name:           src.Name,
		RealName:       src.RealName,
		DynamicToken:   src.DynamicToken,
		RegTime:        src.RegTime,
		LastLoginTime:  src.LastLoginTime,
	}
}

func (s *memberService) parseMemberProfile(src *member.Profile) *proto.SProfile {
	return &proto.SProfile{
		MemberId:   src.MemberId,
		Name:       src.Name,
		Avatar:     src.Avatar,
		Sex:        src.Sex,
		BirthDay:   src.BirthDay,
		Phone:      src.Phone,
		Address:    src.Address,
		Im:         src.Im,
		Email:      src.Email,
		Province:   src.Province,
		City:       src.City,
		District:   src.District,
		Remark:     src.Remark,
		Ext1:       src.Ext1,
		Ext2:       src.Ext2,
		Ext3:       src.Ext3,
		Ext4:       src.Ext4,
		Ext5:       src.Ext5,
		Ext6:       src.Ext6,
		UpdateTime: src.UpdateTime,
	}
}

func (s *memberService) parseComplexMemberDto(src *member.ComplexMember) *proto.SComplexMember {
	return &proto.SComplexMember{
		Name:           src.Name,
		Avatar:         src.Avatar,
		Exp:            int32(src.Exp),
		Level:          int32(src.Level),
		LevelName:      src.LevelName,
		PremiumUser:    int32(src.PremiumUser),
		InviteCode:     src.InviteCode,
		TrustAuthState: int32(src.TrustAuthState),
		TradePwdHasSet: src.TradePwdHasSet,
		UpdateTime:     src.UpdateTime,
	}
}

func (s *memberService) parseAddressDto(src *member.Address) *proto.SAddress {
	return &proto.SAddress{
		ID:             src.ID,
		ConsigneeName:  src.ConsigneeName,
		ConsigneePhone: src.ConsigneePhone,
		Province:       src.Province,
		City:           src.City,
		District:       src.District,
		Area:           src.Area,
		DetailAddress:  src.DetailAddress,
		IsDefault:      int32(src.IsDefault),
	}
}
func round(f float32, n int) float64 {
	return math.Round(float64(f), n)
}
func (s *memberService) parseAccountDto(src *member.Account) *proto.SAccount {
	return &proto.SAccount{
		MemberId:          src.MemberId,
		Integral:          int64(src.Integral),
		FreezeIntegral:    int64(src.FreezeIntegral),
		Balance:           round(src.Balance, 2),
		FreezeBalance:     round(src.FreezeBalance, 2),
		ExpiredBalance:    round(src.ExpiredBalance, 2),
		WalletBalance:     round(src.WalletBalance, 2),
		FreezeWallet:      round(src.FreezeWallet, 2),
		ExpiredWallet:     round(src.ExpiredWallet, 2),
		TotalWalletAmount: round(src.TotalWalletAmount, 2),
		FlowBalance:       round(src.FlowBalance, 2),
		GrowBalance:       round(src.GrowBalance, 2),
		GrowAmount:        round(src.GrowAmount, 2),
		GrowEarnings:      round(src.GrowEarnings, 2),
		GrowTotalEarnings: round(src.GrowTotalEarnings, 2),
		TotalExpense:      round(src.TotalExpense, 2),
		TotalCharge:       round(src.TotalCharge, 2),
		TotalPay:          round(src.TotalPay, 2),
		PriorityPay:       int64(src.PriorityPay),
		UpdateTime:        src.UpdateTime,
	}
}

func (s *memberService) parseMember(src *proto.SMember) *member.Member {
	return &member.Member{
		Id:             int64(src.Id),
		Code:           src.Code,
		Name:           src.Name,
		RealName:       src.RealName,
		User:           src.User,
		Pwd:            src.Pwd,
		Avatar:         src.Avatar,
		TradePwd:       src.TradePwd,
		Exp:            int(src.Exp),
		Level:          int(src.Level),
		InviteCode:     src.InviteCode,
		PremiumUser:    int(src.PremiumUser),
		PremiumExpires: src.PremiumExpires,
		Phone:          src.Phone,
		Email:          src.Email,
		RegFrom:        src.RegFrom,
		RegIp:          src.RegIp,
		Flag:           int(src.Flag),
		State:          int(src.State),
		DynamicToken:   src.DynamicToken,
	}
}

func (s *memberService) parseMemberProfile2(src *proto.SProfile) *member.Profile {
	return &member.Profile{
		MemberId:   src.MemberId,
		Name:       src.Name,
		Avatar:     src.Avatar,
		Sex:        src.Sex,
		BirthDay:   src.BirthDay,
		Phone:      src.Phone,
		Address:    src.Address,
		Im:         src.Im,
		Email:      src.Email,
		Province:   src.Province,
		City:       src.City,
		District:   src.District,
		Remark:     src.Remark,
		Ext1:       src.Ext1,
		Ext2:       src.Ext2,
		Ext3:       src.Ext3,
		Ext4:       src.Ext4,
		Ext5:       src.Ext5,
		Ext6:       src.Ext6,
		UpdateTime: src.UpdateTime,
	}
}

func (s *memberService) parseAddress(src *proto.SAddress) *member.Address {
	return &member.Address{
		ID:             src.ID,
		ConsigneeName:  src.ConsigneeName,
		ConsigneePhone: src.ConsigneePhone,
		Province:       src.Province,
		City:           src.City,
		District:       src.District,
		Area:           src.Area,
		DetailAddress:  src.DetailAddress,
		IsDefault:      int(src.IsDefault),
	}
}
