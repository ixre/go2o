/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-09 20:14
 * description :
 * history :
 */

package dps

import (
	"errors"
	"fmt"
	"github.com/jsix/gof"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/mss/notify"
	"go2o/core/dto"
	"go2o/core/infrastructure/domain"
	"go2o/core/infrastructure/format"
	"go2o/core/query"
	"go2o/core/variable"
	"strings"
	"time"
)

type memberService struct {
	_rep            member.IMemberRep
	_partnerService *merchantService
	_query          *query.MemberQuery
	_orderQuery     *query.OrderQuery
}

func NewMemberService(mchService *merchantService, rep member.IMemberRep,
	q *query.MemberQuery, oq *query.OrderQuery) *memberService {
	ms := &memberService{
		_rep:            rep,
		_query:          q,
		_partnerService: mchService,
		_orderQuery:     oq,
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

// 获取资料
func (ms *memberService) GetProfile(memberId int) *member.Profile {
	m := ms._rep.GetMember(memberId)
	if m != nil {
		v := m.Profile().GetProfile()
		return &v
	}
	return nil
}

// 保存资料
func (ms *memberService) SaveProfile(memberId int, v *member.Profile) error {
	m := ms._rep.GetMember(memberId)
	return m.Profile().SaveProfile(v)
}

// 是否已收藏
func (ms *memberService) Favored(memberId, favType, referId int) bool {
	return ms._rep.CreateMemberById(memberId).
		Favorite().Favored(favType, referId)
}

// 取消收藏
func (ms *memberService) Cancel(memberId, favType, referId int) error {
	return ms._rep.CreateMemberById(memberId).
		Favorite().Cancel(favType, referId)
}

// 收藏商品
func (ms *memberService) FavoriteGoods(memberId, goodsId int) error {
	return ms._rep.CreateMemberById(memberId).
		Favorite().FavoriteGoods(goodsId)
}

// 取消收藏商品
func (ms *memberService) CancelGoodsFavorite(memberId, goodsId int) error {
	return ms._rep.CreateMemberById(memberId).
		Favorite().CancelGoodsFavorite(goodsId)
}

// 收藏店铺
func (ms *memberService) FavoriteShop(memberId, shopId int) error {
	return ms._rep.CreateMemberById(memberId).
		Favorite().FavoriteShop(shopId)
}

// 取消收藏店铺
func (ms *memberService) CancelShopFavorite(memberId, shopId int) error {
	return ms._rep.CreateMemberById(memberId).
		Favorite().CancelShopFavorite(shopId)
}

// 商品是否已收藏
func (ms *memberService) GoodsFavored(memberId, goodsId int) bool {
	return ms._rep.CreateMemberById(memberId).
		Favorite().GoodsFavored(goodsId)
}

// 商店是否已收藏
func (ms *memberService) ShopFavored(memberId, shopId int) bool {
	return ms._rep.CreateMemberById(memberId).
		Favorite().ShopFavored(shopId)
}

/**================ 会员等级 ==================**/
// 获取所有会员等级
func (ms *memberService) GetMemberLevels() []*member.Level {
	return ms._rep.GetManager().LevelManager().GetLevelSet()
}

// 根据编号获取会员等级信息
func (ms *memberService) GetLevelById(id int) *member.Level {
	return ms._rep.GetManager().LevelManager().GetLevelById(id)
}

// 根据可编程字符获取会员等级
func (ms *memberService) GetLevelByProgramSign(sign string) *member.Level {
	return ms._rep.GetManager().LevelManager().GetLevelByProgramSign(sign)
}

// 保存会员等级信息
func (ms *memberService) SaveMemberLevel(v *member.Level) (int, error) {
	return ms._rep.GetManager().LevelManager().SaveLevel(v)
}

// 删除会员等级
func (ms *memberService) DelMemberLevel(levelId int) error {
	return ms._rep.GetManager().LevelManager().DeleteLevel(levelId)
}

// 获取下一个等级
func (ms *memberService) GetNextLevel(levelId int) *member.Level {
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

// 根据会员编号获取会员
func (ms *memberService) GetMember(id int) *member.Member {
	if id <= 0 {
		return nil
	}
	v := ms._rep.GetMember(id)
	if v != nil {
		nv := v.GetValue()
		return &nv
	}
	return nil
}

// 根据用户名获取会员
func (ms *memberService) GetMemberByUser(usr string) *member.Member {
	return ms._rep.GetMemberByUsr(usr)
}

func (ms *memberService) getMember(memberId int) (
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

func (ms *memberService) GetMemberIdByInvitationCode(code string) int {
	return ms._rep.GetMemberIdByInvitationCode(code)
}

// 根据信息获取会员编号
func (ms *memberService) GetMemberIdByBasis(str string, basic int) int {
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
func (ms *memberService) SendCode(memberId int, operation string, msgType int) (string, error) {
	m := ms._rep.GetMember(memberId)
	if m == nil {
		return "", member.ErrNoSuchMember
	}
	return m.SendCheckCode(operation, msgType)
}

// 对比验证码
func (ms *memberService) CompareCode(memberId int, code string) error {
	m := ms._rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.CompareCode(code)
}

// 更改会员用户名
func (ms *memberService) ChangeUsr(id int, usr string) error {
	m := ms._rep.GetMember(id)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.ChangeUsr(usr)
}

// 保存用户
func (ms *memberService) SaveMember(v *member.Member) (int, error) {
	if v.Id > 0 {
		return ms.updateMember(v)
	}
	return -1, errors.New("Create member use \"RegisterMember\" method.")
}

func (ms *memberService) SaveMemberProfile(id int, v *member.Member) (int, error) {
	//todo:
	return -1, nil
}

func (ms *memberService) updateMember(v *member.Member) (int, error) {
	m := ms._rep.GetMember(v.Id)
	if m == nil {
		return -1, member.ErrNoSuchMember
	}
	if err := m.SetValue(v); err != nil {
		return m.GetAggregateRootId(), err
	}
	return m.Save()
}

// 注册会员
func (ms *memberService) RegisterMember(merchantId int, v *member.Member,
	pro *member.Profile, cardId string, invitationCode string) (int, error) {
	invitationId, err := ms._rep.GetManager().PrepareRegister(v, pro, invitationCode)
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
				rl.RefereesId = invitationId
				rl.RegisterMerchantId = merchantId
				rl.CardId = cardId
				err = m.SaveRelation(rl)
			}
		}
		return id, err
	}
	return -1, err
}

func (ms *memberService) GetMemberLevel(memberId int) *member.Level {
	if m := ms._rep.GetMember(memberId); m != nil {
		return m.GetLevel()
	}
	return nil
}

func (ms *memberService) GetRelation(memberId int) *member.Relation {
	return ms._rep.GetRelation(memberId)
}

// 锁定/解锁会员
func (ms *memberService) LockMember(id int) (bool, error) {
	m := ms._rep.GetMember(id)
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
func (ms *memberService) ProfileCompleted(memberId int) bool {
	m := ms._rep.GetMember(memberId)
	if m != nil {
		return m.Profile().ProfileCompleted()
	}
	return false
}

// 重置密码
func (ms *memberService) ResetPassword(memberId int) string {
	m := ms._rep.GetMember(memberId)
	if m != nil {
		newPwd := domain.GenerateRandomIntPwd(6)
		newEncPwd := domain.MemberSha1Pwd(newPwd)
		if m.Profile().ModifyPassword(newEncPwd, "") == nil {
			return newPwd
		}
	}
	return ""
}

// 检查凭据, update:是否更新登录时间
func (ms *memberService) TryLogin(usr, pwd string, update bool) (
	*member.Member, error) {
	usr = strings.ToLower(strings.TrimSpace(usr))

	val := ms._rep.GetMemberByUsr(usr)
	if val == nil {
		val = ms._rep.GetMemberValueByPhone(usr)
	}

	if val == nil {
		return nil, errors.New("会员不存在")
	}

	if val.Pwd != pwd {
		return nil, errors.New("会员用户或密码不正确")
	}

	if val.State == member.StateStopped {
		return nil, errors.New("会员已停用")
	}

	m := ms._rep.GetMember(val.Id)

	if update {
		unix := time.Now().Unix()
		val.LastLoginTime = unix
		val.UpdateTime = unix
		m.Save()
	}

	return val, nil
}

// 检查与现有用户不同的用户是否存在,如存在则返回错误
func (ms *memberService) CheckUsr(usr string, memberId int) error {
	if len(usr) < 6 {
		return member.ErrUsrLength
	}
	if ms._rep.CheckUsrExist(usr, memberId) {
		return member.ErrUsrExist
	}
	return nil
}

// 检查手机号码是否与会员一致
func (ms *memberService) CheckPhone(phone string, memberId int) error {
	return ms._rep.GetManager().CheckPhoneBind(phone, memberId)
}

func (ms *memberService) GetAccount(memberId int) *member.Account {
	m := ms._rep.CreateMember(&member.Member{Id: memberId})
	//m, _ := ms._memberRep.GetMember(memberId)
	//m.AddExp(300)
	return m.GetAccount().GetValue()
}

func (ms *memberService) GetBank(memberId int) *member.BankInfo {
	m := ms._rep.CreateMember(&member.Member{Id: memberId})
	b := m.Profile().GetBank()
	return &b
}

func (ms *memberService) SaveBankInfo(v *member.BankInfo) error {
	m := ms._rep.CreateMember(&member.Member{Id: v.MemberId})
	return m.Profile().SaveBank(v)
}

// 解锁银行卡信息
func (ms *memberService) UnlockBankInfo(memberId int) error {
	m := ms._rep.CreateMember(&member.Member{Id: memberId})
	return m.Profile().UnlockBank()
}

// 实名认证信息
func (ms *memberService) GetTrustedInfo(memberId int) member.TrustedInfo {
	m := ms._rep.GetMember(memberId)
	return m.Profile().GetTrustedInfo()
}

// 保存实名认证信息
func (ms *memberService) SaveTrustedInfo(memberId int, v *member.TrustedInfo) error {
	m := ms._rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.Profile().SaveTrustedInfo(v)
}

// 审核实名认证,若重复审核将返回错误
func (ms *memberService) ReviewTrustedInfo(memberId int, pass bool, remark string) error {
	m := ms._rep.GetMember(memberId)
	return m.Profile().ReviewTrustedInfo(pass, remark)
}

// 获取返现记录
func (ms *memberService) QueryIncomeLog(memberId, begin, end int,
	where, orderBy string) (int, []map[string]interface{}) {
	return ms._query.QueryBalanceLog(memberId, begin, end, where, orderBy)
}

// 获取分页商铺收藏
func (ms *memberService) PagedShopFav(memberId int, begin, end int,
	where string) (int, []*dto.PagedShopFav) {
	return ms._query.PagedShopFav(memberId, begin, end, where)
}

// 获取分页商铺收藏
func (ms *memberService) PagedGoodsFav(memberId int, begin, end int,
	where string) (int, []*dto.PagedGoodsFav) {
	return ms._query.PagedGoodsFav(memberId, begin, end, where)
}

// 获取余额账户分页记录
func (ms *memberService) PagedBalanceAccountLog(memberId, begin, end int,
	where, orderBy string) (int, []map[string]interface{}) {
	return ms._query.PagedBalanceAccountLog(memberId, begin, end, where, orderBy)
}

// 获取赠送账户分页记录
func (ms *memberService) PagedPresentAccountLog(memberId, begin, end int,
	where, orderBy string) (int, []map[string]interface{}) {
	return ms._query.PagedPresentAccountLog(memberId, begin, end, where, orderBy)
}

// 查询分页订单
func (ms *memberService) QueryPagerOrder(memberId, begin, size int, pagination bool,
	where, orderBy string) (num int, rows []*dto.PagedMemberSubOrder) {
	return ms._orderQuery.QueryPagerOrder(memberId, begin, size, pagination, where, orderBy)
}

/*********** 收货地址 ***********/
func (ms *memberService) GetDeliverAddress(memberId int) []*member.DeliverAddress {
	return ms._rep.GetDeliverAddress(memberId)
}

//获取配送地址
func (ms *memberService) GetDeliverAddressById(memberId,
	deliverId int) *member.DeliverAddress {
	m := ms._rep.CreateMember(&member.Member{Id: memberId})
	v := m.Profile().GetDeliver(deliverId).GetValue()
	return &v
}

//保存配送地址
func (ms *memberService) SaveAddress(memberId int, e *member.DeliverAddress) (int, error) {
	m := ms._rep.CreateMember(&member.Member{Id: memberId})
	var v member.IDeliverAddress
	var err error
	if e.Id > 0 {
		v = m.Profile().GetDeliver(e.Id)
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
func (ms *memberService) DeleteDeliverAddress(memberId int, deliverId int) error {
	m := ms._rep.CreateMember(&member.Member{Id: memberId})
	return m.Profile().DeleteDeliver(deliverId)
}

// 修改密码
func (ms *memberService) ModifyPassword(memberId int, oldPwd, newPwd string) error {
	m := ms._rep.GetMember(memberId)
	if m != nil {
		newEncPwd := domain.MemberSha1Pwd(newPwd)
		oldEncPwd := oldPwd
		if oldEncPwd != "" {
			oldEncPwd = domain.MemberSha1Pwd(oldPwd)
		}
		return m.Profile().ModifyPassword(newEncPwd, oldEncPwd)
	}
	return member.ErrNoSuchMember
}

//修改密码
func (ms *memberService) ModifyTradePassword(memberId int, oldPwd, newPwd string) error {
	m := ms._rep.GetMember(memberId)
	if m != nil {
		return m.Profile().ModifyTradePassword(newPwd, oldPwd)
	}
	return member.ErrNoSuchMember
}

//设置余额优先支付
func (ms *memberService) BalancePriorityPay(memberId int, enabled bool) error {
	m := ms._rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().SetPriorityPay(member.AccountBalance, enabled)
}

//判断会员是否由指定会员邀请推荐的
func (ms *memberService) IsInvitation(memberId int, invitationMemberId int) bool {
	m := ms._rep.CreateMember(&member.Member{Id: memberId})
	return m.Invitation().InvitationBy(invitationMemberId)
}

// 获取我邀请的会员及会员邀请的人数
func (ms *memberService) GetMyPagedInvitationMembers(memberId int,
	begin, end int) (total int, rows []*dto.InvitationMember) {
	iv := ms._rep.CreateMember(&member.Member{Id: memberId}).Invitation()
	total, rows = iv.GetInvitationMembers(begin, end)
	if l := len(rows); l > 0 {
		arr := make([]int, l)
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
func (m *memberService) GetReferNum(memberId int, layer int) int {
	return m._query.GetReferNum(memberId, layer)
}

// 获取会员最后更新时间
func (ms *memberService) GetMemberLatestUpdateTime(memberId int) int64 {
	return ms._rep.GetMemberLatestUpdateTime(memberId)
}

func (ms *memberService) GetMemberList(ids []int) []*dto.MemberSummary {
	list := ms._query.GetMemberList(ids)
	for _, v := range list {
		v.Avatar = format.GetResUrl(v.Avatar)
	}
	return list
}

// 获取会员汇总信息
func (ms *memberService) GetMemberSummary(memberId int) *dto.MemberSummary {
	var m member.IMember = ms._rep.GetMember(memberId)
	if m != nil {
		mv := m.GetValue()
		acv := m.GetAccount().GetValue()
		lv := m.GetLevel()
		pro := m.Profile().GetProfile()
		return &dto.MemberSummary{
			Id:                m.GetAggregateRootId(),
			Usr:               mv.Usr,
			Name:              pro.Name,
			Avatar:            format.GetResUrl(pro.Avatar),
			Exp:               mv.Exp,
			Level:             mv.Level,
			LevelName:         lv.Name,
			Integral:          acv.Integral,
			Balance:           acv.Balance,
			PresentBalance:    acv.PresentBalance,
			GrowBalance:       acv.GrowBalance,
			GrowAmount:        acv.GrowAmount,
			GrowEarnings:      acv.GrowEarnings,
			GrowTotalEarnings: acv.GrowTotalEarnings,
			UpdateTime:        mv.UpdateTime,
		}
	}
	return nil
}

// 获取余额变动信息
func (ms *memberService) GetBalanceInfoById(memberId, infoId int) *member.BalanceInfo {
	m := ms._rep.GetMember(memberId)
	if m == nil {
		return nil
	}
	return m.GetAccount().GetBalanceInfo(infoId)
}

// 充值
func (ms *memberService) Charge(memberId, chargeType int, title,
	outerNo string, amount float32, relateUser int) error {
	//todo: ???
	if relateUser == 0 {
		relateUser = 1
	}
	m, err := ms.getMember(memberId)
	if err != nil {
		return err
	}
	return m.GetAccount().ChargeForBalance(chargeType, title,
		outerNo, amount, relateUser)
}

// 增加积分
func (ms *memberService) AddIntegral(memberId int, iType int,
	orderNo string, value int, remark string) error {
	m := ms._rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().AddIntegral(iType, orderNo, value, remark)
}

// 赠送金额充值
func (ms *memberService) PresentBalance(memberId int, title string,
	outerNo string, amount float32, relateUser int) error {
	m, err := ms.getMember(memberId)
	if err != nil {
		return err
	}
	return m.GetAccount().ChargeForPresent(title, outerNo, amount, relateUser)
}

// 赠送金额充值
func (ms *memberService) PresentBalanceByKind(memberId int, kind int, title string,
	outerNo string, amount float32, relateUser int) error {
	m, err := ms.getMember(memberId)
	if err != nil {
		return err
	}
	return m.GetAccount().ChargePresentByKind(kind, title, outerNo, amount, relateUser)
}

// 冻结积分,当new为true不扣除积分,反之扣除积分
func (ms *memberService) FreezesIntegral(memberId int, value int,
	new bool, remark string) error {
	m := ms._rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().FreezesIntegral(value, new, remark)
}

// 解冻积分
func (ms *memberService) UnfreezesIntegral(memberId int,
	value int, remark string) error {
	m := ms._rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().UnfreezesIntegral(value, remark)
}

// 扣减奖金
func (ms *memberService) DiscountPresent(memberId int, title string,
	tradeNo string, amount float32, mustLargeZero bool) error {
	m, err := ms.getMember(memberId)
	if err != nil {
		return err
	}
	return m.GetAccount().DiscountPresent(title, tradeNo, amount,
		member.DefaultRelateUser, mustLargeZero)
}

// 流通账户
func (ms *memberService) ChargeFlowBalance(memberId int, title string,
	tradeNo string, amount float32) error {
	m, err := ms.getMember(memberId)
	if err != nil {
		return err
	}
	return m.GetAccount().ChargeFlowBalance(title, tradeNo, amount)
}

// 验证交易密码
func (ms *memberService) VerifyTradePwd(memberId int, tradePwd string) (bool, error) {
	m := ms.GetMember(memberId)
	if len(m.TradePwd) == 0 {
		return false, member.ErrNotSetTradePwd
	}
	if m.TradePwd != tradePwd {
		return false, member.ErrIncorrectTradePwd
	}
	return true, nil
}

func (ms *memberService) GetPresentLog(memberId int, logId int) *member.PresentLog {
	m := ms._rep.GetMember(memberId)
	return m.GetAccount().GetPresentLog(logId)
}

// 提现并返回提现编号,交易号以及错误信息
func (ms *memberService) SubmitApplyPresentBalance(memberId int, applyType int,
	applyAmount float32, commission float32) (int, string, error) {
	m, err := ms.getMember(memberId)
	if err != nil {
		return 0, "", err
	}

	acc := m.GetAccount()
	var title string
	switch applyType {
	case member.KindPresentTakeOutToBankCard:
		title = "提现到银行卡"
	case member.KindPresentTakeOutToBalance:
		title = "充值账户"
	case member.KindPresentTakeOutToThirdPart:
		title = "充值到第三方账户"
	}
	return acc.RequestTakeOut(applyType, title, applyAmount, commission)
}

// 获取最近的提现
func (ms *memberService) GetLatestApplyCash(memberId int) *member.BalanceInfo {
	return ms._query.GetLatestBalanceInfoByKind(memberId, member.KindPresentTakeOutToBankCard)
}

// 获取最近的提现描述
func (ms *memberService) GetLatestApplyCashText(memberId int) string {
	var latestInfo string
	latestApplyInfo := ms.GetLatestApplyCash(memberId)
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
func (ms *memberService) ConfirmApplyCash(memberId int, infoId int, pass bool, remark string) error {
	m, err := ms.getMember(memberId)
	if err == nil {
		err = m.GetAccount().ConfirmTakeOut(infoId, pass, remark)
	}
	return err
}

// 完成提现
func (ms *memberService) FinishApplyCash(memberId, id int, tradeNo string) error {
	m, err := ms.getMember(memberId)
	if err != nil {
		return err
	}
	return m.GetAccount().FinishTakeOut(id, tradeNo)
}

// 冻结余额
func (ms *memberService) Freeze(memberId int, title string,
	tradeNo string, amount float32, referId int) error {
	m := ms._rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().Freeze(title, tradeNo, amount, referId)
}

// 解冻金额
func (ms *memberService) Unfreeze(memberId int, title string,
	tradeNo string, amount float32, referId int) error {
	m := ms._rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().Unfreeze(title, tradeNo, amount, referId)
}

// 冻结赠送金额
func (ms *memberService) FreezePresent(memberId int, title string,
	tradeNo string, amount float32, referId int) error {
	m := ms._rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().FreezePresent(title, tradeNo, amount, referId)
}

// 解冻赠送金额
func (ms *memberService) UnfreezePresent(memberId int, title string,
	tradeNo string, amount float32, referId int) error {
	m := ms._rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().UnfreezePresent(title, tradeNo, amount, referId)
}

// 将冻结金额标记为失效
func (ms *memberService) FreezeExpired(memberId int, accountKind int, amount float32,
	remark string) error {
	m := ms._rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().FreezeExpired(accountKind, amount, remark)
}

// 转账余额到其他账户
func (ms *memberService) TransferAccounts(accountKind int, fromMember int,
	toMember int, amount float32, csnRate float32, remark string) error {
	m := ms._rep.GetMember(fromMember)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().TransferAccounts(accountKind, toMember,
		amount, csnRate, remark)
}

// 转账余额到其他账户
func (ms *memberService) TransferBalance(memberId int, kind int, amount float32, tradeNo string,
	toTitle, fromTitle string) error {
	m := ms._rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().TransferBalance(kind, amount, tradeNo, toTitle, fromTitle)
}

// 转账返利账户,kind为转账类型，如 KindBalanceTransfer等
// commission手续费
func (ms *memberService) TransferPresent(memberId int, kind int, amount float32, commission float32,
	tradeNo string, toTitle string, fromTitle string) error {
	m := ms._rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().TransferPresent(kind, amount, commission,
		tradeNo, toTitle, fromTitle)
}

// 转账活动账户,kind为转账类型，如 KindBalanceTransfer等
// commission手续费
func (ms *memberService) TransferFlow(memberId int, kind int, amount float32,
	commission float32, tradeNo string, toTitle string, fromTitle string) error {
	m := ms._rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().TransferFlow(kind, amount, commission, tradeNo,
		toTitle, fromTitle)
}

// 将活动金转给其他人
func (ms *memberService) TransferFlowTo(memberId int, toMemberId int, kind int,
	amount float32, commission float32, tradeNo string, toTitle string, fromTitle string) error {
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
func (ms *memberService) GetMemberIdByPhone(phone string) int {
	return ms._query.GetMemberIdByPhone(phone)
}

// 会员推广排名
func (ms *memberService) GetMemberInviRank(merchantId int, allTeam bool, levelComp string, level int,
	startTime int64, endTime int64, num int) []*dto.RankMember {
	return ms._query.GetMemberInviRank(merchantId, allTeam, levelComp, level, startTime, endTime, num)
}

// 生成会员账户人工单据
func (ms *memberService) NewBalanceTicket(merchantId int, memberId int, accountType int,
	tit string, amount float32, relateUser int) (string, error) {
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
	if accountType == member.AccountPresent {
		outerNo = domain.NewTradeNo(merchantId)
		if amount > 0 {
			//增加奖金
			tit2 = "[KF]客服调整-" + variable.AliasPresentAccount
			if len(tit) > 0 {
				tit2 = tit2 + "(" + tit + ")"
			}
			err = acc.ChargeForPresent(tit2, outerNo, amount, relateUser)
		} else {
			//扣减奖金
			tit2 = "[KF]客服扣减-" + variable.AliasPresentAccount
			if len(tit) > 0 {
				tit2 = tit2 + "(" + tit + ")"
			}
			err = acc.DiscountPresent(tit2, outerNo, -amount, relateUser, false)
		}
		return outerNo, err
	}

	if accountType == member.AccountBalance {
		outerNo = domain.NewTradeNo(merchantId)
		if amount > 0 {
			tit2 = "[KF]客服充值"
			if len(tit) > 0 {
				tit2 = tit2 + "(" + tit + ")"
			}
			err = acc.ChargeForBalance(member.ChargeByService,
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
func (ms *memberService) PagedAvailableCoupon(memberId, start, end int) (int, []*dto.SimpleCoupon) {
	return ms._rep.CreateMemberById(memberId).
		GiftCard().PagedAvailableCoupon(start, end)
}

// 已使用的优惠券
func (ms *memberService) PagedAllCoupon(memberId, start, end int) (int, []*dto.SimpleCoupon) {
	return ms._rep.CreateMemberById(memberId).
		GiftCard().PagedAllCoupon(start, end)
}

// 过期的优惠券
func (ms *memberService) PagedExpiresCoupon(memberId, start, end int) (int, []*dto.SimpleCoupon) {
	return ms._rep.CreateMemberById(memberId).
		GiftCard().PagedExpiresCoupon(start, end)
}
