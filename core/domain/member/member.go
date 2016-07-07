/**
 * Copyright 2014 @ z3q.net.
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
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/mss"
	"go2o/core/domain/interface/mss/notify"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/infrastructure/domain"
	"go2o/core/infrastructure/tool/sms"
	"regexp"
	"strings"
	"time"
)

//todo: 依赖商户的 MSS 发送通知消息,应去掉
var _ member.IMember = new(memberImpl)

type memberImpl struct {
	_manager         member.IMemberManager
	_value           *member.Member
	_account         member.IAccount
	_level           *member.Level
	_rep             member.IMemberRep
	_merchantRep     merchant.IMerchantRep
	_relation        *member.Relation
	_invitation      member.IInvitationManager
	_mssRep          mss.IMssRep
	_valRep          valueobject.IValueRep
	_profileManager  member.IProfileManager
	_favoriteManager member.IFavoriteManager
	_giftCardManager member.IGiftCardManager
}

func NewMember(manager member.IMemberManager, val *member.Member, rep member.IMemberRep,
	mp mss.IMssRep, valRep valueobject.IValueRep, merchantRep merchant.IMerchantRep) member.IMember {
	return &memberImpl{
		_manager:     manager,
		_value:       val,
		_rep:         rep,
		_mssRep:      mp,
		_valRep:      valRep,
		_merchantRep: merchantRep,
	}
}

// 获取聚合根编号
func (this *memberImpl) GetAggregateRootId() int {
	return this._value.Id
}

// 会员资料服务
func (this *memberImpl) Profile() member.IProfileManager {
	if this._profileManager == nil {
		this._profileManager = newProfileManagerImpl(this,
			this.GetAggregateRootId(), this._rep, this._valRep)
	}
	return this._profileManager
}

// 会员收藏服务
func (this *memberImpl) Favorite() member.IFavoriteManager {
	if this._favoriteManager == nil {
		this._favoriteManager = newFavoriteManagerImpl(
			this.GetAggregateRootId(), this._rep)
	}
	return this._favoriteManager
}

// 礼品卡服务
func (this *memberImpl) GiftCard() member.IGiftCardManager {
	if this._giftCardManager == nil {
		this._giftCardManager = newGiftCardManagerImpl(
			this.GetAggregateRootId(), this._rep)
	}
	return this._giftCardManager
}

// 邀请管理
func (this *memberImpl) Invitation() member.IInvitationManager {
	if this._invitation == nil {
		this._invitation = &invitationManager{
			_member: this,
		}
	}
	return this._invitation
}

// 获取值
func (this *memberImpl) GetValue() member.Member {
	return *this._value
}

var (
	userRegex  = regexp.MustCompile("^[a-zA-Z0-9_]{6,}$")
	emailRegex = regexp.MustCompile("\\w+([-+.']\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*")
	phoneRegex = regexp.MustCompile("^(13[0-9]|15[0|1|2|3|4|5|6|8|9]|18[0|1|2|3|5|6|7|8|9]|17[0|6|7|8]|14[7])(\\d{8})$")
)

func (this *memberImpl) validate(v *member.Member) error {
	v.Usr = strings.ToLower(strings.TrimSpace(v.Usr)) // 小写并删除空格
	if len([]rune(v.Usr)) < 6 {
		return member.ErrUsrLength
	}
	if !userRegex.MatchString(v.Usr) {
		return member.ErrUsrValidErr
	}
	return nil
}

// 设置值
func (this *memberImpl) SetValue(v *member.Member) error {
	v.Usr = this._value.Usr
	if len(this._value.InvitationCode) == 0 {
		this._value.InvitationCode = v.InvitationCode
	}
	if v.Exp != 0 {
		this._value.Exp = v.Exp
	}
	if v.Level > 0 {
		this._value.Level = v.Level
	}
	if len(v.TradePwd) == 0 {
		this._value.TradePwd = v.TradePwd
	}
	return nil
}

// 发送验证码,并返回验证码
func (this *memberImpl) SendCheckCode(operation string, mssType int) (string, error) {
	const expiresMinutes = 10 //10分钟生效
	code := domain.NewCheckCode()
	this._value.CheckCode = code
	this._value.CheckExpires = time.Now().Add(time.Minute * expiresMinutes).Unix()
	_, err := this.Save()
	if err == nil {
		mgr := this._mssRep.NotifyManager()
		pro := this.Profile().GetProfile()

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
			provider, _ := this._valRep.GetDefaultSmsApiPerm()
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
func (this *memberImpl) CompareCode(code string) error {
	if this._value.CheckCode != strings.TrimSpace(code) {
		return member.ErrCheckCodeError
	}
	if this._value.CheckExpires < time.Now().Unix() {
		return member.ErrCheckCodeExpires
	}
	return nil
}

// 获取账户
func (this *memberImpl) GetAccount() member.IAccount {
	if this._account == nil {
		v := this._rep.GetAccount(this._value.Id)
		return NewAccount(v, this._rep)
	}
	return this._account
}

// 增加经验值
func (this *memberImpl) AddExp(exp int) error {
	this._value.Exp += exp
	_, err := this.Save()
	//判断是否升级
	this.checkUpLevel()

	return err
}

// 获取等级
func (this *memberImpl) GetLevel() *member.Level {
	if this._level == nil {
		this._level = this._manager.LevelManager().
			GetLevelById(this._value.Level)
	}
	return this._level
}

// 检查升级
func (this *memberImpl) checkUpLevel() bool {
	lg := this._manager.LevelManager()
	levelId := lg.GetLevelIdByExp(this._value.Exp)
	if levelId != 0 && this._value.Level < levelId {
		this._value.Level = levelId
		this.Save()
		this._level = nil
		return true
	}
	return false
}

// 获取会员关联
func (this *memberImpl) GetRelation() *member.Relation {
	if this._relation == nil {
		this._relation = this._rep.GetRelation(this._value.Id)
	}
	return this._relation
}

// 更换用户名
func (this *memberImpl) ChangeUsr(usr string) error {
	if usr == this._value.Usr {
		return member.ErrSameUsr
	}
	if len([]rune(usr)) < 6 {
		return member.ErrUsrLength
	}
	if !userRegex.MatchString(usr) {
		return member.ErrUsrValidErr
	}
	if this.usrIsExist(usr) {
		return member.ErrUsrExist
	}
	this._value.Usr = usr
	_, err := this.Save()
	return err
}

// 保存
func (this *memberImpl) Save() (int, error) {
	this._value.UpdateTime = time.Now().Unix() // 更新时间，数据以更新时间触发
	if this._value.Id > 0 {
		return this._rep.SaveMember(this._value)
	}

	if err := this.validate(this._value); err != nil {
		return this.GetAggregateRootId(), err
	}
	return this.create(this._value, nil)
}

// 锁定会员
func (this *memberImpl) Lock() error {
	return this._rep.LockMember(this.GetAggregateRootId(), 0)
}

// 解锁会员
func (this *memberImpl) Unlock() error {
	return this._rep.LockMember(this.GetAggregateRootId(), 1)
}

// 创建会员
func (this *memberImpl) create(m *member.Member, pro *member.Profile) (int, error) {
	//todo: 获取推荐人编号
	//todo: 检测是否有注册权限
	//if err := this._manager.RegisterPerm(this._relation.RefereesId);err != nil{
	//	return -1,err
	//}
	if this.usrIsExist(m.Usr) {
		return -1, member.ErrUsrExist
	}

	t := time.Now().Unix()
	m.State = 1
	m.RegTime = t
	m.LastLoginTime = t
	m.Level = 1
	m.Exp = 1
	m.DynamicToken = m.Pwd
	m.Exp = 0
	if len(m.RegFrom) == 0 {
		m.RegFrom = "API-INTERNAL"
	}
	m.InvitationCode = this.generateInvitationCode() // 创建一个邀请码
	id, err := this._rep.SaveMember(m)
	if id != 0 {
		this._value.Id = id
	}
	return id, err
}

// 创建邀请码
func (this *memberImpl) generateInvitationCode() string {
	var code string
	for {
		code = domain.GenerateInvitationCode()
		if memberId := this._rep.GetMemberIdByInvitationCode(code); memberId == 0 {
			break
		}
	}
	return code
}

// 用户是否已经存在
func (this *memberImpl) usrIsExist(usr string) bool {
	return this._rep.CheckUsrExist(usr, this.GetAggregateRootId())
}

// 创建并初始化
func (this *memberImpl) SaveRelation(r *member.Relation) error {
	this._relation = r
	this._relation.MemberId = this._value.Id
	return this._rep.SaveRelation(this._relation)
}

var _ member.IFavoriteManager = new(favoriteManagerImpl)

// 收藏服务
type favoriteManagerImpl struct {
	_memberId int
	_rep      member.IMemberRep
}

func newFavoriteManagerImpl(memberId int,
	rep member.IMemberRep) member.IFavoriteManager {
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
func (this *favoriteManagerImpl) Favorite(favType, referId int) error {
	if this.Favored(favType, referId) {
		return member.ErrFavored
	}
	return this._rep.Favorite(this._memberId, favType, referId)
}

// 是否已收藏
func (this *favoriteManagerImpl) Favored(favType, referId int) bool {
	return this._rep.Favored(this._memberId, favType, referId)
}

// 取消收藏
func (this *favoriteManagerImpl) Cancel(favType, referId int) error {
	return this._rep.CancelFavorite(this._memberId, favType, referId)
}

// 收藏商品
func (this *favoriteManagerImpl) FavoriteGoods(goodsId int) error {
	return this.Favorite(member.FavTypeGoods, goodsId)
}

// 取消收藏商品
func (this *favoriteManagerImpl) CancelGoodsFavorite(goodsId int) error {
	return this.Cancel(member.FavTypeGoods, goodsId)
}

// 商品是否已收藏
func (this *favoriteManagerImpl) GoodsFavored(goodsId int) bool {
	return this.Favored(member.FavTypeGoods, goodsId)
}

// 收藏店铺
func (this *favoriteManagerImpl) FavoriteShop(shopId int) error {
	return this.Favorite(member.FavTypeShop, shopId)
}

// 取消收藏店铺
func (this *favoriteManagerImpl) CancelShopFavorite(shopId int) error {
	return this.Cancel(member.FavTypeShop, shopId)
}

// 商店是否已收藏
func (this *favoriteManagerImpl) ShopFavored(shopId int) bool {
	return this.Favored(member.FavTypeShop, shopId)
}
