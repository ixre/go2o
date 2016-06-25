/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-09 10:13
 * description :
 * history :
 */

package repository

import (
	"database/sql"
	"fmt"
	"github.com/jsix/gof"
	"github.com/jsix/gof/db"
	"go2o/core"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/mss"
	"go2o/core/domain/interface/valueobject"
	memberImpl "go2o/core/domain/member"
	"go2o/core/dto"
	"go2o/core/variable"
	"strconv"
	"strings"
	"sync"
	"time"
)

var _ member.IMemberRep = new(MemberRep)
var (
	memberManager member.IMemberManager
	memberMux     sync.Mutex
)

type MemberRep struct {
	db.Connector
	_partnerRep merchant.IMerchantRep
	_valRep     valueobject.IValueRep
	_mssRep     mss.IMssRep
}

func NewMemberRep(c db.Connector, mssRep mss.IMssRep,
	valRep valueobject.IValueRep) *MemberRep {
	return &MemberRep{
		Connector: c,
		_mssRep:   mssRep,
		_valRep:   valRep,
	}
}

// 获取管理服务
func (this *MemberRep) GetManager() member.IMemberManager {
	memberMux.Lock()
	if memberManager == nil {
		memberManager = memberImpl.NewMemberManager(this, this._valRep)
	}
	memberMux.Unlock()
	return memberManager
}

// 获取资料或初始化
func (this *MemberRep) GetProfile(memberId int) *member.Profile {
	e := member.Profile{}
	orm := this.Connector.GetOrm()
	if orm.Get(memberId, &e) != nil {
		e.MemberId = memberId
		e.UpdateTime = time.Now().Unix()
		orm.Save(nil, &e)
	}
	return &e
}

// 保存资料
func (this *MemberRep) SaveProfile(v *member.Profile) error {
	_, _, err := this.Connector.GetOrm().Save(v.MemberId, v)
	return err
}

//收藏,typeId 为类型编号, referId为关联的ID
func (this *MemberRep) Favorite(memberId int, favType, referId int) error {
	_, _, err := this.Connector.GetOrm().Save(nil, &member.Favorite{
		MemberId:   memberId,
		FavType:    favType,
		ReferId:    referId,
		UpdateTime: time.Now().Unix(),
	})
	return err
}

//是否已收藏
func (this *MemberRep) Favored(memberId, favType, referId int) bool {
	num := 0
	this.Connector.ExecScalar(`SELECT COUNT(0) FROM mm_favorite
	WHERE member_id=? AND fav_type=? AND refer_id=?`, &num,
		memberId, favType, referId)
	return num > 0
}

//取消收藏
func (this *MemberRep) CancelFavorite(memberId int, favType, referId int) error {
	_, err := this.Connector.GetOrm().Delete(&member.Favorite{},
		"member_id=? AND fav_type=? AND refer_id=?",
		memberId, favType, referId)
	return err
}

// 获取会员等级
func (this *MemberRep) GetMemberLevels_New() []*member.Level {
	list := []*member.Level{}
	this.Connector.GetOrm().Select(&list, "enabled=1")
	return list
}

// 获取等级对应的会员数
func (this *MemberRep) GetMemberNumByLevel_New(id int) int {
	total := 0
	this.Connector.ExecScalar("SELECT COUNT(0) FROM mm_member WHERE level=?", &total, id)
	return total
}

// 删除会员等级
func (this *MemberRep) DeleteMemberLevel_New(id int) error {
	return this.Connector.GetOrm().DeleteByPk(&member.Level{}, id)
}

// 保存会员等级
func (this *MemberRep) SaveMemberLevel_New(v *member.Level) (int, error) {
	orm := this.Connector.GetOrm()
	var err error
	if v.Id > 0 {
		_, _, err = orm.Save(v.Id, v)
	} else {
		var id int64
		_, id, err = orm.Save(nil, v)
		v.Id = int(id)
	}
	return v.Id, err
}

func (this *MemberRep) SetMerchantRep(partnerRep merchant.IMerchantRep) {
	this._partnerRep = partnerRep
}

// 根据用户名获取会员
func (this *MemberRep) GetMemberValueByUsr(usr string) *member.Member {
	e := &member.Member{}
	err := this.Connector.GetOrm().GetBy(e, "usr=?", usr)
	if err != nil {
		return nil
	}
	return e
}

// 根据手机号码获取会员
func (this *MemberRep) GetMemberValueByPhone(phone string) *member.Member {
	e := &member.Member{}
	err := this.Connector.GetOrm().GetBy(e, "phone=?", phone)
	if err != nil {
		return nil
	}
	return e
}

// 获取会员
func (this *MemberRep) GetMember(memberId int) member.IMember {
	e := &member.Member{}
	err := this.Connector.GetOrm().Get(memberId, e)
	if err == nil {
		return this.CreateMember(e)
	}
	return nil
}

func (this *MemberRep) GetMemberIdByUser(user string) int {
	var id int
	this.Connector.ExecScalar("SELECT id FROM mm_member WHERE usr = ?", &id, user)
	return id
}

// 创建会员
func (this *MemberRep) CreateMember(v *member.Member) member.IMember {
	return memberImpl.NewMember(this.GetManager(), v, this,
		this._mssRep, this._valRep, this._partnerRep)
}

// 创建会员,仅作为某些操作使用,不保存
func (this *MemberRep) CreateMemberById(memberId int) member.IMember {
	return this.CreateMember(&member.Member{Id: memberId})
}

// 根据邀请码获取会员编号
func (this *MemberRep) GetMemberIdByInvitationCode(code string) int {
	var memberId int
	this.ExecScalar("SELECT id FROM mm_member WHERE invitation_code=?", &memberId, code)
	return memberId
}

func (this *MemberRep) GetLevel(merchantId, levelValue int) *merchant.MemberLevel {
	var m merchant.MemberLevel
	err := this.Connector.GetOrm().GetBy(&m, "merchant_id=? AND value = ?", merchantId, levelValue)
	if err != nil {
		return nil
	}
	return &m
}

// 获取下一个等级
func (this *MemberRep) GetNextLevel(merchantId, levelVal int) *merchant.MemberLevel {
	var m merchant.MemberLevel
	err := this.Connector.GetOrm().GetBy(&m, "merchant_id=? AND value>? LIMIT 0,1", merchantId, levelVal)
	if err != nil {
		return nil
	}
	return &m
}

// 获取会员等级
func (this *MemberRep) GetMemberLevels(merchantId int) []*merchant.MemberLevel {
	list := []*merchant.MemberLevel{}
	this.Connector.GetOrm().Select(&list,
		"merchant_id=?", merchantId)
	return list
}

// 删除会员等级
func (this *MemberRep) DeleteMemberLevel(merchantId, id int) error {
	_, err := this.Connector.GetOrm().Delete(&merchant.MemberLevel{},
		"id=? AND merchant_id=?", id, merchantId)
	return err
}

// 保存等级
func (this *MemberRep) SaveMemberLevel(merchantId int, v *merchant.MemberLevel) (int, error) {
	orm := this.Connector.GetOrm()
	var err error
	if v.Id > 0 {
		_, _, err = orm.Save(v.Id, v)
	} else {
		_, _, err = orm.Save(nil, v)
		this.Connector.ExecScalar(`SELECT MAX(id) FROM pt_member_level`, &v.Id)

	}
	return v.Id, err
}

// 获取会员最后更新时间
func (this *MemberRep) GetMemberLatestUpdateTime(id int) int64 {
	var updateTime int64
	this.Connector.ExecScalar(`SELECT update_time FROM mm_member where id=?`, &updateTime, id)
	return updateTime
}

// 获取账户
func (this *MemberRep) GetAccount(memberId int) *member.AccountValue {
	e := new(member.AccountValue)
	this.Connector.GetOrm().Get(memberId, e)
	return e
}

// 保存账户，传入会员编号
func (this *MemberRep) SaveAccount(v *member.AccountValue) (int, error) {
	_, _, err := this.Connector.GetOrm().Save(v.MemberId, v)
	this.pushToAccountUpdateQueue(v.MemberId, v.UpdateTime)
	return v.MemberId, err
}

func (this *MemberRep) pushToAccountUpdateQueue(memberId int, updateTime int64) {
	rc := core.GetRedisConn()
	defer rc.Close()
	// 保存最后更新时间
	mutKey := fmt.Sprintf("%s%d", variable.KvAccountUpdateTime, memberId)
	rc.Do("SETEX", mutKey, 3600*400, updateTime)
	// push to tcp notify queue
	rc.Do("RPUSH", variable.KvAccountUpdateTcpNotifyQueue, memberId)
}

// 获取银行信息
func (this *MemberRep) GetBankInfo(memberId int) *member.BankInfo {
	e := new(member.BankInfo)
	this.Connector.GetOrm().Get(memberId, e)
	return e
}

// 保存银行信息
func (this *MemberRep) SaveBankInfo(v *member.BankInfo) error {
	var err error
	_, _, err = this.Connector.GetOrm().Save(v.MemberId, v)
	return err
}

// 保存积分记录
func (this *MemberRep) SaveIntegralLog(l *member.IntegralLog) error {
	orm := this.Connector.GetOrm()
	var err error
	if l.Id > 0 {
		_, _, err = orm.Save(l.Id, l)
	} else {
		_, _, err = orm.Save(nil, l)
	}
	return err
}

// 获取会员关联
func (this *MemberRep) GetRelation(memberId int) *member.MemberRelation {
	var e member.MemberRelation
	if this.Connector.GetOrm().Get(memberId, &e) != nil {
		return nil
	}
	return &e
}

// 获取积分对应的等级
func (this *MemberRep) GetLevelValueByExp(merchantId int, exp int) int {
	var levelId int
	this.Connector.ExecScalar(`SELECT lv.value FROM pt_member_level lv
	 	where lv.merchant_id=? AND lv.require_exp <= ? AND lv.enabled=1
	 	 ORDER BY lv.require_exp DESC LIMIT 0,1`,
		&levelId, merchantId, exp)
	return levelId

}

// 锁定会员
func (this *MemberRep) LockMember(id int, state int) error {
	_, err := this.Connector.ExecNonQuery("update mm_member set state=? WHERE id=?", state, id)
	return err
}

// 保存会员
func (this *MemberRep) SaveMember(v *member.Member) (int, error) {
	if v.Id > 0 {
		rc := core.GetRedisConn()
		defer rc.Close()
		// 保存最后更新时间 todo:
		mutKey := fmt.Sprintf("%s%d", variable.KvMemberUpdateTime, v.Id)
		rc.Do("SETEX", mutKey, 3600*400, v.UpdateTime)
		rc.Do("RPUSH", variable.KvMemberUpdateTcpNotifyQueue, v.Id) // push to tcp notify queue

		// 保存会员信息
		_, _, err := this.Connector.GetOrm().Save(v.Id, v)
		if err == nil {
			rc.Do("RPUSH", variable.KvMemberUpdateQueue, fmt.Sprintf("%d-update", v.Id)) // push to queue
		}
		return v.Id, err
	}
	return this.createMember(v)
}

func (this *MemberRep) createMember(v *member.Member) (int, error) {
	var id int64
	_, id, err := this.Connector.GetOrm().Save(nil, v)
	if err != nil {
		return -1, err
	}
	v.Id = int(id)
	this.initMember(v)

	rc := core.GetRedisConn()
	defer rc.Close()
	rc.Do("RPUSH", variable.KvMemberUpdateQueue,
		fmt.Sprintf("%d-create", v.Id)) // push to queue

	// 更新会员数 todo: 考虑去掉
	var total = 0
	this.Connector.ExecScalar("SELECT COUNT(0) FROM mm_member", &total)
	gof.CurrentApp.Storage().Set(variable.KvTotalMembers, total)

	return v.Id, err
}

func (this *MemberRep) initMember(v *member.Member) {

	orm := this.Connector.GetOrm()
	orm.Save(nil, &member.AccountValue{
		MemberId:    v.Id,
		Balance:     0,
		TotalFee:    0,
		TotalCharge: 0,
		TotalPay:    0,
		UpdateTime:  v.RegTime,
	})

	orm.Save(nil, &member.BankInfo{
		MemberId: v.Id,
		State:    1,
	})

	orm.Save(nil, &member.MemberRelation{
		MemberId:           v.Id,
		CardId:             "",
		RefereesId:         0,
		RegisterMerchantId: 0,
	})
}

// 用户名是否存在
func (this *MemberRep) CheckUsrExist(usr string, memberId int) bool {
	var c int
	this.Connector.ExecScalar("SELECT COUNT(0) FROM mm_member WHERE usr=? AND id<>?",
		&c, usr, memberId)
	return c != 0
}

// 手机号码是否使用
func (this *MemberRep) CheckPhoneBind(phone string, memberId int) bool {
	var c int
	this.Connector.ExecScalar("SELECT COUNT(0) FROM mm_profile WHERE phone=? AND member_id<>?",
		&c, phone, memberId)
	return c != 0
}

// 保存绑定
func (this *MemberRep) SaveRelation(v *member.MemberRelation) error {
	_, _, err := this.Connector.GetOrm().Save(v.MemberId, v)
	return err
}

// 保存地址
func (this *MemberRep) SaveDeliver(v *member.DeliverAddress) (int, error) {
	orm := this.Connector.GetOrm()
	if v.Id <= 0 {
		_, _, err := orm.Save(nil, v)
		this.Connector.ExecScalar("SELECT MAX(id) FROM mm_delivery_addr", &v.Id)
		return v.Id, err
	} else {
		_, _, err := orm.Save(v.Id, v)
		return v.Id, err
	}
}

// 获取全部配送地址
func (this *MemberRep) GetDeliverAddress(memberId int) []*member.DeliverAddress {
	addresses := []*member.DeliverAddress{}
	this.Connector.GetOrm().Select(&addresses, "member_id=?", memberId)
	return addresses
}

// 获取配送地址
func (this *MemberRep) GetSingleDeliverAddress(memberId, deliverId int) *member.DeliverAddress {
	var address member.DeliverAddress
	err := this.Connector.GetOrm().Get(deliverId, &address)

	if err == nil && address.MemberId == memberId {
		return &address
	}
	return nil
}

// 删除配送地址
func (this *MemberRep) DeleteDeliver(memberId, deliverId int) error {
	_, err := this.Connector.ExecNonQuery(
		"DELETE FROM mm_deliver_addr WHERE member_id=? AND id=?",
		memberId, deliverId)
	return err
}

// 邀请
func (this *MemberRep) GetMyInvitationMembers(memberId, begin, end int) (
	total int, rows []*member.Member) {
	arr := []*member.Member{}
	this.Connector.ExecScalar(`SELECT COUNT(0) FROM mm_member WHERE id IN
	 (SELECT member_id FROM mm_relation WHERE invi_member_id=?)`, &total, memberId)
	if total > 0 {

		this.Connector.GetOrm().SelectByQuery(&arr, `SELECT * FROM mm_member
	    WHERE id IN (SELECT member_id FROM mm_relation WHERE invi_member_id=?)
	 ORDER BY level DESC,id LIMIT ?,?`, memberId, begin, end-begin)
	}
	return total, arr
}

// 获取下级会员数量
func (this *MemberRep) GetSubInvitationNum(memberId int, memberIdArr []int) map[int]int {
	if len(memberIdArr) == 0 {
		return map[int]int{}
	}
	var ids []string = make([]string, len(memberIdArr))
	for i, v := range memberIdArr {
		ids[i] = strconv.Itoa(v)
	}
	memberIds := strings.Join(ids, ",")
	var d map[int]int = make(map[int]int)
	err := this.Connector.Query(fmt.Sprintf("SELECT r1.member_id,"+
		"(SELECT COUNT(0) FROM mm_relation r2 WHERE r2.invi_member_id=r1.member_id)"+
		"as num FROM mm_relation r1 WHERE r1.member_id IN(%s)", memberIds),
		func(rows *sql.Rows) {
			var id, num int
			for rows.Next() {
				rows.Scan(&id, &num)
				d[id] = num
			}
			rows.Close()
		})
	handleError(err)
	return d
}

// 获取推荐我的人
func (this *MemberRep) GetInvitationMeMember(memberId int) *member.Member {
	var d *member.Member = new(member.Member)
	err := this.Connector.GetOrm().GetByQuery(d,
		"SELECT * FROM mm_member WHERE id =(SELECT invi_member_id FROM mm_relation  WHERE id=?)",
		memberId)

	if err != nil {
		return nil
	}
	return d
}

// 根据编号获取余额变动信息
func (this *MemberRep) GetBalanceInfo(id int) *member.BalanceInfoValue {
	var e member.BalanceInfoValue
	if err := this.Connector.GetOrm().Get(id, &e); err == nil {
		return &e
	}
	return nil
}

// 根据号码获取余额变动信息
func (this *MemberRep) GetBalanceInfoByNo(tradeNo string) *member.BalanceInfoValue {
	var e member.BalanceInfoValue
	if err := this.Connector.GetOrm().GetBy(&e, "trade_no=?", tradeNo); err == nil {
		return &e
	}
	return nil
}

// 保存余额变动信息
func (this *MemberRep) SaveBalanceInfo(v *member.BalanceInfoValue) (int, error) {
	var err error
	var orm = this.Connector.GetOrm()
	if v.Id > 0 {
		_, _, err = orm.Save(v.Id, v)
	} else {
		_, _, err = orm.Save(nil, v)
		this.Connector.ExecScalar("SELECT MAX(id) FROM mm_balance_info WHERE member_id=?", &v.Id, v.MemberId)
	}
	return v.Id, err
}

// 保存理财账户信息
func (this *MemberRep) SaveGrowAccount(memberId int, balance, totalAmount,
	growEarnings, totalGrowEarnings float32, updateTime int64) error {
	_, err := this.Connector.ExecNonQuery(`UPDATE mm_account SET grow_balance=?,
		grow_amount=?,grow_earnings=?,grow_total_earnings=?,update_time=? where member_id=?`,
		balance, totalAmount, growEarnings, totalGrowEarnings, updateTime, memberId)
	this.pushToAccountUpdateQueue(memberId, updateTime)
	return err
}

// 获取会员分页的优惠券列表
func (this *MemberRep) GetMemberPagedCoupon(memberId, start, end int, where string) (total int, rows []*dto.ValueCoupon) {
	list := []*dto.ValueCoupon{}
	this.Connector.ExecScalar(fmt.Sprintf(`SELECT COUNT(distinct pi.id)
        FROM pm_info pi INNER JOIN pm_coupon c ON c.id = pi.id
	    INNER JOIN pm_coupon_bind pb ON pb.coupon_id=pi.id
	    WHERE member_id=? AND %s`, where), &total, memberId)
	if total > 0 {
		this.Connector.GetOrm().SelectByQuery(&list,
			fmt.Sprintf(`SELECT pi.id,SUM(1) as num,pi.short_name as title,
            code,fee,c.discount,is_used,over_time FROM pm_info pi
             INNER JOIN pm_coupon c ON c.id = pi.id
	        INNER JOIN pm_coupon_bind pb ON pb.coupon_id=pi.id
	        WHERE member_id=? AND %s GROUP BY pi.id order by bind_time DESC LIMIT ?,?`, where),
			memberId, start, end-start)
	}
	return total, list
}
