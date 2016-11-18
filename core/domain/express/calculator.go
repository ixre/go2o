/**
 * Copyright 2015 @ z3q.net.
 * name : calculator
 * author : jarryliu
 * date : 2016-08-08 21:39
 * description :
 * history :
 */
package express

import (
	"go2o/core/domain/interface/express"
)

var _ express.IExpressCalculator = new(expressCalculatorImpl)

type expressComplex struct {
	Template express.IExpressTemplate
	Unit     float32
	Fee      float32
}

// 运费计算实现
type expressCalculatorImpl struct {
	userExpress *userExpressImpl
	tplMap      map[int32]*expressComplex
}

func newExpressCalculator(u *userExpressImpl) *expressCalculatorImpl {
	return &expressCalculatorImpl{
		userExpress: u,
		tplMap:      make(map[int32]*expressComplex),
	}
}

// 添加计算项,tplId为运费模板的编号
func (e *expressCalculatorImpl) Add(tplId int32, unit float32) error {
	ec, ok := e.tplMap[tplId]
	if !ok {
		tpl := e.userExpress.GetTemplate(tplId)
		if tpl == nil {
			return express.ErrNoSuchTemplate
		}
		ec = &expressComplex{
			Template: tpl,
			Fee:      0,
			Unit:     0,
		}
		e.tplMap[tplId] = ec
	}
	ec.Unit += unit
	return nil
}

// 计算运费
func (e *expressCalculatorImpl) Calculate(areaCode string) {
	for _, v := range e.tplMap {
		v.Fee = e.calculate(v.Template, areaCode, v.Unit)
	}
}

// 计算运费
func (e *expressCalculatorImpl) calculate(tpl express.IExpressTemplate,
	areaCode string, basisUnit float32) float32 {
	v := tpl.Value()
	//判断是否免邮
	if v.IsFree == 1 {
		return 0
	}
	//如果单位超出,则多加一个计量单位
	unit := int(basisUnit)
	if float32(unit) != basisUnit {
		unit += 1
	}
	//根据地区规则计算运费
	if areaCode != "" {
		areaSet := tpl.GetAreaExpressTemplateByAreaCode(areaCode)
		if areaSet != nil {
			return e.mathFee(v.Basis, unit, areaSet.FirstUnit, areaSet.FirstFee,
				areaSet.AddUnit, areaSet.AddFee)
		}
	}
	//根据默认规则计算运费
	return e.mathFee(v.Basis, unit, v.FirstUnit, v.FirstFee,
		v.AddUnit, v.AddFee)
}

// 计算快递运费
func (e *expressCalculatorImpl) mathFee(basis int, unit int, firstUnit int,
	firstFee float32, addUnit int, addFee float32) float32 {
	return e.getExpressFee(unit, firstUnit, firstFee, addUnit, addFee)
}

// 根据计量单位和值计算运费
func (e *expressCalculatorImpl) getExpressFee(unit int, firstUnit int,
	firstFee float32, addUnit int, addFee float32) float32 {
	outUnit := unit - firstUnit
	if outUnit > 0 {
		// 如果超过首次计量,则获取超出倍数,叠加计费
		outTimes := outUnit / addUnit
		if outUnit%addUnit > 0 {
			outTimes += 1
		}
		return firstFee + float32(outTimes)*addFee
	}
	return firstFee
}

// 获取累计运费
func (e *expressCalculatorImpl) Total() float32 {
	var total float32 = 0
	feeMap := e.Fee()
	for _, v := range feeMap {
		total += v
	}
	return total
}

// 获取运费模板编号与费用的集合
func (e *expressCalculatorImpl) Fee() map[int32]float32 {
	mp := make(map[int32]float32, len(e.tplMap))
	for k, v := range e.tplMap {
		mp[k] = v.Fee
	}
	return mp
}
