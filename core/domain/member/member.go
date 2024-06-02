/**
 * Copyright 2014 @ 56x.net.
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
	"regexp"
	"strconv"
	"strings"
	"time"

	de "github.com/ixre/go2o/core/domain/interface/domain"
	"github.com/ixre/go2o/core/domain/interface/domain/enum"
	"github.com/ixre/go2o/core/domain/interface/member"
	mss "github.com/ixre/go2o/core/domain/interface/message"
	"github.com/ixre/go2o/core/domain/interface/message/notify"
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/domain/interface/valueobject"
	"github.com/ixre/go2o/core/domain/interface/wallet"
	"github.com/ixre/go2o/core/event/events"
	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/go2o/core/infrastructure/format"
	"github.com/ixre/go2o/core/infrastructure/util/collections"
	"github.com/ixre/gof/domain/eventbus"
	"github.com/ixre/gof/util"
)

// todo: 依赖商户的 MSS 发送通知消息,应去掉
// todo: 会员升级 应单独来处理
var _ member.IMemberAggregateRoot = new(memberImpl)

type memberImpl struct {
	manager         member.IMemberManager
	value           *member.Member
	account         member.IAccount
	level           *member.Level
	repo            member.IMemberRepo
	walletRepo      wallet.IWalletRepo
	relation        *member.InviteRelation
	invitation      member.IInvitationManager
	mssRepo         mss.IMssRepo
	registryRepo    registry.IRegistryRepo
	valueRepo       valueobject.IValueRepo
	profileManager  member.IProfileManager
	favoriteManager member.IFavoriteManager
	giftCardManager member.IGiftCardManager
}

func NewMember(manager member.IMemberManager, val *member.Member,
	rep member.IMemberRepo, walletRepo wallet.IWalletRepo,
	mp mss.IMssRepo, valRepo valueobject.IValueRepo,
	registryRepo registry.IRegistryRepo) member.IMemberAggregateRoot {
	return &memberImpl{
		manager:      manager,
		value:        val,
		repo:         rep,
		mssRepo:      mp,
		walletRepo:   walletRepo,
		valueRepo:    valRepo,
		registryRepo: registryRepo,
	}
}

// 获取聚合根编号
func (m *memberImpl) GetAggregateRootId() int64 {
	return m.value.Id
}

// Complex 会员汇总信息
func (m *memberImpl) Complex() *member.ComplexMember {
	mv := m.GetValue()
	lv := m.GetLevel()
	pf := m.Profile()
	tr := pf.GetTrustedInfo()
	s := &member.ComplexMember{
		Nickname:            mv.Nickname,
		RealName:            mv.RealName,
		Avatar:              format.GetFileFullUrl(mv.Portrait),
		Exp:                 mv.Exp,
		Level:               mv.Level,
		LevelName:           lv.Name,
		TrustAuthState:      tr.ReviewStatus,
		TradePasswordHasSet: mv.TradePassword != "",
		PremiumUser:         mv.PremiumUser,
		Flag:                mv.UserFlag,
		UpdateTime:          mv.UpdateTime,
	}
	return s
}

// Profile 会员资料服务
func (m *memberImpl) Profile() member.IProfileManager {
	if m.profileManager == nil {
		m.profileManager = newProfileManagerImpl(m,
			m.GetAggregateRootId(), m.repo, m.registryRepo, m.valueRepo)
	}
	return m.profileManager
}

// Favorite 会员收藏服务
func (m *memberImpl) Favorite() member.IFavoriteManager {
	if m.favoriteManager == nil {
		m.favoriteManager = newFavoriteManagerImpl(
			m.GetAggregateRootId(), m.repo)
	}
	return m.favoriteManager
}

// GiftCard 礼品卡服务
func (m *memberImpl) GiftCard() member.IGiftCardManager {
	if m.giftCardManager == nil {
		m.giftCardManager = newGiftCardManagerImpl(
			m.GetAggregateRootId(), m.repo)
	}
	return m.giftCardManager
}

// Invitation 邀请管理
func (m *memberImpl) Invitation() member.IInvitationManager {
	if m.invitation == nil {
		m.invitation = &invitationManager{
			member: m,
		}
	}
	return m.invitation
}

// GetValue 获取值
func (m *memberImpl) GetValue() member.Member {
	return *m.value
}

var (
	userRegex  = regexp.MustCompile("^[a-zA-Z0-9_]{6,}$")
	emailRegex = regexp.MustCompile("^[A-Za-z0-9_\\-]+@[a-zA-Z0-9\\-]+(\\.[a-zA-Z0-9]+)+$")
	phoneRegex = regexp.MustCompile("^(13[0-9]|14[5|6|7]|15[0-9]|16[5|6|7|8]|18[0-9]|17[0|1|2|3|4|5|6|7|8]|19[1|8|9])(\\d{8})$")
)

// SendCheckCode 发送验证码,并返回验证码
func (m *memberImpl) SendCheckCode(operation string, mssType int) (string, error) {
	const expiresMinutes = 10 //10分钟生效
	code := domain.NewCheckCode()
	m.value.CheckCode = code
	m.value.CheckExpires = time.Now().Add(time.Minute * expiresMinutes).Unix()
	_, err := m.Save()
	if err == nil {
		// 创建参数
		data := []string{
			operation,
			code,
			strconv.Itoa(expiresMinutes),
		}
		mgr := m.mssRepo.NotifyManager()
		// 根据消息类型发送信息
		switch mssType {
		default:
		case notify.TypeEmailMessage:
			n := mgr.GetNotifyItem("验证邮箱")
			c := &notify.MailMessage{
				Subject: operation + "验证码",
				Body:    n.Content,
			}
			err = mgr.SendEmail(m.value.Email, c, data)
		case notify.TypePhoneMessage:
			// 某些短信平台要求传入模板ID,在这里附加参数
			// re := m.registryRepo.Get(registry.SmsMemberCheckTemplateId)
			// data["templateId"] = re.StringValue()
			// 构造并发送短信
			n := mgr.GetNotifyItem("验证手机")
			c := notify.PhoneMessage(n.Content)
			err = mgr.SendPhoneMessage(m.value.Phone, c, data, "")
		}
	}
	return code, err
}

// CompareCode 对比验证码
func (m *memberImpl) CompareCode(code string) error {
	if m.value.CheckCode != strings.TrimSpace(code) {
		return de.ErrCheckCodeError
	}
	if m.value.CheckExpires < time.Now().Unix() {
		return de.ErrCheckCodeExpires
	}
	return nil
}

// GetAccount 获取账户
func (m *memberImpl) GetAccount() member.IAccount {
	if m.account == nil {
		v := m.repo.GetAccount(m.value.Id)
		if v == nil {
			v = &member.Account{
				MemberId: m.GetAggregateRootId(),
			}
		}
		return newAccount(m, v, m.repo, m.manager, m.walletRepo, m.registryRepo)
	}
	return m.account
}

// AddExp 增加经验值
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

// GetLevel 获取等级
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
			MemberId:     m.GetAggregateRootId(),
			OriginLevel:  origin,
			TargetLevel:  levelId,
			IsFree:       1,
			PaymentId:    0,
			ReviewStatus: int(enum.ReviewConfirm),
			UpgradeMode:  member.LAutoUpgrade,
			CreateTime:   unix,
		}
		_, err = m.repo.SaveLevelUpLog(lvLog)
	}
	return true
}

// ChangeLevel 更改会员等级
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
		MemberId:     m.GetAggregateRootId(),
		OriginLevel:  origin,
		TargetLevel:  level,
		PaymentId:    paymentId,
		ReviewStatus: int(enum.ReviewNotSet),
		UpgradeMode:  member.LServiceAgentUpgrade,
		CreateTime:   unix,
	}
	if paymentId == 0 {
		lvLog.IsFree = 1
	}
	if !review {
		lvLog.ReviewStatus = int(enum.ReviewConfirm)
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

// ContainFlag 是否包含标志
func (m *memberImpl) ContainFlag(f int) bool {
	return m.value.UserFlag&f == f
}

// GrantFlag 标志赋值, 如果flag小于零, 则异或运算
func (m *memberImpl) GrantFlag(flag int) error {
	if flag < 128 {
		return errors.New("disallow grant system flag, flag must large than or equals 128")
	}
	v, err := domain.GrantFlag(m.value.UserFlag, flag)
	if err == nil {
		m.value.UserFlag = v
		_, err = m.Save()
	}
	return err
}

func (m *memberImpl) TestFlag(flag int) bool {
	return domain.TestFlag(m.value.UserFlag, flag)
}

// ReviewLevelUp 审核升级请求
func (m *memberImpl) ReviewLevelUp(id int, pass bool) error {
	l := m.repo.GetLevelUpLog(id)
	if l != nil && l.MemberId == m.GetAggregateRootId() {
		if l.ReviewStatus == int(enum.ReviewPass) {
			return member.ErrLevelUpPass
		}
		if l.ReviewStatus == int(enum.ReviewReject) {
			return member.ErrLevelUpReject
		}
		if l.ReviewStatus == int(enum.ReviewConfirm) {
			return member.ErrLevelUpConfirm
		}
		if time.Now().Unix()-l.CreateTime < 120 {
			return member.ErrLevelUpLaterConfirm
		}
		if pass {
			l.ReviewStatus = int(enum.ReviewPass)
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
			l.ReviewStatus = int(enum.ReviewReject)
			_, err := m.repo.SaveLevelUpLog(l)
			return err
		}
	}
	return member.ErrNoSuchLevelUpLog

}

// ConfirmLevelUp 标记已经处理升级
func (m *memberImpl) ConfirmLevelUp(id int) error {
	l := m.repo.GetLevelUpLog(id)
	if l != nil && l.MemberId == m.GetAggregateRootId() {
		if l.ReviewStatus == int(enum.ReviewConfirm) {
			return member.ErrLevelUpConfirm
		}
		if l.ReviewStatus != int(enum.ReviewPass) {
			return member.ErrLevelUpReject
		}
		l.ReviewStatus = int(enum.ReviewConfirm)
		_, err := m.repo.SaveLevelUpLog(l)
		return err
	}
	return member.ErrNoSuchLevelUpLog
}

// GetRelation 获取会员关联
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
			m.repo.SaveRelation(rel)
		}
		m.relation = rel
	}
	return m.relation
}

// 更换用户名
func (m *memberImpl) ChangeUsername(user string) error {
	user = strings.TrimSpace(user)
	if len(user) == 0 {
		return member.ErrInvalidUsername
	}
	if user == m.value.Username {
		return member.ErrSameUser
	}
	err := m.checkUser(user)
	if err == nil {
		m.value.Username = user
		_, err = m.Save()
	}
	return err
}

// UpdateLoginTime 更新登录时间
func (m *memberImpl) UpdateLoginTime() error {
	unix := time.Now().Unix()
	m.value.LastLoginTime = m.value.LoginTime
	m.value.LoginTime = unix
	m.value.UpdateTime = unix
	_, err := m.Save()
	return err
}

// Save 保存
func (m *memberImpl) Save() (int64, error) {
	m.value.UpdateTime = time.Now().Unix() // 更新时间，数据以更新时间触发
	if m.value.Id > 0 {
		_, err := m.repo.SaveMember(m.value)
		if err == nil {
			go m.pushSaveEvent(false)
		}
		return m.GetAggregateRootId(), err
	}
	return m.create(m.value)
}

func (m *memberImpl) pushSaveEvent(create bool) {
	rl := m.GetRelation()
	eventbus.Publish(&events.MemberPushEvent{
		IsCreate:  create,
		Member:    m.value,
		InviterId: int(rl.InviterId),
	})
}

// Active 激活
func (m *memberImpl) Active() error {
	if m.ContainFlag(member.FlagActive) {
		return member.ErrMemberHasActive
	}
	if m.ContainFlag(member.FlagLocked) {
		return member.ErrMemberLocked
	}
	m.value.UserFlag |= member.FlagActive
	_, err := m.Save()
	return err
}

// Lock 锁定会员
func (m *memberImpl) Lock(minutes int, remark string) error {
	if m.ContainFlag(member.FlagLocked) {
		return nil
	}
	m.value.UserFlag |= member.FlagLocked
	_, err := m.Save()
	if err == nil {
		now := time.Now().Unix()
		ml := &member.MmLockInfo{
			MemberId:   m.GetAggregateRootId(),
			LockTime:   now,
			UnlockTime: now + int64(minutes*60),
			Remark:     remark,
		}
		his := &member.MmLockHistory{
			MemberId: ml.MemberId,
			LockTime: ml.LockTime,
			Duration: minutes,
			Remark:   remark,
		}
		// 永久锁定
		if minutes <= 0 {
			ml.UnlockTime = -1
			his.Duration = -1
		}
		_, err = m.repo.SaveLockInfo(ml)
		if err == nil {
			_, err = m.repo.SaveLockHistory(his)
			// 注册解锁任务
			if ml.UnlockTime > 0 {
				m.repo.RegisterUnlockJob(ml)
			}
		}
	}
	return err
}

// Unlock 解锁会员
func (m *memberImpl) Unlock() error {
	if !m.ContainFlag(member.FlagLocked) {
		return nil
	}
	m.value.UserFlag ^= member.FlagLocked
	_, err := m.Save()
	if err == nil {
		err = m.repo.DeleteLockInfos(m.GetAggregateRootId())
	}
	return err
}

// 根据注册来源计算会员角色身份
func (m *memberImpl) getUserRoleFlag(v *member.Member) int {
	ret := member.RoleUser
	if len(v.RegFrom) != 0 {
		// 根据注册来源设置角色
		v.RegFrom = ""
		if strings.Contains(v.RegFrom, "EMPLOYEE") {
			// 商户职员
			ret |= member.RoleEmployee
		}
		if strings.Contains(v.RegFrom, "EXT1") {
			// 扩展角色1
			ret |= member.RoleEmployee
		}
		if strings.Contains(v.RegFrom, "EXT2") {
			// 扩展角色2
			ret |= member.RoleEmployee
		}
	}
	return ret
}

// 创建会员
func (m *memberImpl) create(v *member.Member) (int64, error) {
	err := m.prepare()
	if err == nil {
		unix := time.Now().Unix()
		v.RegTime = unix
		v.LastLoginTime = unix
		v.Level = 1
		v.Exp = 0
		// 设置VIP用户信息
		v.PremiumUser = member.PremiumNormal
		v.PremiumExpires = 0
		// 创建一个用户编码/邀请码
		v.UserCode = m.generateMemberCode()
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
		return member.ErrUserLength
	}
	if !userRegex.MatchString(user) {
		return member.ErrUserValidErr
	}
	if m.repo.CheckUserExist(user, m.GetAggregateRootId()) {
		return member.ErrUserExist
	}
	return nil
}

// 会员初始化
func (m *memberImpl) memberInit() error {
	// 创建账户
	m.account = newAccount(m, &member.Account{},
		m.repo, m.manager, m.walletRepo, m.registryRepo)
	if _, err := m.account.Save(); err != nil {
		return err
	}
	// 注册后赠送积分
	regPresent := m.registryRepo.Get(registry.MemberRegisterPresentIntegral).IntValue()
	if regPresent > 0 {
		go m.GetAccount().CarryTo(member.AccountIntegral, member.AccountOperateData{
			Title:   "新会员注册赠送积分",
			Amount:  regPresent,
			OuterNo: "-",
			Remark:  "sys",
		}, false, 0)
	}
	go m.pushSaveEvent(true)
	return nil
}

// 检查注册信息是否正确
func (m *memberImpl) prepare() (err error) {
	phoneAsUser := m.registryRepo.Get(registry.MemberRegisterPhoneAsUser).BoolValue()
	mustBindPhone := m.registryRepo.Get(registry.MemberRegisterMustBindPhone).BoolValue()
	// 用户名全小写
	m.value.Username = strings.ToLower(m.value.Username)
	// 验证用户名,如果填写了或非用手机号作为用户名,均验证用户名
	// 使用手机号作为用户名
	if phoneAsUser {
		if m.repo.CheckUserExist(m.value.Phone, 0) {
			return member.ErrPhoneHasBind
		}
		m.value.Username = m.value.Phone
	}
	if len(m.value.Username) == 0 {
		return member.ErrInvalidUsername
	}
	if m.value.Username != "" {
		if err = m.checkUser(m.value.Username); err != nil {
			return err
		}
	}
	// 验证密码
	m.value.Password = strings.TrimSpace(m.value.Password)
	if len(m.value.Password) < 6 {
		return de.ErrPwdStrongLength
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
			return member.ErrInvalidPhone
		}
		if m.checkPhoneBind(m.value.Phone, m.GetAggregateRootId()) != nil {
			return member.ErrPhoneHasBind
		}
	}

	// 验证IM
	//pro.Im = strings.TrimSpace(pro.Im)
	//if perm.NeedIm && len(pro.Im) == 0 {
	//	return 0, errors.New(strings.Replace(member.ErrMissingIM.Error(),
	//		"IM", variable.AliasMemberIM, -1))
	//}
	m.value.Nickname = strings.TrimSpace(m.value.Nickname)
	m.value.RealName = strings.TrimSpace(m.value.RealName)
	//如果未设置昵称,则默认为用户名
	if len(m.value.Nickname) == 0 {
		m.value.Nickname = "User" + m.value.Username
	}
	// 初始化头像
	m.value.Portrait = strings.TrimSpace(m.value.Portrait)
	if len(m.value.Portrait) == 0 {
		m.value.Portrait = "static/init/avatar.png"
	}
	// 验证角色
	if m.value.RoleFlag != 0 && !collections.InArray([]int{
		member.RoleMerchant,
		member.RoleEmployee,
	}, m.value.RoleFlag) {
		return errors.New("用户类型不合法")
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

// 创建用户代码
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

// 绑定邀请人,如果已邀请,force为true时更新
func (m *memberImpl) BindInviter(inviterId int64, force bool) (err error) {
	rl := m.GetRelation()
	if !force && rl.InviterId > 0 {
		return member.ErrExistsInviter
	}
	// 不能绑定自己为推荐人
	if m.GetAggregateRootId() == inviterId {
		return member.ErrInvalidInviter
	}
	// 更改邀请人,在更改邀请人方法里会验证是否绑定下级会员
	if rl.InviterId != inviterId {
		isPush := rl.InviterId > 0 //  仅仅更改推荐人时才会推送信息
		m.relation = nil           // 清除缓存
		err = m.Invitation().UpdateInviter(inviterId, true)
		if err == nil && isPush {
			m.pushSaveEvent(false)
		}
	}
	return err
}

var _ member.IFavoriteManager = new(favoriteManagerImpl)

// 收藏服务
type favoriteManagerImpl struct {
	_memberId int64
	_rep      member.IMemberRepo
}

func newFavoriteManagerImpl(memberId int64,
	rep member.IMemberRepo) member.IFavoriteManager {
	if memberId <= 0 {
		return nil
	}
	return &favoriteManagerImpl{
		_memberId: memberId,
		_rep:      rep,
	}
}

// Favorite 收藏
func (m *favoriteManagerImpl) Favorite(favType int, referId int64) error {
	if m.Favored(favType, referId) {
		return member.ErrFavored
	}
	return m._rep.Favorite(m._memberId, favType, referId)
}

// 是否已收藏
func (m *favoriteManagerImpl) Favored(favType int, referId int64) bool {
	return m._rep.Favored(m._memberId, favType, referId)
}

// 取消收藏
func (m *favoriteManagerImpl) Cancel(favType int, referId int64) error {
	return m._rep.CancelFavorite(m._memberId, favType, referId)
}
