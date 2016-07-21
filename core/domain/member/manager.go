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
	"github.com/jsix/gof/log"
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
func (this *MemberManagerImpl) LevelManager() member.ILevelManager {
	return this._levelManager
}

// 检测注册权限
func (this *MemberManagerImpl) RegisterPerm(invitation bool) error {
	conf := this._valRep.GetRegisterPerm()
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

func (this *MemberManagerImpl) checkInvitationCode(invitationCode string) (int, error) {
	var invitationId int = 0
	if len(invitationCode) > 0 {
		//判断邀请码是否正确
		invitationId = this._rep.GetMemberIdByInvitationCode(invitationCode)
		if invitationId <= 0 {
			return -1, member.ErrInvitationCode
		}
	}
	return invitationId, nil
}

// 检查手机绑定,同时检查手机格式
func (this *MemberManagerImpl) CheckPhoneBind(phone string, memberId int) error {
	if len(phone) <= 0 {
		return member.ErrMissingPhone
	}
	if !phoneRegex.MatchString(phone) {
		return member.ErrBadPhoneFormat
	}
	if b := this._rep.CheckPhoneBind(phone, memberId); b {
		return member.ErrPhoneHasBind
	}
	return nil
}

// 检查注册信息是否正确
func (this *MemberManagerImpl) CheckPostedRegisterInfo(v *member.Member,
	pro *member.Profile, invitationCode string) (invitationId int, err error) {
	perm := this._valRep.GetRegisterPerm()

	//验证用户名
	v.Usr = strings.TrimSpace(v.Usr)
	if len(v.Usr) < 6 {
		return 0, member.ErrUsrLength
	}
	if !userRegex.MatchString(v.Usr) {
		return 0, member.ErrUsrValidErr
	}
	if this._rep.CheckUsrExist(v.Usr, 0) {
		return 0, member.ErrUsrExist
	}

	//验证密码
	v.Pwd = strings.TrimSpace(v.Pwd)
	if len(v.Pwd) < 6 {
		return 0, member.ErrPwdLength
	}

	//验证手机
	pro.Phone = strings.TrimSpace(pro.Phone)
	l := len(pro.Phone)
	if perm.NeedPhone && l == 0 {
		return 0, member.ErrMissingPhone
	}
	if l > 0 {
		if !phoneRegex.MatchString(pro.Phone) {
			return 0, member.ErrBadPhoneFormat
		}
		if b := this._rep.CheckPhoneBind(pro.Phone, v.Id); b {
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
	err = this.RegisterPerm(len(invitationCode) > 0)
	if err == nil {
		invitationId, err = this.checkInvitationCode(invitationCode)
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
		log.Println("----- errr    ")
		l._levels = []*member.Level{
			{
				Name:       "待激活会员",
				RequireExp: 0,
				Enabled:    1,
				IsOfficial: 0,
			},
			{
				Name:       "普通会员",
				RequireExp: 1,
				Enabled:    1,
				IsOfficial: 1,
			},
			{
				Name:       "铜牌会员",
				RequireExp: 100,
				Enabled:    1,
				IsOfficial: 1,
			},
			{
				Name:       "银牌会员",
				RequireExp: 500,
				Enabled:    1,
				IsOfficial: 1,
			},
			{
				Name:       "金牌会员",
				RequireExp: 1200,
				Enabled:    1,
				IsOfficial: 1,
			},
			{
				Name:       "白金会员",
				RequireExp: 1500,
				Enabled:    1,
				IsOfficial: 1,
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

// 获取下一个等级
func (l *levelManagerImpl) GetNextLevelById(id int) *member.Level {
	arr := l.GetLevelSet()
	i := sort.Search(len(arr), func(i int) bool {
		return arr[i].Id >= id
	})
	if i < len(arr)-1 {
		return arr[i+1]
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
	// 如果新增（非初始化）等级自动设置值
	//if v.Id <= 0 && len(this._levelSet) == 0 {
	//    v.Value = this.getMaxLevelValue() + 1
	//}
	l._levels = nil
	if err := l.checkNewLevel(v); err != nil {
		return -1, err
	}
	return l._rep.SaveMemberLevel_New(v)
}

// 新增等级时检查经验值
func (this *levelManagerImpl) checkNewLevel(v *member.Level) error {
	if v.Id <= 0 {
		max := this.getMaxLevelId()
		lv := this.GetLevelById(max)
		if lv != nil && v.RequireExp < lv.RequireExp {
			return errors.New(fmt.Sprintf(
				member.ErrLevelRequireExp.Error(), lv.RequireExp))
		}
	}
	return nil
}

// 获取最大的等级值
func (this *levelManagerImpl) getMaxLevelId() int {
	var k = 0
	for _, v := range this.GetLevelSet() {
		if v.Id > k {
			k = v.Id
		}
	}
	return k
}

// 根据经验值获取等级
func (this *levelManagerImpl) GetLevelIdByExp(exp int) int {
	var lv *member.Level
	var levelVal int
	for i := len(this.GetLevelSet()); i > 0; i-- {
		lv = this.GetLevelSet()[i-1]
		if exp >= lv.RequireExp && lv.Id > levelVal {
			levelVal = lv.Id
		}
	}
	return levelVal
}
