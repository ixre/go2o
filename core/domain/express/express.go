/**
 * Copyright 2015 @ 56x.net.
 * name : express
 * author : jarryliu
 * date : 2016-07-05 15:56
 * description :
 * history :
 */
package express

import (
	"math"
	"strconv"
	"strings"
	"sync"

	"github.com/ixre/go2o/core/domain/interface/express"
	"github.com/ixre/go2o/core/domain/interface/valueobject"
	"github.com/ixre/go2o/core/infrastructure/domain"
)

const (
	areaDelimer = ","
)

var _ express.IUserExpress = new(userExpressImpl)

type userExpressImpl struct {
	userId  int
	arr     []express.IExpressTemplate
	rep     express.IExpressRepo
	valRepo valueobject.IValueRepo
}

func NewUserExpress(userId int, rep express.IExpressRepo,
	valRepo valueobject.IValueRepo) express.IUserExpress {
	return &userExpressImpl{
		userId:  userId,
		rep:     rep,
		valRepo: valRepo,
	}
}

// 获取聚合根编号
func (e *userExpressImpl) GetAggregateRootId() int {
	return e.userId
}

// 创建快递模板
func (e *userExpressImpl) CreateTemplate(t *express.ExpressTemplate) express.IExpressTemplate {
	t.VendorId = int(e.GetAggregateRootId())
	return newExpressTemplate(e, t, e.rep, e.valRepo)
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
	if e.arr == nil {
		list := e.rep.GetUserAllTemplate(int(e.GetAggregateRootId()))
		e.arr = make([]express.IExpressTemplate, len(list))
		for i, v := range list {
			e.arr[i] = e.CreateTemplate(v)
		}
	}
	return e.arr
}

// 删除模板
func (e *userExpressImpl) DeleteTemplate(id int) error {
	for i, v := range e.GetAllTemplate() {
		if v.GetDomainId() == id {
			err := e.rep.DeleteExpressTemplate(
				int(e.GetAggregateRootId()), v.GetDomainId())
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
	_value           *express.ExpressTemplate
	_userExpress     *userExpressImpl
	_repo            express.IExpressRepo
	_regionList      []express.RegionExpressTemplate
	_regionIsChanged bool
	_areaMap         map[string]*express.RegionExpressTemplate
	_mux             sync.Mutex
	_valRepo         valueobject.IValueRepo
}

func newExpressTemplate(u *userExpressImpl, v *express.ExpressTemplate,
	rep express.IExpressRepo, valRepo valueobject.IValueRepo) express.IExpressTemplate {
	return &expressTemplateImpl{
		_value:       v,
		_userExpress: u,
		_repo:        rep,
		_valRepo:     valRepo,
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
	if e._value.VendorId > 0 && v.VendorId != e._value.VendorId {
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

// SetRegionExpress implements express.IExpressTemplate
func (e *expressTemplateImpl) SetRegionExpress(arr *[]express.RegionExpressTemplate) error {
	e._regionList = *arr
	e._regionIsChanged = true
	return nil
}

// 是否启用
func (e *expressTemplateImpl) Enabled() bool {
	return e._value.Enabled == 1
}

// 保存
func (e *expressTemplateImpl) Save() (int, error) {
	id, err := e._repo.SaveExpressTemplate(e._value)
	if err == nil {
		e._value.Id = id
		e._userExpress.arr = nil
		if err == nil && e._regionIsChanged {
			err = e.saveRegionExpress(e._regionList)
			e._regionIsChanged = false
		}
	}
	return id, err
}

// 保存区域快递模板
func (e *expressTemplateImpl) saveRegionExpress(items []express.RegionExpressTemplate) error {
	// 获取存在的项
	old := e._repo.GetExpressTemplateAllAreaSet(e.GetDomainId())
	// 分析当前项目并加入到MAP中
	delList := []int{}
	currMap := make(map[int]*express.RegionExpressTemplate, len(items))
	for _, v := range items {
		currMap[v.Id] = &v
	}
	// 筛选出要删除的项
	for _, v := range old {
		if currMap[v.Id] == nil {
			delList = append(delList, v.Id)
		}
	}
	// 删除项
	for _, v := range delList {
		e._repo.DeleteAreaExpressTemplate(e.GetDomainId(), v)
	}
	// 保存项
	for _, v := range items {
		i, err := e.saveAreaTemplate(&v)
		if err != nil {
			return err
		}
		v.Id = int(i)
	}
	e._regionList = items
	return nil
}

// 保存地区快递模板
func (e *expressTemplateImpl) saveAreaTemplate(t *express.RegionExpressTemplate) (int, error) {
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
	id, err := e._repo.SaveExpressTemplateAreaSet(t)
	if err == nil {
		e._regionList = nil
		e._areaMap = nil
		e.RegionExpress()
	}
	return id, err
}

// 获取地区代码列表
func (e *expressTemplateImpl) getAreaCodeArray(codeList string) []string {
	codeList = strings.Trim(codeList, " ")
	if codeList == "" {
		return nil
	}
	return strings.Split(codeList, areaDelimer)
}

// 初始化地区与运费的映射
func (e *expressTemplateImpl) initAreaMap() {
	e._areaMap = make(map[string]*express.RegionExpressTemplate, len(e._regionList))
	for _, v := range e._regionList {
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
func (e *expressTemplateImpl) RegionExpress() []express.RegionExpressTemplate {
	e._mux.Lock()
	if e._regionList == nil {
		e._regionList = e._repo.GetExpressTemplateAllAreaSet(e.GetDomainId())
		e.initAreaMap()
	}
	e._mux.Unlock()
	return e._regionList
}

// 根据地区编码获取运费模板
func (e *expressTemplateImpl) GetAreaExpressTemplateByAreaCode(areaCode string) *express.RegionExpressTemplate {
	e.RegionExpress()
	return e._areaMap[areaCode]
}

type ExpressRepBase struct {
}

// 将默认的快递服务商保存
func (e *ExpressRepBase) SaveDefaultExpressProviders(rep express.IExpressRepo) []*express.Provider {
	var err error
	for _, v := range express.SupportedExpressProvider {
		if v.Id, err = rep.SaveExpressProvider(v); err != nil {
			domain.HandleError(err, "domain")
		}
	}
	return express.SupportedExpressProvider
}
