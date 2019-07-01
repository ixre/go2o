/**
 * Copyright 2015 @ to2.net.
 * name : conf_manager
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package member

import (
	"errors"
	"fmt"
	"github.com/ixre/gof/util"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/registry"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/variable"
	"sort"
	"strings"
)

var _ member.IMemberManager = new(memberManagerImpl)
var _ member.ILevelManager = new(levelManagerImpl)

type memberManagerImpl struct {
	levelManager member.ILevelManager
	valRepo      valueobject.IValueRepo
	rep          member.IMemberRepo
	registryRepo registry.IRegistryRepo
}

func NewMemberManager(rep member.IMemberRepo,
	valRepo valueobject.IValueRepo, registryRepo registry.IRegistryRepo) member.IMemberManager {
	return &memberManagerImpl{
		levelManager: newLevelManager(rep),
		valRepo:      valRepo,
		rep:          rep,
		registryRepo: registryRepo,
	}
}

// 等级服务
func (m *memberManagerImpl) LevelManager() member.ILevelManager {
	return m.levelManager
}

// 检测注册权限
func (m *memberManagerImpl) registerPerm(regMode int, invitation bool) error {
	if regMode== member.RegisterModeClosed {
		return member.ErrRegOff
	}
	if regMode == member.RegisterModeMustInvitation && !invitation {
		return member.ErrRegMissingInvitationCode
	}
	if regMode == member.RegisterModeMustRedirect && invitation {
		return member.ErrRegOffInvitation
	}
	return nil
}

func (m *memberManagerImpl) checkInvitationCode(invitationCode string) (int64, error) {
	var invitationId int64
	if len(invitationCode) > 0 {
		//判断邀请码是否正确
		invitationId = m.rep.GetMemberIdByInvitationCode(invitationCode)
		if invitationId <= 0 {
			return -1, member.ErrInvitationCode
		}
	}
	return invitationId, nil
}

// 检查手机绑定,同时检查手机格式
func (m *memberManagerImpl) CheckPhoneBind(phone string, memberId int64) error {
	if len(phone) <= 0 {
		return member.ErrMissingPhone
	}
	if m.rep.CheckPhoneBind(phone, memberId) {
		return member.ErrPhoneHasBind
	}
	return nil
}

// 检查邀请注册
func (m *memberManagerImpl) CheckInviteRegister(code string) (inviterId int64, err error) {
	regMode := m.registryRepo.Get(registry.MemberRegisterMode).IntValue()
	isInvite := len(code) > 0
	// 检查注册模式及邀请码
	err = m.registerPerm(regMode, isInvite)
	if err == nil && isInvite {
		//判断邀请码是否正确
		inviterId = m.rep.GetMemberIdByInvitationCode(code)
		if inviterId <= 0 {
			return 0, member.ErrInvitationCode
		}
		return inviterId, nil
	}
	return 0, err
}

// 检查注册信息是否正确
func (m *memberManagerImpl) PrepareRegister(v *member.Member,
	pro *member.Profile, invitationCode string) (invitationId int64, err error) {

	phoneAsUser := m.registryRepo.Get(registry.MemberRegisterPhoneAsUser).BoolValue()
	mustBindPhone := m.registryRepo.Get(registry.MemberRegisterMustBindPhone).BoolValue()
	needIm := m.registryRepo.Get(registry.MemberRegisterNeedIm).BoolValue()
	// 验证用户名,如果填写了或非用手机号作为用户名,均验证用户名
	v.User = strings.TrimSpace(v.User)
	if v.User != "" || !phoneAsUser {
		if len(v.User) < 6 {
			return 0, member.ErrUsrLength
		}
		if !userRegex.MatchString(v.User) {
			return 0, member.ErrUsrValidErr
		}
		if m.rep.CheckUsrExist(v.User, 0) {
			return 0, member.ErrUsrExist
		}
	}

	// 验证密码
	v.Pwd = strings.TrimSpace(v.Pwd)
	if len(v.Pwd) < 6 {
		return 0, member.ErrPwdLength
	}

	// 验证手机
	pro.Phone = strings.TrimSpace(pro.Phone)
	lp := len(pro.Phone)
	if mustBindPhone && lp == 0 {
		return 0, member.ErrMissingPhone
	}
	if lp > 0 {
		checkPhone := m.registryRepo.Get(registry.MemberCheckPhoneFormat).BoolValue()
		if checkPhone && !phoneRegex.MatchString(pro.Phone) {
			return 0, member.ErrBadPhoneFormat
		}
		if m.CheckPhoneBind(pro.Phone, v.Id) != nil {
			return 0, member.ErrPhoneHasBind
		}
	}

	// 使用手机号作为用户名
	if phoneAsUser {
		if m.rep.CheckUsrExist(pro.Phone, 0) {
			return 0, member.ErrPhoneHasBind
		}
		v.User = pro.Phone
	}

	// 验证IM
	pro.Im = strings.TrimSpace(pro.Im)
	if needIm && len(pro.Im) == 0 {
		return 0, errors.New(strings.Replace(member.ErrMissingIM.Error(),
			"IM", variable.AliasMemberIM, -1))
	}

	pro.Name = strings.TrimSpace(pro.Name)
	pro.Avatar = strings.TrimSpace(pro.Avatar)
	if len(pro.Name) == 0 {
		//如果未设置昵称,则默认为用户名
		pro.Name = v.User
	}
	if len(pro.Avatar) == 0 {
		pro.Avatar = "res/no_avatar.gif"
	}
	return invitationId, err
}

// 获取所有买家分组
func (m *memberManagerImpl) GetAllBuyerGroups() []*member.BuyerGroup {
	list := m.rep.SelectMmBuyerGroup("")
	if len(list) == 0 {
		m.initBuyerGroups()
		list = m.rep.SelectMmBuyerGroup("")
	}
	return list
}

// 初始化买家分组
func (m *memberManagerImpl) initBuyerGroups() {
	arr := []*member.BuyerGroup{
		{
			Name:      "默认分组",
			IsDefault: 1,
		},
		{
			Name: "自定义分组1",
		},
		{
			Name: "自定义分组2",
		},
		{
			Name: "自定义分组3",
		},
		{
			Name: "自定义分组4",
		},
		{
			Name: "自定义分组5",
		},
		{
			Name: "自定义分组6",
		},
		{
			Name: "自定义分组7",
		},
		{
			Name: "自定义分组8",
		},
		{
			Name: "自定义分组9",
		},
	}
	for _, v := range arr {
		m.rep.SaveMmBuyerGroup(v)
	}
}

// 获取买家分组
func (m *memberManagerImpl) GetBuyerGroup(id int32) *member.BuyerGroup {
	for _, v := range m.GetAllBuyerGroups() {
		if v.ID == id {
			return v
		}
	}
	return nil
}

// 保存买家分组
func (m *memberManagerImpl) SaveBuyerGroup(v *member.BuyerGroup) (int32, error) {
	return util.I32Err(m.rep.SaveMmBuyerGroup(v))
}

// 等级服务实现
type levelManagerImpl struct {
	rep member.IMemberRepo
}

func newLevelManager(rep member.IMemberRepo) member.ILevelManager {
	impl := &levelManagerImpl{
		rep: rep,
	}
	return impl.init()
}

// 初始化默认等级
func (l *levelManagerImpl) init() member.ILevelManager {
	if len(l.GetLevelSet()) == 0 {
		levels := []*member.Level{
			{
				Name:          "普通会员",
				RequireExp:    0,
				Enabled:       1,
				AllowUpgrade:  1,
				ProgramSignal: "LV.1",
				IsOfficial:    0,
			},
			{
				Name:          "铜牌会员",
				RequireExp:    100,
				Enabled:       1,
				AllowUpgrade:  1,
				ProgramSignal: "LV.2",
				IsOfficial:    1,
			},
			{
				Name:          "银牌会员",
				RequireExp:    500,
				Enabled:       1,
				AllowUpgrade:  1,
				ProgramSignal: "LV.3",
				IsOfficial:    1,
			},
			{
				Name:          "金牌会员",
				RequireExp:    1200,
				Enabled:       1,
				ProgramSignal: "LV.4",
				AllowUpgrade:  1,
				IsOfficial:    1,
			},
			{
				Name:          "钻石会员",
				RequireExp:    1200,
				Enabled:       1,
				ProgramSignal: "LV.5",
				AllowUpgrade:  1,
				IsOfficial:    1,
			},
		}
		// 存储并设置编号
		for _, v := range levels {
			v.ID, _ = l.SaveLevel(v)
		}
	}
	return l
}

// 获取等级设置
func (l *levelManagerImpl) GetLevelSet() []*member.Level {
	return l.rep.GetMemberLevels_New()
}

// 获取等级
func (l *levelManagerImpl) GetLevelById(id int) *member.Level {
	if id == 0 {
		return nil
	}
	arr := l.GetLevelSet()
	if la := len(arr); la > 0 {
		i := sort.Search(la, func(i int) bool {
			return arr[i].ID >= id
		})
		if i < la && arr[i].ID == id {
			return arr[i]
		}
	}
	println(fmt.Sprintf("level = %#v", arr))
	panic(errors.New(fmt.Sprintf("no such member level id as %d", id)))
}

// 根据可编程字符获取会员等级
func (l *levelManagerImpl) GetLevelByProgramSign(sign string) *member.Level {
	for _, v := range l.GetLevelSet() {
		if v.ProgramSignal == sign {
			return v
		}
	}
	return nil
}

// 获取下一个等级
func (l *levelManagerImpl) GetNextLevelById(id int) *member.Level {
	arr := l.GetLevelSet()
	if la := len(arr); la > 0 {
		i := sort.Search(la, func(i int) bool {
			return arr[i].ID >= id
		})
		// 获取一下个等级,如果等级未启用,则升级下一个等级
		for j := 1; j < la-i; j++ {
			a := arr[i+j]
			if a.Enabled == 1 {
				return a
			}
		}
	}
	return nil //已经是最高级
}

// 删除等级
func (l *levelManagerImpl) DeleteLevel(id int) error {
	lv := l.GetLevelById(id)
	if lv != nil {
		// 获取等级对应的会员数, 如果 > 0不允许删除
		// todo: 也可以更新到下一个等级
		if n := l.rep.GetMemberNumByLevel_New(id); n > 0 {
			return member.ErrLevelUsed
		}
		return l.rep.DeleteMemberLevel_New(id)
	}
	return nil
}

// 保存等级
func (l *levelManagerImpl) SaveLevel(v *member.Level) (int, error) {
	v.ProgramSignal = strings.TrimSpace(v.ProgramSignal)
	if !l.checkProgramSignal(v.ProgramSignal, v.ID) {
		return -1, member.ErrExistsSameProgramSignalLevel
	}
	err := l.checkLevelExp(v)
	if err == nil {
		return l.rep.SaveMemberLevel_New(v)
	}
	return v.ID, err
}

// 判断等级与等级可编程签名是否一致
func (l *levelManagerImpl) checkProgramSignal(sign string, id int) bool {
	if sign != "" {
		for _, v := range l.GetLevelSet() {
			if v.ProgramSignal == sign {
				return id == v.ID
			}
		}
	}
	return true
}

// 新增等级时检查经验值
func (m *levelManagerImpl) checkLevelExp(lv *member.Level) error {
	// 新增判断经验值
	if lv.ID <= 0 {
		max := m.getMaxLevelId()
		lvMax := m.GetLevelById(max)
		if lvMax != nil && lv.RequireExp < lvMax.RequireExp {
			return member.ErrMustMoreThanMaxLevel
		}
		return nil
	}

	// 保存时检查经验值,必须大于前一个等级或小于后一个等级
	arr := m.GetLevelSet()
	la := len(arr)
	for i, v := range arr {
		// 如果为保存等级
		if lv.ID > 0 && v.ID == lv.ID {
			err := m.checkBetweenRequireExp(arr, i, la, lv.RequireExp)
			if err != nil {
				return err
			}
			break
		}
	}

	return nil
}

// 检查保存等级区间经验值
func (l *levelManagerImpl) checkBetweenRequireExp(arr []*member.Level,
	i int, la int, exp int) error {
	// 如果小于前一个等级
	if i > 0 && arr[i-1].RequireExp > exp {
		return member.ErrLessThanLevelRequireExp
	}
	// 如果大于后一个等级
	if i < la-1 && arr[i+1].RequireExp < exp {
		return member.ErrMoreThanLevelRequireExp
	}
	return nil
}

// 获取最高已启用的等级
func (l *levelManagerImpl) GetHighestLevel() *member.Level {
	var lv *member.Level
	for _, v := range l.GetLevelSet() {
		if v.Enabled != 1 {
			continue
		}
		if lv == nil {
			lv = v
		} else if v.ID > lv.ID {
			lv = v
		}
	}
	return lv
}

// 获取最大的等级值
func (l *levelManagerImpl) getMaxLevelId() int {
	var k int
	for _, v := range l.GetLevelSet() {
		if v.ID > k {
			k = v.ID
		}
	}
	return k
}

// 根据经验值获取等级
func (l *levelManagerImpl) GetLevelIdByExp(exp int) int {
	var lv *member.Level
	var levelVal int
	arr := l.GetLevelSet()
	for i := len(arr); i > 0; i-- {
		lv = arr[i-1]
		if exp >= lv.RequireExp && lv.ID > levelVal {
			levelVal = lv.ID
		}
	}
	return levelVal
}
