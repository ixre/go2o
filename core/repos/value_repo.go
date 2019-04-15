/**
 * Copyright 2015 @ z3q.net.
 * name : value_rep
 * author : jarryliu
 * date : 2016-05-27 15:32
 * description :
 * history :
 */
package repos

//todo: 因配置缓存与本地存储问题,子系统不能分布式部署。

import (
	"database/sql"
	"errors"
	"github.com/ixre/gof"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/storage"
	"github.com/ixre/gof/util"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/infrastructure/tool/sms"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var _ valueobject.IValueRepo = new(valueRepo)
var (
	valueRepCacheKey = "go2o:repo:value-rep:cache"
)

type valueRepo struct {
	db.Connector
	o     orm.Orm
	kvMap map[string]int32
	kvMux *sync.RWMutex

	storage         storage.Interface
	wxConf          *valueobject.WxApiConfig
	wxGob           *util.GobFile
	rpConf          *valueobject.RegisterPerm
	rpGob           *util.GobFile
	numConf         *valueobject.GlobNumberConf
	numGob          *util.GobFile
	globMchConf     *valueobject.PlatformConf
	mchGob          *util.GobFile
	globRegistry    *valueobject.Registry
	rstGob          *util.GobFile
	globMchSaleConf *valueobject.GlobMchSaleConf
	mscGob          *util.GobFile
	smsConf         valueobject.SmsApiSet
	smsGob          *util.GobFile
	moAppConf       *valueobject.MoAppConf
	moAppGob        *util.GobFile
	tplConf         *valueobject.TemplateConf
	tplGob          *util.GobFile
	areaCache       map[int32][]*valueobject.Area
	areaMux         sync.Mutex

	confRegistry *gof.Registry
}

func NewValueRepo(confPath string, conn db.Connector, storage storage.Interface) valueobject.IValueRepo {
	confRegistry, err := gof.NewRegistry(confPath, ":")
	if err != nil {
		log.Println("[ Go2o][ Crash]: can't load registry,", err.Error())
		os.Exit(1)
	}
	return &valueRepo{
		Connector:    conn,
		o:            conn.GetOrm(),
		storage:      storage,
		kvMux:        &sync.RWMutex{},
		rstGob:       util.NewGobFile("conf/core/registry"),
		wxGob:        util.NewGobFile("conf/core/wx_api"),
		rpGob:        util.NewGobFile("conf/core/register_perm"),
		numGob:       util.NewGobFile("conf/core/number_conf"),
		mchGob:       util.NewGobFile("conf/core/pm_conf"),
		mscGob:       util.NewGobFile("conf/core/mch_sale_conf"),
		smsGob:       util.NewGobFile("conf/core/sms_conf"),
		tplGob:       util.NewGobFile("conf/core/tpl_conf"),
		moAppGob:     util.NewGobFile("conf/core/mo_app"),
		confRegistry: confRegistry,
	}
}

func (r *valueRepo) checkReload() error {
	i, err := r.storage.GetInt(valueRepCacheKey)
	if i == 0 || err != nil {
		r.wxConf = nil
		r.numConf = nil
		r.rpConf = nil
		r.smsConf = nil
		r.globMchConf = nil
		r.globMchSaleConf = nil
		r.globRegistry = nil
	}
	return r.storage.Set(valueRepCacheKey, 1)
}

func (r *valueRepo) signReload() {
	r.storage.Set(valueRepCacheKey, 0)
}

// 加载所有的键
func (r *valueRepo) loadAllKeys() {
	r.kvMux.Lock()
	r.kvMap = make(map[string]int32)
	list := r.selectSysKv("")
	for _, v := range list {
		r.kvMap[v.Key] = v.ID
		r.storage.Set("go2o:repo:kv:"+v.Key, v.Value)
	}
	r.kvMux.Unlock()
}

// 根据条件获取键值
func (r *valueRepo) selectSysKv(where string, v ...interface{}) []*valueobject.SysKeyValue {
	list := make([]*valueobject.SysKeyValue, 0)
	err := r.o.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:SysKv")
	}
	return list
}

// 检查KEY与编号MAP
func (r *valueRepo) checkKvMap() {
	if r.kvMap == nil {
		r.loadAllKeys()
	}
}

// 根据键获取值
func (r *valueRepo) GetValue(key string) string {
	r.checkKvMap()
	r.kvMux.RLock()
	id, ok := r.kvMap[key]
	r.kvMux.RUnlock()
	if ok {
		rdsKey := "go2o:repo:kv:" + key
		val, err := r.storage.GetString(rdsKey)
		if err != nil {
			e := valueobject.SysKeyValue{}
			err := r.o.Get(id, &e)
			if err != nil && err != sql.ErrNoRows {
				log.Println("[ Orm][ Error]:", err.Error(), "; Entity:SysKv")
			}
			val = e.Value
			if err == nil {
				r.storage.Set(rdsKey, val)
			}
		}
		return val
	}
	return ""
}

// 根据前缀获取值
func (r *valueRepo) GetValues(prefix string) map[string]string {
	r.checkKvMap()
	result := make(map[string]string)
	for k := range r.kvMap {
		if strings.HasPrefix(k, prefix) {
			result[k] = r.GetValue(k)
		}
	}
	return result
}

// Save SysKv
func (r *valueRepo) SetValue(key string, v interface{}) error {
	r.checkKvMap()
	r.kvMux.RLock()
	id, ok := r.kvMap[key]
	r.kvMux.RUnlock()
	kv := &valueobject.SysKeyValue{
		ID:         id,
		Key:        key,
		Value:      util.Str(v),
		UpdateTime: time.Now().Unix(),
	}
	id2, err := orm.Save(r.o, kv, int(kv.ID))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:SysKv")
	}
	if err == nil {
		id = int32(id2)
		r.storage.Set("go2o:repo:kv:"+kv.Key, kv.Value)
		if !ok {
			r.kvMux.Lock()
			r.kvMap[key] = id
			r.kvMux.Unlock()
		}
	}
	return err
}

// Delete SysKv
func (r *valueRepo) DeleteValue(key string) error {
	err := r.o.DeleteByPk(valueobject.SysKeyValue{}, key)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:SysKv")
	}
	if err == nil {
		r.kvMux.Lock()
		delete(r.kvMap, key)
		r.kvMux.Unlock()
		r.storage.Del("go2o:repo:kv:" + key)
	}
	return err
}

// 获取微信接口配置
func (r *valueRepo) GetWxApiConfig() valueobject.WxApiConfig {
	r.checkReload()
	if r.wxConf == nil {
		r.wxConf = &valueobject.WxApiConfig{}
		r.wxGob.Unmarshal(r.wxConf)
	}
	return *r.wxConf
}

// 保存微信接口配置
func (r *valueRepo) SaveWxApiConfig(v *valueobject.WxApiConfig) error {
	if v != nil {
		defer r.signReload()
		//todo: 检查证书文件是否存在
		r.wxConf = v
		return r.wxGob.Save(r.wxConf)
	}
	return errors.New("nil value")
}

// 获取注册权限
func (r *valueRepo) GetRegisterPerm() valueobject.RegisterPerm {
	r.checkReload()
	if r.rpConf == nil {
		v := defaultRegisterPerm
		r.rpConf = &v
		r.rpGob.Unmarshal(r.rpConf)
	}
	return *r.rpConf
}

// 保存注册权限
func (r *valueRepo) SaveRegisterPerm(v *valueobject.RegisterPerm) error {
	if v != nil {
		defer r.signReload()
		// 如果要验证手机，则必须开启填写手机
		if v.MustBindPhone {
			v.NeedPhone = true
		}
		r.rpConf = v
		return r.rpGob.Save(r.rpConf)
	}
	return nil
}

// 获取全局系统销售设置
func (r *valueRepo) GetGlobNumberConf() valueobject.GlobNumberConf {
	r.checkReload()
	if r.numConf == nil {
		v := DefaultGlobNumberConf
		r.numConf = &v
		r.numGob.Unmarshal(r.numConf)
	}
	return *r.numConf
}

// 保存全局系统销售设置
func (r *valueRepo) SaveGlobNumberConf(v *valueobject.GlobNumberConf) error {
	if v != nil {
		defer r.signReload()
		r.numConf = v
		return r.numGob.Save(r.numConf)
	}
	return nil
}

// 获取平台设置
func (r *valueRepo) GetPlatformConf() valueobject.PlatformConf {
	r.checkReload()
	if r.globMchConf == nil {
		v := DefaultPlatformConf
		r.globMchConf = &v
		r.mchGob.Unmarshal(r.globMchConf)
	}
	return *r.globMchConf
}

// 保存平台设置
func (r *valueRepo) SavePlatformConf(v *valueobject.PlatformConf) error {
	if v != nil {
		defer r.signReload()
		r.globMchConf = v
		return r.mchGob.Save(r.globMchConf)
	}
	return nil
}

// 获取模板配置
func (r *valueRepo) GetTemplateConf() valueobject.TemplateConf {
	r.checkReload()
	if r.tplConf == nil {
		v2 := DefaultTemplateConf
		r.tplConf = &v2
		r.tplGob.Unmarshal(r.tplConf)
	}
	return *r.tplConf
}

// 保存模板配置
func (r *valueRepo) SaveTemplateConf(t *valueobject.TemplateConf) error {
	if t != nil {
		defer r.signReload()
		r.tplConf = t
		return r.tplGob.Save(r.tplConf)
	}
	return nil
}

// 获取移动应用设置
func (r *valueRepo) GetMoAppConf() valueobject.MoAppConf {
	r.checkReload()
	if r.moAppConf == nil {
		v2 := DefaultMoAppConf
		r.moAppConf = &v2
		r.moAppGob.Unmarshal(r.moAppConf)
	}
	return *r.moAppConf
}

// 保存移动应用设置
func (r *valueRepo) SaveMoAppConf(v *valueobject.MoAppConf) error {
	if r != nil {
		defer r.signReload()
		r.moAppConf = v
		return r.moAppGob.Save(r.moAppConf)
	}
	return nil
}

// 获取数据存储
func (r *valueRepo) GetRegistry() valueobject.Registry {
	v := r.getRegistry()
	return *v
}

// 保存数据存储
func (r *valueRepo) SaveRegistry(v *valueobject.Registry) error {
	if r != nil {
		defer r.signReload()
		r.globRegistry = v
		return r.rstGob.Save(r.globRegistry)
	}
	return nil
}

// 获取数据存储
func (r *valueRepo) getRegistry() *valueobject.Registry {
	r.checkReload()
	if r.globRegistry == nil {
		v2 := DefaultRegistry
		r.globRegistry = &v2
		r.rstGob.Unmarshal(r.globRegistry)
	}
	return r.globRegistry
}

func (r *valueRepo) GetsRegistry(keys []string) []string {
	if strings.Index(keys[0], ":") != -1 {
		return r.getsRegistryNew(keys)
	}
	v := r.getRegistry()
	mp := make([]string, len(keys))
	for i, key := range keys {
		d, ok := v.RegistryData[key]
		if ok {
			mp[i] = d
		} else {
			mp[i] = "no value in registry"
		}
	}
	return mp
}
func (r *valueRepo) getsRegistryNew(keys []string) []string {
	mp := make([]string, len(keys))
	for i, k := range keys {
		v := r.confRegistry.Get(k)
		mp[i] = util.Str(v)
	}
	return mp
}
func (r *valueRepo) getsRegistryMapNew(keys []string) map[string]string {
	mp := map[string]string{}
	for _, k := range keys {
		v := r.confRegistry.Get(k)
		mp[k] = util.Str(v)
	}
	return mp
}

// 根据键获取数据值
func (r *valueRepo) GetsRegistryMap(keys []string) map[string]string {
	if strings.Index(keys[0], ":") != -1 {
		return r.getsRegistryMapNew(keys)
	}
	v := r.getRegistry()
	mp := map[string]string{}
	for _, key := range keys {
		d, ok := v.RegistryData[key]
		if ok {
			mp[key] = d
		} else {
			mp[key] = "no value in registry"
		}
	}
	return mp
}

// 保存数据值
func (r *valueRepo) SavesRegistry(values map[string]string) error {
	v := r.getRegistry()
	if v != nil {
		defer r.signReload()
		for k, val := range values {
			v.RegistryData[k] = val
		}
		r.globRegistry = v
		return r.rstGob.Save(r.globRegistry)
	}
	return nil
}

// 获取全局商户销售设置
func (r *valueRepo) GetGlobMchSaleConf() valueobject.GlobMchSaleConf {
	r.checkReload()
	if r.globMchSaleConf == nil {
		v := DefaultGlobMchSaleConf
		r.globMchSaleConf = &v
		r.mscGob.Unmarshal(r.globMchSaleConf)
	}
	return *r.globMchSaleConf
}

// 保存全局商户销售设置
func (r *valueRepo) SaveGlobMchSaleConf(v *valueobject.GlobMchSaleConf) error {
	if v != nil {
		defer r.signReload()
		r.globMchSaleConf = v
		return r.mscGob.Save(r.globMchSaleConf)
	}
	return nil
}

// 获取短信设置
func (r *valueRepo) GetSmsApiSet() valueobject.SmsApiSet {
	r.checkReload()
	if r.smsConf == nil {
		r.smsConf = defaultSmsConf
		r.smsGob.Unmarshal(&r.smsConf)
	}
	return r.smsConf
}

// 保存短信API
func (r *valueRepo) SaveSmsApiPerm(provider int, v *valueobject.SmsApiPerm) error {
	if _, ok := r.GetSmsApiSet()[provider]; !ok {
		return errors.New("系统不支持的短信接口")
	}
	err := sms.CheckSmsApiPerm(provider, v)
	if err == nil {
		if v.Default {
			// 取消其他接口的默认选项
			for p, c := range r.smsConf {
				if p == provider {
					c.Default = true
				} else {
					c.Default = false
				}
			}
		} else {
			//检验是否取消了正在使用的短信接口
			if i, _ := r.GetDefaultSmsApiPerm(); i == provider {
				return errors.New("系统应启用一个短信接口")
			}
		}
		defer r.signReload()
		r.smsConf[provider] = v
		err = r.smsGob.Save(r.smsConf)
	}
	return err
}

// 获取默认的短信API
func (r *valueRepo) GetDefaultSmsApiPerm() (int, *valueobject.SmsApiPerm) {
	for i, v := range r.GetSmsApiSet() {
		if v.Default {
			return i, v
		}
	}
	panic(errors.New("至少为系统设置一个短信接口"))
}

// 获取下级区域
func (r *valueRepo) GetChildAreas(code int32) []*valueobject.Area {
	r.areaMux.Lock()
	defer r.areaMux.Unlock()
	if r.areaCache == nil {
		r.areaCache = make(map[int32][]*valueobject.Area)
	}
	if v, ok := r.areaCache[code]; ok {
		return v
	}
	var v []*valueobject.Area
	err := r.Connector.GetOrm().Select(&v, "code <> 0 AND parent=?", code)
	if err == nil {
		r.areaCache[code] = v
	}
	return v
}

// 获取区域名称
func (r *valueRepo) GetAreaName(code int32) string {
	if code <= 0 {
		return ""
	}
	strId := strconv.Itoa(int(code))
	key := "go2o:repo:area:name-" + strId
	name, err := r.storage.GetString(key)
	if err != nil {
		err = r.Connector.ExecScalar("SELECT name FROM china_area WHERE code=?", &name, strId)
		if err == nil {
			name = strings.TrimSpace(name)
			if name == "市辖区" || name == "市辖县" || name == "县" {
				name = ""
			}
			r.storage.Set(key, name)
		}
	}
	return name
}

// 获取地区名称
func (r *valueRepo) GetAreaNames(codeArr []int32) []string {
	arr := make([]string, len(codeArr))
	for i, v := range codeArr {
		arr[i] = r.GetAreaName(v)
	}
	if len(codeArr) >= 3 {
		if arr[1] == "市辖区" || arr[1] == "市辖县" || arr[1] == "县" {
			return []string{arr[0], arr[2]}
		}
	}
	return arr
}

// 获取省市区字符串
func (r *valueRepo) GetAreaString(province, city, district int32) string {
	names := r.GetAreaNames([]int32{province, city, district})
	return strings.Join(names, " ")
}

// 获取省市区字符串
func (r *valueRepo) AreaString(province, city, district int32, detail string) string {
	names := r.GetAreaNames([]int32{province, city, district})
	prefix := []byte(strings.Join(names, ""))
	if len(prefix) != 0 && len(detail) != 0 {
		i := strings.IndexFunc(detail, func(r rune) bool {
			return r == '县' || r == '区'
		})
		if i == -1 {
			i = strings.IndexRune(detail, '市')
			if i == -1 {
				i = strings.IndexRune(detail, '省')
			}
		}
		prefix = append(prefix, detail[i+1:]...)
	}
	return string(prefix)
}
