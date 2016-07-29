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
	"github.com/jsix/gof/storage"
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
	Storage storage.Interface
	db.Connector
	_partnerRep merchant.IMerchantRep
	_valRep     valueobject.IValueRep
	_mssRep     mss.IMssRep
}

func NewMemberRep(sto storage.Interface, c db.Connector, mssRep mss.IMssRep,
	valRep valueobject.IValueRep) *MemberRep {
	return &MemberRep{
		Storage:   sto,
		Connector: c,
		_mssRep:   mssRep,
		_valRep:   valRep,
	}
}

// 获取管理服务
func (m *MemberRep) GetManager() member.IMemberManager {
	memberMux.Lock()
	if memberManager == nil {
		memberManager = memberImpl.NewMemberManager(m, m._valRep)
	}
	memberMux.Unlock()
	return memberManager
}

// 获取资料或初始化
func (m *MemberRep) GetProfile(memberId int) *member.Profile {
	e := &member.Profile{}
	key := m.getProfileCk(memberId)
	if m.Storage.Get(key, &e) != nil {
		if m.Connector.GetOrm().Get(memberId, e) != nil {
			e.MemberId = memberId
		} else {
			m.Storage.Set(key, *e)
		}
	}
	return e
}

// 保存资料
func (m *MemberRep) SaveProfile(v *member.Profile) error {
	_, _, err := m.Connector.GetOrm().Save(v.MemberId, v)
	if err == nil {
		err = m.Storage.Set(m.getProfileCk(v.MemberId), *v)
	}
	return err
}

//收藏,typeId 为类型编号, referId为关联的ID
func (m *MemberRep) Favorite(memberId int, favType, referId int) error {
	_, _, err := m.Connector.GetOrm().Save(nil, &member.Favorite{
		MemberId:   memberId,
		FavType:    favType,
		ReferId:    referId,
		UpdateTime: time.Now().Unix(),
	})
	return err
}

//是否已收藏
func (m *MemberRep) Favored(memberId, favType, referId int) bool {
	num := 0
	m.Connector.ExecScalar(`SELECT COUNT(0) FROM mm_favorite
	WHERE member_id=? AND fav_type=? AND refer_id=?`, &num,
		memberId, favType, referId)
	return num > 0
}

//取消收藏
func (m *MemberRep) CancelFavorite(memberId int, favType, referId int) error {
	_, err := m.Connector.GetOrm().Delete(&member.Favorite{},
		"member_id=? AND fav_type=? AND refer_id=?",
		memberId, favType, referId)
	return err
}

// 获取会员等级
func (m *MemberRep) GetMemberLevels_New() []*member.Level {
	list := []*member.Level{}
	m.Connector.GetOrm().Select(&list, "1=1 ORDER BY id ASC")
	return list
}

// 获取等级对应的会员数
func (m *MemberRep) GetMemberNumByLevel_New(id int) int {
	total := 0
	m.Connector.ExecScalar("SELECT COUNT(0) FROM mm_member WHERE level=?", &total, id)
	return total
}

// 删除会员等级
func (m *MemberRep) DeleteMemberLevel_New(id int) error {
	return m.Connector.GetOrm().DeleteByPk(&member.Level{}, id)
}

// 保存会员等级
func (m *MemberRep) SaveMemberLevel_New(v *member.Level) (int, error) {
	orm := m.Connector.GetOrm()
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

func (m *MemberRep) SetMerchantRep(partnerRep merchant.IMerchantRep) {
	m._partnerRep = partnerRep
}

// 根据用户名获取会员
func (m *MemberRep) GetMemberValueByUsr(usr string) *member.Member {
	e := &member.Member{}
	err := m.Connector.GetOrm().GetBy(e, "usr=?", usr)
	if err != nil {
		return nil
	}
	return e
}

// 根据手机号码获取会员
func (m *MemberRep) GetMemberValueByPhone(phone string) *member.Member {
	e := &member.Member{}
	err := m.GetOrm().GetByQuery(e, `SELECT * FROM mm_member
		INNER JOIN mm_profile ON mm_profile.member_id = mm_member.id
		 WHERE phone=?`, phone)
	if err != nil {
		return nil
	}
	return e
}

// 根据手机号获取会员编号
func (m *MemberRep) GetMemberIdByPhone(phone string) int {
	id := -1
	m.Connector.ExecScalar("SELECT member_id FROM mm_profile WHERE phone=?", &id, phone)
	return id
}

// 根据邮箱地址获取会员编号
func (m *MemberRep) GetMemberIdByEmail(email string) int {
	id := -1
	m.Connector.ExecScalar("SELECT member_id FROM mm_profile WHERE email=?", &id, email)
	return id
}

func (m *MemberRep) getMemberCk(memberId int) string {
	return fmt.Sprintf("go2o:rep:member:%d", memberId)
}
func (m *MemberRep) getAccountCk(memberId int) string {
	return fmt.Sprintf("go2o:rep:mm-acc:%d", memberId)
}
func (m *MemberRep) getProfileCk(memberId int) string {
	return fmt.Sprintf("go2o:rep:mm-pro:%d", memberId)
}
func (m *MemberRep) getTrustCk(memberId int) string {
	return fmt.Sprintf("go2o:rep:mm-trust:%d", memberId)
}
func (m *MemberRep) getGlobLevelsCk() string {
	return "go2o:rep:mm-lv"
}

// 获取会员
func (m *MemberRep) GetMember(memberId int) member.IMember {
	e := &member.Member{}
	key := m.getMemberCk(memberId)
	if err := m.Storage.Get(key, &e); err != nil {
		//log.Println("-- mm",err)
		if m.Connector.GetOrm().Get(memberId, e) != nil {
			return nil
		}
		m.Storage.Set(key, *e)
	} else {
		//log.Println(fmt.Sprintf("--- member: %d > %#v",memberId,e))
	}
	return m.CreateMember(e)
}

// 保存会员
func (m *MemberRep) SaveMember(v *member.Member) (int, error) {
	if v.Id > 0 {
		rc := core.GetRedisConn()
		defer rc.Close()
		// 保存最后更新时间
		// todo: del
		mutKey := fmt.Sprintf("%s%d", variable.KvMemberUpdateTime, v.Id)
		rc.Do("SETEX", mutKey, 3600*400, v.UpdateTime)
		rc.Do("RPUSH", variable.KvMemberUpdateTcpNotifyQueue, v.Id) // push to tcp notify queue

		// 保存会员信息
		_, _, err := m.Connector.GetOrm().Save(v.Id, v)

		if err == nil {
			// 存储到缓存中
			err = m.Storage.Set(m.getMemberCk(v.Id), *v)
			// 存储到队列
			rc.Do("RPUSH", variable.KvMemberUpdateQueue, fmt.Sprintf("%d-update", v.Id))
		}
		return v.Id, err
	}

	return m.createMember(v)
}

func (m *MemberRep) createMember(v *member.Member) (int, error) {
	var id int64
	_, id, err := m.Connector.GetOrm().Save(nil, v)
	if err != nil {
		return -1, err
	}
	v.Id = int(id)
	m.initMember(v)

	rc := core.GetRedisConn()
	defer rc.Close()
	rc.Do("RPUSH", variable.KvMemberUpdateQueue,
		fmt.Sprintf("%d-create", v.Id)) // push to queue

	// 更新会员数 todo: 考虑去掉
	var total = 0
	m.Connector.ExecScalar("SELECT COUNT(0) FROM mm_member", &total)
	gof.CurrentApp.Storage().Set(variable.KvTotalMembers, total)

	return v.Id, err
}

func (m *MemberRep) initMember(v *member.Member) {
	orm := m.Connector.GetOrm()
	orm.Save(nil, &member.Account{
		MemberId:         v.Id,
		Balance:          0,
		TotalConsumption: 0,
		TotalCharge:      0,
		TotalPay:         0,
		UpdateTime:       v.RegTime,
	})

	orm.Save(nil, &member.BankInfo{
		MemberId: v.Id,
		State:    1,
	})

	orm.Save(nil, &member.Relation{
		MemberId:           v.Id,
		CardId:             "",
		RefereesId:         0,
		RegisterMerchantId: 0,
	})
}

func (m *MemberRep) GetMemberIdByUser(user string) int {
	var id int
	m.Connector.ExecScalar("SELECT id FROM mm_member WHERE usr = ?", &id, user)
	return id
}

// 创建会员
func (m *MemberRep) CreateMember(v *member.Member) member.IMember {
	return memberImpl.NewMember(m.GetManager(), v, m,
		m._mssRep, m._valRep, m._partnerRep)
}

// 创建会员,仅作为某些操作使用,不保存
func (m *MemberRep) CreateMemberById(memberId int) member.IMember {
	return m.CreateMember(&member.Member{Id: memberId})
}

// 根据邀请码获取会员编号
func (m *MemberRep) GetMemberIdByInvitationCode(code string) int {
	var memberId int
	m.ExecScalar("SELECT id FROM mm_member WHERE invitation_code=?", &memberId, code)
	return memberId
}

func (m *MemberRep) GetLevel(merchantId, levelValue int) *merchant.MemberLevel {
	e := merchant.MemberLevel{}
	err := m.Connector.GetOrm().GetBy(&e, "merchant_id=? AND value = ?", merchantId, levelValue)
	if err != nil {
		return nil
	}
	return &e
}

// 获取下一个等级
func (m *MemberRep) GetNextLevel(merchantId, levelVal int) *merchant.MemberLevel {
	e := merchant.MemberLevel{}
	err := m.Connector.GetOrm().GetBy(&e, "merchant_id=? AND value>? LIMIT 0,1", merchantId, levelVal)
	if err != nil {
		return nil
	}
	return &e
}

// 获取会员等级
func (m *MemberRep) GetMemberLevels(merchantId int) []*merchant.MemberLevel {
	list := []*merchant.MemberLevel{}
	m.Connector.GetOrm().Select(&list,
		"merchant_id=?", merchantId)
	return list
}

// 删除会员等级
func (m *MemberRep) DeleteMemberLevel(merchantId, id int) error {
	_, err := m.Connector.GetOrm().Delete(&merchant.MemberLevel{},
		"id=? AND merchant_id=?", id, merchantId)
	return err
}

// 保存等级
func (m *MemberRep) SaveMemberLevel(merchantId int, v *merchant.MemberLevel) (int, error) {
	orm := m.Connector.GetOrm()
	var err error
	if v.Id > 0 {
		_, _, err = orm.Save(v.Id, v)
	} else {
		_, _, err = orm.Save(nil, v)
		m.Connector.ExecScalar(`SELECT MAX(id) FROM pt_member_level`, &v.Id)

	}
	return v.Id, err
}

// 获取会员最后更新时间
func (m *MemberRep) GetMemberLatestUpdateTime(id int) int64 {
	var updateTime int64
	m.Connector.ExecScalar(`SELECT update_time FROM mm_member where id=?`, &updateTime, id)
	return updateTime
}

// 获取账户
func (m *MemberRep) GetAccount(memberId int) *member.Account {
	e := &member.Account{}
	key := m.getAccountCk(memberId)
	if m.Storage.Get(key, &e) != nil {
		if m.Connector.GetOrm().Get(memberId, e) != nil {
			return nil
		}
		m.Storage.Set(key, *e)
	} else {
		//log.Println(fmt.Sprintf("--- account: %d > %#v",memberId,e))
	}
	return e
}

// 保存账户，传入会员编号
func (m *MemberRep) SaveAccount(v *member.Account) (int, error) {
	_, _, err := m.Connector.GetOrm().Save(v.MemberId, v)
	if err == nil {
		m.pushToAccountUpdateQueue(v.MemberId, v.UpdateTime)
		m.Storage.Set(m.getAccountCk(v.MemberId), *v)
	}
	return v.MemberId, err
}

func (m *MemberRep) pushToAccountUpdateQueue(memberId int, updateTime int64) {
	rc := core.GetRedisConn()
	defer rc.Close()
	// 保存最后更新时间
	mutKey := fmt.Sprintf("%s%d", variable.KvAccountUpdateTime, memberId)
	rc.Do("SETEX", mutKey, 3600*400, updateTime)
	// push to tcp notify queue
	rc.Do("RPUSH", variable.KvAccountUpdateTcpNotifyQueue, memberId)
}

// 获取银行信息
func (m *MemberRep) GetBankInfo(memberId int) *member.BankInfo {
	e := new(member.BankInfo)
	m.Connector.GetOrm().Get(memberId, e)
	return e
}

// 保存银行信息
func (m *MemberRep) SaveBankInfo(v *member.BankInfo) error {
	var err error
	_, _, err = m.Connector.GetOrm().Save(v.MemberId, v)
	return err
}

// 保存积分记录
func (m *MemberRep) SaveIntegralLog(l *member.IntegralLog) error {
	orm := m.Connector.GetOrm()
	var err error
	if l.Id > 0 {
		_, _, err = orm.Save(l.Id, l)
	} else {
		_, _, err = orm.Save(nil, l)
	}
	return err
}

func (m *MemberRep) getRelationCk(memberId int) string {
	return fmt.Sprintf("go2o:rep:mm-relation:%d", memberId)
}

// 获取会员关联
func (m *MemberRep) GetRelation(memberId int) *member.Relation {
	e := &member.Relation{}
	key := m.getRelationCk(memberId)
	if m.Storage.Get(key, e) != nil {
		if m.Connector.GetOrm().Get(memberId, e) != nil {
			return nil
		}
		m.Storage.Set(key, *e)
	}
	return e
}

// 获取积分对应的等级
func (m *MemberRep) GetLevelValueByExp(merchantId int, exp int) int {
	var levelId int
	m.Connector.ExecScalar(`SELECT lv.value FROM pt_member_level lv
	 	where lv.merchant_id=? AND lv.require_exp <= ? AND lv.enabled=1
	 	 ORDER BY lv.require_exp DESC LIMIT 0,1`,
		&levelId, merchantId, exp)
	return levelId

}

// 锁定会员
func (m *MemberRep) LockMember(id int, state int) error {
	_, err := m.Connector.ExecNonQuery("update mm_member set state=? WHERE id=?", state, id)
	return err
}

// 用户名是否存在
func (m *MemberRep) CheckUsrExist(usr string, memberId int) bool {
	var c int
	m.Connector.ExecScalar("SELECT COUNT(0) FROM mm_member WHERE usr=? AND id<>?",
		&c, usr, memberId)
	return c != 0
}

// 手机号码是否使用
func (m *MemberRep) CheckPhoneBind(phone string, memberId int) bool {
	var c int
	m.Connector.ExecScalar("SELECT COUNT(0) FROM mm_profile WHERE phone=? AND member_id<>?",
		&c, phone, memberId)
	return c != 0
}

// 保存绑定
func (m *MemberRep) SaveRelation(v *member.Relation) error {
	_, _, err := m.Connector.GetOrm().Save(v.MemberId, v)
	if err == nil {
		err = m.Storage.Set(m.getRelationCk(v.MemberId), *v)
	}
	return err
}

// 保存地址
func (m *MemberRep) SaveDeliver(v *member.DeliverAddress) (int, error) {
	orm := m.Connector.GetOrm()
	if v.Id <= 0 {
		_, _, err := orm.Save(nil, v)
		m.Connector.ExecScalar("SELECT MAX(id) FROM mm_delivery_addr", &v.Id)
		return v.Id, err
	} else {
		_, _, err := orm.Save(v.Id, v)
		return v.Id, err
	}
}

// 获取全部配送地址
func (m *MemberRep) GetDeliverAddress(memberId int) []*member.DeliverAddress {
	addresses := []*member.DeliverAddress{}
	m.Connector.GetOrm().Select(&addresses, "member_id=?", memberId)
	return addresses
}

// 获取配送地址
func (m *MemberRep) GetSingleDeliverAddress(memberId, deliverId int) *member.DeliverAddress {
	var address member.DeliverAddress
	err := m.Connector.GetOrm().Get(deliverId, &address)

	if err == nil && address.MemberId == memberId {
		return &address
	}
	return nil
}

// 删除配送地址
func (m *MemberRep) DeleteDeliver(memberId, deliverId int) error {
	_, err := m.Connector.ExecNonQuery(
		"DELETE FROM mm_deliver_addr WHERE member_id=? AND id=?",
		memberId, deliverId)
	return err
}

// 邀请
func (m *MemberRep) GetMyInvitationMembers(memberId, begin, end int) (
	total int, rows []*dto.InvitationMember) {
	arr := []*dto.InvitationMember{}
	m.Connector.ExecScalar(`SELECT COUNT(0) FROM mm_member WHERE id IN
	 (SELECT member_id FROM mm_relation WHERE invi_member_id=?)`, &total, memberId)
	if total > 0 {
		m.Connector.Query(`SELECT m.id,m.usr,m.level,p.avatar,p.name,p.phone,p.im FROM
            (SELECT id,usr,level FROM mm_member WHERE id IN (SELECT member_id FROM
             mm_relation WHERE invi_member_id=?) ORDER BY level DESC,id LIMIT ?,?) m
             INNER JOIN mm_profile p ON p.member_id = m.id ORDER BY level DESC,id`,
			func(rs *sql.Rows) {
				for rs.Next() {
					e := &dto.InvitationMember{}
					rs.Scan(&e.MemberId, &e.User, &e.Level, &e.Avatar, &e.NickName, &e.Phone, &e.Im)
					arr = append(arr, e)
				}
			}, memberId, begin, end-begin)
	}
	return total, arr
}

// 获取下级会员数量
func (m *MemberRep) GetSubInvitationNum(memberId int, memberIdArr []int) map[int]int {
	if len(memberIdArr) == 0 {
		return map[int]int{}
	}
	var ids []string = make([]string, len(memberIdArr))
	for i, v := range memberIdArr {
		ids[i] = strconv.Itoa(v)
	}
	memberIds := strings.Join(ids, ",")
	var d map[int]int = make(map[int]int)
	err := m.Connector.Query(fmt.Sprintf("SELECT r1.member_id,"+
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
func (m *MemberRep) GetInvitationMeMember(memberId int) *member.Member {
	var d *member.Member = new(member.Member)
	err := m.Connector.GetOrm().GetByQuery(d,
		"SELECT * FROM mm_member WHERE id =(SELECT invi_member_id FROM mm_relation  WHERE id=?)",
		memberId)

	if err != nil {
		return nil
	}
	return d
}

// 根据编号获取余额变动信息
func (m *MemberRep) GetBalanceInfo(id int) *member.BalanceInfo {
	var e member.BalanceInfo
	if err := m.Connector.GetOrm().Get(id, &e); err == nil {
		return &e
	}
	return nil
}

// 根据号码获取余额变动信息
func (m *MemberRep) GetBalanceInfoByNo(tradeNo string) *member.BalanceInfo {
	var e member.BalanceInfo
	if err := m.Connector.GetOrm().GetBy(&e, "trade_no=?", tradeNo); err == nil {
		return &e
	}
	return nil
}

// 保存余额变动信息
func (m *MemberRep) SaveBalanceInfo(v *member.BalanceInfo) (int, error) {
	var err error
	var orm = m.Connector.GetOrm()
	if v.Id > 0 {
		_, _, err = orm.Save(v.Id, v)
	} else {
		_, _, err = orm.Save(nil, v)
		m.Connector.ExecScalar("SELECT MAX(id) FROM mm_balance_info WHERE member_id=?", &v.Id, v.MemberId)
	}
	return v.Id, err
}

// 保存理财账户信息
func (m *MemberRep) SaveGrowAccount(memberId int, balance, totalAmount,
	growEarnings, totalGrowEarnings float32, updateTime int64) error {
	_, err := m.Connector.ExecNonQuery(`UPDATE mm_account SET grow_balance=?,
		grow_amount=?,grow_earnings=?,grow_total_earnings=?,update_time=? where member_id=?`,
		balance, totalAmount, growEarnings, totalGrowEarnings, updateTime, memberId)
	m.pushToAccountUpdateQueue(memberId, updateTime)
	return err
}

// 获取会员分页的优惠券列表
func (m *MemberRep) GetMemberPagedCoupon(memberId, start, end int, where string) (total int, rows []*dto.SimpleCoupon) {
	list := []*dto.SimpleCoupon{}
	m.Connector.ExecScalar(fmt.Sprintf(`SELECT COUNT(distinct pi.id)
        FROM pm_info pi INNER JOIN pm_coupon c ON c.id = pi.id
	    INNER JOIN pm_coupon_bind pb ON pb.coupon_id=pi.id
	    WHERE member_id=? AND %s`, where), &total, memberId)
	if total > 0 {
		m.Connector.GetOrm().SelectByQuery(&list,
			fmt.Sprintf(`SELECT pi.id,SUM(1) as num,pi.short_name as title,
            code,fee,c.discount,is_used,over_time FROM pm_info pi
             INNER JOIN pm_coupon c ON c.id = pi.id
	        INNER JOIN pm_coupon_bind pb ON pb.coupon_id=pi.id
	        WHERE member_id=? AND %s GROUP BY pi.id order by bind_time DESC LIMIT ?,?`, where),
			memberId, start, end-start)
	}
	return total, list
}
