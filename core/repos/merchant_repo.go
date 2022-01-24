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
	"github.com/ixre/go2o/core/domain/interface/item"
	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/domain/interface/merchant"
	"github.com/ixre/go2o/core/domain/interface/merchant/shop"
	"github.com/ixre/go2o/core/domain/interface/merchant/user"
	"github.com/ixre/go2o/core/domain/interface/merchant/wholesaler"
	"github.com/ixre/go2o/core/domain/interface/mss"
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/domain/interface/valueobject"
	"github.com/ixre/go2o/core/domain/interface/wallet"
	merchantImpl "github.com/ixre/go2o/core/domain/merchant"
	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/storage"
	"log"
	"strings"
	"sync"
	"time"
)

var _ merchant.IMerchantRepo = new(merchantRepo)

type merchantRepo struct {
	db.Connector
	o             orm.Orm
	storage       storage.Interface
	manager       merchant.IMerchantManager
	_wsRepo       wholesaler.IWholesaleRepo
	_itemRepo     item.IItemRepo
	_userRepo     user.IUserRepo
	_mssRepo      mss.IMssRepo
	_shopRepo     shop.IShopRepo
	_valRepo      valueobject.IValueRepo
	_memberRepo   member.IMemberRepo
	_walletRepo   wallet.IWalletRepo
	_registryRepo registry.IRegistryRepo
	mux           *sync.RWMutex
}

func NewMerchantRepo(o orm.Orm, storage storage.Interface,
	wsRepo wholesaler.IWholesaleRepo, itemRepo item.IItemRepo,
	shopRepo shop.IShopRepo, userRepo user.IUserRepo, memberRepo member.IMemberRepo, mssRepo mss.IMssRepo,
	walletRepo wallet.IWalletRepo, valRepo valueobject.IValueRepo, registryRepo registry.IRegistryRepo) merchant.IMerchantRepo {
	return &merchantRepo{
		Connector:     o.Connector(),
		o:             o,
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

func (m *merchantRepo) cleanCache(mchId int64) {
	key := m.getMchCacheKey(mchId)
	m.storage.Delete(key)
	PrefixDel(m.storage, key+":*")
}

func (m *merchantRepo) getMchCacheKey(mchId int64) string {
	return fmt.Sprintf("go2o:repo:mch:%d", mchId)
}

func (m *merchantRepo) GetMerchant(id int) merchant.IMerchant {
	e := merchant.Merchant{}
	key := m.getMchCacheKey(int64(id))
	err := m.storage.Get(key, &e)
	if err != nil {
		// 获取并缓存到列表中
		err = m.o.Get(id, &e)
		if err != nil {
			return nil
		}
		m.storage.Set(key, e)
	}
	return m.CreateMerchant(&e)
}

// 根据登录用户名获取商户
func (m *merchantRepo) GetMerchantByLoginUser(user string) merchant.IMerchant {
	e := merchant.Merchant{}
	if err := m.o.GetBy(&e, "login_user=$1", user); err == nil {
		return m.CreateMerchant(&e)
	}
	return nil
}

// 获取账户
func (m *merchantRepo) GetAccount(mchId int) *merchant.Account {
	e := merchant.Account{}
	err := m.o.Get(mchId, &e)
	if err == nil {
		return &e
	}
	if err == sql.ErrNoRows {
		e.MchId = int64(mchId)
		e.UpdateTime = time.Now().Unix()
		orm.Save(m.o, &e, 0)
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
	m.Connector.ExecScalar(`SELECT COUNT(0) FROM mch_merchant WHERE login_user= $1 AND id <> $2 LIMIT 1`,
		&row, user, id)
	return row > 0
}

// CheckMemberBind 验证会员是否绑定商户
func (m *merchantRepo) CheckMemberBind(memberId int64, mchId int64) bool {
	var row int
	m.Connector.ExecScalar(`SELECT COUNT(0) FROM mch_merchant
		WHERE member_id = $1 AND id <> $2`,
		&row, memberId, mchId)
	return row > 0
}

// 保存
func (m *merchantRepo) SaveMerchant(v *merchant.Merchant) (int, error) {
	id, err := orm.I64(orm.Save(m.o, v, int(v.Id)))
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
	if m.o.Get(mchId, saleConf) == nil {
		return saleConf
	}
	return nil
}

func (m *merchantRepo) SaveMerchantSaleConf(v *merchant.SaleConf) error {
	var err error
	if v.MerchantId > 0 {
		_, _, err = m.o.Save(v.MerchantId, v)
	} else {
		_, _, err = m.o.Save(nil, v)
	}
	return err
}

// 保存API信息
func (m *merchantRepo) SaveApiInfo(v *merchant.ApiInfo) (err error) {
	if m.GetApiInfo(v.MerchantId) == nil {
		_, err = orm.Save(m.o, v, 0)
	} else {
		_, err = orm.Save(m.o, v, int(v.MerchantId))
	}
	return err
}

// 获取API信息
func (m *merchantRepo) GetApiInfo(mchId int) *merchant.ApiInfo {
	var d = new(merchant.ApiInfo)
	if err := m.o.Get(mchId, d); err == nil {
		return d
	}
	return nil
}

// 根据API编号获取商户编号
func (m *merchantRepo) GetMerchantIdByApiId(apiId string) int64 {
	var mchId int64
	m.ExecScalar("SELECT mch_id FROM mch_api_info WHERE api_id= $1", &mchId, apiId)
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
	err := m.Connector.ExecScalar("SELECT COUNT(0) FROM pt_"+indent+
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
	err := m.o.GetBy(&e, "merchant_id= $1 AND value = $2", mchId, levelValue)
	if err != nil {
		return nil
	}
	return &e
}

// 获取下一个等级
func (m *merchantRepo) GetNextLevel(mchId, levelVal int32) *merchant.MemberLevel {
	e := merchant.MemberLevel{}
	err := m.o.GetBy(&e, "merchant_id= $1 AND value> $2 LIMIT 1", mchId, levelVal)
	if err != nil {
		return nil
	}
	return &e
}

// 获取会员等级
func (m *merchantRepo) GetMemberLevels(mchId int64) []*merchant.MemberLevel {
	var list []*merchant.MemberLevel
	m.o.Select(&list,
		"merchant_id= $1", mchId)
	return list
}

// 删除会员等级
func (m *merchantRepo) DeleteMemberLevel(mchId, id int32) error {
	_, err := m.o.Delete(&merchant.MemberLevel{},
		"id= $1 AND merchant_id= $2", id, mchId)
	return err
}

// 保存等级
func (m *merchantRepo) SaveMemberLevel(mchId int64, v *merchant.MemberLevel) (int32, error) {
	return orm.I32(orm.Save(m.o, v, int(v.Id)))
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
	orm := m.o
	var err error
	if v.MchId > 0 {
		_, _, err = orm.Save(v.MchId, v)
	}
	return err
}

// Get MchEnterpriseInfo
func (m *merchantRepo) GetMchEnterpriseInfo(mchId int) *merchant.EnterpriseInfo {
	e := merchant.EnterpriseInfo{}
	err := m.o.GetBy(&e, "mch_id= $1", mchId)
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
	id, err := orm.Save(m.o, v, int(v.ID))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:MchEnterpriseInfo")
	}
	return id, err
}

// Get MchBuyerGroupSetting
func (m *merchantRepo) GetMchBuyerGroupByGroupId(mchId, groupId int32) *merchant.MchBuyerGroupSetting {
	e := merchant.MchBuyerGroupSetting{}
	err := m.o.GetBy(&e, "mch_id= $1 AND group_id= $2", mchId, groupId)
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
	err := m.o.Select(&list, "mch_id= $1", mchId)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:MchBuyerGroupSetting")
	}
	return list
}

// Save MchBuyerGroupSetting
func (m *merchantRepo) SaveMchBuyerGroup(v *merchant.MchBuyerGroupSetting) (int, error) {
	id, err := orm.Save(m.o, v, int(v.ID))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:MchBuyerGroupSetting")
	}
	return id, err
}

// Get MchTradeConf
func (m *merchantRepo) GetMchTradeConf(primary interface{}) *merchant.TradeConf {
	e := merchant.TradeConf{}
	err := m.o.Get(primary, &e)
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
	err := m.o.GetBy(&e, where, v...)
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
	err := m.o.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:MchTradeConf")
	}
	return list
}

// Save MchTradeConf
func (m *merchantRepo) SaveMchTradeConf(v *merchant.TradeConf) (int, error) {
	id, err := orm.Save(m.o, v, int(v.ID))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:MchTradeConf")
	}
	return id, err
}

// Delete MchTradeConf
func (m *merchantRepo) DeleteMchTradeConf(primary interface{}) error {
	err := m.o.DeleteByPk(merchant.TradeConf{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:MchTradeConf")
	}
	return err
}

// Batch Delete MchTradeConf
func (m *merchantRepo) BatchDeleteMchTradeConf(where string, v ...interface{}) (int64, error) {
	r, err := m.o.Delete(merchant.TradeConf{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:MchTradeConf")
	}
	return r, err
}
