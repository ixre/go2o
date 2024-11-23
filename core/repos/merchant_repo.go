/**
 * Copyright 2014 @ 56x.net.
 * name :
 * author : jarryliu
 * date : 2013-12-12 17:16
 * description :
 * history :
 */

package repos

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/ixre/go2o/core/domain/interface/approval"
	"github.com/ixre/go2o/core/domain/interface/invoice"
	"github.com/ixre/go2o/core/domain/interface/item"
	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/domain/interface/merchant"
	"github.com/ixre/go2o/core/domain/interface/merchant/shop"
	"github.com/ixre/go2o/core/domain/interface/merchant/staff"
	"github.com/ixre/go2o/core/domain/interface/merchant/user"
	"github.com/ixre/go2o/core/domain/interface/merchant/wholesaler"
	mss "github.com/ixre/go2o/core/domain/interface/message"
	rbac "github.com/ixre/go2o/core/domain/interface/rabc"
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/domain/interface/sys"
	"github.com/ixre/go2o/core/domain/interface/valueobject"
	"github.com/ixre/go2o/core/domain/interface/wallet"
	merchantImpl "github.com/ixre/go2o/core/domain/merchant"
	"github.com/ixre/go2o/core/infrastructure/fw"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/storage"
)

var _ merchant.IMerchantRepo = new(merchantRepo)
var mchMerchantDaoImplMapped = false

type merchantRepo struct {
	fw.BaseRepository[merchant.Merchant]
	authRepo      fw.BaseRepository[merchant.Authenticate]
	billRepo      fw.Repository[merchant.MerchantBill]
	_settleRepo   fw.Repository[merchant.SettleConf]
	Connector     db.Connector
	_orm          orm.Orm
	storage       storage.Interface
	manager       merchant.IMerchantManager
	_wsRepo       wholesaler.IWholesaleRepo
	_itemRepo     item.IItemRepo
	_userRepo     user.IUserRepo
	_employeeRepo staff.IStaffRepo
	_mssRepo      mss.IMessageRepo
	_shopRepo     shop.IShopRepo
	_valRepo      valueobject.IValueRepo
	_memberRepo   member.IMemberRepo
	_sysRepo      sys.ISystemRepo
	_walletRepo   wallet.IWalletRepo
	_registryRepo registry.IRegistryRepo
	_invoiceRepo  invoice.IInvoiceRepo
	_approvalRepo approval.IApprovalRepository
	_rbacRepo     rbac.IRbacRepository
	mux           *sync.RWMutex
}

// GetBalanceAccountLog implements merchant.IMerchantRepo

func NewMerchantRepo(o orm.Orm, on fw.ORM, storage storage.Interface,
	wsRepo wholesaler.IWholesaleRepo, itemRepo item.IItemRepo,
	shopRepo shop.IShopRepo, userRepo user.IUserRepo,
	employeeRepo staff.IStaffRepo,
	memberRepo member.IMemberRepo,
	sysRepo sys.ISystemRepo,
	mssRepo mss.IMessageRepo,
	walletRepo wallet.IWalletRepo,
	valRepo valueobject.IValueRepo,
	registryRepo registry.IRegistryRepo,
	invoiceRepo invoice.IInvoiceRepo,
	approvalRepo approval.IApprovalRepository,
	rbacRepo rbac.IRbacRepository,
) merchant.IMerchantRepo {
	if !mchMerchantDaoImplMapped {
		// 映射实体
		o.Mapping(merchant.Merchant{}, "mch_merchant")
		o.Mapping(merchant.Authenticate{}, "mch_authenticate")
		mchMerchantDaoImplMapped = true
	}
	r := &merchantRepo{
		Connector:     o.Connector(),
		_orm:          o,
		storage:       storage,
		_wsRepo:       wsRepo,
		_itemRepo:     itemRepo,
		_userRepo:     userRepo,
		_employeeRepo: employeeRepo,
		_mssRepo:      mssRepo,
		_shopRepo:     shopRepo,
		_sysRepo:      sysRepo,
		_valRepo:      valRepo,
		_memberRepo:   memberRepo,
		_walletRepo:   walletRepo,
		_registryRepo: registryRepo,
		_invoiceRepo:  invoiceRepo,
		_approvalRepo: approvalRepo,
		_rbacRepo:     rbacRepo,
		mux:           &sync.RWMutex{},
	}
	r.ORM = on
	r.authRepo.ORM = on
	r.billRepo = &fw.BaseRepository[merchant.MerchantBill]{ORM: on}
	r._settleRepo = &fw.BaseRepository[merchant.SettleConf]{ORM: on}
	return r
}

func (m *merchantRepo) BillRepo() fw.Repository[merchant.MerchantBill] {
	return m.billRepo
}

func (m *merchantRepo) SettleRepo() fw.Repository[merchant.SettleConf] {
	return m._settleRepo
}

// 获取商户管理器
func (m *merchantRepo) GetManager() merchant.IMerchantManager {
	if m.manager == nil {
		m.manager = merchantImpl.NewMerchantManager(m, m._valRepo)
	}
	return m.manager
}

func (m *merchantRepo) CreateMerchant(v *merchant.Merchant) merchant.IMerchantAggregateRoot {
	return merchantImpl.NewMerchant(v,
		m.storage,
		m, m._wsRepo,
		m._itemRepo,
		m._shopRepo,
		m._userRepo,
		m._employeeRepo,
		m._memberRepo,
		m._sysRepo,
		m._walletRepo,
		m._valRepo,
		m._registryRepo,
		m._invoiceRepo,
		m._approvalRepo,
		m._rbacRepo)
}

func (m *merchantRepo) cleanCache(mchId int64) {
	key := m.getMchCacheKey(mchId)
	m.storage.Delete(key)
	PrefixDel(m.storage, key+":*")
}

func (m *merchantRepo) getMchCacheKey(mchId int64) string {
	return fmt.Sprintf("go2o:repo:mch:%d", mchId)
}

func (m *merchantRepo) GetMerchant(id int) merchant.IMerchantAggregateRoot {
	e := merchant.Merchant{}
	key := m.getMchCacheKey(int64(id))
	err := m.storage.Get(key, &e)
	if err != nil {
		// 获取并缓存到列表中
		err = m._orm.Get(id, &e)
		if err != nil {
			return nil
		}
		m.storage.Set(key, e)
	}
	return m.CreateMerchant(&e)
}

// 根据登录用户名获取商户
func (m *merchantRepo) GetMerchantByUsername(user string) merchant.IMerchantAggregateRoot {
	e := merchant.Merchant{}
	if err := m._orm.GetBy(&e, "username=$1", user); err == nil {
		return m.CreateMerchant(&e)
	}
	return nil
}

// 获取账户
func (m *merchantRepo) GetAccount(mchId int) *merchant.Account {
	e := merchant.Account{}
	err := m._orm.Get(mchId, &e)
	if err == nil {
		return &e
	}
	// 初始化一个钱包账户
	if err == sql.ErrNoRows {
		e.MchId = mchId
		e.UpdateTime = int(time.Now().Unix())
		orm.Save(m._orm, &e, 0)
		return &e
	}
	return nil
}

// 获取合作商主要的域名主机
func (m *merchantRepo) GetMerchantMajorHost(mchId int) string {
	//todo:
	var host string
	m.Connector.ExecScalar(`SELECT host FROM pt_siteconf WHERE mch_id= $1 LIMIT 1`,
		&host, mchId)
	return host
}

// 验证商户用户名是否存在
func (m *merchantRepo) CheckUserExists(user string, id int) bool {
	var row int
	err := m.Connector.ExecScalar(`SELECT COUNT(*) FROM mch_merchant WHERE username= $1 AND id <> $2 LIMIT 1`,
		&row, user, id)
	if err != nil {
		panic(err)
	}
	return row > 0
}

// CheckMemberBind 验证会员是否绑定商户
func (m *merchantRepo) CheckMemberBind(memberId int, mchId int) bool {
	var row int
	m.Connector.ExecScalar(`SELECT COUNT(1) FROM mch_merchant
		WHERE member_id = $1 AND id <> $2`,
		&row, memberId, mchId)
	return row > 0
}

// 保存
func (m *merchantRepo) SaveMerchant(v *merchant.Merchant) (int, error) {
	id, err := orm.I64(orm.Save(m._orm, v, int(v.Id)))
	if err == nil {
		m.cleanCache(id)
	}
	return int(id), err
}

// 获取商户的编号
func (m *merchantRepo) GetMerchantsId() []int32 {
	dst := []int32{}
	var i int32

	m.Connector.Query("SELECT id FROM mch_merchant", func(rows *sql.Rows) {
		for rows.Next() {
			rows.Scan(&i)
			dst = append(dst, i)
		}
	})
	return dst
}

// 获取销售配置
func (m *merchantRepo) GetMerchantSaleConf(mchId int64) *merchant.SaleConf {
	//10%分成
	//0.2,         #上级
	//0.1,         #上上级
	//0.8          #消费者自己
	var saleConf = new(merchant.SaleConf)
	if m._orm.Get(mchId, saleConf) == nil {
		return saleConf
	}
	return nil
}

func (m *merchantRepo) SaveMerchantSaleConf(v *merchant.SaleConf) error {
	var err error
	if v.MchId > 0 {
		_, _, err = m._orm.Save(v.MchId, v)
	} else {
		_, _, err = m._orm.Save(nil, v)
	}
	return err
}

// 保存API信息
func (m *merchantRepo) SaveApiInfo(v *merchant.ApiInfo) (err error) {
	if m.GetApiInfo(v.MerchantId) == nil {
		_, err = orm.Save(m._orm, v, 0)
	} else {
		_, err = orm.Save(m._orm, v, int(v.MerchantId))
	}
	return err
}

// 获取API信息
func (m *merchantRepo) GetApiInfo(mchId int) *merchant.ApiInfo {
	var d = new(merchant.ApiInfo)
	if err := m._orm.Get(mchId, d); err == nil {
		return d
	}
	return nil
}

// 根据API编号获取商户编号
func (m *merchantRepo) GetMerchantIdByApiId(apiId string) int64 {
	var mchId int64
	m.Connector.ExecScalar("SELECT mch_id FROM mch_api_info WHERE api_id= $1", &mchId, apiId)
	return mchId
}

// 获取键值
func (m *merchantRepo) GetKeyValue(mchId int, indent string, k string) string {
	var v string
	m.Connector.ExecScalar(
		fmt.Sprintf("SELECT value FROM pt_%s WHERE merchant_id= $1 AND `key`= $2", indent),
		&v, mchId, k)
	return v
}

// 设置键值
func (m *merchantRepo) SaveKeyValue(mchId int, indent string, k, v string, updateTime int64) error {
	i, err := m.Connector.ExecNonQuery(
		fmt.Sprintf("UPDATE pt_%s SET value= $1,update_time= $2 WHERE merchant_id= $3 AND `key`= $4", indent),
		v, updateTime, mchId, k)
	if i == 0 {
		_, err = m.Connector.ExecNonQuery(
			fmt.Sprintf("INSERT INTO pt_%s(merchant_id,`key`,value,update_time)VALUES($1,$2,$3,$4)", indent),
			mchId, k, v, updateTime)
	}
	return err
}

// 获取多个键值
func (m *merchantRepo) GetKeyMap(mchId int, indent string, k []string) map[string]string {
	mp := make(map[string]string)
	var k1, v1 string
	m.Connector.Query(fmt.Sprintf("SELECT `key`,value FROM pt_%s WHERE merchant_id= $1 AND `key` IN ($2)", indent),
		func(rows *sql.Rows) {
			for rows.Next() {
				rows.Scan(&k1, &v1)
				mp[k1] = v1
			}
		}, mchId, strings.Join(k, ","))
	return mp
}

// 检查是否包含值的键数量,keyStr为键模糊匹配
func (m *merchantRepo) CheckKvContainValue(mchId int, indent string, value string, keyStr string) int {
	var i int
	err := m.Connector.ExecScalar("SELECT COUNT(1) FROM pt_"+indent+
		" WHERE merchant_id= $1 AND value= $2 AND `key` LIKE '%"+
		keyStr+"%'", &i, mchId, value)
	if err != nil {
		return 999
	}
	return i
}

// 根据关键字获取字典
func (m *merchantRepo) GetKeyMapByChar(mchId int, indent string, keyword string) map[string]string {
	mp := make(map[string]string)
	var k1, v1 string
	m.Connector.Query("SELECT `key`,value FROM pt_"+indent+
		" WHERE merchant_id= $1 AND `key` LIKE '%"+keyword+"%'",
		func(rows *sql.Rows) {
			for rows.Next() {
				rows.Scan(&k1, &v1)
				mp[k1] = v1
			}
		}, mchId)
	return mp
}

func (m *merchantRepo) GetLevel(mchId, levelValue int32) *merchant.MemberLevel {
	e := merchant.MemberLevel{}
	err := m._orm.GetBy(&e, "merchant_id= $1 AND value = $2", mchId, levelValue)
	if err != nil {
		return nil
	}
	return &e
}

// 获取下一个等级
func (m *merchantRepo) GetNextLevel(mchId, levelVal int32) *merchant.MemberLevel {
	e := merchant.MemberLevel{}
	err := m._orm.GetBy(&e, "merchant_id= $1 AND value> $2 LIMIT 1", mchId, levelVal)
	if err != nil {
		return nil
	}
	return &e
}

// 获取会员等级
func (m *merchantRepo) GetMemberLevels(mchId int64) []*merchant.MemberLevel {
	var list []*merchant.MemberLevel
	m._orm.Select(&list,
		"merchant_id= $1", mchId)
	return list
}

// 删除会员等级
func (m *merchantRepo) DeleteMemberLevel(mchId, id int32) error {
	_, err := m._orm.Delete(&merchant.MemberLevel{},
		"id= $1 AND merchant_id= $2", id, mchId)
	return err
}

// 保存等级
func (m *merchantRepo) SaveMemberLevel(mchId int64, v *merchant.MemberLevel) (int32, error) {
	return orm.I32(orm.Save(m._orm, v, int(v.Id)))
}

//	func (m *merchantRepo) UpdateMechOfflineRate(id int, rate float32, return_rate float32) error {
//		_, err := m.Connector.ExecNonQuery("UPDATE mch_merchant SET offline_rate= ? ,return_rate= ? WHERE  id= ?", rate, return_rate, id)
//		return err
//	}
//
//	func (m *merchantRepo) GetOfflineRate(id int32) (float32, float32, error) {
//		var rate float32
//		var return_rate float32
//		err := m.Connector.ExecScalar("SELECT  offline_rate FROM mch_merchant WHERE id= ?", &rate, id)
//		m.Connector.ExecScalar("SELECT  return_rate  FROM mch_merchant WHERE id= ?", &return_rate, id)
//		return rate, return_rate, err
//	}
//
// 保存会员账户
func (m *merchantRepo) SaveAccount(v *merchant.Account) (int, error) {
	orm := m._orm
	var err error
	if v.MchId > 0 {
		_, _, err = orm.Save(v.MchId, v)
	}
	return int(v.MchId), err
}

// Get MchBuyerGroupSetting
func (m *merchantRepo) GetMchBuyerGroupByGroupId(mchId, groupId int32) *merchant.MchBuyerGroupSetting {
	e := merchant.MchBuyerGroupSetting{}
	err := m._orm.GetBy(&e, "mch_id= $1 AND group_id= $2", mchId, groupId)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:MchBuyerGroupSetting")
	}
	return nil
}

// Select MchBuyerGroupSetting
func (m *merchantRepo) SelectMchBuyerGroup(mchId int64) []*merchant.MchBuyerGroupSetting {
	var list []*merchant.MchBuyerGroupSetting
	err := m._orm.Select(&list, "mch_id= $1", mchId)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:MchBuyerGroupSetting")
	}
	return list
}

// Save MchBuyerGroupSetting
func (m *merchantRepo) SaveMchBuyerGroup(v *merchant.MchBuyerGroupSetting) (int, error) {
	id, err := orm.Save(m._orm, v, int(v.ID))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:MchBuyerGroupSetting")
	}
	return id, err
}

// Get MchTradeConf
func (m *merchantRepo) GetMchTradeConf(primary interface{}) *merchant.TradeConf {
	e := merchant.TradeConf{}
	err := m._orm.Get(primary, &e)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:MchTradeConf")
	}
	return nil
}

// GetBy MchTradeConf
func (m *merchantRepo) GetMchTradeConfBy(where string, v ...interface{}) *merchant.TradeConf {
	e := merchant.TradeConf{}
	err := m._orm.GetBy(&e, where, v...)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:MchTradeConf")
	}
	return nil
}

// Select MchTradeConf
func (m *merchantRepo) SelectMchTradeConf(where string, v ...interface{}) []*merchant.TradeConf {
	var list []*merchant.TradeConf
	err := m._orm.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:MchTradeConf")
	}
	return list
}

// Save MchTradeConf
func (m *merchantRepo) SaveMchTradeConf(v *merchant.TradeConf) (int, error) {
	id, err := orm.Save(m._orm, v, int(v.ID))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:MchTradeConf")
	}
	return id, err
}

// Delete MchTradeConf
func (m *merchantRepo) DeleteMchTradeConf(primary interface{}) error {
	err := m._orm.DeleteByPk(merchant.TradeConf{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:MchTradeConf")
	}
	return err
}

// Batch Delete MchTradeConf
func (m *merchantRepo) BatchDeleteMchTradeConf(where string, v ...interface{}) (int64, error) {
	r, err := m._orm.Delete(merchant.TradeConf{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:MchTradeConf")
	}
	return r, err
}

func (m *merchantRepo) GetBalanceAccountLog(id int) *merchant.BalanceLog {
	e := merchant.BalanceLog{}
	err := m._orm.Get(id, &e)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:MchBalanceLog")
	}
	return nil
}

// SaveBalanceAccountLog implements merchant.IMerchantRepo
func (m *merchantRepo) SaveBalanceAccountLog(v *merchant.BalanceLog) (int, error) {
	id, err := orm.Save(m._orm, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:MchBalanceLog")
	}
	return id, err
}

// GetMerchantByMemberId implements merchant.IMerchantRepo
func (m *merchantRepo) GetMerchantByMemberId(memberId int) merchant.IMerchantAggregateRoot {
	v := merchant.Merchant{}
	err := m._orm.GetBy(&v, "member_id= $1", memberId)
	if err == nil {
		return m.CreateMerchant(&v)
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:MchMerchant")
	}
	return nil
}

// SaveAuthenticate Save 商户认证信息
func (m *merchantRepo) SaveAuthenticate(v *merchant.Authenticate) (int, error) {
	dst, err := m.authRepo.Save(v)
	return dst.Id, err
}

// GetMerchantAuthenticate implements merchant.IMerchantRepo.
func (m *merchantRepo) GetMerchantAuthenticate(mchId int, version int) *merchant.Authenticate {
	return m.authRepo.FindBy("mch_id = ? AND version= ?", mchId, version)
}

func (m *merchantRepo) DeleteOthersAuthenticate(mchId int, id int) error {
	_, err := m.authRepo.DeleteBy("mch_id = $1 AND id <> $2", mchId, id)
	return err
}
