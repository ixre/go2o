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
	"go2o/core/domain/interface/member"
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
}

func NewMemberService(mchService *merchantService, rep member.IMemberRep,
	q *query.MemberQuery) *memberService {
	return &memberService{
		_rep:            rep,
		_query:          q,
		_partnerService: mchService,
	}
}

// 获取资料
func (this *memberService) GetProfile(memberId int) member.Profile {
	m := this._rep.GetMember(memberId)
	if m != nil {
		return m.Profile().GetProfile()
	}
	return member.Profile{}
}

// 保存资料
func (this *memberService) SaveProfile(memberId int, v *member.Profile) error {
	m := this._rep.GetMember(memberId)
	return m.Profile().SaveProfile(v)
}

// 是否已收藏
func (this *memberService) Favored(memberId, favType, referId int) bool {
	return this._rep.CreateMemberById(memberId).
		Favorite().Favored(favType, referId)
}

// 取消收藏
func (this *memberService) Cancel(memberId, favType, referId int) error {
	return this._rep.CreateMemberById(memberId).
		Favorite().Cancel(favType, referId)
}

// 收藏商品
func (this *memberService) FavoriteGoods(memberId, goodsId int) error {
	return this._rep.CreateMemberById(memberId).
		Favorite().FavoriteGoods(goodsId)
}

// 取消收藏商品
func (this *memberService) CancelGoodsFavorite(memberId, goodsId int) error {
	return this._rep.CreateMemberById(memberId).
		Favorite().CancelGoodsFavorite(goodsId)
}

// 收藏店铺
func (this *memberService) FavoriteShop(memberId, shopId int) error {
	return this._rep.CreateMemberById(memberId).
		Favorite().FavoriteShop(shopId)
}

// 取消收藏店铺
func (this *memberService) CancelShopFavorite(memberId, shopId int) error {
	return this._rep.CreateMemberById(memberId).
		Favorite().CancelShopFavorite(shopId)
}

// 商品是否已收藏
func (this *memberService) GoodsFavored(memberId, goodsId int) bool {
	return this._rep.CreateMemberById(memberId).
		Favorite().GoodsFavored(goodsId)
}

// 商店是否已收藏
func (this *memberService) ShopFavored(memberId, shopId int) bool {
	return this._rep.CreateMemberById(memberId).
		Favorite().ShopFavored(shopId)
}

/**================ 会员等级 ==================**/
// 获取所有会员等级
func (this *memberService) GetMemberLevels() []*member.Level {
	return this._rep.GetManager().LevelManager().GetLevelSet()
}

// 根据编号获取会员等级信息
func (this *memberService) GetMemberLevelById(id int) *member.Level {
	return this._rep.GetManager().LevelManager().GetLevelById(id)
}

// 保存会员等级信息
func (this *memberService) SaveMemberLevel(v *member.Level) (int, error) {
	return this._rep.GetManager().LevelManager().SaveLevel(v)
}

// 删除会员等级
func (this *memberService) DelMemberLevel(levelId int) error {
	return this._rep.GetManager().LevelManager().DeleteLevel(levelId)
}

// 获取下一个等级
func (this *memberService) GetNextLevel(levelId int) *member.Level {
	return this._rep.GetManager().LevelManager().GetNextLevelById(levelId)
}

func (this *memberService) GetMember(id int) *member.Member {
	if id <= 0 {
		return nil
	}
	v := this._rep.GetMember(id)
	if v != nil {
		nv := v.GetValue()
		return &nv
	}
	return nil
}

func (this *memberService) getMember(memberId int) (
	member.IMember, error) {
	if memberId <= 0 {
		return nil, member.ErrNoSuchMember
	}
	m := this._rep.GetMember(memberId)
	if m == nil {
		return m, member.ErrNoSuchMember
	}
	return m, nil
}

func (this *memberService) GetMemberIdByInvitationCode(code string) int {
	return this._rep.GetMemberIdByInvitationCode(code)
}

// 更改会员用户名
func (this *memberService) ChangeUsr(id int, usr string) error {
	m := this._rep.GetMember(id)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.ChangeUsr(usr)
}

// 保存用户
func (this *memberService) SaveMember(v *member.Member) (int, error) {
	if v.Id > 0 {
		return this.updateMember(v)
	}
	return -1, errors.New("Create member use \"RegisterMember\" method.")
}

func (this *memberService) SaveMemberProfile(id int, v *member.Member) (int, error) {
	//todo:
	return -1, nil
}

func (this *memberService) updateMember(v *member.Member) (int, error) {
	m := this._rep.GetMember(v.Id)
	if m == nil {
		return -1, member.ErrNoSuchMember
	}
	if err := m.SetValue(v); err != nil {
		return m.GetAggregateRootId(), err
	}
	return m.Save()
}

// 注册会员
func (this *memberService) RegisterMember(merchantId int, v *member.Member,
	pro *member.Profile, cardId string, invitationCode string) (int, error) {

	//先验证手机
	if len(pro.Phone) > 0 {
		if b := this._rep.CheckPhoneBind(pro.Phone, v.Id); b {
			return -1, member.ErrPhoneHasBind
		}
	}

	// todo: 检测注册权限,这里应删除。应用到 member的create方法
	err := this._rep.GetManager().RegisterPerm(len(invitationCode) > 0)
	if err == nil {
		var invitationId int = 0
		if len(invitationCode) > 0 {
			//判断邀请码是否正确
			invitationId = this.GetMemberIdByInvitationCode(invitationCode)
			if invitationId <= 0 {
				return -1, member.ErrInvitationCode
			}
		}
		m := this._rep.CreateMember(v) //创建会员
		id, err := m.Save()
		if err == nil {
			pro.MemberId = id
			if len(pro.Name) == 0 { //如果未设置昵称,则默认为用户名
				pro.Name = "用户:" + v.Usr
			}
			if len(pro.Avatar) == 0 {
				pro.Avatar = "res/no_avatar.gif"
			}
			err = m.Profile().SaveProfile(pro)
			if err == nil {
				// 保存关联信息
				rl := m.GetRelation()
				rl.RefereesId = invitationId
				rl.RegisterMerchantId = merchantId
				rl.CardId = cardId
				return id, m.SaveRelation(rl)
			}
		}
		return id, err
	}
	return -1, err
}

func (this *memberService) GetLevel(memberId int) *member.Level {
	if m := this._rep.GetMember(memberId); m != nil {
		return m.GetLevel()
	}
	return nil
}

func (this *memberService) GetRelation(memberId int) *member.Relation {
	return this._rep.GetRelation(memberId)
}

// 锁定/解锁会员
func (this *memberService) LockMember(id int) (bool, error) {
	m := this._rep.GetMember(id)
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
func (this *memberService) ProfileCompleted(memberId int) bool {
	m := this._rep.GetMember(memberId)
	if m != nil {
		return m.Profile().ProfileCompleted()
	}
	return false
}

// 重置密码
func (this *memberService) ResetPassword(memberId int) string {
	m := this._rep.GetMember(memberId)
	if m != nil {
		newPwd := domain.GenerateRandomIntPwd(6)
		newEncPwd := domain.MemberSha1Pwd(newPwd)
		if m.Profile().ModifyPassword(newEncPwd, "") == nil {
			return newPwd
		}
	}
	return ""
}

// 检查凭据, update:是否更新登陆时间
func (this *memberService) TryLogin(usr, pwd string, update bool) (
	*member.Member, error) {
	usr = strings.ToLower(strings.TrimSpace(usr))

	val := this._rep.GetMemberValueByUsr(usr)
	if val == nil {
		val = this._rep.GetMemberValueByPhone(usr)
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

	m := this._rep.GetMember(val.Id)

	if update {
		unix := time.Now().Unix()
		val.LastLoginTime = unix
		val.UpdateTime = unix
		m.Save()
	}

	return val, nil
}

// 检查与现有用户不同的用户是否存在,如存在则返回错误
func (this *memberService) CheckUsr(usr string, memberId int) error {
	if len(usr) < 6 {
		return member.ErrUsrLength
	}
	if this._rep.CheckUsrExist(usr, memberId) {
		return member.ErrUsrExist
	}
	return nil
}

func (this *memberService) GetAccount(memberId int) *member.Account {
	m := this._rep.CreateMember(&member.Member{Id: memberId})
	//m, _ := this._memberRep.GetMember(memberId)
	//m.AddExp(300)
	return m.GetAccount().GetValue()
}

func (this *memberService) GetBank(memberId int) *member.BankInfo {
	m := this._rep.CreateMember(&member.Member{Id: memberId})
	b := m.Profile().GetBank()
	return &b
}

func (this *memberService) SaveBankInfo(v *member.BankInfo) error {
	m := this._rep.CreateMember(&member.Member{Id: v.MemberId})
	return m.Profile().SaveBank(v)
}

// 解锁银行卡信息
func (this *memberService) UnlockBankInfo(memberId int) error {
	m := this._rep.CreateMember(&member.Member{Id: memberId})
	return m.Profile().UnlockBank()
}

// 实名认证信息
func (this *memberService) GetTrustedInfo(memberId int) member.TrustedInfo {
	m := this._rep.GetMember(memberId)
	return m.Profile().GetTrustedInfo()
}

// 保存实名认证信息
func (this *memberService) SaveTrustedInfo(memberId int, v *member.TrustedInfo) error {
	m := this._rep.GetMember(memberId)
	err := m.Profile().SaveTrustedInfo(v)
	//todo: 取消自动认证
	if err == nil {
		err = this.ReviewTrustedInfo(memberId, true, "")
	}
	return err
}

// 审核实名认证,若重复审核将返回错误
func (this *memberService) ReviewTrustedInfo(memberId int, pass bool, remark string) error {
	m := this._rep.GetMember(memberId)
	return m.Profile().ReviewTrustedInfo(pass, remark)
}

// 获取返现记录
func (this *memberService) QueryIncomeLog(memberId, begin, end int,
	where, orderBy string) (num int, rows []map[string]interface{}) {
	return this._query.QueryBalanceLog(memberId, begin, end, where, orderBy)
}

// 查询分页订单
func (this *memberService) QueryPagerOrder(memberId, page, size int,
	where, orderBy string) (num int, rows []map[string]interface{}) {
	return this._query.QueryPagerOrder(memberId, page, size, where, orderBy)
}

/*********** 收货地址 ***********/
func (this *memberService) GetDeliverAddress(memberId int) []*member.DeliverAddress {
	return this._rep.GetDeliverAddress(memberId)
}

//获取配送地址
func (this *memberService) GetDeliverAddressById(memberId,
	deliverId int) *member.DeliverAddress {
	m := this._rep.CreateMember(&member.Member{Id: memberId})
	v := m.Profile().GetDeliver(deliverId).GetValue()
	return &v
}

//保存配送地址
func (this *memberService) SaveDeliverAddress(memberId int, e *member.DeliverAddress) (int, error) {
	m := this._rep.CreateMember(&member.Member{Id: memberId})
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
func (this *memberService) DeleteDeliverAddress(memberId int, deliverId int) error {
	m := this._rep.CreateMember(&member.Member{Id: memberId})
	return m.Profile().DeleteDeliver(deliverId)
}

func (this *memberService) ModifyPassword(memberId int, oldPwd, newPwd string) error {
	m := this._rep.GetMember(memberId)
	if m != nil {
		return m.Profile().ModifyPassword(newPwd, oldPwd)
	}
	return member.ErrNoSuchMember
}

//修改密码
func (this *memberService) ModifyTradePassword(memberId int, oldPwd, newPwd string) error {
	m := this._rep.GetMember(memberId)
	if m != nil {
		return m.Profile().ModifyTradePassword(newPwd, oldPwd)
	}
	return member.ErrNoSuchMember
}

//判断会员是否由指定会员邀请推荐的
func (this *memberService) IsInvitation(memberId int, invitationMemberId int) bool {
	m := this._rep.CreateMember(&member.Member{Id: memberId})
	return m.Invitation().InvitationBy(invitationMemberId)
}

// 获取我邀请的会员及会员邀请的人数
func (this *memberService) GetMyPagedInvitationMembers(memberId int,
	begin, end int) (total int, rows []*member.Member, num map[int]int) {
	iv := this._rep.CreateMember(&member.Member{Id: memberId}).Invitation()
	total, rows = iv.GetInvitationMembers(begin, end)
	if l := len(rows); l > 0 {
		arr := make([]int, l)
		for i := 0; i < l; i++ {
			arr[i] = rows[i].Id
		}
		num = iv.GetSubInvitationNum(arr)
	}
	return total, rows, num
}

// 获取会员最后更新时间
func (this *memberService) GetMemberLatestUpdateTime(memberId int) int64 {
	return this._rep.GetMemberLatestUpdateTime(memberId)
}

func (this *memberService) GetMemberList(ids []int) []*dto.MemberSummary {
	list := this._query.GetMemberList(ids)
	for _, v := range list {
		v.Avatar = format.GetResUrl(v.Avatar)
	}
	return list
}

// 获取会员汇总信息
func (this *memberService) GetMemberSummary(memberId int) *dto.MemberSummary {
	var m member.IMember = this._rep.GetMember(memberId)
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
func (this *memberService) GetBalanceInfoById(memberId, infoId int) *member.BalanceInfo {
	m := this._rep.GetMember(memberId)
	if m == nil {
		return nil
	}
	return m.GetAccount().GetBalanceInfo(infoId)
}

// 充值
func (this *memberService) Charge(memberId, chargeType int, title, tradeNo string, amount float32) error {
	m, err := this.getMember(memberId)
	if err != nil {
		return err
	}
	return m.GetAccount().ChargeBalance(chargeType, title, tradeNo, amount)
}

// 赠送金额充值
func (this *memberService) PresentBalance(memberId int, title string, tradeNo string, amount float32) error {
	m, err := this.getMember(memberId)
	if err != nil {
		return err
	}
	return m.GetAccount().PresentBalance(title, tradeNo, amount)
}

// 扣减奖金
func (this *memberService) DiscountPresent(memberId int, title string, tradeNo string, amount float32, mustLargeZero bool) error {
	m, err := this.getMember(memberId)
	if err != nil {
		return err
	}
	return m.GetAccount().DiscountPresent(title, tradeNo, amount, mustLargeZero)
}

// 流通账户
func (this *memberService) ChargeFlowBalance(memberId int, title string, tradeNo string, amount float32) error {
	m, err := this.getMember(memberId)
	if err != nil {
		return err
	}
	return m.GetAccount().ChargeFlowBalance(title, tradeNo, amount)
}

// 验证交易密码
func (this *memberService) VerifyTradePwd(memberId int, tradePwd string) (bool, error) {
	m := this.GetMember(memberId)
	if len(m.TradePwd) == 0 {
		return false, member.ErrNotSetTradePwd
	}
	if m.TradePwd != tradePwd {
		return false, member.ErrIncorrectTradePwd
	}
	return true, nil
}

// 提现并返回提现编号,交易号以及错误信息
func (this *memberService) SubmitApplyPresentBalance(memberId int, applyType int,
	applyAmount float32, commission float32) (int, string, error) {
	m, err := this.getMember(memberId)
	if err != nil {
		return 0, "", err
	}

	acc := m.GetAccount()
	var title string
	switch applyType {
	case member.TypeApplyCashToBank:
		title = "提现到银行卡"
	case member.TypeApplyCashToCharge:
		title = "充值账户"
	case member.TypeApplyCashToServiceProvider:
		title = "充值到第三方账户"
	}
	return acc.RequestApplyCash(applyType, title, applyAmount, commission)
}

// 获取最近的提现
func (this *memberService) GetLatestApplyCash(memberId int) *member.BalanceInfo {
	return this._query.GetLatestBalanceInfoByKind(memberId, member.KindBalanceApplyCash)
}

// 获取最近的提现描述
func (this *memberService) GetLatestApplyCashText(memberId int) string {
	var latestInfo string
	latestApplyInfo := this.GetLatestApplyCash(memberId)
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
func (this *memberService) ConfirmApplyCash(memberId int, infoId int, pass bool, remark string) error {
	m, err := this.getMember(memberId)
	if err != nil {
		return err
	}
	return m.GetAccount().ConfirmApplyCash(infoId, pass, remark)
}

// 完成提现
func (this *memberService) FinishApplyCash(memberId, id int, tradeNo string) error {
	m, err := this.getMember(memberId)
	if err != nil {
		return err
	}
	return m.GetAccount().FinishApplyCash(id, tradeNo)
}

// 冻结余额
func (this *memberService) Freezes(memberId int, title string,
	tradeNo string, amount float32, referId int) error {
	m := this._rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().Freezes(title, tradeNo, amount, referId)
}

// 解冻金额
func (this *memberService) Unfreezes(memberId int, title string,
	tradeNo string, amount float32, referId int) error {
	m := this._rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().Unfreezes(title, tradeNo, amount, referId)
}

// 冻结赠送金额
func (this *memberService) FreezesPresent(memberId int, title string,
	tradeNo string, amount float32, referId int) error {
	m := this._rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().FreezesPresent(title, tradeNo, amount, referId)
}

// 解冻赠送金额
func (this *memberService) UnfreezesPresent(memberId int, title string,
	tradeNo string, amount float32, referId int) error {
	m := this._rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().UnfreezesPresent(title, tradeNo, amount, referId)
}

// 转账余额到其他账户
func (this *memberService) TransferBalance(memberId int, kind int, amount float32, tradeNo string,
	toTitle, fromTitle string) error {
	m := this._rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().TransferBalance(kind, amount, tradeNo, toTitle, fromTitle)
}

// 转账返利账户,kind为转账类型，如 KindBalanceTransfer等
// commission手续费
func (this *memberService) TransferPresent(memberId int, kind int, amount float32, commission float32,
	tradeNo string, toTitle string, fromTitle string) error {
	m := this._rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().TransferPresent(kind, amount, commission,
		tradeNo, toTitle, fromTitle)
}

// 转账活动账户,kind为转账类型，如 KindBalanceTransfer等
// commission手续费
func (this *memberService) TransferFlow(memberId int, kind int, amount float32,
	commission float32, tradeNo string, toTitle string, fromTitle string) error {
	m := this._rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().TransferFlow(kind, amount, commission, tradeNo,
		toTitle, fromTitle)
}

// 将活动金转给其他人
func (this *memberService) TransferFlowTo(memberId int, toMemberId int, kind int,
	amount float32, commission float32, tradeNo string, toTitle string, fromTitle string) error {
	m := this._rep.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	return m.GetAccount().TransferFlowTo(toMemberId, kind, amount,
		commission, tradeNo, toTitle, fromTitle)
}

// 根据用户或手机筛选会员
func (this *memberService) FilterMemberByUsrOrPhone(key string) []*dto.SimpleMember {
	return this._query.FilterMemberByUsrOrPhone(key)
}

// 会员推广排名
func (this *memberService) GetMemberInviRank(merchantId int, allTeam bool, levelComp string, level int,
	startTime int64, endTime int64, num int) []*dto.RankMember {
	return this._query.GetMemberInviRank(merchantId, allTeam, levelComp, level, startTime, endTime, num)
}

// 生成会员账户人工单据
func (this *memberService) NewBalanceTicket(merchantId int, memberId int, kind int,
	tit string, amount float32) (string, error) {
	var err error
	var tradeNo string
	if amount == 0 {
		return "", member.ErrIncorrectAmount
	}
	m := this._rep.GetMember(memberId)
	if m == nil {
		return "", member.ErrNoSuchMember
	}
	acc := m.GetAccount()
	var tit2 string
	if kind == member.KindBalancePresent {
		tradeNo = domain.NewTradeNo(merchantId)
		if amount > 0 {
			//增加奖金
			tit2 = "[KF]客服调整-" + variable.AliasPresentAccount
			if len(tit) > 0 {
				tit2 = tit2 + "(" + tit + ")"
			}
			err = acc.PresentBalance(tit2, tradeNo, amount)
		} else {
			//扣减奖金
			tit2 = "[KF]客服扣减-" + variable.AliasPresentAccount
			if len(tit) > 0 {
				tit2 = tit2 + "(" + tit + ")"
			}
			err = acc.DiscountPresent(tit2, tradeNo, -amount, false)
		}
	}

	if kind == member.KindBalanceCharge {
		tradeNo = domain.NewTradeNo(merchantId)
		if amount > 0 {
			tit2 = "[KF]客服充值卡"
			if len(tit) > 0 {
				tit2 = tit2 + "(" + tit + ")"
			}
			err = acc.ChargeBalance(member.TypeBalanceServiceCharge, tit2, tradeNo, amount)
		} else {
			tit2 = "[KF]客服扣减"
			if len(tit) > 0 {
				tit2 = tit2 + "(" + tit + ")"
			}
			err = acc.DiscountBalance(tit2, tradeNo, -amount)
		}
	}

	return tradeNo, err
}

//********* 促销  **********//

// 可用的优惠券分页数据
func (this *memberService) PagedAvailableCoupon(memberId, start, end int) (int, []*dto.SimpleCoupon) {
	return this._rep.CreateMemberById(memberId).
		GiftCard().PagedAvailableCoupon(start, end)
}

// 已使用的优惠券
func (this *memberService) PagedAllCoupon(memberId, start, end int) (int, []*dto.SimpleCoupon) {
	return this._rep.CreateMemberById(memberId).
		GiftCard().PagedAllCoupon(start, end)
}

// 过期的优惠券
func (this *memberService) PagedExpiresCoupon(memberId, start, end int) (int, []*dto.SimpleCoupon) {
	return this._rep.CreateMemberById(memberId).
		GiftCard().PagedExpiresCoupon(start, end)
}
