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
	"github.com/jsix/gof/storage"
	"github.com/jsix/gof/util"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/infrastructure/tool/sms"
	"strconv"
	"strings"
	"sync"
)

var _ valueobject.IValueRep = new(valueRep)
var (
	valueRepCacheKey = "go2o:rep:value-rep:cache"
)

type valueRep struct {
	db.Connector
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

func NewValueRep(conn db.Connector, storage storage.Interface) valueobject.IValueRep {
	return &valueRep{
		Connector: conn,
		storage:   storage,
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

func (vp *valueRep) checkReload() error {
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

func (vp *valueRep) signReload() {
	vp.storage.Set(valueRepCacheKey, 0)
}

// 获取微信接口配置
func (vp *valueRep) GetWxApiConfig() valueobject.WxApiConfig {
	vp.checkReload()
	if vp._wxConf == nil {
		vp._wxConf = &valueobject.WxApiConfig{}
		vp._wxGob.Unmarshal(vp._wxConf)
	}
	return *vp._wxConf
}

// 保存微信接口配置
func (vp *valueRep) SaveWxApiConfig(v *valueobject.WxApiConfig) error {
	if v != nil {
		defer vp.signReload()
		//todo: 检查证书文件是否存在
		vp._wxConf = v
		return vp._wxGob.Save(vp._wxConf)
	}
	return errors.New("nil value")
}

// 获取注册权限
func (vp *valueRep) GetRegisterPerm() valueobject.RegisterPerm {
	vp.checkReload()
	if vp._rpConf == nil {
		v := defaultRegisterPerm
		vp._rpConf = &v
		vp._rpGob.Unmarshal(vp._rpConf)
	}
	return *vp._rpConf
}

// 保存注册权限
func (vp *valueRep) SaveRegisterPerm(v *valueobject.RegisterPerm) error {
	if v != nil {
		defer vp.signReload()
		vp._rpConf = v
		return vp._rpGob.Save(vp._rpConf)
	}
	return nil
}

// 获取全局系统销售设置
func (vp *valueRep) GetGlobNumberConf() valueobject.GlobNumberConf {
	vp.checkReload()
	if vp._numConf == nil {
		v := DefaultGlobNumberConf
		vp._numConf = &v
		vp._numGob.Unmarshal(vp._numConf)
	}
	return *vp._numConf
}

// 保存全局系统销售设置
func (vp *valueRep) SaveGlobNumberConf(v *valueobject.GlobNumberConf) error {
	if v != nil {
		defer vp.signReload()
		vp._numConf = v
		return vp._numGob.Save(vp._numConf)
	}
	return nil
}

// 获取平台设置
func (vp *valueRep) GetPlatformConf() valueobject.PlatformConf {
	vp.checkReload()
	if vp._globMchConf == nil {
		v := DefaultPlatformConf
		vp._globMchConf = &v
		vp._mchGob.Unmarshal(vp._globMchConf)
	}
	return *vp._globMchConf
}

// 保存平台设置
func (vp *valueRep) SavePlatformConf(v *valueobject.PlatformConf) error {
	if v != nil {
		defer vp.signReload()
		vp._globMchConf = v
		return vp._mchGob.Save(vp._globMchConf)
	}
	return nil
}

// 获取模板配置
func (v *valueRep) GetTemplateConf() valueobject.TemplateConf {
	v.checkReload()
	if v._tplConf == nil {
		v2 := DefaultTemplateConf
		v._tplConf = &v2
		v._tplGob.Unmarshal(v._tplConf)
	}
	return *v._tplConf
}

// 保存模板配置
func (v *valueRep) SaveTemplateConf(t *valueobject.TemplateConf) error {
	if t != nil {
		defer v.signReload()
		v._tplConf = t
		return v._tplGob.Save(v._tplConf)
	}
	return nil
}

// 获取移动应用设置
func (v *valueRep) GetMoAppConf() valueobject.MoAppConf {
	v.checkReload()
	if v._moAppConf == nil {
		v2 := DefaultMoAppConf
		v._moAppConf = &v2
		v._moAppGob.Unmarshal(v._moAppConf)
	}
	return *v._moAppConf
}

// 保存移动应用设置
func (v *valueRep) SaveMoAppConf(r *valueobject.MoAppConf) error {
	if r != nil {
		defer v.signReload()
		v._moAppConf = r
		return v._moAppGob.Save(v._moAppConf)
	}
	return nil
}

// 获取数据存储
func (v *valueRep) GetRegistry() valueobject.Registry {
	v.checkReload()
	if v._globRegistry == nil {
		v2 := DefaultRegistry
		v._globRegistry = &v2
		v._rstGob.Unmarshal(v._globRegistry)
	}
	return *v._globRegistry
}

// 保存数据存储
func (v *valueRep) SaveRegistry(r *valueobject.Registry) error {
	if r != nil {
		defer v.signReload()
		v._globRegistry = r
		return v._rstGob.Save(v._globRegistry)
	}
	return nil
}

// 获取全局商户销售设置
func (vp *valueRep) GetGlobMchSaleConf() valueobject.GlobMchSaleConf {
	vp.checkReload()
	if vp._globMchSaleConf == nil {
		v := DefaultGlobMchSaleConf
		vp._globMchSaleConf = &v
		vp._mscGob.Unmarshal(vp._globMchSaleConf)
	}
	return *vp._globMchSaleConf
}

// 保存全局商户销售设置
func (vp *valueRep) SaveGlobMchSaleConf(v *valueobject.GlobMchSaleConf) error {
	if v != nil {
		defer vp.signReload()
		vp._globMchSaleConf = v
		return vp._mscGob.Save(vp._globMchSaleConf)
	}
	return nil
}

// 获取短信设置
func (vp *valueRep) GetSmsApiSet() valueobject.SmsApiSet {
	vp.checkReload()
	if vp._smsConf == nil {
		vp._smsConf = defaultSmsConf
		vp._smsGob.Unmarshal(&vp._smsConf)
	}
	return vp._smsConf
}

// 保存短信API
func (vp *valueRep) SaveSmsApiPerm(provider int, s *valueobject.SmsApiPerm) error {
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
func (vp *valueRep) GetDefaultSmsApiPerm() (int, *valueobject.SmsApiPerm) {
	for i, v := range vp.GetSmsApiSet() {
		if v.Default {
			return i, v
		}
	}
	panic(errors.New("至少为系统设置一个短信接口"))
}

// 获取下级区域
func (vp *valueRep) GetChildAreas(id int32) []*valueobject.Area {
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
func (vp *valueRep) GetAreaNames(id []int32) []string {
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
func (vp *valueRep) GetAreaString(province, city, district int32) string {
	names := vp.GetAreaNames([]int32{province, city, district})
	if names[1] == "市辖区" || names[1] == "市辖县" || names[1] == "县" {
		return strings.Join([]string{names[0], names[2]}, " ")
	}
	return strings.Join(names, " ")
}
