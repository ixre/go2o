package express

import (
	"go2o/core/domain/interface/express"
)

var _ express.IExpressCalculator = new(expressCalculatorImpl)

// 运费，单位使用克/毫升，计算时按千克/升来计算
// 界面设置也应按照千克/升来设置
type expressComplex struct {
	Template express.IExpressTemplate
	Unit     int
	Fee      float64
}

// 运费计算实现
type expressCalculatorImpl struct {
	userExpress *userExpressImpl
	tplMap      map[int]*expressComplex
}

func newExpressCalculator(u *userExpressImpl) *expressCalculatorImpl {
	return &expressCalculatorImpl{
		userExpress: u,
		tplMap:      make(map[int]*expressComplex),
	}
}

// 添加计算项,tplId为运费模板的编号
func (e *expressCalculatorImpl) Add(tplId int, unit int) error {
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
	areaCode string, basisUnit int) float64 {
	v := tpl.Value()
	//如果免邮或计价数值为零，则不计算运费
	if v.IsFree == 1 || basisUnit == 0 {
		return 0
	}
	var trimUnit float32 //实际计价单位数量(转换后的数量)
	var finalUnit int    //最终计价单位数量

	switch v.Basis {
	case express.BasisByNumber:
		finalUnit = basisUnit
	case express.BasisByVolume:
		// ML -> L ; 1L = 500ML
		trimUnit = float32(basisUnit) / 500
	case express.BasisByWeight:
		// G  -> KG; 1KG = 1000g
		trimUnit = float32(basisUnit) / 1000
	}
	//如果单位超出,则多加一个计量单位
	if trimUnit > 0 {
		finalUnit = int(trimUnit)
		if trimUnit-float32(finalUnit) > 0 {
			finalUnit += 1
		}
	}

	//根据地区规则计算运费
	if areaCode != "" {
		areaSet := tpl.GetAreaExpressTemplateByAreaCode(areaCode)
		if areaSet != nil {
			return e.mathFee(v.Basis, finalUnit, int(areaSet.FirstUnit),
				float64(areaSet.FirstFee), int(areaSet.AddUnit), float64(areaSet.AddFee))
		}
	}
	//根据默认规则计算运费
	return e.mathFee(v.Basis, finalUnit, v.FirstUnit, v.FirstFee,
		v.AddUnit, v.AddFee)
}

// 计算快递运费
func (e *expressCalculatorImpl) mathFee(basis int, unit, firstUnit int,
	firstFee float64, addUnit int, addFee float64) float64 {
	return e.getExpressFee(unit, firstUnit, firstFee, addUnit, addFee)
}

// 根据计量单位和值计算运费
func (e *expressCalculatorImpl) getExpressFee(unit, firstUnit int,
	firstFee float64, addUnit int, addFee float64) float64 {
	outUnit := unit - firstUnit
	if outUnit > 0 {
		// 如果超过首次计量,则获取超出倍数,叠加计费
		outTimes := outUnit / addUnit
		if outUnit%addUnit > 0 {
			outTimes += 1
		}
		return firstFee + float64(outTimes)*addFee
	}
	return firstFee
}

// 获取累计运费
func (e *expressCalculatorImpl) Total() float64 {
	var total float64 = 0
	feeMap := e.Fee()
	for _, v := range feeMap {
		total += v
	}
	return total
}

// 获取运费模板编号与费用的集合
func (e *expressCalculatorImpl) Fee() map[int]float64 {
	mp := make(map[int]float64, len(e.tplMap))
	for k, v := range e.tplMap {
		mp[k] = v.Fee
	}
	return mp
}
