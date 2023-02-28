package impl

/**
 * Copyright 2014 @ 56x.net.
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
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	de "github.com/ixre/go2o/core/domain/interface/domain"
	"github.com/ixre/go2o/core/domain/interface/domain/enum"
	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/domain/interface/valueobject"
	"github.com/ixre/go2o/core/domain/interface/wallet"
	"github.com/ixre/go2o/core/dto"
	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/go2o/core/infrastructure/format"
	"github.com/ixre/go2o/core/infrastructure/regex"
	"github.com/ixre/go2o/core/module"
	"github.com/ixre/go2o/core/query"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/go2o/core/variable"
	api "github.com/ixre/gof/jwt-api"
	"github.com/ixre/gof/types"
	"github.com/ixre/gof/types/typeconv"
	"github.com/ixre/gof/util"
)

var _ proto.MemberServiceServer = new(memberService)

type memberService struct {
	repo         member.IMemberRepo
	registryRepo registry.IRegistryRepo
	mchService   *merchantService
	query        *query.MemberQuery
	orderQuery   *query.OrderQuery
	valRepo      valueobject.IValueRepo
	serviceUtil
	proto.UnimplementedMemberServiceServer
}

func NewMemberService(mchService *merchantService, repo member.IMemberRepo,
	registryRepo registry.IRegistryRepo,
	q *query.MemberQuery, oq *query.OrderQuery,
	valRepo valueobject.IValueRepo) *memberService {
	s := &memberService{
		repo:         repo,
		registryRepo: registryRepo,
		query:        q,
		mchService:   mchService,
		orderQuery:   oq,
		valRepo:      valRepo,
	}
	return s
	//return _s.init()
}

// FindMember 交换会员编号
func (s *memberService) FindMember(_ context.Context, r *proto.FindMemberRequest) (*proto.Int64, error) {
	var memberId int64
	switch r.Cred {
	default:
	case proto.ECredentials_USER:
		memberId = s.repo.GetMemberIdByUser(r.Value)
	case proto.ECredentials_CODE:
		memberId = s.repo.GetMemberIdByCode(r.Value)
	case proto.ECredentials_PHONE:
		memberId = s.repo.GetMemberIdByPhone(r.Value)
	case proto.ECredentials_EMAIL:
		memberId = s.repo.GetMemberIdByEmail(r.Value)
	}
	return &proto.Int64{Value: memberId}, nil
}

//func (_s *memberService) init() *memberService {
//	db := gof.CurrentApp.Db()
//	var list []*member.Member
//	db.o.Select(&list, "")
//	for _, v := range list {
//		im := _s.repo.CreateMember(v)
//		if rl := im.GetRelation(); rl != nil {
//			im.BindInviter(rl.InviterId, true)
//		}
//		//if len(v.InviteCode) < 6 {
//		//	im := _s.repo.CreateMember(v)
//		//	v.InviteCode = _s.generateInviteCode()
//		//	im.SetValue(v)
//		//	im.Save()
//		//}
//	}
//	return _s
//}

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

// GetMember 根据会员编号获取会员
func (s *memberService) GetMember(_ context.Context, id *proto.MemberIdRequest) (*proto.SMember, error) {
	iv := s.repo.GetMember(id.MemberId)
	if iv != nil {
		v := iv.GetValue()
		if len(v.TradePassword) == 0 {
			v.UserFlag |= member.FlagNoTradePasswd
		}
		return s.parseMemberDto(&v), nil
	}
	return nil, member.ErrNoSuchMember
}

// GetProfile 获取资料
func (s *memberService) GetProfile(_ context.Context, id *proto.MemberIdRequest) (*proto.SProfile, error) {
	m := s.repo.GetMember(id.MemberId)
	if m != nil {
		v := m.Profile().GetProfile()
		return s.parseMemberProfile(&v), nil
	}
	return nil, member.ErrNoSuchMember
}

// SaveProfile 保存资料
func (s *memberService) SaveProfile(_ context.Context, v *proto.SProfile) (*proto.Result, error) {
	if v.MemberId > 0 {
		v2 := s.parseMemberProfile2(v)
		m := s.repo.GetMember(v.MemberId)
		if m != nil {
			err := m.Profile().SaveProfile(v2)
			return s.error(err), nil
		}
	}
	return s.error(member.ErrNoSuchMember), nil
}

// Premium 升级为高级会员
func (s *memberService) Premium(_ context.Context, r *proto.PremiumRequest) (*proto.Result, error) {
	m := s.repo.GetMember(r.MemberId)
	if m == nil {
		return s.result(member.ErrNoSuchMember), nil
	}
	err := m.Premium(int(r.Value), r.Expires)
	return s.result(err), nil
}

// CheckToken 检查会员的会话Token是否正确
func (s *memberService) CheckToken(_ context.Context, r *proto.CheckTokenRequest) (*proto.Bool, error) {
	md := module.Get(module.MM).(*module.MemberModule)
	return &proto.Bool{
		Value: md.CheckToken(r.MemberId, r.Token),
	}, nil
}

// GetToken 获取会员的会员Token,reset表示是否重置会员的token
func (s *memberService) GetToken(_ context.Context, r *proto.GetTokenRequest) (*proto.String, error) {
	pubToken := ""
	md := module.Get(module.MM).(*module.MemberModule)
	if !r.Reset_ {
		pubToken = md.GetToken(r.MemberId)
	}
	if r.Reset_ || (pubToken == "" && r.MemberId > 0) {
		m := s.getMemberValue(r.MemberId)
		if m != nil {
			return &proto.String{Value: md.ResetToken(r.MemberId, m.Password)}, nil
		}
	}
	return &proto.String{Value: pubToken}, nil
}

// RemoveToken 移除会员的Token
func (s *memberService) RemoveToken(_ context.Context, id *proto.MemberIdRequest) (*proto.Empty, error) {
	md := module.Get(module.MM).(*module.MemberModule)
	md.RemoveToken(id.MemberId)
	return &proto.Empty{}, nil
}

// ChangePhone 更改手机号码，不验证手机格式
func (s *memberService) ChangePhone(_ context.Context, r *proto.ChangePhoneRequest) (*proto.Result, error) {
	err := s.changePhone(r.MemberId, r.Phone)
	return s.result(err), nil
}

// ChangeNickname 更改昵称
func (s *memberService) ChangeNickname(_ context.Context, req *proto.ChangeNicknameRequest) (*proto.Result, error) {
	m := s.repo.GetMember(req.MemberId)
	if m == nil {
		return s.error(member.ErrNoSuchMember), nil
	}
	err := m.Profile().ChangeNickname(req.Nickname, req.LimitTime)
	return s.result(err), nil
}

// ChangeInviterId 更改邀请人
func (s *memberService) SetInviter(_ context.Context, r *proto.SetInviterRequest) (*proto.Result, error) {
	im := s.repo.GetMember(r.MemberId)
	if im == nil {
		return s.result(member.ErrNoSuchMember), nil
	}
	inviterId := s.repo.GetMemberIdByCode(r.InviterCode)
	if inviterId <= 0 {
		return s.result(member.ErrInvalidInviter), nil
	}
	err := im.BindInviter(inviterId, r.AllowChange)
	return s.result(err), nil
}

// RemoveFavorite 取消收藏
func (s *memberService) RemoveFavorite(_ context.Context, r *proto.FavoriteRequest) (rs *proto.Result, err error) {
	f := s.repo.CreateMemberById(r.MemberId).Favorite()
	switch r.FavoriteType {
	case proto.FavoriteType_SHOP:
		err = f.Cancel(member.FavTypeShop, r.ReferId)
	case proto.FavoriteType_GOOGS:
		err = f.Cancel(member.FavTypeGoods, r.ReferId)
	}
	if err != nil {
		return s.error(err), nil
	}
	return s.success(nil), nil
}

func (s *memberService) Favorite(_ context.Context, r *proto.FavoriteRequest) (rs *proto.Result, err error) {
	m := s.repo.CreateMemberById(r.MemberId)
	if m == nil {
		return s.error(member.ErrNoSuchMember), nil
	}
	f := m.Favorite()
	switch r.FavoriteType {
	case proto.FavoriteType_SHOP:
		err = f.Favorite(member.FavTypeShop, r.ReferId)
	case proto.FavoriteType_GOOGS:
		err = f.Favorite(member.FavTypeGoods, r.ReferId)
	}
	if err != nil {
		return s.error(err), nil
	}
	return s.success(nil), nil
}

// IsFavored 是否已收藏
func (s *memberService) IsFavored(c context.Context, r *proto.FavoriteRequest) (*proto.Bool, error) {
	m := s.repo.CreateMemberById(r.MemberId)
	if m == nil {
		return &proto.Bool{}, nil
	}
	f := m.Favorite()
	t := member.FavTypeGoods
	switch r.FavoriteType {
	case proto.FavoriteType_SHOP:
		t = member.FavTypeShop
	case proto.FavoriteType_GOOGS:
		t = member.FavTypeGoods
	}
	b := f.Favored(t, r.ReferId)
	return &proto.Bool{Value: b}, nil
}

// 获取所有会员等级
func (s *memberService) GetMemberLevels() []*member.Level {
	return s.repo.GetManager().LevelManager().GetLevelSet()
}

// GetLevels 等级列表
func (s *memberService) GetLevels(_ context.Context, empty *proto.Empty) (*proto.SMemberLevelListResponse, error) {
	var arr []*proto.SMemberLevel
	list := s.repo.GetManager().LevelManager().GetLevelSet()
	for _, v := range list {
		arr = append(arr, s.parseLevelDto(v))
	}
	return &proto.SMemberLevelListResponse{Value: arr}, nil
}

// GetMemberLevel 根据编号获取会员等级信息
func (s *memberService) GetMemberLevel(_ context.Context, i *proto.Int32) (*proto.SMemberLevel, error) {
	lv := s.repo.GetManager().LevelManager().GetLevelById(int(i.Value))
	if lv != nil {
		return s.parseLevelDto(lv), nil
	}
	return nil, member.ErrNoSuchLevelUpLog
}

// SaveMemberLevel 保存会员等级信息
func (s *memberService) SaveMemberLevel(_ context.Context, level *proto.SMemberLevel) (*proto.Result, error) {
	lv := &member.Level{
		Id:            int(level.Id),
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

// GetLevelBySign 根据SIGN获取等级
func (s *memberService) GetLevelBySign(_ context.Context, sign *proto.String) (*proto.SMemberLevel, error) {
	lv := s.repo.GetManager().LevelManager().GetLevelByProgramSign(sign.Value)
	if lv != nil {
		return s.parseLevelDto(lv), nil
	}
	return nil, member.ErrNoSuchLevel
}

// DeleteMemberLevel 删除会员等级
func (s *memberService) DeleteMemberLevel(_ context.Context, levelId *proto.Int64) (*proto.Result, error) {
	err := s.repo.GetManager().LevelManager().DeleteLevel(int(levelId.Value))
	return s.result(err), nil
}

// GetHighestLevel 获取启用中的最大等级,用于判断是否可以升级
func (s *memberService) GetHighestLevel() member.Level {
	lv := s.repo.GetManager().LevelManager().GetHighestLevel()
	if lv != nil {
		return *lv
	}
	return member.Level{}
}

func (s *memberService) GetWalletLog(_ context.Context, r *proto.WalletLogRequest) (*proto.WalletLogResponse, error) {
	m := s.repo.GetMember(r.MemberId)
	v := m.GetAccount().GetWalletLog(r.LogId)
	return &proto.WalletLogResponse{
		LogId:       v.Id,
		MemberId:    r.MemberId,
		OuterNo:     v.OuterNo,
		Kind:        int32(v.Kind),
		Title:       v.Subject,
		Amount:      float64(v.ChangeValue),
		TradeFee:    float64(v.ProcedureFee),
		ReviewState: int32(v.ReviewState),
		Remark:      v.Remark,
		CreateTime:  v.CreateTime,
		UpdateTime:  v.UpdateTime,
		RelateUser:  int64(v.OperatorUid),
	}, nil
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

// SendCode 发送会员验证码消息, 并返回验证码, 验证码通过data.code获取
func (s *memberService) SendCode(_ context.Context, r *proto.SendCodeRequest) (*proto.Result, error) {
	m := s.repo.GetMember(r.MemberId)
	if m == nil {
		return s.error(member.ErrNoSuchMember), nil
	}
	code, err := m.SendCheckCode(r.Operation, int(r.MsgType))
	if err != nil {
		return s.error(err), nil
	}
	return s.success(map[string]string{"code": code}), nil
}

// CompareCode 比较验证码是否正确
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

// ChangeUsername 更改会员用户名
func (s *memberService) ChangeUsername(_ context.Context, r *proto.ChangeUsernameRequest) (*proto.Result, error) {
	var err error
	m := s.repo.GetMember(int64(r.MemberId))
	if m == nil {
		err = member.ErrNoSuchMember
	} else {
		if err = m.ChangeUsername(r.Username); err == nil {
			return s.success(nil), nil
		}
	}
	return s.result(err), nil
}

// MemberLevelInfo 获取会员等级信息
func (s *memberService) MemberLevelInfo(_ context.Context, id *proto.MemberIdRequest) (*proto.SMemberLevelInfo, error) {
	level := &proto.SMemberLevelInfo{Level: -1}
	im := s.repo.GetMember(id.MemberId)
	if im != nil {
		v := im.GetValue()
		level.Exp = int32(v.Exp)
		level.Level = int32(v.Level)
		lv := im.GetLevel()
		level.LevelName = lv.Name
		level.ProgramSignal = lv.ProgramSignal
		nextLv := s.repo.GetManager().LevelManager().GetNextLevelById(lv.Id)
		if nextLv == nil {
			level.NextLevel = -1
		} else {
			level.NextLevel = int32(nextLv.Id)
			level.NextLevelName = nextLv.Name
			level.NextProgramSignal = nextLv.ProgramSignal
			level.RequireExp = int32(nextLv.RequireExp - v.Exp)
		}
	}
	return level, nil
}

// ChangeLevel 更改会员等级
func (s *memberService) ChangeLevel(_ context.Context, r *proto.ChangeLevelRequest) (*proto.Result, error) {
	if len(r.LevelCode) > 0 {
		if r.Level != 0 {
			return s.error(errors.New("levelCode和level不能同时设置")), nil
		}
		lv := s.repo.GetManager().LevelManager().GetLevelByProgramSign(r.LevelCode)
		if lv == nil {
			return s.error(fmt.Errorf("no such level, code=%s", r.LevelCode)), nil
		}
		r.Level = int32(lv.Id)
	}
	m := s.repo.GetMember(r.MemberId)
	var err error
	if m == nil {
		err = member.ErrNoSuchMember
	} else {
		err = m.ChangeLevel(int(r.Level), int(r.PaymentOrderId), r.Review)
	}
	return s.result(err), nil
}

func (s *memberService) ReviewLevelUpRequest(_ context.Context, r *proto.LevelUpReviewRequest) (*proto.Result, error) {
	m := s.repo.GetMember(r.MemberId)
	var err error
	if m == nil {
		err = member.ErrNoSuchMember
	} else {
		err = m.ReviewLevelUp(int(r.RequestId), r.ReviewPass)
	}
	return s.result(err), nil
}

func (s *memberService) ConfirmLevelUpRequest(_ context.Context, r *proto.LevelUpConfirmRequest) (*proto.Result, error) {
	m := s.repo.GetMember(r.MemberId)
	var err error
	if m == nil {
		err = member.ErrNoSuchMember
	} else {
		err = m.ConfirmLevelUp(int(r.RequestId))
	}
	return s.result(err), nil
}

// ChangeHeadPortrait 上传会员头像
func (s *memberService) ChangeHeadPortrait(_ context.Context, r *proto.ChangePortraitRequest) (*proto.Result, error) {
	m := s.repo.GetMember(r.MemberId)
	if m == nil {
		return s.error(member.ErrNoSuchMember), nil
	}
	err := m.Profile().ChangeHeadPortrait(r.PortraitUrl)
	return s.result(err), nil
}

// Register 注册会员
func (s *memberService) Register(_ context.Context, r *proto.RegisterMemberRequest) (*proto.RegisterResponse, error) {
	if len(r.Password) != 32 {
		return &proto.RegisterResponse{
			ErrCode: 1,
			ErrMsg:  de.ErrNotMD5Format.Error(),
		}, nil
	}
	salt := util.RandString(6)
	v := &member.Member{
		Username: r.Username,
		Salt:     salt,
		Password: domain.Sha1Pwd(r.Password, salt),
		Nickname: r.Nickname,
		RealName: "",
		Portrait: "", //todo: default avatar
		Phone:    r.Phone,
		Email:    r.Email,
		RegFrom:  r.RegFrom,
		RegIp:    r.RegIp,
		UserFlag: int(r.Flag),
	}
	// 验证邀请码
	inviterId, err := s.repo.GetManager().CheckInviteRegister(r.InviterCode)
	if err != nil {
		return &proto.RegisterResponse{
			ErrCode: 2,
			ErrMsg:  err.Error(),
		}, nil
	}
	// 创建会员
	m := s.repo.CreateMember(v)
	id, err := m.Save()

	if err == nil {
		// 保存关联信息
		err = m.BindInviter(inviterId, true)
		//m := _s.repo.CreateMember(v) //创建会员
		//id, err := m.Save()
		//if err == nil {
		//	pro.Gender = 1
		//	pro.MemberId = id
		//	//todo: 如果注册失败，则删除。应使用SQL-TRANSFER
		//	if err = m.Profile().SaveProfile(pro); err != nil {
		//		_s.repo.DeleteMember(id)
		//}
	}
	ret := &proto.RegisterResponse{MemberId: id}
	if err != nil {
		ret.ErrCode = 1
		ret.ErrMsg = err.Error()
	}
	return ret, nil
}

// GetInviter 获取会员邀请关系
func (s *memberService) GetInviter(_ context.Context, id *proto.MemberIdRequest) (*proto.MemberInviterResponse, error) {
	r := s.repo.GetRelation(id.MemberId)
	if r != nil {
		ret := &proto.MemberInviterResponse{
			InviterId: r.InviterId,
			InviterD2: r.InviterD2,
			InviterD3: r.InviterD3,
		}
		if r.InviterId > 0 {
			if mm := s.repo.GetMember(r.InviterId); mm != nil {
				mv := mm.GetValue()
				ret.InviterUsername = mv.Username
				ret.InviterNickname = mv.Nickname
				ret.InviterPortrait = mv.Portrait
				ret.InviterPhone = mv.Phone
			}
		}
		return ret, nil
	}
	return &proto.MemberInviterResponse{}, nil
}

// GetInviteCount 获取会员邀请数量
func (s *memberService) GetInviteCount(_ context.Context, req *proto.MemberIdRequest) (*proto.MemberInviteCountResponse, error) {
	memberId := int(req.MemberId)
	if memberId > 0 {
		f := func(level int) int32 {
			return int32(s.repo.GetInvitationCount(int(req.MemberId), level))
		}
		return &proto.MemberInviteCountResponse{
			FirstLevelCount: f(1),
			SecondCount:     f(2),
			ThridCount:      f(3),
		}, nil
	}
	return &proto.MemberInviteCountResponse{}, nil

}

// Active 激活会员
func (s *memberService) Active(_ context.Context, id *proto.MemberIdRequest) (*proto.Result, error) {
	m := s.repo.GetMember(id.MemberId)
	if m == nil {
		return s.error(member.ErrNoSuchMember), nil
	}
	if err := m.Active(); err != nil {
		return s.error(err), nil
	}
	return s.success(nil), nil
}

// Lock 锁定/解锁会员
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
func (s *memberService) Unlock(_ context.Context, id *proto.MemberIdRequest) (*proto.Result, error) {
	m := s.repo.GetMember(id.MemberId)
	if m == nil {
		return s.error(member.ErrNoSuchMember), nil
	}
	if err := m.Unlock(); err != nil {
		return s.error(err), nil
	}
	return s.success(nil), nil
}

// CheckProfileCompleted 判断资料是否完善
func (s *memberService) CheckProfileCompleted(_ context.Context, memberId *proto.Int64) (*proto.Bool, error) {
	m := s.repo.GetMember(memberId.Value)
	if m != nil {
		return &proto.Bool{Value: m.Profile().ProfileCompleted()}, nil
	}
	return &proto.Bool{}, nil
}

// CheckProfileComplete 判断资料是否完善
func (s *memberService) CheckProfileComplete(_ context.Context, id *proto.MemberIdRequest) (*proto.Result, error) {
	m := s.repo.GetMember(id.MemberId)
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

// ChangePassword 更改密码
func (s *memberService) ChangePassword(_ context.Context, r *proto.ChangePasswordRequest) (*proto.Result, error) {
	m := s.repo.GetMember(r.MemberId)
	if m == nil {
		return s.error(member.ErrNoSuchMember), nil
	}
	v := m.GetValue()
	pwd := r.NewPassword
	old := r.OriginPassword
	if l := len(r.NewPassword); l != 32 {
		return s.error(de.ErrNotMD5Format), nil
	} else {
		pwd = domain.MemberSha1Pwd(pwd, v.Salt)
	}
	if l := len(old); l > 0 && l != 32 {
		return s.error(de.ErrNotMD5Format), nil
	} else {
		old = domain.MemberSha1Pwd(old, v.Salt)
	}
	log.Println("--password", pwd, v.Password, v.Salt, typeconv.MustJson(v))
	err := m.Profile().ChangePassword(pwd, old)
	if err != nil {
		return s.error(err), nil
	}
	return s.success(nil), nil
}

// ChangeTradePassword 更改交易密码
func (s *memberService) ChangeTradePassword(_ context.Context, r *proto.ChangePasswordRequest) (*proto.Result, error) {
	m := s.repo.GetMember(r.MemberId)
	if m == nil {
		return s.error(member.ErrNoSuchMember), nil
	}
	pwd, old := r.NewPassword, r.OriginPassword
	v := m.GetValue()
	if l := len(pwd); l != 32 {
		return s.error(de.ErrNotMD5Format), nil
	} else {
		pwd = domain.TradePassword(pwd, v.Salt)
	}
	if l := len(old); l > 0 && l != 32 {
		return s.error(de.ErrNotMD5Format), nil
	} else {
		old = domain.TradePassword(old, v.Salt)
	}
	err := m.Profile().ChangeTradePassword(pwd, old)
	if err != nil {
		return s.error(err), nil
	}
	return s.success(nil), nil
}

// 登录，返回结果(Result_)和会员编号(Id);
// Result值为：-1:会员不存在; -2:账号密码不正确; -3:账号被停用
func (s *memberService) tryLogin(user string, pwd string, update bool) (v *member.Member, errCode int32, err error) {
	user = strings.ToLower(user)
	if len(pwd) != 32 {
		return nil, 4, de.ErrNotMD5Format
	}
	memberId := s.repo.GetMemberIdByUser(user)
	if memberId <= 0 {
		// 用户名不正确时,尝试匹配手机号
		if regex.IsPhone(user) {
			memberId = s.repo.GetMemberIdByPhone(user)
		}
	}
	if memberId <= 0 {
		return nil, 2, de.ErrCredential // 用户不存在,也返回用户或密码不正确
	}
	im := s.repo.GetMember(memberId)
	val := im.GetValue()

	if s := domain.Sha1Pwd(pwd, val.Salt); s != val.Password {
		return nil, 1, de.ErrCredential
	}
	if val.UserFlag&member.FlagLocked == member.FlagLocked {
		return nil, 3, member.ErrMemberLocked
	}
	if update {
		go im.UpdateLoginTime()
	}
	return &val, 0, nil
}

// CheckLogin 登录，返回结果(Result_)和会员编号(Id);
// Result值为：-1:会员不存在; -2:账号密码不正确; -3:账号被停用
func (s *memberService) CheckLogin(_ context.Context, r *proto.LoginRequest) (*proto.LoginResponse, error) {
	v, code, err := s.tryLogin(r.Username, r.Password, r.Update)
	ret := &proto.LoginResponse{
		ErrCode: code,
	}
	if err != nil {
		ret.ErrMsg = err.Error()
		return ret, nil
	} else {
		ret.MemberId = v.Id
		ret.UserCode = v.UserCode
	}
	return ret, nil
}

// GrantAccessToken 发放访问令牌
func (s *memberService) GrantAccessToken(_ context.Context, request *proto.GrantAccessTokenRequest) (*proto.GrantAccessTokenResponse, error) {
	now := time.Now().Unix()
	if request.ExpiresTime <= now {
		return &proto.GrantAccessTokenResponse{
			Error: fmt.Sprintf("令牌有效时间已过有效期: value=%d", request.ExpiresTime),
		}, nil
	}
	im := s.repo.GetMember(request.MemberId)
	if im == nil {
		return &proto.GrantAccessTokenResponse{Error: member.ErrNoSuchMember.Error()}, nil
	}
	// 创建token并返回
	claims := api.CreateClaims(strconv.Itoa(int(request.MemberId)), "go2o",
		"go2o-api-jwt", request.ExpiresTime).(jwt.MapClaims)
	jwtSecret, err := s.registryRepo.GetValue(registry.SysJWTSecret)
	if err != nil {
		log.Println("[ GO2O][ ERROR]: grant access token error ", err.Error())
		return &proto.GrantAccessTokenResponse{Error: err.Error()}, nil
	}
	token, err := api.AccessToken(claims, []byte(jwtSecret))
	if err != nil {
		log.Println("[ GO2O][ ERROR]: grant access token error ", err.Error())
		return &proto.GrantAccessTokenResponse{Error: err.Error()}, nil
	}
	return &proto.GrantAccessTokenResponse{
		AccessToken: token,
	}, nil
}

// CheckAccessToken 检查令牌是否有效
func (s *memberService) CheckAccessToken(_ context.Context, request *proto.CheckAccessTokenRequest) (*proto.CheckAccessTokenResponse, error) {
	if len(request.AccessToken) == 0 {
		return &proto.CheckAccessTokenResponse{Error: "令牌不能为空"}, nil
	}
	jwtSecret, err := s.registryRepo.GetValue(registry.SysJWTSecret)
	if err != nil {
		log.Println("[ GO2O][ ERROR]: check access token error ", err.Error())
		return &proto.CheckAccessTokenResponse{Error: err.Error()}, nil
	}
	// 去掉"Bearer "
	if len(request.AccessToken) > 6 &&
		strings.HasPrefix(request.AccessToken, "Bearer") {
		request.AccessToken = request.AccessToken[7:]
	}
	// 转换token
	dstClaims := jwt.MapClaims{} // 可以用实现了Claim接口的自定义结构
	tk, err := jwt.ParseWithClaims(request.AccessToken, &dstClaims, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if tk == nil {
		return &proto.CheckAccessTokenResponse{Error: "令牌无效"}, nil
	}
	if !dstClaims.VerifyIssuer("go2o", true) ||
		dstClaims["sub"] != "go2o-api-jwt" {
		return &proto.CheckAccessTokenResponse{Error: "未知颁发者的令牌"}, nil
	}
	// 令牌过期时间
	exp := int64(dstClaims["exp"].(float64))
	// 判断是否有效
	if !tk.Valid {
		ve, _ := err.(*jwt.ValidationError)
		if ve.Errors&jwt.ValidationErrorExpired != 0 {
			return &proto.CheckAccessTokenResponse{
				Error:            "令牌已过期",
				IsExpires:        true,
				TokenExpiresTime: exp,
			}, nil
		}
		return &proto.CheckAccessTokenResponse{Error: "令牌无效:" + ve.Error()}, nil
	}
	aud := int64(typeconv.MustInt(dstClaims["aud"]))
	// 如果设置了续期参数
	if exp <= request.CheckExpireTime {
		return s.renewAccessToken(request, aud, exp), nil
	}
	return &proto.CheckAccessTokenResponse{
		MemberId:         aud,
		TokenExpiresTime: exp,
	}, nil
}

// renewAccessToken 续签令牌
func (s *memberService) renewAccessToken(request *proto.CheckAccessTokenRequest,
	aud int64, exp int64) *proto.CheckAccessTokenResponse {
	if request.RenewExpiresTime < request.CheckExpireTime {
		return &proto.CheckAccessTokenResponse{
			Error: "令牌续期过期时间必须在检测过期时间之后",
		}
	}
	ret, _ := s.GrantAccessToken(context.TODO(), &proto.GrantAccessTokenRequest{
		MemberId:    aud,
		ExpiresTime: request.RenewExpiresTime,
	})
	if len(ret.Error) > 0 {
		return &proto.CheckAccessTokenResponse{
			Error: ret.Error,
		}
	}
	return &proto.CheckAccessTokenResponse{
		IsExpires:        false,
		TokenExpiresTime: exp,
		MemberId:         aud,
		RenewAccessToken: ret.AccessToken,
	}
}

// VerifyTradePassword 检查交易密码
func (s *memberService) VerifyTradePassword(_ context.Context, r *proto.VerifyPasswordRequest) (*proto.Result, error) {
	m := s.repo.GetMember(r.MemberId)
	if m == nil {
		return s.result(member.ErrNoSuchMember), nil
	}
	mv := m.GetValue()
	if mv.TradePassword == "" {
		return s.error(member.ErrNotSetTradePassword), nil
	}
	if len(r.Password) != 32 {
		return s.error(de.ErrNotMD5Format), nil
	}
	if encPwd := domain.TradePassword(r.Password, mv.Salt); mv.TradePassword != encPwd {
		return s.error(member.ErrIncorrectTradePassword), nil
	}
	return s.success(nil), nil
}

// 检查与现有用户不同的用户是否存在,如存在则返回错误
//func (_s *memberService) CheckUser(user string, memberId int64) error {
//	if len(user) < 6 {
//		return member.ErrUserLength
//	}
//	if _s.repo.CheckUserExist(user, memberId) {
//		return member.ErrUserExist
//	}
//	return nil
//}

// GetAccount 获取会员账户
func (s *memberService) GetAccount(_ context.Context, id *proto.MemberIdRequest) (*proto.SAccount, error) {
	m := s.repo.CreateMember(&member.Member{Id: id.MemberId})
	acc := m.GetAccount()
	if acc != nil {
		return s.parseAccountDto(acc.GetValue()), nil
	}
	return nil, member.ErrNoSuchMember
}

// 获取上级邀请人会员编号数组
func (s *memberService) InviterArray(_ context.Context, r *proto.DepthRequest) (*proto.InviterIdListResponse, error) {
	m := s.repo.CreateMember(&member.Member{Id: r.MemberId})
	var arr []int64
	if m != nil {
		arr = m.Invitation().InviterArray(r.MemberId, int(r.Depth))
	}
	return &proto.InviterIdListResponse{
		Value: arr,
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
		Value: s.query.GetInviteArray(r.MemberId, where),
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

// 解锁银行卡信息
func (s *memberService) RemoveBankCard(_ context.Context, r *proto.BankCardRequest) (*proto.Result, error) {
	m := s.repo.CreateMember(&member.Member{Id: r.MemberId})
	err := m.Profile().RemoveBankCard(r.BankCardNo)
	return s.result(err), nil
}

// 获取收款码
func (s *memberService) ReceiptsCodes(_ context.Context, id *proto.MemberIdRequest) (*proto.SReceiptsCodeListResponse, error) {
	m := s.repo.GetMember(id.MemberId)
	if m == nil {
		return &proto.SReceiptsCodeListResponse{
			Value: make([]*proto.SReceiptsCode, 0),
		}, nil
	}
	arr := m.Profile().ReceiptsCodes()
	list := make([]*proto.SReceiptsCode, len(arr))
	for i, v := range arr {
		list[i] = &proto.SReceiptsCode{
			Identity:       v.Identity,
			ReceipterName:  v.Name,
			ReceiptAccount: v.AccountId,
			CodeImageUrl:   v.CodeUrl,
			State:          int32(v.State),
		}
	}
	return &proto.SReceiptsCodeListResponse{Value: list}, nil
}

// 保存收款码
func (s *memberService) SaveReceiptsCode(_ context.Context, r *proto.ReceiptsCodeSaveRequest) (*proto.Result, error) {
	m := s.repo.GetMember(r.MemberId)
	if m == nil {
		return s.error(member.ErrNoSuchMember), nil
	}
	v := &member.ReceiptsCode{
		Identity:  r.Code.Identity,
		Name:      r.Code.ReceipterName,
		AccountId: r.Code.ReceiptAccount,
		CodeUrl:   r.Code.CodeImageUrl,
		State:     int(r.Code.State),
	}
	if err := m.Profile().SaveReceiptsCode(v); err != nil {
		return s.error(err), nil
	}
	return s.success(nil), nil
}

// 获取银行卡
func (s *memberService) GetBankCards(_ context.Context, id *proto.MemberIdRequest) (*proto.BankCardListResponse, error) {
	m := s.repo.CreateMember(&member.Member{Id: id.MemberId})
	b := m.Profile().GetBankCards()
	arr := make([]*proto.SBankCardInfo, len(b))
	for i, v := range b {
		arr[i] = &proto.SBankCardInfo{
			BankName:    v.BankName,
			AccountName: v.AccountName,
			AccountNo:   v.BankAccount,
			BankId:      int32(v.BankId),
			BankCode:    v.BankCode,
			AuthCode:    v.AuthCode,
			Network:     v.Network,
			State:       int32(v.State),
			UpdateTime:  v.CreateTime,
		}
	}
	return &proto.BankCardListResponse{
		Value: arr,
	}, nil
}

// 保存银行卡
func (s *memberService) AddBankCard(_ context.Context, r *proto.BankCardAddRequest) (*proto.Result, error) {
	m := s.repo.CreateMember(&member.Member{Id: r.MemberId})
	var v = &member.BankCard{
		MemberId:    r.MemberId,
		BankAccount: r.Value.AccountNo,
		AccountName: r.Value.AccountName,
		BankId:      int(r.Value.BankId),
		BankName:    r.Value.BankName,
		BankCode:    r.Value.BankCode,
		Network:     r.Value.Network,
		AuthCode:    r.Value.AuthCode,
		State:       int16(r.Value.State),
	}
	err := m.Profile().AddBankCard(v)
	return s.result(err), nil
}

// 实名认证信息
func (s *memberService) GetTrustInfo(_ context.Context, id *proto.MemberIdRequest) (*proto.STrustedInfo, error) {
	t := &member.TrustedInfo{}
	m := s.repo.GetMember(id.MemberId)
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

/*********** 收货地址 ***********/

// GetAddressList 获取会员的收货地址
func (s *memberService) GetAddressList(_ context.Context, id *proto.MemberIdRequest) (*proto.AddressListResponse, error) {
	src := s.repo.GetDeliverAddress(id.MemberId)
	var arr []*proto.SAddress
	for _, v := range src {
		arr = append(arr, s.parseAddressDto(v))
	}
	return &proto.AddressListResponse{Value: arr}, nil
}

// GetAddress 获取配送地址
func (s *memberService) GetAddress(_ context.Context, r *proto.GetAddressRequest) (*proto.SAddress, error) {
	m := s.repo.CreateMember(&member.Member{Id: r.MemberId})
	pro := m.Profile()
	var addr member.IDeliverAddress
	if r.AddressId > 0 {
		addr = pro.GetAddress(r.AddressId)
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
	return nil, member.ErrNoSuchAddress
}

// SaveAddress 保存配送地址
func (s *memberService) SaveAddress(_ context.Context, r *proto.SaveAddressRequest) (*proto.SaveAddressResponse, error) {
	e := s.parseAddress(r.Value)
	e.MemberId = r.MemberId
	if r.MemberId <= 0 {
		return &proto.SaveAddressResponse{ErrCode: 1, ErrMsg: member.ErrNoSuchMember.Error()}, nil
	}
	m := s.repo.CreateMember(&member.Member{Id: r.MemberId})
	var v member.IDeliverAddress
	ret := &proto.SaveAddressResponse{}
	if e.Id > 0 {
		v = m.Profile().GetAddress(e.Id)
	} else {
		v = m.Profile().CreateDeliver(e)
	}
	err := v.SetValue(e)
	if err == nil {
		ret.AddressId, err = v.Save()
		// 设置默认收货地址
		if e.IsDefault == 1 && err == nil {
			err = m.Profile().SetDefaultAddress(v.GetDomainId())
		}
	}
	if err != nil {
		ret.ErrCode = 1
		ret.ErrMsg = err.Error()
	}
	return ret, nil
}

// DeleteAddress 删除配送地址
func (s *memberService) DeleteAddress(_ context.Context, r *proto.AddressIdRequest) (*proto.Result, error) {
	m := s.repo.CreateMember(&member.Member{Id: r.MemberId})
	err := m.Profile().DeleteAddress(r.AddressId)
	if err != nil {
		return s.error(err), nil
	}
	return s.success(nil), nil
}

// SetPayPriority 设置余额优先支付
func (s *memberService) SetPayPriority(_ context.Context, r *proto.PayPriorityRequest) (*proto.Result, error) {
	m := s.repo.GetMember(r.OwnerId)
	if m == nil {
		return s.result(member.ErrNoSuchMember), nil
	}
	var accountTid member.AccountType
	switch r.Account {
	case proto.PaymentAccountType_PA_BALANCE:
		accountTid = member.AccountBalance
	case proto.PaymentAccountType_PA_WALLET:
		accountTid = member.AccountWallet
	case proto.PaymentAccountType_PA_QUICK_PAY:
		return s.error(errors.New("暂时不支持")), nil
	}
	err := m.GetAccount().SetPriorityPay(accountTid, true)
	return s.error(err), nil
}

// IsInvitation 判断会员是否由指定会员邀请推荐的
func (s *memberService) IsInvitation(c context.Context, r *proto.IsInvitationRequest) (*proto.Bool, error) {
	m := s.repo.CreateMember(&member.Member{Id: r.MemberId})
	b := m.Invitation().InvitationBy(r.InviterId)
	return &proto.Bool{Value: b}, nil
}

// GetMyPagedInvitationMembers 获取我邀请的会员及会员邀请的人数
func (s *memberService) GetMyPagedInvitationMembers(_ context.Context, r *proto.MemberInvitationPagingRequest) (*proto.MemberInvitationPagingResponse, error) {
	iv := s.repo.CreateMember(&member.Member{Id: r.MemberId}).Invitation()
	total, rows := iv.GetInvitationMembers(int(r.Begin), int(r.End))
	ret := &proto.MemberInvitationPagingResponse{
		Total: int64(total),
		Data:  make([]*proto.SInvitationMember, 0),
	}
	if l := len(rows); l > 0 {
		arr := make([]int32, l)
		for i := 0; i < l; i++ {
			arr[i] = rows[i].MemberId
		}
		num := iv.GetSubInvitationNum(arr)
		for i := 0; i < l; i++ {
			rows[i].InvitationNum = num[rows[i].MemberId]
			rows[i].Portrait = format.GetFileFullUrl(rows[i].Portrait)
			ret.Data = append(ret.Data, &proto.SInvitationMember{
				MemberId: int64(rows[i].MemberId),
				Username: rows[i].Username,
				Level:    rows[i].Level,
				Portrait: rows[i].Portrait,
				Nickname: rows[i].Nickname,
				Phone:    rows[i].Phone,
				RegTime:  rows[i].RegTime,
				//Im:            rows[i].Im,
				InvitationNum: int32(rows[i].InvitationNum),
			})
		}
	}
	return ret, nil
}

// GetMemberLatestUpdateTime 获取会员最后更新时间
func (s *memberService) GetMemberLatestUpdateTime(memberId int64) int64 {
	return s.repo.GetMemberLatestUpdateTime(memberId)
}

// GrantFlag 标志赋值, 如果flag小于零, 则异或运算
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

// Complex 获取会员汇总信息
func (s *memberService) Complex(_ context.Context, id *proto.MemberIdRequest) (*proto.SComplexMember, error) {
	m := s.repo.GetMember(id.MemberId)
	if m != nil {
		x := m.Complex()
		return s.parseComplexMemberDto(x), nil
	}
	return nil, member.ErrNoSuchMember
}

func (s *memberService) Freeze(_ context.Context, r *proto.AccountFreezeRequest) (*proto.AccountFreezeResponse, error) {
	m := s.repo.GetMember(r.MemberId)
	if m == nil {
		return &proto.AccountFreezeResponse{ErrCode: 1, ErrMsg: member.ErrNoSuchMember.Error()}, nil
	}
	id, err := m.GetAccount().Freeze(member.AccountType(r.AccountType),
		member.AccountOperateData{
			Title:   r.Title,
			Amount:  int(r.Amount),
			OuterNo: r.OuterNo,
			Remark:  r.Remark,
		}, 0)
	if err != nil {
		return &proto.AccountFreezeResponse{ErrCode: 1, ErrMsg: err.Error()}, nil
	}
	return &proto.AccountFreezeResponse{LogId: int64(id)}, nil
}

func (s *memberService) Unfreeze(_ context.Context, r *proto.AccountUnfreezeRequest) (*proto.Result, error) {
	m := s.repo.GetMember(r.MemberId)
	if m == nil {
		return s.error(member.ErrNoSuchMember), nil
	}
	err := m.GetAccount().Unfreeze(member.AccountType(r.AccountType),
		member.AccountOperateData{
			Title:   r.Title,
			Amount:  int(r.Amount),
			OuterNo: r.OuterNo,
			Remark:  r.Remark,
		}, 0)
	if err != nil {
		return s.error(err), nil
	}
	return s.success(nil), nil
}

// AccountCharge 充值,account为账户类型,kind为业务类型
func (s *memberService) AccountCharge(_ context.Context, r *proto.AccountChangeRequest) (*proto.Result, error) {
	var err error
	m := s.repo.CreateMember(&member.Member{Id: r.MemberId})
	acc := m.GetAccount()
	if acc == nil {
		err = member.ErrNoSuchMember
	} else {
		err = acc.Charge(member.AccountType(r.AccountType), r.Title, int(r.Amount), r.OuterNo, r.Remark)
	}
	return s.result(err), nil
}

// AccountCarryTo 账户入账
func (s *memberService) AccountCarryTo(_ context.Context, r *proto.AccountCarryRequest) (*proto.AccountCarryResponse, error) {
	m := s.repo.CreateMember(&member.Member{Id: r.MemberId})
	if m == nil {
		return &proto.AccountCarryResponse{
			ErrCode: 1,
			ErrMsg:  member.ErrNoSuchMember.Error(),
		}, nil
	}
	acc := m.GetAccount()
	if acc == nil {
		return &proto.AccountCarryResponse{
			ErrCode: 1,
			ErrMsg:  member.ErrNoSuchMember.Error(),
		}, nil
	}
	id, err := acc.CarryTo(member.AccountType(r.AccountType),
		member.AccountOperateData{
			Title:   r.Title,
			Amount:  int(r.Amount),
			OuterNo: r.OuterNo,
			Remark:  r.Remark,
		}, r.Freeze, int(r.ProcedureFee))
	if err != nil {
		return &proto.AccountCarryResponse{
			ErrCode: 1,
			ErrMsg:  err.Error(),
		}, nil
	}
	return &proto.AccountCarryResponse{
		LogId: int64(id),
	}, nil
}

// AccountDiscount 账户抵扣
func (s *memberService) AccountDiscount(_ context.Context, r *proto.AccountChangeRequest) (*proto.Result, error) {
	m, err := s.getMember(r.MemberId)
	if err == nil {
		acc := m.GetAccount()
		err = acc.Discount(member.AccountType(r.AccountType), r.Title, int(r.Amount), r.OuterNo, r.Remark)
	}
	return s.result(err), nil
}

// AccountConsume 账户消耗
func (s *memberService) AccountConsume(_ context.Context, r *proto.AccountChangeRequest) (*proto.Result, error) {
	m, err := s.getMember(r.MemberId)
	if err == nil {
		acc := m.GetAccount()
		err = acc.Consume(member.AccountType(r.AccountType), r.Title, int(r.Amount), r.OuterNo, r.Remark)
	}
	return s.result(err), nil
}

// AccountRefund 账户退款
func (s *memberService) AccountRefund(_ context.Context, r *proto.AccountChangeRequest) (*proto.Result, error) {
	m, err := s.getMember(r.MemberId)
	if err == nil {
		acc := m.GetAccount()
		err = acc.Refund(member.AccountType(r.AccountType), r.Title, int(r.Amount), r.OuterNo, r.Remark)
	}
	return s.result(err), nil
}

// AccountAdjust 调整账户
func (s *memberService) AccountAdjust(_ context.Context, r *proto.AccountAdjustRequest) (*proto.Result, error) {
	m, err := s.getMember(r.MemberId)
	if err == nil {
		tit := "系统冲正"
		// 人工冲正带[KF]字样
		if r.ManualAdjust {
			tit = "[KF]系统冲正"
		}
		acc := m.GetAccount()
		err = acc.Adjust(member.AccountType(r.Account), tit, int(r.Value), r.Remark, r.RelateUser)
	}
	return s.result(err), nil
}

// B4EAuth !银行四要素认证
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

// Withdraw 提现并返回提现编号,交易号以及错误信息
func (s *memberService) Withdraw(_ context.Context, r *proto.WithdrawRequest) (*proto.WithdrawalResponse, error) {
	m, err := s.getMember(r.MemberId)
	if err != nil {
		return &proto.WithdrawalResponse{ErrCode: 1, ErrMsg: err.Error()}, nil
	}
	title := ""
	kind := 0
	switch int(r.WithdrawalKind) {
	case int(proto.EWithdrawalKind_WithdrawToBankCard):
		title = "提现到银行卡"
		kind = wallet.KWithdrawToBankCard
	case int(proto.EWithdrawalKind_WithdrawToThirdPart):
		title = "充值到第三方账户"
		kind = wallet.KWithdrawToThirdPart
	case int(proto.EWithdrawalKind_WithdrawByExchange):
		title = "提现到余额"
		kind = wallet.KWithdrawExchange
	}
	acc := m.GetAccount()
	_, tradeNo, err := acc.RequestWithdrawal(kind, title,
		int(r.Amount), int(r.ProcedureFee), r.AccountNo)
	if err != nil {
		return &proto.WithdrawalResponse{ErrCode: 1, ErrMsg: err.Error()}, nil
	}
	return &proto.WithdrawalResponse{
		ErrCode: 0,
		ErrMsg:  "",
		TradeNo: tradeNo,
	}, nil
}

func (s *memberService) QueryWithdrawalLog(_ context.Context, r *proto.WithdrawalLogRequest) (*proto.WithdrawalLogResponse, error) {
	//todo: 这里只返回了一条
	latestApplyInfo := s.query.GetLatestWalletLogByKind(r.MemberId,
		wallet.KWithdrawToBankCard)
	//if latestApplyInfo != nil {
	//	var sText string
	//	switch latestApplyInfo.ReviewState {
	//	case enum.ReviewAwaiting:
	//		sText = "已申请"
	//	case enum.ReviewPass:
	//		sText = "已审核,等待打款"
	//	case enum.ReviewReject:
	//		sText = "被退回"
	//	case enum.ReviewConfirm:
	//		sText = "已完成"
	//	}
	//	if latestApplyInfo.Amount < 0 {
	//		latestApplyInfo.Amount = -latestApplyInfo.Amount
	//	}
	//	latestInfo := fmt.Sprintf(`<b>最近提现：</b>%_s&nbsp;申请提现%_s ，状态：<span class="status">%_s</span>。`,
	//		time.Unix(latestApplyInfo.CreateTime, 0).Format("2006-01-02 15:04"),
	//		format.FormatFloat(latestApplyInfo.Amount),
	//		sText)
	//}
	ret := &proto.WithdrawalLogResponse{Data: make([]*proto.WithdrawalLog, 0)}
	if latestApplyInfo != nil {
		ret.Data = append(ret.Data, &proto.WithdrawalLog{
			Id:           latestApplyInfo.Id,
			OuterNo:      latestApplyInfo.OuterNo,
			Kind:         int32(latestApplyInfo.Kind),
			Title:        latestApplyInfo.Title,
			Amount:       latestApplyInfo.Amount,
			ProcedureFee: latestApplyInfo.ProcedureFee,
			RelateUser:   latestApplyInfo.RelateUser,
			ReviewState:  latestApplyInfo.ReviewState,
			Remark:       latestApplyInfo.Remark,
			SubmitTime:   latestApplyInfo.CreateTime,
			UpdateTime:   latestApplyInfo.UpdateTime,
		})
	}
	return ret, nil
}

// ReviewWithdrawal 确认提现
func (s *memberService) ReviewWithdrawal(_ context.Context, r *proto.ReviewWithdrawalRequest) (*proto.Result, error) {
	m, err := s.getMember(r.MemberId)
	if err == nil {
		err = m.GetAccount().ReviewWithdrawal(r.LogId, r.Pass, r.Remark)
	}
	return s.error(err), nil
}

// 完成提现
func (s *memberService) FinishWithdrawal(_ context.Context, r *proto.FinishWithdrawalRequest) (*proto.Result, error) {
	var err error
	m, err := s.getMember(r.MemberId)
	if err == nil {
		err = m.GetAccount().FinishWithdrawal(r.InfoId, r.TradeNo)
	}
	return s.error(err), nil
}

// AccountTransfer 转账余额到其他账户
func (s *memberService) AccountTransfer(_ context.Context, r *proto.AccountTransferRequest) (*proto.Result, error) {
	var err error
	m := s.repo.GetMember(r.FromMemberId)
	if m == nil {
		err = member.ErrNoSuchMember
	} else {
		var account member.AccountType
		switch r.TransferAccount {
		case proto.EAccountType_AccountBalance:
			account = member.AccountBalance
		case proto.EAccountType_AccountWallet:
			account = member.AccountWallet
		case proto.EAccountType_AccountIntegral:
			account = member.AccountIntegral
		}
		err = m.GetAccount().TransferAccount(account, r.ToMemberId,
			int(r.Amount), int(r.ProcedureFee), r.Remark)
	}
	return s.error(err), nil
}

// GetMemberInviRank 会员推广排名
func (s *memberService) GetMemberInviRank(mchId int64, allTeam bool,
	levelComp string, level int, startTime int64, endTime int64,
	num int) []*dto.RankMember {
	return s.query.GetMemberInviRank(mchId, allTeam, levelComp, level, startTime, endTime, num)
}

//********* 促销  **********//

// QueryCoupons 查询优惠券
func (s *memberService) QueryCoupons(_ context.Context, r *proto.MemberCouponPagingRequest) (*proto.MemberCouponListResponse, error) {
	cp := s.repo.CreateMemberById(r.MemberId).GiftCard()
	begin, end := int(r.Begin), int(r.End)
	var total int
	var list []*dto.SimpleCoupon
	switch r.State {
	case proto.PagingCouponState_CS_AVAILABLE:
		total, list = cp.PagedAvailableCoupon(begin, end)
	case proto.PagingCouponState_CS_EXPIRED:
		total, list = cp.PagedExpiresCoupon(begin, end)
	default:
		total, list = cp.PagedAllCoupon(begin, end)
	}
	ret := &proto.MemberCouponListResponse{
		Total: int64(total),
		Data:  make([]*proto.SMemberCoupon, total),
	}
	for i, v := range list {
		ret.Data[i] = &proto.SMemberCoupon{
			CouponId:    int64(v.Id),
			Number:      int32(v.Num),
			Title:       v.Title,
			Code:        v.Code,
			DiscountFee: int32(v.Fee),
			Discount:    int32(v.Discount),
			IsUsed:      v.IsUsed == 1,
			GetTime:     0, //todo: ???
			OverTime:    v.OverTime,
		}
	}
	return ret, nil
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
		Id:            int32(src.Id),
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
		Username:       src.Username,
		UserCode:       src.UserCode,
		Exp:            int64(src.Exp),
		Level:          int32(src.Level),
		PremiumUser:    int32(src.PremiumUser),
		PremiumExpires: src.PremiumExpires,
		RegIp:          src.RegIp,
		RegFrom:        src.RegFrom,
		State:          int32(src.State),
		UserFlag:       int32(src.UserFlag),
		Portrait:       src.Portrait,
		Phone:          src.Phone,
		Email:          src.Email,
		Nickname:       src.Nickname,
		RealName:       src.RealName,
		RegTime:        src.RegTime,
		LastLoginTime:  src.LastLoginTime,
	}
}

func (s *memberService) parseMemberProfile(src *member.Profile) *proto.SProfile {
	return &proto.SProfile{
		MemberId:   src.MemberId,
		Nickname:   src.Name,
		Portrait:   src.Avatar,
		Gender:     src.Gender,
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
		Nickname:            src.Nickname,
		Portrait:            src.Avatar,
		Exp:                 int32(src.Exp),
		Level:               int32(src.Level),
		LevelName:           src.LevelName,
		PremiumUser:         int32(src.PremiumUser),
		TrustAuthState:      int32(src.TrustAuthState),
		TradePasswordHasSet: src.TradePasswordHasSet,
		UpdateTime:          src.UpdateTime,
	}
}

func (s *memberService) parseAddressDto(src *member.ConsigneeAddress) *proto.SAddress {
	return &proto.SAddress{
		AddressId:      src.Id,
		ConsigneeName:  src.ConsigneeName,
		ConsigneePhone: src.ConsigneePhone,
		Province:       src.Province,
		City:           src.City,
		District:       src.District,
		Area:           src.Area,
		DetailAddress:  src.DetailAddress,
		IsDefault:      src.IsDefault == 1,
	}
}

func (s *memberService) parseAccountDto(src *member.Account) *proto.SAccount {
	return &proto.SAccount{
		Integral:            int64(src.Integral),
		FreezeIntegral:      int64(src.FreezeIntegral),
		Balance:             src.Balance,
		FreezeBalance:       src.FreezeBalance,
		ExpiredBalance:      src.ExpiredBalance,
		WalletBalance:       src.WalletBalance,
		WalletCode:          src.WalletCode,
		WalletFreezedAmount: src.FreezeWallet,
		WalletExpiredAmount: src.ExpiredWallet,
		TotalWalletAmount:   src.TotalWalletAmount,
		FlowBalance:         src.FlowBalance,
		GrowBalance:         src.GrowBalance,
		GrowAmount:          src.GrowAmount,
		GrowEarnings:        src.GrowEarnings,
		GrowTotalEarnings:   src.GrowTotalEarnings,
		TotalExpense:        src.TotalExpense,
		TotalCharge:         src.TotalCharge,
		TotalPay:            src.TotalPay,
		PriorityPay:         int32(src.PriorityPay),
		UpdateTime:          src.UpdateTime,
	}
}

func (s *memberService) parseMemberProfile2(src *proto.SProfile) *member.Profile {
	return &member.Profile{
		MemberId:   src.MemberId,
		Name:       src.Nickname,
		Avatar:     src.Portrait,
		Gender:     src.Gender,
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

func (s *memberService) parseAddress(src *proto.SAddress) *member.ConsigneeAddress {
	return &member.ConsigneeAddress{
		Id:             src.AddressId,
		ConsigneeName:  src.ConsigneeName,
		ConsigneePhone: src.ConsigneePhone,
		Province:       src.Province,
		City:           src.City,
		District:       src.District,
		Area:           src.Area,
		DetailAddress:  src.DetailAddress,
		IsDefault:      types.ElseInt(src.IsDefault, 1, 0),
	}
}

// BindOAuthApp 绑定第三方应用
func (m *memberService) BindOAuthApp(_ context.Context, req *proto.SMemberOAuthAccount) (*proto.Result, error) {
	mm := m.repo.GetMember(req.MemberId)
	if mm == nil {
		return m.error(member.ErrNoSuchMember), nil
	}
	err := mm.Profile().BindOAuthApp(req.AppCode, req.OpenId, req.AuthToken)
	return m.error(err), nil
}

// GetOAuthBindInfo 获取第三方应用绑定信息
func (m *memberService) GetOAuthBindInfo(_ context.Context, req *proto.MemberOAuthRequest) (*proto.SMemberOAuthAccount, error) {
	mm := m.repo.GetMember(req.MemberId)
	if mm == nil {
		return &proto.SMemberOAuthAccount{}, nil
	}
	bind := mm.Profile().GetOAuthBindInfo(req.AppCode)
	if bind == nil {
		return &proto.SMemberOAuthAccount{}, nil
	}
	return &proto.SMemberOAuthAccount{
		MemberId:    req.MemberId,
		AppCode:     req.AppCode,
		OpenId:      bind.OpenId,
		AuthToken:   bind.AuthToken,
		PortraitUrl: bind.HeadImgUrl,
	}, nil
}

// UnbindOAuthApp 解除第三方应用绑定
func (m *memberService) UnbindOAuthApp(_ context.Context, req *proto.MemberOAuthRequest) (*proto.Result, error) {
	mm := m.repo.GetMember(req.MemberId)
	if mm == nil {
		return m.error(member.ErrNoSuchMember), nil
	}
	err := mm.Profile().UnbindOAuthApp(req.AppCode)
	return m.error(err), nil
}
