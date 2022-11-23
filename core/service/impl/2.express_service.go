/**
 * Copyright 2015 @ 56x.net.
 * name : 2.express_service.go
 * author : jarryliu
 * date : 2016-07-05 18:57
 * description :
 * history :
 */
package impl

import (
	"context"

	"github.com/ixre/go2o/core/domain/interface/express"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/types"
)

var _ proto.ExpressServiceServer = new(expressService)

type expressService struct {
	_rep express.IExpressRepo
	serviceUtil
	proto.UnimplementedExpressServiceServer
}

// 获取快递服务
func NewExpressService(rep express.IExpressRepo) *expressService {
	return &expressService{
		_rep: rep,
	}
}

// 获取快递公司
func (e *expressService) GetExpressProvider(_ context.Context, name *proto.IdOrName) (*proto.SExpressProvider, error) {
	var v *express.Provider
	if name.Id > 0 {
		v = e._rep.GetExpressProvider(int32(name.Id))
	} else {
		//v = e._rep.GetExpressProviderByName(name.Name)
	}
	if v != nil {
		return e.parseProviderDto(v), nil
	}
	return nil, nil
}

// 保存快递公司
func (e *expressService) SaveExpressProvider(_ context.Context, r *proto.SExpressProvider) (*proto.Result, error) {
	v := e.parseProvider(r)
	_, err := e._rep.SaveExpressProvider(v)
	return e.error(err), nil
}

// 获取卖家的快递公司
func (e *expressService) GetProviders(_ context.Context, _ *proto.Empty) (*proto.ExpressProviderListResponse, error) {
	var arr []*proto.SExpressProvider
	list := e._rep.GetExpressProviders()
	for _, v := range list {
		if v.Enabled == 1 {
			arr = append(arr, e.parseProviderDto(v))
		}
	}
	return &proto.ExpressProviderListResponse{
		Value: arr,
	}, nil
}

// 保存快递模板
func (e *expressService) SaveTemplate(_ context.Context, r *proto.SExpressTemplate) (*proto.SaveTemplateResponse, error) {
	u := e._rep.GetUserExpress(int(r.SellerId))
	v := e.parseExpressTemplate(r)
	var ie express.IExpressTemplate
	if r.Id > 0 {
		ie = u.GetTemplate(int(r.Id))
	} else {
		ie = u.CreateTemplate(&express.ExpressTemplate{
			VendorId: int(r.SellerId),
		})
	}
	var id int
	err := ie.Set(v)
	if err == nil {
		id, err = ie.Save()
	}
	ret := &proto.SaveTemplateResponse{
		TemplateId: int64(id),
	}
	if err != nil {
		ret.ErrCode = 1
		ret.ErrMsg = err.Error()
	}
	return ret, nil
}

// 获取快递模板
func (e *expressService) GetTemplate(_ context.Context, id *proto.ExpressTemplateId) (*proto.SExpressTemplate, error) {
	u := e._rep.GetUserExpress(int(id.SellerId))
	t := u.GetTemplate(int(id.TemplateId))
	if t != nil {
		v := t.Value()
		return e.parseExpressTemplateDto(&v), nil
	}
	return nil, nil
}

// 获取所有的快递模板
func (e *expressService) GetAllTemplate(userId int32) []*express.ExpressTemplate {
	u := e._rep.GetUserExpress(int(userId))
	list := u.GetAllTemplate()
	arr := make([]*express.ExpressTemplate, len(list))
	for i, v := range list {
		v2 := v.Value()
		arr[i] = &v2
	}
	return arr
}

// 获取可有的快递模板
func (e *expressService) GetTemplates(_ context.Context, r *proto.GetTemplatesRequest) (*proto.ExpressTemplateListResponse, error) {
	u := e._rep.GetUserExpress(int(r.SellerId))
	list := u.GetAllTemplate()
	var arr []*proto.SExpressTemplate
	for _, v := range list {
		v2 := v.Value()
		if v2.Enabled == 1 {
			arr = append(arr, e.parseExpressTemplateDto(&v2))
		}
	}
	return &proto.ExpressTemplateListResponse{
		Value: arr,
	}, nil
}

// 删除模板
func (e *expressService) DeleteTemplate(_ context.Context, id *proto.ExpressTemplateId) (*proto.Result, error) {
	u := e._rep.GetUserExpress(int(id.SellerId))
	err := u.DeleteTemplate(int(id.TemplateId))
	return e.error(err), nil
}

// 保存地区快递模板
func (e *expressService) SaveAreaTemplate(_ context.Context, r *proto.SaveAreaExpTemplateRequest) (*proto.Result, error) {
	u := e._rep.GetUserExpress(int(r.SellerId))
	t := u.GetTemplate(int(r.TemplateId))
	var err error
	if t == nil {
		err = express.ErrNoSuchTemplate
	} else {
		v := e.parseAreaTemplate(r.Value)
		v.TemplateId = int32(r.TemplateId)
		_, err = t.SaveAreaTemplate(v)
	}
	return e.error(err), nil
}

// 删除模板地区设定
func (e *expressService) DeleteAreaTemplate(_ context.Context, id *proto.AreaTemplateId) (*proto.Result, error) {
	u := e._rep.GetUserExpress(int(id.SellerId))
	t := u.GetTemplate(int(id.TemplateId))
	var err error
	if t == nil {
		err = express.ErrNoSuchTemplate
	} else {
		err = t.DeleteAreaSet(int32(id.AreaTemplateId))
	}
	return e.error(err), nil
}

//// 获取快递费,传入地区编码，根据单位值，如总重量。
//func (e *expressService) GetExpressFee(userId int32,templateId int32,
//	areaCode string, basisUnit float32) float32 {
//	u := e.repo.GetUserExpress(userId)
//	return u.GetExpressFee(templateId, areaCode, basisUnit)
//}

// 根据地区编码获取运费模板
func (e *expressService) GetAreaExpressTemplateByAreaCode(userId int64,
	templateId int32, areaCode string) *express.ExpressAreaTemplate {
	u := e._rep.GetUserExpress(int(userId))
	t := u.GetTemplate(int(templateId))
	if t != nil {
		return t.GetAreaExpressTemplateByAreaCode(areaCode)
	}
	return nil
}

// 根据编号获取地区的运费模板
func (e *expressService) GetAreaExpressTemplate(userId int64,
	templateId int32, id int32) *express.ExpressAreaTemplate {
	u := e._rep.GetUserExpress(int(userId))
	t := u.GetTemplate(int(templateId))
	if t != nil {
		return t.GetAreaExpressTemplate(id)
	}
	return nil
}

// 获取所有的地区快递模板
func (e *expressService) GetAllAreaTemplate(userId int64,
	templateId int32) []express.ExpressAreaTemplate {
	u := e._rep.GetUserExpress(int(userId))
	t := u.GetTemplate(int(templateId))
	if t != nil {
		return t.GetAllAreaTemplate()
	}
	return []express.ExpressAreaTemplate{}
}

func (e *expressService) parseProviderDto(v *express.Provider) *proto.SExpressProvider {
	return &proto.SExpressProvider{
		Id:        int64(v.Id),
		Name:      v.Name,
		Letter:    v.FirstLetter,
		GroupFlag: v.GroupFlag,
		Code:      v.Code,
		ApiCode:   v.ApiCode,
		Enabled:   v.Enabled == 1,
	}
}

func (e *expressService) parseProvider(r *proto.SExpressProvider) *express.Provider {
	return &express.Provider{
		Id:          int32(r.Id),
		Name:        r.Name,
		FirstLetter: r.Letter,
		GroupFlag:   r.GroupFlag,
		Code:        r.Code,
		ApiCode:     r.ApiCode,
		Enabled:     types.ElseInt(r.Enabled, 1, 0),
	}
}

func (e *expressService) parseExpressTemplate(r *proto.SExpressTemplate) *express.ExpressTemplate {
	return &express.ExpressTemplate{
		Id:        int(r.Id),
		VendorId:  int(r.SellerId),
		Name:      r.Name,
		IsFree:    types.ElseInt(r.IsFree, 1, 0),
		Basis:     int(r.Basis),
		FirstUnit: int(r.FirstUnit),
		FirstFee:  r.FirstFee,
		AddUnit:   int(r.AddUnit),
		AddFee:    r.AddFee,
		Enabled:   types.ElseInt(r.Enabled, 1, 0),
	}
}

func (e *expressService) parseExpressTemplateDto(v *express.ExpressTemplate) *proto.SExpressTemplate {
	return &proto.SExpressTemplate{
		Id:        int64(v.Id),
		SellerId:  int64(v.VendorId),
		Name:      v.Name,
		IsFree:    v.IsFree == 1,
		Basis:     int32(v.Basis),
		FirstUnit: int32(v.FirstUnit),
		FirstFee:  v.FirstFee,
		AddUnit:   int32(v.AddUnit),
		AddFee:    v.AddFee,
		Enabled:   v.Enabled == 1,
	}
}

func (e *expressService) parseAreaTemplate(v *proto.SExpressAreaTemplate) *express.ExpressAreaTemplate {
	return &express.ExpressAreaTemplate{
		Id:        int32(v.Id),
		CodeList:  v.CodeList,
		NameList:  v.NameList,
		FirstUnit: v.FirstUnit,
		FirstFee:  v.FirstFee,
		AddUnit:   v.AddUnit,
		AddFee:    v.AddFee,
	}
}
