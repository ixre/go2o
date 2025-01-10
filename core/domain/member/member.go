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
	"net/url"
	"strconv"
	"strings"
	"time"

	de "github.com/ixre/go2o/core/domain/interface/domain"
	"github.com/ixre/go2o/core/domain/interface/domain/enum"
	"github.com/ixre/go2o/core/domain/interface/member"
	mss "github.com/ixre/go2o/core/domain/interface/message"
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/domain/interface/sys"
	"github.com/ixre/go2o/core/domain/interface/valueobject"
	"github.com/ixre/go2o/core/domain/interface/wallet"
	"github.com/ixre/go2o/core/event/events"
	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/go2o/core/infrastructure/fw/collections"
	"github.com/ixre/go2o/core/infrastructure/logger"
	"github.com/ixre/go2o/core/infrastructure/regex"
	"github.com/ixre/gof/domain/eventbus"
	"github.com/ixre/gof/util"
)

// todo: 依赖商户的 MSS 发送通知消息,应去掉
// todo: 会员升级 应单独来处理
var _ member.IMemberAggregateRoot = new(memberImpl)

type memberImpl struct {
	manager         member.IMemberManager
	value           *member.Member
	_extra          *member.ExtraField
	account         member.IAccount
	level           *member.Level
	repo            member.IMemberRepo
	walletRepo      wallet.IWalletRepo
	_systemRepo     sys.ISystemRepo
	relation        *member.InviteRelation
	invitation      member.IInvitationManager
	mssRepo         mss.IMessageRepo
	registryRepo    registry.IRegistryRepo
	valueRepo       valueobject.IValueRepo
	profileManager  member.IProfileManager
	favoriteManager member.IFavoriteManager
	giftCardManager member.IGiftCardManager
}

func NewMember(manager member.IMemberManager, val *member.Member,
	rep member.IMemberRepo, walletRepo wallet.IWalletRepo,
	mp mss.IMessageRepo, valRepo valueobject.IValueRepo,
	registryRepo registry.IRegistryRepo, systemRepo sys.ISystemRepo) member.IMemberAggregateRoot {
	return &memberImpl{
		manager:      manager,
		value:        val,
		repo:         rep,
		mssRepo:      mp,
		walletRepo:   walletRepo,
		valueRepo:    valRepo,
		registryRepo: registryRepo,
		_systemRepo:  systemRepo,
	}
}

// 获取聚合根编号
func (m *memberImpl) GetAggregateRootId() int {
	return m.value.Id
}

// Complex 会员汇总信息
func (m *memberImpl) Complex() *member.ComplexMember {
	mv := m.GetValue()
	lv := m.GetLevel()
	pf := m.Profile()
	tr := pf.GetCertificationInfo()
	extra := m.Extra()
	s := &member.ComplexMember{
		Nickname:            mv.Nickname,
		RealName:            mv.RealName,
		Avatar:              mv.ProfilePhoto,
		Exp:                 extra.Exp,
		Level:               mv.Level,
		LevelName:           lv.Name,
		TrustAuthState:      tr.ReviewStatus,
		TradePasswordHasSet: mv.TradePassword != "",
		PremiumUser:         mv.PremiumUser,
		Flag:                mv.UserFlag,
		UpdateTime:          int64(mv.UpdateTime),
	}
	return s
}

// Profile 会员资料服务
func (m *memberImpl) Profile() member.IProfileManager {
	if m.profileManager == nil {
		m.profileManager = newProfileManagerImpl(m,
			int64(m.GetAggregateRootId()), m.repo, m.registryRepo, m.valueRepo)
	}
	return m.profileManager
}

// Favorite 会员收藏服务
func (m *memberImpl) Favorite() member.IFavoriteManager {
	if m.favoriteManager == nil {
		m.favoriteManager = newFavoriteManagerImpl(
			int64(m.GetAggregateRootId()), m.repo)
	}
	return m.favoriteManager
}

// GiftCard 礼品卡服务
func (m *memberImpl) GiftCard() member.IGiftCardManager {
	if m.giftCardManager == nil {
		m.giftCardManager = newGiftCardManagerImpl(
			int64(m.GetAggregateRootId()), m.repo)
	}
	return m.giftCardManager
}

// Invitation 邀请管理
func (m *memberImpl) Invitation() member.IInvitationManager {
	if m.invitation == nil {
		m.invitation = &invitationManager{
			member:      m,
			_memberRepo: m.repo,
		}
	}
	return m.invitation
}

// GetValue 获取值
func (m *memberImpl) GetValue() member.Member {
	return *m.value
}

func (m *memberImpl) Extra() member.ExtraField {
	return *m.getExtra()
}

// getExtra 获取扩展字段
func (m *memberImpl) getExtra() *member.ExtraField {
	if m._extra == nil {
		m._extra = m.repo.ExtraRepo().FindBy("member_id=?", m.GetAggregateRootId())
		if m._extra == nil {
			m._extra = &member.ExtraField{
				Id:                 0,
				MemberId:           m.GetAggregateRootId(),
				Exp:                0,
				RegIp:              "",
				RegFrom:            "",
				RegTime:            0,
				CheckCode:          "",
				CheckExpires:       0,
				PersonalServiceUid: 0,
				LoginTime:          0,
				LastLoginTime:      0,
				UpdateTime:         0,
			}
		}
	}
	return m._extra
}

// SendCheckCode 发送验证码,并返回验证码
func (m *memberImpl) SendCheckCode(operation string, mssType int) (string, error) {
	extra := m.getExtra()
	const expiresMinutes = 10 //10分钟生效
	code := domain.NewCheckCode(4)
	extra.CheckCode = code
	extra.CheckExpires = int(time.Now().Add(time.Minute * expiresMinutes).Unix())
	_, err := m.repo.ExtraRepo().Save(extra)
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
		case mss.TypeEmail:
			n := mgr.GetNotifyItem("邮箱验证码")
			c := &mss.MailMessage{
				Subject: operation + "验证码",
				Body:    n.Content,
			}
			err = mgr.SendEmail(m.value.Email, c, data, "")
		case mss.TypeSMS:
			// 某些短信平台要求传入模板ID,在这里附加参数
			// re := m.registryRepo.Get(registry.SmsMemberCheckTemplateId)
			// data["templateId"] = re.StringValue()
			// 构造并发送短信
			n := mgr.GetNotifyItem("短信验证码")
			c := mss.PhoneMessage(n.Content)
			err = mgr.SendPhoneMessage(m.value.Phone, c, data, "")
		}
	}
	return code, err
}

// CompareCode 对比验证码
func (m *memberImpl) CompareCode(code string) error {
	extra := m.getExtra()
	if extra.CheckCode != strings.TrimSpace(code) {
		return de.ErrCheckCodeError
	}
	if extra.CheckExpires < int(time.Now().Unix()) {
		return de.ErrCheckCodeExpires
	}
	return nil
}

// GetAccount 获取账户
func (m *memberImpl) GetAccount() member.IAccount {
	if m.account == nil {
		v := m.repo.GetAccount(m.GetAggregateRootId())
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
	extra := m.getExtra()
	extra.Exp += exp
	_, err := m.repo.ExtraRepo().Save(extra)
	if err == nil {
		// 判断是否升级
		m.checkLevelUp()
	}
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
		m.value.PremiumExpires = int(expires)
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
	extra := m.getExtra()
	levelId := lg.GetLevelIdByExp(extra.Exp)
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
	m.value.UpdateTime = int(unix)
	_, err := m.Save()
	if err == nil {
		m.level = nil
		lvLog := &member.LevelUpLog{
			MemberId:     int(m.GetAggregateRootId()),
			OriginLevel:  origin,
			TargetLevel:  levelId,
			IsFree:       1,
			PaymentId:    0,
			ReviewStatus: int(enum.ReviewCompleted),
			UpgradeMode:  member.LAutoUpgrade,
			CreateTime:   int(unix),
		}
		_, _ = m.repo.SaveLevelUpLog(lvLog)
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
		MemberId:     int(m.GetAggregateRootId()),
		OriginLevel:  origin,
		TargetLevel:  level,
		PaymentId:    paymentId,
		ReviewStatus: int(enum.ReviewNone),
		UpgradeMode:  member.LServiceAgentUpgrade,
		CreateTime:   int(unix),
	}
	if paymentId == 0 {
		lvLog.IsFree = 1
	}
	if !review {
		lvLog.ReviewStatus = int(enum.ReviewCompleted)
	}
	_, err := m.repo.SaveLevelUpLog(lvLog)
	if err == nil && !review {
		if err = m.updateLevel(level); err != nil {
			return err
		}
		_, err = m.Save()
		if err == nil {
			m.level = nil
			// 更新经验值
			extra := m.getExtra()
			extra.Exp = lv.RequireExp
			extra.UpdateTime = int(unix)
			_, err = m.repo.ExtraRepo().Save(extra)
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
	if l != nil && l.MemberId == int(m.GetAggregateRootId()) {
		if l.ReviewStatus == int(enum.ReviewApproved) {
			return member.ErrLevelUpPass
		}
		if l.ReviewStatus == int(enum.ReviewRejected) {
			return member.ErrLevelUpReject
		}
		if l.ReviewStatus == int(enum.ReviewCompleted) {
			return member.ErrLevelUpConfirm
		}
		if time.Now().Unix()-int64(l.CreateTime) < 120 {
			return member.ErrLevelUpLaterConfirm
		}
		if pass {
			l.ReviewStatus = int(enum.ReviewApproved)
			_, err := m.repo.SaveLevelUpLog(l)
			if err == nil {
				lv := m.manager.LevelManager().GetLevelById(l.TargetLevel)
				if err = m.updateLevel(l.TargetLevel); err != nil {
					return err
				}
				_, err = m.Save()
				if err == nil {
					extra := m.getExtra()
					extra.Exp = lv.RequireExp
					extra.UpdateTime = int(time.Now().Unix())
					_, err = m.repo.ExtraRepo().Save(extra)
				}
			}
			return err
		} else {
			l.ReviewStatus = int(enum.ReviewRejected)
			_, err := m.repo.SaveLevelUpLog(l)
			return err
		}
	}
	return member.ErrNoSuchLevelUpLog

}

// ConfirmLevelUp 标记已经处理升级
func (m *memberImpl) ConfirmLevelUp(id int) error {
	l := m.repo.GetLevelUpLog(id)
	if l != nil && l.MemberId == int(m.GetAggregateRootId()) {
		if l.ReviewStatus == int(enum.ReviewCompleted) {
			return member.ErrLevelUpConfirm
		}
		if l.ReviewStatus != int(enum.ReviewApproved) {
			return member.ErrLevelUpReject
		}
		l.ReviewStatus = int(enum.ReviewCompleted)
		_, err := m.repo.SaveLevelUpLog(l)
		return err
	}
	return member.ErrNoSuchLevelUpLog
}

// GetRelation 获取会员关联
func (m *memberImpl) GetRelation() *member.InviteRelation {
	if m.relation == nil {
		rel := m.repo.GetRelation(int64(m.GetAggregateRootId()))
		if rel == nil {
			rel = &member.InviteRelation{
				MemberId:  int(m.GetAggregateRootId()),
				CardNo:    "",
				InviterId: 0,
				RegMchId:  0,
				InviterD2: 0,
				InviterD3: 0,
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
	extra := m.getExtra()
	extra.LastLoginTime = extra.LoginTime
	extra.LoginTime = int(unix)
	extra.UpdateTime = int(unix)
	_, err := m.repo.ExtraRepo().Save(extra)
	return err
}

// Save 保存
func (m *memberImpl) Save() (int64, error) {
	m.value.UpdateTime = int(time.Now().Unix()) // 更新时间，数据以更新时间触发
	if m.value.Id > 0 {
		_, err := m.repo.SaveMember(m.value)
		if err == nil {
			go m.pushSaveEvent(false)
		}
		return int64(m.value.Id), err
	}
	return 0, errors.New("member not registration")
}

func (m *memberImpl) pushSaveEvent(create bool) {
	rl := m.GetRelation()
	regFrom := ""
	if create {
		// 推送注册来源
		extra := m.getExtra()
		regFrom = extra.RegFrom
	}
	eventbus.Dispatch(&events.MemberPushEvent{
		IsCreate:  create,
		Member:    m.value,
		InviterId: int(rl.InviterId),
		RegFrom:   regFrom,
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
	return nil
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
			MemberId:   int64(m.value.Id),
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
		err = m.repo.DeleteLockInfos(int64(m.GetAggregateRootId()))
	}
	return err
}

// 根据注册来源计算会员角色身份
func (m *memberImpl) getUserRoleFlag(v *member.ExtraField) int {
	ret := member.RoleUser
	if len(v.RegFrom) != 0 {
		// 根据注册来源设置角色
		v.RegFrom = ""
		if strings.Contains(v.RegFrom, "EMPLOYEE") {
			// 商户职员
			ret |= member.RoleMchStaff
		}
		if strings.Contains(v.RegFrom, "EXT1") {
			// 扩展角色1
			ret |= member.RoleMchStaff
		}
		if strings.Contains(v.RegFrom, "EXT2") {
			// 扩展角色2
			ret |= member.RoleMchStaff
		}
	}
	return ret
}

func (m *memberImpl) updateLevel(levelId int) error {
	lm := m.manager.LevelManager()
	var level *member.Level
	if levelId <= 0 {
		level = lm.GetInitialLevel()
	} else {
		level = lm.GetLevelById(levelId)
	}
	if level == nil {
		return member.ErrLevelNotExist.Apply(levelId)
	}
	m.value.Level = level.Id
	if level.IsOfficial == 1 && !m.ContainFlag(member.FlagActive) {
		// 如果为正式等级,则必须激活
		return m.Active()
	}
	return nil
}

// 创建会员
func (m *memberImpl) SubmitRegistration(data *member.SubmitRegistrationData) error {
	v := m.value
	err := m.prepare()
	if err == nil {
		unix := time.Now().Unix()
		v.CreateTime = int(unix)
		// 初始化等级
		m.updateLevel(0)
		// 初始化国家信息
		if v.CountryCode == "" {
			v.CountryCode = "CN"
		}
		// 设置VIP用户信息
		v.PremiumUser = member.PremiumNormal
		v.PremiumExpires = 0
		// 创建一个用户编码/邀请码
		v.UserCode = m.generateMemberCode()
		id, err1 := m.repo.SaveMember(v)

		if err1 == nil {
			m.value.Id = int(id)
			if data.InviterId > 0 {
				// 如果邀请人编号存在,则绑定邀请人
				m.BindInviter(data.InviterId, true)
			}
			go m.memberInit(data)
		} else {
			err = err1
		}
	}
	return err
}

// 验证用户名
func (m *memberImpl) checkUser(user string) error {
	if len([]rune(user)) < 6 {
		return member.ErrUserLength
	}
	if !regex.IsUser(user) {
		return member.ErrUserValidErr
	}
	if m.repo.CheckUserExist(user, int64(m.GetAggregateRootId())) {
		return member.ErrUserExist
	}
	return nil
}

// 会员初始化
func (m *memberImpl) memberInit(data *member.SubmitRegistrationData) error {
	// 创建账户
	m.account = newAccount(m, &member.Account{},
		m.repo, m.manager, m.walletRepo, m.registryRepo)
	_, err := m.account.Save()
	if err != nil {
		return err
	}
	// 创建初始化数据
	unix := int(time.Now().Unix())
	extra := &member.ExtraField{
		Id:            0,
		MemberId:      int(m.GetAggregateRootId()),
		Exp:           0,
		RegIp:         data.RegIp,
		RegionCode:    0,
		RegFrom:       data.RegFrom,
		LoginTime:     unix,
		LastLoginTime: unix,
		UpdateTime:    unix,
	}
	_, err = m.repo.ExtraRepo().Save(extra)
	if err == nil {
		// 如果为中国大陆IP,则记录IP信息
		go m.updateRegionInfo(extra)
	} else {
		return err
	}
	// 注册后赠送积分
	regPresent := m.registryRepo.Get(registry.MemberRegisterPresentIntegral).IntValue()
	if regPresent > 0 {
		go m.GetAccount().CarryTo(member.AccountIntegral, member.AccountOperateData{
			TransactionTitle:   "新会员注册赠送积分",
			Amount:             regPresent,
			OuterTransactionNo: "-",
			TransactionRemark:  "sys",
		}, false, 0)
	}
	go m.pushSaveEvent(true)
	return nil
}

// 更新会员区域信息
func (m *memberImpl) updateRegionInfo(extra *member.ExtraField) {
	if m.value.CountryCode == "CN" && extra.RegIp != "" {
		// 如果为中国大陆IP,则记录IP信息
		ir := m._systemRepo.GetSystemAggregateRoot()
		region, err := ir.Location().FindRegionByIp(extra.RegIp)
		if err == nil && region != nil {
			extra.RegionCode = region.Code
			_, err := m.repo.ExtraRepo().Save(extra)
			if err != nil {
				logger.Error("更新会员区域信息失败! memberId:%d, ip:%s, error:%s",
					m.GetAggregateRootId(), extra.RegIp, err.Error())
			}
		}
	}
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
		if checkPhone && !regex.IsPhone(m.value.Phone) {
			return member.ErrInvalidPhone
		}
		if m.checkPhoneBind(m.value.Phone, int64(m.GetAggregateRootId())) != nil {
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
		m.value.Nickname = "用户" + util.RandString(4)
	}
	// 初始化头像
	m.value.ProfilePhoto = strings.TrimSpace(m.value.ProfilePhoto)
	if len(m.value.ProfilePhoto) == 0 {
		// 使用默认头像
		re, _ := m.registryRepo.GetValue(registry.MemberDefaultProfilePhoto)
		if len(strings.TrimSpace(re)) == 0 {
			// 如果未设置,则用系统内置头像
			prefix, _ := m.registryRepo.GetValue(registry.FileServerUrl)
			re, _ = url.JoinPath(prefix, "static/init/avatar.jpg")
		}
		m.value.ProfilePhoto = re
	}
	// 验证角色
	if m.value.RoleFlag != 0 && !collections.InArray([]int{
		member.RoleMchStaff,
		member.RoleExt1,
		member.RoleExt2,
		member.RoleExt3,
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
	if m.repo.CheckPhoneBind(phone, int(memberId)) {
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
func (m *memberImpl) BindInviter(inviterId int, force bool) (err error) {
	rl := m.GetRelation()
	if !force && rl.InviterId > 0 {
		return member.ErrExistsInviter
	}
	// 不能绑定自己为推荐人
	if int(m.GetAggregateRootId()) == inviterId {
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

// 绑定邀请人,如果已邀请,force为true时更新
func (m *memberImpl) BindMerchantId(mchId int, force bool) (err error) {
	rl := m.GetRelation()
	if !force && rl.RegMchId > 0 {
		return errors.New("商户已绑定")
	}
	if rl.RegMchId != mchId {
		rl.RegMchId = mchId
		return m.repo.SaveRelation(rl)
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
