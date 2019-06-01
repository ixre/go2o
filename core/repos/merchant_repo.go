/**
 * Copyright 2014 @ z3q.net.
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
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/storage"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/domain/interface/merchant/user"
	"go2o/core/domain/interface/merchant/wholesaler"
	"go2o/core/domain/interface/mss"
	"go2o/core/domain/interface/registry"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/domain/interface/wallet"
	merchantImpl "go2o/core/domain/merchant"
	"go2o/core/infrastructure/domain"
	"log"
	"strings"
	"sync"
	"time"
)

var _ merchant.IMerchantRepo = new(merchantRepo)

type merchantRepo struct {
	db.Connector
	_orm          orm.Orm
	storage       storage.Interface
	manager       merchant.IMerchantManager
	_wsRepo       wholesaler.IWholesaleRepo
	_itemRepo     item.IGoodsItemRepo
	_userRepo     user.IUserRepo
	_mssRepo      mss.IMssRepo
	_shopRepo     shop.IShopRepo
	_valRepo      valueobject.IValueRepo
	_memberRepo   member.IMemberRepo
	_walletRepo   wallet.IWalletRepo
	_registryRepo registry.IRegistryRepo
	mux           *sync.RWMutex
}

func NewMerchantRepo(c db.Connector, storage storage.Interface,
	wsRepo wholesaler.IWholesaleRepo, itemRepo item.IGoodsItemRepo,
	shopRepo shop.IShopRepo, userRepo user.IUserRepo, memberRepo member.IMemberRepo, mssRepo mss.IMssRepo,
	walletRepo wallet.IWalletRepo, valRepo valueobject.IValueRepo, registryRepo registry.IRegistryRepo) merchant.IMerchantRepo {
	return &merchantRepo{
		Connector:     c,
		_orm:          c.GetOrm(),
		storage:       storage,
		_wsRepo:       wsRepo,
		_itemRepo:     itemRepo,
		_userRepo:     userRepo,
		_mssRepo:      mssRepo,
		_shopRepo:     shopRepo,
		_valRepo:      valRepo,
		_memberRepo:   memberRepo,
		_walletRepo:   walletRepo,
		_registryRepo: registryRepo,
		mux:           &sync.RWMutex{},
	}
}

// 获取商户管理器
func (m *merchantRepo) GetManager() merchant.IMerchantManager {
	if m.manager == nil {
		m.manager = merchantImpl.NewMerchantManager(m, m._valRepo)
	}
	return m.manager
}

// 创建会员申请商户密钥
func (m *merchantRepo) CreateSignUpToken(memberId int64) string {
	mKey := fmt.Sprintf("go2o:repo:mch:signup:mm-%d", memberId)
	if token, err := m.storage.GetString(mKey); err == nil {
		return token
	}
	for {
		token := domain.NewSecret(0)[8:14]
		key := "go2o:repo:mch:signup:tk-" + token
		if _, err := m.storage.GetInt(key); err != nil {
			seconds := int64(time.Hour * 12)
			m.storage.SetExpire(key, memberId, seconds)
			m.storage.SetExpire(mKey, token, seconds)
			return token
		}
	}
	return ""
}

// 根据商户申请密钥获取会员编号
func (m *merchantRepo) GetMemberFromSignUpToken(token string) int64 {
	key := "go2o:repo:mch:signup:tk-" + token
	id, err := m.storage.GetInt64(key)
	if err == nil {
		return id
	}
	return -1
}

func (m *merchantRepo) CreateMerchant(v *merchant.Merchant) merchant.IMerchant {
	return merchantImpl.NewMerchant(v, m, m._wsRepo, m._itemRepo,
		m._shopRepo, m._userRepo, m._memberRepo, m._walletRepo, m._valRepo, m._registryRepo)
}

func (m *merchantRepo) cleanCache(mchId int32) {
	key := m.getMchCacheKey(mchId)
	m.storage.Del(key)
	PrefixDel(m.storage, key+":*")
}

func (m *merchantRepo) getMchCacheKey(mchId int32) string {
	return fmt.Sprintf("go2o:repo:mch:%d", mchId)
}

func (m *merchantRepo) GetMerchant(id int32) merchant.IMerchant {
	e := merchant.Merchant{}
	key := m.getMchCacheKey(id)
	err := m.storage.Get(key, &e)
	if err != nil {
		// 获取并缓存到列表中
		err = m.Connector.GetOrm().Get(id, &e)
		if err != nil {
			return nil
		}
		m.storage.Set(key, e)
	}
	return m.CreateMerchant(&e)
}

// 获取账户
func (m *merchantRepo) GetAccount(mchId int32) *merchant.Account {
	e := merchant.Account{}
	err := m.Connector.GetOrm().Get(mchId, &e)
	if err == nil {
		return &e
	}
	if err == sql.ErrNoRows {
		e.MchId = mchId
		e.UpdateTime = time.Now().Unix()
		orm.Save(m.Connector.GetOrm(), &e, 0)
		return &e
	}
	return nil
}

// 获取合作商主要的域名主机
func (m *merchantRepo) GetMerchantMajorHost(mchId int32) string {
	//todo:
	var host string
	m.Connector.ExecScalar(`SELECT host FROM pt_siteconf WHERE mch_id= $1 LIMIT 1`,
		&host, mchId)
	return host
}

// 保存
func (m *merchantRepo) SaveMerchant(v *merchant.Merchant) (int32, error) {
	id, err := orm.I32(orm.Save(m.GetOrm(), v, int(v.ID)))
	if err == nil {
		m.cleanCache(id)
	}
	return id, err
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
func (m *merchantRepo) GetMerchantSaleConf(mchId int32) *merchant.SaleConf {
	//10%分成
	//0.2,         #上级
	//0.1,         #上上级
	//0.8          #消费者自己
	var saleConf *merchant.SaleConf = new(merchant.SaleConf)
	if m.Connector.GetOrm().Get(mchId, saleConf) == nil {
		return saleConf
	}
	return nil
}

func (m *merchantRepo) SaveMerchantSaleConf(v *merchant.SaleConf) error {
	var err error
	if v.MerchantId > 0 {
		_, _, err = m.Connector.GetOrm().Save(v.MerchantId, v)
	} else {
		_, _, err = m.Connector.GetOrm().Save(nil, v)
	}
	return err
}

// 保存API信息
func (m *merchantRepo) SaveApiInfo(v *merchant.ApiInfo) error {
	_, err := orm.Save(m.GetOrm(), v, int(v.MerchantId))
	return err
}

// 获取API信息
func (m *merchantRepo) GetApiInfo(mchId int32) *merchant.ApiInfo {
	var d *merchant.ApiInfo = new(merchant.ApiInfo)
	if err := m.GetOrm().Get(mchId, d); err == nil {
		return d
	}
	return nil
}

// 根据API编号获取商户编号
func (m *merchantRepo) GetMerchantIdByApiId(apiId string) int32 {
	var mchId int32
	m.ExecScalar("SELECT mch_id FROM mch_api_info WHERE api_id= $1", &mchId, apiId)
	return mchId
}

// 获取键值
func (m *merchantRepo) GetKeyValue(mchId int32, indent string, k string) string {
	var v string
	m.Connector.ExecScalar(
		fmt.Sprintf("SELECT value FROM pt_%s WHERE merchant_id= $1 AND `key`= $2", indent),
		&v, mchId, k)
	return v
}

// 设置键值
func (m *merchantRepo) SaveKeyValue(mchId int32, indent string, k, v string, updateTime int64) error {
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
func (m *merchantRepo) GetKeyMap(mchId int32, indent string, k []string) map[string]string {
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
func (m *merchantRepo) CheckKvContainValue(mchId int32, indent string, value string, keyStr string) int {
	var i int
	err := m.Connector.ExecScalar("SELECT COUNT(0) FROM pt_"+indent+
		" WHERE merchant_id= $1 AND value= $2 AND `key` LIKE '%"+
		keyStr+"%'", &i, mchId, value)
	if err != nil {
		return 999
	}
	return i
}

// 根据关键字获取字典
func (m *merchantRepo) GetKeyMapByChar(mchId int32, indent string, keyword string) map[string]string {
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
	err := m.Connector.GetOrm().GetBy(&e, "merchant_id= $1 AND value = $2", mchId, levelValue)
	if err != nil {
		return nil
	}
	return &e
}

// 获取下一个等级
func (m *merchantRepo) GetNextLevel(mchId, levelVal int32) *merchant.MemberLevel {
	e := merchant.MemberLevel{}
	err := m.Connector.GetOrm().GetBy(&e, "merchant_id= $1 AND value> $2 LIMIT 1", mchId, levelVal)
	if err != nil {
		return nil
	}
	return &e
}

// 获取会员等级
func (m *merchantRepo) GetMemberLevels(mchId int32) []*merchant.MemberLevel {
	var list []*merchant.MemberLevel
	m.Connector.GetOrm().Select(&list,
		"merchant_id= $1", mchId)
	return list
}

// 删除会员等级
func (m *merchantRepo) DeleteMemberLevel(mchId, id int32) error {
	_, err := m.Connector.GetOrm().Delete(&merchant.MemberLevel{},
		"id= $1 AND merchant_id= $2", id, mchId)
	return err
}

// 保存等级
func (m *merchantRepo) SaveMemberLevel(mchId int32, v *merchant.MemberLevel) (int32, error) {
	return orm.I32(orm.Save(m.GetOrm(), v, int(v.Id)))
}

//
//func (m *merchantRepo) UpdateMechOfflineRate(id int, rate float32, return_rate float32) error {
//	_, err := m.Connector.ExecNonQuery("UPDATE mch_merchant SET offline_rate= ? ,return_rate= ? WHERE  id= ?", rate, return_rate, id)
//	return err
//}
//
//func (m *merchantRepo) GetOfflineRate(id int32) (float32, float32, error) {
//	var rate float32
//	var return_rate float32
//	err := m.Connector.ExecScalar("SELECT  offline_rate FROM mch_merchant WHERE id= ?", &rate, id)
//	m.Connector.ExecScalar("SELECT  return_rate  FROM mch_merchant WHERE id= ?", &return_rate, id)
//	return rate, return_rate, err
//}
//
// 保存销售配置
func (m *merchantRepo) UpdateAccount(v *merchant.Account) error {
	orm := m.Connector.GetOrm()
	var err error
	if v.MchId > 0 {
		_, _, err = orm.Save(v.MchId, v)
	}
	return err
}

// Get MchEnterpriseInfo
func (m *merchantRepo) GetMchEnterpriseInfo(mchId int32) *merchant.EnterpriseInfo {
	e := merchant.EnterpriseInfo{}
	err := m._orm.GetBy(&e, "mch_id= $1", mchId)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:MchEnterpriseInfo")
	}
	return nil
}

// Save MchEnterpriseInfo
func (m *merchantRepo) SaveMchEnterpriseInfo(v *merchant.EnterpriseInfo) (int, error) {
	id, err := orm.Save(m._orm, v, int(v.ID))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:MchEnterpriseInfo")
	}
	return id, err
}

// Get MchBuyerGroup
func (m *merchantRepo) GetMchBuyerGroupByGroupId(mchId, groupId int32) *merchant.MchBuyerGroup {
	e := merchant.MchBuyerGroup{}
	err := m._orm.GetBy(&e, "mch_id= $1 AND group_id= $2", mchId, groupId)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:MchBuyerGroup")
	}
	return nil
}

// Select MchBuyerGroup
func (m *merchantRepo) SelectMchBuyerGroup(mchId int32) []*merchant.MchBuyerGroup {
	var list []*merchant.MchBuyerGroup
	err := m._orm.Select(&list, "mch_id= $1", mchId)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:MchBuyerGroup")
	}
	return list
}

// Save MchBuyerGroup
func (m *merchantRepo) SaveMchBuyerGroup(v *merchant.MchBuyerGroup) (int, error) {
	id, err := orm.Save(m._orm, v, int(v.ID))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:MchBuyerGroup")
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
