/**
 * Copyright 2014 @ 56x.net.
 * name :
 * author : jarryliu
 * date : 2013-12-09 10:13
 * description :
 * history :
 */

package repos

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/domain/interface/message"
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/domain/interface/valueobject"
	"github.com/ixre/go2o/core/domain/interface/wallet"
	memberImpl "github.com/ixre/go2o/core/domain/member"
	"github.com/ixre/go2o/core/dto"
	"github.com/ixre/go2o/core/infrastructure/format"
	"github.com/ixre/go2o/core/infrastructure/tool"
	"github.com/ixre/go2o/core/msq"
	"github.com/ixre/go2o/core/variable"
	"github.com/ixre/gof"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/storage"
	"github.com/ixre/gof/util"
)

var _ member.IMemberRepo = new(MemberRepoImpl)
var (
	memberManager member.IMemberManager
	memberMux     sync.Mutex
)

type MemberRepoImpl struct {
	storage storage.Interface
	db.Connector
	walletRepo   wallet.IWalletRepo
	valueRepo    valueobject.IValueRepo
	mssRepo      mss.IMssRepo
	registryRepo registry.IRegistryRepo
	o            orm.Orm
}

func NewMemberRepo(sto storage.Interface, o orm.Orm,
	walletRepo wallet.IWalletRepo, mssRepo mss.IMssRepo,
	valRepo valueobject.IValueRepo, registryRepo registry.IRegistryRepo) *MemberRepoImpl {
	return &MemberRepoImpl{
		storage:      sto,
		Connector:    o.Connector(),
		o:            o,
		mssRepo:      mssRepo,
		walletRepo:   walletRepo,
		valueRepo:    valRepo,
		registryRepo: registryRepo,
	}
}

// 获取管理服务
func (m *MemberRepoImpl) GetManager() member.IMemberManager {
	memberMux.Lock()
	if memberManager == nil {
		memberManager = memberImpl.NewMemberManager(m, m.valueRepo, m.registryRepo)
	}
	memberMux.Unlock()
	return memberManager
}

// 获取资料或初始化
func (m *MemberRepoImpl) GetProfile(memberId int64) *member.Profile {
	e := &member.Profile{}
	key := m.getProfileCk(memberId)
	if m.storage.Get(key, &e) != nil {
		if err := m.o.Get(memberId, e); err != nil {
			if err == sql.ErrNoRows {
				//todo: 没有资料应该到会员注册时候创建
				e.MemberId = memberId
				orm.Save(m.o, e, 0)
			}
		} else {
			m.storage.Set(key, *e)
		}
	}
	return e
}

// 保存资料
func (m *MemberRepoImpl) SaveProfile(v *member.Profile) error {
	_, _, err := m.o.Save(v.MemberId, v)
	if err == nil {
		err = m.storage.Set(m.getProfileCk(v.MemberId), *v)
	}
	return err
}

// 收藏,typeId 为类型编号, referId为关联的ID
func (m *MemberRepoImpl) Favorite(memberId int64, favType int, referId int64) error {
	_, _, err := m.o.Save(nil, &member.Favorite{
		MemberId:   memberId,
		FavType:    favType,
		ReferId:    referId,
		CreateTime: time.Now().Unix(),
	})
	return err
}

// 是否已收藏
func (m *MemberRepoImpl) Favored(memberId int64, favType int, referId int64) bool {
	num := 0
	m.Connector.ExecScalar(`SELECT COUNT(1) FROM mm_favorite
	WHERE member_id= $1 AND fav_type= $2 AND refer_id= $3`, &num,
		memberId, favType, referId)
	return num > 0
}

// 取消收藏
func (m *MemberRepoImpl) CancelFavorite(memberId int64, favType int, referId int64) error {
	_, err := m.o.Delete(&member.Favorite{},
		"member_id= $1 AND fav_type= $2 AND refer_id= $3",
		memberId, favType, referId)
	return err
}

var (
	globLevels []*member.Level
)

// 获取会员等级
func (m *MemberRepoImpl) GetMemberLevels_New() []*member.Level {
	const key = "go2o:repo:level:glob:cache"
	i, err := m.storage.GetInt(key)
	load := err != nil || i != 1 || globLevels == nil
	if load {
		list := make([]*member.Level, 0)
		m.o.Select(&list, "1=1 ORDER BY id ASC")
		globLevels = list
		m.storage.Set(key, 1)
	}
	return globLevels
}

// 获取等级对应的会员数
func (m *MemberRepoImpl) GetMemberNumByLevel_New(id int) int {
	total := 0
	m.Connector.ExecScalar("SELECT COUNT(1) FROM mm_member WHERE level= $1", &total, id)
	return total
}

// 删除会员等级
func (m *MemberRepoImpl) DeleteMemberLevel_New(id int) error {
	err := m.o.DeleteByPk(&member.Level{}, id)
	if err == nil {
		globLevels = nil
		PrefixDel(m.storage, "go2o:repo:level:*")
	}
	return err
}

// 保存会员等级
func (m *MemberRepoImpl) SaveMemberLevel_New(v *member.Level) (int, error) {
	id, err := orm.I32(orm.Save(m.o, v, int(v.Id)))
	if err == nil {
		globLevels = nil
		PrefixDel(m.storage, "go2o:repo:level:*")
	}
	return int(id), err
}

// 根据用户名获取会员
func (m *MemberRepoImpl) GetMemberByUser(user string) *member.Member {
	e := &member.Member{}
	err := m.o.GetBy(e, "username = $1", user)
	if err == nil {
		return e
	}
	return nil
}

func (m *MemberRepoImpl) getId(field string, value string) int64 {
	key := fmt.Sprintf("go2o:member:id:%s-%s", field, value)
	id, err := m.storage.GetInt(key)
	if err != nil {
		m.Connector.ExecScalar(fmt.Sprintf("SELECT id FROM mm_member WHERE %s=$1 LIMIT 1", field), &id, value)
		if id > 0 {
			m.storage.SetExpire(key, id, 48*3600)
		}
	}
	return int64(id)
}

// 根据编码获取会员
func (m *MemberRepoImpl) GetMemberIdByCode(code string) int64 {
	return m.getId("user_code", code)
}

// 根据手机号码获取会员
func (m *MemberRepoImpl) GetMemberValueByPhone(phone string) *member.Member {
	e := &member.Member{}
	err := m.o.GetBy(e, "phone = $1", phone)
	if err == nil {
		return e
	}
	return nil
}

// 根据手机号获取会员编号
func (m *MemberRepoImpl) GetMemberIdByPhone(phone string) int64 {
	return int64(m.getId("phone", phone))
}

// 根据邮箱地址获取会员编号
func (m *MemberRepoImpl) GetMemberIdByEmail(email string) int64 {
	return int64(m.getId("email", email))
}

func (m *MemberRepoImpl) getMemberCk(memberId int64) string {
	return fmt.Sprintf("go2o:repo:mm:inf:%d", memberId)
}
func (m *MemberRepoImpl) getAccountCk(memberId int64) string {
	return fmt.Sprintf("go2o:repo:mm:%d:acc", memberId)
}
func (m *MemberRepoImpl) getProfileCk(memberId int64) string {
	return fmt.Sprintf("go2o:repo:mm:pro:%d", memberId)
}
func (m *MemberRepoImpl) getTrustCk(memberId int64) string {
	return fmt.Sprintf("go2o:repo:mm:trust:%d", memberId)
}
func (m *MemberRepoImpl) getGlobLevelsCk() string {
	return "go2o:repo:mm-lv"
}

// 获取会员
func (m *MemberRepoImpl) GetMember(memberId int64) member.IMember {
	e := &member.Member{}
	key := m.getMemberCk(memberId)
	if err := m.storage.Get(key, &e); err != nil {
		//log.("-- mm",err)
		err = m.o.Get(memberId, e)
		if err != nil {
			return nil
		}
		m.storage.SetExpire(key, *e, DefaultCacheSeconds*3)
	}
	return m.CreateMember(e)
}

// 保存会员
func (m *MemberRepoImpl) SaveMember(v *member.Member) (int64, error) {
	if v.Id > 0 {
		//rc := core.GetRedisConn()
		//defer rc.Close()
		//// 保存最后更新时间
		//mutKey := fmt.Sprintf("%s%d", variable.KvMemberUpdateTime, v.Id)
		//rc.Do("SETEX", mutKey, 3600*400, v.UpdateTime)
		//rc.Do("RPUSH", variable.KvMemberUpdateTcpNotifyQueue, v.Id) // push to tcp notify queue

		// 保存会员信息
		_, _, err := m.o.Save(v.Id, v)

		if err == nil {
			// 存储到缓存中
			err = m.storage.Set(m.getMemberCk(v.Id), *v)
			// 存储到队列
			//rc.Do("RPUSH", variable.KvMemberUpdateQueue, fmt.Sprintf("%d-update", v.Id))

			// 推送消息
			go msq.Push(msq.MemberUpdated, "update|"+strconv.Itoa(int(v.Id)))
		}
		return v.Id, err
	}
	return m.createMember(v)
}

func (m *MemberRepoImpl) createMember(v *member.Member) (int64, error) {
	var id int64
	_, id, err := m.o.Save(nil, v)
	if err != nil {
		return -1, err
	}
	v.Id = id
	// 推送消息
	go msq.Push(msq.MemberUpdated, "create|"+strconv.Itoa(int(v.Id)))
	//rc := core.GetRedisConn()
	//defer rc.Close()
	//rc.Do("RPUSH", variable.KvMemberUpdateQueue,
	//	fmt.Sprintf("%d-create", v.Id)) // push to queue

	// 更新会员数 todo: 考虑去掉
	var total = 0
	m.Connector.ExecScalar("SELECT COUNT(1) FROM mm_member", &total)
	gof.CurrentApp.Storage().Set(variable.KvTotalMembers, total)

	return v.Id, err
}

// 删除会员
func (m *MemberRepoImpl) DeleteMember(id int64) error {
	m.storage.Delete(m.getMemberCk(id))
	_, err := m.ExecNonQuery("delete from mm_member where id = $1", id)
	sql := `
    /* 清理会员 */
     delete from mm_profile where member_id NOT IN (select id from mm_member) and member_id > 0;
     delete from mm_bank where member_id NOT IN(SELECT id FROM mm_member) and member_id > 0;
     delete from mm_account where member_id NOT IN(SELECT id FROM mm_member) and member_id > 0;
     delete from mm_relation where member_id NOT IN(SELECT id FROM mm_member) and member_id > 0;
     delete from mm_integral_log where member_id NOT IN (SELECT id FROM mm_member) and id > 0;
     delete from pay_order where buy_user NOT IN(SELECT id FROM mm_member) and id > 0;
     delete from mm_levelup where member_id NOT IN(SELECT id FROM mm_member) and id > 0;
    `
	for _, v := range strings.Split(sql, ";") {
		if v = strings.TrimSpace(v); len(v) > 5 {
			_, err := m.ExecNonQuery(v)
			if err != nil {
				log.Println("执行清理出错:", err, " sql:", v)
			}
		}
	}
	return err
}

func (m *MemberRepoImpl) GetMemberIdByUser(user string) int64 {
	return m.getId("username", user)
}

// 创建会员
func (m *MemberRepoImpl) CreateMember(v *member.Member) member.IMember {
	return memberImpl.NewMember(m.GetManager(), v, m,
		m.walletRepo, m.mssRepo, m.valueRepo, m.registryRepo)
}

// 创建会员,仅作为某些操作使用,不保存
func (m *MemberRepoImpl) CreateMemberById(memberId int64) member.IMember {
	if memberId <= 0 {
		return nil
	}
	return m.CreateMember(&member.Member{Id: memberId})
}

// 获取会员最后更新时间
func (m *MemberRepoImpl) GetMemberLatestUpdateTime(memberId int64) int64 {
	var updateTime int64
	m.Connector.ExecScalar(`SELECT update_time FROM mm_member where id= $1`,
		&updateTime, memberId)
	return updateTime
}

// GetAccount 获取账户
func (m *MemberRepoImpl) GetAccount(memberId int64) *member.Account {
	e := &member.Account{}
	key := m.getAccountCk(memberId)
	if m.storage.Get(key, &e) != nil {
		if m.o.Get(memberId, e) != nil {
			return nil
		}
		m.storage.Set(key, *e)
	}
	return e
}

// SaveAccount 保存账户，传入会员编号
func (m *MemberRepoImpl) SaveAccount(v *member.Account) (int64, error) {
	var err error
	if m.GetAccount(v.MemberId) == nil {
		_, _, err = m.o.Save(nil, v)
	} else {
		_, _, err = m.o.Save(v.MemberId, v)
		if err == nil {
			go m.pushToAccountUpdateQueue(v.MemberId, v.UpdateTime)
		}
	}
	if err == nil {
		m.storage.Set(m.getAccountCk(v.MemberId), *v)
	}
	return v.MemberId, err
}

func (m *MemberRepoImpl) pushToAccountUpdateQueue(memberId int64, updateTime int64) {
	//rc := core.GetRedisConn()
	//defer rc.Close()
	//// 保存最后更新时间
	//mutKey := fmt.Sprintf("%s%d", variable.KvAccountUpdateTime, memberId)
	//rc.Do("SETEX", mutKey, 3600*400, updateTime)
	//// push to tcp notify queue
	//rc.Do("RPUSH", variable.KvAccountUpdateTcpNotifyQueue, memberId)
}

// 获取银行信息
func (m *MemberRepoImpl) BankCards(memberId int64) []member.BankCard {
	var arr = make([]member.BankCard, 0)
	m.o.Select(&arr, "member_id=$1", memberId)
	return arr
}

// 保存银行信息
func (m *MemberRepoImpl) SaveBankCard(v *member.BankCard) error {
	_, err := orm.Save(m.o, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:MmBankCard")
	}
	return err
}

func (m *MemberRepoImpl) RemoveBankCard(memberId int64, cardNo string) error {
	_, err := m.o.Delete(&member.BankCard{},
		"member_id=$1 AND bank_account=$2", memberId, cardNo)
	return err
}

func (m *MemberRepoImpl) ReceiptsCodes(memberId int64) []member.ReceiptsCode {
	list := make([]member.ReceiptsCode, 0)
	err := m.o.Select(&list, "member_id=$1", memberId)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:MmReceiptsCode")
	}
	return list
}

func (m *MemberRepoImpl) SaveReceiptsCode(v *member.ReceiptsCode, memberId int64) (int, error) {
	id, err := orm.Save(m.o, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:MmReceiptsCode")
	}
	return id, err
}

// SaveIntegralLog 保存积分记录
func (m *MemberRepoImpl) SaveIntegralLog(v *member.IntegralLog) error {
	id, err := orm.Save(m.o, v, int(v.Id))
	if err == nil {
		v.Id = id
	}
	return err
}

// SaveBalanceLog 保存余额日志
func (m *MemberRepoImpl) SaveBalanceLog(v *member.BalanceLog) (int32, error) {
	id, err := orm.Save(m.o, v, int(v.Id))
	if err == nil {
		v.Id = int64(id)
	}
	return int32(id), err
}

// 保存钱包账户日志
func (m *MemberRepoImpl) SaveWalletAccountLog(v *member.WalletAccountLog) (int32, error) {
	return orm.I32(orm.Save(m.o, v, int(v.Id)))
}

func (m *MemberRepoImpl) GetWalletLog(id int32) *member.WalletAccountLog {
	e := member.WalletAccountLog{}
	if err := m.o.Get(id, &e); err != nil {
		return nil
	}
	return &e
}

// 获取会员提现次数键
func (m *MemberRepoImpl) getMemberTakeOutTimesKey(memberId int64) string {
	return fmt.Sprintf("sys:go2o:repo:mm:take-out-times:%d", memberId)
}

// 增加会员当天提现次数
func (m *MemberRepoImpl) AddTodayTakeOutTimes(memberId int64) error {
	times := m.GetTodayTakeOutTimes(memberId)
	key := m.getMemberTakeOutTimesKey(memberId)
	// 保存到当天结束
	t := time.Now()
	d := (24-t.Hour())*3600 + (60-t.Minute())*60 + (60 - t.Second())
	return m.storage.SetExpire(key, times+1, int64(d))
}

// 获取会员每日提现次数
func (m *MemberRepoImpl) GetTodayTakeOutTimes(memberId int64) int {
	key := m.getMemberTakeOutTimesKey(memberId)
	applyTimes, _ := m.storage.GetInt(key)
	return applyTimes

	total := 0
	b, e := tool.GetStartEndUnix(time.Now())
	err := m.ExecScalar(`SELECT COUNT(1) FROM mm_wallet_log WHERE
        member_id= $1 AND kind IN($2,$3) AND create_time BETWEEN $4 AND $5`, &total,
		memberId, wallet.KWithdrawToBankCard,
		wallet.KWithdrawToThirdPart, b, e)
	if err != nil {
		handleError(err)
	}
	return total
}

func (m *MemberRepoImpl) getRelationCk(memberId int64) string {
	return fmt.Sprintf("go2o:repo:mm:%d:rel", memberId)
}

// 获取会员关联
func (m *MemberRepoImpl) GetRelation(memberId int64) *member.InviteRelation {
	e := member.InviteRelation{}
	key := m.getRelationCk(memberId)
	if m.storage.Get(key, &e) != nil {
		if err := m.o.Get(memberId, &e); err != nil {
			return nil
		}
		m.storage.Set(key, e)
	}
	return &e
}

// 获取会员邀请的会员编号列表
func (m *MemberRepoImpl) GetInviteChildren(id int64) []int64 {
	arr := make([]int64, 0)
	var d int64
	m.Query("SELECT member_id FROM mm_relation WHERE inviter_id=$1", func(r *sql.Rows) {
		for r.Next() {
			r.Scan(&d)
			arr = append(arr, d)
		}
	}, id)
	return arr
}

// 获取积分对应的等级
func (m *MemberRepoImpl) GetLevelValueByExp(mchId int64, exp int64) int {
	var levelId int
	m.Connector.ExecScalar(`SELECT lv.value FROM pt_member_level lv
	 	where lv.merchant_id= $1 AND lv.require_exp <= $2 AND lv.enabled=1
	 	 ORDER BY lv.require_exp DESC LIMIT 1`,
		&levelId, mchId, exp)
	return levelId

}

// 用户名是否存在
func (m *MemberRepoImpl) CheckUserExist(user string, memberId int64) bool {
	var c int
	m.Connector.ExecScalar("SELECT id FROM mm_member WHERE username= $1 AND id <> $2 LIMIT 1",
		&c, user, memberId)
	return c != 0
}

// 手机号码是否使用
func (m *MemberRepoImpl) CheckPhoneBind(phone string, memberId int64) bool {
	var c int
	m.Connector.ExecScalar("SELECT COUNT(1) FROM mm_member WHERE phone= $1 AND id <> $2",
		&c, phone, memberId)
	return c != 0
}

// 手机号码是否使用
func (m *MemberRepoImpl) CheckNicknameIsUse(phone string, memberId int64) bool {
	var c int
	m.Connector.ExecScalar("SELECT COUNT(1) FROM mm_member WHERE nickname= $1 AND id <> $2",
		&c, phone, memberId)
	return c != 0
}

// 保存绑定
func (m *MemberRepoImpl) SaveRelation(v *member.InviteRelation) (err error) {
	rel := m.GetRelation(v.MemberId)
	if rel == nil {
		_, _, err = m.o.Save(nil, v)
	} else {
		_, _, err = m.o.Save(v.MemberId, v)
	}
	if err == nil {
		err = m.storage.Set(m.getRelationCk(v.MemberId), *v)
	}
	return err
}

// 获取会员升级记录
func (m *MemberRepoImpl) GetLevelUpLog(id int) *member.LevelUpLog {
	e := member.LevelUpLog{}
	if m.o.Get(id, &e) == nil {
		return &e
	}
	return nil
}

// SaveLevelUpLog 保存会员升级记录
func (m *MemberRepoImpl) SaveLevelUpLog(v *member.LevelUpLog) (int32, error) {
	return orm.I32(orm.Save(m.o, v, int(v.Id)))
}

// SaveDeliver 保存地址
func (m *MemberRepoImpl) SaveDeliver(v *member.ConsigneeAddress) (int64, error) {
	return orm.I64(orm.Save(m.o, v, int(v.Id)))
}

// GetDeliverAddress 获取全部配送地址
func (m *MemberRepoImpl) GetDeliverAddress(memberId int64) []*member.ConsigneeAddress {
	var addresses []*member.ConsigneeAddress
	m.o.Select(&addresses, "member_id= $1 ORDER BY id ASC", memberId)
	return addresses
}

// GetSingleDeliverAddress 获取配送地址
func (m *MemberRepoImpl) GetSingleDeliverAddress(memberId, deliverId int64) *member.ConsigneeAddress {
	var address member.ConsigneeAddress
	err := m.o.Get(deliverId, &address)

	if err == nil && address.MemberId == memberId {
		return &address
	}
	return nil
}

// DeleteAddress 删除配送地址
func (m *MemberRepoImpl) DeleteAddress(memberId, deliverId int64) error {
	_, err := m.Connector.ExecNonQuery(
		"DELETE FROM mm_deliver_addr WHERE member_id= $1 AND id= $2",
		memberId, deliverId)
	return err
}

// GetMyInvitationMembers 邀请
func (m *MemberRepoImpl) GetMyInvitationMembers(memberId int64, begin, end int) (
	total int, rows []*dto.InvitationMember) {
	var arr []*dto.InvitationMember
	m.Connector.ExecScalar(`SELECT COUNT(1) FROM mm_member WHERE id IN
	 (SELECT member_id FROM mm_relation WHERE inviter_id= $1)`, &total, memberId)
	if total > 0 {
		m.Connector.Query(`SELECT id,username,level,portrait,real_name,phone,reg_time FROM mm_member 
				WHERE id IN (SELECT member_id FROM
             mm_relation WHERE inviter_id= $1)
             ORDER BY level DESC,id LIMIT $3 OFFSET $2`,
			func(rs *sql.Rows) {
				for rs.Next() {
					e := &dto.InvitationMember{}
					rs.Scan(&e.MemberId, &e.Username, &e.Level, &e.Portrait, &e.Nickname, &e.Phone, &e.RegTime)
					arr = append(arr, e)
				}
			}, memberId, begin, end-begin)
	}
	return total, arr
}

// 获取下级会员数量
func (m *MemberRepoImpl) GetSubInvitationNum(memberId int64, memberIdArr []int32) map[int32]int {
	if len(memberIdArr) == 0 {
		return map[int32]int{}
	}
	memberIds := format.I32ArrStrJoin(memberIdArr)
	var d = make(map[int32]int)
	err := m.Connector.Query(fmt.Sprintf("SELECT r1.member_id,"+
		"(SELECT COUNT(1) FROM mm_relation r2 WHERE r2.inviter_id = r1.member_id)"+
		"as num FROM mm_relation r1 WHERE r1.member_id IN(%s)", memberIds),
		func(rows *sql.Rows) {
			var id int32
			var num int
			for rows.Next() {
				rows.Scan(&id, &num)
				d[id] = num
			}
		})
	handleError(err)
	return d
}

// 获取推荐我的人
func (m *MemberRepoImpl) GetInvitationMeMember(memberId int64) *member.Member {
	var d = new(member.Member)
	err := m.o.GetByQuery(d,
		"SELECT * FROM mm_member WHERE id = (SELECT inviter_id FROM mm_relation  WHERE id= $1)",
		memberId)

	if err != nil {
		return nil
	}
	return d
}

// 保存余额变动信息
func (m *MemberRepoImpl) SaveFlowAccountInfo(v *member.FlowAccountLog) (int32, error) {
	return orm.I32(orm.Save(m.o, v, int(v.Id)))
}

// 保存理财账户信息
func (m *MemberRepoImpl) SaveGrowAccount(memberId int64, balance, totalAmount,
	growEarnings, totalGrowEarnings float32, updateTime int64) error {
	_, err := m.Connector.ExecNonQuery(`UPDATE mm_account SET grow_balance= $1,
		grow_amount= $2,grow_earnings= $3,grow_total_earnings= $4,update_time= $5 where member_id= $6`,
		balance, totalAmount, growEarnings, totalGrowEarnings, updateTime, memberId)
	//清除缓存
	m.storage.Delete(m.getAccountCk(memberId))
	//加入通知队列
	m.pushToAccountUpdateQueue(memberId, updateTime)
	return err
}

// 获取会员分页的优惠券列表
func (m *MemberRepoImpl) GetMemberPagedCoupon(memberId int64, start, end int, where string) (total int, rows []*dto.SimpleCoupon) {
	list := []*dto.SimpleCoupon{}
	m.Connector.ExecScalar(fmt.Sprintf(`SELECT COUNT(distinct pi.id)
        FROM pm_info pi INNER JOIN pm_coupon c ON c.id = pi.id
	    INNER JOIN pm_coupon_bind pb ON pb.coupon_id=pi.id
	    WHERE member_id= $1 AND %s`, where), &total, memberId)
	if total > 0 {
		m.o.SelectByQuery(&list,
			fmt.Sprintf(`SELECT pi.id,SUM(1) as num,pi.short_name as title,
            code,fee,c.discount,is_used,over_time FROM pm_info pi
             INNER JOIN pm_coupon c ON c.id = pi.id
	        INNER JOIN pm_coupon_bind pb ON pb.coupon_id=pi.id
	        WHERE member_id= $1 AND %s GROUP BY pi.id order by bind_time DESC LIMIT $3 OFFSET $2`, where),
			memberId, start, end-start)
	}
	return total, list
}

// Select MmBuyerGroup
func (m *MemberRepoImpl) SelectMmBuyerGroup(where string, v ...interface{}) []*member.BuyerGroup {
	var list []*member.BuyerGroup
	err := m.o.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:MmBuyerGroup")
	}
	return list
}

// Save MmBuyerGroup
func (m *MemberRepoImpl) SaveMmBuyerGroup(v *member.BuyerGroup) (int, error) {
	id, err := orm.Save(m.o, v, int(v.ID))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:MmBuyerGroup")
	}
	return id, err
}

func (m *MemberRepoImpl) SaveLockHistory(v *member.MmLockHistory) (int, error) {
	id, err := orm.Save(m.o, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:MmLockHistory")
	}
	return id, err
}

func (m *MemberRepoImpl) SaveLockInfo(v *member.MmLockInfo) (int, error) {
	id, err := orm.Save(m.o, v, v.Id)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:MmLockInfo")
	}
	return id, err
}

// 存储自动解锁信息到任务队列
func (m *MemberRepoImpl) RegisterUnlockJob(v *member.MmLockInfo) {
	slice := util.GetMinuteSlice(time.Unix(v.UnlockTime, 0), 1)
	key := fmt.Sprintf("%s:%s:%d", variable.KvMemberAutoUnlock, slice, v.MemberId)
	m.storage.SetExpire(key, v.MemberId, v.UnlockTime-v.LockTime+120)
}

func (m *MemberRepoImpl) DeleteLockInfos(memberId int64) error {
	_, err := m.o.Delete(member.MmLockInfo{}, "member_id=$1", memberId)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:MmLockInfo")
	}
	return err
}
