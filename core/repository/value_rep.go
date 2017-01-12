/**
 * Copyright 2015 @ z3q.net.
 * name : value_rep
 * author : jarryliu
 * date : 2016-05-27 15:32
 * description :
 * history :
 */
package repository

//todo: 因配置缓存与本地存储问题,子系统不能分布式部署。

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jsix/gof/db"
	"github.com/jsix/gof/db/orm"
	"github.com/jsix/gof/storage"
	"github.com/jsix/gof/util"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/infrastructure/tool/sms"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

var _ valueobject.IValueRepo = new(valueRepo)
var (
	valueRepCacheKey = "go2o:rep:value-rep:cache"
)

type valueRepo struct {
	db.Connector
	_orm   orm.Orm
	_kvMap map[string]int32
	_kvMux *sync.RWMutex

	storage          storage.Interface
	_wxConf          *valueobject.WxApiConfig
	_wxGob           *util.GobFile
	_rpConf          *valueobject.RegisterPerm
	_rpGob           *util.GobFile
	_numConf         *valueobject.GlobNumberConf
	_numGob          *util.GobFile
	_globMchConf     *valueobject.PlatformConf
	_mchGob          *util.GobFile
	_globRegistry    *valueobject.Registry
	_rstGob          *util.GobFile
	_globMchSaleConf *valueobject.GlobMchSaleConf
	_mscGob          *util.GobFile
	_smsConf         valueobject.SmsApiSet
	_smsGob          *util.GobFile
	_moAppConf       *valueobject.MoAppConf
	_moAppGob        *util.GobFile
	_tplConf         *valueobject.TemplateConf
	_tplGob          *util.GobFile
	_areaCache       map[int32][]*valueobject.Area
	_areaMux         sync.Mutex
}

func NewValueRepo(conn db.Connector, storage storage.Interface) valueobject.IValueRepo {
	return &valueRepo{
		Connector: conn,
		_orm:      conn.GetOrm(),
		storage:   storage,
		_kvMux:    &sync.RWMutex{},
		_rstGob:   util.NewGobFile("conf/core/registry"),
		_wxGob:    util.NewGobFile("conf/core/wx_api"),
		_rpGob:    util.NewGobFile("conf/core/register_perm"),
		_numGob:   util.NewGobFile("conf/core/number_conf"),
		_mchGob:   util.NewGobFile("conf/core/pm_conf"),
		_mscGob:   util.NewGobFile("conf/core/mch_sale_conf"),
		_smsGob:   util.NewGobFile("conf/core/sms_conf"),
		_tplGob:   util.NewGobFile("conf/core/tpl_conf"),
		_moAppGob: util.NewGobFile("conf/core/mo_app"),
	}
}

func (vp *valueRepo) checkReload() error {
	i, err := vp.storage.GetInt(valueRepCacheKey)
	if i == 0 || err != nil {
		vp._wxConf = nil
		vp._numConf = nil
		vp._rpConf = nil
		vp._smsConf = nil
		vp._globMchConf = nil
		vp._globMchSaleConf = nil
		vp._globRegistry = nil
	}
	return vp.storage.Set(valueRepCacheKey, 1)
}

func (vp *valueRepo) signReload() {
	vp.storage.Set(valueRepCacheKey, 0)
}

// 加载所有的键
func (s *valueRepo) loadAllKeys() {
	s._kvMux.Lock()
	s._kvMap = make(map[string]int32)
	list := s.selectSysKv("")
	for _, v := range list {
		s._kvMap[v.Key] = v.ID
		s.storage.Set("go2o:rep:kv:"+v.Key, v.Value)
	}
	s._kvMux.Unlock()
}

// 根据条件获取键值
func (s *valueRepo) selectSysKv(where string, v ...interface{}) []*valueobject.SysKeyValue {
	list := []*valueobject.SysKeyValue{}
	err := s._orm.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:SysKv")
	}
	return list
}

// 检查KEY与编号MAP
func (s *valueRepo) checkKvMap() {
	if s._kvMap == nil {
		s.loadAllKeys()
	}
}

// 根据键获取值
func (s *valueRepo) GetValue(key string) string {
	s.checkKvMap()
	s._kvMux.RLock()
	id, ok := s._kvMap[key]
	s._kvMux.RUnlock()
	if ok {
		rdsKey := "go2o:rep:kv:" + key
		r, err := s.storage.GetString(rdsKey)
		if err != nil {
			e := valueobject.SysKeyValue{}
			err := s._orm.Get(id, &e)
			if err != nil && err != sql.ErrNoRows {
				log.Println("[ Orm][ Error]:", err.Error(), "; Entity:SysKv")
			}
			r = e.Value
			if err == nil {
				s.storage.Set(rdsKey, r)
			}
		}
		return r
	}
	return ""
}

// 根据前缀获取值
func (s *valueRepo) GetValues(prefix string) map[string]string {
	s.checkKvMap()
	result := make(map[string]string)
	for k, _ := range s._kvMap {
		if strings.HasPrefix(k, prefix) {
			result[k] = s.GetValue(k)
		}
	}
	return result
}

// Save SysKv
func (s *valueRepo) SetValue(key string, v interface{}) error {
	s.checkKvMap()
	s._kvMux.RLock()
	id, ok := s._kvMap[key]
	s._kvMux.RUnlock()
	kv := &valueobject.SysKeyValue{
		ID:         id,
		Key:        key,
		Value:      util.Str(v),
		UpdateTime: time.Now().Unix(),
	}
	id2, err := orm.Save(s._orm, kv, int(kv.ID))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:SysKv")
	}
	if err == nil {
		id = int32(id2)
		s.storage.Set("go2o:rep:kv:"+kv.Key, kv.Value)
		if !ok {
			s._kvMux.Lock()
			s._kvMap[key] = id
			s._kvMux.Unlock()
		}
	}
	return err
}

// Delete SysKv
func (s *valueRepo) DeleteValue(key string) error {
	err := s._orm.DeleteByPk(valueobject.SysKeyValue{}, key)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:SysKv")
	}
	if err == nil {
		s._kvMux.Lock()
		delete(s._kvMap, key)
		s._kvMux.Unlock()
		s.storage.Del("go2o:rep:kv:" + key)
	}
	return err
}

// 获取微信接口配置
func (vp *valueRepo) GetWxApiConfig() valueobject.WxApiConfig {
	vp.checkReload()
	if vp._wxConf == nil {
		vp._wxConf = &valueobject.WxApiConfig{}
		vp._wxGob.Unmarshal(vp._wxConf)
	}
	return *vp._wxConf
}

// 保存微信接口配置
func (vp *valueRepo) SaveWxApiConfig(v *valueobject.WxApiConfig) error {
	if v != nil {
		defer vp.signReload()
		//todo: 检查证书文件是否存在
		vp._wxConf = v
		return vp._wxGob.Save(vp._wxConf)
	}
	return errors.New("nil value")
}

// 获取注册权限
func (vp *valueRepo) GetRegisterPerm() valueobject.RegisterPerm {
	vp.checkReload()
	if vp._rpConf == nil {
		v := defaultRegisterPerm
		vp._rpConf = &v
		vp._rpGob.Unmarshal(vp._rpConf)
	}
	return *vp._rpConf
}

// 保存注册权限
func (vp *valueRepo) SaveRegisterPerm(v *valueobject.RegisterPerm) error {
	if v != nil {
		defer vp.signReload()
		vp._rpConf = v
		return vp._rpGob.Save(vp._rpConf)
	}
	return nil
}

// 获取全局系统销售设置
func (vp *valueRepo) GetGlobNumberConf() valueobject.GlobNumberConf {
	vp.checkReload()
	if vp._numConf == nil {
		v := DefaultGlobNumberConf
		vp._numConf = &v
		vp._numGob.Unmarshal(vp._numConf)
	}
	return *vp._numConf
}

// 保存全局系统销售设置
func (vp *valueRepo) SaveGlobNumberConf(v *valueobject.GlobNumberConf) error {
	if v != nil {
		defer vp.signReload()
		vp._numConf = v
		return vp._numGob.Save(vp._numConf)
	}
	return nil
}

// 获取平台设置
func (vp *valueRepo) GetPlatformConf() valueobject.PlatformConf {
	vp.checkReload()
	if vp._globMchConf == nil {
		v := DefaultPlatformConf
		vp._globMchConf = &v
		vp._mchGob.Unmarshal(vp._globMchConf)
	}
	return *vp._globMchConf
}

// 保存平台设置
func (vp *valueRepo) SavePlatformConf(v *valueobject.PlatformConf) error {
	if v != nil {
		defer vp.signReload()
		vp._globMchConf = v
		return vp._mchGob.Save(vp._globMchConf)
	}
	return nil
}

// 获取模板配置
func (v *valueRepo) GetTemplateConf() valueobject.TemplateConf {
	v.checkReload()
	if v._tplConf == nil {
		v2 := DefaultTemplateConf
		v._tplConf = &v2
		v._tplGob.Unmarshal(v._tplConf)
	}
	return *v._tplConf
}

// 保存模板配置
func (v *valueRepo) SaveTemplateConf(t *valueobject.TemplateConf) error {
	if t != nil {
		defer v.signReload()
		v._tplConf = t
		return v._tplGob.Save(v._tplConf)
	}
	return nil
}

// 获取移动应用设置
func (v *valueRepo) GetMoAppConf() valueobject.MoAppConf {
	v.checkReload()
	if v._moAppConf == nil {
		v2 := DefaultMoAppConf
		v._moAppConf = &v2
		v._moAppGob.Unmarshal(v._moAppConf)
	}
	return *v._moAppConf
}

// 保存移动应用设置
func (v *valueRepo) SaveMoAppConf(r *valueobject.MoAppConf) error {
	if r != nil {
		defer v.signReload()
		v._moAppConf = r
		return v._moAppGob.Save(v._moAppConf)
	}
	return nil
}

// 获取数据存储
func (v *valueRepo) GetRegistry() valueobject.Registry {
	v.checkReload()
	if v._globRegistry == nil {
		v2 := DefaultRegistry
		v._globRegistry = &v2
		v._rstGob.Unmarshal(v._globRegistry)
	}
	return *v._globRegistry
}

// 保存数据存储
func (v *valueRepo) SaveRegistry(r *valueobject.Registry) error {
	if r != nil {
		defer v.signReload()
		v._globRegistry = r
		return v._rstGob.Save(v._globRegistry)
	}
	return nil
}

// 获取全局商户销售设置
func (vp *valueRepo) GetGlobMchSaleConf() valueobject.GlobMchSaleConf {
	vp.checkReload()
	if vp._globMchSaleConf == nil {
		v := DefaultGlobMchSaleConf
		vp._globMchSaleConf = &v
		vp._mscGob.Unmarshal(vp._globMchSaleConf)
	}
	return *vp._globMchSaleConf
}

// 保存全局商户销售设置
func (vp *valueRepo) SaveGlobMchSaleConf(v *valueobject.GlobMchSaleConf) error {
	if v != nil {
		defer vp.signReload()
		vp._globMchSaleConf = v
		return vp._mscGob.Save(vp._globMchSaleConf)
	}
	return nil
}

// 获取短信设置
func (vp *valueRepo) GetSmsApiSet() valueobject.SmsApiSet {
	vp.checkReload()
	if vp._smsConf == nil {
		vp._smsConf = defaultSmsConf
		vp._smsGob.Unmarshal(&vp._smsConf)
	}
	return vp._smsConf
}

// 保存短信API
func (vp *valueRepo) SaveSmsApiPerm(provider int, s *valueobject.SmsApiPerm) error {
	if _, ok := vp.GetSmsApiSet()[provider]; !ok {
		return errors.New("系统不支持的短信接口")
	}
	err := sms.CheckSmsApiPerm(provider, s)
	if err == nil {
		if s.Default {
			// 取消其他接口的默认选项
			for p, v := range vp._smsConf {
				if p == provider {
					v.Default = true
				} else {
					v.Default = false
				}
			}
		} else {
			//检验是否取消了正在使用的短信接口
			if i, _ := vp.GetDefaultSmsApiPerm(); i == provider {
				return errors.New("系统应启用一个短信接口")
			}
		}
		defer vp.signReload()
		vp._smsConf[provider] = s
		err = vp._smsGob.Save(vp._smsConf)
	}
	return err
}

// 获取默认的短信API
func (vp *valueRepo) GetDefaultSmsApiPerm() (int, *valueobject.SmsApiPerm) {
	for i, v := range vp.GetSmsApiSet() {
		if v.Default {
			return i, v
		}
	}
	panic(errors.New("至少为系统设置一个短信接口"))
}

// 获取下级区域
func (vp *valueRepo) GetChildAreas(id int32) []*valueobject.Area {
	vp._areaMux.Lock()
	defer vp._areaMux.Unlock()
	if vp._areaCache == nil {
		vp._areaCache = make(map[int32][]*valueobject.Area)
	}
	if v, ok := vp._areaCache[id]; ok {
		return v
	}
	v := []*valueobject.Area{}
	err := vp.Connector.GetOrm().Select(&v, "code <> 0 AND parent=?", id)
	if err == nil {
		vp._areaCache[id] = v
	}
	return v
}

// 获取地区名称
func (vp *valueRepo) GetAreaNames(id []int32) []string {
	strArr := make([]string, len(id))
	for i, v := range id {
		strArr[i] = strconv.Itoa(int(v))
	}
	i := 0
	vp.Connector.Query(fmt.Sprintf(
		"SELECT name FROM china_area WHERE code IN (%s)",
		strings.Join(strArr, ",")),
		func(rows *sql.Rows) {
			for rows.Next() {
				rows.Scan(&strArr[i])
				strArr[i] = strings.TrimSpace(strArr[i]) //去除空格
				i++
			}
		})
	return strArr
}

// 获取省市区字符串
func (vp *valueRepo) GetAreaString(province, city, district int32) string {
	names := vp.GetAreaNames([]int32{province, city, district})
	if names[1] == "市辖区" || names[1] == "市辖县" || names[1] == "县" {
		return strings.Join([]string{names[0], names[2]}, " ")
	}
	return strings.Join(names, " ")
}
