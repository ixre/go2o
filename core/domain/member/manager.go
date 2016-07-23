/**
 * Copyright 2015 @ z3q.net.
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
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/variable"
	"sort"
	"strings"
)

var _ member.IMemberManager = new(MemberManagerImpl)
var _ member.ILevelManager = new(levelManagerImpl)

type MemberManagerImpl struct {
	_levelManager member.ILevelManager
	_valRep       valueobject.IValueRep
	_rep          member.IMemberRep
}

func NewMemberManager(rep member.IMemberRep,
	valRep valueobject.IValueRep) member.IMemberManager {
	return &MemberManagerImpl{
		_levelManager: newLevelManager(rep),
		_valRep:       valRep,
		_rep:          rep,
	}
}

// 等级服务
func (l *MemberManagerImpl) LevelManager() member.ILevelManager {
	return l._levelManager
}

// 检测注册权限
func (l *MemberManagerImpl) RegisterPerm(invitation bool) error {
	conf := l._valRep.GetRegisterPerm()
	if conf.RegisterMode == member.RegisterModeClosed {
		return member.ErrRegOff
	}
	if conf.RegisterMode == member.RegisterModeMustInvitation && !invitation {
		return member.ErrRegMustInvitation
	}
	if conf.RegisterMode == member.RegisterModeMustRedirect && invitation {
		return member.ErrRegOffInvitation
	} else if conf.RegisterMode == member.RegisterModeNormal {

	}
	return nil
}

func (l *MemberManagerImpl) checkInvitationCode(invitationCode string) (int, error) {
	var invitationId int = 0
	if len(invitationCode) > 0 {
		//判断邀请码是否正确
		invitationId = l._rep.GetMemberIdByInvitationCode(invitationCode)
		if invitationId <= 0 {
			return -1, member.ErrInvitationCode
		}
	}
	return invitationId, nil
}

// 检查手机绑定,同时检查手机格式
func (l *MemberManagerImpl) CheckPhoneBind(phone string, memberId int) error {
	if len(phone) <= 0 {
		return member.ErrMissingPhone
	}
	if !phoneRegex.MatchString(phone) {
		return member.ErrBadPhoneFormat
	}
	if b := l._rep.CheckPhoneBind(phone, memberId); b {
		return member.ErrPhoneHasBind
	}
	return nil
}

// 检查注册信息是否正确
func (l *MemberManagerImpl) CheckPostedRegisterInfo(v *member.Member,
	pro *member.Profile, invitationCode string) (invitationId int, err error) {
	perm := l._valRep.GetRegisterPerm()

	//验证用户名
	v.Usr = strings.TrimSpace(v.Usr)
	if len(v.Usr) < 6 {
		return 0, member.ErrUsrLength
	}
	if !userRegex.MatchString(v.Usr) {
		return 0, member.ErrUsrValidErr
	}
	if l._rep.CheckUsrExist(v.Usr, 0) {
		return 0, member.ErrUsrExist
	}

	//验证密码
	v.Pwd = strings.TrimSpace(v.Pwd)
	if len(v.Pwd) < 6 {
		return 0, member.ErrPwdLength
	}

	//验证手机
	pro.Phone = strings.TrimSpace(pro.Phone)
	lp := len(pro.Phone)
	if perm.NeedPhone && lp == 0 {
		return 0, member.ErrMissingPhone
	}
	if lp > 0 {
		if !phoneRegex.MatchString(pro.Phone) {
			return 0, member.ErrBadPhoneFormat
		}
		if b := l._rep.CheckPhoneBind(pro.Phone, v.Id); b {
			return 0, member.ErrPhoneHasBind
		}
	}

	//验证IM
	pro.Im = strings.TrimSpace(pro.Im)
	if perm.NeedIm && len(pro.Im) == 0 {
		return 0, errors.New(strings.Replace(member.ErrMissingIM.Error(),
			"IM", variable.AliasMemberIM, -1))
	}

	// 检查验证码
	err = l.RegisterPerm(len(invitationCode) > 0)
	if err == nil {
		invitationId, err = l.checkInvitationCode(invitationCode)
	}
	return invitationId, err
}

// 等级服务实现
type levelManagerImpl struct {
	_rep    member.IMemberRep
	_levels []*member.Level //可用的等级
}

func newLevelManager(rep member.IMemberRep) member.ILevelManager {
	impl := &levelManagerImpl{
		_rep: rep,
	}
	return impl.init()
}

// 初始化默认等级
func (l *levelManagerImpl) init() member.ILevelManager {
	if len(l.GetLevelSet()) == 0 {
		l._levels = []*member.Level{
			{
				Name:          "待激活会员",
				RequireExp:    0,
				Enabled:       1,
				ProgramSignal: "M_DJH",
				IsOfficial:    0,
			},
			{
				Name:          "普通会员",
				RequireExp:    1,
				Enabled:       1,
				ProgramSignal: "M_PT",
				IsOfficial:    1,
			},
			{
				Name:          "铜牌会员",
				RequireExp:    100,
				Enabled:       1,
				ProgramSignal: "M_TP",
				IsOfficial:    1,
			},
			{
				Name:          "银牌会员",
				RequireExp:    500,
				Enabled:       1,
				ProgramSignal: "M_YP",
				IsOfficial:    1,
			},
			{
				Name:          "金牌会员",
				RequireExp:    1200,
				Enabled:       1,
				ProgramSignal: "M_JP",
				IsOfficial:    1,
			},
			{
				Name:          "白金会员",
				RequireExp:    1500,
				Enabled:       1,
				ProgramSignal: "M_BJ",
				IsOfficial:    1,
			},
		}
		// 存储并设置编号
		for _, v := range l._levels {
			v.Id, _ = l.SaveLevel(v)
		}
	}
	return l
}

// 获取等级设置
func (l *levelManagerImpl) GetLevelSet() []*member.Level {
	if l._levels == nil {
		// 已经排好序
		l._levels = l._rep.GetMemberLevels_New()
	}
	return l._levels
}

// 获取等级
func (l *levelManagerImpl) GetLevelById(id int) *member.Level {
	if id == 0 {
		return nil
	}
	arr := l.GetLevelSet()
	if la := len(arr); la > 0 {
		i := sort.Search(la, func(i int) bool {
			return arr[i].Id >= id
		})
		if i < la && arr[i].Id == id {
			return arr[i]
		}
	}
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
			return arr[i].Id >= id
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
	pos := 0
	for i, v := range l.GetLevelSet() {
		if v.Id == id {
			pos = i
			break
		}
	}
	if pos > 0 {
		// 获取等级对应的会员数, 如果 > 0不允许删除
		// todo: 也可以更新到下一个等级
		if n := l._rep.GetMemberNumByLevel_New(id); n > 0 {
			return member.ErrLevelUsed
		}
		l._levels = append(l._levels[:pos],
			l._levels[pos+1:]...)
		return l._rep.DeleteMemberLevel_New(id)
	}
	return nil
}

// 保存等级
func (l *levelManagerImpl) SaveLevel(v *member.Level) (int, error) {
	v.ProgramSignal = strings.TrimSpace(v.ProgramSignal)
	if !l.checkProgramSignal(v.ProgramSignal, v.Id) {
		return -1, member.ErrExistsSameProgramSignalLevel
	}
	err := l.checkLevelExp(v)
	if err == nil {
		v.Id, err = l._rep.SaveMemberLevel_New(v)
		if err == nil {
			l._levels = nil
		}
	}
	return v.Id, err
}

// 判断等级与等级可编程签名是否一致
func (l *levelManagerImpl) checkProgramSignal(sign string, id int) bool {
	if sign != "" {
		for _, v := range l.GetLevelSet() {
			if v.ProgramSignal == sign {
				return id == v.Id
			}
		}
	}
	return true
}

// 新增等级时检查经验值
func (m *levelManagerImpl) checkLevelExp(lv *member.Level) error {
	// 新增判断经验值
	if lv.Id <= 0 {
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
		if lv.Id > 0 && v.Id == lv.Id {
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
		} else if v.Id > lv.Id {
			lv = v
		}
	}
	return lv
}

// 获取最大的等级值
func (l *levelManagerImpl) getMaxLevelId() int {
	var k = 0
	for _, v := range l.GetLevelSet() {
		if v.Id > k {
			k = v.Id
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
		if exp >= lv.RequireExp && lv.Id > levelVal {
			levelVal = lv.Id
		}
	}
	return levelVal
}
