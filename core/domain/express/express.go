/**
 * Copyright 2015 @ z3q.net.
 * name : express
 * author : jarryliu
 * date : 2016-07-05 15:56
 * description :
 * history :
 */
package express

import (
	"go2o/core/domain/interface/express"
	"go2o/core/infrastructure/domain"
	"strings"
	"sync"
)

const (
	areaDelimer = ","
)

var _ express.IUserExpress = new(userExpressImpl)

type userExpressImpl struct {
	_userId int
	_arr    []express.IExpressTemplate
	_rep    express.IExpressRep
}

func NewUserExpress(userId int, rep express.IExpressRep) express.IUserExpress {
	return &userExpressImpl{
		_userId: userId,
		_rep:    rep,
	}
}

// 获取聚合根编号
func (this *userExpressImpl) GetAggregateRootId() int {
	return this._userId
}

// 创建快递模板
func (this *userExpressImpl) CreateTemplate(t *express.ExpressTemplate) express.IExpressTemplate {
	t.UserId = this.GetAggregateRootId()
	return newExpressTemplate(this, t, this._rep)
}

// 获取快递模板
func (this *userExpressImpl) GetTemplate(id int) express.IExpressTemplate {
	for _, v := range this.GetAllTemplate() {
		if v.GetDomainId() == id {
			return v
		}
	}
	return nil
}

// 获取所有的快递模板
func (this *userExpressImpl) GetAllTemplate() []express.IExpressTemplate {
	if this._arr == nil {
		list := this._rep.GetUserAllTemplate(this.GetAggregateRootId())
		this._arr = make([]express.IExpressTemplate, len(list))
		for i, v := range list {
			this._arr[i] = this.CreateTemplate(v)
		}
	}
	return this._arr
}

// 删除模板
func (this *userExpressImpl) DeleteTemplate(id int) error {
	for i, v := range this.GetAllTemplate() {
		if v.GetDomainId() == id {
			err := this._rep.DeleteExpressTemplate(
				this.GetAggregateRootId(), v.GetDomainId())
			if err == nil {
				this._arr = append(this._arr[:i], this._arr[i+1:]...)
			}
			return err
		}
	}
	return nil
}

// 计算快递运费
func (this *userExpressImpl) mathFee(unit int, firstUnit int,
	firstFee float32, addUnix int, addFee float32) float32 {
	outUnit := unit - firstUnit
	if outUnit > 0 {
		// 如果超过首次计量,则获取超出倍数,叠加计费
		outTimes := outUnit / addUnix
		if outUnit%addUnix > 0 {
			outTimes += 1
		}
		return firstFee + float32(outTimes)*addFee
	}
	return firstFee

}

// 获取快递费,传入地区编码，根据单位值，如总重量。
func (this *userExpressImpl) GetExpressFee(templateId int, areaCode string, unit int) float32 {
	tpl := this.GetTemplate(templateId)
	v := tpl.Value()
	//判断是否免邮
	if v.IsFree == 1 {
		return 0
	}
	//根据地区规则计算运费
	areaSet := tpl.GetAreaExpressTemplateByAreaCode(areaCode)
	if areaSet != nil {
		return this.mathFee(unit, areaSet.FirstUnit, areaSet.FirstFee,
			areaSet.AddUnit, areaSet.AddFee)
	}
	//根据默认规则计算运费
	return this.mathFee(unit, v.FirstUnit, v.FirstFee,
		v.AddUnit, v.AddFee)
}

var _ express.IExpressTemplate = new(expressTemplateImpl)

type expressTemplateImpl struct {
	_value       *express.ExpressTemplate
	_userExpress *userExpressImpl
	_rep         express.IExpressRep
	_areaList    []express.ExpressAreaTemplate
	_areaMap     map[string]*express.ExpressAreaTemplate
	_mux         sync.Mutex
}

func newExpressTemplate(u *userExpressImpl, v *express.ExpressTemplate,
	rep express.IExpressRep) express.IExpressTemplate {
	return &expressTemplateImpl{
		_value:       v,
		_userExpress: u,
		_rep:         rep,
	}
}

// 获取领域对象编号
func (this *expressTemplateImpl) GetDomainId() int {
	return this._value.Id
}

// 获取快递模板数据
func (this *expressTemplateImpl) Value() express.ExpressTemplate {
	return *this._value
}

// 设置地区的快递模板
func (this *expressTemplateImpl) Set(v *express.ExpressTemplate) error {
	if this._value.UserId != v.UserId ||
		v.Basis <= 0 || len(v.Name) == 0 ||
		v.AddFee <= 0 || v.AddUnit <= 0 ||
		v.FirstUnit <= 0 {
		return express.ErrNotFullExpressTemplate
	}
	this._value = v
	return nil
}

// 保存
func (this *expressTemplateImpl) Save() (int, error) {
	id, err := this._rep.SaveExpressTemplate(this._value)
	if err == nil {
		this._value.Id = id
		this._userExpress._arr = nil
	}
	return id, err
}

// 保存地区快递模板
func (this *expressTemplateImpl) SaveAreaTemplate(t *express.ExpressAreaTemplate) (int, error) {
	this.GetAllAreaTemplate()
	arr := this.getAreaCodeArray(t.CodeList)
	if arr == nil {
		return 0, express.ErrExpressTemplateMissingAreaCode
	}
	for _, code := range arr {
		v, ok := this._areaMap[code]
		if ok && v.Id != t.Id {
			return 0, express.ErrExistsAreaTemplateSet
		}
	}
	// 保存,如果未出错,则更新缓存
	id, err := this._rep.SaveExpressTemplateAreaSet(t)
	if err == nil {
		this._areaList = nil
		this._areaMap = nil
		this.GetAllAreaTemplate()
	}
	return id, err
}

func (this *expressTemplateImpl) getAreaCodeArray(codeList string) []string {
	codeList = strings.Trim(codeList, " ")
	if codeList == "" {
		return nil
	}
	return strings.Split(codeList, areaDelimer)
}

// 初始化地区与运费的映射
func (this *expressTemplateImpl) initAreaMap() {
	this._areaMap = make(map[string]*express.ExpressAreaTemplate, len(this._areaList))
	for _, v := range this._areaList {
		arr := this.getAreaCodeArray(v.CodeList)
		if arr == nil {
			continue
		}
		// 遍历地区编号
		for _, code := range arr {
			if code == "" {
				continue
			}
			if _, ok := this._areaMap[code]; !ok {
				this._areaMap[code] = &v
			}
		}
	}
}

// 获取所有的地区快递模板
func (this *expressTemplateImpl) GetAllAreaTemplate() []express.ExpressAreaTemplate {
	this._mux.Lock()
	if this._areaList == nil {
		this._areaList = this._rep.GetExpressTemplateAllAreaSet(this.GetDomainId())
		this.initAreaMap()
	}
	this._mux.Unlock()
	return this._areaList
}

// 根据地区编码获取运费模板
func (this *expressTemplateImpl) GetAreaExpressTemplateByAreaCode(areaCode string) *express.ExpressAreaTemplate {
	this.GetAllAreaTemplate()
	return this._areaMap[areaCode]
}

// 根据编号获取地区的运费模板
func (this *expressTemplateImpl) GetAreaExpressTemplate(id int) *express.ExpressAreaTemplate {
	for _, v := range this.GetAllAreaTemplate() {
		if v.Id == id {
			return &v
		}
	}
	return nil
}

type ExpressRepBase struct {
}

func (this *ExpressRepBase) SaveDefaultExpressProviders(rep express.IExpressRep) []*express.ExpressProvider {
	var err error
	for _, v := range express.SupportedExpressProvider {
		if v.Id, err = rep.SaveExpressProvider(v); err != nil {
			domain.HandleError(err, "domain")
		}
	}
	return express.SupportedExpressProvider
}
