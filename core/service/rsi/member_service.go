/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-09 20:14
 * description :
 * history :
 */

package rsi

import (
	"errors"
	"fmt"
	"github.com/jsix/gof"
	dm "go2o/core/domain"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/mss/notify"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/dto"
	"go2o/core/infrastructure/domain"
	"go2o/core/infrastructure/format"
	"go2o/core/module"
	"go2o/core/query"
	"go2o/core/service/thrift/idl/gen-go/define"
	"go2o/core/service/thrift/parser"
	"go2o/core/variable"
	"log"
	"strings"
	"time"
)

var _ define.MemberService = new(memberService)

type memberService struct {
	_rep            member.IMemberRepo
	_partnerService *merchantService
	_query          *query.MemberQuery
	_orderQuery     *query.OrderQuery
	valRepo         valueobject.IValueRepo
}

func NewMemberService(mchService *merchantService, rep member.IMemberRepo,
	q *query.MemberQuery, oq *query.OrderQuery, valRepo valueobject.IValueRepo) *memberService {
	ms := &memberService{
		_rep:            rep,
		_query:          q,
		_partnerService: mchService,
		_orderQuery:     oq,
		valRepo:         valRepo,
	}
	return ms
	//return ms.init()
}

func (ms *memberService) init() *memberService {
	db := gof.CurrentApp.Db()
	list := []*member.Relation{}
	db.GetOrm().Select(&list, "")
	for _, v := range list {
		ms._rep.GetMember(v.MemberId).SaveRelation(v)
	}
	return ms
}

// 根据会员编号获取会员
func (ms *memberService) getValueMember(memberId int64) *member.Member {
	if memberId > 0 {
		v := ms._rep.GetMember(memberId)
		if v != nil {
			vv := v.GetValue()
			return &vv
		}
	}
	return nil
}

// 根据会员编号获取会员
func (ms *memberService) GetMember(id int64) (*define.Member, error) {
	v := ms.getValueMember(id)
	if v != nil {
		return parser.MemberDto(v), nil
	}
	return nil, nil
}

// 根据用户名获取会员
func (ms *memberService) GetMemberByUser(usr string) (*define.Member, error) {
	v := ms._rep.GetMemberByUsr(usr)
	if v != nil {
		return parser.MemberDto(v), nil
	}
	return nil, nil
}

// 获取资料
func (ms *memberService) GetProfile(memberId int64) (*define.Profile, error) {
	m := ms._rep.GetMember(memberId)
	if m != nil {
		v := m.Profile().GetProfile()
		return parser.MemberProfile(&v), nil
	}
	return nil, nil
}

// 保存资料
func (ms *memberService) SaveProfile(v *define.Profile) error {
	if v.MemberId > 0 {
		v2 := parser.MemberProfile2(v)
		m := ms._rep.GetMember(v.MemberId)
		if m == nil {
			return member.ErrNoSuchMember
		}
		return m.Profile().SaveProfile(v2)
	}
	return nil
}

// 升级为高级会员
func (ms *memberService) Premium(memberId int64, v int32, expires int64) (*define.Result_, error) {
	m := ms._rep.GetMember(memberId)
	if m == nil {
		return parser.I64Result(memberId, member.ErrNoSuchMember), nil
	}
	err := m.Premium(v, expires)
	return parser.I64Result(memberId, err), nil
}

// 检查会员的会话Token是否正确
func (ms *memberService) CheckToken(memberId int64, token string) (r bool, err error) {
	md := module.Get(module.M_MM).(*module.MemberModule)
	return md.CheckToken(memberId, token), nil
}

// 获取会员的会员Token,reset表示是否重置会员的token
func (ms *memberService) GetToken(memberId int64, reset bool) (r string, err error) {
	pubToken := ""
	md := module.Get(module.M_MM).(*module.MemberModule)
	if !reset {
		pubToken = md.GetToken(memberId)
	}
	if reset || (pubToken == "" && memberId > 0) {
		m := ms.getValueMember(memberId)
		if m != nil {
			return md.ResetToken(memberId, m.Pwd), nil
		}
	}
	return pubToken, nil
}

// 移除会员的Token
func (ms *memberService) RemoveToken(memberId int64) (err error) {
	md := module.Get(module.M_MM).(*module.MemberModule)
	md.RemoveToken(memberId)
	return nil
}

// 更改手机号码，不验证手机格式
func (ms *memberService) ChangePhone(memberId int64, phone string) error {
	m := ms._rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.Profile().ChangePhone(phone)
}

// 是否已收藏
func (ms *memberService) Favored(memberId int64, favType int, referId int32) bool {
	return ms._rep.CreateMemberById(memberId).
		Favorite().Favored(favType, referId)
}

// 取消收藏
func (ms *memberService) Cancel(memberId int64, favType int, referId int32) error {
	return ms._rep.CreateMemberById(memberId).
		Favorite().Cancel(favType, referId)
}

// 收藏商品
func (ms *memberService) FavoriteGoods(memberId int64, goodsId int32) error {
	return ms._rep.CreateMemberById(memberId).
		Favorite().Favorite(member.FavTypeGoods, goodsId)
}

// 取消收藏商品
func (ms *memberService) CancelGoodsFavorite(memberId int64, goodsId int32) error {
	return ms._rep.CreateMemberById(memberId).
		Favorite().Cancel(member.FavTypeGoods, goodsId)
}

// 收藏店铺
func (ms *memberService) FavoriteShop(memberId int64, shopId int32) error {
	return ms._rep.CreateMemberById(memberId).
		Favorite().Favorite(member.FavTypeShop, shopId)
}

// 取消收藏店铺
func (ms *memberService) CancelShopFavorite(memberId int64, shopId int32) error {
	return ms._rep.CreateMemberById(memberId).
		Favorite().Cancel(member.FavTypeShop, shopId)
}

// 商品是否已收藏
func (ms *memberService) GoodsFavored(memberId int64, goodsId int32) bool {
	return ms._rep.CreateMemberById(memberId).
		Favorite().Favored(member.FavTypeGoods, goodsId)
}

// 商店是否已收藏
func (ms *memberService) ShopFavored(memberId int64, shopId int32) bool {
	return ms._rep.CreateMemberById(memberId).
		Favorite().Favored(member.FavTypeShop, shopId)
}

/**================ 会员等级 ==================**/
// 获取所有会员等级
func (ms *memberService) GetMemberLevels() []*member.Level {
	return ms._rep.GetManager().LevelManager().GetLevelSet()
}

// 根据编号获取会员等级信息
func (ms *memberService) GetLevelById(id int32) *member.Level {
	return ms._rep.GetManager().LevelManager().GetLevelById(id)
}

// 根据可编程字符获取会员等级
func (ms *memberService) GetLevelByProgramSign(sign string) *member.Level {
	return ms._rep.GetManager().LevelManager().GetLevelByProgramSign(sign)
}

// 保存会员等级信息
func (ms *memberService) SaveMemberLevel(v *member.Level) (int32, error) {
	return ms._rep.GetManager().LevelManager().SaveLevel(v)
}

// 删除会员等级
func (ms *memberService) DelMemberLevel(levelId int32) error {
	return ms._rep.GetManager().LevelManager().DeleteLevel(levelId)
}

// 获取下一个等级
func (ms *memberService) GetNextLevel(levelId int32) *member.Level {
	return ms._rep.GetManager().LevelManager().GetNextLevelById(levelId)
}

// 获取启用中的最大等级,用于判断是否可以升级
func (ms *memberService) GetHighestLevel() member.Level {
	lv := ms._rep.GetManager().LevelManager().GetHighestLevel()
	if lv != nil {
		return *lv
	}
	return member.Level{}
}

func (ms *memberService) GetPresentLog(memberId int64, logId int32) *member.PresentLog {
	m := ms._rep.GetMember(memberId)
	return m.GetAccount().GetPresentLog(logId)
}

func (ms *memberService) getMember(memberId int64) (
	member.IMember, error) {
	if memberId <= 0 {
		return nil, member.ErrNoSuchMember
	}
	m := ms._rep.GetMember(memberId)
	if m == nil {
		return m, member.ErrNoSuchMember
	}
	return m, nil
}

func (ms *memberService) GetMemberIdByInvitationCode(code string) int64 {
	return ms._rep.GetMemberIdByInvitationCode(code)
}

// 根据信息获取会员编号
func (ms *memberService) GetMemberIdByBasis(str string, basic int) int64 {
	switch basic {
	default:
	case notify.TypePhoneMessage:
		return ms._rep.GetMemberIdByPhone(str)
	case notify.TypeEmailMessage:
		return ms._rep.GetMemberIdByEmail(str)
	}
	return -1
}

// 发送验证码
func (ms *memberService) SendCode(memberId int64, operation string, msgType int) (string, error) {
	m := ms._rep.GetMember(memberId)
	if m == nil {
		return "", member.ErrNoSuchMember
	}
	return m.SendCheckCode(operation, msgType)
}

// 对比验证码
func (ms *memberService) CompareCode(memberId int64, code string) error {
	m := ms._rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.CompareCode(code)
}

// 更改会员用户名
func (ms *memberService) ChangeUsr(memberId int64, usr string) error {
	m := ms._rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.ChangeUsr(usr)
}

// 更改会员等级
func (ms *memberService) ChangeLevel(memberId int64, levelId int32,
	review bool, paymentId int32) error {
	m := ms._rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.ChangeLevel(levelId, paymentId, review)
}

// 上传会员头像
func (ms *memberService) SetAvatar(memberId int64, avatar string) error {
	m := ms._rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.Profile().SetAvatar(avatar)
}

// 保存用户
func (ms *memberService) SaveMember(v *define.Member) (int64, error) {
	if v.ID > 0 {
		return ms.updateMember(v)
	}
	return -1, errors.New("Create member use \"RegisterMember\" method.")
}

func (ms *memberService) updateMember(v *define.Member) (int64, error) {
	m := ms._rep.GetMember(v.ID)
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
func (ms *memberService) RegisterMember(mchId int32, v1 *define.Member,
	pro1 *define.Profile, cardId string, invitationCode string) (int64, error) {
	if v1 == nil || pro1 == nil {
		return 0, errors.New("missing data")
	}
	v := parser.Member(v1)
	pro := parser.MemberProfile2(pro1)
	invitationId, err := ms._rep.GetManager().PrepareRegister(
		v, pro, invitationCode)
	if err == nil {
		m := ms._rep.CreateMember(v) //创建会员
		id, err := m.Save()
		if err == nil {
			pro.Sex = 1
			pro.MemberId = id
			//todo: 如果注册失败，则删除。应使用SQL-TRANSFER
			if err = m.Profile().SaveProfile(pro); err != nil {
				ms._rep.DeleteMember(id)
			} else {
				// 保存关联信息
				rl := m.GetRelation()
				rl.InviterId = invitationId
				rl.RegMchId = mchId
				rl.CardCard = cardId
				err = m.SaveRelation(rl)
			}
		}
		return id, err
	}
	return -1, err
}

// 获取会员等级
func (ms *memberService) GetMemberLevel(memberId int64) *member.Level {
	m := ms._rep.GetMember(memberId)
	if m == nil {
		return nil
	}
	return m.GetLevel()
}

func (ms *memberService) GetRelation(memberId int64) *member.Relation {
	return ms._rep.GetRelation(memberId)
}

// 锁定/解锁会员
func (ms *memberService) LockMember(memberId int64) (bool, error) {
	m := ms._rep.GetMember(memberId)
	if m == nil {
		return false, member.ErrNoSuchMember
	}

	state := m.GetValue().State
	if state == 1 {
		return false, m.Lock()
	}
	return true, m.Unlock()
}

// 判断资料是否完善
func (ms *memberService) ProfileCompleted(memberId int64) bool {
	m := ms._rep.GetMember(memberId)
	if m != nil {
		return m.Profile().ProfileCompleted()
	}
	return false
}

// 重置密码
func (ms *memberService) ResetPassword(memberId int64) string {
	m := ms._rep.GetMember(memberId)
	if m != nil {
		newPwd := domain.GenerateRandomIntPwd(6)
		newEncPwd := domain.MemberSha1Pwd(newPwd)
		if err := m.Profile().ModifyPassword(newEncPwd, ""); err == nil {
			return newPwd
		} else {
			log.Println("--- 重置密码:", err)
		}
	}
	return ""
}

// 重置交易密码
func (ms *memberService) ResetTradePwd(memberId int64) string {
	m := ms._rep.GetMember(memberId)
	if m != nil {
		newPwd := domain.GenerateRandomIntPwd(6)
		newEncPwd := domain.TradePwd(newPwd)
		if err := m.Profile().ModifyTradePassword(newEncPwd, ""); err == nil {
			return newPwd
		} else {
			log.Println("--- 重置交易密码:", err)
		}
	}
	return ""
}

// 修改密码
func (ms *memberService) ModifyPassword(memberId int64, newPwd, oldPwd string) error {
	m := ms._rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	//return m.Profile().ModifyPassword(newPwd, oldPwd)

	// 兼容旧加密的密码
	pro := m.Profile()
	err := pro.ModifyPassword(newPwd, oldPwd)
	if err == dm.ErrPwdOldPwdNotRight {
		err = pro.ModifyPassword(newPwd, domain.Sha1(oldPwd))
	}
	return err
}

//修改密码,传入密文密码
func (ms *memberService) ModifyTradePassword(memberId int64,
	oldPwd, newPwd string) error {
	m := ms._rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.Profile().ModifyTradePassword(newPwd, oldPwd)
}

// 登录，返回结果(Result)和会员编号(Id);
// Result值为：-1:会员不存在; -2:账号密码不正确; -3:账号被停用
func (ms *memberService) Login(usr string, pwd string, update bool) (r *define.Result64, err error) {

	usr = strings.ToLower(strings.TrimSpace(usr))
	val := ms._rep.GetMemberByUsr(usr)

	if val == nil {
		val = ms._rep.GetMemberValueByPhone(usr)
	}
	r = &define.Result64{}
	if val == nil {
		r.Message = member.ErrNoSuchMember.Error()
		return r, nil
	}
	if val.Pwd != pwd {
		//todo: 兼容旧密码
		if val.Pwd != domain.Sha1(pwd) {
			r.Message = member.ErrCredential.Error()
			return r, nil
		}
	}
	if val.State == member.StateStopped {
		r.Message = member.ErrMemberDisabled.Error()
		return r, nil
	}
	r.ID = val.Id
	r.Result_ = true
	if update {
		m := ms._rep.GetMember(val.Id)
		m.UpdateLoginTime()
	}
	return r, nil
}

// 检查与现有用户不同的用户是否存在,如存在则返回错误
func (ms *memberService) CheckUsr(usr string, memberId int64) error {
	if len(usr) < 6 {
		return member.ErrUsrLength
	}
	if ms._rep.CheckUsrExist(usr, memberId) {
		return member.ErrUsrExist
	}
	return nil
}

// 检查手机号码是否与会员一致
func (ms *memberService) CheckPhone(phone string, memberId int64) error {
	return ms._rep.GetManager().CheckPhoneBind(phone, memberId)
}

// 获取会员账户
func (ms *memberService) GetAccount(memberId int64) (*define.Account, error) {
	m := ms._rep.CreateMember(&member.Member{Id: memberId})
	acc := m.GetAccount()
	if acc != nil {
		return parser.AccountDto(acc.GetValue()), nil
	}
	return nil, nil
}

// 获取邀请人会员编号数组
func (ms *memberService) InviterArray(memberId int64, depth int32) (r []int64, err error) {
	m := ms._rep.CreateMember(&member.Member{Id: memberId})
	if m != nil {
		return m.Invitation().InviterArray(memberId, depth), nil
	}
	return []int64{}, nil
}

func (ms *memberService) GetBank(memberId int64) *member.BankInfo {
	m := ms._rep.CreateMember(&member.Member{Id: memberId})
	b := m.Profile().GetBank()
	return &b
}

func (ms *memberService) SaveBankInfo(v *member.BankInfo) error {
	m := ms._rep.CreateMember(&member.Member{Id: v.MemberId})
	return m.Profile().SaveBank(v)
}

// 解锁银行卡信息
func (ms *memberService) UnlockBankInfo(memberId int64) error {
	m := ms._rep.CreateMember(&member.Member{Id: memberId})
	return m.Profile().UnlockBank()
}

// 实名认证信息
func (ms *memberService) GetTrustedInfo(memberId int64) member.TrustedInfo {
	m := ms._rep.GetMember(memberId)
	if m == nil {
		return member.TrustedInfo{}
	}
	return m.Profile().GetTrustedInfo()
}

// 保存实名认证信息
func (ms *memberService) SaveTrustedInfo(memberId int64, v *member.TrustedInfo) error {
	m := ms._rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.Profile().SaveTrustedInfo(v)
}

// 审核实名认证,若重复审核将返回错误
func (ms *memberService) ReviewTrustedInfo(memberId int64, pass bool, remark string) error {
	m := ms._rep.GetMember(memberId)
	return m.Profile().ReviewTrustedInfo(pass, remark)
}

// 获取返现记录
func (ms *memberService) QueryIncomeLog(memberId int64, begin, end int,
	where, orderBy string) (int, []map[string]interface{}) {
	return ms._query.QueryBalanceLog(memberId, begin, end, where, orderBy)
}

// 获取分页商铺收藏
func (ms *memberService) PagedShopFav(memberId int64, begin, end int,
	where string) (int, []*dto.PagedShopFav) {
	return ms._query.PagedShopFav(memberId, begin, end, where)
}

// 获取分页商铺收藏
func (ms *memberService) PagedGoodsFav(memberId int64, begin, end int,
	where string) (int, []*dto.PagedGoodsFav) {
	return ms._query.PagedGoodsFav(memberId, begin, end, where)
}

// 获取余额账户分页记录
func (ms *memberService) PagedBalanceAccountLog(memberId int64, begin, end int,
	where, orderBy string) (int, []map[string]interface{}) {
	return ms._query.PagedBalanceAccountLog(memberId, begin, end, where, orderBy)
}

// 获取钱包账户分页记录
func (ms *memberService) PagedWalletAccountLog(memberId int64, begin, end int,
	where, orderBy string) (int, []map[string]interface{}) {
	return ms._query.PagedWalletAccountLog(memberId, begin, end, where, orderBy)
}

// 查询分页订单
func (ms *memberService) QueryPagerOrder(memberId int64, begin, size int, pagination bool,
	where, orderBy string) (num int, rows []*dto.PagedMemberSubOrder) {
	return ms._orderQuery.QueryPagerOrder(memberId, begin, size, pagination, where, orderBy)
}

// 查询分页订单
func (ms *memberService) QueryPagerTradeOrder(memberId int64, begin, size int,
	pagination bool, where, orderBy string) (num int, rows []*define.ComplexOrder) {
	return ms._orderQuery.QueryPagerTradeOrder(memberId, begin, size, pagination, where, orderBy)
}

/*********** 收货地址 ***********/
func (ms *memberService) GetAddressList(memberId int64) []*member.Address {
	return ms._rep.GetDeliverAddress(memberId)
}

//获取配送地址
func (ms *memberService) GetAddress(memberId int64, addrId int64) (
	*define.Address, error) {
	m := ms._rep.CreateMember(&member.Member{Id: memberId})
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
		d.Area = ms.valRepo.GetAreaString(
			v.Province, v.City, v.District)
		return d, nil
	}
	return nil, nil
}

//保存配送地址
func (ms *memberService) SaveAddress(memberId int64, e *member.Address) (int64, error) {
	m := ms._rep.CreateMember(&member.Member{Id: memberId})
	var v member.IDeliverAddress
	var err error
	if e.Id > 0 {
		v = m.Profile().GetAddress(e.Id)
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
func (ms *memberService) DeleteAddress(memberId int64, deliverId int64) error {
	m := ms._rep.CreateMember(&member.Member{Id: memberId})
	return m.Profile().DeleteAddress(deliverId)
}

//设置余额优先支付
func (ms *memberService) BalancePriorityPay(memberId int64, enabled bool) error {
	m := ms._rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().SetPriorityPay(member.AccountBalance, enabled)
}

//判断会员是否由指定会员邀请推荐的
func (ms *memberService) IsInvitation(memberId int64, invitationMemberId int64) bool {
	m := ms._rep.CreateMember(&member.Member{Id: memberId})
	return m.Invitation().InvitationBy(invitationMemberId)
}

// 获取我邀请的会员及会员邀请的人数
func (ms *memberService) GetMyPagedInvitationMembers(memberId int64,
	begin, end int) (total int, rows []*dto.InvitationMember) {
	iv := ms._rep.CreateMember(&member.Member{Id: memberId}).Invitation()
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
func (m *memberService) GetReferNum(memberId int64, layer int) int {
	return m._query.GetReferNum(memberId, layer)
}

// 获取会员最后更新时间
func (ms *memberService) GetMemberLatestUpdateTime(memberId int64) int64 {
	return ms._rep.GetMemberLatestUpdateTime(memberId)
}

func (ms *memberService) GetMemberList(ids []int64) []*dto.MemberSummary {
	list := ms._query.GetMemberList(ids)
	for _, v := range list {
		v.Avatar = format.GetResUrl(v.Avatar)
	}
	return list
}

// 获取会员汇总信息
func (ms *memberService) Complex(memberId int64) (*define.ComplexMember, error) {
	m := ms._rep.GetMember(memberId)
	if m != nil {
		s := m.Complex()
		return parser.ComplexMemberDto(s), nil
	}
	return nil, nil
}

// 获取余额变动信息
func (ms *memberService) GetBalanceInfoById(memberId int64, infoId int32) *member.BalanceInfo {
	m := ms._rep.GetMember(memberId)
	if m == nil {
		return nil
	}
	return m.GetAccount().GetBalanceInfo(infoId)
}

// 增加积分
func (ms *memberService) AddIntegral(memberId int64, iType int,
	orderNo string, value int64, remark string) error {
	m := ms._rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().AddIntegral(iType, orderNo, value, remark)
}

// 充值,account为账户类型,kind为业务类型
func (ms *memberService) ChargeAccount(memberId int64, account int32,
	kind int32, title, outerNo string, amount float64, relateUser int64) (*define.Result_, error) {
	var err error
	m := ms._rep.CreateMember(&member.Member{Id: memberId})
	acc := m.GetAccount()
	if acc == nil {
		err = member.ErrNoSuchMember
	} else {
		err = acc.Charge(account, kind, title, outerNo, float32(amount), relateUser)
	}
	return parser.Result(0, err), nil
}

// 冻结积分,当new为true不扣除积分,反之扣除积分
func (ms *memberService) FreezesIntegral(memberId int64, value int64,
	new bool, remark string) error {
	m := ms._rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().FreezesIntegral(value, new, remark)
}

// 解冻积分
func (ms *memberService) UnfreezesIntegral(memberId int64,
	value int64, remark string) error {
	m := ms._rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().UnfreezesIntegral(value, remark)
}

// 抵扣账户
func (ms *memberService) DiscountAccount(memberId int64, account int32, title string,
	outerNo string, amount float64, relateUser int64, mustLargeZero bool) (r *define.Result_, err error) {
	m, err := ms.getMember(memberId)
	if err == nil {
		acc := m.GetAccount()
		switch int(account) {
		case member.AccountWallet:
			err = acc.DiscountWallet(title, outerNo, float32(amount),
				member.DefaultRelateUser, mustLargeZero)
		}
	}
	return parser.I64Result(memberId, err), nil
}

// 扣减奖金
func (ms *memberService) DiscountWallet(memberId int64, title string,
	tradeNo string, amount float32, mustLargeZero bool) error {
	m, err := ms.getMember(memberId)
	if err != nil {
		return err
	}
	return m.GetAccount().DiscountWallet(title, tradeNo, amount,
		member.DefaultRelateUser, mustLargeZero)
}

// 流通账户
func (ms *memberService) ChargeFlowBalance(memberId int64, title string,
	tradeNo string, amount float32) error {
	m, err := ms.getMember(memberId)
	if err != nil {
		return err
	}
	return m.GetAccount().ChargeFlowBalance(title, tradeNo, amount)
}

// 验证交易密码
func (ms *memberService) VerifyTradePwd(memberId int64, tradePwd string) (bool, error) {
	im, err := ms.getMember(memberId)
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
func (ms *memberService) SubmitTakeOutRequest(memberId int64, takeKind int32,
	applyAmount float32, commission float32) (int32, string, error) {
	m, err := ms.getMember(memberId)
	if err != nil {
		return 0, "", err
	}

	acc := m.GetAccount()
	var title string
	switch takeKind {
	case member.KindWalletTakeOutToBankCard:
		title = "提现到银行卡"
	case member.KindWalletTakeOutToBalance:
		title = "充值账户"
	case member.KindWalletTakeOutToThirdPart:
		title = "充值到第三方账户"
	}
	return acc.RequestTakeOut(takeKind, title, applyAmount, commission)
}

// 获取最近的提现
func (ms *memberService) GetLatestTakeOut(memberId int64) *member.BalanceInfo {
	return ms._query.GetLatestBalanceInfoByKind(memberId,
		member.KindWalletTakeOutToBankCard)
}

// 获取最近的提现描述
func (ms *memberService) GetLatestApplyCashText(memberId int64) string {
	var latestInfo string
	latestApplyInfo := ms.GetLatestTakeOut(memberId)
	if latestApplyInfo != nil {
		var sText string
		switch latestApplyInfo.State {
		case 0:
			sText = "已申请"
		case 1:
			sText = "已审核,等待打款"
		case 2:
			sText = "被退回"
		case 3:
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
func (a *memberService) ConfirmTakeOutRequest(memberId int64,
	infoId int32, pass bool, remark string) error {
	m, err := a.getMember(memberId)
	if err == nil {
		err = m.GetAccount().ConfirmTakeOut(infoId, pass, remark)
	}
	return err
}

// 完成提现
func (ms *memberService) FinishTakeOutRequest(memberId int64, id int32, tradeNo string) error {
	m, err := ms.getMember(memberId)
	if err != nil {
		return err
	}
	return m.GetAccount().FinishTakeOut(id, tradeNo)
}

// 冻结余额
func (ms *memberService) Freeze(memberId int64, title string,
	tradeNo string, amount float32, referId int64) error {
	m := ms._rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().Freeze(title, tradeNo, amount, referId)
}

// 解冻金额
func (ms *memberService) Unfreeze(memberId int64, title string,
	tradeNo string, amount float32, referId int64) error {
	m := ms._rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().Unfreeze(title, tradeNo, amount, referId)
}

// 冻结赠送金额
func (ms *memberService) FreezeWallet(memberId int64, title string,
	tradeNo string, amount float32, referId int64) error {
	m := ms._rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().FreezeWallet(title, tradeNo, amount, referId)
}

// 解冻赠送金额
func (ms *memberService) UnfreezeWallet(memberId int64, title string,
	tradeNo string, amount float32, referId int64) error {
	m := ms._rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().UnfreezeWallet(title, tradeNo, amount, referId)
}

// 将冻结金额标记为失效
func (ms *memberService) FreezeExpired(memberId int64, accountKind int, amount float32,
	remark string) error {
	m := ms._rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().FreezeExpired(accountKind, amount, remark)
}

// 转账余额到其他账户
func (ms *memberService) TransferAccount(accountKind int, fromMember int64,
	toMember int64, amount float32, csnRate float32, remark string) error {
	m := ms._rep.GetMember(fromMember)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().TransferAccount(accountKind, toMember,
		amount, csnRate, remark)
}

// 转账余额到其他账户
func (ms *memberService) TransferBalance(memberId int64, kind int32, amount float32, tradeNo string,
	toTitle, fromTitle string) error {
	m := ms._rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().TransferBalance(kind, amount, tradeNo, toTitle, fromTitle)
}

// 转账返利账户,kind为转账类型，如 KindBalanceTransfer等
// commission手续费
func (ms *memberService) TransferWallet(memberId int64, kind int32, amount float32, commission float32,
	tradeNo string, toTitle string, fromTitle string) error {
	m := ms._rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().TransferWallet(kind, amount, commission,
		tradeNo, toTitle, fromTitle)
}

// 转账活动账户,kind为转账类型，如 KindBalanceTransfer等
// commission手续费
func (ms *memberService) TransferFlow(memberId int64, kind int32, amount float32,
	commission float32, tradeNo string, toTitle string, fromTitle string) error {
	m := ms._rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().TransferFlow(kind, amount, commission, tradeNo,
		toTitle, fromTitle)
}

// 将活动金转给其他人
func (ms *memberService) TransferFlowTo(memberId int64, toMemberId int64, kind int32,
	amount float32, commission float32, tradeNo string, toTitle string,
	fromTitle string) error {
	m := ms._rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().TransferFlowTo(toMemberId, kind, amount,
		commission, tradeNo, toTitle, fromTitle)
}

// 根据用户或手机筛选会员
func (ms *memberService) FilterMemberByUsrOrPhone(key string) []*dto.SimpleMember {
	return ms._query.FilterMemberByUsrOrPhone(key)
}

// 根据用户名货手机获取会员
func (ms *memberService) GetMemberByUserOrPhone(key string) *dto.SimpleMember {
	return ms._query.GetMemberByUsrOrPhone(key)
}

// 根据手机获取会员编号
func (ms *memberService) GetMemberIdByPhone(phone string) int64 {
	return ms._query.GetMemberIdByPhone(phone)
}

// 会员推广排名
func (ms *memberService) GetMemberInviRank(mchId int32, allTeam bool,
	levelComp string, level int, startTime int64, endTime int64,
	num int) []*dto.RankMember {
	return ms._query.GetMemberInviRank(mchId, allTeam, levelComp, level, startTime, endTime, num)
}

// 生成会员账户人工单据
func (ms *memberService) NewBalanceTicket(mchId int32, memberId int64, accountType int,
	tit string, amount float32, relateUser int64) (string, error) {
	//todo: 暂时不记录人员,等支持系统多用户后再传入
	if relateUser <= 0 {
		relateUser = 1
	}
	var err error
	var outerNo string
	if amount == 0 {
		return "", member.ErrIncorrectAmount
	}
	m := ms._rep.GetMember(memberId)
	if m == nil {
		return "", member.ErrNoSuchMember
	}
	acc := m.GetAccount()
	var tit2 string
	if accountType == member.AccountWallet {
		outerNo = domain.NewTradeNo(int(mchId))
		if amount > 0 {
			//增加奖金
			tit2 = "[KF]客服调整-" + variable.AliasWalletAccount
			if len(tit) > 0 {
				tit2 = tit2 + "(" + tit + ")"
			}
			err = acc.Charge(member.AccountWallet,
				member.KindWalletServiceAdd,
				tit2, outerNo, amount, relateUser)
		} else {
			//扣减奖金
			tit2 = "[KF]客服扣减-" + variable.AliasWalletAccount
			if len(tit) > 0 {
				tit2 = tit2 + "(" + tit + ")"
			}
			err = acc.DiscountWallet(tit2, outerNo, -amount, relateUser, false)
		}
		return outerNo, err
	}

	if accountType == member.AccountBalance {
		outerNo = domain.NewTradeNo(int(mchId))
		if amount > 0 {
			tit2 = "[KF]客服充值"
			if len(tit) > 0 {
				tit2 = tit2 + "(" + tit + ")"
			}
			err = acc.Charge(member.AccountBalance,
				member.KindBalanceServiceCharge,
				tit2, outerNo, amount, relateUser)
		} else {
			tit2 = "[KF]客服扣减"
			if len(tit) > 0 {
				tit2 = tit2 + "(" + tit + ")"
			}
			err = acc.DiscountBalance(tit2, outerNo, -amount, relateUser)
		}
		return outerNo, err
	}

	return outerNo, err
}

//********* 促销  **********//

// 可用的优惠券分页数据
func (ms *memberService) PagedAvailableCoupon(memberId int64, start, end int) (int, []*dto.SimpleCoupon) {
	return ms._rep.CreateMemberById(memberId).
		GiftCard().PagedAvailableCoupon(start, end)
}

// 已使用的优惠券
func (ms *memberService) PagedAllCoupon(memberId int64, start, end int) (int, []*dto.SimpleCoupon) {
	return ms._rep.CreateMemberById(memberId).
		GiftCard().PagedAllCoupon(start, end)
}

// 过期的优惠券
func (ms *memberService) PagedExpiresCoupon(memberId int64, start, end int) (int, []*dto.SimpleCoupon) {
	return ms._rep.CreateMemberById(memberId).
		GiftCard().PagedExpiresCoupon(start, end)
}
