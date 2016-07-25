/**
 * Copyright 2015 @ z3q.net.
 * name : value_rep
 * author : jarryliu
 * date : 2016-05-27 15:32
 * description :
 * history :
 */
package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jsix/gof/db"
	"github.com/jsix/gof/util"
	"go2o/core/domain/interface/valueobject"
	"strconv"
	"strings"
	"sync"
)

var _ valueobject.IValueRep = new(valueRep)

type valueRep struct {
	db.Connector
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
	_areaCache       map[int][]*valueobject.Area
	_areaMux         sync.Mutex
}

func NewValueRep(conn db.Connector) valueobject.IValueRep {
	return &valueRep{
		Connector: conn,
		_rstGob:   util.NewGobFile("conf/core/registry"),
		_wxGob:    util.NewGobFile("conf/core/wx_api"),
		_rpGob:    util.NewGobFile("conf/core/register_perm"),
		_numGob:   util.NewGobFile("conf/core/number_conf"),
		_mchGob:   util.NewGobFile("conf/core/mch_conf"),
		_mscGob:   util.NewGobFile("conf/core/mch_sale_conf"),
		_smsGob:   util.NewGobFile("conf/core/sms_conf"),
	}
}

// 获取微信接口配置
func (vp *valueRep) GetWxApiConfig() valueobject.WxApiConfig {
	if vp._wxConf == nil {
		vp._wxConf = &valueobject.WxApiConfig{}
		vp._wxGob.Unmarshal(vp._wxConf)
	}
	return *vp._wxConf
}

// 保存微信接口配置
func (vp *valueRep) SaveWxApiConfig(v *valueobject.WxApiConfig) error {
	if v != nil {
		//todo: 检查证书文件是否存在
		vp._wxConf = v
		return vp._wxGob.Save(vp._wxConf)
	}
	return errors.New("nil value")
}

// 获取注册权限
func (vp *valueRep) GetRegisterPerm() valueobject.RegisterPerm {
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
		vp._rpConf = v
		return vp._rpGob.Save(vp._rpConf)
	}
	return nil
}

// 获取全局系统销售设置
func (vp *valueRep) GetGlobNumberConf() valueobject.GlobNumberConf {
	if vp._numConf == nil {
		v := defaultGlobNumberConf
		vp._numConf = &v
		vp._numGob.Unmarshal(vp._numConf)
	}
	return *vp._numConf
}

// 保存全局系统销售设置
func (vp *valueRep) SaveGlobNumberConf(v *valueobject.GlobNumberConf) error {
	if v != nil {
		vp._numConf = v
		return vp._numGob.Save(vp._numConf)
	}
	return nil
}

// 获取平台设置
func (vp *valueRep) GetPlatformConf() valueobject.PlatformConf {
	if vp._globMchConf == nil {
		v := defaultPlatformConf
		vp._globMchConf = &v
		vp._mchGob.Unmarshal(vp._globMchConf)
	}
	return *vp._globMchConf
}

// 保存平台设置
func (vp *valueRep) SavePlatformConf(v *valueobject.PlatformConf) error {
	if v != nil {
		vp._globMchConf = v
		return vp._mchGob.Save(vp._globMchConf)
	}
	return nil
}

// 获取数据存储
func (v *valueRep) GetRegistry() valueobject.Registry {
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
		v._globRegistry = r
		return v._rstGob.Save(v._globRegistry)
	}
	return nil
}

// 获取全局商户销售设置
func (vp *valueRep) GetGlobMchSaleConf() valueobject.GlobMchSaleConf {
	if vp._globMchSaleConf == nil {
		v := defaultGlobMchSaleConf
		vp._globMchSaleConf = &v
		vp._mscGob.Unmarshal(vp._globMchSaleConf)
	}
	return *vp._globMchSaleConf
}

// 保存全局商户销售设置
func (vp *valueRep) SaveGlobMchSaleConf(v *valueobject.GlobMchSaleConf) error {
	if v != nil {
		vp._globMchSaleConf = v
		return vp._mscGob.Save(vp._globMchSaleConf)
	}
	return nil
}

// 获取短信设置
func (vp *valueRep) GetSmsApiSet() valueobject.SmsApiSet {
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
	vp._smsConf[provider] = s
	return vp._smsGob.Save(vp._smsConf)
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
func (vp *valueRep) GetChildAreas(id int) []*valueobject.Area {
	vp._areaMux.Lock()
	defer vp._areaMux.Unlock()
	if vp._areaCache == nil {
		vp._areaCache = make(map[int][]*valueobject.Area)
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
func (vp *valueRep) GetAreaNames(id []int) []string {
	strArr := make([]string, len(id))
	for i, v := range id {
		strArr[i] = strconv.Itoa(v)
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
