/**
 * Copyright 2015 @ z3q.net.
 * name : express_service.go
 * author : jarryliu
 * date : 2016-07-05 18:57
 * description :
 * history :
 */
package dps

import "go2o/core/domain/interface/express"

type expressService struct {
	_rep express.IExpressRep
}

// 获取快递服务
func NewExpressService(rep express.IExpressRep) *expressService {
	return &expressService{
		_rep: rep,
	}
}

// 获取快递公司
func (this *expressService) GetExpressProvider(id int) *express.ExpressProvider {
	return this._rep.GetExpressProvider(id)
}

// 保存快递公司
func (this *expressService) SaveExpressProvider(v *express.ExpressProvider) (int, error) {
	return this._rep.SaveExpressProvider(v)
}

// 保存快递模板
func (this *expressService) SaveTemplate(userId int, v *express.ExpressTemplate) (int, error) {
	u := this._rep.GetUserExpress(userId)
	var e express.IExpressTemplate
	if v.Id > 0 {
		e = u.GetTemplate(v.Id)
	} else {
		e = u.CreateTemplate(&express.ExpressTemplate{})
	}
	err := e.Set(v)
	if err == nil {
		v.Id, err = e.Save()
	}
	return v.Id, err
}

// 获取快递模板
func (this *expressService) GetTemplate(userId, id int) *express.ExpressTemplate {
	u := this._rep.GetUserExpress(userId)
	t := u.GetTemplate(id)
	if t != nil {
		v := t.Value()
		return &v
	}
	return nil
}

// 获取所有的快递模板
func (this *expressService) GetAllTemplate(userId int) []*express.ExpressTemplate {
	u := this._rep.GetUserExpress(userId)
	list := u.GetAllTemplate()
	arr := make([]*express.ExpressTemplate, len(list))
	for i, v := range list {
		v2 := v.Value()
		arr[i] = &v2
	}
	return arr
}

// 删除模板
func (this *expressService) DeleteTemplate(userId int, id int) error {
	u := this._rep.GetUserExpress(userId)
	return u.DeleteTemplate(id)
}

// 获取快递费,传入地区编码，根据单位值，如总重量。
func (this *expressService) GetExpressFee(userId int, templateId int,
	areaCode string, unit int) float32 {
	u := this._rep.GetUserExpress(userId)
	return u.GetExpressFee(templateId, areaCode, unit)
}

// 根据地区编码获取运费模板
func (this *expressService) GetAreaExpressTemplateByAreaCode(userId int,
	templateId int, areaCode string) *express.ExpressAreaTemplate {
	u := this._rep.GetUserExpress(userId)
	t := u.GetTemplate(templateId)
	if t != nil {
		return t.GetAreaExpressTemplateByAreaCode(areaCode)
	}
	return nil
}

// 根据编号获取地区的运费模板
func (this *expressService) GetAreaExpressTemplate(userId int,
	templateId int, id int) *express.ExpressAreaTemplate {
	u := this._rep.GetUserExpress(userId)
	t := u.GetTemplate(templateId)
	if t != nil {
		return t.GetAreaExpressTemplate(id)
	}
	return nil
}

// 保存地区快递模板
func (this *expressService) SaveAreaTemplate(userId int,
	templateId int, v *express.ExpressAreaTemplate) (int, error) {
	u := this._rep.GetUserExpress(userId)
	t := u.GetTemplate(templateId)
	if t != nil {
		return t.SaveAreaTemplate(v)
	}
	return 0, nil
}

// 获取所有的地区快递模板
func (this *expressService) GetAllAreaTemplate(userId int,
	templateId int) []express.ExpressAreaTemplate {
	u := this._rep.GetUserExpress(userId)
	t := u.GetTemplate(templateId)
	if t != nil {
		return t.GetAllAreaTemplate()
	}
	return []express.ExpressAreaTemplate{}
}
