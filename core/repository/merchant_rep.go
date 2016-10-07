/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-12 17:16
 * description :
 * history :
 */

package repository

import (
	"database/sql"
	"fmt"
	"github.com/jsix/gof/db"
	"github.com/jsix/gof/db/orm"
	"github.com/jsix/gof/storage"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/domain/interface/merchant/user"
	"go2o/core/domain/interface/mss"
	"go2o/core/domain/interface/valueobject"
	merchantImpl "go2o/core/domain/merchant"
	"go2o/core/infrastructure/domain"
	"strings"
	"sync"
	"time"
)

var _ merchant.IMerchantRep = new(merchantRep)

type merchantRep struct {
	db.Connector
	storage    storage.Interface
	manager    merchant.IMerchantManager
	_userRep   user.IUserRep
	_mssRep    mss.IMssRep
	_shopRep   shop.IShopRep
	_valRep    valueobject.IValueRep
	_memberRep member.IMemberRep
	mux        *sync.RWMutex
}

func NewMerchantRep(c db.Connector, storage storage.Interface, shopRep shop.IShopRep,
	userRep user.IUserRep, memberRep member.IMemberRep, mssRep mss.IMssRep,
	valRep valueobject.IValueRep) merchant.IMerchantRep {
	return &merchantRep{
		Connector:  c,
		storage:    storage,
		_userRep:   userRep,
		_mssRep:    mssRep,
		_shopRep:   shopRep,
		_valRep:    valRep,
		_memberRep: memberRep,
		mux:        &sync.RWMutex{},
	}
}

// 获取商户管理器
func (m *merchantRep) GetManager() merchant.IMerchantManager {
	if m.manager == nil {
		m.manager = merchantImpl.NewMerchantManager(m, m._valRep)
	}
	return m.manager
}

// 创建会员申请商户密钥
func (m *merchantRep) CreateSignUpToken(memberId int) string {
	mKey := fmt.Sprintf("go2o:rep:mch:signup:mm-%d", memberId)
	if token, err := m.storage.GetString(mKey); err == nil {
		return token
	}
	for {
		token := domain.NewSecret(0)[8:14]
		key := "go2o:rep:mch:signup:tk-" + token
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
func (m *merchantRep) GetMemberFromSignUpToken(token string) int {
	key := "go2o:rep:mch:signup:tk-" + token
	id, err := m.storage.GetInt(key)
	if err == nil {
		return id
	}
	return -1
}

func (m *merchantRep) CreateMerchant(v *merchant.Merchant) merchant.IMerchant {
	return merchantImpl.NewMerchant(v, m, m._shopRep, m._userRep,
		m._memberRep, m._valRep)
}

func (m *merchantRep) cleanCache(mchId int) {
	key := m.getMchCacheKey(mchId)
	m.storage.Del(key)
	PrefixDel(m.storage, key+":*")
}

func (m *merchantRep) getMchCacheKey(mchId int) string {
	return fmt.Sprintf("go2o:rep:mch:%d", mchId)
}

func (m *merchantRep) GetMerchant(id int) merchant.IMerchant {
	e := merchant.Merchant{}
	key := m.getMchCacheKey(id)
	if m.storage.Get(key, &e) != nil {
		// 获取并缓存到列表中
		err := m.Connector.GetOrm().Get(id, &e)
		if err == nil {
			m.storage.Set(key, e)
		}
	}
	return m.CreateMerchant(&e)
}

// 获取账户
func (m *merchantRep) GetAccount(mchId int) *merchant.Account {
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
func (m *merchantRep) GetMerchantMajorHost(merchantId int) string {
	//todo:
	var host string
	m.Connector.ExecScalar(`SELECT host FROM pt_siteconf WHERE merchant_id=? LIMIT 0,1`,
		&host, merchantId)
	return host
}

// 保存
func (m *merchantRep) SaveMerchant(v *merchant.Merchant) (int, error) {
	id, err := orm.Save(m.GetOrm(), v, v.Id)
	if err == nil {
		m.cleanCache(id)
	}
	return id, err
}

// 获取商户的编号
func (m *merchantRep) GetMerchantsId() []int {
	dst := []int{}
	var i int

	m.Connector.Query("SELECT id FROM mch_merchant", func(rows *sql.Rows) {
		for rows.Next() {
			rows.Scan(&i)
			dst = append(dst, i)
		}
		rows.Close()
	})
	return dst
}

// 获取销售配置
func (m *merchantRep) GetMerchantSaleConf(merchantId int) *merchant.SaleConf {
	//10%分成
	//0.2,         #上级
	//0.1,         #上上级
	//0.8          #消费者自己
	var saleConf *merchant.SaleConf = new(merchant.SaleConf)
	if m.Connector.GetOrm().Get(merchantId, saleConf) == nil {
		return saleConf
	}
	return nil
}

func (m *merchantRep) SaveMerchantSaleConf(v *merchant.SaleConf) error {
	var err error
	if v.MerchantId > 0 {
		_, _, err = m.Connector.GetOrm().Save(v.MerchantId, v)
	} else {
		_, _, err = m.Connector.GetOrm().Save(nil, v)
	}
	return err
}

// 保存API信息
func (m *merchantRep) SaveApiInfo(v *merchant.ApiInfo) error {
	orm := m.Connector.GetOrm()
	i, _, err := orm.Save(v.MerchantId, v)
	if i == 0 {
		_, _, err = orm.Save(nil, v)
	}
	return err
}

// 获取API信息
func (m *merchantRep) GetApiInfo(merchantId int) *merchant.ApiInfo {
	var d *merchant.ApiInfo = new(merchant.ApiInfo)
	if err := m.GetOrm().Get(merchantId, d); err == nil {
		return d
	}
	return nil
}

// 根据API编号获取商户编号
func (m *merchantRep) GetMerchantIdByApiId(apiId string) int {
	var merchantId int
	m.ExecScalar("SELECT merchant_id FROM mch_api_info WHERE api_id=?", &merchantId, apiId)
	return merchantId
}

// 获取键值
func (m *merchantRep) GetKeyValue(merchantId int, indent string, k string) string {
	var v string
	m.Connector.ExecScalar(
		fmt.Sprintf("SELECT value FROM pt_%s WHERE merchant_id=? AND `key`=?", indent),
		&v, merchantId, k)
	return v
}

// 设置键值
func (m *merchantRep) SaveKeyValue(merchantId int, indent string, k, v string, updateTime int64) error {
	i, err := m.Connector.ExecNonQuery(
		fmt.Sprintf("UPDATE pt_%s SET value=?,update_time=? WHERE merchant_id=? AND `key`=?", indent),
		v, updateTime, merchantId, k)
	if i == 0 {
		_, err = m.Connector.ExecNonQuery(
			fmt.Sprintf("INSERT INTO pt_%s(merchant_id,`key`,value,update_time)VALUES(?,?,?,?)", indent),
			merchantId, k, v, updateTime)
	}
	return err
}

// 获取多个键值
func (m *merchantRep) GetKeyMap(merchantId int, indent string, k []string) map[string]string {
	mp := make(map[string]string)
	var k1, v1 string
	m.Connector.Query(fmt.Sprintf("SELECT `key`,value FROM pt_%s WHERE merchant_id=? AND `key` IN (?)", indent),
		func(rows *sql.Rows) {
			for rows.Next() {
				rows.Scan(&k1, &v1)
				mp[k1] = v1
			}
		}, merchantId, strings.Join(k, ","))
	return mp
}

// 检查是否包含值的键数量,keyStr为键模糊匹配
func (m *merchantRep) CheckKvContainValue(merchantId int, indent string, value string, keyStr string) int {
	var i int
	err := m.Connector.ExecScalar("SELECT COUNT(0) FROM pt_"+indent+
		" WHERE merchant_id=? AND value=? AND `key` LIKE '%"+
		keyStr+"%'", &i, merchantId, value)
	if err != nil {
		return 999
	}
	return i
}

// 根据关键字获取字典
func (m *merchantRep) GetKeyMapByChar(merchantId int, indent string, keyword string) map[string]string {
	mp := make(map[string]string)
	var k1, v1 string
	m.Connector.Query("SELECT `key`,value FROM pt_"+indent+
		" WHERE merchant_id=? AND `key` LIKE '%"+keyword+"%'",
		func(rows *sql.Rows) {
			for rows.Next() {
				rows.Scan(&k1, &v1)
				mp[k1] = v1
			}
		}, merchantId)
	return mp
}

func (m *merchantRep) GetLevel(merchantId, levelValue int) *merchant.MemberLevel {
	e := merchant.MemberLevel{}
	err := m.Connector.GetOrm().GetBy(&e, "merchant_id=? AND value = ?", merchantId, levelValue)
	if err != nil {
		return nil
	}
	return &e
}

// 获取下一个等级
func (m *merchantRep) GetNextLevel(merchantId, levelVal int) *merchant.MemberLevel {
	e := merchant.MemberLevel{}
	err := m.Connector.GetOrm().GetBy(&e, "merchant_id=? AND value>? LIMIT 0,1", merchantId, levelVal)
	if err != nil {
		return nil
	}
	return &e
}

// 获取会员等级
func (m *merchantRep) GetMemberLevels(merchantId int) []*merchant.MemberLevel {
	list := []*merchant.MemberLevel{}
	m.Connector.GetOrm().Select(&list,
		"merchant_id=?", merchantId)
	return list
}

// 删除会员等级
func (m *merchantRep) DeleteMemberLevel(merchantId, id int) error {
	_, err := m.Connector.GetOrm().Delete(&merchant.MemberLevel{},
		"id=? AND merchant_id=?", id, merchantId)
	return err
}

// 保存等级
func (m *merchantRep) SaveMemberLevel(merchantId int, v *merchant.MemberLevel) (int, error) {
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

func (m *merchantRep) UpdateMechOfflineRate(id int, rate float32) error {
	_, err := m.Connector.ExecNonQuery("UPDATE mch_merchant SET offline_rate=?  WHERE  id=?", rate, id)
	return err
}

func (m *merchantRep) GetOfflineRate(id int) (float32, error) {
	var rate float32
	err := m.Connector.ExecScalar("SELECT  offline_rate  FROM mch_merchant WHERE id=?", &rate, id)
	return rate, err
}
