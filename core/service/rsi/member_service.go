package rsi

/**
 * Copyright 2014 @ to2.net.
 * name :
 * author : jarryliu
 * date : 2013-12-09 20:14
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
	"github.com/ixre/gof/util"
	"go2o/core/domain/interface/enum"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/mss/notify"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/dto"
	"go2o/core/infrastructure/domain"
	"go2o/core/infrastructure/format"
	"go2o/core/module"
	"go2o/core/query"
	"go2o/core/service/auto_gen/rpc/member_service"
	"go2o/core/service/auto_gen/rpc/order_service"
	"go2o/core/service/auto_gen/rpc/ttype"
	"go2o/core/service/thrift/parser"
	"go2o/core/variable"
	"log"
	"strconv"
	"strings"
	"time"
)

var _ member_service.MemberService = new(memberService)

type memberService struct {
	repo       member.IMemberRepo
	mchService *merchantService
	query      *query.MemberQuery
	orderQuery *query.OrderQuery
	valRepo    valueobject.IValueRepo
	serviceUtil
}

func NewMemberService(mchService *merchantService, repo member.IMemberRepo,
	q *query.MemberQuery, oq *query.OrderQuery, valRepo valueobject.IValueRepo) *memberService {
	ms := &memberService{
		repo:       repo,
		query:      q,
		mchService: mchService,
		orderQuery: oq,
		valRepo:    valRepo,
	}
	return ms
	//return m.init()
}

func (s *memberService) init() *memberService {
	db := gof.CurrentApp.Db()
	var list []*member.InviteRelation
	db.GetOrm().Select(&list, "")
	//for _, v := range list {
	//	s.repo.GetMember(v.MemberId).saveRelation(v)
	//}
	return s
}

// 根据会员编号获取会员
func (s *memberService) getValueMember(memberId int64) *member.Member {
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
func (s *memberService) GetMember(ctx context.Context, id int64) (*member_service.SMember, error) {
	v := s.getValueMember(id)
	if v != nil {
		return parser.MemberDto(v), nil
	}
	return nil, nil
}

// 根据用户名获取会员
func (s *memberService) GetMemberByUser(ctx context.Context, user string) (*member_service.SMember, error) {
	v := s.repo.GetMemberByUser(user)
	if v != nil {
		return parser.MemberDto(v), nil
	}
	return nil, nil
}

// 获取资料
func (s *memberService) GetProfile(ctx context.Context, memberId int64) (*member_service.SProfile, error) {
	m := s.repo.GetMember(memberId)
	if m != nil {
		v := m.Profile().GetProfile()
		return parser.MemberProfile(&v), nil
	}
	return nil, nil
}

// 保存资料
func (s *memberService) SaveProfile(v *member_service.SProfile) error {
	if v.MemberId > 0 {
		v2 := parser.MemberProfile2(v)
		m := s.repo.GetMember(v.MemberId)
		if m == nil {
			return member.ErrNoSuchMember
		}
		return m.Profile().SaveProfile(v2)
	}
	return nil
}

// 升级为高级会员
func (s *memberService) Premium(ctx context.Context, memberId int64, v int32, expires int64) (*ttype.Result_, error) {
	m := s.repo.GetMember(memberId)
	if m == nil {
		return s.result(member.ErrNoSuchMember), nil
	}
	err := m.Premium(int(v), expires)
	return s.result(err), nil
}

// 检查会员的会话Token是否正确
func (s *memberService) CheckToken(ctx context.Context, memberId int64, token string) (r bool, err error) {
	md := module.Get(module.M_MM).(*module.MemberModule)
	return md.CheckToken(memberId, token), nil
}

// 获取会员的会员Token,reset表示是否重置会员的token
func (s *memberService) GetToken(ctx context.Context, memberId int64, reset bool) (r string, err error) {
	pubToken := ""
	md := module.Get(module.M_MM).(*module.MemberModule)
	if !reset {
		pubToken = md.GetToken(memberId)
	}
	if reset || (pubToken == "" && memberId > 0) {
		m := s.getValueMember(memberId)
		if m != nil {
			return md.ResetToken(memberId, m.Pwd), nil
		}
	}
	return pubToken, nil
}

// 移除会员的Token
func (s *memberService) RemoveToken(ctx context.Context, memberId int64) (err error) {
	md := module.Get(module.M_MM).(*module.MemberModule)
	md.RemoveToken(memberId)
	return nil
}

// 更改手机号码，不验证手机格式
func (s *memberService) ChangePhone(ctx context.Context, memberId int64, phone string) (*ttype.Result_, error) {
	err := s.changePhone(memberId, phone)
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

/**================ 会员等级 ==================**/
// 获取所有会员等级
func (s *memberService) GetMemberLevels() []*member.Level {
	return s.repo.GetManager().LevelManager().GetLevelSet()
}

// 等级列表
func (s *memberService) LevelList(ctx context.Context) ([]*member_service.SLevel, error) {
	var arr []*member_service.SLevel
	list := s.repo.GetManager().LevelManager().GetLevelSet()
	for _, v := range list {
		arr = append(arr, parser.LevelDto(v))
	}
	return arr, nil
}

// 根据编号获取会员等级信息
func (s *memberService) GetLevel(ctx context.Context, id int32) (*member_service.SLevel, error) {
	lv := s.repo.GetManager().LevelManager().GetLevelById(int(id))
	if lv != nil {
		return parser.LevelDto(lv), nil
	}
	return nil, nil
}

// 根据SIGN获取等级
func (s *memberService) GetLevelBySign(ctx context.Context, sign string) (*member_service.SLevel, error) {
	lv := s.repo.GetManager().LevelManager().GetLevelByProgramSign(sign)
	if lv != nil {
		return parser.LevelDto(lv), nil
	}
	return nil, nil
}

// 根据可编程字符获取会员等级
func (s *memberService) GetLevelByProgramSign(sign string) *member.Level {
	return s.repo.GetManager().LevelManager().GetLevelByProgramSign(sign)
}

// 保存会员等级信息
func (s *memberService) SaveMemberLevel(v *member.Level) (int32, error) {
	n, err := s.repo.GetManager().LevelManager().SaveLevel(v)
	return int32(n), err
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

func (s *memberService) GetWalletLog(memberId int64, logId int32) *member.MWalletLog {
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

func (s *memberService) GetMemberIdByInvitationCode(code string) int64 {
	return s.repo.GetMemberIdByInvitationCode(code)
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

// 发送验证码
func (s *memberService) SendCode(memberId int64, operation string, msgType int) (string, error) {
	m := s.repo.GetMember(memberId)
	if m == nil {
		return "", member.ErrNoSuchMember
	}
	return m.SendCheckCode(operation, msgType)
}

// 对比验证码
func (s *memberService) CompareCode(memberId int64, code string) error {
	m := s.repo.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.CompareCode(code)
}

// 更改会员用户名
func (s *memberService) ChangeUsr(ctx context.Context, memberId int64, user string) (*ttype.Result_, error) {
	err := s.changeUsr(int(memberId), user)
	return s.result(err), nil
}

// 更改会员等级
func (s *memberService) UpdateLevel(ctx context.Context, memberId int64, level int32,
	review bool, paymentOrderId int64) (r *ttype.Result_, err error) {
	m := s.repo.GetMember(memberId)
	if m == nil {
		err = member.ErrNoSuchMember
	} else {
		err = m.ChangeLevel(int(level), int(paymentOrderId), review)
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
func (s *memberService) SaveMember(v *member_service.SMember) (int64, error) {
	if v.ID > 0 {
		return s.updateMember(v)
	}
	return -1, errors.New("Create member use \"RegisterMember\" method.")
}

func (s *memberService) updateMember(v *member_service.SMember) (int64, error) {
	m := s.repo.GetMember(int64(v.ID))
	if m == nil {
		return -1, member.ErrNoSuchMember
	}
	mv := parser.Member(v)
	if err := m.SetValue(mv); err != nil {
		return m.GetAggregateRootId(), err
	}
	return m.Save()
}

// 注册会员
func (s *memberService) RegisterMemberV2(ctx context.Context, user string, pwd string,
	flag int32, name string, phone string, email string, avatar string,
	extend map[string]string) (r *ttype.Result_, err error) {
	inviteCode := extend["invite_code"]
	inviterId, err := s.repo.GetManager().CheckInviteRegister(inviteCode)
	if err != nil {
		return s.error(err), nil
	}
	v := &member.Member{
		User:    user,
		Pwd:     domain.Sha1Pwd(pwd),
		Name:    name,
		Avatar:  avatar,
		Phone:   phone,
		Email:   email,
		RegFrom: extend["reg_from"],
		RegIp:   extend["reg_ip"],
		Flag:    int(flag),
	}
	m := s.repo.CreateMember(v) //创建会员
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
		//	} else {
		//		// 保存关联信息
		//		rl := m.GetRelation()
		//		rl.InviterId = invitationId
		//		rl.RegMchId = mchId
		//		rl.CardCard = cardId
		//		err = m.saveRelation(rl)
		//	}
		//}
		return s.success(map[string]string{
			"member_id": util.Str(id),
		}), nil
	}
	return s.error(err), nil
}

// 注册会员
func (s *memberService) RegisterMember1(mchId int32, v1 *member_service.SMember,
	pro1 *member_service.SProfile, cardId string, invitationCode string) (int64, error) {
	if v1 == nil || pro1 == nil {
		return 0, errors.New("missing data")
	}
	v := parser.Member(v1)
	pro := parser.MemberProfile2(pro1)
	invitationId, err := s.repo.GetManager().PrepareRegister(
		v, pro, invitationCode)
	if err == nil {
		m := s.repo.CreateMember(v) //创建会员
		id, err := m.Save()
		if err == nil {
			pro.Sex = 1
			pro.MemberId = id
			//todo: 如果注册失败，则删除。应使用SQL-TRANSFER
			if err = m.Profile().SaveProfile(pro); err != nil {
				s.repo.DeleteMember(id)
			} else {
				// 保存关联信息
				err = m.BindInviter(invitationId, true)
			}
		}
		return id, err
	}
	return -1, err
}

// 获取会员等级
func (s *memberService) GetMemberLevel(memberId int64) *member.Level {
	m := s.repo.GetMember(memberId)
	if m == nil {
		return nil
	}
	return m.GetLevel()
}

func (s *memberService) GetRelation(memberId int64) *member.InviteRelation {
	return s.repo.GetRelation(memberId)
}

// 激活会员
func (s *memberService) Active(ctx context.Context, memberId int64) (r *ttype.Result_, err error) {
	m := s.repo.GetMember(memberId)
	if m == nil {
		return s.error(member.ErrNoSuchMember), nil
	}
	if err := m.Active(); err != nil {
		return s.error(err), nil
	}
	return s.success(nil), nil
}

// 锁定/解锁会员
func (s *memberService) ToggleLock(ctx context.Context, memberId int64) (r *ttype.Result_, err error) {
	m := s.repo.GetMember(memberId)
	if m == nil {
		return s.error(member.ErrNoSuchMember), nil
	}
	state := m.GetValue().State
	if state == 1 {
		err = m.Lock()
	} else {
		err = m.Unlock()
	}
	if err != nil {
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
func (s *memberService) CheckProfileComplete(ctx context.Context, memberId int64) (r *ttype.Result_, e error) {
	m := s.repo.GetMember(memberId)
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

// 重置密码
func (s *memberService) ResetPassword(memberId int64) string {
	m := s.repo.GetMember(memberId)
	if m != nil {
		newPwd := domain.GenerateRandomIntPwd(6)
		newEncPwd := domain.MemberSha1Pwd(domain.Md5(newPwd))
		if err := m.Profile().ModifyPassword(newEncPwd, ""); err == nil {
			return newPwd
		} else {
			log.Println("--- 重置密码:", err)
		}
	}
	return ""
}

// 重置交易密码
func (s *memberService) ResetTradePwd(memberId int64) string {
	m := s.repo.GetMember(memberId)
	if m != nil {
		newPwd := domain.GenerateRandomIntPwd(6)
		newEncPwd := domain.TradePwd(domain.Md5(newPwd))
		if err := m.Profile().ModifyTradePassword(newEncPwd, ""); err == nil {
			return newPwd
		} else {
			log.Println("--- 重置交易密码:", err)
		}
	}
	return ""
}

// 修改密码
func (s *memberService) ModifyPassword(memberId int64, newPwd, oldPwd string) error {
	m := s.repo.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.Profile().ModifyPassword(newPwd, oldPwd)
}

//修改密码,传入密文密码
func (s *memberService) ModifyTradePassword(memberId int64,
	oldPwd, newPwd string) error {
	m := s.repo.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.Profile().ModifyTradePassword(newPwd, oldPwd)
}

// 登录，返回结果(Result_)和会员编号(ID);
// Result值为：-1:会员不存在; -2:账号密码不正确; -3:账号被停用
func (s *memberService) testLogin(user string, pwd string) (id int64, err error) {
	user = strings.ToLower(strings.TrimSpace(user))
	val := s.repo.GetMemberByUser(user)
	if val == nil {
		//todo: 界面加上使用手机号码登陆
		//val = m.repo.GetMemberValueByPhone(user)
	}
	if val == nil {
		return 0, member.ErrNoSuchMember
	}
	if val.Pwd != pwd {
		return 0, member.ErrCredential
	}
	if val.Flag&member.FlagLocked == member.FlagLocked {
		return 0, member.ErrMemberLocked
	}
	return val.Id, nil
}

// 登录，返回结果(Result_)和会员编号(ID);
// Result值为：-1:会员不存在; -2:账号密码不正确; -3:账号被停用
func (s *memberService) CheckLogin(ctx context.Context, user string, pwd string, update bool) (*ttype.Result_, error) {
	id, err := s.testLogin(user, pwd)
	if update && err == nil {
		m := s.repo.GetMember(id)
		err = m.UpdateLoginTime()
	}
	r := s.result(err)
	r.Data = map[string]string{
		"member_id": strconv.Itoa(int(id)),
	}
	return r, nil
}

// 检查交易密码
func (s *memberService) CheckTradePwd(ctx context.Context, id int64, tradePwd string) (r *ttype.Result_, err error) {
	m := s.repo.GetMember(id)
	if m == nil {
		return s.result(member.ErrNoSuchMember), nil
	}
	mv := m.GetValue()
	if mv.TradePwd == "" {
		return s.result(member.ErrNotSetTradePwd), nil
	}
	if mv.TradePwd != tradePwd {
		return s.result(member.ErrIncorrectTradePwd), nil
	}
	return s.success(nil), nil
}

// 检查与现有用户不同的用户是否存在,如存在则返回错误
func (s *memberService) CheckUsr(user string, memberId int64) error {
	if len(user) < 6 {
		return member.ErrUsrLength
	}
	if s.repo.CheckUsrExist(user, memberId) {
		return member.ErrUsrExist
	}
	return nil
}

// 检查手机号码是否与会员一致
func (s *memberService) CheckPhone(phone string, memberId int64) error {
	return s.repo.GetManager().CheckPhoneBind(phone, memberId)
}

// 获取会员账户
func (s *memberService) GetAccount(ctx context.Context, memberId int64) (*member_service.SAccount, error) {
	m := s.repo.CreateMember(&member.Member{Id: memberId})
	acc := m.GetAccount()
	if acc != nil {
		return parser.AccountDto(acc.GetValue()), nil
	}
	return nil, nil
}

// 获取上级邀请人会员编号数组
func (s *memberService) InviterArray(ctx context.Context, memberId int64, depth int32) (r []int64, err error) {
	m := s.repo.CreateMember(&member.Member{Id: memberId})
	if m != nil {
		return m.Invitation().InviterArray(memberId, int(depth)), nil
	}
	return []int64{}, nil
}

// 按条件获取荐指定等级会员的数量
func (s *memberService) GetInviterQuantity(ctx context.Context, memberId int64, data map[string]string) (int32, error) {
	where := ""
	if data != nil && len(data) > 0 {
		where = s.parseGetInviterDataParams(data)
	}
	return s.query.GetInviterQuantity(memberId, where), nil
}

// 按条件获取荐指定等级会员的列表
func (s *memberService) GetInviterArray(ctx context.Context, memberId int64, data map[string]string) ([]int64, error) {
	where := ""
	if data != nil && len(data) > 0 {
		where = s.parseGetInviterDataParams(data)
	}
	return s.query.GetInviterArray(memberId, where), nil
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

// 实名认证信息
func (s *memberService) GetTrustInfo(ctx context.Context, memberId int64) (*member_service.STrustedInfo, error) {
	t := member.TrustedInfo{}
	m := s.repo.GetMember(memberId)
	if m != nil {
		t = m.Profile().GetTrustedInfo()
	}
	return parser.TrustedInfoDto(&t), nil
}

// 保存实名认证信息
func (s *memberService) SaveTrustedInfo(memberId int64, v *member.TrustedInfo) error {
	m := s.repo.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.Profile().SaveTrustedInfo(v)
}

// 审核实名认证,若重复审核将返回错误
func (s *memberService) ReviewTrustedInfo(memberId int64, pass bool, remark string) error {
	m := s.repo.GetMember(memberId)
	return m.Profile().ReviewTrustedInfo(pass, remark)
}

// 获取分页商铺收藏
func (s *memberService) PagedShopFav(memberId int64, begin, end int,
	where string) (int, []*dto.PagedShopFav) {
	return s.query.PagedShopFav(memberId, begin, end, where)
}

// 获取分页商铺收藏
func (s *memberService) PagedGoodsFav(memberId int64, begin, end int,
	where string) (int, []*dto.PagedGoodsFav) {
	return s.query.PagedGoodsFav(memberId, begin, end, where)
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
	where, orderBy string) (int, []*order_service.SComplexOrder) {
	return s.orderQuery.PagedTradeOrderOfBuyer(buyerId, begin, size, pagination, where, orderBy)
}

/*********** 收货地址 ***********/

// 获取会员的收货地址
func (s *memberService) GetAddressList(ctx context.Context, memberId int64) ([]*member_service.SAddress, error) {
	src := s.repo.GetDeliverAddress(memberId)
	var arr []*member_service.SAddress
	for _, v := range src {
		arr = append(arr, parser.AddressDto(v))
	}
	return arr, nil
}

//获取配送地址
func (s *memberService) GetAddress(ctx context.Context, memberId int64, addrId int64) (
	*member_service.SAddress, error) {
	m := s.repo.CreateMember(&member.Member{Id: memberId})
	pro := m.Profile()
	var addr member.IDeliverAddress
	if addrId > 0 {
		addr = pro.GetAddress(addrId)
	} else {
		addr = pro.GetDefaultAddress()
	}
	if addr != nil {
		v := addr.GetValue()
		d := parser.AddressDto(&v)
		d.Area = s.valRepo.GetAreaString(
			v.Province, v.City, v.District)
		return d, nil
	}
	return nil, nil
}

//保存配送地址
func (s *memberService) SaveAddress(memberId int64, src *member_service.SAddress) (int64, error) {
	e := parser.Address(src)
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

// 查询有邀请关系的会员数量
func (s *memberService) GetReferNum(memberId int64, layer int) int {
	return s.query.GetReferNum(memberId, layer)
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

// 获取会员汇总信息
func (s *memberService) Complex(ctx context.Context, memberId int64) (*member_service.SComplexMember, error) {
	m := s.repo.GetMember(memberId)
	if m != nil {
		s := m.Complex()
		return parser.ComplexMemberDto(s), nil
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
func (s *memberService) AccountCharge(ctx context.Context, memberId int64, account int32,
	title string, amount int32, outerNo string, remark string) (*ttype.Result_, error) {
	var err error
	m := s.repo.CreateMember(&member.Member{Id: memberId})
	acc := m.GetAccount()
	if acc == nil {
		err = member.ErrNoSuchMember
	} else {
		err = acc.Charge(account, title, int(amount), outerNo, remark)
	}
	return s.result(err), nil
}

// 账户抵扣
func (s *memberService) AccountDiscount(ctx context.Context, memberId int64, account int32, title string,
	amount int32, outerNo string, remark string) (r *ttype.Result_, err error) {
	m, err := s.getMember(memberId)
	if err == nil {
		acc := m.GetAccount()
		err = acc.Discount(int(account), title, int(amount), outerNo, remark)
	}
	return s.result(err), nil
}

// 账户消耗
func (s *memberService) AccountConsume(ctx context.Context, memberId int64, account int32, title string,
	amount int32, outerNo string, remark string) (r *ttype.Result_, err error) {
	m, err := s.getMember(memberId)
	if err == nil {
		acc := m.GetAccount()
		err = acc.Consume(int(account), title, int(amount), outerNo, remark)
	}
	return s.result(err), nil
}

// 账户消耗
func (s *memberService) AccountRefund(ctx context.Context, memberId int64, account int32, title string,
	amount int32, outerNo string, remark string) (r *ttype.Result_, err error) {
	m, err := s.getMember(memberId)
	if err == nil {
		acc := m.GetAccount()
		err = acc.Refund(int(account), title, int(amount), outerNo, remark)
	}
	return s.result(err), nil
}

// 调整账户
func (s *memberService) AccountAdjust(ctx context.Context, memberId int64, account int32,
	amount int32, relateUser int64, remark string) (r *ttype.Result_, err error) {
	m, err := s.getMember(memberId)
	if err == nil {
		tit := "[KF]系统冲正"
		if amount > 0 {
			tit = "[KF]系统充值"
		}
		acc := m.GetAccount()
		err = acc.Adjust(int(account), tit, int(amount), remark, relateUser)
	}
	return s.result(err), nil
}

// !银行四要素认证
func (s *memberService) B4EAuth(ctx context.Context, memberId int64, action string, data map[string]string) (r *ttype.Result_, err error) {
	mod := module.Get(module.M_B4E).(*module.Bank4E)
	if action == "get" {
		data := mod.GetBasicInfo(memberId)
		d, err := json.Marshal(data)
		if err != nil {
			return s.error(err), nil
		}
		return s.success(map[string]string{"data": string(d)}), nil
	}
	if action == "update" {
		err := mod.UpdateInfo(memberId,
			data["real_name"],
			data["id_card"],
			data["phone"],
			data["bank_account"])
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
func (s *memberService) FilterMemberByUsrOrPhone(key string) []*dto.SimpleMember {
	return s.query.FilterMemberByUsrOrPhone(key)
}

// 根据用户名货手机获取会员
func (s *memberService) GetMemberByUserOrPhone(key string) *dto.SimpleMember {
	return s.query.GetMemberByUserOrPhone(key)
}

// 根据手机获取会员编号
func (s *memberService) GetMemberIdByPhone(phone string) int64 {
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

// 更改用户名
func (s *memberService) changeUsr(memberId int, user string) error {
	m := s.repo.GetMember(int64(memberId))
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.ChangeUsr(user)
}

// 更改手机号
func (s *memberService) changePhone(memberId int64, phone string) error {
	m := s.repo.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.Profile().ChangePhone(phone)
}
