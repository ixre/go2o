/**
 * Copyright 2014 @ to2.net.
 * name :
 * author : jarryliu
 * date : 2013-12-09 10:12
 * description :
 * history :
 */

package member

//todo: 要注意UpdateTime的更新

import (
	"errors"
	"github.com/ixre/gof/util"
	"go2o/core/domain/interface/enum"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/mss"
	"go2o/core/domain/interface/mss/notify"
	"go2o/core/domain/interface/registry"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/infrastructure/domain"
	"go2o/core/infrastructure/format"
	"go2o/core/infrastructure/tool/sms"
	"go2o/core/msq"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//todo: 依赖商户的 MSS 发送通知消息,应去掉
//todo: 会员升级 应单独来处理
var _ member.IMember = new(memberImpl)

type memberImpl struct {
	manager         member.IMemberManager
	value           *member.Member
	account         member.IAccount
	level           *member.Level
	repo            member.IMemberRepo
	relation        *member.InviteRelation
	invitation      member.IInvitationManager
	mssRepo         mss.IMssRepo
	registryRepo    registry.IRegistryRepo
	valueRepo       valueobject.IValueRepo
	profileManager  member.IProfileManager
	favoriteManager member.IFavoriteManager
	giftCardManager member.IGiftCardManager
}

func (m *memberImpl) ContainFlag(f int) bool {
	return m.value.Flag&f == f
}

func NewMember(manager member.IMemberManager, val *member.Member, rep member.IMemberRepo,
	mp mss.IMssRepo, valRepo valueobject.IValueRepo, registryRepo registry.IRegistryRepo) member.IMember {
	return &memberImpl{
		manager:      manager,
		value:        val,
		repo:         rep,
		mssRepo:      mp,
		valueRepo:    valRepo,
		registryRepo: registryRepo,
	}
}

// 获取聚合根编号
func (m *memberImpl) GetAggregateRootId() int64 {
	return m.value.Id
}

// 会员汇总信息
func (m *memberImpl) Complex() *member.ComplexMember {
	mv := m.GetValue()
	lv := m.GetLevel()
	pf := m.Profile()
	tr := pf.GetTrustedInfo()
	s := &member.ComplexMember{
		Name:           mv.Name,
		Avatar:         format.GetResUrl(mv.Avatar),
		Exp:            mv.Exp,
		Level:          mv.Level,
		LevelName:      lv.Name,
		InvitationCode: mv.InvitationCode,
		TrustAuthState: tr.ReviewState,
		PremiumUser:    mv.PremiumUser,
		Flag:           mv.Flag,
		UpdateTime:     mv.UpdateTime,
	}
	return s
}

// 会员资料服务
func (m *memberImpl) Profile() member.IProfileManager {
	if m.profileManager == nil {
		m.profileManager = newProfileManagerImpl(m,
			m.GetAggregateRootId(), m.repo, m.registryRepo, m.valueRepo)
	}
	return m.profileManager
}

// 会员收藏服务
func (m *memberImpl) Favorite() member.IFavoriteManager {
	if m.favoriteManager == nil {
		m.favoriteManager = newFavoriteManagerImpl(
			m.GetAggregateRootId(), m.repo)
	}
	return m.favoriteManager
}

// 礼品卡服务
func (m *memberImpl) GiftCard() member.IGiftCardManager {
	if m.giftCardManager == nil {
		m.giftCardManager = newGiftCardManagerImpl(
			m.GetAggregateRootId(), m.repo)
	}
	return m.giftCardManager
}

// 邀请管理
func (m *memberImpl) Invitation() member.IInvitationManager {
	if m.invitation == nil {
		m.invitation = &invitationManager{
			member: m,
		}
	}
	return m.invitation
}

// 获取值
func (m *memberImpl) GetValue() member.Member {
	return *m.value
}

var (
	userRegex  = regexp.MustCompile("^[a-zA-Z0-9_]{6,}$")
	emailRegex = regexp.MustCompile("^[A-Za-z0-9_\\-]+@[a-zA-Z0-9\\-]+(\\.[a-zA-Z0-9]+)+$")
	phoneRegex = regexp.MustCompile("^(13[0-9]|14[5|6|7]|15[0-9]|16[6|8]|18[0-9]|17[0|1|2|3|4|5|6|7|8]|19[9|8])(\\d{8})$")
)

// 设置值
func (m *memberImpl) SetValue(v *member.Member) error {
	v.User = m.value.User
	if len(m.value.InvitationCode) == 0 {
		m.value.InvitationCode = v.InvitationCode
	}
	if v.Level > 0 {
		m.value.Level = v.Level
	}
	if len(v.TradePwd) == 0 {
		m.value.TradePwd = v.TradePwd
	}
	return nil
}

// 发送验证码,并返回验证码
func (m *memberImpl) SendCheckCode(operation string, mssType int) (string, error) {
	const expiresMinutes = 10 //10分钟生效
	code := domain.NewCheckCode()
	m.value.CheckCode = code
	m.value.CheckExpires = time.Now().Add(time.Minute * expiresMinutes).Unix()
	_, err := m.Save()
	if err == nil {
		mgr := m.mssRepo.NotifyManager()
		pro := m.Profile().GetProfile()

		// 创建参数
		data := map[string]interface{}{
			"code":      code,
			"operation": operation,
			"minutes":   expiresMinutes,
		}

		// 根据消息类型发送信息
		switch mssType {
		case notify.TypePhoneMessage:
			// 某些短信平台要求传入模板ID,在这里附加参数
			provider, _ := m.valueRepo.GetDefaultSmsApiPerm()
			data = sms.AppendCheckPhoneParams(provider, data)

			// 构造并发送短信
			n := mgr.GetNotifyItem("验证手机")
			c := notify.PhoneMessage(n.Content)
			err = mgr.SendPhoneMessage(pro.Phone, c, data)

		default:
		case notify.TypeEmailMessage:
			n := mgr.GetNotifyItem("验证邮箱")
			c := &notify.MailMessage{
				Subject: operation + "验证码",
				Body:    n.Content,
			}
			err = mgr.SendEmail(pro.Phone, c, data)
		}
	}
	return code, err
}

// 对比验证码
func (m *memberImpl) CompareCode(code string) error {
	if m.value.CheckCode != strings.TrimSpace(code) {
		return member.ErrCheckCodeError
	}
	if m.value.CheckExpires < time.Now().Unix() {
		return member.ErrCheckCodeExpires
	}
	return nil
}

// 获取账户
func (m *memberImpl) GetAccount() member.IAccount {
	if m.account == nil {
		v := m.repo.GetAccount(m.value.Id)
		if v == nil {
			v = &member.Account{
				MemberId: m.GetAggregateRootId(),
			}
		}
		return NewAccount(m, v, m.repo, m.manager, m.registryRepo)
	}
	return m.account
}

// 增加经验值
func (m *memberImpl) AddExp(exp int) error {
	m.value.Exp += exp
	_, err := m.Save()
	m.checkLevelUp() //判断是否升级
	return err
}

// 升级为高级会员

func (m *memberImpl) Premium(v int, expires int64) error {
	switch v {
	case member.PremiumNormal:
		m.value.PremiumUser = v
		m.value.PremiumExpires = 0
	case member.PremiumGold, member.PremiumWhiteGold, member.PremiumSuper:
		m.value.PremiumUser = v
		m.value.PremiumExpires = expires
	default:
		return member.ErrPremiumValue
	}
	_, err := m.Save()
	return err
}

// 获取等级
func (m *memberImpl) GetLevel() *member.Level {
	if m.level == nil {
		m.level = m.manager.LevelManager().
			GetLevelById(m.value.Level)
	}
	return m.level
}

// 检查升级
func (m *memberImpl) checkLevelUp() bool {
	lg := m.manager.LevelManager()
	levelId := lg.GetLevelIdByExp(m.value.Exp)
	if levelId == 0 {
		return false
	}
	// 判断是否大于当前等级
	if m.value.Level > levelId {
		return false
	}
	// 判断等级是否启用
	lv := lg.GetLevelById(levelId)
	if lv.Enabled == 0 || lv.AllowUpgrade == 0 {
		return false
	}
	origin := m.value.Level
	unix := time.Now().Unix()
	m.value.Level = levelId
	m.value.UpdateTime = unix
	_, err := m.Save()
	if err == nil {
		m.level = nil
		lvLog := &member.LevelUpLog{
			MemberId:    m.GetAggregateRootId(),
			OriginLevel: origin,
			TargetLevel: levelId,
			IsFree:      1,
			PaymentId:   0,
			ReviewState: int(enum.ReviewConfirm),
			UpgradeMode: member.LAutoUpgrade,
			CreateTime:  unix,
		}
		_, err = m.repo.SaveLevelUpLog(lvLog)
	}
	return true
}

// 更改会员等级
func (m *memberImpl) ChangeLevel(level int, paymentId int, review bool) error {
	lg := m.manager.LevelManager()
	lv := lg.GetLevelById(level)
	// 判断等级是否启用
	if lv == nil || lv.Enabled == 0 {
		return member.ErrLevelDisabled
	}
	origin := m.value.Level
	unix := time.Now().Unix()
	lvLog := &member.LevelUpLog{
		MemberId:    m.GetAggregateRootId(),
		OriginLevel: origin,
		TargetLevel: level,
		PaymentId:   paymentId,
		ReviewState: int(enum.ReviewNotSet),
		UpgradeMode: member.LServiceAgentUpgrade,
		CreateTime:  unix,
	}
	if paymentId == 0 {
		lvLog.IsFree = 1
	}
	if !review {
		lvLog.ReviewState = int(enum.ReviewConfirm)
	}
	_, err := m.repo.SaveLevelUpLog(lvLog)
	if err == nil && !review {
		m.value.Exp = lv.RequireExp
		m.value.Level = level
		m.value.UpdateTime = unix
		_, err = m.Save()
		if err == nil {
			m.level = nil
		}
	}
	return err
}

// 标志赋值, 如果flag小于零, 则异或运算
func (m *memberImpl) GrantFlag(flag int) error {
	f := int(math.Abs(float64(flag)))
	if f&(f-1) != 0 {
		return errors.New("not right flag value")
	}
	if f < 128 {
		return errors.New("disallow grant system flag, flag must large than or equals 128")
	}
	own := m.value.Flag&f == f
	if flag > 0 {
		if own {
			return errors.New("member has granted flag:" + strconv.Itoa(flag))
		}
		m.value.Flag |= flag
	} else {
		if !own {
			return errors.New("member not grant flag:" + strconv.Itoa(flag))
		}
		m.value.Flag ^= f
	}
	_, err := m.Save()
	return err
}

// 审核升级请求
func (m *memberImpl) ReviewLevelUp(id int, pass bool) error {
	l := m.repo.GetLevelUpLog(int32(id))
	if l != nil && l.MemberId == m.GetAggregateRootId() {
		if l.ReviewState == int(enum.ReviewPass) {
			return member.ErrLevelUpPass
		}
		if l.ReviewState == int(enum.ReviewReject) {
			return member.ErrLevelUpReject
		}
		if l.ReviewState == int(enum.ReviewConfirm) {
			return member.ErrLevelUpConfirm
		}
		if time.Now().Unix()-l.CreateTime < 120 {
			return member.ErrLevelUpLaterConfirm
		}
		if pass {
			l.ReviewState = int(enum.ReviewPass)
			_, err := m.repo.SaveLevelUpLog(l)
			if err == nil {
				lv := m.manager.LevelManager().GetLevelById(l.TargetLevel)
				m.value.Exp = lv.RequireExp
				m.value.Level = l.TargetLevel
				m.value.UpdateTime = time.Now().Unix()
				_, err = m.Save()
			}
			return err
		} else {
			l.ReviewState = int(enum.ReviewReject)
			_, err := m.repo.SaveLevelUpLog(l)
			return err
		}
	}
	return member.ErrNoSuchLevelUpLog

}

// 标记已经处理升级
func (m *memberImpl) ConfirmLevelUp(id int32) error {
	l := m.repo.GetLevelUpLog(id)
	if l != nil && l.MemberId == m.GetAggregateRootId() {
		if l.ReviewState == int(enum.ReviewConfirm) {
			return member.ErrLevelUpConfirm
		}
		if l.ReviewState != int(enum.ReviewPass) {
			return member.ErrLevelUpReject
		}
		l.ReviewState = int(enum.ReviewConfirm)
		_, err := m.repo.SaveLevelUpLog(l)
		return err
	}
	return member.ErrNoSuchLevelUpLog
}

// 获取会员关联
func (m *memberImpl) GetRelation() *member.InviteRelation {
	if m.relation == nil {
		rel := m.repo.GetRelation(m.GetAggregateRootId())
		if rel == nil {
			rel = &member.InviteRelation{
				MemberId:  m.GetAggregateRootId(),
				CardCard:  "",
				InviterId: 0,
				RegMchId:  0,
			}
		}
		m.relation = rel
	}
	return m.relation
}

// 更换用户名
func (m *memberImpl) ChangeUsr(user string) error {
	if user == m.value.User {
		return member.ErrSameUsr
	}
	err := m.checkUser(m.value.User)
	if err == nil {
		m.value.User = user
		_, err = m.Save()
	}
	return err
}

// 更新登录时间
func (m *memberImpl) UpdateLoginTime() error {
	unix := time.Now().Unix()
	m.value.LastLoginTime = m.value.LoginTime
	m.value.LoginTime = unix
	m.value.UpdateTime = unix
	_, err := m.Save()
	return err
}

// 保存
func (m *memberImpl) Save() (int64, error) {
	m.value.UpdateTime = time.Now().Unix() // 更新时间，数据以更新时间触发
	if m.value.Id > 0 {
		return m.repo.SaveMember(m.value)
	}
	return m.create(m.value)
}

// 激活
func (m *memberImpl) Active() error {
	if m.ContainFlag(member.FlagActive) {
		return member.ErrMemberHasActive
	}
	if m.ContainFlag(member.FlagLocked) {
		return member.ErrMemberLocked
	}
	m.value.Flag |= member.FlagActive
	_, err := m.Save()
	return err
}

// 锁定会员
func (m *memberImpl) Lock() error {
	if m.ContainFlag(member.FlagLocked) {
		return nil
	}
	m.value.Flag |= member.FlagLocked
	_, err := m.Save()
	return err
}

// 解锁会员
func (m *memberImpl) Unlock() error {
	if !m.ContainFlag(member.FlagLocked) {
		return nil
	}
	m.value.Flag ^= member.FlagLocked
	_, err := m.Save()
	return err
}

// 创建会员
func (m *memberImpl) create(v *member.Member) (int64, error) {
	err := m.prepare()
	if err == nil {
		unix := time.Now().Unix()
		v.State = 1
		v.RegTime = unix
		v.LastLoginTime = unix
		v.Level = 1
		v.Exp = 0
		v.DynamicToken = ""
		if len(v.RegFrom) == 0 {
			v.RegFrom = ""
		}
		// 设置VIP用户信息
		v.PremiumUser = member.PremiumNormal
		v.PremiumExpires = 0
		// 创建一个用户编码
		v.Code = m.generateMemberCode()
		// 创建一个邀请码
		v.InvitationCode = m.generateInvitationCode()
		id, err1 := m.repo.SaveMember(v)
		if err1 == nil {
			m.value.Id = id
			go m.memberInit()
		} else {
			err = err1
		}
	}
	return m.GetAggregateRootId(), err
}

// 验证用户名
func (m *memberImpl) checkUser(user string) error {
	if len([]rune(user)) < 6 {
		return member.ErrUsrLength
	}
	if !userRegex.MatchString(user) {
		return member.ErrUsrValidErr
	}
	if m.repo.CheckUsrExist(user, m.GetAggregateRootId()) {
		return member.ErrUsrExist
	}
	return nil
}

// 会员初始化
func (m *memberImpl) memberInit() error {
	// 创建账户
	m.account = NewAccount(m, &member.Account{}, m.repo, m.manager, m.registryRepo)
	if _, err := m.account.Save(); err != nil {
		return err
	}
	// 注册后赠送积分
	regPresent := m.registryRepo.Get(registry.MemberRegisterPresentIntegral).IntValue()
	if regPresent > 0 {
		go m.GetAccount().Charge(member.AccountIntegral, "新会员注册赠送积分",
			regPresent, "-", "sys")
	}
	return nil
}

// 检查注册信息是否正确
func (m *memberImpl) prepare() (err error) {

	phoneAsUser := m.registryRepo.Get(registry.MemberRegisterPhoneAsUser).BoolValue()
	mustBindPhone := m.registryRepo.Get(registry.MemberRegisterMustBindPhone).BoolValue()
	// 验证用户名,如果填写了或非用手机号作为用户名,均验证用户名
	m.value.User = strings.TrimSpace(m.value.User)
	if m.value.User != "" || !phoneAsUser {
		if err = m.checkUser(m.value.User); err != nil {
			return err
		}
	}
	// 验证密码
	m.value.Pwd = strings.TrimSpace(m.value.Pwd)
	if len(m.value.Pwd) < 6 {
		return member.ErrPwdLength
	}
	// 验证手机
	m.value.Phone = strings.TrimSpace(m.value.Phone)
	lp := len(m.value.Phone)
	if mustBindPhone && lp == 0 {
		return member.ErrMissingPhone
	}
	if lp > 0 {
		checkPhone := m.registryRepo.Get(registry.MemberCheckPhoneFormat).BoolValue()
		if checkPhone && !phoneRegex.MatchString(m.value.Phone) {
			return member.ErrBadPhoneFormat
		}
		if m.checkPhoneBind(m.value.Phone, m.GetAggregateRootId()) != nil {
			return member.ErrPhoneHasBind
		}
	}
	// 使用手机号作为用户名
	if phoneAsUser {
		if m.repo.CheckUsrExist(m.value.Phone, 0) {
			return member.ErrPhoneHasBind
		}
		m.value.User = m.value.Phone
	}

	// 验证IM
	//pro.Im = strings.TrimSpace(pro.Im)
	//if perm.NeedIm && len(pro.Im) == 0 {
	//	return 0, errors.New(strings.Replace(member.ErrMissingIM.Error(),
	//		"IM", variable.AliasMemberIM, -1))
	//}
	m.value.Name = strings.TrimSpace(m.value.Name)

	//如果未设置昵称,则默认为用户名
	if len(m.value.Name) == 0 {
		m.value.Name = "User" + m.value.User
	}
	m.value.Avatar = strings.TrimSpace(m.value.Avatar)
	if len(m.value.Avatar) == 0 {
		m.value.Avatar = "res/no_avatar.gif"
	}
	return err
}

// 检查手机绑定,同时检查手机格式
func (m *memberImpl) checkPhoneBind(phone string, memberId int64) error {
	if len(phone) <= 0 {
		return member.ErrMissingPhone
	}
	if m.repo.CheckPhoneBind(phone, memberId) {
		return member.ErrPhoneHasBind
	}
	return nil
}

func (m *memberImpl) generateMemberCode() string {
	var code string
	for {
		code = util.RandString(6)
		if memberId := m.repo.GetMemberIdByCode(code); memberId == 0 {
			break
		}
	}
	return code
}

// 创建邀请码
func (m *memberImpl) generateInvitationCode() string {
	var code string
	for {
		code = domain.GenerateInvitationCode()
		if memberId := m.repo.GetMemberIdByInvitationCode(code); memberId == 0 {
			break
		}
	}
	return code
}

// 更新邀请关系
func (m *memberImpl) updateDepthInvite(r *member.InviteRelation) {
	if r.InviterId > 0 {
		arr := m.Invitation().InviterArray(r.InviterId, 2)
		r.InviterD2 = arr[0]
		r.InviterD3 = arr[1]
	}
}

// 保存推荐关系
func (m *memberImpl) saveRelation(r *member.InviteRelation) error {
	m.relation = r
	m.relation.MemberId = m.value.Id
	m.updateDepthInvite(m.relation)
	err := m.repo.SaveRelation(m.relation)
	if err == nil {
		// 推送关系更新消息
		go msq.PushDelay(msq.MemberRelationUpdated, strconv.Itoa(int(m.GetAggregateRootId())), "", 1000)
	}
	return err
}

// 绑定邀请人,如果已邀请,force为true时更新
func (m *memberImpl) BindInviter(memberId int64, force bool) error {
	if memberId > 0 {
		if rm := m.repo.GetMember(memberId); rm == nil {
			return member.ErrNoValidInviter
		}
	}
	rl := m.GetRelation()
	if true || rl.InviterId != memberId {
		rl.InviterId = memberId
		return m.saveRelation(rl)
	}
	return nil
}

var _ member.IFavoriteManager = new(favoriteManagerImpl)

// 收藏服务
type favoriteManagerImpl struct {
	_memberId int64
	_rep      member.IMemberRepo
}

func newFavoriteManagerImpl(memberId int64,
	rep member.IMemberRepo) member.IFavoriteManager {
	if memberId == 0 {
		//如果会员不存在,则不应创建服务
		panic(errors.New("member not exists"))
	}
	return &favoriteManagerImpl{
		_memberId: memberId,
		_rep:      rep,
	}
}

// 收藏
func (m *favoriteManagerImpl) Favorite(favType int, referId int32) error {
	if m.Favored(favType, referId) {
		return member.ErrFavored
	}
	return m._rep.Favorite(m._memberId, favType, referId)
}

// 是否已收藏
func (m *favoriteManagerImpl) Favored(favType int, referId int32) bool {
	return m._rep.Favored(m._memberId, favType, referId)
}

// 取消收藏
func (m *favoriteManagerImpl) Cancel(favType int, referId int32) error {
	return m._rep.CancelFavorite(m._memberId, favType, referId)
}
