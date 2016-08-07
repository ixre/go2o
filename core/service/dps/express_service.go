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
func (e *expressService) GetExpressProvider(id int) *express.ExpressProvider {
	return e._rep.GetExpressProvider(id)
}

// 保存快递公司
func (e *expressService) SaveExpressProvider(v *express.ExpressProvider) (int, error) {
	return e._rep.SaveExpressProvider(v)
}

// 获取可用的快递公司
func (e *expressService) GetEnabledProviders() []*express.ExpressProvider {
	arr := []*express.ExpressProvider{}
	list := e._rep.GetExpressProviders()
	for _, v := range list {
		if v.Enabled == 1 {
			arr = append(arr, v)
		}
	}
	return arr
}

// 保存快递模板
func (ec *expressService) SaveTemplate(userId int, v *express.ExpressTemplate) (int, error) {
	u := ec._rep.GetUserExpress(userId)
	var e express.IExpressTemplate
	if v.Id > 0 {
		e = u.GetTemplate(v.Id)
	} else {
		e = u.CreateTemplate(&express.ExpressTemplate{
			UserId: userId,
		})
	}
	err := e.Set(v)
	if err == nil {
		v.Id, err = e.Save()
	}
	return v.Id, err
}

// 获取快递模板
func (e *expressService) GetTemplate(userId, id int) *express.ExpressTemplate {
	u := e._rep.GetUserExpress(userId)
	t := u.GetTemplate(id)
	if t != nil {
		v := t.Value()
		return &v
	}
	return nil
}

// 获取所有的快递模板
func (e *expressService) GetAllTemplate(userId int) []*express.ExpressTemplate {
	u := e._rep.GetUserExpress(userId)
	list := u.GetAllTemplate()
	arr := make([]*express.ExpressTemplate, len(list))
	for i, v := range list {
		v2 := v.Value()
		arr[i] = &v2
	}
	return arr
}

// 删除模板
func (e *expressService) DeleteTemplate(userId int, id int) error {
	u := e._rep.GetUserExpress(userId)
	return u.DeleteTemplate(id)
}

// 删除模板地区设定
func (e *expressService) DeleteTemplateAreaSet(userId, id, areaSetId int) error {
	u := e._rep.GetUserExpress(userId)
	t := u.GetTemplate(id)
	if t == nil {
		return express.ErrNoSuchTemplate
	}
	return t.DeleteAreaSet(areaSetId)
}

// 获取快递费,传入地区编码，根据单位值，如总重量。
func (e *expressService) GetExpressFee(userId int, templateId int,
	areaCode string, unit int) float32 {
	u := e._rep.GetUserExpress(userId)
	return u.GetExpressFee(templateId, areaCode, unit)
}

// 根据地区编码获取运费模板
func (e *expressService) GetAreaExpressTemplateByAreaCode(userId int,
	templateId int, areaCode string) *express.ExpressAreaTemplate {
	u := e._rep.GetUserExpress(userId)
	t := u.GetTemplate(templateId)
	if t != nil {
		return t.GetAreaExpressTemplateByAreaCode(areaCode)
	}
	return nil
}

// 根据编号获取地区的运费模板
func (e *expressService) GetAreaExpressTemplate(userId int,
	templateId int, id int) *express.ExpressAreaTemplate {
	u := e._rep.GetUserExpress(userId)
	t := u.GetTemplate(templateId)
	if t != nil {
		return t.GetAreaExpressTemplate(id)
	}
	return nil
}

// 保存地区快递模板
func (e *expressService) SaveAreaTemplate(userId int,
	templateId int, v *express.ExpressAreaTemplate) (int, error) {
	u := e._rep.GetUserExpress(userId)
	t := u.GetTemplate(templateId)
	if t != nil {
		return t.SaveAreaTemplate(v)
	}
	return 0, nil
}

// 获取所有的地区快递模板
func (e *expressService) GetAllAreaTemplate(userId int,
	templateId int) []express.ExpressAreaTemplate {
	u := e._rep.GetUserExpress(userId)
	t := u.GetTemplate(templateId)
	if t != nil {
		return t.GetAllAreaTemplate()
	}
	return []express.ExpressAreaTemplate{}
}
