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
	"go2o/core/domain/interface/valueobject"
	"go2o/core/infrastructure/domain"
	"strconv"
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
	_valRep valueobject.IValueRep
}

func NewUserExpress(userId int, rep express.IExpressRep,
	valRep valueobject.IValueRep) express.IUserExpress {
	return &userExpressImpl{
		_userId: userId,
		_rep:    rep,
		_valRep: valRep,
	}
}

// 获取聚合根编号
func (e *userExpressImpl) GetAggregateRootId() int {
	return e._userId
}

// 创建快递模板
func (e *userExpressImpl) CreateTemplate(t *express.ExpressTemplate) express.IExpressTemplate {
	t.UserId = e.GetAggregateRootId()
	return newExpressTemplate(e, t, e._rep, e._valRep)
}

// 获取快递模板
func (e *userExpressImpl) GetTemplate(id int) express.IExpressTemplate {
	for _, v := range e.GetAllTemplate() {
		if v.GetDomainId() == id {
			return v
		}
	}
	return nil
}

// 获取所有的快递模板
func (e *userExpressImpl) GetAllTemplate() []express.IExpressTemplate {
	if e._arr == nil {
		list := e._rep.GetUserAllTemplate(e.GetAggregateRootId())
		e._arr = make([]express.IExpressTemplate, len(list))
		for i, v := range list {
			e._arr[i] = e.CreateTemplate(v)
		}
	}
	return e._arr
}

// 删除模板
func (e *userExpressImpl) DeleteTemplate(id int) error {
	for i, v := range e.GetAllTemplate() {
		if v.GetDomainId() == id {
			err := e._rep.DeleteExpressTemplate(
				e.GetAggregateRootId(), v.GetDomainId())
			if err == nil {
				e._arr = append(e._arr[:i], e._arr[i+1:]...)
			}
			return err
		}
	}
	return nil
}

// 计算快递运费
func (e *userExpressImpl) mathFee(basis int, unit int, firstUnit int,
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
func (e *userExpressImpl) GetExpressFee(templateId int, areaCode string, unit int) float32 {
	tpl := e.GetTemplate(templateId)
	if tpl == nil {
		//return 0
		//todo: 仅仅为测试,如果未指定快递模板,则用默认的第一个
		if len(e.GetAllTemplate()) > 0 {
			tpl = e.GetTemplate(e.GetAllTemplate()[0].GetDomainId())
		} else {
			return 0
		}
	}
	v := tpl.Value()

	//log.Println("--------", e._userId, unit,v)

	//判断是否免邮
	if v.IsFree == 1 {
		return 0
	}
	//根据地区规则计算运费
	areaSet := tpl.GetAreaExpressTemplateByAreaCode(areaCode)
	if areaSet != nil {
		return e.mathFee(v.Basis, unit, areaSet.FirstUnit, areaSet.FirstFee,
			areaSet.AddUnit, areaSet.AddFee)
	}
	//根据默认规则计算运费
	return e.mathFee(v.Basis, unit, v.FirstUnit, v.FirstFee,
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
	_valRep      valueobject.IValueRep
}

func newExpressTemplate(u *userExpressImpl, v *express.ExpressTemplate,
	rep express.IExpressRep, valRep valueobject.IValueRep) express.IExpressTemplate {
	return &expressTemplateImpl{
		_value:       v,
		_userExpress: u,
		_rep:         rep,
		_valRep:      valRep,
	}
}

// 获取领域对象编号
func (e *expressTemplateImpl) GetDomainId() int {
	return e._value.Id
}

// 获取快递模板数据
func (e *expressTemplateImpl) Value() express.ExpressTemplate {
	return *e._value
}

func (e *expressTemplateImpl) checkValue(v *express.ExpressTemplate) error {
	if v.Name == "" {
		return express.ErrExpressTemplateName
	}
	if e._value.UserId > 0 && v.UserId != e._value.UserId {
		return express.ErrUserNotMatch
	}
	// 如果不包邮,检查相关设置
	if v.IsFree == 0 {
		if v.Basis <= 0 {
			return express.ErrExpressBasis
		}
		if v.FirstUnit <= 0 {
			return express.ErrFirstUnitNotSet
		}

		if v.AddUnit <= 0 {
			return express.ErrAddUnitNotSet
		}
		if v.AddFee <= 0 {
			return express.ErrAddFee
		}
	}
	return nil
}

// 设置地区的快递模板
func (e *expressTemplateImpl) Set(v *express.ExpressTemplate) error {
	v.Name = strings.TrimSpace(v.Name)
	if v.FirstFee < 0 {
		v.FirstFee = -v.FirstFee
	}
	err := e.checkValue(v)
	if err == nil {
		e._value = v
	}
	return err
}

// 保存
func (e *expressTemplateImpl) Save() (int, error) {
	id, err := e._rep.SaveExpressTemplate(e._value)
	if err == nil {
		e._value.Id = id
		e._userExpress._arr = nil
	}
	return id, err
}

// 保存地区快递模板
func (e *expressTemplateImpl) SaveAreaTemplate(t *express.ExpressAreaTemplate) (int, error) {
	e.GetAllAreaTemplate()
	arr := e.getAreaCodeArray(t.CodeList)
	if arr == nil {
		return 0, express.ErrExpressTemplateMissingAreaCode
	}
	intArr := make([]int, len(arr))
	for i, code := range arr {
		v, ok := e._areaMap[code]
		if ok && v.Id != t.Id {
			return 0, express.ErrExistsAreaTemplateSet
		}
		intArr[i], _ = strconv.Atoi(code)
	}
	// 获取对应的中文名称
	names := e._valRep.GetAreaNames(intArr)
	t.NameList = strings.Join(names, ",")

	// 保存,如果未出错,则更新缓存
	id, err := e._rep.SaveExpressTemplateAreaSet(t)
	if err == nil {
		e._areaList = nil
		e._areaMap = nil
		e.GetAllAreaTemplate()
	}
	return id, err
}

func (e *expressTemplateImpl) getAreaCodeArray(codeList string) []string {
	codeList = strings.Trim(codeList, " ")
	if codeList == "" {
		return nil
	}
	return strings.Split(codeList, areaDelimer)
}

// 初始化地区与运费的映射
func (e *expressTemplateImpl) initAreaMap() {
	e._areaMap = make(map[string]*express.ExpressAreaTemplate, len(e._areaList))
	for _, v := range e._areaList {
		arr := e.getAreaCodeArray(v.CodeList)
		if arr == nil {
			continue
		}
		// 遍历地区编号
		for _, code := range arr {
			if code == "" {
				continue
			}
			if _, ok := e._areaMap[code]; !ok {
				e._areaMap[code] = &v
			}
		}
	}
}

// 获取所有的地区快递模板
func (e *expressTemplateImpl) GetAllAreaTemplate() []express.ExpressAreaTemplate {
	e._mux.Lock()
	if e._areaList == nil {
		e._areaList = e._rep.GetExpressTemplateAllAreaSet(e.GetDomainId())
		e.initAreaMap()
	}
	e._mux.Unlock()
	return e._areaList
}

// 删除模板地区设定
func (e *expressTemplateImpl) DeleteAreaSet(areaSetId int) error {
	e.GetAllAreaTemplate()
	if e.GetAreaExpressTemplate(areaSetId) != nil {
		err := e._rep.DeleteAreaExpressTemplate(e.GetDomainId(), areaSetId)
		if err == nil {
			e._areaList = nil
			e._areaMap = nil
		}
		return err
	}
	return nil
}

// 根据地区编码获取运费模板
func (e *expressTemplateImpl) GetAreaExpressTemplateByAreaCode(areaCode string) *express.ExpressAreaTemplate {
	e.GetAllAreaTemplate()
	return e._areaMap[areaCode]
}

// 根据编号获取地区的运费模板
func (e *expressTemplateImpl) GetAreaExpressTemplate(id int) *express.ExpressAreaTemplate {
	for _, v := range e.GetAllAreaTemplate() {
		if v.Id == id {
			return &v
		}
	}
	return nil
}

type ExpressRepBase struct {
}

// 将默认的快递服务商保存
func (e *ExpressRepBase) SaveDefaultExpressProviders(rep express.IExpressRep) []*express.ExpressProvider {
	var err error
	for _, v := range express.SupportedExpressProvider {
		if v.Id, err = rep.SaveExpressProvider(v); err != nil {
			domain.HandleError(err, "domain")
		}
	}
	return express.SupportedExpressProvider
}
