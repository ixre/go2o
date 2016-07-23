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
func (this *valueRep) GetWxApiConfig() valueobject.WxApiConfig {
	if this._wxConf == nil {
		this._wxConf = &valueobject.WxApiConfig{}
		this._wxGob.Unmarshal(this._wxConf)
	}
	return *this._wxConf
}

// 保存微信接口配置
func (this *valueRep) SaveWxApiConfig(v *valueobject.WxApiConfig) error {
	if v != nil {
		//todo: 检查证书文件是否存在
		this._wxConf = v
		return this._wxGob.Save(this._wxConf)
	}
	return errors.New("nil value")
}

// 获取注册权限
func (this *valueRep) GetRegisterPerm() valueobject.RegisterPerm {
	if this._rpConf == nil {
		v := defaultRegisterPerm
		this._rpConf = &v
		this._rpGob.Unmarshal(this._rpConf)
	}
	return *this._rpConf
}

// 保存注册权限
func (this *valueRep) SaveRegisterPerm(v *valueobject.RegisterPerm) error {
	if v != nil {
		this._rpConf = v
		return this._rpGob.Save(this._rpConf)
	}
	return nil
}

// 获取全局系统销售设置
func (this *valueRep) GetGlobNumberConf() valueobject.GlobNumberConf {
	if this._numConf == nil {
		v := defaultGlobNumberConf
		this._numConf = &v
		this._numGob.Unmarshal(this._numConf)
	}
	return *this._numConf
}

// 保存全局系统销售设置
func (this *valueRep) SaveGlobNumberConf(v *valueobject.GlobNumberConf) error {
	if v != nil {
		this._numConf = v
		return this._numGob.Save(this._numConf)
	}
	return nil
}

// 获取平台设置
func (this *valueRep) GetPlatformConf() valueobject.PlatformConf {
	if this._globMchConf == nil {
		v := defaultPlatformConf
		this._globMchConf = &v
		this._mchGob.Unmarshal(this._globMchConf)
	}
	return *this._globMchConf
}

// 保存平台设置
func (this *valueRep) SavePlatformConf(v *valueobject.PlatformConf) error {
	if v != nil {
		this._globMchConf = v
		return this._mchGob.Save(this._globMchConf)
	}
	return nil
}

// 获取数据存储
func (v *valueRep) GetRegistry() valueobject.Registry {
	if v._globRegistry == nil {
		v2 := defaultRegistry
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
func (this *valueRep) GetGlobMchSaleConf() valueobject.GlobMchSaleConf {
	if this._globMchSaleConf == nil {
		v := defaultGlobMchSaleConf
		this._globMchSaleConf = &v
		this._mscGob.Unmarshal(this._globMchSaleConf)
	}
	return *this._globMchSaleConf
}

// 保存全局商户销售设置
func (this *valueRep) SaveGlobMchSaleConf(v *valueobject.GlobMchSaleConf) error {
	if v != nil {
		this._globMchSaleConf = v
		return this._mscGob.Save(this._globMchSaleConf)
	}
	return nil
}

// 获取短信设置
func (this *valueRep) GetSmsApiSet() valueobject.SmsApiSet {
	if this._smsConf == nil {
		this._smsConf = defaultSmsConf
		this._smsGob.Unmarshal(&this._smsConf)
	}
	return this._smsConf
}

// 保存短信API
func (this *valueRep) SaveSmsApiPerm(provider int, s *valueobject.SmsApiPerm) error {
	if _, ok := this.GetSmsApiSet()[provider]; !ok {
		return errors.New("系统不支持的短信接口")
	}

	if s.Default {
		// 取消其他接口的默认选项
		for p, v := range this._smsConf {
			if p == provider {
				v.Default = true
			} else {
				v.Default = false
			}
		}
	} else {
		//检验是否取消了正在使用的短信接口
		if i, _ := this.GetDefaultSmsApiPerm(); i == provider {
			return errors.New("系统应启用一个短信接口")
		}
	}
	this._smsConf[provider] = s
	return this._smsGob.Save(this._smsConf)
}

// 获取默认的短信API
func (this *valueRep) GetDefaultSmsApiPerm() (int, *valueobject.SmsApiPerm) {
	for i, v := range this.GetSmsApiSet() {
		if v.Default {
			return i, v
		}
	}
	panic(errors.New("至少为系统设置一个短信接口"))
}

// 获取下级区域
func (this *valueRep) GetChildAreas(id int) []*valueobject.Area {
	this._areaMux.Lock()
	defer this._areaMux.Unlock()
	if this._areaCache == nil {
		this._areaCache = make(map[int][]*valueobject.Area)
	}
	if v, ok := this._areaCache[id]; ok {
		return v
	}
	v := []*valueobject.Area{}
	err := this.Connector.GetOrm().Select(&v, "code <> 0 AND parent=?", id)
	if err == nil {
		this._areaCache[id] = v
	}
	return v
}

// 获取地区名称
func (this *valueRep) GetAreaNames(id []int) []string {
	strArr := make([]string, len(id))
	for i, v := range id {
		strArr[i] = strconv.Itoa(v)
	}
	i := 0
	this.Connector.Query(fmt.Sprintf(
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
