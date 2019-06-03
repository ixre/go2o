/**
 * Copyright 2015 @ to2.net.
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
	"math"
	"strconv"
	"strings"
	"sync"
)

const (
	areaDelimer = ","
)

var _ express.IUserExpress = new(userExpressImpl)

type userExpressImpl struct {
	userId  int32
	arr     []express.IExpressTemplate
	rep     express.IExpressRepo
	valRepo valueobject.IValueRepo
}

func NewUserExpress(userId int32, rep express.IExpressRepo,
	valRepo valueobject.IValueRepo) express.IUserExpress {
	return &userExpressImpl{
		userId:  userId,
		rep:     rep,
		valRepo: valRepo,
	}
}

// 获取聚合根编号
func (e *userExpressImpl) GetAggregateRootId() int32 {
	return e.userId
}

// 创建快递模板
func (e *userExpressImpl) CreateTemplate(t *express.ExpressTemplate) express.IExpressTemplate {
	t.UserId = e.GetAggregateRootId()
	return newExpressTemplate(e, t, e.rep, e.valRepo)
}

// 获取快递模板
func (e *userExpressImpl) GetTemplate(id int32) express.IExpressTemplate {
	for _, v := range e.GetAllTemplate() {
		if v.GetDomainId() == id {
			return v
		}
	}
	return nil
}

// 获取所有的快递模板
func (e *userExpressImpl) GetAllTemplate() []express.IExpressTemplate {
	if e.arr == nil {
		list := e.rep.GetUserAllTemplate(e.GetAggregateRootId())
		e.arr = make([]express.IExpressTemplate, len(list))
		for i, v := range list {
			e.arr[i] = e.CreateTemplate(v)
		}
	}
	return e.arr
}

// 删除模板
func (e *userExpressImpl) DeleteTemplate(id int32) error {
	for i, v := range e.GetAllTemplate() {
		if v.GetDomainId() == id {
			err := e.rep.DeleteExpressTemplate(
				e.GetAggregateRootId(), v.GetDomainId())
			if err == nil {
				e.arr = append(e.arr[:i], e.arr[i+1:]...)
			}
			return err
		}
	}
	return nil
}

// 创建运费计算器
func (e *userExpressImpl) CreateCalculator() express.IExpressCalculator {
	return newExpressCalculator(e)
}

var _ express.IExpressTemplate = new(expressTemplateImpl)

type expressTemplateImpl struct {
	_value       *express.ExpressTemplate
	_userExpress *userExpressImpl
	_rep         express.IExpressRepo
	_areaList    []express.ExpressAreaTemplate
	_areaMap     map[string]*express.ExpressAreaTemplate
	_mux         sync.Mutex
	_valRepo     valueobject.IValueRepo
}

func newExpressTemplate(u *userExpressImpl, v *express.ExpressTemplate,
	rep express.IExpressRepo, valRepo valueobject.IValueRepo) express.IExpressTemplate {
	return &expressTemplateImpl{
		_value:       v,
		_userExpress: u,
		_rep:         rep,
		_valRepo:     valRepo,
	}
}

// 获取领域对象编号
func (e *expressTemplateImpl) GetDomainId() int32 {
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
		// 不要求首重费用，因为可能存在首重免费、超出收费的情况
		if v.FirstUnit <= 0 {
			return express.ErrFirstUnitNotSet
		}

		if v.AddUnit <= 0 {
			return express.ErrAddUnitNotSet
		}
		if v.AddFee <= 0 || math.IsNaN(float64(v.AddFee)) {
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

// 是否启用
func (e *expressTemplateImpl) Enabled() bool {
	return e._value.Enabled == 1
}

// 保存
func (e *expressTemplateImpl) Save() (int32, error) {
	id, err := e._rep.SaveExpressTemplate(e._value)
	if err == nil {
		e._value.Id = id
		e._userExpress.arr = nil
	}
	return id, err
}

// 保存地区快递模板
func (e *expressTemplateImpl) SaveAreaTemplate(t *express.ExpressAreaTemplate) (int32, error) {
	e.GetAllAreaTemplate()
	arr := e.getAreaCodeArray(t.CodeList)
	if arr == nil {
		return 0, express.ErrExpressTemplateMissingAreaCode
	}
	intArr := make([]int32, len(arr))
	for i, code := range arr {
		v, ok := e._areaMap[code]
		if ok && v.Id != t.Id {
			return 0, express.ErrExistsAreaTemplateSet
		}
		i2, _ := strconv.Atoi(code)
		intArr[i] = int32(i2)
	}
	// 获取对应的中文名称
	names := e._valRepo.GetAreaNames(intArr)
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
func (e *expressTemplateImpl) DeleteAreaSet(areaSetId int32) error {
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
func (e *expressTemplateImpl) GetAreaExpressTemplate(id int32) *express.ExpressAreaTemplate {
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
func (e *ExpressRepBase) SaveDefaultExpressProviders(rep express.IExpressRepo) []*express.ExpressProvider {
	var err error
	for _, v := range express.SupportedExpressProvider {
		if v.Id, err = rep.SaveExpressProvider(v); err != nil {
			domain.HandleError(err, "domain")
		}
	}
	return express.SupportedExpressProvider
}
